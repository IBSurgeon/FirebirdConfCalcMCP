package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	defaultCredFile = "password_api.txt"
	serverBuildPkg  = "./cmd/firebird-conf-calc-mcp"
)

func TestE2E_MCPServer(t *testing.T) {
	if os.Getenv("CC_E2E") == "0" {
		t.Skip("e2e disabled (CC_E2E=0)")
	}

	root := moduleRoot(t)
	credPath := resolveCredentialsPath(root)
	if _, err := os.Stat(credPath); err != nil {
		t.Skipf("credentials file not found: %s", credPath)
	}

	bin := buildServerBinary(t, root)
	session := connectServer(t, bin, credPath)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	t.Run("ListTools", func(t *testing.T) {
		want := []string{
			"calculate_firebird_config",
			"write_firebird_configs",
			"get_server_info",
		}
		var got []string
		for tool, err := range session.Tools(ctx, nil) {
			if err != nil {
				t.Fatalf("Tools() error = %v", err)
			}
			got = append(got, tool.Name)
		}
		slices.Sort(got)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Fatalf("tools = %v, want %v", got, want)
		}
	})

	t.Run("GetServerInfo", func(t *testing.T) {
		info := callToolMap(t, ctx, session, "get_server_info", nil)
		if info["current_version"] == nil {
			t.Fatalf("missing current_version: %#v", info)
		}
		tools, ok := info["current_tools"].([]any)
		if !ok || len(tools) == 0 {
			t.Fatalf("missing current_tools: %#v", info["current_tools"])
		}
		if info["credential_file"] == nil {
			t.Fatalf("missing credential_file: %#v", info)
		}
	})

	t.Run("CalculateFirebirdConfig", func(t *testing.T) {
		result, err := callTool(ctx, session, "calculate_firebird_config", map[string]any{
			"server_version":      "fb3",
			"server_architecture": "Classic",
			"cores":               8,
			"ram":                 16,
			"count_users":         100,
			"page_size":           4096,
			"size_db":             100,
			"name_main_db":        "mainDB",
			"path_to_main_db":     "c:/data/main.fdb",
		})
		if err != nil {
			if isCredentialError(err) {
				t.Skipf("skipping live API test: credentials in %s were rejected by cc.ib-aid.com (%v)", credPath, err)
			}
			t.Fatalf("calculate_firebird_config failed: %v", err)
		}

		firebirdConf, _ := result["firebird_conf"].(string)
		if strings.TrimSpace(firebirdConf) == "" {
			t.Fatalf("empty firebird_conf: %#v", result)
		}
		if apiVersion, _ := result["api_version"].(string); apiVersion == "" {
			t.Fatalf("missing api_version: %#v", result)
		}
	})

	t.Run("WriteFirebirdConfigsDryRun", func(t *testing.T) {
		outDir := t.TempDir()
		write := callToolMap(t, ctx, session, "write_firebird_configs", map[string]any{
			"output_dir":    outDir,
			"firebird_conf": "ServerMode = Classic\nDefaultDbCachePages = 2048\n",
			"databases_conf": "{\n    DefaultDbCachePages = 250\n}\n",
			"dry_run":       true,
		})

		if dryRun, _ := write["dry_run"].(bool); !dryRun {
			t.Fatalf("dry_run = %#v, want true", write["dry_run"])
		}
		written, ok := write["written"].([]any)
		if !ok || len(written) == 0 {
			t.Fatalf("expected written paths: %#v", write["written"])
		}
	})
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	cmd := exec.Command("go", "env", "GOMOD")
	cmd.Dir = "."
	out, err := cmd.Output()
	if err == nil {
		modPath := strings.TrimSpace(string(out))
		if modPath != "" && modPath != "/dev/null" {
			return filepath.Dir(modPath)
		}
	}

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine module root")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), ".."))
}

func resolveCredentialsPath(root string) string {
	if path := os.Getenv("CC_CREDENTIALS_FILE"); path != "" {
		if filepath.IsAbs(path) {
			return path
		}
		return filepath.Join(root, path)
	}
	return filepath.Join(root, defaultCredFile)
}

func buildServerBinary(t *testing.T, root string) string {
	t.Helper()
	name := "firebird-conf-calc-mcp"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	bin := filepath.Join(t.TempDir(), name)

	cmd := exec.Command("go", "build", "-o", bin, serverBuildPkg)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build server binary: %v\n%s", err, out)
	}
	return bin
}

func connectServer(t *testing.T, bin, credPath string) *mcp.ClientSession {
	t.Helper()
	cmd := exec.Command(bin, "--credentials", credPath)
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "firebird-conf-calc-mcp-e2e",
		Version: "1.0.0",
	}, nil)

	session, err := client.Connect(context.Background(), &mcp.CommandTransport{Command: cmd}, nil)
	if err != nil {
		t.Fatalf("connect to MCP server: %v", err)
	}
	t.Cleanup(func() { session.Close() })
	return session
}

func callToolMap(t *testing.T, ctx context.Context, session *mcp.ClientSession, name string, args map[string]any) map[string]any {
	t.Helper()
	result, err := callTool(ctx, session, name, args)
	if err != nil {
		t.Fatalf("CallTool(%s) error = %v", name, err)
	}
	return result
}

func callTool(ctx context.Context, session *mcp.ClientSession, name string, args map[string]any) (map[string]any, error) {
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      name,
		Arguments: args,
	})
	if err != nil {
		return nil, err
	}
	if result.IsError {
		return nil, fmt.Errorf("%s", toolErrorText(result))
	}

	if result.StructuredContent != nil {
		switch v := result.StructuredContent.(type) {
		case map[string]any:
			return v, nil
		default:
			data, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("marshal structured content: %w", err)
			}
			var out map[string]any
			if err := json.Unmarshal(data, &out); err != nil {
				return nil, fmt.Errorf("unmarshal structured content: %w", err)
			}
			return out, nil
		}
	}

	for _, content := range result.Content {
		if text, ok := content.(*mcp.TextContent); ok {
			var out map[string]any
			if err := json.Unmarshal([]byte(text.Text), &out); err == nil {
				return out, nil
			}
		}
	}
	return nil, fmt.Errorf("no structured content in tool result")
}

func isCredentialError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "login or password") || strings.Contains(msg, "passapi")
}

func toolErrorText(result *mcp.CallToolResult) string {
	var parts []string
	for _, content := range result.Content {
		if text, ok := content.(*mcp.TextContent); ok {
			parts = append(parts, text.Text)
		}
	}
	if len(parts) == 0 {
		return fmt.Sprintf("%#v", result)
	}
	return strings.Join(parts, "; ")
}
