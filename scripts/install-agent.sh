#!/bin/bash
# Copyright (c) 2025 kk
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT
#
# kkArtifact Agent Installation Script for Unix-like systems (Linux, macOS, BSD)
# This script automatically downloads and installs the kkartifact-agent binary

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Server URL is automatically injected by the server when serving this script
# This allows the script to work with simple "curl URL | bash" format
# Priority: 1) Injected SERVER_URL, 2) SERVER_URL env var, 3) Default
# SERVER_URL will be replaced by the server at runtime
SERVER_URL="${SERVER_URL:-__SERVER_URL__}"
# If still contains placeholder, try environment variable or default
if [ "$SERVER_URL" = "__SERVER_URL__" ]; then
    SERVER_URL="${SERVER_URL_ENV:-http://localhost:8080}"
fi

# Detect platform and architecture
detect_platform() {
    local os=""
    local arch=""
    
    # Detect OS
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        FreeBSD*)   os="freebsd" ;;
        OpenBSD*)   os="openbsd" ;;
        *)          os="unknown" ;;
    esac
    
    # Detect architecture
    case "$(uname -m)" in
        x86_64|amd64)   arch="amd64" ;;
        aarch64|arm64)  arch="arm64" ;;
        armv7l|armv6l) arch="arm" ;;
        *)              arch="unknown" ;;
    esac
    
    echo "${os}/${arch}"
}

# Get agent version info from server
get_version_info() {
    local version_url="${SERVER_URL}/api/v1/downloads/agent/version"
    if command -v curl >/dev/null 2>&1; then
        curl -s "${version_url}"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "${version_url}"
    else
        echo "Error: curl or wget is required but not found" >&2
        exit 1
    fi
}

# Download binary
download_binary() {
    local url="$1"
    local output="$2"
    
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "${output}" "${url}"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "${output}" "${url}"
    else
        echo "Error: curl or wget is required but not found" >&2
        exit 1
    fi
}

