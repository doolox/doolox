#Requires -Version 5.1
<#
.SYNOPSIS
    Build Doolox for Windows (native build, no `make` required).

.DESCRIPTION
    Requires Go 1.26+ and a C compiler (MinGW-w64 gcc), because Doolox links
    SQLite through cgo.

.EXAMPLE
    .\build.ps1
#>
[CmdletBinding()]
param()

$ErrorActionPreference = "Stop"
Set-Location (Split-Path -Parent $MyInvocation.MyCommand.Path)

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Error "Go is not installed or not on PATH. Install it from https://go.dev/dl/ and try again."
    exit 1
}

if (-not (Get-Command gcc -ErrorAction SilentlyContinue)) {
    Write-Error @"
No C compiler (gcc) found on PATH. Doolox links SQLite via cgo, so a C compiler
is required. Install MinGW-w64 (https://www.mingw-w64.org/) and add gcc to PATH.
"@
    exit 1
}

New-Item -ItemType Directory -Force -Path "bin" | Out-Null

Write-Host "Building Doolox for windows/amd64..."
$env:CGO_ENABLED = "1"
go build -trimpath -ldflags "-s -w -H=windowsgui" -o "bin\doolox-windows-amd64.exe" .
if ($LASTEXITCODE -ne 0) { throw "Windows build failed" }

Write-Host "Done. Output: .\bin\doolox-windows-amd64.exe"
