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
# Priority: 1) server_url env var (lowercase, highest priority), 2) SERVER_URL env var (uppercase, backward compatibility), 3) Injected SERVER_URL, 4) Default localhost
# The server replaces __SERVER_URL__ with the actual server URL when serving the script
# Note: Server uses strings.ReplaceAll, so ALL occurrences of __SERVER_URL__ are replaced
# We use a workaround: check if value looks like a URL (contains ://) to detect server injection
if ($env:server_url) {
    # server_url env var is set (lowercase, highest priority)
    $SERVER_URL = $env:server_url
} elseif ($env:SERVER_URL) {
    # SERVER_URL env var is set (uppercase, backward compatibility)
    $SERVER_URL = $env:SERVER_URL
} else {
    # Neither env var is set, use default value
    # Server will replace __SERVER_URL__ with actual URL (e.g., http://packages.slileisure.com)
    $SERVER_URL = "__SERVER_URL__"
    # If SERVER_URL_ENV is set, use it
    if ($env:SERVER_URL_ENV) {
        $SERVER_URL = $env:SERVER_URL_ENV
    # Check if value looks like a URL (contains ://) - means server injected it
    } elseif ($SERVER_URL -match "://") {
        # SERVER_URL contains ://, so it was replaced by server, use it as-is
        # $SERVER_URL already set to injected value
    } else {
        # SERVER_URL doesn't look like a URL, still contains placeholder, use localhost
        $SERVER_URL = "http://localhost:8080"
    }
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

# Download binary
function Download-Binary {
    param(
        [string]$Url,
        [string]$OutputPath
    )
    
    try {
        # Use System.Net.WebClient for reliable binary file download
        # WebClient.DownloadFile will throw an exception for non-2xx status codes
        $webClient = New-Object System.Net.WebClient
        try {
            # Download file (will throw exception if status code is not 2xx)
            $webClient.DownloadFile($Url, $OutputPath)
            
            # Verify file was created and has reasonable size (at least 1KB for a binary)
            if (-not (Test-Path $OutputPath) -or (Get-Item $OutputPath).Length -lt 1024) {
                Write-Host "Error: Downloaded file is too small or empty (may be an error response)" -ForegroundColor Red
                if (Test-Path $OutputPath) {
                    # Check if it's a JSON error
                    $firstLine = Get-Content $OutputPath -TotalCount 1 -ErrorAction SilentlyContinue
                    if ($firstLine -and $firstLine.Trim().StartsWith("{")) {
                        Write-Host "Error: Server returned JSON error instead of binary file" -ForegroundColor Red
                        Write-Host "Response content:"
                        Get-Content $OutputPath -TotalCount 10 -ErrorAction SilentlyContinue
                    }
                    Remove-Item $OutputPath -Force
                }
                exit 1
            }
            
        } finally {
            $webClient.Dispose()
        }
        
    } catch [System.Net.WebException] {
        # Handle HTTP errors (404, 500, etc.)
        $statusCode = 0
        $errorContent = ""
        
        if ($_.Exception.Response) {
            $statusCode = [int]$_.Exception.Response.StatusCode
            Write-Host "Error: Download failed with HTTP status $statusCode" -ForegroundColor Red
            
            # Try to read error response body
            try {
                $errorStream = $_.Exception.Response.GetResponseStream()
                $reader = New-Object System.IO.StreamReader($errorStream)
                $errorContent = $reader.ReadToEnd()
                $reader.Close()
                $errorStream.Close()
                
                if ($errorContent) {
                    Write-Host "Server returned error:" -ForegroundColor Red
                    Write-Host $errorContent
                }
            } catch {
                # Ignore errors reading error response
            }
        } else {
            Write-Host "Error: Failed to download binary: $_" -ForegroundColor Red
        }
        
        # Clean up any partial download
        if (Test-Path $OutputPath) {
            Remove-Item $OutputPath -Force -ErrorAction SilentlyContinue
        }
        exit 1
        
    } catch {
        Write-Host "Error: Failed to download binary: $_" -ForegroundColor Red
        
        # Clean up any partial download
        if (Test-Path $OutputPath) {
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
    
    # Construct binary filename directly based on platform
    $filename = "kkartifact-agent-${os}-${arch}.exe"
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
    
    # Install (move to target location) - force overwrite existing version
    if (Test-Path $installPath) {
        Remove-Item $installPath -Force -ErrorAction SilentlyContinue
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
