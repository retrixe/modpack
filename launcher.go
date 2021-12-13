//go:build launcher

package main

import (
	"os"
	"os/exec"
	"path/filepath"

	_ "embed"
)

//go:embed modpack.exe
var modpackExe []byte

//go:embed webview.dll
var webviewDll []byte

//go:embed WebView2Loader.dll
var webview2LoaderDll []byte

func main() {
	// Extract modpack.exe to %LocalAppData%.
	folder, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(filepath.Join(folder, "modpack"), os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filepath.Join(folder, "modpack", "modpack.exe"), modpackExe, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filepath.Join(folder, "modpack", "webview.dll"), webviewDll, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filepath.Join(folder, "modpack", "WebView2Loader.dll"), webview2LoaderDll, os.ModePerm)
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(filepath.Join(folder, "modpack", "modpack.exe"))
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	cmd.Process.Release()
	if err != nil {
		panic(err)
	}
}
