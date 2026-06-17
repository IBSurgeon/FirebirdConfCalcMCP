BINARY_NAME := firebird-conf-calc-mcp
MODULE      := github.com/IBSurgeon/FirebirdConfCalcMCP
VERSION     ?= $(shell git describe --tags --always --dirty 2>nul || echo 1.0.0-dev)
COMMIT      ?= $(shell git rev-parse --short HEAD 2>nul || echo unknown)
DATE        ?= $(shell powershell -NoProfile -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ'" 2>nul || date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -s -w \
	-X $(MODULE)/internal/version.Version=$(VERSION) \
	-X $(MODULE)/internal/version.Commit=$(COMMIT) \
	-X $(MODULE)/internal/version.Date=$(DATE)

.PHONY: build test clean release-snapshot

build:
	@mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) ./cmd/firebird-conf-calc-mcp

test:
	go test ./...

clean:
	rm -rf bin dist

release-snapshot:
	goreleaser build --snapshot --clean
