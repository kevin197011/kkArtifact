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
    
    local http_code=""
    local temp_output=$(mktemp)
    trap "rm -f ${temp_output}" EXIT
    
    if command -v curl >/dev/null 2>&1; then
        # Use -w to write HTTP code to stderr, download to temp file first
        http_code=$(curl -L -w "%{http_code}" -o "${temp_output}" -s "${url}" 2>&1 | tail -n1)
        # Move temp file to output if successful
        if [ "${http_code}" = "200" ]; then
            mv "${temp_output}" "${output}"
        else
            mv "${temp_output}" "${output}" 2>/dev/null || true
        fi
    elif command -v wget >/dev/null 2>&1; then
        if wget -O "${temp_output}" "${url}" 2>&1 | grep -q "200 OK"; then
            http_code="200"
            mv "${temp_output}" "${output}"
        else
            http_code="000"
            mv "${temp_output}" "${output}" 2>/dev/null || true
        fi
    else
        echo "Error: curl or wget is required but not found" >&2
        exit 1
    fi
    
    # Check HTTP status code
    if [ "${http_code}" != "200" ]; then
        echo -e "${RED}Error: Download failed with HTTP status ${http_code}${NC}" >&2
        # Check if output is JSON error response
        if [ -f "${output}" ] && head -n1 "${output}" 2>/dev/null | grep -q "^{"; then
            echo "Server returned error:" >&2
            cat "${output}" >&2
            echo "" >&2
        fi
        rm -f "${output}"
        exit 1
    fi
    
    # Verify downloaded file is not JSON (error response)
    if [ -f "${output}" ] && head -n1 "${output}" 2>/dev/null | grep -q "^{"; then
        echo -e "${RED}Error: Server returned JSON error instead of binary file${NC}" >&2
        echo "Response content:" >&2
        head -n5 "${output}" >&2
        rm -f "${output}"
        exit 1
    fi
    
    # Verify file is not empty and has reasonable size (at least 1KB for a binary)
    if [ ! -s "${output}" ] || [ $(stat -f%z "${output}" 2>/dev/null || stat -c%s "${output}" 2>/dev/null || echo 0) -lt 1024 ]; then
        echo -e "${RED}Error: Downloaded file is too small or empty (may be an error response)${NC}" >&2
        if [ -f "${output}" ]; then
            echo "File content (first 10 lines):" >&2
            head -n10 "${output}" >&2
        fi
        rm -f "${output}"
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
        # Check if server_url needs to be updated (if it's localhost or placeholder)
        local current_url=$(grep -E "^server_url:" "${config_file}" 2>/dev/null | sed 's/^server_url:[[:space:]]*//' | tr -d '"' || echo "")
        if [ -z "${current_url}" ] || [ "${current_url}" = "http://localhost:8080" ] || [ "${current_url}" = "__SERVER_URL__" ]; then
            # Update server_url if it's using default/placeholder value
            if command -v sed >/dev/null 2>&1; then
                # Use sed to update server_url line
                if grep -q "^server_url:" "${config_file}"; then
                    sed -i "s|^server_url:.*|server_url: ${SERVER_URL}|" "${config_file}"
                    echo -e "${GREEN}✓ Updated server_url in global config file: ${config_file}${NC}"
                    echo "  Updated to: ${SERVER_URL}"
                else
                    # Add server_url if it doesn't exist
                    sed -i "/^# kkArtifact Agent Global Configuration/a\\server_url: ${SERVER_URL}" "${config_file}"
                    echo -e "${GREEN}✓ Added server_url to global config file: ${config_file}${NC}"
                fi
            else
                echo -e "${YELLOW}Global config file exists but server_url may need updating${NC}"
                echo "  Current: ${current_url:-not set}"
                echo "  Should be: ${SERVER_URL}"
                echo "  Please update manually if needed."
            fi
        else
            echo -e "${YELLOW}Global config file already exists: ${config_file}${NC}"
            echo "  Current server_url: ${current_url}"
            echo "  Skipping update to preserve existing settings."
        fi
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
