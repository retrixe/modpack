//go:build !clionly

package main

import (
	_ "embed"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/sqweek/dialog"
	"github.com/webview/webview"
)

var w webview.WebView

// LOW-TODO: Bundle Roboto font, don't depend on the internet for this?
const html = `
<html lang="en">
<head>
  <meta charset="UTF-8">
  <!-- Use minimum-scale=1 to enable GPU rasterization -->
  <meta
    name='viewport'
    content='user-scalable=0, initial-scale=1, minimum-scale=1, width=device-width, height=device-height'
  />
	<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Roboto:wght@200;300;400;500;700;900&display=swap">
	<style>
	body {
		margin: 0;
		font-family: -apple-system,BlinkMacSystemFont,"Segoe UI",
		  Ubuntu,Cantarell,Oxygen-Sans,"Helvetica Neue",Arial,Roboto,sans-serif;
	}
  </style>
</head>
<body><div id="app"></div><script>initiateReact()</script></body>
</html>
`

//go:embed ui/dist/main.js
var js string

var guiDialogQueryResponse bool
var guiDialogQueryResponseMutex sync.Mutex

func runGui() {
	debug := false
	if val, exists := os.LookupEnv("DEBUG"); exists {
		debug = val == "true"
	}
	w = webview.New(debug)
	defer w.Destroy()
	w.SetSize(640, 480, webview.HintNone) // 540/360
	w.SetTitle("ibu's mod installer " + modpackVersion)
	// Bind a function to initiate React via webview.Eval.
	w.Bind("initiateReact", func() { w.Eval(js) })
	w.Bind("changeVersion", func(name string) {
		selectedVersionMutex.Lock()
		minecraftFolderMutex.Lock()
		defer selectedVersionMutex.Unlock()
		defer minecraftFolderMutex.Unlock()
		selectedVersion = name
	})
	w.Bind("toggleInstallFabric", func() {
		installFabricOptMutex.Lock()
		defer installFabricOptMutex.Unlock()
		installFabricOpt = !installFabricOpt
	})
	w.Bind("respondQuery", func(response bool) {
		guiDialogQueryResponse = response
		guiDialogQueryResponseMutex.Unlock()
	})
	w.Bind("updateMinecraftFolder", func(directory string) {
		selectedVersionMutex.Lock()
		minecraftFolderMutex.Lock()
		defer selectedVersionMutex.Unlock()
		defer minecraftFolderMutex.Unlock()
		minecraftFolder = directory
		checkUpdatableAndUpdateVersion()
	})
	w.Bind("promptForFolder", func() {
		directory, err := dialog.Directory().Title("Select Minecraft game directory").Browse()
		if err != nil {
			setError(err.Error())
			return
		}
		selectedVersionMutex.Lock()
		minecraftFolderMutex.Lock()
		defer selectedVersionMutex.Unlock()
		defer minecraftFolderMutex.Unlock()
		minecraftFolder = directory
		checkUpdatableAndUpdateVersion()
		folder := strings.ReplaceAll(strings.ReplaceAll(directory, "\\", "\\\\"), "\"", "\\\"")
		w.Eval("window.setMinecraftFolderState(\"" + folder + "\")")
	})
	w.Bind("installMods", func() { go initiateInstall() })
	w.Navigate("data:text/html," + strings.ReplaceAll(string(html), "+", "%2B"))
	w.Run()
}

func initiateInstall() {
	selectedVersionMutex.Lock()
	installFabricOptMutex.Lock()
	minecraftFolderMutex.Lock()
	defer selectedVersionMutex.Unlock()
	defer installFabricOptMutex.Unlock()
	defer minecraftFolderMutex.Unlock()
	defer w.Dispatch(checkUpdatableAndUpdateVersion)
	w.Dispatch(func() {
		setError("")             // This is already set in the JavaScript on mount, but whatever.
		setMessage("Working...") // This is already set in the JavaScript on mount, but whatever.
		setInProgress(true)
	})
	// TODO: If there was an upgrade, then it should show what mods were updated and what mods weren't
	err := installMods(func(msg string) { w.Dispatch(func() { setMessage(msg) }) }, queryUser)
	w.Dispatch(func() {
		if err != nil && err.Error() != "Cancelled" {
			log.Println(err)
			setError(err.Error())
		} else if installFabricOpt {
			setMessage("Done! You can now launch Minecraft, select the latest \"ibu\\'s modpack\" profile," +
				" and enjoy! See the FAQ if you need any more information.")
		} else {
			setMessage("Done! You can now launch Minecraft, select the latest fabric-loader or quilt-loader" +
				" version, and enjoy! See the FAQ if you need any more information.")
		}
		setInProgress(false)
	})
}

func queryUser(query string) bool {
	guiDialogQueryResponseMutex.Lock()
	w.Dispatch(func() { w.Eval("window.setQueryState(`" + query + "`)") })
	// This waits for the mutex to unlock.
	guiDialogQueryResponseMutex.Lock()
	defer guiDialogQueryResponseMutex.Unlock()
	return guiDialogQueryResponse
}

func checkUpdatableAndUpdateVersion() {
	updatable := getUpdatableVersions()
	versionArrayJson := ""
	for _, version := range updatable {
		versionArrayJson += "\"" + version + "\","
	}
	w.Eval("setUpdatableVersionsState([" + versionArrayJson + "])")
	w.Eval("setMinecraftVersionState(\"" + selectedVersion + "\")")
}

func setMessage(content string) {
	w.Dispatch(func() { w.Eval("window.setMessageState('" + content + "')") })
}

func setError(content string) {
	w.Dispatch(func() { w.Eval("window.setErrorState('" + content + "')") })
}

// setInProgress disables buttons and shows the progress bar.
func setInProgress(inProgress bool) {
	w.Dispatch(func() {
		if inProgress {
			w.Eval("window.setInProgressState(true)")
		} else {
			w.Eval("window.setInProgressState(false)")
		}
	})
}
