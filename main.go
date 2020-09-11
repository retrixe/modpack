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

const modpackVersion = "1.0.1"

var selectedVersion = "1.16.2"
var selectedVersionMutex sync.Mutex
var installFabricOpt = true
var installFabricOptMutex sync.Mutex

var w webview.WebView

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "--version" {
		println("modpack version " + modpackVersion)
		return
	} else if len(os.Args) >= 2 && os.Args[1] == "install" {
		InteractiveCliInstall()
		return
	}

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
	err := installMods(setProgress)
	if err != nil {
		handleError(err)
		return
	}
	enableButtons()
	hideProgress()
	showMessage()
}

func installMods(updateProgress func(string)) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	minecraftFolder := filepath.Join(home, ".minecraft") // TODO advanced opts
	if runtime.GOOS == "darwin" {
		minecraftFolder = filepath.Join(home, "Library", "Application Support", "minecraft")
	} else if runtime.GOOS == "windows" {
		minecraftFolder = filepath.Join(home, "AppData", "Roaming", ".minecraft")
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
		err = unzipFile(file, filepath.Join(minecraftFolder, "versions"))
		if err != nil {
			return err
		}
	}
	updateProgress("Downloading my mods for " + selectedVersion + "...")
	file, err := downloadMods(modVersion.URL)
	if err != nil {
		return err
	}
	// Check if there's already a mod folder.
	_, err = os.Stat(filepath.Join(minecraftFolder, "mods"))
	modsExist := false // Check if the mod folder contains the same version of mods as our pack.
	if err == nil {
		modsExist = getInstalledModsVersion(minecraftFolder) == getMajorMinecraftVersion(selectedVersion)
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err == nil && !modsExist {
		os.Rename(filepath.Join(minecraftFolder, "mods"), filepath.Join(minecraftFolder, "oldmodfolder"))
		updateProgress("Renamed old mods folder to oldmodfolder!") // todo more explicit
	} else if err != nil && os.IsNotExist(err) {
		updateProgress("Creating mods folder...")
		if err = os.MkdirAll(filepath.Join(minecraftFolder, "mods"), os.ModePerm); err != nil {
			return err
		}
	}
	updateProgress("Unzipping my mods...")
	if modsExist {
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
		if &modsData != nil {
			if err = moveOldMods(modsData, minecraftFolder, r); err != nil {
				return err
			}
		}
	} else {
		err = unzipFile(file, filepath.Join(minecraftFolder, "mods"))
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile( // Write the modsversion.txt.
		filepath.Join(minecraftFolder, "mods", "modsversion.txt"), []byte(selectedVersion), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func handleError(err error) {
	log.Println(err)
	setError(err.Error())
	hideProgress()
	enableButtons()
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

func getInstalledModsVersion(location string) string {
	file, err := os.Open(filepath.Join(location, "mods", "modsversion.txt"))
	defer file.Close()
	if err != nil {
		return ""
	}
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return ""
	}
	return getMajorMinecraftVersion(string(contents))
}

func getMajorMinecraftVersion(version string) string {
	lastIndex := strings.LastIndex(version, ".")
	if lastIndex == -1 || strings.Index(version, ".") == lastIndex {
		return version
	}
	return version[:lastIndex]
}

func moveOldMods(modsData ModsData, minecraftFolder string, r *zip.Reader) error {
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
	return nil
}

func unzipFile(zipFile []byte, location string) error {
	// Uses: os, io, strings, filepath, zip, bytes
	r, err := zip.NewReader(bytes.NewReader(zipFile), int64(len(zipFile)))
	if err != nil {
		return err
	}
	for _, f := range r.File {
		if f.Name == "mods.json" { // Ignore /mods.json during extraction.
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

// ModsData ... JSON containing data on mods inside a zip.
type ModsData struct {
	Mods    map[string]string `json:"mods"`
	OldMods map[string]string `json:"oldmods"`
}
