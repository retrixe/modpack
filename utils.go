package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Lock minecraftFolder before calling.
func getUpdatableVersions() []string {
	folder := minecraftFolder
	if folder == "" || folder == ".minecraft" {
		home, err := os.UserHomeDir()
		if err != nil {
			return []string{}
		}
		if runtime.GOOS == "darwin" {
			folder = filepath.Join(home, "Library", "Application Support", "minecraft")
		} else if runtime.GOOS == "windows" {
			folder = filepath.Join(home, "AppData", "Roaming", ".minecraft")
		} else {
			folder = filepath.Join(home, ".minecraft")
		}
	}
	contents, err := os.ReadDir(filepath.Join(folder, "mods"))
	var modsVersionTxt *InstalledModsInfo
	if err == nil {
		modsVersionTxt = getInstalledModsVersion(filepath.Join(folder, "mods"))
	}
	if modsVersionTxt != nil {
		return []string{getMajorMinecraftVersion(modsVersionTxt.Version)}
	} else {
		// Check all subfolders for available updates.
		versions := []string{}
		for _, file := range contents {
			if strings.HasPrefix(file.Name(), "=") && file.IsDir() {
				modsInfo := getInstalledModsVersion(filepath.Join(folder, "mods", file.Name()))
				if modsInfo != nil &&
					!includes(versions, getMajorMinecraftVersion(modsInfo.Version)) &&
					// We don't support upgrading mods in sub-folders.
					!includes(fabricVersions, getMajorMinecraftVersion(modsInfo.Version)) {
					versions = append(versions, getMajorMinecraftVersion(modsInfo.Version))
				}
			}
		}
		return versions
	}
}

func includes[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func getInstalledModsVersion(location string) *InstalledModsInfo {
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
