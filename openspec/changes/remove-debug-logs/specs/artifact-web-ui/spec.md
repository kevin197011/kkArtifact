## MODIFIED Requirements

### Requirement: Web UI Logging Hygiene
The Web UI SHALL NOT emit debug logs (e.g. `console.log`, `console.debug`) during normal runtime. Any optional debug logging MUST be explicitly enabled and MUST NOT expose secrets.

#### Scenario: Normal browsing produces no debug logs
- **WHEN** a user navigates to `/audit-logs`, `/tokens`, and `/` (inventory)
- **THEN** the Web UI does not emit `console.log` / `console.debug` output during normal operation

#### Scenario: Error handling is user-facing
- **WHEN** an API request fails on a page
- **THEN** the user is informed via UI feedback (e.g. toast/message)
- **AND** the Web UI does not rely on `console.error` for user-visible error reporting

