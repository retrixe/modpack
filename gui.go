// +build !clionly

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
//go:embed faq.html
var Faq string

// HTML is the HTML for the main page.
//go:embed modpack.html
var HTML string

var guiDialogQueryResponse bool
var guiDialogQueryResponseMutex sync.Mutex

func runGui() {
	debug := false
	if val, exists := os.LookupEnv("DEBUG"); exists {
		debug = val == "true"
	}
	w = webview.New(debug)
	defer w.Destroy()
	w.SetSize(540, 360, webview.HintNone)
	w.SetTitle("ibu's mod installer")
	w.Bind("changeVersion", func(name string) {
		selectedVersionMutex.Lock()
		minecraftFolderMutex.Lock()
		defer selectedVersionMutex.Unlock()
		defer minecraftFolderMutex.Unlock()
		selectedVersion = name
		if areModsUpdatable() == getMajorMinecraftVersion(selectedVersion) {
			w.Eval("document.getElementById('install').innerHTML = 'Update'")
		} else {
			w.Eval("document.getElementById('install').innerHTML = 'Install'")
		}
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
		if areModsUpdatable() == getMajorMinecraftVersion(selectedVersion) {
			w.Eval("document.getElementById('install').innerHTML = 'Update'")
		} else {
			w.Eval("document.getElementById('install').innerHTML = 'Install'")
		}
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
		if areModsUpdatable() == getMajorMinecraftVersion(selectedVersion) {
			w.Eval("document.getElementById('install').innerHTML = 'Update'")
		} else {
			w.Eval("document.getElementById('install').innerHTML = 'Install'")
		}
		folder := strings.ReplaceAll(strings.ReplaceAll(directory, "\\", "\\\\"), "\"", "\\\"")
		w.Eval("document.getElementById('gamedir-input').value = \"" + folder + "\"")
	})
	w.Bind("installMods", func() { go initiateInstall() })
	w.Bind("showFaq", func() { w.Navigate("data:text/html," + string(Faq)) })
	w.Bind("showGui", func() { w.Navigate("data:text/html," + strings.ReplaceAll(string(HTML), "+", "%2B")) })
	w.Navigate("data:text/html," + strings.ReplaceAll(string(HTML), "+", "%2B"))
	w.Run()
}

func initiateInstall() {
	selectedVersionMutex.Lock()
	installFabricOptMutex.Lock()
	minecraftFolderMutex.Lock()
	defer selectedVersionMutex.Unlock()
	defer installFabricOptMutex.Unlock()
	defer minecraftFolderMutex.Unlock()
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
	guiDialogQueryResponseMutex.Lock()
	defer guiDialogQueryResponseMutex.Unlock()
	return guiDialogQueryResponse
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
