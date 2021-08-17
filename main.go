package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const modpackVersion = "1.3.0"

var selectedVersion = "1.17.1"
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
	if minecraftFolder == "" {
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
	if installFabricOpt {
		s := modVersion.Fabric
		if s == "latest" {
			updateProgress("Querying latest Fabric version...")
			s, err = getLatestFabric()
			if err != nil {
				return err
			}
		}
		updateProgress("Downloading Fabric...")
		file, err := downloadFabric(selectedVersion, s)
		if err != nil {
			return err
		}
		updateProgress("Installing Fabric...")
		err = unzipFile(file, filepath.Join(minecraftFolder, "versions"), nil, nil)
		if err != nil {
			return err
		}
	}

	// Check if there's already a mod folder.
	modsFolder := filepath.Join(minecraftFolder, "mods")
	_, err = os.Stat(modsFolder)
	var modsVersionTxt *ModsVersionTxt
	if err == nil {
		modsVersionTxt = getInstalledModsVersion(minecraftFolder)
	}
	incompatModsExist := modsVersionTxt == nil || (modsVersionTxt != nil &&
		modsVersionTxt.Version != getMajorMinecraftVersion(selectedVersion))
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err != nil && os.IsNotExist(err) {
		updateProgress("Creating mods folder...")
		if err = os.MkdirAll(modsFolder, os.ModePerm); err != nil {
			return err
		}
	} else if err == nil && incompatModsExist {
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
	}

	updateProgress("Downloading my mods for " + selectedVersion + "...")
	file, err := downloadFile(modVersion.URL)
	if err != nil {
		return err
	}

	// Install/update the mods.
	updateProgress("Installing mods...")
	modsversionTxt := selectedVersion + "\n"
	// TODO: Improve error handling by reverting changes made?
	if incompatModsExist {
		modsData, err := readModsJsonFromZip(file)
		if err != nil {
			return err
		}
		err = unzipFile(file, modsFolder, []string{"mods.json"}, nil)
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
	} else {
		modsData, err := readModsJsonFromZip(file)
		if err != nil {
			return err
		}
		var modsToInstall []string
		// Compare modsVersionTxt with mods.json to get a list of new mods.
		for modName, modFilename := range modsData.Mods {
			found := false
			for _, installedMod := range modsVersionTxt.InstalledMods {
				if installedMod == modName {
					found = true
				}
			}
			if !found {
				modsToInstall = append(modsToInstall, modFilename)
			}
		}
		// Discover old mods which need to be moved.
		err = os.MkdirAll(filepath.Join(modsFolder, "oldmods"), os.ModePerm)
		if err != nil {
			return err
		}
		for modFilename, modName := range modsData.OldMods {
			modFilePath := filepath.Join(modsFolder, modFilename)
			stat, err := os.Stat(modFilePath)
			found := err == nil && !stat.IsDir()
			if found {
				err := os.Rename(modFilePath, filepath.Join(modsFolder, "oldmods", modFilename))
				if err != nil {
					return err
				}
				if _, modExists := modsData.Mods[modName]; modExists {
					modsToInstall = append(modsToInstall, modsData.Mods[modName])
				}
			}
		}
		// Unzip only new mods.
		err = unzipFile(file, modsFolder, nil, modsToInstall)
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
		filepath.Join(modsFolder, "modsversion.txt"), []byte(modsversionTxt), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// Lock minecraftFolder before calling.
func areModsUpdatable() bool {
	folder := minecraftFolder
	if folder == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return false
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
	var modsVersionTxt *ModsVersionTxt
	if err == nil {
		modsVersionTxt = getInstalledModsVersion(folder)
	}
	return modsVersionTxt != nil && modsVersionTxt.Version == getMajorMinecraftVersion(selectedVersion)
}

func getInstalledModsVersion(location string) *ModsVersionTxt {
	file, err := os.Open(filepath.Join(location, "mods", "modsversion.txt"))
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
	return &ModsVersionTxt{
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

func unzipFile(zipFile []byte, location string, exclude []string, include []string) error {
	// Uses: os, io, strings, filepath, zip, bytes
	r, err := zip.NewReader(bytes.NewReader(zipFile), int64(len(zipFile)))
	if err != nil {
		return err
	}
	for _, f := range r.File {
		toContinue := len(include) > 0
		for _, excluded := range exclude {
			if excluded == f.FileInfo().Name() {
				toContinue = true
			}
		}
		for _, included := range include {
			if included == f.FileInfo().Name() {
				toContinue = false
			}
		}
		if toContinue {
			continue
		}
		fpath := filepath.Join(location, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(location)+string(os.PathSeparator)) {
			continue // "%s: illegal file path"
		}
		// Create folders.
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		// Create parent folder of file if needed.
		err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
		if err != nil {
			return err
		}
		// Open target file.
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		// Open file in zip.
		rc, err := f.Open()
		if err != nil {
			return err
		}
		// Copy file from zip to disk.
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
		outFile.Close()
		rc.Close()
	}
	return nil
}

// ModsData is a JSON containing data on mods inside a zip.
type ModsData struct {
	Mods    map[string]string `json:"mods"`
	OldMods map[string]string `json:"oldmods"`
}

// ModsVersionTxt contains the contents of modsversion.txt.
type ModsVersionTxt struct {
	Version       string
	InstalledMods []string
}
