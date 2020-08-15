package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/webview/webview"
)

var selectedVersion = "1.16.1"
var selectedVersionMutex sync.Mutex
var installFabricOpt = true
var installFabricOptMutex sync.Mutex

var w webview.WebView

func main() {
	/*
		html, err := ioutil.ReadFile("modpack.html")
		if err != nil {
			log.Panicln("Unable to open the GUI HTML!")
		}
		faq, err := ioutil.ReadFile("faq.html")
		if err != nil {
			log.Panicln("Unable to open the GUI HTML!")
		}
	*/
	w = webview.New(false)
	defer w.Destroy()
	w.SetSize(600, 300, webview.HintNone)
	w.SetTitle("ibu's mod installer")
	w.Navigate("data:text/html," + string(HTML))
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
	w.Bind("installMods", initiateInstall)
	w.Bind("showFaq", func() { w.Navigate("data:text/html," + string(Faq)) })
	w.Bind("showGui", func() { w.Navigate("data:text/html," + string(HTML)) })
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
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
		setError(err.Error())
		hideProgress()
		return
	}
	minecraftFolder := filepath.Join(home, ".minecraft") // TODO advanced opts, also gettng newest fabrc
	if runtime.GOOS == "darwin" {
		minecraftFolder = filepath.Join(home, "Library", "Application Support", "minecraft")
	} else if runtime.GOOS == "windows" {
		minecraftFolder = filepath.Join(home, "AppData", "Roaming", ".minecraft")
	}
	setProgress("Querying latest mod versions...")
	modVersion, err := getModVersions(selectedVersion)
	if err != nil {
		log.Println(err)
		setError(err.Error())
		hideProgress()
		return
	}
	if installFabricOpt {
		s := modVersion.Fabric
		if s == "latest" {
			setProgress("Querying latest Fabric version...")
			s, err = getLatestFabric()
			if err != nil {
				log.Println(err)
				setError(err.Error())
				hideProgress()
				return
			}
		}
		setProgress("Downloading Fabric...")
		file, err := downloadFabric(selectedVersion, s)
		if err != nil {
			log.Println(err)
			setError(err.Error())
			hideProgress()
			return
		}
		setProgress("Installing Fabric...")
		err = unzipFile(file, filepath.Join(minecraftFolder, "versions"))
		if err != nil {
			log.Println(err)
			setError(err.Error())
			hideProgress()
			return
		}
	}
	setProgress("Downloading my mods for " + selectedVersion + "...")
	file, err := downloadMods(modVersion.URL)
	if err != nil {
		log.Println(err)
		setError(err.Error())
		hideProgress()
		return
	}
	// Check if there's already a mod folder.
	_, err = os.Stat(filepath.Join(minecraftFolder, "mods"))
	if err != nil && !os.IsNotExist(err) {
		log.Println(err)
		setError(err.Error())
		hideProgress()
		return
	} else if err == nil {
		os.Rename(filepath.Join(minecraftFolder, "mods"), filepath.Join(minecraftFolder, "oldmodfolder"))
		setProgress("Renamed old mods folder to oldmodfolder!") // todo more explicit
	}
	setProgress("Creating mods folder...")
	if err = os.MkdirAll(filepath.Join(minecraftFolder, "mods"), os.ModePerm); err != nil {
		log.Println(err)
		setError(err.Error())
		hideProgress()
		return
	}
	setProgress("Unzipping my mods...")
	err = unzipFile(file, filepath.Join(minecraftFolder, "mods"))
	if err != nil {
		log.Println(err)
		setError(err.Error())
		hideProgress()
		return
	}
	hideProgress()
	showMessage()
}

func showProgress() {
	w.Eval("document.getElementById('progress').removeAttribute('style'); " +
		"document.getElementById('progress-display').removeAttribute('style')")
}
func setProgress(content string) {
	w.Eval("document.getElementById('progress').textContent = '" + content + "'") // TODO show %
}
func hideProgress() {
	setProgress("")
	w.Eval("document.getElementById('progress').setAttribute('style', 'display: none;'); " +
		"document.getElementById('progress-display').setAttribute('style', 'display: none;')")
}

func setError(content string) {
	w.Eval("document.getElementById('error').removeAttribute('style'); " +
		"document.getElementById('error').textContent = 'Error: " + content + "'")
}
func hideError() {
	setError("")
	w.Eval("document.getElementById('error').setAttribute('style', 'display: none;');")
}

func showMessage() {
	w.Eval("document.getElementById('message').removeAttribute('style')")
}
func hideMessage() {
	w.Eval("document.getElementById('message').setAttribute('style', 'display: none;');")
}

func getLatestFabric() (string, error) {
	resp, err := http.Get("https://maven.fabricmc.net/net/fabricmc/fabric-loader/maven-metadata.xml")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var versions FabricVersionResponse
	xml.NewDecoder(resp.Body).Decode(&versions)
	return versions.Versioning.Latest, nil
}

func downloadFabric(version string, fabricVersion string) ([]byte, error) {
	resp, err := http.Get("https://fabricmc.net/download/vanilla?intermediary=" + version +
		"&loader=" + url.QueryEscape(fabricVersion))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func getModVersions(version string) (*ModVersion, error) {
	url := "https://mythicmc.info/modpack/modpack.json"
	if val, exists := os.LookupEnv("MODS_JSON_URL"); exists {
		url = val
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	var res map[string]ModVersion
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	ver := res[version]
	return &ver, nil
}

func downloadMods(url string) ([]byte, error) {
	/*
		url := "https://cdn.discordapp.com/attachments/402428932839833601/744123637291941888/1161_mods.zip"
		if version == "1.14.4" {
			url = "https://cdn.discordapp.com/attachments/402428932839833601/744123637291941888/1161_mods.zip"
		} else if version == "1.15.2" {
			url = "https://cdn.discordapp.com/attachments/402428932839833601/744123637291941888/1161_mods.zip"
		}
	*/
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func unzipFile(zipFile []byte, location string) error {
	// Uses: os, io, strings, filepath, zip, bytes
	r, err := zip.NewReader(bytes.NewReader(zipFile), int64(len(zipFile)))
	if err != nil {
		return err
	}
	for _, f := range r.File {
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

// FabricVersionResponse ... Response from querying Fabric's Maven API.
type FabricVersionResponse struct {
	XMLName    xml.Name       `xml:"metadata"`
	GroupID    string         `xml:"groupId"`
	ArtifactID string         `xml:"artifactId"`
	Versioning FabricVersions `xml:"versioning"`
}

// FabricVersions ... Contains latest Fabric version as well as list of Fabric versions.
type FabricVersions struct {
	XMLName  xml.Name             `xml:"versioning"`
	Latest   string               `xml:"latest"`
	Release  string               `xml:"release"`
	Versions []FabricVersionNames `xml:"versions"`
}

// FabricVersionNames ... List of Fabric versions.
type FabricVersionNames struct {
	XMLName xml.Name `xml:"versions"`
	Version string   `xml:"version"`
}

// ModVersion ... JSON containing version mappings of mods.
type ModVersion struct {
	Fabric string `json:"fabric"`
	URL    string `json:"url"`
}
