package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// InteractiveCliInstall installs mods from the command line.
func InteractiveCliInstall() {
	// Lock Mutexes.
	selectedVersionMutex.Lock()
	installFabricOptMutex.Lock()
	minecraftFolderMutex.Lock()
	defer selectedVersionMutex.Unlock()
	defer installFabricOptMutex.Unlock()
	defer minecraftFolderMutex.Unlock()

	println("ibu's mod installer for Fabric 1.14.4+ - CLI")
	println("FAQ: Open the GUI, or use https://mythicmc.org/modpack/faq.html")
	println("")

	// Detect if update is possible.
	updatable := areModsUpdatable()
	if updatable != "" {
		update := takeInput(
			"It looks like you already have "+updatable+" version installed. "+
				"Would you like to update your mods for "+updatable+" to the latest version? [yes/no]",
			&Inputs{"y", "yes", "n", "no"},
		)
		if update == "y" || update == "yes" {
			confirm := takeInput("Confirm? [yes/no]", &Inputs{"y", "yes", "n", "no"})
			if confirm == "y" || confirm == "yes" {
				selectedVersion = updatable
				err := installMods(func(s string) { println(s) }, func(s string) bool {
					response := takeInput(s+" [yes/no]", &Inputs{"y", "yes", "n", "no"})
					if response == "y" || response == "yes" {
						return true
					}
					return false
				})
				if err != nil {
					log.Println(err)
					print("An error occurred. Press Enter on keyboard to close.")
					inputHalt()
					os.Exit(1)
				} else {
					print("Installation succeeded. Press Enter on keyboard to close.")
					inputHalt()
					return
				}
			} else {
				println("Alright, sending you to the installation prompt instead of updating your mods.")
			}
		} else {
			println("Alright, sending you to the installation prompt instead of updating your mods.")
		}
	}

	// Take inputs.
	selectedVersion = getMajorMinecraftVersion(takeInput(
		"Version of Minecraft to use? [1.14.4/1.15.2/1.16.5/1.17.1/1.18.1]",
		&Inputs{"1.14.4", "1.15.2", "1.16.5", "1.17.1", "1.18.1", "1.14", "1.15", "1.16", "1.17", "1.18"},
	))
	println("")
	installFabric := takeInput("Should the modpack install Fabric? [yes/no]", &Inputs{"y", "yes", "n", "no"})
	if installFabric == "y" || installFabric == "yes" {
		installFabricOpt = true
	} else {
		installFabricOpt = false
	}
	println("")
	minecraftFolderYes := takeInput("Do you want to install to custom .minecraft folder? [yes/no]", &Inputs{"y", "yes", "n", "no"})
	if minecraftFolderYes == "y" || minecraftFolderYes == "yes" {
		minecraftFolder = takeInput("Enter path to .minecraft folder:", nil)
	}
	println("")

	// Confirm.
	println("Installing mods for " + selectedVersion + " (Install Fabric: " + strconv.FormatBool(installFabricOpt) + ")")
	confirm := takeInput("Confirm? [yes/no]", &Inputs{"y", "yes", "n", "no"})
	if confirm == "n" || confirm == "no" {
		print("Installation cancelled! Exiting...")
		inputHalt()
		os.Exit(1)
	}

	// Install the mods.
	err := installMods(func(s string) { println(s) }, func(s string) bool {
		response := takeInput(s+" [yes/no]", &Inputs{"y", "yes", "n", "no"})
		if response == "y" || response == "yes" {
			return true
		}
		return false
	})
	if err != nil {
		log.Println(err)
		print("An error occurred. Press Enter on keyboard to close.")
		inputHalt()
		os.Exit(1)
	}
	print("Installation succeeded. Press Enter on keyboard to close.")
	inputHalt()
}

func inputHalt() {
	var input string
	fmt.Scanln(&input)
}

func takeInput(prompt string, inputs *Inputs) string {
	for {
		print(prompt + " ")
		var input string
		fmt.Scanln(&input)
		if inputs != nil && inputs.IsValidInput(input) {
			return input
		} else if inputs == nil {
			return input
		}
		println("Invalid input! Possible values: " + strings.Join(*inputs, ", "))
	}
}

// Inputs is a type that defines a set of valid inputs.
type Inputs []string

// IsValidInput checks if an input in an array matches.
func (inputs *Inputs) IsValidInput(input string) bool {
	for _, val := range *inputs {
		if val == input {
			return true
		}
	}
	return false
}

// Abstraction because yes.
// GetValidInputs gets a list of valid inputs.
// func (inputs *Inputs) GetValidInputs() []string { return *inputs }

// InputChecker is an interface that represents inputs of all types.
// type InputChecker interface {
// 	IsValidInput(input string) bool
// 	GetValidInputs() []string
// }
