## MODIFIED Requirements

### Requirement: Server Debug Logging Gate
The server SHALL gate debug-only logs behind an explicit debug switch. In normal operation, the server SHALL avoid noisy debug output and SHALL NOT print sensitive information (tokens, credentials).

#### Scenario: Debug logs are disabled by default
- **WHEN** the server runs without debug enabled
- **THEN** debug-only logs are not emitted

#### Scenario: Debug logs are enabled explicitly
- **WHEN** debug mode is enabled (e.g. `DEBUG=true`)
- **THEN** debug logs MAY be emitted to assist troubleshooting
- **AND** logs MUST still avoid printing secrets (full tokens, credentials)

