# Design: Add Agent Install Scripts and Update UI Text

## Context

The kkArtifact agent is a CLI tool that users need to install on their systems. Currently:
- Users download binaries manually from the web UI
- They must manually set permissions (`chmod +x`) on Unix systems
- They must manually move binaries to system paths
- The UI shows "下载 Agent 客户端工具" with a version tag, which is confusing

## Goals

1. Change UI text from "下载" (download) to "安装" (install) for clarity
2. Remove version tag from the main header (version info can be shown elsewhere)
3. Provide one-click install scripts for Unix-like systems (Linux, macOS)
4. Provide one-click install scripts for Windows (PowerShell)
5. Automate the entire installation process (download, permissions, path setup)

## Non-Goals

- Uninstall scripts (can be added later if needed)
- Update scripts (agent already has `update` command)
- Package manager integration (apt, yum, homebrew, etc.) - can be added later
- GUI installer - CLI scripts only

## Decisions

### Decision 1: Separate Scripts for Unix and Windows

**What**: Create two separate install scripts:
- `install-agent.sh` for Unix-like systems (Linux, macOS, BSD)
- `install-agent.ps1` for Windows (PowerShell)

**Why**:
- Different platforms require different approaches
- Unix scripts use bash/sh, Windows uses PowerShell
- Clear separation makes maintenance easier
- Users can easily identify the correct script for their platform

**Alternatives considered**:
- Single script with platform detection: More complex, harder to maintain
- Batch file for Windows: PowerShell is more modern and feature-rich

### Decision 2: Script Location and Serving

**What**: Store install scripts in `server/static/scripts/` and serve via API endpoint `/api/v1/downloads/scripts/install-agent.{sh|ps1}`

**Why**:
- Centralized location for all static files
- Consistent with agent binary serving pattern
- Easy to update scripts without code changes
- Can be cached and versioned

**Alternatives considered**:
- Embed scripts in Go code: Harder to update, requires rebuild
- Store in web-ui public folder: Not accessible from API, harder to version

### Decision 3: Installation Path

**What**: 
- Unix: Install to `/usr/local/bin/kkartifact-agent` (system-wide) or `~/.local/bin/kkartifact-agent` (user-local if no sudo)
- Windows: Install to `%LOCALAPPDATA%\kkartifact\kkartifact-agent.exe` (user-local)

**Why**:
- `/usr/local/bin` is standard for system-wide Unix binaries
- `~/.local/bin` is standard for user-local Unix binaries (no sudo required)
- `%LOCALAPPDATA%` is standard Windows user application data location
- User-local installation doesn't require admin privileges

**Alternatives considered**:
- Always require sudo: Less user-friendly, may fail
- Always user-local: May not be in PATH, requires manual PATH setup

### Decision 4: Script Functionality

**What**: Install scripts should:
1. Detect platform/architecture
2. Download correct binary from server
3. Set executable permissions (Unix)
4. Install to appropriate location
5. Verify installation
6. Provide usage instructions

**Why**:
- Complete automation reduces user errors
- Platform detection ensures correct binary selection
- Verification confirms successful installation
- Usage instructions help users get started

### Decision 5: UI Text Changes

**What**: 
- Change header from "下载 Agent 客户端工具" to "安装 agent 客户端工具"
- Remove version tag from header (keep version info in description or elsewhere)
- Add install script download buttons alongside binary download buttons

**Why**:
- "安装" (install) is clearer than "下载" (download) - users want to install, not just download
- Removing version tag from header reduces visual clutter
- Install script buttons provide one-click installation option
- Binary download buttons remain for users who prefer manual installation

## Risks / Trade-offs

### Risk 1: Script Security
**Risk**: Users downloading and executing scripts from the internet raises security concerns.

**Mitigation**:
- Scripts should be simple and readable
- Use HTTPS for script downloads
- Provide script checksums for verification
- Document security best practices

### Risk 2: Platform Detection
**Risk**: Incorrect platform detection may download wrong binary.

**Mitigation**:
- Use standard platform detection methods (`uname`, `$OS`, etc.)
- Provide clear error messages if detection fails
- Allow manual platform selection as fallback

### Risk 3: Permission Issues
**Risk**: Installation may fail due to permission issues (no sudo, no write access).

**Mitigation**:
- Try system-wide installation first, fallback to user-local
- Provide clear error messages with solutions
- Document permission requirements

### Risk 4: PATH Configuration
**Risk**: Installed binary may not be in user's PATH.

**Mitigation**:
- Use standard installation paths that are typically in PATH
- Provide instructions for adding to PATH if needed
- Test on common systems to ensure PATH compatibility

## Implementation Plan

1. **Create install scripts**:
   - `scripts/install-agent.sh` for Unix systems
   - `scripts/install-agent.ps1` for Windows

2. **Add script serving endpoints**:
   - Update `server/internal/api/download_handlers.go` to serve scripts
   - Store scripts in `server/static/scripts/`

3. **Update frontend UI**:
   - Change "下载 Agent 客户端工具" to "安装 agent 客户端工具"
   - Remove version tag from header
   - Add install script download buttons

4. **Test installation**:
   - Test on Linux (various distributions)
   - Test on macOS
   - Test on Windows (various versions)

## Open Questions

1. Should we provide script checksums for verification? (Recommended but not critical)
2. Should we support custom installation paths? (Not needed initially)
3. Should we add uninstall scripts? (Can be added later if needed)
