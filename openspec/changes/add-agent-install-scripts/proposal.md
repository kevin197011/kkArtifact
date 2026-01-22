# Change: Add Agent Install Scripts and Update UI Text

## Why

Currently, users need to manually download agent binaries, set permissions, and configure paths. This process is error-prone and requires multiple manual steps:

1. **Manual download**: Users must download the correct binary for their platform
2. **Permission setup**: Unix users must manually run `chmod +x`
3. **Path configuration**: Users must manually move binaries to system paths
4. **UI text clarity**: "下载 Agent 客户端工具" with version tag is confusing - users want to "install" not just "download"

By providing one-click install scripts for Unix and Windows systems, we can:
- Automate the entire installation process
- Reduce user errors and support burden
- Provide a better user experience with clear "安装" (install) terminology
- Support both Unix-like systems (Linux, macOS) and Windows with platform-specific scripts

## What Changes

- **MODIFIED**: Frontend UI text - Change "下载 Agent 客户端工具" to "安装 agent 客户端工具" and remove version tag from header
- **ADDED**: Unix install script (`install-agent.sh`) - One-click installation for Linux/macOS
- **ADDED**: Windows install script (`install-agent.ps1`) - One-click installation for Windows
- **ADDED**: Script download endpoints - Server API to serve install scripts
- **MODIFIED**: Frontend download section - Add install script download buttons alongside binary downloads

## Impact

- **Affected specs**: 
  - `artifact-web-ui` - UI text and download section updates
  - `artifact-agent` - Install script functionality
- **Affected code**:
  - `web-ui/src/pages/InventoryPage.tsx` - Update UI text and add install script links
  - `server/internal/api/download_handlers.go` - Add install script endpoints
  - `scripts/install-agent.sh` - New Unix install script
  - `scripts/install-agent.ps1` - New Windows install script
- **Breaking changes**: None - this is an additive change
- **Migration**: No migration needed

## Benefits

1. **Improved UX**: Clear "安装" (install) terminology instead of "下载" (download)
2. **Automation**: One-command installation reduces manual steps
3. **Cross-platform**: Support for both Unix and Windows systems
4. **Error reduction**: Automated scripts prevent common installation mistakes
5. **Consistency**: Standardized installation process across platforms
