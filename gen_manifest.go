package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Manifest struct {
	SDKs []Asset `json:"sdks"`
	DDIs []DDI   `json:"ddis"`
}

type Asset struct {
	Name     string `json:"name"`
	Platform string `json:"platform"`
	URL      string `json:"url"`
	Hash     string `json:"hash,omitempty"`
}

type DDI struct {
	Version  string `json:"version"`
	Platform string `json:"platform"`
	URL      string `json:"url"`
	SigURL   string `json:"sig_url"`
	Hash     string `json:"hash,omitempty"`
}

func calculateMD5(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	manifest := Manifest{}
	// Usamos la URL raw de GitHub apuntando a la rama main
	baseURL := "https://raw.githubusercontent.com/nikitacontreras/gios-platform-assets/main"
	platforms := []string{"iPhoneOS", "AppleTVOS", "WatchOS"}

	// Scan SDKs per platform
	for _, p := range platforms {
		dir := filepath.Join("sdk", p)
		files, _ := ioutil.ReadDir(dir)
		for _, f := range files {
			// Si es un directorio .sdk o un archivo .tar.xz/.zip
			if (f.IsDir() && strings.HasSuffix(f.Name(), ".sdk")) || strings.HasSuffix(f.Name(), ".tar.xz") {
				name := f.Name()
				// Si es carpeta, el link debe ser al archivo comprimido que gios espera
				if f.IsDir() {
					name = name + ".tar.xz"
				}
				
				manifest.SDKs = append(manifest.SDKs, Asset{
					Name:     f.Name(),
					Platform: p,
					URL:      fmt.Sprintf("%s/sdk/%s/%s", baseURL, p, name),
				})
			}
		}

		// Scan DDIs per platform
		ddiDir := filepath.Join("ddi", p)
		ddiFiles, _ := ioutil.ReadDir(ddiDir)
		for _, f := range ddiFiles {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".zip") {
				fullPath := filepath.Join(ddiDir, f.Name())
				hash := calculateMD5(fullPath)
				version := strings.TrimSuffix(f.Name(), ".zip")
				
				manifest.DDIs = append(manifest.DDIs, DDI{
					Version:  version,
					Platform: p,
					URL:      fmt.Sprintf("%s/ddi/%s/%s", baseURL, p, f.Name()),
					SigURL:   fmt.Sprintf("%s/ddi/%s/%s.signature", baseURL, p, f.Name()),
					Hash:     hash,
				})
			}
		}
	}

	data, _ := json.MarshalIndent(manifest, "", "  ")
	ioutil.WriteFile("assets.json", data, 0644)
	fmt.Println("[+] assets.json updated with Raw GitHub URLs!")
}
