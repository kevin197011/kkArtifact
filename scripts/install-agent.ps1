# Copyright (c) 2025 kk
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT
#
# kkArtifact Agent Installation Script for Windows
# This script automatically downloads and installs the kkartifact-agent binary

#Requires -Version 5.1

$ErrorActionPreference = "Stop"

# Server URL is automatically injected by the server when serving this script
# This allows the script to work with simple "irm URL | iex" format
# Priority: 1) Injected SERVER_URL, 2) SERVER_URL env var, 3) Default
# SERVER_URL will be replaced by the server at runtime
$SERVER_URL = if ($env:SERVER_URL) { $env:SERVER_URL } else { "__SERVER_URL__" }
# If still contains placeholder, try environment variable or default
if ($SERVER_URL -eq "__SERVER_URL__") {
    $SERVER_URL = if ($env:SERVER_URL_ENV) { $env:SERVER_URL_ENV } else { "http://localhost:8080" }
}

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
        $response = Invoke-WebRequest -Uri $Url -OutFile $OutputPath -ErrorAction Stop
        
        # Verify HTTP status code
        if ($response.StatusCode -ne 200) {
            Write-Host "Error: Download failed with HTTP status $($response.StatusCode)" -ForegroundColor Red
            if (Test-Path $OutputPath) {
                Remove-Item $OutputPath -Force
            }
            exit 1
        }
        
        # Verify downloaded file is not JSON (error response)
        if (Test-Path $OutputPath) {
            $firstLine = Get-Content $OutputPath -TotalCount 1 -ErrorAction SilentlyContinue
            if ($firstLine -and $firstLine.Trim().StartsWith("{")) {
                Write-Host "Error: Server returned JSON error instead of binary file" -ForegroundColor Red
                Write-Host "Response content:"
                Get-Content $OutputPath -TotalCount 5
                Remove-Item $OutputPath -Force
                exit 1
            }
        }
    } catch {
        Write-Host "Error: Failed to download binary: $_" -ForegroundColor Red
        if (Test-Path $OutputPath) {
            $content = Get-Content $OutputPath -Raw -ErrorAction SilentlyContinue
            if ($content -and $content.Trim().StartsWith("{")) {
                Write-Host "Server returned error response:" -ForegroundColor Red
                Write-Host $content
            }
            Remove-Item $OutputPath -Force -ErrorAction SilentlyContinue
        }
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
    $filename = $null
    if ($versionInfo.binaries -and $versionInfo.binaries.Count -gt 0) {
        $binary = $versionInfo.binaries | Where-Object { $_.platform -eq $platform } | Select-Object -First 1
        if ($binary -and $binary.filename) {
            $filename = $binary.filename
        }
    }
    
    # Fallback: construct filename manually if not found
    if (-not $filename) {
        $filename = "kkartifact-agent-${os}-${arch}.exe"
        Write-Host "Warning: Version info not available or binaries array is empty, using default filename: $filename" -ForegroundColor Yellow
    }
    
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
    
    # Create global configuration file
    New-GlobalConfig
    
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

# Create global configuration file
function New-GlobalConfig {
    # Windows global config location: C:\ProgramData\kkArtifact\config.yml
    $configDir = Join-Path $env:ProgramData "kkArtifact"
    $configFile = Join-Path $configDir "config.yml"
    
    # Check if we can write to ProgramData (requires admin privileges)
    try {
        if (-not (Test-Path $configDir)) {
            New-Item -ItemType Directory -Path $configDir -Force | Out-Null
            Write-Host "Created global config directory: $configDir"
        }
        
        # Create or update config file
        if (Test-Path $configFile) {
            # Check if server_url needs to be updated (if it's localhost or placeholder)
            $content = Get-Content $configFile -Raw -ErrorAction SilentlyContinue
            $currentUrl = ""
            if ($content) {
                $match = [regex]::Match($content, '(?m)^server_url:\s*(.+)$')
                if ($match.Success) {
                    $currentUrl = $match.Groups[1].Value.Trim()
                }
            }
            
            if ([string]::IsNullOrEmpty($currentUrl) -or $currentUrl -eq "http://localhost:8080" -or $currentUrl -eq "__SERVER_URL__") {
                # Update server_url if it's using default/placeholder value
                $newContent = $content -replace '(?m)^server_url:\s*.+$', "server_url: $SERVER_URL"
                if ($newContent -ne $content) {
                    Set-Content -Path $configFile -Value $newContent -Encoding UTF8
                    Write-Host "✓ Updated server_url in global config file: $configFile" -ForegroundColor Green
                    Write-Host "  Updated to: $SERVER_URL"
                } else {
                    # Add server_url if it doesn't exist
                    $newContent = $content + "`nserver_url: $SERVER_URL`n"
                    Set-Content -Path $configFile -Value $newContent -Encoding UTF8
                    Write-Host "✓ Added server_url to global config file: $configFile" -ForegroundColor Green
                }
            } else {
                Write-Host "Global config file already exists: $configFile" -ForegroundColor Yellow
                Write-Host "  Current server_url: $currentUrl"
                Write-Host "  Skipping update to preserve existing settings."
            }
        } else {
            $configContent = @"
# kkArtifact Agent Global Configuration
# This file is automatically created by the installation script
# You can modify this file to change global settings

server_url: $SERVER_URL
# token: YOUR_TOKEN_HERE  # Uncomment and set your token here
# concurrency: 50          # Number of concurrent uploads/downloads (default: 50)
# ignore: []               # Global ignore patterns
"@
            Set-Content -Path $configFile -Value $configContent -Encoding UTF8
            Write-Host "✓ Created global config file: $configFile" -ForegroundColor Green
            Write-Host "  You can edit this file to set your token and other global settings."
        }
    } catch {
        Write-Host "Note: Cannot create global config at $configFile (requires admin privileges)" -ForegroundColor Yellow
        Write-Host "You can create it manually later with:"
        Write-Host "  New-Item -ItemType Directory -Path `"$configDir`" -Force"
        Write-Host "  @'"
        Write-Host "server_url: $SERVER_URL"
        Write-Host "'@ | Set-Content -Path `"$configFile`" -Encoding UTF8"
    }
}

# Run main function
Main
