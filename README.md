# Firebird Configuration Calculator MCP

MCP server for the [IBSurgeon Configuration Calculator for Firebird](https://cc.ib-aid.com/). Generates optimized `firebird.conf` and `databases.conf` files via API and can write them to disk with timestamped backups.

Repository: [github.com/IBSurgeon/FirebirdConfCalcMCP](https://github.com/IBSurgeon/FirebirdConfCalcMCP)  
Releases: [github.com/IBSurgeon/FirebirdConfCalcMCP/releases](https://github.com/IBSurgeon/FirebirdConfCalcMCP/releases)

## Prerequisites

- Go **1.25+** (only if building from source)
- A registered account at [cc.ib-aid.com](https://cc.ib-aid.com/) with an API password (widget: *Get Firebird configuration via API*)
- Credentials file (local only, never commit to git):

```
user: your@email.com
password: your_api_password
```

Accepted key aliases: `user` / `mailLogin` / `login` / `email` and `password` / `passApi` / `pass`.

## Installation

### Windows

1. Download `firebird-conf-calc-mcp_*_windows_amd64.zip` from [Releases](https://github.com/IBSurgeon/FirebirdConfCalcMCP/releases)
2. Extract to e.g. `C:\Tools\firebird-conf-calc-mcp\`
3. Create `password_api.txt` in the same folder (or pass `--credentials` path)

### Linux

```bash
curl -sL -o /tmp/fcc-mcp.tar.gz \
  "$(curl -s https://api.github.com/repos/IBSurgeon/FirebirdConfCalcMCP/releases/latest \
    | grep browser_download_url | grep linux_amd64 | cut -d '"' -f 4)"
tar xzf /tmp/fcc-mcp.tar.gz -C /tmp
sudo mv /tmp/firebird-conf-calc-mcp /usr/local/bin/
chmod +x /usr/local/bin/firebird-conf-calc-mcp
```

### macOS

```bash
curl -sL -o /tmp/fcc-mcp.tar.gz \
  "$(curl -s https://api.github.com/repos/IBSurgeon/FirebirdConfCalcMCP/releases/latest \
    | grep browser_download_url | grep darwin_arm64 | cut -d '"' -f 4)"
tar xzf /tmp/fcc-mcp.tar.gz -C /tmp
sudo mv /tmp/firebird-conf-calc-mcp /usr/local/bin/
```

### From source

```bash
git clone https://github.com/IBSurgeon/FirebirdConfCalcMCP.git
cd FirebirdConfCalcMCP
go test ./...
make build
```

Or:

```bash
go install github.com/IBSurgeon/FirebirdConfCalcMCP/cmd/firebird-conf-calc-mcp@latest
```

## AI client setup

This MCP server uses **stdio** transport: the client launches `firebird-conf-calc-mcp` as a local subprocess. That works out of the box in desktop/IDE clients listed below as **Direct (stdio)**.

**ChatGPT** and **Grok** require a **public HTTPS** MCP endpoint (HTTP/SSE). This release is stdio-only, so use a bridge (see [Remote clients](#remote-clients-chatgpt-grok)) or an MCP-capable desktop client instead.

### Shared configuration

Use absolute paths to the binary and `password_api.txt`. Example for Windows:

```json
{
  "mcpServers": {
    "firebird-config-calculator": {
      "command": "C:\\Tools\\firebird-conf-calc-mcp\\firebird-conf-calc-mcp.exe",
      "args": [
        "--credentials",
        "C:\\Tools\\firebird-conf-calc-mcp\\password_api.txt"
      ]
    }
  }
}
```

Linux / macOS — set `command` to `/usr/local/bin/firebird-conf-calc-mcp` and an absolute credentials path.

After editing config files, restart the client (or reload MCP servers where supported).

### Client compatibility

| Client | Direct (stdio) | Config location | Notes |
|--------|----------------|-----------------|-------|
| [Cursor](https://cursor.com) | Yes | `.cursor/mcp.json` or **Settings → MCP** | Hot reload supported |
| [Claude Desktop](https://claude.ai/download) | Yes | See below | Full quit required after config change |
| [ChatGPT](https://chatgpt.com) | No (remote only) | **Settings → Connectors** | Developer Mode + public URL |
| [Grok](https://grok.com) | No (remote only) | [grok.com/connectors](https://grok.com/connectors) | Custom connector + public HTTPS |
| [DeepSeek](https://www.deepseek.com) | Via other clients | — | No built-in custom MCP UI; use Cursor/Claude/Cline below |
| [Qwen Code](https://github.com/QwenLM/qwen-code) | Yes | `~/.qwen/settings.json` | `qwen mcp add` or manual `mcpServers` |
| [Qwen-Agent](https://github.com/QwenLM/Qwen-Agent) | Yes | Python `tools` config | For custom agents / scripts |

---

### Cursor

Add to `.cursor/mcp.json` or **Settings → MCP**:

```json
{
  "mcpServers": {
    "firebird-config-calculator": {
      "command": "C:\\Tools\\firebird-conf-calc-mcp\\firebird-conf-calc-mcp.exe",
      "args": [
        "--credentials",
        "C:\\Tools\\firebird-conf-calc-mcp\\password_api.txt"
      ]
    }
  }
}
```

Reload MCP servers or restart Cursor. Confirm three tools are listed.

---

### Claude Desktop

1. Install [Claude Desktop](https://claude.ai/download).
2. Open the MCP config file:
   - **Windows:** `%APPDATA%\Claude\claude_desktop_config.json`
   - **macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`
   - **Linux:** `~/.config/Claude/claude_desktop_config.json`

   Shortcut: **Settings → Developer → Edit Config**.

3. Add under `mcpServers` (merge with existing entries if any):

```json
{
  "mcpServers": {
    "firebird-config-calculator": {
      "command": "C:\\Tools\\firebird-conf-calc-mcp\\firebird-conf-calc-mcp.exe",
      "args": [
        "--credentials",
        "C:\\Tools\\firebird-conf-calc-mcp\\password_api.txt"
      ]
    }
  }
}
```

4. **Fully quit** Claude Desktop (not just close the window) and reopen.
5. Check **Settings → Developer** — server should show as running with tools available in chat.

Docs: [Connect to local MCP servers](https://modelcontextprotocol.io/docs/develop/connect-local-servers)

---

### ChatGPT

ChatGPT custom connectors need a **remotely reachable** MCP server (Streamable HTTP or SSE), not a local stdio binary.

**If you have a public HTTPS endpoint** (self-hosted or tunneled):

1. **ChatGPT Plus / Pro / Team / Enterprise** — enable **Developer Mode**:
   - **Settings → Connectors → Advanced settings → Developer Mode** (on)
2. **Connectors → Create**:
   - **Name:** Firebird Config Calculator
   - **MCP server URL:** your public HTTPS MCP URL (e.g. `https://your-host/mcp`)
   - **Authentication:** as required by your deployment (OAuth is common for ChatGPT connectors)
3. In a chat, open **Tools / Connectors** and enable the connector.

**Local binary (this repo):** expose stdio via a bridge, for example [mcp-remote](https://www.npmjs.com/package/mcp-remote) or your own HTTP wrapper, then point ChatGPT at the tunnel URL. Stdio cannot be pasted directly into ChatGPT.

Docs: [OpenAI — custom MCP connectors](https://community.openai.com/t/how-to-set-up-a-remote-mcp-server-and-connect-it-to-chatgpt-deep-research/1278375)

---

### Grok

Grok **Bring Your Own MCP** also requires a **public HTTPS** server (Streamable HTTP or SSE).

1. Open [grok.com/connectors](https://grok.com/connectors).
2. **New Connector → Custom** (or **Other** on Business/Enterprise admin console).
3. Enter your MCP server URL and complete authentication.
4. Enable the connector in a Grok chat.

**Local binary:** run a tunnel (e.g. ngrok, Cloudflare Tunnel) to an HTTP/SSE front-end for this server, or host the MCP service on a VM with HTTPS.

Docs: [xAI — Custom MCP connectors](https://docs.x.ai/grok/connectors)

---

### DeepSeek

The [DeepSeek chat](https://chat.deepseek.com) web app does **not** currently offer a user-facing “add custom MCP server” setting like Cursor or Claude Desktop.

**Practical options:**

1. **Use an MCP client with this server** (recommended) — Cursor, Claude Desktop, Cline, Continue, etc. You can still select DeepSeek as the chat model in clients that support custom model endpoints, while tools come from this MCP server.
2. **Cline / VS Code** — add the shared `mcpServers` block to Cline MCP settings (`cline_mcp_settings.json`).
3. **Continue** — add a stdio server entry in `~/.continue/config.json` under `mcpServers`.

Example (Cline / compatible JSON clients):

```json
{
  "mcpServers": {
    "firebird-config-calculator": {
      "command": "/usr/local/bin/firebird-conf-calc-mcp",
      "args": ["--credentials", "/path/to/password_api.txt"]
    }
  }
}
```

---

### Qwen

**Qwen Code** (CLI / IDE agent) supports local stdio MCP servers.

**Option A — CLI:**

```bash
qwen mcp add firebird-config-calculator \
  /usr/local/bin/firebird-conf-calc-mcp \
  --args "--credentials" "/path/to/password_api.txt"
```

**Option B — `~/.qwen/settings.json`:**

```json
{
  "mcpServers": {
    "firebird-config-calculator": {
      "command": "/usr/local/bin/firebird-conf-calc-mcp",
      "args": [
        "--credentials",
        "/path/to/password_api.txt"
      ]
    }
  }
}
```

**Qwen-Agent** (Python agents, powers parts of [Qwen Chat](https://chat.qwen.ai) tooling):

```bash
pip install -U "qwen-agent[mcp]"
```

```python
from qwen_agent.agents import Assistant

tools = [{
    "mcpServers": {
        "firebird-config-calculator": {
            "command": "/usr/local/bin/firebird-conf-calc-mcp",
            "args": ["--credentials", "/path/to/password_api.txt"],
        }
    }
}]

bot = Assistant(llm=llm_cfg, function_list=tools)
```

Docs: [Qwen Code — MCP servers](https://qwenlm.github.io/qwen-code-docs/en/developers/tools/mcp-server/), [Qwen-Agent MCP guide](https://qwenlm-qwen-agent.mintlify.app/guides/mcp-integration)

---

### Remote clients (ChatGPT, Grok)

This v1 release ships **stdio only**. For ChatGPT or Grok you need either:

- A **hosted** build of this server behind HTTPS (future release), or
- A **local bridge** that runs `firebird-conf-calc-mcp` and exposes Streamable HTTP/SSE on a tunnel URL.

Until an official HTTP transport is added, the simplest path is **Cursor** or **Claude Desktop** on the same machine as the binary.


## Tools

| Tool | Description |
|------|-------------|
| `calculate_firebird_config` | Call Configuration Calculator API; return config text |
| `write_firebird_configs` | Save configs to disk with timestamped backups |
| `get_server_info` | Server version, tools list, update notification |

### Usage workflow

1. **Calculate** — call `calculate_firebird_config` with server parameters:

```json
{
  "server_version": "fb3",
  "server_architecture": "Classic",
  "cores": 8,
  "ram": 16,
  "count_users": 100,
  "page_size": 4096,
  "size_db": 100,
  "name_main_db": "mainDB",
  "path_to_main_db": "c:/data/main.fdb"
}
```

2. **Review** returned `firebird_conf` and `databases_conf`.

3. **Write** — call `write_firebird_configs` (use `dry_run: true` first to preview):

```json
{
  "output_dir": "C:/Program Files/Firebird/Firebird_3_0",
  "firebird_conf": "<from step 1>",
  "databases_conf": "<from step 1>",
  "dry_run": false
}
```

4. **Apply on server** ([IBSurgeon guide](https://ib-aid.com/en/articles/how-to-use-configuration-calculator-for-firebird)):

   - `gstat -h yourdb.fdb` — Page Buffers must be 0; run `gfix -buff 0` if not
   - Stop Firebird service
   - MCP creates timestamped backups (`firebird.conf.bak.YYYYMMDD-HHMMSS`) before overwrite
   - Start Firebird

### Parameter notes by version/architecture

Numeric tool parameters (`cores`, `ram`, `count_users`, `size_db`, `page_size`) must be **JSON integers**, not strings.

- **FB/HQ 2.5 Classic/SuperClassic**: `ram`, `count_users`, `page_size`
- **FB/HQ 2.5 SuperServer**: version and architecture only
- **FB/HQ 3.0–5.0 Classic/SuperClassic**: `ram`, `count_users`, `cores`, `page_size`
- **FB/HQ 3.0–5.0 SuperServer**: all parameters

## Updating

1. Download the latest release for your OS from [Releases](https://github.com/IBSurgeon/FirebirdConfCalcMCP/releases)
2. Replace the binary
3. Restart your MCP client (Cursor, Claude Desktop, etc.)
4. Call `get_server_info` — if outdated, response includes e.g. *"Current version 1.0.1. Newer version 1.2.0 exists — update to get more tools: …"*

Check version from CLI:

```bash
firebird-conf-calc-mcp --version
```

## Environment variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `CC_CREDENTIALS_FILE` | `password_api.txt` | Path to credentials file |
| `CC_ALLOW_ANY_OUTPUT_DIR` | unset | Set to `1` to allow writes with `..` in path |
| `CC_UPDATE_CHECK` | enabled | Set to `0` to disable remote version check |
| `CC_VERSION_MANIFEST_URL` | GitHub latest `version.json` | Override update manifest URL |

## Build

```bash
go test ./...
make build
```

Binary: `bin/firebird-conf-calc-mcp` (`.exe` on Windows)

Cross-platform release archives (Windows, Linux, macOS):

```bash
make build-all
# or
powershell -File scripts/build-all.ps1 -Version 1.0.1
```

Output: `dist/firebird-conf-calc-mcp_<version>_<os>_<arch>.zip` or `.tar.gz`

### End-to-end test

Runs the MCP server as a subprocess, connects via stdio, and calls the live Configuration Calculator API using [`password_api.txt`](password_api.txt):

```bash
make e2e
# or
go test -v -timeout 3m ./e2e/...
```

Skips automatically if `password_api.txt` is missing (e.g. in CI). Set `CC_E2E=0` to disable. Override credentials path with `CC_CREDENTIALS_FILE`.

## Troubleshooting

- **MCP not connecting**: use absolute paths; rebuild binary; restart client (see [AI client setup](#ai-client-setup))
- **API error**: verify credentials; check required parameters for your Firebird version ([API docs](https://ib-aid.com/api-to-create-firebird-configurations))
- **Type validation error on numeric fields**: pass `cores`, `ram`, `count_users`, `size_db`, and `page_size` as integers (e.g. `8`, not `"8"`)
- **Write rejected**: ensure `output_dir` exists; try `dry_run: true` first

## References

- [Configuration Calculator portal](https://cc.ib-aid.com/)
- [API documentation](https://ib-aid.com/api-to-create-firebird-configurations)
- [How to apply generated configs](https://ib-aid.com/en/articles/how-to-use-configuration-calculator-for-firebird)

## License

MIT — see [LICENSE](LICENSE). Copyright (c) 2026 IBSurgeon Ltd.
