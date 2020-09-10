package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// InteractiveCliInstall ... Installs mods from the command line.
func InteractiveCliInstall() {
	// Lock Mutexes.
	selectedVersionMutex.Lock()
	installFabricOptMutex.Lock()
	defer selectedVersionMutex.Unlock()
	defer installFabricOptMutex.Unlock()

	fmt.Println("ibu's mod installer for Fabric 1.14.4+ - CLI")
	fmt.Println("FAQ: Open the GUI, or use https://mythicmc.info/modpack/faq.html")
	fmt.Println("")

	// Take inputs.
	selectedVersion = takeInput("Version of Minecraft to use? [1.14.4/1.15.2/1.16.2]", &Inputs{"1.14.4", "1.15.2", "1.16.2"})
	fmt.Println("")
	installFabric := takeInput("Should the modpack install Fabric? [yes/no]", &Inputs{"y", "yes", "n", "no"})
	if installFabric == "y" || installFabric == "yes" {
		installFabricOpt = true
	} else {
		installFabricOpt = false
	}
	fmt.Println("")

	// Confirm.
	fmt.Println("Installing mods for " + selectedVersion + " (Install Fabric: " + strconv.FormatBool(installFabricOpt) + ")")
	confirm := takeInput("Confirm? [yes/no]", &Inputs{"y", "yes", "n", "no"})
	if confirm == "n" || confirm == "no" {
		fmt.Println("Installation cancelled! Exiting...")
		return
	}

	// Install the mods.
	home, err := os.UserHomeDir()
	if err != nil {
		log.Panicln(err)
	}
	minecraftFolder := filepath.Join(home, ".minecraft") // TODO advanced opts
	if runtime.GOOS == "darwin" {
		minecraftFolder = filepath.Join(home, "Library", "Application Support", "minecraft")
	} else if runtime.GOOS == "windows" {
		minecraftFolder = filepath.Join(home, "AppData", "Roaming", ".minecraft")
	}
	fmt.Println("Querying latest mod versions...")
	modVersion, err := getModVersions(selectedVersion)
	if err != nil {
		log.Panicln(err)
	}
	if installFabricOpt {
		s := modVersion.Fabric
		if s == "latest" {
			fmt.Println("Querying latest Fabric version...")
			s, err = getLatestFabric()
			if err != nil {
				log.Panicln(err)
			}
		}
		fmt.Println("Downloading Fabric...")
		file, err := downloadFabric(selectedVersion, s)
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println("Installing Fabric...")
		err = unzipFile(file, filepath.Join(minecraftFolder, "versions"))
		if err != nil {
			log.Panicln(err)
		}
	}
	fmt.Println("Downloading my mods for " + selectedVersion + "...")
	file, err := downloadMods(modVersion.URL)
	if err != nil {
		log.Panicln(err)
	}
	// Check if there's already a mod folder.
	_, err = os.Stat(filepath.Join(minecraftFolder, "mods"))
	modsExist := false // Check if the mod folder contains the same version of mods as our pack.
	if err == nil {
		modsExist = getInstalledModsVersion(minecraftFolder) == getMajorMinecraftVersion(selectedVersion)
	}
	if err != nil && !os.IsNotExist(err) {
		log.Panicln(err)
	} else if err == nil && !modsExist {
		os.Rename(filepath.Join(minecraftFolder, "mods"), filepath.Join(minecraftFolder, "oldmodfolder"))
		fmt.Println("Renamed old mods folder to oldmodfolder!")
	} else if err != nil && os.IsNotExist(err) {
		fmt.Println("Creating mods folder...")
		if err = os.MkdirAll(filepath.Join(minecraftFolder, "mods"), os.ModePerm); err != nil {
			log.Panicln(err)
		}
	}
	fmt.Println("Unzipping my mods...")
	if modsExist {
		r, err := zip.NewReader(bytes.NewReader(file), int64(len(file)))
		if err != nil {
			log.Panicln(err)
		}
		// Read mods.json.
		var modsData ModsData
		for _, f := range r.File {
			if filepath.Base(f.Name) == "mods.json" {
				modsJSON, err := f.Open()
				if err != nil {
					log.Panicln(err)
				}
				json.NewDecoder(modsJSON).Decode(&modsData)
				break
			}
		}
		if &modsData != nil {
			if err = moveOldMods(modsData, minecraftFolder, r); err != nil {
				log.Panicln(err)
			}
		}
	} else {
		err = unzipFile(file, filepath.Join(minecraftFolder, "mods"))
		if err != nil {
			log.Panicln(err)
		}
	}
	err = ioutil.WriteFile( // Write the modsversion.txt.
		filepath.Join(minecraftFolder, "mods", "modsversion.txt"), []byte(selectedVersion), os.ModePerm)
	if err != nil {
		log.Panicln(err)
	}
}

func takeInput(prompt string, inputs *Inputs) string {
	for {
		fmt.Print(prompt + " ")
		var input string
		fmt.Scanln(&input)
		if inputs.IsValidInput(input) {
			return input
		}
		fmt.Println("Invalid input! Possible values: " + strings.Join(*inputs, ", "))
	}
}

// Inputs ... A type that defines a set of valid inputs.
type Inputs []string

// IsValidInput ... Checks if an input in an array matches.
func (inputs *Inputs) IsValidInput(input string) bool {
	for _, val := range *inputs {
		if val == input {
			return true
		}
	}
	return false
}

// Abstraction because yes.
// GetValidInputs ... Get a list of valid inputs.
// func (inputs *Inputs) GetValidInputs() []string { return *inputs }

// InputChecker ... An interface that represents inputs of all types.
// type InputChecker interface {
// 	IsValidInput(input string) bool
// 	GetValidInputs() []string
// }
