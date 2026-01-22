# Artifact Web UI Specification

## MODIFIED Requirements

### Requirement: Agent Download Section

The frontend SHALL provide an agent installation section that allows users to install the agent client tool on their systems.

#### Scenario: Display install section with updated text
- **WHEN** the user views the inventory page
- **THEN** the section header displays "安装 agent 客户端工具" (Install Agent Client Tool)
- **AND** the version tag is not shown in the header (version info may be shown elsewhere)
- **AND** install script download buttons are available for Unix and Windows platforms

#### Scenario: Provide install script downloads
- **WHEN** the user views the agent installation section
- **THEN** install script download buttons are displayed for:
  - Unix-like systems (`install-agent.sh`)
  - Windows systems (`install-agent.ps1`)
- **AND** binary download buttons remain available for manual installation
- **AND** clicking install script buttons downloads the appropriate script

## ADDED Requirements

### Requirement: Install Script Download Links

The frontend SHALL provide download links for platform-specific install scripts that automate agent installation.

#### Scenario: Unix install script download
- **WHEN** the user clicks the Unix install script button
- **THEN** the `install-agent.sh` script is downloaded
- **AND** the script can be executed to install the agent

#### Scenario: Windows install script download
- **WHEN** the user clicks the Windows install script button
- **THEN** the `install-agent.ps1` script is downloaded
- **AND** the script can be executed to install the agent
