package server

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const serverName = "firebird-config-calculator"

var toolNames = []string{
	"calculate_firebird_config",
	"write_firebird_configs",
	"get_server_info",
}

func CurrentTools() []string {
	out := make([]string, len(toolNames))
	copy(out, toolNames)
	return out
}

func RegisterTools(s *mcp.Server, app *App) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "calculate_firebird_config",
		Description: "Call the IBSurgeon Configuration Calculator API and return optimized firebird.conf and databases.conf content.",
	}, app.CalculateFirebirdConfig)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "write_firebird_configs",
		Description: "Write firebird.conf and databases.conf to disk with timestamped backups before overwrite. Use dry_run to preview paths.",
	}, app.WriteFirebirdConfigs)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_server_info",
		Description: "Return MCP server version, available tools, API version, and update notification if a newer release exists.",
	}, app.GetServerInfo)
}

func Run(ctx context.Context, app *App) error {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    serverName,
		Version: "v" + app.Version,
	}, nil)
	RegisterTools(s, app)
	return s.Run(ctx, &mcp.StdioTransport{})
}
