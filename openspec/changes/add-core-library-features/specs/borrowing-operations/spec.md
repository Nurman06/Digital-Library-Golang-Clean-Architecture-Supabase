## ADDED Requirements

### Requirement: Book Checkout
The system SHALL allow authenticated users to borrow available books with proper validation.

#### Scenario: Successful book checkout
- **WHEN** an authenticated user with available borrowing capacity requests to borrow an available book
- **THEN** the system creates a borrow record, updates book availability, calculates due date, and returns checkout details

#### Scenario: Checkout unavailable book
- **WHEN** a user attempts to borrow a book with no available copies
- **THEN** the system returns a 409 Conflict error indicating the book is not available

#### Scenario: Checkout when borrowing limit reached
- **WHEN** a user attempts to borrow a book when they have reached their maximum borrowing limit (e.g., 5 books)
- **THEN** the system returns a 409 Conflict error indicating the borrowing limit has been reached

#### Scenario: Checkout with overdue books
- **WHEN** a user with overdue books attempts to borrow another book
- **THEN** the system returns a 403 Forbidden error indicating overdue books must be returned first

#### Scenario: Checkout with unpaid late fees
- **WHEN** a user with unpaid late fees attempts to borrow a book
- **THEN** the system returns a 403 Forbidden error indicating late fees must be paid first

#### Scenario: Checkout by suspended user
- **WHEN** a suspended user attempts to borrow a book
- **THEN** the system returns a 403 Forbidden error

### Requirement: Book Return
The system SHALL allow users to return borrowed books and update system state accordingly.

#### Scenario: Return book on time
- **WHEN** a user returns a borrowed book before or on the due date
- **THEN** the system updates the borrow record status to returned, updates book availability, and records return date

#### Scenario: Return book late
- **WHEN** a user returns a borrowed book after the due date
- **THEN** the system calculates late fees, updates the borrow record with fees owed, and updates book availability

#### Scenario: Return non-borrowed book
- **WHEN** a user attempts to return a book they have not borrowed
- **THEN** the system returns a 404 Not Found error

#### Scenario: Return already returned book
- **WHEN** a user attempts to return a book that has already been returned
- **THEN** the system returns a 409 Conflict error

#### Scenario: Return with damage report
- **WHEN** a librarian processes a return and reports damage
- **THEN** the system records damage information, may assess damage fees, and updates book copy status

### Requirement: Book Renewal
The system SHALL allow users to renew borrowed books if eligible.

#### Scenario: Successful book renewal
- **WHEN** a user requests to renew a borrowed book that is eligible for renewal (not overdue, no reservations)
- **THEN** the system extends the due date by the renewal period and updates the borrow record

#### Scenario: Renew overdue book
- **WHEN** a user attempts to renew an overdue book
- **THEN** the system returns a 403 Forbidden error indicating overdue books cannot be renewed

#### Scenario: Renew book with reservations
- **WHEN** a user attempts to renew a book that has been reserved by another user
- **THEN** the system returns a 409 Conflict error indicating the book cannot be renewed due to reservation

#### Scenario: Renew at maximum renewal limit
- **WHEN** a user attempts to renew a book that has already been renewed the maximum number of times (e.g., 2 renewals)
- **THEN** the system returns a 409 Conflict error indicating renewal limit reached

#### Scenario: Renew non-borrowed book
- **WHEN** a user attempts to renew a book they have not borrowed
- **THEN** the system returns a 404 Not Found error

### Requirement: Due Date Calculation
The system SHALL automatically calculate due dates based on book type and library policy.

#### Scenario: Calculate standard due date
- **WHEN** a book is checked out
- **THEN** the system calculates the due date as 14 days from checkout date for standard books

#### Scenario: Calculate extended due date for reference books
- **WHEN** a reference book is checked out
- **THEN** the system calculates the due date as 30 days from checkout date

#### Scenario: Calculate renewal due date
- **WHEN** a book is renewed
- **THEN** the system extends the due date by the standard borrowing period from the current due date

