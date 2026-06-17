package configwriter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	firebirdConfName   = "firebird.conf"
	databasesConfName  = "databases.conf"
)

type Params struct {
	OutputDir     string `json:"output_dir" jsonschema:"Directory where configuration files will be written"`
	FirebirdConf  string `json:"firebird_conf" jsonschema:"Full firebird.conf content"`
	DatabasesConf string `json:"databases_conf,omitempty" jsonschema:"Full databases.conf content (optional)"`
	DryRun        bool   `json:"dry_run,omitempty" jsonschema:"If true, return planned paths without writing files"`
}

type Result struct {
	OutputDir string   `json:"output_dir"`
	Written   []string `json:"written"`
	Backups   []string `json:"backups"`
	DryRun    bool     `json:"dry_run"`
}

func Write(params Params) (*Result, error) {
	if strings.TrimSpace(params.OutputDir) == "" {
		return nil, fmt.Errorf("output_dir is required")
	}
	if strings.TrimSpace(params.FirebirdConf) == "" {
		return nil, fmt.Errorf("firebird_conf is required")
	}

	outputDir, err := validateOutputDir(params.OutputDir)
	if err != nil {
		return nil, err
	}

	result := &Result{
		OutputDir: outputDir,
		DryRun:    params.DryRun,
		Written:   []string{},
		Backups:   []string{},
	}

	files := []struct {
		name    string
		content string
	}{
		{name: firebirdConfName, content: params.FirebirdConf},
	}
	if strings.TrimSpace(params.DatabasesConf) != "" {
		files = append(files, struct {
			name    string
			content string
		}{name: databasesConfName, content: params.DatabasesConf})
	}

	for _, f := range files {
		target := filepath.Join(outputDir, f.name)
		if params.DryRun {
			result.Written = append(result.Written, target)
			if fileExists(target) {
				result.Backups = append(result.Backups, backupPath(target))
			}
			continue
		}

		if fileExists(target) {
			backup, err := createBackup(target)
			if err != nil {
				return nil, err
			}
			result.Backups = append(result.Backups, backup)
		}

		if err := atomicWrite(target, f.content); err != nil {
			return nil, err
		}
		result.Written = append(result.Written, target)
	}

	return result, nil
}

func validateOutputDir(dir string) (string, error) {
	clean := filepath.Clean(dir)
	if clean == "." || clean == "" {
		return "", fmt.Errorf("output_dir must be an absolute or explicit directory path")
	}

	abs, err := filepath.Abs(clean)
	if err != nil {
		return "", fmt.Errorf("resolve output_dir: %w", err)
	}

	if os.Getenv("CC_ALLOW_ANY_OUTPUT_DIR") != "1" {
		if strings.Contains(clean, "..") {
			return "", fmt.Errorf("output_dir must not contain parent directory references (set CC_ALLOW_ANY_OUTPUT_DIR=1 to override)")
		}
	}

	if err := os.MkdirAll(abs, 0o755); err != nil {
		return "", fmt.Errorf("create output_dir: %w", err)
	}
	return abs, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func backupPath(path string) string {
	ts := time.Now().Format("20060102-150405")
	return path + ".bak." + ts
}

func createBackup(path string) (string, error) {
	backup := backupPath(path)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %s for backup: %w", path, err)
	}
	if err := os.WriteFile(backup, data, 0o644); err != nil {
		return "", fmt.Errorf("write backup %s: %w", backup, err)
	}
	return backup, nil
}

func atomicWrite(path, content string) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if _, err := tmp.WriteString(content); err != nil {
		tmp.Close()
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("rename to %s: %w", path, err)
	}
	return nil
}
