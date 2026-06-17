BINARY_NAME := firebird-conf-calc-mcp
MODULE      := github.com/IBSurgeon/FirebirdConfCalcMCP
VERSION     ?= $(shell git describe --tags --always --dirty 2>nul || echo 1.0.0-dev)
COMMIT      ?= $(shell git rev-parse --short HEAD 2>nul || echo unknown)
DATE        ?= $(shell powershell -NoProfile -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ'" 2>nul || date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -s -w \
	-X $(MODULE)/internal/version.Version=$(VERSION) \
	-X $(MODULE)/internal/version.Commit=$(COMMIT) \
	-X $(MODULE)/internal/version.Date=$(DATE)

.PHONY: build build-all test clean release-snapshot e2e

build:
	@mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) ./cmd/firebird-conf-calc-mcp

build-all:
	@powershell -NoProfile -ExecutionPolicy Bypass -File scripts/build-all.ps1 -Version $(VERSION)

test:
	go test ./...

e2e:
	go test -v -timeout 3m ./e2e/...

clean:
	rm -rf bin dist

release-snapshot:
	goreleaser build --snapshot --clean
