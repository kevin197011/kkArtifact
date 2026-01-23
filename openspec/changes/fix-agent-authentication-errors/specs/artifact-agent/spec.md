## MODIFIED Requirements

### Requirement: Token Configuration and Validation
The agent SHALL correctly parse, clean, and validate authentication tokens from configuration files and command-line arguments. Tokens SHALL be validated before making any authenticated API requests, and validation failures SHALL provide clear, actionable error messages.

#### Scenario: Parse token from global config with whitespace
- **WHEN** the global config file contains a token with leading/trailing whitespace: `token: " YVza5JEjUgzaWzULeQM198aXzCzVmJZw0oNjY5Db_rc= "`
- **THEN** the agent removes all whitespace and uses the clean token value
- **AND** the token is successfully validated before API calls

#### Scenario: Parse token from config with inline comment
- **WHEN** the config file contains a token with an inline comment: `token: YVza5JEjUgzaWzULeQM198aXzCzVmJZw0oNjY5Db_rc=  # Uncomment and set your token here`
- **THEN** the agent removes the comment and uses only the token value
- **AND** the token is successfully validated before API calls

#### Scenario: Parse token with newlines
- **WHEN** the config file contains a token with newlines or carriage returns
- **THEN** the agent removes all newline characters and uses the clean token
- **AND** the token is successfully validated before API calls

#### Scenario: Validate token format before API call
- **WHEN** the agent is about to make an authenticated API request
- **THEN** the agent validates the token format (base64 URL encoding pattern)
- **AND** if the token format is invalid, the agent fails with a clear error message before making the request
- **AND** the error message includes: masked token preview, config file paths, and format validation status

#### Scenario: Handle invalid token format
- **WHEN** the token contains invalid characters (not base64 URL encoding)
- **THEN** the agent fails with a clear error message: "Token format is invalid. Tokens must be base64 URL encoded strings."
- **AND** the error message includes the config file path where the token was found

#### Scenario: Handle empty token
- **WHEN** the token is empty or not set in any config file
- **THEN** the agent fails with a clear error message before making API calls
- **AND** the error message lists all config file paths that were checked

#### Scenario: Pull command with invalid token
- **WHEN** the user runs `kkartifact-agent pull` with an invalid token
- **THEN** the agent validates the token before making any API requests
- **AND** if validation fails, the agent exits with a clear error message
- **AND** the error message includes guidance on how to fix the issue

#### Scenario: Push command with invalid token
- **WHEN** the user runs `kkartifact-agent push` with an invalid token
- **THEN** the agent validates the token before making any API requests
- **AND** if validation fails, the agent exits with a clear error message
- **AND** the error message includes guidance on how to fix the issue

#### Scenario: Enhanced error messages for 401 responses
- **WHEN** the agent receives a 401 Unauthorized response from the server
- **THEN** the error message includes:
  - Masked token preview (first 5 and last 5 characters)
  - Config file paths where the token was found
  - Token length and format validation status
  - Guidance on checking token validity in the server
