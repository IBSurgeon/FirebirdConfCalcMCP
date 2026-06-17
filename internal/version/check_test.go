package version

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseManifest(t *testing.T) {
	data := []byte(`{
		"latest_version": "1.1.0",
		"release_url": "https://example.com/v1.1.0",
		"tools": ["a", "b"],
		"release_notes": "notes"
	}`)
	m, err := ParseManifest(data)
	if err != nil {
		t.Fatalf("ParseManifest() error = %v", err)
	}
	if m.LatestVersion != "1.1.0" {
		t.Fatalf("LatestVersion = %q", m.LatestVersion)
	}
}

func TestCheckerUpdateAvailable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{
			"latest_version": "1.1.0",
			"release_url": "https://example.com/release",
			"tools": ["calculate_firebird_config", "write_firebird_configs", "get_server_info", "new_tool"],
			"release_notes": "new tool"
		}`))
	}))
	defer srv.Close()

	checker := &Checker{
		manifestURL: srv.URL,
		client:      srv.Client(),
		cacheTTL:    0,
	}
	info := checker.Check("1.0.0", CurrentToolsForTest())
	if !info.UpdateAvailable {
		t.Fatal("expected update available")
	}
	if len(info.NewTools) != 1 || info.NewTools[0] != "new_tool" {
		t.Fatalf("new tools = %v", info.NewTools)
	}
}

func TestCheckerUpToDate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{
			"latest_version": "1.0.0",
			"release_url": "https://example.com/release",
			"tools": ["calculate_firebird_config"],
			"release_notes": ""
		}`))
	}))
	defer srv.Close()

	checker := &Checker{
		manifestURL: srv.URL,
		client:      srv.Client(),
		cacheTTL:    0,
	}
	info := checker.Check("1.0.0", []string{"calculate_firebird_config"})
	if info.UpdateAvailable {
		t.Fatal("did not expect update")
	}
}

func CurrentToolsForTest() []string {
	return []string{
		"calculate_firebird_config",
		"write_firebird_configs",
		"get_server_info",
	}
}
