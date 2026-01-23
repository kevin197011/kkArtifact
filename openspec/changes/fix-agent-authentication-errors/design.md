## Context

Users are experiencing 401 Unauthorized errors when using agent commands, despite having valid tokens configured. Investigation reveals potential issues with:

1. **Token Parsing**: YAML parsing might include whitespace, comments, or encoding issues
2. **Token Format**: Tokens might have invisible characters or encoding problems
3. **Token Validation**: Lack of early validation leads to unclear error messages
4. **Error Messages**: Current error messages don't provide enough debugging information

## Goals / Non-Goals

### Goals
- Ensure tokens are correctly parsed and cleaned from all sources (global config, local config, command-line)
- Validate token format before making API calls
- Provide clear, actionable error messages when authentication fails
- Handle all edge cases in token parsing (whitespace, encoding, comments)
- Ensure consistent token format between agent and server

### Non-Goals
- Changing token storage format (remains base64 URL encoding)
- Changing authentication mechanism (remains Bearer token)
- Adding token refresh/rotation (out of scope)
- Changing server-side token validation algorithm (bcrypt remains)

## Decisions

### Decision 1: Comprehensive Token Cleaning
**What**: Enhance `cleanTokenValue` to handle all edge cases including:
- All whitespace types (spaces, tabs, newlines, carriage returns, zero-width spaces)
- Encoding issues (UTF-8 BOM, non-printable characters)
- Inline comments in YAML
- Base64 URL encoding validation

**Why**: 
- YAML parsing can include unexpected characters
- User input might have formatting issues
- Base64 tokens should not contain certain characters

**Implementation**: 
- Use `strings.TrimSpace` and regex to remove all whitespace
- Validate token matches base64 URL encoding pattern: `^[A-Za-z0-9_-]+$`
- Check token length (typically 32-44 characters for base64)

### Decision 2: Early Token Validation
**What**: Validate tokens in `client.New()` and before creating API clients in commands.

**Why**:
- Fail fast with clear error messages
- Avoid making API calls with invalid tokens
- Provide better user experience

**Implementation**:
- Add `ValidateToken(token string) error` function
- Call validation in `client.New()` and fail if invalid
- Call validation in command handlers before creating clients

### Decision 3: Enhanced Error Messages
**What**: Include masked token preview, config paths, and validation status in error messages.

**Why**:
- Help users debug authentication issues
- Provide actionable guidance
- Don't expose full tokens for security

**Implementation**:
- Mask tokens: show first 5 and last 5 characters (e.g., `YVza5...b_rc=`)
- Include config file paths checked
- Include token length and format validation status

### Decision 4: Server-Side Token Trimming
**What**: Trim whitespace from tokens in server-side authentication.

**Why**:
- Handle edge cases where agent might send tokens with whitespace
- Ensure consistent token comparison
- Defense in depth

**Implementation**:
- Trim token after extracting from Authorization header
- Apply before token comparison

## Risks / Trade-offs

### Risk: Over-aggressive Token Cleaning
**Mitigation**: 
- Only remove characters that are definitely not part of valid tokens
- Validate token format after cleaning
- Preserve valid base64 characters

### Risk: Breaking Valid Tokens
**Mitigation**:
- Test with various token formats
- Only clean characters that are definitely problematic
- Validate before and after cleaning

### Trade-off: Error Message Detail vs Security
- Show enough information to debug (masked token, config paths)
- Don't expose full tokens or sensitive information
- Balance helpfulness with security

## Migration Plan

1. **Agent**: Enhance token cleaning and validation
2. **Agent**: Add validation before API calls
3. **Agent**: Improve error messages
4. **Server**: Add token trimming in authentication
5. **Testing**: Test with various edge cases
6. **Documentation**: Update error message documentation

## Open Questions

- Should we add a `--debug` flag to show more token information?
  - **Decision**: Not in this change, but can be added later if needed

- Should we validate token format against a regex?
  - **Decision**: Yes, validate base64 URL encoding pattern
