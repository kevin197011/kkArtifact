## MODIFIED Requirements

### Requirement: Token Authentication and Validation
The server SHALL correctly validate authentication tokens from Authorization headers, handling edge cases such as whitespace and encoding issues. Authentication failures SHALL be logged appropriately without exposing sensitive token information.

#### Scenario: Validate token with whitespace
- **WHEN** the server receives an Authorization header with a token containing leading/trailing whitespace: `Bearer  YVza5JEjUgzaWzULeQM198aXzCzVmJZw0oNjY5Db_rc= `
- **THEN** the server trims whitespace from the token before validation
- **AND** the token is successfully validated against stored hashes

#### Scenario: Handle missing Authorization header
- **WHEN** a request is made without an Authorization header
- **THEN** the server returns 401 Unauthorized with error message: `{"error": "missing authorization header"}`
- **AND** the error is logged appropriately

#### Scenario: Handle invalid Authorization format
- **WHEN** a request is made with an Authorization header that doesn't start with "Bearer "
- **THEN** the server returns 401 Unauthorized with error message: `{"error": "invalid authorization format"}`
- **AND** the error is logged appropriately

#### Scenario: Handle invalid token
- **WHEN** a request is made with a token that doesn't match any stored token hash
- **THEN** the server returns 401 Unauthorized with error message: `{"error": "unauthorized"}`
- **AND** the error is logged with masked token information (first 5 and last 5 characters) for debugging

#### Scenario: Token validation with edge cases
- **WHEN** the server receives tokens with various edge cases (whitespace, encoding issues)
- **THEN** the server handles these cases gracefully by trimming and normalizing before validation
- **AND** valid tokens are successfully authenticated regardless of minor formatting differences
