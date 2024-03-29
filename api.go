package main

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
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

// QuiltMavenRepositoryURL points to the base URL of Quilt's Maven repo.
const QuiltMavenRepositoryURL = "https://maven.quiltmc.org/repository/release/"

// FabricMavenRepositoryURL points to the base URL of Fabric's Maven repo.
const FabricMavenRepositoryURL = "https://maven.fabricmc.net/"

func getLatestFabric(quilt bool) (string, error) {
	url := FabricMavenRepositoryURL + "net/fabricmc/fabric-loader/maven-metadata.xml"
	if quilt {
		url = QuiltMavenRepositoryURL + "org/quiltmc/quilt-loader/maven-metadata.xml"
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var versions FabricVersionResponse
	xml.NewDecoder(resp.Body).Decode(&versions)
	if quilt {
		latest := ""
		for _, v := range versions.Versioning.Versions {
			if latest < v && !strings.Contains(v, "-") {
				latest = v
			}
		}
		return latest, nil
	}
	return versions.Versioning.Latest, nil
}

func downloadFabric(version string, fabricVersion string) ([]byte, error) {
	return downloadFile("https://meta.fabricmc.net/v2/versions/loader/" + version + "/" +
		url.QueryEscape(fabricVersion) + "/profile/zip")
}

func downloadQuilt(version string, quiltVersion string) ([]byte, error) {
	f, err := downloadFile("https://meta.quiltmc.org/v3/versions/loader/" + version + "/" +
		url.QueryEscape(quiltVersion) + "/profile/json")
	if err != nil {
		return nil, err
	}
	sep := strings.Split(string(f), "org.quiltmc:hashed")
	sep[1] = strings.Replace(string(sep[1]), QuiltMavenRepositoryURL, FabricMavenRepositoryURL, 1)
	return []byte(strings.Join(sep, "net.fabricmc:intermediary")), nil
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
	XMLName  xml.Name `xml:"versioning"`
	Latest   string   `xml:"latest"`
	Release  string   `xml:"release"`
	Versions []string `xml:"versions>version"`
}

// FabricVersionNames is a list of Fabric versions.
type FabricVersionNames struct {
	XMLName xml.Name `xml:"version"`
	Version string   `xml:"version"`
}

// ModVersion is a JSON containing version mappings of mods.
type ModVersion struct {
	FullVersion string `json:"fullVersion"`
	Fabric      string `json:"fabric"`
	URL         string `json:"url"`
}
