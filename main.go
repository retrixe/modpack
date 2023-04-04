//go:build !launcher

package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const modpackVersion = "1.6.0"
const defaultVersion = "1.19"

// LOW-TODO: Create a special profile that loads mods from a special folder?

var selectedVersion = defaultVersion
var selectedVersionMutex sync.Mutex
var installFabricOpt = true
var installFabricOptMutex sync.Mutex
var minecraftFolder = ""
var minecraftFolderMutex sync.Mutex

func main() {
	if len(os.Args) >= 2 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		println("modpack version " + modpackVersion)
		return
	} else if len(os.Args) >= 2 && os.Args[1] == "install" {
		InteractiveCliInstall()
		return
	}
	runGui()
}

func installMods(updateProgress func(string), queryUser func(string) bool) error {
	if minecraftFolder == "" || minecraftFolder == ".minecraft" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		if runtime.GOOS == "darwin" {
			minecraftFolder = filepath.Join(home, "Library", "Application Support", "minecraft")
		} else if runtime.GOOS == "windows" {
			minecraftFolder = filepath.Join(home, "AppData", "Roaming", ".minecraft")
		} else {
			minecraftFolder = filepath.Join(home, ".minecraft")
		}
	}
	updateProgress("Querying latest mod versions...")
	modVersion, err := getModVersions(selectedVersion)
	if err != nil {
		return err
	}
	useQuilt := strings.HasPrefix(modVersion.Fabric, "quilt:")
	version := modVersion.FullVersion
	if version == "" { // some old compatibility if
		version = selectedVersion
	}
	if installFabricOpt {
		s := modVersion.Fabric
		loaderName := "Fabric"
		if useQuilt {
			loaderName = "Quilt"
			s = s[6:]
		}
		if s == "latest" {
			updateProgress("Querying latest " + loaderName + " version...")
			s, err = getLatestFabric(loaderName == "Quilt")
			if err != nil {
				return err
			}
		}
		updateProgress("Downloading " + loaderName + "...")
		if loaderName == "Quilt" {
			file, err := downloadQuilt(version, s)
			if err != nil {
				return err
			}
			updateProgress("Installing Quilt...")
			versionName := "quilt-loader-" + s + "-" + version
			versionFolder := filepath.Join(minecraftFolder, "versions", versionName)
			err = os.MkdirAll(versionFolder, os.ModePerm)
			if err != nil {
				return err
			}
			err = os.WriteFile(filepath.Join(versionFolder, versionName+".json"), file, os.ModePerm)
			if err != nil {
				return err
			}

		} else {
			file, err := downloadFabric(version, s)
			if err != nil {
				return err
			}
			updateProgress("Installing Fabric...")
			err = unzipFile(file, filepath.Join(minecraftFolder, "versions"), nil, nil)
			if err != nil {
				return err
			}
		}
	}

	// Check if there's already a mod folder.
	modsFolder := filepath.Join(minecraftFolder, "mods")
	modsFolderContents, err := os.ReadDir(modsFolder)
	var installedModsInfo *InstalledModsInfo
	if err == nil {
		installedModsInfo = getInstalledModsVersion(modsFolder)
	}
	upgradeSupported := installedModsInfo != nil &&
		getMajorMinecraftVersion(installedModsInfo.Version) == selectedVersion
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err != nil && os.IsNotExist(err) {
		updateProgress("Creating mods folder...")
		if err = os.MkdirAll(modsFolder, os.ModePerm); err != nil {
			return err
		}
	} else if err == nil && hasAnyJarFile(modsFolderContents) && !upgradeSupported {
		_, err = os.Stat(filepath.Join(minecraftFolder, "oldmods"))
		if err == nil || !os.IsNotExist(err) {
			return errors.New("mods folder and oldmods folder exist, user must remove/rename either folder")
		}
		answer := queryUser(`A mods folder is already present which does not seem to be created by this pack.
Would you like to rename it to oldmods?`)
		if !answer {
			return errors.New("mods folder already exists, and user refused to rename it")
		}
		os.Rename(modsFolder, filepath.Join(minecraftFolder, "oldmods"))
	} else if useQuilt && installedModsInfo != nil {
		updateProgress("Moving mods into Quilt version-specific folder...")
		// These mods match the currently selected version of modpack, migrate them into a subfolder.
		tempFolder := filepath.Join(minecraftFolder, "~mods")
		os.Rename(modsFolder, tempFolder)
		os.Mkdir(filepath.Join(minecraftFolder, "mods"), os.ModePerm)
		os.Rename(tempFolder, filepath.Join(minecraftFolder, "mods", "="+version))
	}

	// Check mods version-specific subfolder.
	if useQuilt {
		upgradeSupported = true
		existingModsSubfolder := ""
		for _, file := range modsFolderContents {
			if strings.HasPrefix(file.Name(), "="+version) &&
				existingModsSubfolder < file.Name() {
				existingModsSubfolder = file.Name()
				break
			}
		}
		if existingModsSubfolder != "" {
			// Check if this subfolder is managed.
			installedModsInfo = getInstalledModsVersion(filepath.Join(modsFolder, existingModsSubfolder))
			// If it is, has a compatible minor version, but a different patch version, then copy mods.
			if installedModsInfo != nil &&
				selectedVersion == getMajorMinecraftVersion(installedModsInfo.Version) &&
				existingModsSubfolder != "="+version {
				// Make new mods folder.
				oldModsFolder := filepath.Join(modsFolder, existingModsSubfolder)
				newModsFolder := filepath.Join(modsFolder, "="+version)
				os.MkdirAll(newModsFolder, os.ModePerm)
				// Copy old mods from the old folder to the new one.
				oldModsFolderContents, err := os.ReadDir(oldModsFolder)
				if err != nil {
					return err
				}
				for _, file := range oldModsFolderContents {
					if !file.IsDir() {
						input, err := os.ReadFile(filepath.Join(oldModsFolder, file.Name()))
						if err != nil {
							return err
						}
						err = os.WriteFile(filepath.Join(newModsFolder, file.Name()), input, os.ModePerm)
						if err != nil {
							return err
						}
					}
				}
			} else if existingModsSubfolder == "="+version && installedModsInfo == nil {
				// Rename the folder if it matches the one we need.
				neededName := "mods/=" + version
				newName := "mods/.old " + existingModsSubfolder
				_, err = os.Stat(filepath.Join(minecraftFolder, "mods", ".old "+existingModsSubfolder))
				if err == nil || !os.IsNotExist(err) {
					return errors.New(neededName + " folder and " + newName + " folder exist, user must remove/rename either folder")
				}
				answer := queryUser(`A ` + neededName + ` folder is already present which does not seem to be created by this pack.
Would you like to rename it to ` + newName + `?`)
				if !answer {
					return errors.New(neededName + " folder already exists, and user refused to rename it")
				}
				os.Rename(filepath.Join(modsFolder, existingModsSubfolder),
					filepath.Join(modsFolder, ".old "+existingModsSubfolder))
				upgradeSupported = false
			}
		} else {
			os.MkdirAll(filepath.Join(modsFolder, "="+version), os.ModePerm)
			upgradeSupported = false
		}
	}

	updateProgress("Downloading my mods for " + version + "...")
	file, err := downloadFile(modVersion.URL)
	if err != nil {
		return err
	}

	// Install/update the mods.
	updateProgress("Installing mods...")
	modFolder := modsFolder
	if useQuilt {
		modFolder = filepath.Join(modFolder, "="+version)
	}
	modsversionTxt := selectedVersion + "\n"
	if !upgradeSupported { // The mods folder no longer exists.
		modsData, err := readModsJsonFromZip(file)
		if err != nil {
			return err
		}
		err = unzipFile(file, modFolder, []string{"mods.json"}, nil)
		if err != nil {
			return err
		} else if modsData != nil {
			// Get all the mods that were installed and put them in modsversion.txt
			mods := make([]string, 0, len(modsData.Mods))
			for mod := range modsData.Mods {
				mods = append(mods, mod)
			}
			modsversionTxt = selectedVersion + "\n" + strings.Join(mods, ",") + "\n"
		}
	} else { // Use the upgrade mechanism.
		modsData, err := readModsJsonFromZip(file)
		if err != nil {
			return err
		}
		var modsToInstall []string = []string{}
		// Compare modsVersionTxt with mods.json to get a list of new mods.
		for modName, modFilename := range modsData.Mods {
			found := false
			for _, installedMod := range installedModsInfo.InstalledMods {
				if installedMod == modName {
					found = true
				}
			}
			if !found {
				modsToInstall = append(modsToInstall, modFilename)
			}
		}
		// Discover old mods which need to be moved.
		err = os.MkdirAll(filepath.Join(modFolder, "oldmods"), os.ModePerm)
		if err != nil {
			return err
		}
		for modFilename, modName := range modsData.OldMods {
			modFilePath := filepath.Join(modFolder, modFilename)
			stat, err := os.Stat(modFilePath)
			found := err == nil && !stat.IsDir()
			if found {
				err := os.Rename(modFilePath, filepath.Join(modFolder, "oldmods", modFilename))
				if err != nil {
					return err
				}
				if _, modExists := modsData.Mods[modName]; modExists {
					modsToInstall = append(modsToInstall, modsData.Mods[modName])
				}
			}
		}
		// Unzip only new mods.
		err = unzipFile(file, modFolder, nil, modsToInstall)
		if err != nil {
			return err
		} else if modsData != nil {
			// Get all the mods that were installed and put them in modsversion.txt
			mods := make([]string, 0, len(modsData.Mods))
			for mod := range modsData.Mods {
				mods = append(mods, mod)
			}
			modsversionTxt = selectedVersion + "\n" + strings.Join(mods, ",") + "\n"
		}
	}
	err = os.WriteFile( // Write the modsversion.txt.
		filepath.Join(modFolder, "modsversion.txt"), []byte(modsversionTxt), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// Lock minecraftFolder before calling.
func areModsUpdatable() string { // TODO: Support Quilt subfolders
	folder := minecraftFolder
	if folder == "" || folder == ".minecraft" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		if runtime.GOOS == "darwin" {
			folder = filepath.Join(home, "Library", "Application Support", "minecraft")
		} else if runtime.GOOS == "windows" {
			folder = filepath.Join(home, "AppData", "Roaming", ".minecraft")
		} else {
			folder = filepath.Join(home, ".minecraft")
		}
	}
	_, err := os.Stat(filepath.Join(folder, "mods"))
	var modsVersionTxt *InstalledModsInfo
	if err == nil {
		modsVersionTxt = getInstalledModsVersion(filepath.Join(folder, "mods"))
	}
	if modsVersionTxt != nil {
		return modsVersionTxt.Version
	} else {
		return ""
	}
}

