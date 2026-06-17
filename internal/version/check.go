package version

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/mod/semver"
)

type UpdateInfo struct {
	UpdateAvailable bool     `json:"update_available"`
	LatestVersion   string   `json:"latest_version,omitempty"`
	ReleaseURL      string   `json:"release_url,omitempty"`
	UpdateMessage   string   `json:"update_message"`
	NewTools        []string `json:"new_tools,omitempty"`
}

type Checker struct {
	manifestURL string
	client      *http.Client
	cacheTTL    time.Duration

	mu       sync.Mutex
	cached   *Manifest
	cachedAt time.Time
}

func NewChecker() *Checker {
	url := os.Getenv("CC_VERSION_MANIFEST_URL")
	if url == "" {
		url = DefaultManifestURL
	}
	return &Checker{
		manifestURL: url,
		client:      &http.Client{Timeout: 10 * time.Second},
		cacheTTL:    24 * time.Hour,
	}
}

func (c *Checker) Check(currentVersion string, currentTools []string) UpdateInfo {
	current := normalizeSemver(currentVersion)
	display := strings.TrimPrefix(current, "v")
	info := UpdateInfo{
		UpdateMessage: fmt.Sprintf("Current version %s. You are on the latest release.", display),
	}

	if os.Getenv("CC_UPDATE_CHECK") == "0" {
		info.UpdateMessage = fmt.Sprintf("Current version %s. Update check disabled.", display)
		return info
	}

	manifest, err := c.getManifest()
	if err != nil || manifest == nil {
		info.UpdateMessage = fmt.Sprintf("Current version %s. Could not check for updates.", display)
		return info
	}

	latest := normalizeSemver(manifest.LatestVersion)
	if !semver.IsValid(current) || !semver.IsValid(latest) {
		info.UpdateMessage = fmt.Sprintf("Current version %s. Could not compare versions.", display)
		return info
	}

	if semver.Compare(current, latest) >= 0 {
		return info
	}

	info.UpdateAvailable = true
	info.LatestVersion = strings.TrimPrefix(latest, "v")
	info.ReleaseURL = manifest.ReleaseURL
	info.NewTools = diffTools(currentTools, manifest.Tools)

	if len(info.NewTools) > 0 {
		info.UpdateMessage = fmt.Sprintf(
			"Current version %s. Newer version %s exists — update to get more tools: %s",
			display,
			info.LatestVersion,
			strings.Join(info.NewTools, ", "),
		)
	} else {
		info.UpdateMessage = fmt.Sprintf(
			"Current version %s. Newer version %s exists — update to get the latest release.",
			display,
			info.LatestVersion,
		)
	}
	return info
}

func (c *Checker) getManifest() (*Manifest, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cached != nil && time.Since(c.cachedAt) < c.cacheTTL {
		return c.cached, nil
	}

	manifest, err := FetchManifest(c.client, c.manifestURL)
	if err != nil {
		return nil, err
	}
	c.cached = manifest
	c.cachedAt = time.Now()
	return manifest, nil
}

func normalizeSemver(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}
	return v
}

func diffTools(current, latest []string) []string {
	have := make(map[string]struct{}, len(current))
	for _, t := range current {
		have[t] = struct{}{}
	}
	var missing []string
	for _, t := range latest {
		if _, ok := have[t]; !ok {
			missing = append(missing, t)
		}
	}
	return missing
}
