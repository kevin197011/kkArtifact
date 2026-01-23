# Change: Add Server URL Environment Variable Override for Install Scripts

## Why

Currently, the install scripts (`install-agent.sh` and `install-agent.ps1`) use a server-injected URL when downloaded from the server API. However, users may need to override this URL in certain scenarios:

1. **Custom API endpoints**: Users may want to point the agent to a different API endpoint than the one serving the script
2. **Testing/Development**: Developers may need to test against different server instances
3. **Proxy/Network configurations**: Users behind proxies or with complex network setups may need to specify a different URL
4. **Flexibility**: Providing environment variable override gives users more control over the installation process

The current implementation only checks for `SERVER_URL` environment variable, but the user's request shows they want to use `server_url` (lowercase) as the environment variable name, which is more consistent with the configuration file format.

## What Changes

- **MODIFIED**: Install scripts (`install-agent.sh` and `install-agent.ps1`) - Add support for `server_url` environment variable override
- **MODIFIED**: Server URL resolution priority - Update priority order to: 1) `server_url` env var, 2) `SERVER_URL` env var, 3) Server-injected URL, 4) Default localhost
- **ADDED**: Documentation - Update README to document the environment variable override option

## Impact

- **Affected specs**: 
  - `artifact-agent` - Install script behavior
- **Affected code**:
  - `scripts/install-agent.sh` - Add `server_url` environment variable support
  - `scripts/install-agent.ps1` - Add `server_url` environment variable support
  - `README.md` - Document the environment variable override option
- **Breaking changes**: None - this is an additive change that maintains backward compatibility
- **Migration**: No migration needed

## Benefits

1. **Flexibility**: Users can override the server URL without modifying the script
2. **Consistency**: Using `server_url` (lowercase) matches the configuration file format
3. **Backward compatibility**: Existing `SERVER_URL` environment variable still works
4. **Better UX**: Supports use cases like `curl ... | server_url="https://custom-endpoint.com" bash`
