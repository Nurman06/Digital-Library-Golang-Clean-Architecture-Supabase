## ADDED Requirements

### Requirement: User Registration
The system SHALL allow new users to register with email and password credentials.

#### Scenario: Successful user registration
- **WHEN** a new user submits valid registration data (email, password, name)
- **THEN** the system creates a new user account with Member role and returns user details with a 201 Created status

#### Scenario: Registration with existing email
- **WHEN** a user attempts to register with an email that already exists
- **THEN** the system returns a 409 Conflict error indicating the email is already registered

#### Scenario: Registration with invalid email format
- **WHEN** a user submits an invalid email format during registration
- **THEN** the system returns a 400 Bad Request error with validation details

#### Scenario: Registration with weak password
- **WHEN** a user submits a password that does not meet security requirements (minimum 8 characters, mix of letters and numbers)
- **THEN** the system returns a 400 Bad Request error with password requirements

### Requirement: User Login
The system SHALL authenticate users and provide JWT tokens for session management.

#### Scenario: Successful login with valid credentials
- **WHEN** a registered user submits correct email and password
- **THEN** the system validates credentials and returns a JWT access token with user details

#### Scenario: Login with invalid credentials
- **WHEN** a user submits incorrect email or password
- **THEN** the system returns a 401 Unauthorized error without revealing which credential was incorrect

#### Scenario: Login with non-existent email
- **WHEN** a user attempts to login with an email that is not registered
- **THEN** the system returns a 401 Unauthorized error

#### Scenario: Login for suspended user
- **WHEN** a suspended user attempts to login
- **THEN** the system returns a 403 Forbidden error indicating the account is suspended

### Requirement: JWT Token Management
The system SHALL issue and validate JWT tokens for authenticated requests.

#### Scenario: Generate JWT token on login
- **WHEN** a user successfully authenticates
- **THEN** the system generates a JWT token containing user ID, email, and role with appropriate expiration time

#### Scenario: Validate JWT token on protected endpoints
- **WHEN** a request is made to a protected endpoint with a valid JWT token
- **THEN** the system validates the token and allows access to the resource

#### Scenario: Reject expired JWT token
- **WHEN** a request is made with an expired JWT token
- **THEN** the system returns a 401 Unauthorized error indicating token expiration

#### Scenario: Reject invalid JWT token
- **WHEN** a request is made with a malformed or tampered JWT token
- **THEN** the system returns a 401 Unauthorized error

### Requirement: Role-Based Access Control
The system SHALL enforce role-based permissions for different user types (Admin, Librarian, Member).

#### Scenario: Admin access to all operations
- **WHEN** an Admin user attempts any operation in the system
- **THEN** the system grants access to all CRUD operations and administrative functions

#### Scenario: Librarian access to catalog management
- **WHEN** a Librarian user attempts to manage books and borrowing operations
- **THEN** the system grants access to book CRUD, borrowing operations, but denies user management

#### Scenario: Member access to read-only catalog
- **WHEN** a Member user attempts to view books and manage their own borrowing
- **THEN** the system grants read access to catalog and personal borrowing operations, but denies book modifications

#### Scenario: Unauthorized role access attempt
- **WHEN** a user attempts an operation not permitted for their role
- **THEN** the system returns a 403 Forbidden error

### Requirement: User Profile Management
The system SHALL allow users to view and update their own profile information.

#### Scenario: Get own user profile
- **WHEN** an authenticated user requests their own profile
- **THEN** the system returns complete profile information including name, email, role, and borrowing limits

#### Scenario: Update own profile
- **WHEN** an authenticated user updates their profile information (name, contact details)
- **THEN** the system updates the profile and returns the updated data

#### Scenario: Update email to existing email
- **WHEN** a user attempts to change their email to one already in use
- **THEN** the system returns a 409 Conflict error

#### Scenario: Access another user's profile
- **WHEN** a Member user attempts to access another user's profile
- **THEN** the system returns a 403 Forbidden error (only Admin can view other profiles)

### Requirement: Password Management
The system SHALL securely store and manage user passwords.

#### Scenario: Hash password on registration
- **WHEN** a user registers with a password
- **THEN** the system hashes the password using bcrypt before storing in the database

#### Scenario: Validate password on login
- **WHEN** a user attempts to login
- **THEN** the system compares the provided password with the hashed password using secure comparison

#### Scenario: Change password
- **WHEN** an authenticated user submits a password change request with current and new password
- **THEN** the system validates the current password, hashes the new password, and updates the user record

#### Scenario: Change password with incorrect current password
- **WHEN** a user attempts to change password with incorrect current password
- **THEN** the system returns a 401 Unauthorized error

### Requirement: User Account Status Management
The system SHALL support different user account statuses (Active, Suspended, Expired).

#### Scenario: Suspend user account
- **WHEN** an Admin suspends a user account
- **THEN** the system updates the user status to suspended and prevents login

#### Scenario: Reactivate suspended account
- **WHEN** an Admin reactivates a suspended user account
- **THEN** the system updates the user status to active and allows login

#### Scenario: Check account status on login
- **WHEN** a user attempts to login
- **THEN** the system verifies the account is active before issuing a token

### Requirement: Admin User Management
The system SHALL allow Admin users to manage other user accounts.

#### Scenario: List all users
- **WHEN** an Admin requests a list of all users
- **THEN** the system returns paginated user data with filtering and sorting options

#### Scenario: Update user role
- **WHEN** an Admin changes a user's role
- **THEN** the system updates the role and the new permissions take effect immediately

#### Scenario: Delete user account
- **WHEN** an Admin deletes a user account with no active borrows
- **THEN** the system soft-deletes the user account

#### Scenario: Prevent deletion of user with active borrows
- **WHEN** an Admin attempts to delete a user with active borrow records
- **THEN** the system returns a 409 Conflict error

### Requirement: Authentication Middleware
The system SHALL provide middleware to protect endpoints requiring authentication.

#### Scenario: Protect endpoint with authentication
- **WHEN** a request is made to a protected endpoint without authentication token
- **THEN** the system returns a 401 Unauthorized error

#### Scenario: Extract user context from token
- **WHEN** a valid authenticated request is made
- **THEN** the system extracts user information from JWT token and makes it available to request handlers

#### Scenario: Handle missing authorization header
- **WHEN** a request to a protected endpoint has no Authorization header
- **THEN** the system returns a 401 Unauthorized error with appropriate message