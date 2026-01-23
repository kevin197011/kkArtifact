# Artifact Agent Specification

## MODIFIED Requirements

### Requirement: Agent Installation Scripts

The system SHALL provide platform-specific installation scripts that automate the download, installation, and configuration of the agent binary.

#### Scenario: Unix install script execution
- **WHEN** a user executes `install-agent.sh` on a Unix-like system (Linux, macOS)
- **THEN** the script detects the platform and architecture
- **AND** downloads the correct agent binary from the server
- **AND** sets executable permissions on the binary
- **AND** installs the binary to `/usr/local/bin/kkartifact-agent` (system-wide) or `~/.local/bin/kkartifact-agent` (user-local if no sudo)
- **AND** verifies the installation was successful
- **AND** displays usage instructions

#### Scenario: Windows install script execution
- **WHEN** a user executes `install-agent.ps1` on Windows
- **THEN** the script detects the platform and architecture
- **AND** downloads the correct agent binary from the server
- **AND** installs the binary to `%LOCALAPPDATA%\kkartifact\kkartifact-agent.exe`
- **AND** verifies the installation was successful
- **AND** displays usage instructions

#### Scenario: Script platform detection
- **WHEN** an install script is executed
- **THEN** the script correctly identifies the operating system (Linux, macOS, Windows)
- **AND** correctly identifies the architecture (amd64, arm64)
- **AND** selects the appropriate binary filename for download
- **AND** handles detection failures with clear error messages

#### Scenario: Installation path selection
- **WHEN** the Unix install script is executed
- **THEN** the script attempts system-wide installation to `/usr/local/bin/` first
- **AND** if system-wide installation fails (no sudo), falls back to user-local installation to `~/.local/bin/`
- **AND** provides clear feedback about the installation location

#### Scenario: Script download endpoint
- **WHEN** a client requests `/api/v1/downloads/scripts/install-agent.sh` or `/api/v1/downloads/scripts/install-agent.ps1`
- **THEN** the server returns the appropriate install script file
- **AND** the response includes correct Content-Type headers
- **AND** the script file is served with appropriate security headers

#### Scenario: Server URL environment variable override (Unix)
- **WHEN** a user executes `curl ... | server_url="https://custom-endpoint.com" bash` on Unix
- **THEN** the script uses `server_url` environment variable value as the server URL
- **AND** the `server_url` value takes precedence over server-injected URL
- **AND** the `server_url` value takes precedence over `SERVER_URL` environment variable
- **AND** the configuration file is created with the `server_url` value

#### Scenario: Server URL environment variable override (Windows)
- **WHEN** a user sets `$env:server_url="https://custom-endpoint.com"` and executes the PowerShell script on Windows
- **THEN** the script uses `$env:server_url` environment variable value as the server URL
- **AND** the `$env:server_url` value takes precedence over server-injected URL
- **AND** the `$env:server_url` value takes precedence over `$env:SERVER_URL` environment variable
- **AND** the configuration file is created with the `$env:server_url` value

#### Scenario: Server URL priority order
- **WHEN** an install script is executed
- **THEN** the script resolves server URL in the following priority order:
  1. `server_url` environment variable (lowercase, highest priority)
  2. `SERVER_URL` environment variable (uppercase, for backward compatibility)
  3. Server-injected URL (from script download endpoint)
  4. Default `http://localhost:8080` (fallback)
- **AND** the first available value in this priority order is used

#### Scenario: Backward compatibility with SERVER_URL
- **WHEN** a user sets `SERVER_URL="https://custom-endpoint.com"` (uppercase) and executes the install script
- **THEN** the script uses `SERVER_URL` environment variable value as the server URL
- **AND** the configuration file is created with the `SERVER_URL` value
- **AND** the behavior matches previous implementation for backward compatibility
