package server

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/IBSurgeon/FirebirdConfCalcMCP/internal/calculator"
	"github.com/IBSurgeon/FirebirdConfCalcMCP/internal/configwriter"
	"github.com/IBSurgeon/FirebirdConfCalcMCP/internal/credentials"
	"github.com/IBSurgeon/FirebirdConfCalcMCP/internal/version"
)

type App struct {
	Version         string
	CredentialsPath string
	Calculator      *calculator.Client
	VersionChecker  *version.Checker
}

type ServerInfoOutput struct {
	CurrentVersion  string   `json:"current_version"`
	Commit          string   `json:"commit"`
	BuildDate       string   `json:"build_date"`
	APIVersion      string   `json:"api_version"`
	CredentialFile  string   `json:"credential_file"`
	CurrentTools    []string `json:"current_tools"`
	UpdateAvailable bool     `json:"update_available"`
	LatestVersion   string   `json:"latest_version,omitempty"`
	UpdateMessage   string   `json:"update_message"`
	NewTools        []string `json:"new_tools,omitempty"`
	ReleaseURL      string   `json:"release_url,omitempty"`
	RepositoryURL   string   `json:"repository_url"`
}

func (a *App) loadCredentials() (credentials.Credentials, error) {
	path, err := filepath.Abs(a.CredentialsPath)
	if err != nil {
		path = a.CredentialsPath
	}
	return credentials.Load(path)
}

func (a *App) CalculateFirebirdConfig(ctx context.Context, _ *mcp.CallToolRequest, input calculator.CalculateParams) (*mcp.CallToolResult, calculator.Result, error) {
	if err := input.Validate(); err != nil {
		return toolError(err), calculator.Result{}, nil
	}

	creds, err := a.loadCredentials()
	if err != nil {
		return toolError(err), calculator.Result{}, nil
	}

	result, err := a.Calculator.Calculate(ctx, input.ToRequest(creds.MailLogin, creds.PassAPI))
	if err != nil {
		return toolError(err), calculator.Result{}, nil
	}
	return nil, *result, nil
}

func (a *App) WriteFirebirdConfigs(_ context.Context, _ *mcp.CallToolRequest, input configwriter.Params) (*mcp.CallToolResult, configwriter.Result, error) {
	result, err := configwriter.Write(input)
	if err != nil {
		return toolError(err), configwriter.Result{}, nil
	}
	return nil, *result, nil
}

func (a *App) GetServerInfo(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, ServerInfoOutput, error) {
	tools := CurrentTools()
	update := a.VersionChecker.Check(a.Version, tools)

	credPath := a.CredentialsPath
	if abs, err := filepath.Abs(credPath); err == nil {
		credPath = abs
	}

	out := ServerInfoOutput{
		CurrentVersion:  a.Version,
		Commit:          version.Commit,
		BuildDate:       version.Date,
		APIVersion:      version.APIVersion,
		CredentialFile:  credPath,
		CurrentTools:    tools,
		UpdateAvailable: update.UpdateAvailable,
		LatestVersion:   update.LatestVersion,
		UpdateMessage:   update.UpdateMessage,
		NewTools:        update.NewTools,
		ReleaseURL:      update.ReleaseURL,
		RepositoryURL:   version.RepositoryURL,
	}
	return nil, out, nil
}

func toolError(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)},
		},
	}
}
