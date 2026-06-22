#!/bin/sh
# Entry point for the Firebird Configuration Calculator MCP Docker container
#
# Expects credentials in one of:
#   - /app/credentials/password_api.txt (mounted volume or pre-populated)
#   - CC_EMAIL and CC_PASSWORD environment variables

set -e

# Write credentials file from environment variables if password_api.txt is empty/missing
if [ ! -s /app/credentials/password_api.txt ]; then
    if [ -n "$CC_EMAIL" ] && [ -n "$CC_PASSWORD" ]; then
        echo "user: $CC_EMAIL" > /app/credentials/password_api.txt
        echo "password: $CC_PASSWORD" >> /app/credentials/password_api.txt
        echo "entrypoint: wrote credentials from env vars"
    else
        echo "WARNING: No credentials found."
        echo "Provide CC_EMAIL + CC_PASSWORD env vars or mount password_api.txt"
    fi
fi

# Start supergateway wrapping the MCP server via stdio
exec npx supergateway \
    --stdio "/usr/local/bin/firebird-conf-calc-mcp --credentials /app/credentials/password_api.txt" \
    --port "${PORT:-8000}"
