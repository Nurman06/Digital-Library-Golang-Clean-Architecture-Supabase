## ADDED Requirements

### Requirement: User Borrowing History
The system SHALL maintain a complete history of all borrowing transactions for each user.

#### Scenario: View personal borrowing history
- **WHEN** an authenticated user requests their borrowing history
- **THEN** the system returns a paginated list of all their past and current borrow records with book details

#### Scenario: Filter history by status
- **WHEN** a user filters their borrowing history by status (active, returned, overdue)
- **THEN** the system returns only borrow records matching the specified status

#### Scenario: Filter history by date range
- **WHEN** a user requests borrowing history for a specific date range
- **THEN** the system returns borrow records with checkout dates within the specified range

#### Scenario: Sort history by date
- **WHEN** a user requests their borrowing history
- **THEN** the system allows sorting by checkout date, return date, or due date in ascending or descending order

### Requirement: Book Borrowing History
The system SHALL track the complete borrowing history for each book.

#### Scenario: View book's borrowing history
- **WHEN** a librarian or admin requests a book's borrowing history
- **THEN** the system returns all past borrow records for that book with anonymized or full user details based on role

#### Scenario: Track book circulation statistics
- **WHEN** viewing a book's history
- **THEN** the system displays total times borrowed, average borrowing duration, and return rate

#### Scenario: Identify frequent borrowers of a book
- **WHEN** an admin analyzes a book's history
- **THEN** the system identifies users who have borrowed the book multiple times

### Requirement: Borrow Record Details
The system SHALL provide comprehensive details for each borrow record in history.

#### Scenario: View complete borrow record
- **WHEN** a user or admin views a specific borrow record
- **THEN** the system displays book details, user details, checkout date, due date, return date, late fees, and status

#### Scenario: Include book copy information
- **WHEN** viewing a borrow record
- **THEN** the system includes which specific copy of the book was borrowed (copy ID)

#### Scenario: Track renewal history
- **WHEN** a borrow record has been renewed
- **THEN** the system displays renewal count and dates of each renewal

#### Scenario: Display late fee details
- **WHEN** a returned book had late fees
- **THEN** the system shows late fee amount, calculation basis, payment status, and any waivers

### Requirement: Admin Borrowing History Access
The system SHALL allow administrators to access borrowing history across all users.

#### Scenario: View all borrowing history
- **WHEN** an admin requests system-wide borrowing history
- **THEN** the system returns paginated borrow records for all users with filtering and sorting options

#### Scenario: Search history by user
- **WHEN** an admin searches borrowing history by user email or ID
- **THEN** the system returns all borrow records for that specific user

#### Scenario: Search history by book
- **WHEN** an admin searches borrowing history by book title or ISBN
- **THEN** the system returns all borrow records for that specific book

#### Scenario: Export borrowing history
- **WHEN** an admin requests to export borrowing history data
- **THEN** the system generates a downloadable report in CSV or JSON format

### Requirement: Borrowing Statistics
The system SHALL generate statistics and analytics from borrowing history.

#### Scenario: Calculate user borrowing statistics
- **WHEN** viewing a user's profile or history
- **THEN** the system displays total books borrowed, current active borrows, overdue count, and on-time return rate

#### Scenario: Calculate library circulation statistics
- **WHEN** an admin requests library statistics
- **THEN** the system provides total checkouts, returns, active borrows, and overdue books for specified time periods

#### Scenario: Identify most borrowed books
- **WHEN** an admin requests popular books report
- **THEN** the system ranks books by number of times borrowed in a specified period

#### Scenario: Track late fee collection
- **WHEN** an admin views financial statistics
- **THEN** the system displays total late fees assessed, collected, and waived for specified time periods

### Requirement: History Data Retention
The system SHALL retain borrowing history according to data retention policies.

#### Scenario: Retain completed borrow records
- **WHEN** a book is returned
- **THEN** the system permanently retains the borrow record for audit and statistical purposes

#### Scenario: Archive old records
- **WHEN** borrow records exceed the active retention period (e.g., 7 years)
- **THEN** the system archives records while maintaining data integrity for compliance

#### Scenario: Prevent history deletion
- **WHEN** an attempt is made to delete a borrow record
- **THEN** the system prevents deletion and only allows soft-delete or archival for data integrity

### Requirement: History Privacy Controls
The system SHALL respect user privacy when displaying borrowing history.

#### Scenario: User views only own history
- **WHEN** a regular member user requests borrowing history
- **THEN** the system only displays their own borrow records, not other users' records

#### Scenario: Anonymize history for privacy
- **WHEN** displaying book borrowing history to non-admin users
- **THEN** the system anonymizes borrower information to protect user privacy

#### Scenario: Admin access with audit trail
- **WHEN** an admin accesses another user's borrowing history
- **THEN** the system logs the access for audit purposes

### Requirement: History Search and Filtering
The system SHALL provide comprehensive search and filtering capabilities for borrowing history.

#### Scenario: Full-text search in history
- **WHEN** a user searches their borrowing history with a text query
- **THEN** the system searches across book titles, authors, and ISBNs and returns matching records

#### Scenario: Filter by multiple criteria
- **WHEN** a user applies multiple filters (status, date range, book category)
- **THEN** the system returns records matching all specified criteria

#### Scenario: Quick filter for overdue items
- **WHEN** a user selects the overdue filter
- **THEN** the system immediately displays all their currently overdue borrow records

### Requirement: History Timeline View
The system SHALL provide a chronological timeline view of borrowing activity.

#### Scenario: Display chronological timeline
- **WHEN** a user views their borrowing history in timeline mode
- **THEN** the system displays events (checkouts, returns, renewals) in chronological order with visual indicators

#### Scenario: Highlight important events
- **WHEN** viewing timeline
- **THEN** the system highlights significant events like overdue dates, late fees assessed, and renewals

#### Scenario: Group by time periods
- **WHEN** viewing history timeline
- **THEN** the system allows grouping by day, week, month, or year for easier navigation

### Requirement: Borrowing Pattern Analysis
The system SHALL analyze borrowing patterns for insights and recommendations.

#### Scenario: Identify reading preferences
- **WHEN** analyzing a user's borrowing history
- **THEN** the system identifies preferred categories, authors, and genres based on past borrows

#### Scenario: Detect unusual borrowing patterns
- **WHEN** monitoring borrowing activity
- **THEN** the system flags unusual patterns (e.g., sudden spike in borrows, frequent overdue books) for review

#### Scenario: Generate borrowing trends report
- **WHEN** an admin requests trend analysis
- **THEN** the system generates reports showing borrowing trends over time, seasonal patterns, and popular periods