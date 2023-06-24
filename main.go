//go:build !launcher

package main

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const modpackVersion = "1.7.0"
const defaultVersion = "1.20"

var selectedVersion = defaultVersion
var selectedVersionMutex sync.Mutex
var installFabricOpt = true
var installFabricOptMutex sync.Mutex
var minecraftFolder = ""
var minecraftFolderMutex sync.Mutex

var fabricVersions = []string{"1.14", "1.15", "1.16", "1.17"}

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

// ModsData is a JSON containing data on mods inside a zip.
type ModsData struct {
	Mods    map[string]string `json:"mods"`
	OldMods map[string]string `json:"oldmods"`
}

// InstalledModsInfo contains the contents of mods/modsversion.txt or mods/=<version>/modsversion.txt
type InstalledModsInfo struct {
	Version       string
	InstalledMods []string
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
			// Create launcher profile.
			err = addProfileToLauncher(minecraftFolder, versionName, selectedVersion)
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
		upgradeSupported = false
		// First check if there is a compatible managed sub-folder.
		existingModsSubfolder := ""
		for _, file := range modsFolderContents {
			if strings.HasPrefix(file.Name(), "="+selectedVersion) &&
				existingModsSubfolder < file.Name() { // We need the newest managed one.
				modsInfo := getInstalledModsVersion(filepath.Join(modsFolder, file.Name()))
				if modsInfo != nil && getMajorMinecraftVersion(modsInfo.Version) == selectedVersion {
					upgradeSupported = true
					installedModsInfo = modsInfo
					existingModsSubfolder = file.Name()
				}
			}
		}
		// Check if a mod folder for the current version already exists.
		modFolderForCurrentVersionAlreadyExists := false
		_, err = os.Stat(filepath.Join(modsFolder, "="+version))
		if err == nil || !os.IsNotExist(err) {
			modFolderForCurrentVersionAlreadyExists = true
		}
		// If a supported mod folder exists, and it doesn't match the current version, OR
		// there is no supported mod folder, then rename any existing folder named =version.
		if (upgradeSupported && existingModsSubfolder != "="+version) || !upgradeSupported {
			if modFolderForCurrentVersionAlreadyExists {
				_, err = os.Stat(filepath.Join(minecraftFolder, "mods", ".old ="+version))
				if err == nil || !os.IsNotExist(err) {
					return errors.New("mods/=" + version + " folder and mods/.old =" + version + " folder exist, user must remove/rename either folder")
				}
				answer := queryUser(`A mods/=` + version + ` folder is already present which does not seem to be created by this pack.
Would you like to rename it to mods/.old =` + version + `?`)
				if !answer {
					return errors.New("mods/=" + version + " folder already exists, and user refused to rename it")
				}
				err := os.Rename(filepath.Join(modsFolder, "="+version),
					filepath.Join(modsFolder, ".old ="+version))
				if err != nil {
					return err
				}
			}
			// If a supported mod folder exists, then rename it to =version.
			if upgradeSupported {
				os.Rename(filepath.Join(modsFolder, existingModsSubfolder),
					filepath.Join(modsFolder, "="+version))
			}
		}
		err := os.MkdirAll(filepath.Join(modsFolder, "="+version), os.ModePerm)
		if err != nil {
			return err
		}
	}

	updateProgress("Downloading mods for " + version + "...")
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
		quiltIgnoreFile, err := os.Create(filepath.Join(modFolder, "oldmods", "quilt_loader_ignored"))
		if err != nil {
			return err
		}
		quiltIgnoreFile.Close()
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
