## ADDED Requirements

### Requirement: Real-Time Availability Status
The system SHALL maintain and display real-time availability status for all books and copies.

#### Scenario: Display book availability
- **WHEN** a user views a book's details
- **THEN** the system displays current availability status (available, borrowed, reserved, lost, damaged)

#### Scenario: Show available copy count
- **WHEN** viewing a book with multiple copies
- **THEN** the system displays total copies and number of currently available copies

#### Scenario: Update availability on checkout
- **WHEN** a book copy is checked out
- **THEN** the system immediately updates the availability status to borrowed and decrements available count

#### Scenario: Update availability on return
- **WHEN** a book copy is returned
- **THEN** the system immediately updates the availability status to available and increments available count

### Requirement: Copy-Level Tracking
The system SHALL track availability at the individual copy level for books with multiple copies.

#### Scenario: Track individual copy status
- **WHEN** viewing a book's copies
- **THEN** the system displays status for each physical copy (copy ID, barcode, status, location)

#### Scenario: Assign specific copy on checkout
- **WHEN** a user checks out a book with multiple available copies
- **THEN** the system assigns a specific copy and tracks which copy is borrowed

#### Scenario: Track copy location
- **WHEN** viewing a copy's details
- **THEN** the system displays the copy's current location (shelf location or "on loan" with due date)

#### Scenario: Handle damaged or lost copies
- **WHEN** a copy is marked as damaged or lost
- **THEN** the system updates that copy's status and excludes it from available count

### Requirement: Availability Status Types
The system SHALL support multiple availability status types for comprehensive tracking.

#### Scenario: Available status
- **WHEN** a book copy is not borrowed and in good condition
- **THEN** the system marks it as "available" and includes it in borrowable inventory

#### Scenario: Borrowed status
- **WHEN** a book copy is currently checked out
- **THEN** the system marks it as "borrowed" with borrower information and due date

#### Scenario: Reserved status
- **WHEN** a book copy is held for a user with reservation
- **THEN** the system marks it as "reserved" with reservation holder and expiration date

#### Scenario: Lost status
- **WHEN** a book copy is reported lost
- **THEN** the system marks it as "lost" and removes it from available and total active inventory

#### Scenario: Damaged status
- **WHEN** a book copy is damaged and under repair
- **THEN** the system marks it as "damaged" and temporarily removes it from borrowable inventory

#### Scenario: In processing status
- **WHEN** a new book copy is being cataloged
- **THEN** the system marks it as "in processing" until ready for circulation

### Requirement: Availability Queries
The system SHALL provide efficient queries for checking book availability.

#### Scenario: Check if book is available
- **WHEN** a user or system checks if a book has any available copies
- **THEN** the system returns a boolean result based on current availability status

#### Scenario: Get next available date
- **WHEN** all copies of a book are borrowed
- **THEN** the system calculates and displays the earliest expected return date

#### Scenario: Find available books in category
- **WHEN** searching for available books in a specific category
- **THEN** the system returns only books with at least one available copy

#### Scenario: Bulk availability check
- **WHEN** checking availability for multiple books simultaneously
- **THEN** the system efficiently returns availability status for all requested books

### Requirement: Availability History
The system SHALL maintain a history of availability changes for auditing and analytics.

#### Scenario: Track status change history
- **WHEN** a book copy's status changes
- **THEN** the system records the status change with timestamp and reason

#### Scenario: View availability timeline
- **WHEN** an admin views a copy's history
- **THEN** the system displays a timeline of all status changes

#### Scenario: Calculate availability metrics
- **WHEN** analyzing book performance
- **THEN** the system calculates metrics like availability rate, average borrowed duration, and utilization

### Requirement: Availability Notifications
The system SHALL provide availability-based notification capabilities.

#### Scenario: Notify when book becomes available
- **WHEN** a reserved book becomes available
- **THEN** the system identifies users with active reservations for notification

#### Scenario: Hold period for reserved books
- **WHEN** a reserved book becomes available
- **THEN** the system holds it for the reserved user for a specified period (e.g., 3 days)

#### Scenario: Release after hold expiration
- **WHEN** the hold period expires without checkout
- **THEN** the system releases the book back to available status and notifies next reservation

