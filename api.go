package main

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"os"
)

func getModVersions(version string) (*ModVersion, error) {
	url := "https://mythicmc.org/modpack/modpack.json"
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

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
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
	resp, err := http.Get("https://meta.fabricmc.net/v2/versions/loader/" + version + "/" +
		url.QueryEscape(fabricVersion) + "/profile/zip")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// FabricVersionResponse is the response from querying Fabric's Maven API.
type FabricVersionResponse struct {
	XMLName    xml.Name       `xml:"metadata"`
	GroupID    string         `xml:"groupId"`
	ArtifactID string         `xml:"artifactId"`
	Versioning FabricVersions `xml:"versioning"`
}

// FabricVersions contains the latest Fabric version as well as list of Fabric versions.
type FabricVersions struct {
	XMLName  xml.Name             `xml:"versioning"`
	Latest   string               `xml:"latest"`
	Release  string               `xml:"release"`
	Versions []FabricVersionNames `xml:"versions"`
}

// FabricVersionNames is a list of Fabric versions.
type FabricVersionNames struct {
	XMLName xml.Name `xml:"versions"`
	Version string   `xml:"version"`
}

// ModVersion is a JSON containing version mappings of mods.
type ModVersion struct {
	FullVersion string `json:"fullVersion"`
	Fabric      string `json:"fabric"`
	URL         string `json:"url"`
}
