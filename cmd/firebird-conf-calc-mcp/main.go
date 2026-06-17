package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/IBSurgeon/FirebirdConfCalcMCP/internal/calculator"
	"github.com/IBSurgeon/FirebirdConfCalcMCP/internal/credentials"
	"github.com/IBSurgeon/FirebirdConfCalcMCP/internal/server"
	"github.com/IBSurgeon/FirebirdConfCalcMCP/internal/version"
)

func main() {
	credentialsFlag := flag.String("credentials", "", "Path to credentials file (user/password)")
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	credPath := credentials.ResolvePath(*credentialsFlag)

	app := &server.App{
		Version:         version.Version,
		CredentialsPath: credPath,
		Calculator:      calculator.NewClient(""),
		VersionChecker:  version.NewChecker(),
	}

	if *showVersion {
		printVersion(app)
		os.Exit(0)
	}

	if err := server.Run(context.Background(), app); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func printVersion(app *server.App) {
	tools := server.CurrentTools()
	update := app.VersionChecker.Check(version.Version, tools)

	info := map[string]any{
		"current_version": version.Version,
		"commit":          version.Commit,
		"build_date":      version.Date,
		"current_tools":   tools,
		"update_message":  update.UpdateMessage,
	}
	if update.UpdateAvailable {
		info["update_available"] = true
		info["latest_version"] = update.LatestVersion
		info["release_url"] = update.ReleaseURL
		if len(update.NewTools) > 0 {
			info["new_tools"] = update.NewTools
		}
	}

	enc := json.NewEncoder(os.Stderr)
	enc.SetIndent("", "  ")
	_ = enc.Encode(info)
}
