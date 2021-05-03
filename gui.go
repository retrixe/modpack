// +build !clionly

package main

import (
	"log"

	"github.com/webview/webview"
)

var w webview.WebView

func runGui() {
	w = webview.New(false)
	defer w.Destroy()
	w.SetSize(600, 300, webview.HintNone)
	w.SetTitle("ibu's mod installer")
	w.Bind("changeVersion", func(name string) {
		selectedVersionMutex.Lock()
		defer selectedVersionMutex.Unlock()
		selectedVersion = name
	})
	w.Bind("toggleInstallFabric", func() {
		installFabricOptMutex.Lock()
		defer installFabricOptMutex.Unlock()
		installFabricOpt = !installFabricOpt
	})
	w.Bind("installMods", func() { go initiateInstall() })
	w.Bind("showFaq", func() { w.Navigate("data:text/html," + string(Faq)) })
	w.Bind("showGui", func() { w.Navigate("data:text/html," + string(HTML)) })
	w.Navigate("data:text/html," + string(HTML))
	w.Run()
}

func initiateInstall() {
	selectedVersionMutex.Lock()
	installFabricOptMutex.Lock()
	defer selectedVersionMutex.Unlock()
	defer installFabricOptMutex.Unlock()
	hideMessage()
	hideError()
	showProgress()
	disableButtons()
	err := installMods(setProgress, queryUser)
	if err != nil {
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
	// w.Dispatch(func() {
	// 	w.Eval("modal.open(); document.getElementById('query').textContent = '" + query + "'")
	// 	// TODO: How do we get back close/open?
	// })
	return true
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
		w.Eval("document.getElementById('progress').textContent = '" + content + "'") // TODO show %
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
