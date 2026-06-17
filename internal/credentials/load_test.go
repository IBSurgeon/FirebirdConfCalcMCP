package credentials

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "creds.txt")
	content := "user: test@example.com\npassword: secret123\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	creds, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if creds.MailLogin != "test@example.com" {
		t.Fatalf("MailLogin = %q", creds.MailLogin)
	}
	if creds.PassAPI != "secret123" {
		t.Fatalf("PassAPI = %q", creds.PassAPI)
	}
}

func TestLoadMissingUser(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "creds.txt")
	if err := os.WriteFile(path, []byte("password: x\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected error for missing user")
	}
}

func TestResolvePath(t *testing.T) {
	if got := ResolvePath("/custom/path.txt"); got != "/custom/path.txt" {
		t.Fatalf("ResolvePath flag = %q", got)
	}
	t.Setenv("CC_CREDENTIALS_FILE", "")
	if got := ResolvePath(""); got != "password_api.txt" {
		t.Fatalf("ResolvePath default = %q", got)
	}
}
