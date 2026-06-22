FROM node:20-alpine AS build

WORKDIR /app

# Install dependencies for download
RUN apk add --no-cache curl tar

# Download the latest release of the MCP server (Linux amd64)
RUN LATEST=$(curl -s https://api.github.com/repos/IBSurgeon/FirebirdConfCalcMCP/releases/latest) && \
    URL=$(echo "$LATEST" | grep browser_download_url | grep linux_amd64 | head -1 | cut -d'"' -f4) && \
    curl -sL -o /tmp/fcc-mcp.tar.gz "$URL" && \
    tar xzf /tmp/fcc-mcp.tar.gz -C /tmp && \
    mv /tmp/firebird-conf-calc-mcp /app/firebird-conf-calc-mcp && \
    chmod +x /app/firebird-conf-calc-mcp

# ========== Runtime ==========
FROM node:20-alpine

WORKDIR /app

# Copy the MCP server binary
COPY --from=build /app/firebird-conf-calc-mcp /usr/local/bin/firebird-conf-calc-mcp

# Install supergateway for Streamable HTTP transport
RUN npm install -g supergateway

# Create directory for credentials
RUN mkdir -p /app/credentials

# Entry point script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8000

HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
  CMD wget -qO- http://localhost:8000/mcp || exit 1

ENTRYPOINT ["/entrypoint.sh"]
