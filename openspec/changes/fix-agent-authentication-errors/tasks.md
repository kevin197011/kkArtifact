## 1. Agent Token Cleaning and Validation
- [x] 1.1 Enhance `cleanTokenValue` function to handle all edge cases:
  - Remove all types of whitespace (spaces, tabs, newlines, carriage returns)
  - Handle encoding issues (UTF-8 BOM, zero-width characters)
  - Validate token format (base64 URL encoding pattern)
  - Remove any trailing/leading invisible characters
- [x] 1.2 Add token format validation function to check if token matches expected pattern
- [x] 1.3 Add token length validation (minimum/maximum length checks)
- [x] 1.4 Apply enhanced cleaning to all token sources (global config, local config, command-line)

## 2. Agent Token Validation Before API Calls
- [x] 2.1 Add comprehensive token validation in `client.New()` to fail early if token is invalid
- [x] 2.2 Add token validation in `pull` command before creating API client
- [x] 2.3 Add token validation in `push` command before creating API client
- [x] 2.4 Add token validation in all other commands that require authentication (update command allows empty token for public endpoints)
- [x] 2.5 Ensure all validation errors provide clear guidance on how to fix the issue

## 3. Improved Error Messages
- [x] 3.1 Enhance error messages in `GetManifest` to include:
  - Masked token preview (first 5 and last 5 characters)
  - Config file paths checked
  - Token length and format validation status
- [x] 3.2 Enhance error messages in `GetLatestVersion` with same details
- [x] 3.3 Enhance error messages in `InitUpload`, `UploadFile`, `FinishUpload` with same details
- [x] 3.4 Add error messages in `DownloadFile` and other download operations
- [x] 3.5 Ensure all 401 errors include actionable guidance

## 4. Server-Side Token Validation Improvements
- [x] 4.1 Add token trimming in server-side token extraction (handle whitespace)
- [x] 4.2 Improve error messages in `AuthenticateAPIToken` to be more descriptive (added maskTokenForLogging)
- [x] 4.3 Add logging for authentication failures (without exposing full tokens) - via maskTokenForLogging function
- [x] 4.4 Ensure token comparison handles edge cases (case sensitivity, encoding) - token trimming handles whitespace

## 5. Testing and Validation
- [ ] 5.1 Test token cleaning with various edge cases:
  - Tokens with leading/trailing whitespace
  - Tokens with newlines in YAML
  - Tokens with inline comments
  - Tokens with special characters
- [ ] 5.2 Test authentication with valid tokens in various formats
- [ ] 5.3 Test error messages for invalid/missing tokens
- [ ] 5.4 Verify both pull and push operations work correctly after fixes
- [ ] 5.5 Test token validation with tokens from different sources (global, local, command-line)