### Requirement: Availability Dashboard
The system SHALL provide a dashboard view of overall availability metrics.

#### Scenario: Display availability summary
- **WHEN** a librarian views the availability dashboard
- **THEN** the system displays total books, available copies, borrowed copies, and utilization rate

#### Scenario: Show low availability alerts
- **WHEN** viewing the dashboard
- **THEN** the system highlights books with low availability (e.g., high demand books with few available copies)

#### Scenario: Track circulation rate
- **WHEN** analyzing library performance
- **THEN** the system calculates and displays circulation rate (checkouts per book per time period)

### Requirement: Availability Filters and Search
The system SHALL integrate availability filtering in search and browse features.

#### Scenario: Filter catalog by availability
- **WHEN** browsing the catalog
- **THEN** users can filter to show only currently available books

#### Scenario: Sort by availability
- **WHEN** viewing search results
- **THEN** users can sort books with available copies first

#### Scenario: Show availability badge
- **WHEN** displaying books in lists or search results
- **THEN** the system shows a visual indicator (badge) for availability status

### Requirement: Reservation Queue Management
The system SHALL manage reservation queues based on availability.

#### Scenario: Track reservation queue position
- **WHEN** a user reserves an unavailable book
- **THEN** the system assigns a queue position and displays estimated wait time

#### Scenario: Automatic queue advancement
- **WHEN** a book becomes available
- **THEN** the system automatically notifies the next user in the reservation queue

#### Scenario: Display queue length
- **WHEN** viewing an unavailable book
- **THEN** the system displays the number of users in the reservation queue

### Requirement: Availability Predictions
The system SHALL provide predictive availability information based on historical data.

#### Scenario: Estimate return date
- **WHEN** all copies are borrowed
- **THEN** the system displays estimated return dates based on due dates

#### Scenario: Predict availability patterns
- **WHEN** analyzing a popular book
- **THEN** the system predicts typical availability patterns (e.g., "usually available on weekdays")

#### Scenario: Recommend alternative books
- **WHEN** a book is unavailable
- **THEN** the system suggests similar available books as alternatives

### Requirement: Multi-Location Availability
The system SHALL track availability across multiple library locations if applicable.

#### Scenario: Show availability by location
- **WHEN** a library has multiple branches
- **THEN** the system displays which locations have available copies

#### Scenario: Transfer between locations
- **WHEN** a book is requested from another location
- **THEN** the system tracks the transfer and updates availability at both locations

#### Scenario: Location-specific search
- **WHEN** a user searches with location filter
- **THEN** the system returns only books available at the specified location

### Requirement: Availability API
The system SHALL provide API endpoints for availability queries and updates.

#### Scenario: Get book availability via API
- **WHEN** an external system queries book availability
- **THEN** the system returns current availability status and copy count via REST API

#### Scenario: Bulk availability query
- **WHEN** requesting availability for multiple books
- **THEN** the system provides efficient batch query endpoint

#### Scenario: Real-time availability updates
- **WHEN** availability changes occur
- **THEN** the system can push real-time updates to subscribed clients (optional WebSocket support)

### Requirement: Availability Cache Management
The system SHALL efficiently cache availability data for performance.

#### Scenario: Cache availability status
- **WHEN** querying frequently accessed book availability
- **THEN** the system serves cached data with appropriate TTL (Time To Live)

#### Scenario: Invalidate cache on status change
- **WHEN** a book's availability status changes
- **THEN** the system immediately invalidates relevant cache entries

#### Scenario: Warm cache for popular books
- **WHEN** the system starts or during off-peak hours
- **THEN** it pre-loads availability data for popular books into cache

### Requirement: Availability Reporting
The system SHALL generate reports on availability and circulation patterns.

#### Scenario: Generate availability report
- **WHEN** a librarian requests an availability report
- **THEN** the system generates a report showing availability rates by category, time period, and location

#### Scenario: Identify underutilized books
- **WHEN** analyzing collection
- **THEN** the system identifies books that are rarely borrowed despite being available

#### Scenario: Highlight high-demand books
- **WHEN** reviewing collection needs
- **THEN** the system identifies books with high demand and frequent unavailability for acquisition decisions