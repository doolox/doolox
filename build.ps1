#Requires -Version 5.1
<#
.SYNOPSIS
    Build Doolox for Windows (native build, no `make` required).

.DESCRIPTION
    Requires Go 1.26+.

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

New-Item -ItemType Directory -Force -Path "bin" | Out-Null

Write-Host "Building Doolox for windows/amd64..."
$env:CGO_ENABLED = "0"
go build -trimpath -ldflags "-s -w -H=windowsgui" -o "bin\doolox-windows-amd64.exe" .
if ($LASTEXITCODE -ne 0) { throw "Windows build failed" }

Write-Host "Done. Output: .\bin\doolox-windows-amd64.exe"
