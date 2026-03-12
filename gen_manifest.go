// Triggering workflow update
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
	SigURL   string `json:"sig_url,omitempty"`
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
		fmt.Printf("[!] Warning: Could not open %s for hashing: %v\n", path, err)
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
	// Usamos la URL de los release assets para las descargas reales
	// assets.json se mantendrá en raw github para consulta rápida
	releaseBaseURL := "https://github.com/nikitacontreras/gios-platform-assets/releases/latest/download"
	platforms := []string{"iPhoneOS", "AppleTVOS", "WatchOS"}

	fmt.Println("[gios-assets] Generating manifest from Release assets...")

	// Scan SDKs per platform
	for _, p := range platforms {
		dir := filepath.Join("sdk", p)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			continue
		}
		
		for _, f := range files {
			// El generador ahora espera encontrar archivos .tar.xz (comprimidos previamente en el CI)
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".tar.xz") {
				prettyName := strings.TrimSuffix(f.Name(), ".tar.xz")
				fullPath := filepath.Join(dir, f.Name())
				hash := calculateMD5(fullPath)
				
				asset := Asset{
					Name:     prettyName,
					Platform: p,
					URL:      fmt.Sprintf("%s/%s", releaseBaseURL, f.Name()),
					Hash:     hash,
				}

				// Verificar si existe firma .signature
				sigFile := fullPath + ".signature"
				if _, err := os.Stat(sigFile); err == nil {
					asset.SigURL = fmt.Sprintf("%s/%s.signature", releaseBaseURL, f.Name())
				}

				manifest.SDKs = append(manifest.SDKs, asset)
				fmt.Printf(" [+] Indexing SDK: %s (MD5: %s)\n", f.Name(), hash)
			}
		}

		// Scan DDIs per platform
		ddiDir := filepath.Join("ddi", p)
		ddiFiles, err := ioutil.ReadDir(ddiDir)
		if err != nil {
			continue
		}

		for _, f := range ddiFiles {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".zip") {
				fullPath := filepath.Join(ddiDir, f.Name())
				hash := calculateMD5(fullPath)
				version := strings.TrimSuffix(f.Name(), ".zip")
				
				manifest.DDIs = append(manifest.DDIs, DDI{
					Version:  version,
					Platform: p,
					URL:      fmt.Sprintf("%s/%s", releaseBaseURL, f.Name()),
					SigURL:   fmt.Sprintf("%s/%s.signature", releaseBaseURL, f.Name()),
					Hash:     hash,
				})
				fmt.Printf(" [+] Indexing DDI: %s (MD5: %s)\n", f.Name(), hash)
			}
		}
	}

	data, _ := json.MarshalIndent(manifest, "", "  ")
	ioutil.WriteFile("assets.json", data, 0644)
	fmt.Println("\n[Success] assets.json updated based on Release Assets!")
}
