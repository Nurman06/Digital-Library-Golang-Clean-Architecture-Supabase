## ADDED Requirements

### Requirement: Book Creation
The system SHALL allow authorized users (Admin, Librarian) to create new book entries in the catalog with complete metadata.

#### Scenario: Create book with valid data
- **WHEN** an authorized user submits a book creation request with title, author, ISBN, category, publication year, and description
- **THEN** the system creates a new book entry and returns the book ID with a 201 Created status

#### Scenario: Create book with duplicate ISBN
- **WHEN** an authorized user attempts to create a book with an ISBN that already exists
- **THEN** the system returns a 409 Conflict error indicating the ISBN is already in use

#### Scenario: Create book with invalid data
- **WHEN** an authorized user submits incomplete or invalid book data (missing required fields, invalid ISBN format)
- **THEN** the system returns a 400 Bad Request error with validation details

#### Scenario: Unauthorized book creation attempt
- **WHEN** a non-authorized user (Member role or unauthenticated) attempts to create a book
- **THEN** the system returns a 403 Forbidden error

### Requirement: Book Retrieval
The system SHALL allow all users to retrieve book information from the catalog.

#### Scenario: Get book by ID
- **WHEN** a user requests a book by its unique ID
- **THEN** the system returns the complete book details including metadata and availability status

#### Scenario: Get non-existent book
- **WHEN** a user requests a book with an ID that does not exist
- **THEN** the system returns a 404 Not Found error

#### Scenario: List all books with pagination
- **WHEN** a user requests a list of books with page number and page size parameters
- **THEN** the system returns a paginated list of books with total count and navigation metadata

### Requirement: Book Update
The system SHALL allow authorized users (Admin, Librarian) to update existing book information.

#### Scenario: Update book metadata
- **WHEN** an authorized user submits updated book information for an existing book
- **THEN** the system updates the book record and returns the updated book data

#### Scenario: Update book with duplicate ISBN
- **WHEN** an authorized user attempts to change a book's ISBN to one that already exists for another book
- **THEN** the system returns a 409 Conflict error

#### Scenario: Update non-existent book
- **WHEN** an authorized user attempts to update a book that does not exist
- **THEN** the system returns a 404 Not Found error

#### Scenario: Unauthorized book update attempt
- **WHEN** a non-authorized user attempts to update a book
- **THEN** the system returns a 403 Forbidden error

### Requirement: Book Deletion
The system SHALL allow authorized users (Admin only) to delete books from the catalog with proper validation.

#### Scenario: Delete book with no active borrows
- **WHEN** an Admin user deletes a book that has no active borrow records
- **THEN** the system soft-deletes the book (marks as deleted) and returns a 204 No Content status

#### Scenario: Delete book with active borrows
- **WHEN** an Admin user attempts to delete a book that has active borrow records
- **THEN** the system returns a 409 Conflict error indicating the book cannot be deleted while borrowed

#### Scenario: Delete non-existent book
- **WHEN** an Admin user attempts to delete a book that does not exist
- **THEN** the system returns a 404 Not Found error

#### Scenario: Unauthorized book deletion attempt
- **WHEN** a non-Admin user attempts to delete a book
- **THEN** the system returns a 403 Forbidden error

### Requirement: Book Copy Management
The system SHALL support multiple physical copies of the same book with individual tracking.

#### Scenario: Add book copy
- **WHEN** an authorized user adds a new copy of an existing book
- **THEN** the system creates a new copy record with a unique copy ID and available status

#### Scenario: List book copies
- **WHEN** a user requests all copies of a specific book
- **THEN** the system returns a list of all copies with their individual status (available, borrowed, reserved, lost)

#### Scenario: Update copy status
- **WHEN** an authorized user updates the status of a book copy (e.g., marking as lost or damaged)
- **THEN** the system updates the copy status and adjusts availability accordingly

### Requirement: Book Metadata Validation
The system SHALL validate all book metadata according to defined business rules.

#### Scenario: Validate ISBN format
- **WHEN** a user submits a book with an ISBN
- **THEN** the system validates the ISBN format (ISBN-10 or ISBN-13) and rejects invalid formats

#### Scenario: Validate required fields
- **WHEN** a user submits a book creation or update request
- **THEN** the system validates that all required fields (title, author, ISBN) are present and non-empty

#### Scenario: Validate publication year
- **WHEN** a user submits a book with a publication year
- **THEN** the system validates the year is not in the future and is a reasonable historical date

### Requirement: Book Categorization
The system SHALL support organizing books into categories and genres.

#### Scenario: Assign category to book
- **WHEN** an authorized user creates or updates a book with a category
- **THEN** the system assigns the category and allows filtering by that category

#### Scenario: List books by category
- **WHEN** a user requests books filtered by a specific category
- **THEN** the system returns all books in that category with pagination support

#### Scenario: Support multiple categories
- **WHEN** a book is assigned to multiple categories or genres
- **THEN** the system stores all category associations and allows filtering by any of them