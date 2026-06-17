package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Manifest struct {
	LatestVersion string   `json:"latest_version"`
	ReleaseURL    string   `json:"release_url"`
	Tools         []string `json:"tools"`
	ReleaseNotes  string   `json:"release_notes"`
}

func ParseManifest(data []byte) (*Manifest, error) {
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse version manifest: %w", err)
	}
	if m.LatestVersion == "" {
		return nil, fmt.Errorf("version manifest missing latest_version")
	}
	return &m, nil
}

func FetchManifest(client *http.Client, url string) (*Manifest, error) {
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("version manifest HTTP %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return ParseManifest(data)
}
