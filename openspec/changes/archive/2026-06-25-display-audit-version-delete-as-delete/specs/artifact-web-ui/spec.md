## ADDED Requirements

### Requirement: Audit Log Operation Type Labels
The web UI SHALL display human-readable labels for audit log operation types. The operation type `version_delete` SHALL be displayed as 「删除」 on the audit logs page and on the dashboard where audit entries are shown.

#### Scenario: version_delete displayed as 删除 on audit logs page
- **WHEN** the user navigates to the audit logs page
- **AND** the list contains an audit log entry with operation `version_delete`
- **THEN** the operation column for that entry displays 「删除」 (not the raw value `version_delete`)

#### Scenario: version_delete displayed as 删除 on dashboard
- **WHEN** the user views the dashboard recent activity (audit log entries)
- **AND** an entry has operation `version_delete`
- **THEN** the operation for that entry is displayed as 「删除」 (not the raw value `version_delete`)