# Main installation function
main() {
    echo -e "${GREEN}kkArtifact Agent Installation Script${NC}"
    echo "=========================================="
    echo ""
    
    # Detect platform
    local platform=$(detect_platform)
    local os=$(echo "${platform}" | cut -d'/' -f1)
    local arch=$(echo "${platform}" | cut -d'/' -f2)
    
    if [ "${os}" = "unknown" ] || [ "${arch}" = "unknown" ]; then
        echo -e "${RED}Error: Unsupported platform: ${platform}${NC}" >&2
        echo "Please download the binary manually from the web UI." >&2
        exit 1
    fi
    
    echo -e "Detected platform: ${GREEN}${platform}${NC}"
    
    # Get version info
    echo "Fetching agent version information..."
    local version_info=$(get_version_info)
    if [ -z "${version_info}" ]; then
        echo -e "${RED}Error: Failed to fetch version information from server${NC}" >&2
        echo "Please check that the server is running at ${SERVER_URL}" >&2
        exit 1
    fi
    
    # Extract binary filename for this platform
    local filename=""
    if command -v jq >/dev/null 2>&1; then
        # Check if binaries array exists and is not empty
        local binaries_count=$(echo "${version_info}" | jq -r '.binaries | if type == "array" then length else 0 end' 2>/dev/null || echo "0")
        if [ "${binaries_count}" -gt 0 ] 2>/dev/null; then
            filename=$(echo "${version_info}" | jq -r ".binaries[]? | select(.platform == \"${platform}\") | .filename" 2>/dev/null | head -n1)
        fi
        # If jq returned null or empty, fallback to manual construction
        if [ -z "${filename}" ] || [ "${filename}" = "null" ] || [ "${filename}" = "" ]; then
            filename=""
        fi
    elif command -v python3 >/dev/null 2>&1; then
        filename=$(echo "${version_info}" | python3 -c "import sys, json; data=json.load(sys.stdin); binaries = data.get('binaries', []); print(next((b['filename'] for b in binaries if b.get('platform') == '${platform}'), ''))" 2>/dev/null || echo "")
    fi
    
    # Fallback: construct filename manually if not found
    if [ -z "${filename}" ] || [ "${filename}" = "null" ] || [ "${filename}" = "" ]; then
        filename="kkartifact-agent-${os}-${arch}"
        echo -e "${YELLOW}Warning: Version info not available or binaries array is empty, using default filename: ${filename}${NC}"
    fi
    
    echo -e "Target binary: ${GREEN}${filename}${NC}"
    
    # Determine installation path
    local install_path=""
    local install_dir=""
    
    # Try system-wide installation first
    if [ -w "/usr/local/bin" ]; then
        install_dir="/usr/local/bin"
        install_path="${install_dir}/kkartifact-agent"
        echo "Installing to system-wide location: ${install_path}"
    else
        # Fallback to user-local installation
        install_dir="${HOME}/.local/bin"
        install_path="${install_dir}/kkartifact-agent"
        mkdir -p "${install_dir}"
        echo "Installing to user-local location: ${install_path}"
        echo -e "${YELLOW}Note: Make sure ${install_dir} is in your PATH${NC}"
    fi
    
    # Create temporary file for download
    local temp_file=$(mktemp)
    trap "rm -f ${temp_file}" EXIT
    
    # Download binary
    local download_url="${SERVER_URL}/api/v1/downloads/agent/${filename}"
    echo "Downloading ${filename}..."
    download_binary "${download_url}" "${temp_file}"
    
    # Verify download
    if [ ! -f "${temp_file}" ] || [ ! -s "${temp_file}" ]; then
        echo -e "${RED}Error: Download failed or file is empty${NC}" >&2
        exit 1
    fi
    
    # Make executable
    chmod +x "${temp_file}"
    
    # Install (move to target location)
    if [ -f "${install_path}" ]; then
        echo -e "${YELLOW}Warning: ${install_path} already exists. Overwriting...${NC}"
    fi
    
    mv "${temp_file}" "${install_path}"
    
    # Verify installation
    if [ ! -f "${install_path}" ] || [ ! -x "${install_path}" ]; then
        echo -e "${RED}Error: Installation failed${NC}" >&2
        exit 1
    fi
    
    # Get version
    local agent_version=$("${install_path}" version 2>/dev/null | head -n1 || echo "unknown")
    
    # Create global configuration file
    create_global_config
    
    echo ""
    echo -e "${GREEN}✓ Installation successful!${NC}"
    echo ""
    echo "Installed: ${install_path}"
    echo "Version: ${agent_version}"
    echo ""
    echo "Usage:"
    echo "  ${install_path} --help"
    echo ""
    echo "To use globally, make sure the installation directory is in your PATH:"
    if [ "${install_dir}" = "/usr/local/bin" ]; then
        echo "  (Already in standard PATH)"
    else
        echo "  Add to ~/.bashrc or ~/.zshrc:"
        echo "    export PATH=\"\${HOME}/.local/bin:\${PATH}\""
    fi
}

# Create global configuration file
create_global_config() {
    local config_dir="/etc/kkArtifact"
    local config_file="${config_dir}/config.yml"
    
    # Check if we can write to /etc (requires root/sudo)
    if [ ! -w "/etc" ]; then
        echo -e "${YELLOW}Note: Cannot create global config at ${config_file} (requires root privileges)${NC}"
        echo "You can create it manually later with:"
        echo "  sudo mkdir -p ${config_dir}"
        echo "  sudo tee ${config_file} > /dev/null <<EOF"
        echo "server_url: ${SERVER_URL}"
        echo "EOF"
        return 0
    fi
    
    # Create config directory if it doesn't exist
    if [ ! -d "${config_dir}" ]; then
        mkdir -p "${config_dir}"
        echo "Created global config directory: ${config_dir}"
    fi
    
    # Create or update config file
    if [ -f "${config_file}" ]; then
        echo -e "${YELLOW}Global config file already exists: ${config_file}${NC}"
        echo "Skipping config file creation to preserve existing settings."
    else
        cat > "${config_file}" <<EOF
# kkArtifact Agent Global Configuration
# This file is automatically created by the installation script
# You can modify this file to change global settings

server_url: ${SERVER_URL}
# token: YOUR_TOKEN_HERE  # Uncomment and set your token here
# concurrency: 50          # Number of concurrent uploads/downloads (default: 50)
# ignore: []               # Global ignore patterns
EOF
        chmod 644 "${config_file}"
        echo -e "${GREEN}✓ Created global config file: ${config_file}${NC}"
        echo "  You can edit this file to set your token and other global settings."
    fi
}

# Run main function
main "$@"
