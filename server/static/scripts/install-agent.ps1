# Copyright (c) 2025 kk
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT
#
# kkArtifact Agent Installation Script for Windows
# This script automatically downloads and installs the kkartifact-agent binary

#Requires -Version 5.1

$ErrorActionPreference = "Stop"

# Default server URL (can be overridden with SERVER_URL environment variable)
$SERVER_URL = if ($env:SERVER_URL) { $env:SERVER_URL } else { "http://localhost:8080" }

# Detect platform and architecture
function Get-Platform {
    $os = "windows"
    $arch = ""
    
    # Detect architecture
    switch ($env:PROCESSOR_ARCHITECTURE) {
        "AMD64" { $arch = "amd64" }
        "ARM64" { $arch = "arm64" }
        default {
            # Try alternative method
            $machine = (Get-WmiObject Win32_ComputerSystem).SystemType
            if ($machine -match "x64") {
                $arch = "amd64"
            } elseif ($machine -match "ARM") {
                $arch = "arm64"
            } else {
                $arch = "amd64"  # Default fallback
            }
        }
    }
    
    return "${os}/${arch}"
}

# Get agent version info from server
function Get-VersionInfo {
    $versionUrl = "${SERVER_URL}/api/v1/downloads/agent/version"
    
    try {
        $response = Invoke-RestMethod -Uri $versionUrl -Method Get -ErrorAction Stop
        return $response
    } catch {
        Write-Host "Error: Failed to fetch version information from server" -ForegroundColor Red
        Write-Host "Please check that the server is running at ${SERVER_URL}" -ForegroundColor Red
        exit 1
    }
}

# Download binary
function Download-Binary {
    param(
        [string]$Url,
        [string]$OutputPath
    )
    
    try {
        Invoke-WebRequest -Uri $Url -OutFile $OutputPath -ErrorAction Stop
    } catch {
        Write-Host "Error: Failed to download binary: $_" -ForegroundColor Red
        exit 1
    }
}

# Main installation function
function Main {
    Write-Host "kkArtifact Agent Installation Script" -ForegroundColor Green
    Write-Host "=====================================" -ForegroundColor Green
    Write-Host ""
    
    # Detect platform
    $platform = Get-Platform
    $os, $arch = $platform -split '/'
    
    Write-Host "Detected platform: $platform" -ForegroundColor Green
    Write-Host ""
    
    # Get version info
    Write-Host "Fetching agent version information..."
    $versionInfo = Get-VersionInfo
    
    # Extract binary filename for this platform
    $binary = $versionInfo.binaries | Where-Object { $_.platform -eq $platform } | Select-Object -First 1
    
    if (-not $binary) {
        Write-Host "Error: No binary available for platform $platform" -ForegroundColor Red
        exit 1
    }
    
    $filename = $binary.filename
    Write-Host "Target binary: $filename" -ForegroundColor Green
    Write-Host ""
    
    # Determine installation path
    $installDir = Join-Path $env:LOCALAPPDATA "kkartifact"
    $installPath = Join-Path $installDir "kkartifact-agent.exe"
    
    Write-Host "Installing to: $installPath"
    
    # Create installation directory
    if (-not (Test-Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    }
    
    # Create temporary file for download
    $tempFile = Join-Path $env:TEMP "kkartifact-agent-install.tmp"
    
    # Download binary
    $downloadUrl = "${SERVER_URL}/api/v1/downloads/agent/${filename}"
    Write-Host "Downloading $filename..."
    Download-Binary -Url $downloadUrl -OutputPath $tempFile
    
    # Verify download
    if (-not (Test-Path $tempFile) -or (Get-Item $tempFile).Length -eq 0) {
        Write-Host "Error: Download failed or file is empty" -ForegroundColor Red
        exit 1
    }
    
    # Install (move to target location)
    if (Test-Path $installPath) {
        Write-Host "Warning: $installPath already exists. Overwriting..." -ForegroundColor Yellow
        Remove-Item $installPath -Force
    }
    
    Move-Item -Path $tempFile -Destination $installPath -Force
    
    # Verify installation
    if (-not (Test-Path $installPath)) {
        Write-Host "Error: Installation failed" -ForegroundColor Red
        exit 1
    }
    
    # Get version
    $agentVersion = "unknown"
    try {
        $versionOutput = & $installPath version 2>&1
        if ($versionOutput) {
            $agentVersion = ($versionOutput | Select-Object -First 1).ToString().Trim()
        }
    } catch {
        # Ignore version check errors
    }
    
    Write-Host ""
    Write-Host "✓ Installation successful!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Installed: $installPath"
    Write-Host "Version: $agentVersion"
    Write-Host ""
    Write-Host "Usage:"
    Write-Host "  $installPath --help"
    Write-Host ""
    Write-Host "Note: Add the installation directory to your PATH to use 'kkartifact-agent' command:"
    Write-Host "  [Environment]::SetEnvironmentVariable('Path', `$env:Path + ';$installDir', 'User')"
    Write-Host ""
}

# Run main function
Main
