## ADDED Requirements

### Requirement: Token-Based Authentication
The system SHALL authenticate all API requests using Bearer tokens provided in the Authorization header. Tokens SHALL be stored securely (hashed) and validated on each request.

#### Scenario: Authenticate with valid token
- **WHEN** a request includes `Authorization: Bearer <token>` with a valid token
- **THEN** the request is authenticated
- **AND** the token's permissions are checked for the requested operation

#### Scenario: Reject invalid token
- **WHEN** a request includes an invalid or expired token
- **THEN** the request is rejected with 401 Unauthorized

### Requirement: Token Scope Model
The system SHALL support three token scope levels: Global (all projects/apps), Project (all apps in a project), and App (single app). Tokens SHALL have a scope that determines their access boundaries.

#### Scenario: Global token access
- **WHEN** a token has Global scope
- **THEN** the token can access all projects and all apps
- **AND** the token can perform operations on any project/app combination

#### Scenario: Project token access
- **WHEN** a token has Project scope for project `myproject`
- **THEN** the token can access all apps within `myproject`
- **AND** the token cannot access apps in other projects

#### Scenario: App token access
- **WHEN** a token has App scope for project `myproject` and app `myapp`
- **THEN** the token can only access `myproject/myapp`
- **AND** the token cannot access other apps or projects

### Requirement: Token Permissions
Tokens SHALL have permissions that determine what operations they can perform. Permissions SHALL include: push (upload artifacts), pull (download artifacts), promote (mark versions as promoted), and admin (full access including token management).

#### Scenario: Token with push permission
- **WHEN** a token has push permission for a project/app
- **THEN** the token can upload artifacts
- **AND** the token cannot pull, promote, or perform admin operations

#### Scenario: Token with pull permission
- **WHEN** a token has pull permission for a project/app
- **THEN** the token can download artifacts
- **AND** the token cannot push, promote, or perform admin operations

#### Scenario: Token with promote permission
- **WHEN** a token has promote permission for a project/app
- **THEN** the token can mark versions as promoted
- **AND** the token can pull artifacts but cannot push

#### Scenario: Token with admin permission
- **WHEN** a token has admin permission
- **THEN** the token can perform all operations including token management
- **AND** the token can create, update, and revoke other tokens

### Requirement: Permission Validation
The system SHALL validate token permissions before allowing operations. Permission checks SHALL consider token scope, requested operation, and target project/app.

#### Scenario: Validate permission for push
- **WHEN** a push operation is requested for project `myproject` and app `myapp`
- **THEN** the system checks if the token has push permission for that scope
- **AND** if the token scope is App-level, it must match `myproject/myapp`
- **AND** if the token scope is Project-level, it must match `myproject`
- **AND** if the token scope is Global, permission is granted

#### Scenario: Deny operation without permission
- **WHEN** a token lacks the required permission for an operation
- **THEN** the operation is rejected with 403 Forbidden
- **AND** an error message indicates the required permission

### Requirement: Token Storage
Tokens SHALL be stored securely with their values hashed using bcrypt or argon2. The system SHALL never return token values in plaintext after creation.

#### Scenario: Store token securely
- **WHEN** a new token is created
- **THEN** the token value is hashed using bcrypt or argon2
- **AND** only the hash is stored in the database
- **AND** the plaintext token is returned only once during creation

#### Scenario: Validate token from hash
- **WHEN** a token is provided in an API request
- **THEN** the system hashes the provided token
- **AND** compares it with stored hashes to find a match

### Requirement: Token Expiration
Tokens SHALL support optional expiration dates. Expired tokens SHALL be rejected during authentication.

#### Scenario: Token with expiration
- **WHEN** a token has an expiration date set
- **THEN** requests using the token after expiration are rejected
- **AND** the response indicates the token has expired

#### Scenario: Token without expiration
- **WHEN** a token does not have an expiration date
- **THEN** the token remains valid until explicitly revoked

### Requirement: Token Revocation
The system SHALL support revoking tokens, which immediately invalidates them for all future requests.

#### Scenario: Revoke token
- **WHEN** a token is revoked via admin API
- **THEN** the token is marked as revoked in the database
- **AND** all subsequent requests using that token are rejected with 401 Unauthorized

