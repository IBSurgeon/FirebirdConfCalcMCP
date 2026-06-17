param(
    [string]$Version = "1.0.1"
)

$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$go = if ($env:GO_BIN) { $env:GO_BIN } else { "go" }
$commit = (& git -C $root rev-parse --short HEAD 2>$null)
if (-not $commit) { $commit = "unknown" }
$date = (Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ")
$ldflags = "-s -w -X github.com/IBSurgeon/FirebirdConfCalcMCP/internal/version.Version=$Version -X github.com/IBSurgeon/FirebirdConfCalcMCP/internal/version.Commit=$commit -X github.com/IBSurgeon/FirebirdConfCalcMCP/internal/version.Date=$date"

$dist = Join-Path $root "dist"
New-Item -ItemType Directory -Force -Path $dist | Out-Null

$targets = @(
    @{ goos = "windows"; goarch = "amd64"; name = "firebird-conf-calc-mcp_${Version}_windows_amd64.zip" },
    @{ goos = "linux"; goarch = "amd64"; name = "firebird-conf-calc-mcp_${Version}_linux_amd64.tar.gz" },
    @{ goos = "linux"; goarch = "arm64"; name = "firebird-conf-calc-mcp_${Version}_linux_arm64.tar.gz" },
    @{ goos = "darwin"; goarch = "amd64"; name = "firebird-conf-calc-mcp_${Version}_darwin_amd64.tar.gz" },
    @{ goos = "darwin"; goarch = "arm64"; name = "firebird-conf-calc-mcp_${Version}_darwin_arm64.tar.gz" }
)

foreach ($t in $targets) {
    $env:GOOS = $t.goos
    $env:GOARCH = $t.goarch
    $env:CGO_ENABLED = "0"

    $stage = Join-Path $dist ("stage_" + $t.goos + "_" + $t.goarch)
    if (Test-Path $stage) { Remove-Item -Recurse -Force $stage }
    New-Item -ItemType Directory -Force -Path $stage | Out-Null

    $binName = "firebird-conf-calc-mcp"
    if ($t.goos -eq "windows") { $binName += ".exe" }
    $binPath = Join-Path $stage $binName

    Write-Host "Building $($t.goos)/$($t.goarch) -> $($t.name)"
    Push-Location $root
    & $go build -ldflags $ldflags -o $binPath ./cmd/firebird-conf-calc-mcp
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
    Pop-Location

    Copy-Item (Join-Path $root "LICENSE") $stage
    Copy-Item (Join-Path $root "README.md") $stage
    Copy-Item (Join-Path $root "version.json") $stage

    $archive = Join-Path $dist $t.name
    if (Test-Path $archive) { Remove-Item -Force $archive }

    if ($t.name.EndsWith(".zip")) {
        Compress-Archive -Path (Join-Path $stage "*") -DestinationPath $archive
    } else {
        tar -czf $archive -C $stage .
    }

    Remove-Item -Recurse -Force $stage
}

Write-Host ""
Write-Host "Built version ${Version}:"
Get-ChildItem $dist -Filter "firebird-conf-calc-mcp_$Version*" | Format-Table Name, Length, LastWriteTime
