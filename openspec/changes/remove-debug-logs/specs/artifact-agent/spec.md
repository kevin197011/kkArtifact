## MODIFIED Requirements

### Requirement: Agent Debug Output Control
The agent SHALL NOT emit verbose debug dumps (request/response detail blocks) during normal operation. Verbose dumps MUST only be emitted when explicitly enabled (env var and/or flag). The agent MUST still provide clear error messages without requiring debug mode.

#### Scenario: Default agent output is clean
- **WHEN** a user runs `kkartifact-agent pull` or `kkartifact-agent push` in normal mode
- **THEN** the agent does not print verbose request/response dump blocks to stderr

#### Scenario: Debug mode enables verbose dumps
- **WHEN** debug mode is explicitly enabled for the agent
- **THEN** the agent MAY print verbose request/response dumps for troubleshooting
- **AND** dumps MUST NOT include full tokens (only masked token preview)

