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

// Faq is the HTML for the FAQ page.
//go:embed src/faq.html
var Faq string

// HTML is the HTML for the main page.
//go:embed src/modpack.html
var HTML string

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
	w.SetTitle("ibu's mod installer")
	// Bind a function to initiate React via webview.Eval.
	w.Bind("initiateReact", func() { w.Eval(js) })
	w.Bind("changeVersion", func(name string) {
		selectedVersionMutex.Lock()
		minecraftFolderMutex.Lock()
		defer selectedVersionMutex.Unlock()
		defer minecraftFolderMutex.Unlock()
		selectedVersion = name
		/* OLD: if areModsUpdatable() == selectedVersion {
			w.Eval("document.getElementById('install').innerHTML = 'Update'")
		} else {
			w.Eval("document.getElementById('install').innerHTML = 'Install'")
		} */
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
		// OLD: w.Eval("document.getElementById('gamedir-input').value = \"" + folder + "\"")
		w.Eval("window.setMinecraftFolderState(\"" + folder + "\")")
	})
	w.Bind("installMods", func() { go initiateInstall() })
	w.Bind("showFaq", func() { w.Navigate("data:text/html," + string(Faq)) })
	w.Bind("showGui", func() { w.Navigate("data:text/html," + strings.ReplaceAll(string(HTML), "+", "%2B")) })
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
	hideMessage()
	hideError()
	showProgress()
	disableButtons()
	err := installMods(setProgress, queryUser)
	if err != nil && err.Error() != "Cancelled" {
		handleError(err)
		return
	}
	enableButtons()
	hideProgress()
	showMessage()
}

func handleError(err error) {
	log.Println(err)
	setError(err.Error())
	hideProgress()
	enableButtons()
}

func queryUser(query string) bool {
	guiDialogQueryResponseMutex.Lock()
	w.Dispatch(func() {
		w.Eval("document.getElementById('query').textContent = `" + query + "`")
		w.Eval("M.Modal.getInstance(document.getElementById('modal1')).open()")
	})
	// This waits for the mutex to unlock.
	guiDialogQueryResponseMutex.Lock()
	defer guiDialogQueryResponseMutex.Unlock()
	return guiDialogQueryResponse
}

func checkUpdatableAndUpdateVersion() {
	updatable := areModsUpdatable()
	if updatable != "" {
		selectedVersion = updatable
		// OLD: w.Eval("document.getElementById('select-version').value = '" + selectedVersion + "'")
		// OLD: w.Eval("document.getElementById('install').innerHTML = 'Update'")
		w.Eval("setMinecraftVersionState('" + selectedVersion + "')")
		w.Eval("setUpdatableVersionState('" + selectedVersion + "')")
	} else {
		selectedVersion = defaultVersion
		// OLD: w.Eval("document.getElementById('select-version').value = '" + defaultVersion + "'")
		// OLD: w.Eval("document.getElementById('install').innerHTML = 'Install'")
		w.Eval("setMinecraftVersionState('" + selectedVersion + "')")
		w.Eval("setUpdatableVersionState('')")
	}
}

func disableButtons() {
	w.Dispatch(func() {
		w.Eval(`
      document.getElementById('faq').setAttribute('disabled', 'disabled')
			document.getElementById('install').setAttribute('disabled', 'disabled')
      document.getElementById('install-fabric').setAttribute('disabled', 'disabled')
			document.getElementById('select-version').setAttribute('disabled', 'disabled')
		`)
	})
}
func enableButtons() {
	w.Dispatch(func() {
		w.Eval(`
      document.getElementById('faq').removeAttribute('disabled')
			document.getElementById('install').removeAttribute('disabled')
      document.getElementById('install-fabric').removeAttribute('disabled')
			document.getElementById('select-version').removeAttribute('disabled')
		`)
	})
}

func showProgress() {
	w.Dispatch(func() {
		w.Eval("document.getElementById('progress').removeAttribute('style'); " +
			"document.getElementById('progress-display').removeAttribute('style')")
	})
}
func setProgress(content string) {
	w.Dispatch(func() {
		w.Eval("document.getElementById('progress').textContent = '" + content + "'")
	})
}
func hideProgress() {
	setProgress("")
	w.Dispatch(func() {
		w.Eval("document.getElementById('progress').setAttribute('style', 'display: none;'); " +
			"document.getElementById('progress-display').setAttribute('style', 'display: none;')")
	})
}

func setError(content string) {
	w.Dispatch(func() {
		w.Eval("document.getElementById('error').removeAttribute('style'); " +
			"document.getElementById('error').textContent = 'Error: " + content + "'")
	})
}
func hideError() {
	setError("")
	w.Dispatch(func() {
		w.Eval("document.getElementById('error').setAttribute('style', 'display: none;');")
	})
}

func showMessage() {
	w.Dispatch(func() {
		w.Eval("document.getElementById('message').removeAttribute('style')")
	})
}
func hideMessage() {
	w.Dispatch(func() {
		w.Eval("document.getElementById('message').setAttribute('style', 'display: none;');")
	})
}
