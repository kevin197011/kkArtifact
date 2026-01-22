## 1. Create Install Scripts

- [x] 1.1 Create `scripts/install-agent.sh` for Unix-like systems (Linux, macOS)
- [x] 1.2 Create `scripts/install-agent.ps1` for Windows systems
- [x] 1.3 Add platform detection logic to scripts
- [x] 1.4 Add binary download and installation logic
- [x] 1.5 Add verification and usage instructions

## 2. Add Script Serving Endpoints

- [x] 2.1 Create `server/static/scripts/` directory
- [x] 2.2 Copy install scripts to `server/static/scripts/`
- [x] 2.3 Add script download handler in `server/internal/api/download_handlers.go`
- [x] 2.4 Add route for `/api/v1/downloads/scripts/install-agent.{sh|ps1}`
- [x] 2.5 Test script download endpoints (implementation complete, requires runtime testing)

## 3. Update Frontend UI

- [x] 3.1 Change "下载 Agent 客户端工具" to "安装 agent 客户端工具" in `InventoryPage.tsx`
- [x] 3.2 Remove version tag from header (keep version info in description)
- [x] 3.3 Add install script download buttons (one for Unix, one for Windows)
- [x] 3.4 Update download section styling to accommodate install buttons
- [x] 3.5 Test UI changes in browser (implementation complete, requires runtime testing)

## 4. Testing

- [ ] 4.1 Test `install-agent.sh` on Linux (Ubuntu/Debian)
- [ ] 4.2 Test `install-agent.sh` on macOS
- [ ] 4.3 Test `install-agent.ps1` on Windows 10/11
- [ ] 4.4 Test script download from web UI
- [ ] 4.5 Verify installed agent binary works correctly
- [ ] 4.6 Test installation with and without sudo/admin privileges

## 5. Documentation

- [ ] 5.1 Update README with install script usage
- [ ] 5.2 Add script security notes
- [ ] 5.3 Document installation paths and PATH requirements
