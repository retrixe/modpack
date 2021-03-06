package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const modpackVersion = "1.2.1"

var selectedVersion = "1.16.5"
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
		_, err = unzipFile(file, filepath.Join(minecraftFolder, "versions"))
		if err != nil {
			return err
		}
	}
	updateProgress("Downloading my mods for " + selectedVersion + "...")
	file, err := downloadFile(modVersion.URL)
	if err != nil {
		return err
	}

	// Check if there's already a mod folder.
	_, err = os.Stat(filepath.Join(minecraftFolder, "mods"))
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
		if err = os.MkdirAll(filepath.Join(minecraftFolder, "mods"), os.ModePerm); err != nil {
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
		os.Rename(filepath.Join(minecraftFolder, "mods"), filepath.Join(minecraftFolder, "oldmods"))
	}

	// Install/update the mods.
	updateProgress("Installing mods...")
	modsversionTxt := selectedVersion + "\n"
	if !incompatModsExist {
		r, err := zip.NewReader(bytes.NewReader(file), int64(len(file)))
		if err != nil {
			return err
		}
		// Read mods.json.
		var modsData ModsData
		for _, f := range r.File {
			if filepath.Base(f.Name) == "mods.json" {
				modsJSON, err := f.Open()
				if err != nil {
					return err
				}
				json.NewDecoder(modsJSON).Decode(&modsData)
				break
			}
		}
		if modsData.Mods != nil {
			if err = moveOldMods(modsData, minecraftFolder, r, modsVersionTxt); err != nil {
				return err
			}
			// Get all the mods that were installed and put them in modsversion.txt
			mods := make([]string, 0, len(modsData.Mods))
			for mod := range modsData.Mods {
				mods = append(mods, mod)
			}
			modsversionTxt = selectedVersion + "\n" + strings.Join(mods, ",") + "\n"
		}
	} else {
		modsData, err := unzipFile(file, filepath.Join(minecraftFolder, "mods"))
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
	err = ioutil.WriteFile( // Write the modsversion.txt.
		filepath.Join(minecraftFolder, "mods", "modsversion.txt"), []byte(modsversionTxt), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func getInstalledModsVersion(location string) *ModsVersionTxt {
	file, err := os.Open(filepath.Join(location, "mods", "modsversion.txt"))
	if err != nil {
		return nil
	}
	defer file.Close()
	contents, err := ioutil.ReadAll(file)
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

func moveOldMods(modsData ModsData, minecraftFolder string, r *zip.Reader, m *ModsVersionTxt) error {
	location := filepath.Join(minecraftFolder, "mods")
	err := os.MkdirAll(filepath.Join(location, "oldmods"), os.ModePerm)
	if err != nil {
		return err
	}
	for key, val := range modsData.OldMods {
		if _, err := os.Stat(filepath.Join(location, key)); err == nil {
			err := os.Rename(filepath.Join(location, key), filepath.Join(location, "oldmods", key))
			if err != nil {
				return err
			}
			mod := modsData.Mods[val]
			for _, f := range r.File {
				if filepath.Base(f.Name) == mod {
					modFile, err := f.Open()
					if err != nil {
						return err
					}
					fpath := filepath.Join(location, mod)
					outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
					if err != nil {
						return err
					}
					io.Copy(outFile, modFile)
					break
				}
			}
		}
	}
	// Install mods newly added to the pack.
	if m != nil && len(m.InstalledMods) > 0 {
		for mod, filename := range modsData.Mods {
			// Check if it's in InstalledMods, if it isn't, then install it.
			installed := false
			for _, installedMod := range m.InstalledMods {
				print(installedMod)
				if installedMod == mod {
					installed = true
				}
			}
			if !installed { // Then install it.
				for _, f := range r.File {
					if filepath.Base(f.Name) == filename {
						modFile, err := f.Open()
						if err != nil {
							return err
						}
						fpath := filepath.Join(location, filename)
						outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
						if err != nil {
							return err
						}
						io.Copy(outFile, modFile)
						break
					}
				}
			}
		}
	}

	return nil
}

func unzipFile(zipFile []byte, location string) (*ModsData, error) {
	// Uses: os, io, strings, filepath, zip, bytes
	r, err := zip.NewReader(bytes.NewReader(zipFile), int64(len(zipFile)))
	if err != nil {
		return nil, err
	}
	var modsData *ModsData = nil
	for _, f := range r.File {
		if f.Name == "mods.json" { // Ignore /mods.json during extraction.
			modsJSON, err := f.Open()
			if err != nil {
				return nil, err
			}
			var decode ModsData
			json.NewDecoder(modsJSON).Decode(&decode)
			modsData = &decode
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
			return nil, err
		}
		// Open target file.
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return nil, err
		}
		// Open file in zip.
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		// Copy file from zip to disk.
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return nil, err
		}
		outFile.Close()
		rc.Close()
	}
	return modsData, nil
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
