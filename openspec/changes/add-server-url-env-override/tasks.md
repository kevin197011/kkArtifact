## 1. Implementation

- [x] 1.1 Update `install-agent.sh` to support `server_url` environment variable
  - [x] 1.1.1 Add check for `server_url` env var (lowercase) with highest priority
  - [x] 1.1.2 Maintain backward compatibility with `SERVER_URL` (uppercase)
  - [x] 1.1.3 Update priority order: `server_url` > `SERVER_URL` > injected URL > default
  - [x] 1.1.4 Test with `curl ... | server_url="https://test.com" bash`

- [x] 1.2 Update `install-agent.ps1` to support `server_url` environment variable
  - [x] 1.2.1 Add check for `$env:server_url` (lowercase) with highest priority
  - [x] 1.2.2 Maintain backward compatibility with `$env:SERVER_URL` (uppercase)
  - [x] 1.2.3 Update priority order: `server_url` > `SERVER_URL` > injected URL > default
  - [x] 1.2.4 Test with PowerShell environment variable override

- [x] 1.3 Update server static scripts
  - [x] 1.3.1 Copy updated scripts to `server/static/scripts/` directory

- [x] 1.4 Update documentation
  - [x] 1.4.1 Update README.md to document `server_url` environment variable override
  - [x] 1.4.2 Add example usage: `curl ... | server_url="https://custom.com" bash`

## 2. Testing

- [ ] 2.1 Test Unix script with `server_url` override
  - [ ] 2.1.1 Verify `server_url` takes precedence over injected URL
  - [ ] 2.1.2 Verify `server_url` takes precedence over `SERVER_URL`
  - [ ] 2.1.3 Verify backward compatibility with `SERVER_URL`
  - [ ] 2.1.4 Verify fallback to injected URL when env vars not set

- [ ] 2.2 Test PowerShell script with `server_url` override
  - [ ] 2.2.1 Verify `$env:server_url` takes precedence over injected URL
  - [ ] 2.2.2 Verify `$env:server_url` takes precedence over `$env:SERVER_URL`
  - [ ] 2.2.3 Verify backward compatibility with `$env:SERVER_URL`
  - [ ] 2.2.4 Verify fallback to injected URL when env vars not set

- [ ] 2.3 Test configuration file creation
  - [ ] 2.3.1 Verify config file uses `server_url` env var value when set
  - [ ] 2.3.2 Verify config file uses injected URL when env vars not set

## 3. Validation

- [x] 3.1 Run OpenSpec validation
  - [x] 3.1.1 Execute `openspec validate add-server-url-env-override --strict`
  - [x] 3.1.2 Fix any validation errors

- [x] 3.2 Code review
  - [x] 3.2.1 Review script changes for correctness
  - [x] 3.2.2 Review documentation updates
  - [x] 3.2.3 Verify backward compatibility