func getInstalledModsVersion(location string) *InstalledModsInfo { // TODO: Support Quilt subfolders
	file, err := os.Open(filepath.Join(location, "modsversion.txt"))
	if err != nil {
		return nil
	}
	defer file.Close()
	contents, err := io.ReadAll(file)
	if err != nil {
		return nil
	}
	// Extract only the major version for now, comprehensive update system later.
	stringContents := string(contents)
	installedMods := make([]string, 0)
	firstNewlineIndex := strings.Index(stringContents, "\n")
	firstLine := stringContents
	if firstNewlineIndex != -1 {
		firstLine = stringContents[:firstNewlineIndex]
		nextNewlineIndex := firstNewlineIndex + strings.Index(stringContents[firstNewlineIndex+1:], "\n")
		if nextNewlineIndex != -1 {
			installedMods = strings.Split(stringContents[firstNewlineIndex+1:nextNewlineIndex+1], ",")
		}
	}
	return &InstalledModsInfo{
		Version:       getMajorMinecraftVersion(firstLine),
		InstalledMods: installedMods,
	}
}

func getMajorMinecraftVersion(version string) string {
	lastIndex := strings.LastIndex(version, ".")
	if lastIndex == -1 || strings.Index(version, ".") == lastIndex {
		return version
	}
	return version[:lastIndex]
}

func hasAnyJarFile(files []fs.DirEntry) bool {
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".jar") {
			return true
		}
	}
	return false
}

func readModsJsonFromZip(zipFile []byte) (*ModsData, error) {
	// Uses zip and bytes
	r, err := zip.NewReader(bytes.NewReader(zipFile), int64(len(zipFile)))
	if err != nil {
		return nil, err
	}
	var modsData *ModsData = nil
	for _, f := range r.File {
		if f.FileInfo().Name() == "mods.json" {
			modsJSON, err := f.Open()
			if err != nil {
				return nil, err
			}
			var decode ModsData
			err = json.NewDecoder(modsJSON).Decode(&decode)
			if err != nil {
				return nil, err
			}
			modsData = &decode
			break
		}
	}
	return modsData, nil
}

// ModsData is a JSON containing data on mods inside a zip.
type ModsData struct {
	Mods    map[string]string `json:"mods"`
	OldMods map[string]string `json:"oldmods"`
}

// InstalledModsInfo contains the contents of mods/modsversion.txt or mods/=<version>/modpack.txt
type InstalledModsInfo struct {
	Version       string
	InstalledMods []string
}
