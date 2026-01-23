# Change: Fix Agent Authentication Errors

## Why

Users are experiencing persistent 401 Unauthorized errors when using `kkartifact-agent pull` and potentially `push` commands, despite having valid tokens configured in `/etc/kkArtifact/config.yml`. The errors indicate that tokens are not being properly validated by the server, even though:

- Tokens are correctly stored in the configuration file
- Token values appear to be valid base64-encoded strings
- The agent is sending Authorization headers with Bearer tokens

This change addresses root causes of authentication failures by:
- Improving token parsing and cleaning to handle edge cases (whitespace, encoding issues)
- Adding comprehensive token validation and debugging
- Ensuring consistent token format between agent and server
- Adding detailed error messages to help diagnose authentication issues

## What Changes

- **MODIFIED**: Enhanced token cleaning in agent config loading to handle all edge cases (whitespace, newlines, encoding issues)
- **MODIFIED**: Added token validation before API calls in all agent operations (pull, push, and other commands)
- **MODIFIED**: Improved error messages to include token debugging information (masked token preview, config file paths)
- **ADDED**: Token format validation to ensure tokens match expected base64 URL encoding pattern
- **ADDED**: Server-side token validation improvements to handle edge cases in token comparison
- **ADDED**: Comprehensive logging for authentication failures (without exposing full tokens)

**BREAKING**: None - this is a bug fix that improves error handling and validation.

## Impact

- **Affected specs**: `artifact-agent` (authentication and configuration), `artifact-api` (token validation)
- **Affected code**: 
  - `agent/internal/config/config.go` - Enhanced token cleaning and validation
  - `agent/internal/client/client.go` - Improved token validation and error messages
  - `agent/internal/cli/pull.go` - Enhanced token validation before API calls
  - `agent/internal/cli/push.go` - Enhanced token validation before API calls
  - `server/internal/auth/middleware.go` - Improved token validation and error handling
- **User impact**: Users will get clearer error messages when authentication fails, and edge cases causing 401 errors will be resolved