### Requirement: Late Fee Calculation
The system SHALL calculate late fees for overdue books according to library policy.

#### Scenario: Calculate late fee for standard book
- **WHEN** a book is returned after the due date
- **THEN** the system calculates late fee as $0.50 per day overdue with a maximum cap of $10.00

#### Scenario: Calculate late fee for different book types
- **WHEN** a reference book is returned late
- **THEN** the system applies the appropriate late fee rate for that book type (e.g., $1.00 per day)

#### Scenario: No late fee for on-time return
- **WHEN** a book is returned on or before the due date
- **THEN** the system does not assess any late fees

#### Scenario: Late fee maximum cap
- **WHEN** late fees exceed the maximum cap
- **THEN** the system caps the late fee at the maximum amount regardless of overdue days

### Requirement: Borrowing Limit Enforcement
The system SHALL enforce borrowing limits per user role and status.

#### Scenario: Check borrowing limit on checkout
- **WHEN** a user attempts to checkout a book
- **THEN** the system verifies the user has not reached their maximum concurrent borrowing limit

#### Scenario: Different limits by role
- **WHEN** checking borrowing limits
- **THEN** the system applies role-specific limits (Member: 5 books, Librarian: 10 books, Admin: unlimited)

#### Scenario: Temporary limit increase
- **WHEN** an Admin grants a temporary borrowing limit increase to a user
- **THEN** the system allows borrowing up to the increased limit for the specified period

### Requirement: Borrow Record Management
The system SHALL maintain accurate records of all borrowing transactions.

#### Scenario: Create borrow record on checkout
- **WHEN** a book is checked out
- **THEN** the system creates a borrow record with user ID, book copy ID, checkout date, due date, and status

#### Scenario: Update borrow record on return
- **WHEN** a book is returned
- **THEN** the system updates the borrow record with return date, late fees (if any), and status

#### Scenario: Track borrow record status
- **WHEN** querying borrow records
- **THEN** the system maintains status values (active, returned, overdue, renewed)

#### Scenario: Link borrow records to specific book copies
- **WHEN** a book with multiple copies is borrowed
- **THEN** the system tracks which specific copy was borrowed in the borrow record

### Requirement: Overdue Book Notifications
The system SHALL identify and track overdue books for notification purposes.

#### Scenario: Identify overdue books
- **WHEN** the system checks for overdue books
- **THEN** it identifies all borrow records where the current date is past the due date and status is active

#### Scenario: Calculate overdue days
- **WHEN** checking an overdue book
- **THEN** the system calculates the number of days overdue from the due date to current date

#### Scenario: List user's overdue books
- **WHEN** a user or admin requests overdue books for a specific user
- **THEN** the system returns all active borrow records past their due date with overdue days

### Requirement: Book Reservation
The system SHALL allow users to reserve books that are currently unavailable.

#### Scenario: Reserve unavailable book
- **WHEN** a user requests to reserve a book with no available copies
- **THEN** the system creates a reservation with the user's position in the queue

#### Scenario: Notify when reserved book available
- **WHEN** a reserved book becomes available
- **THEN** the system identifies the next user in the reservation queue for notification

#### Scenario: Cancel reservation
- **WHEN** a user cancels their book reservation
- **THEN** the system removes the reservation and advances the queue

#### Scenario: Reservation expiration
- **WHEN** a reserved book becomes available and the user does not checkout within the hold period (e.g., 3 days)
- **THEN** the system moves to the next reservation in the queue

### Requirement: Librarian Override Functions
The system SHALL provide librarians with override capabilities for special cases.

#### Scenario: Override borrowing limit
- **WHEN** a librarian approves a checkout that exceeds normal borrowing limits
- **THEN** the system allows the checkout with an override flag and records the librarian's authorization

#### Scenario: Waive late fees
- **WHEN** a librarian waives late fees for a borrow record
- **THEN** the system updates the borrow record to remove or reduce late fees and records the waiver

#### Scenario: Manual due date adjustment
- **WHEN** a librarian manually adjusts a due date for valid reasons
- **THEN** the system updates the due date and records the adjustment reason