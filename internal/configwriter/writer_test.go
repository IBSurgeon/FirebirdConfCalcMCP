package configwriter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteNewFiles(t *testing.T) {
	dir := t.TempDir()
	result, err := Write(Params{
		OutputDir:     dir,
		FirebirdConf:  "ServerMode = Classic\n",
		DatabasesConf: "{ DefaultDbCachePages = 250 }\n",
	})
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if len(result.Written) != 2 {
		t.Fatalf("written = %v", result.Written)
	}
	if _, err := os.Stat(filepath.Join(dir, firebirdConfName)); err != nil {
		t.Fatalf("firebird.conf missing: %v", err)
	}
}

func TestWriteCreatesBackup(t *testing.T) {
	dir := t.TempDir()
	existing := filepath.Join(dir, firebirdConfName)
	if err := os.WriteFile(existing, []byte("old content\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := Write(Params{
		OutputDir:    dir,
		FirebirdConf: "new content\n",
	})
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if len(result.Backups) != 1 {
		t.Fatalf("backups = %v", result.Backups)
	}
	if !strings.Contains(result.Backups[0], ".bak.") {
		t.Fatalf("backup path = %s", result.Backups[0])
	}

	data, err := os.ReadFile(existing)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "new content\n" {
		t.Fatalf("content = %q", string(data))
	}
}

func TestDryRun(t *testing.T) {
	dir := t.TempDir()
	existing := filepath.Join(dir, firebirdConfName)
	if err := os.WriteFile(existing, []byte("old\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := Write(Params{
		OutputDir:    dir,
		FirebirdConf: "new\n",
		DryRun:       true,
	})
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if len(result.Written) != 1 {
		t.Fatalf("written = %v", result.Written)
	}

	data, err := os.ReadFile(existing)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "old\n" {
		t.Fatal("dry run should not modify file")
	}
}

func TestSkipEmptyDatabasesConf(t *testing.T) {
	dir := t.TempDir()
	result, err := Write(Params{
		OutputDir:    dir,
		FirebirdConf: "ServerMode = Classic\n",
	})
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if len(result.Written) != 1 {
		t.Fatalf("written = %v", result.Written)
	}
	if _, err := os.Stat(filepath.Join(dir, databasesConfName)); !os.IsNotExist(err) {
		t.Fatal("databases.conf should not be created")
	}
}
