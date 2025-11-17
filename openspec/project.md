# Project Context

## Purpose

This is a Digital Library Management System built with Golang, implementing Clean Architecture principles and using Supabase as the backend infrastructure. The system manages books, users, and borrowing/lending operations with the following core features:

- **Book Catalog Management**: Complete CRUD operations for books with metadata (title, author, ISBN, category, etc.)
- **User Authentication & Authorization**: Secure user registration, login, and role-based access control
- **Borrowing/Lending Operations**: Track book checkouts, returns, due dates, and late fees
- **Borrowing History**: Maintain complete history of all borrowing transactions
- **Search & Discovery**: Advanced search functionality for finding books by various criteria
- **Availability Tracking**: Real-time book availability status

## Tech Stack

### Backend
- **Go (Golang)**: Primary programming language for backend services
- **Supabase**: Backend-as-a-Service platform providing:
  - PostgreSQL database
  - Authentication & authorization
  - Real-time subscriptions
  - Storage for book covers/documents
  - Row Level Security (RLS)

### Architecture & Frameworks
- **Clean Architecture**: Separation of concerns with clear boundaries between layers
- **Go Standard Library**: HTTP server, routing, and core functionality
- **Go Modules**: Dependency management

### Database
- **PostgreSQL**: Relational database via Supabase
- **Supabase Client Library**: Go SDK for Supabase integration

### Development Tools
- **Git**: Version control
- **GitHub**: Code hosting and collaboration
- **Go Test**: Unit and integration testing

## Project Conventions

### Code Style

- Follow official Go code style guidelines and idioms
- Use `gofmt` for consistent code formatting
- Use `golint` and `go vet` for static analysis
- **Naming Conventions**:
  - Package names: lowercase, single word (e.g., `book`, `user`, `repository`)
  - Interfaces: descriptive names ending with behavior (e.g., `BookRepository`, `UserService`)
  - Exported functions: PascalCase (e.g., `GetBookByID`, `CreateUser`)
  - Private functions: camelCase (e.g., `validateISBN`, `hashPassword`)
  - Constants: PascalCase or SCREAMING_SNAKE_CASE for exported/package-level constants
- **Error Handling**: Always check and handle errors explicitly; use custom error types for domain errors
- **Comments**: Provide package-level documentation and doc comments for all exported types and functions

### Architecture Patterns

The project follows **Clean Architecture** with four distinct layers:

1. **Entities Layer** (`internal/entity/` or `domain/`):
   - Core business models and rules
   - Independent of external frameworks and databases
   - Examples: `Book`, `User`, `BorrowRecord`

2. **Use Cases Layer** (`internal/usecase/` or `service/`):
   - Application business rules and orchestration
   - Interfaces for repositories and external services
   - Examples: `BookUseCase`, `BorrowingUseCase`

3. **Interface Adapters Layer** (`internal/adapter/` or `handler/`, `repository/`):
   - Controllers/Handlers for HTTP requests
   - Repository implementations for data access
   - Presenters for formatting responses
   - Examples: `BookHandler`, `PostgresBookRepository`

4. **Frameworks & Drivers Layer** (`internal/infrastructure/` or `pkg/`):
   - External frameworks and tools
   - Database connections
   - Web frameworks
   - Examples: `SupabaseClient`, `HTTPServer`

**Dependency Rule**: Dependencies point inward. Inner layers know nothing about outer layers.

**Key Patterns**:
- Repository Pattern for data access abstraction
- Dependency Injection for loose coupling
- Interface segregation for testability
- Factory pattern for object creation where appropriate

### Testing Strategy

- **Unit Tests**: Test individual functions and methods in isolation
  - Use table-driven tests for comprehensive coverage
  - Mock external dependencies using interfaces
  - Aim for >80% code coverage for business logic
  
- **Integration Tests**: Test interactions between layers
  - Test repository implementations with test database
  - Test use cases with mocked repositories
  
- **Test Organization**:
  - Test files named `*_test.go` alongside implementation
  - Use `testify` package for assertions and mocking if needed
  - Separate unit tests from integration tests using build tags

- **Test Database**: Use Supabase test project or local PostgreSQL for integration tests

### Git Workflow

- **Branching Strategy**:
  - `main`: Production-ready code, protected branch
  - `develop`: Integration branch for features (if using GitFlow)
  - `feature/*`: New features (e.g., `feature/book-search`)
  - `bugfix/*`: Bug fixes (e.g., `bugfix/borrowing-date-validation`)
  - `hotfix/*`: Critical production fixes

- **Commit Conventions** (Conventional Commits):
  - `feat:` New features (e.g., `feat: add book search by author`)
  - `fix:` Bug fixes (e.g., `fix: correct due date calculation`)
  - `docs:` Documentation changes
  - `refactor:` Code refactoring without behavior changes
  - `test:` Adding or updating tests
  - `chore:` Maintenance tasks (e.g., dependency updates)

- **Pull Request Process**:
  1. Create feature branch from `main` or `develop`
  2. Implement changes with tests
  3. Ensure all tests pass and code is formatted
  4. Create PR with descriptive title and summary
  5. Address review comments
  6. Squash commits if needed before merging

## Domain Context

### Core Domain Concepts

**Book Management**:
- Books have unique identifiers (ISBN) and metadata
- Books can have multiple copies (physical/digital)
- Categories/genres for classification
- Book availability status (available, borrowed, reserved, lost)

**User Management**:
- Different user roles: Admin, Librarian, Member
- User profiles with contact information
- Borrowing limits per user (e.g., max 5 books)
- User status (active, suspended, expired)

**Borrowing Operations**:
- Checkout process with due dates (typically 14-30 days)
- Return process with condition checking
- Late fee calculation based on overdue days
- Reservation system for unavailable books
- Renewal options (if allowed)

**Business Rules**:
- Users cannot borrow more than their limit
- Overdue books prevent new borrowing
- Books must be returned before renewal
- Late fees must be paid before new borrowing
- Admins and librarians have elevated privileges

## Important Constraints

### Technical Constraints
- Must use Supabase for all data persistence and authentication
- Go version compatibility (specify minimum version, e.g., Go 1.21+)
- RESTful API design principles
- Stateless service design for horizontal scalability
- Response time: API endpoints should respond within 500ms for 95th percentile

### Business Constraints
- User privacy: Personal data must be protected per GDPR/local regulations
- Data retention: Borrowing history retained for audit purposes (e.g., 7 years)
- Late fees: Configurable per library policy
- Borrowing period: Configurable per book type or category
- Maximum concurrent borrowings per user: Configurable

### Security Constraints
- Authentication required for all operations except public book browsing
- Role-based access control (RBAC) enforced at API level
- Supabase Row Level Security (RLS) policies for data protection
- Input validation and sanitization to prevent injection attacks
- Secure password storage using bcrypt or similar
- API rate limiting to prevent abuse

## External Dependencies

### Supabase Services
- **Authentication**: User registration, login, JWT token management
- **PostgreSQL Database**: Primary data store with ACID guarantees
- **Storage**: Book cover images and digital book files
- **Real-time**: Optional real-time updates for book availability
- **Edge Functions**: Serverless functions for background tasks (if needed)

### Third-Party Integrations (Future/Optional)
- **Email Service**: Notifications for due dates, overdue books (e.g., SendGrid, AWS SES)
- **ISBN Lookup API**: Automatic book metadata retrieval (e.g., Open Library API, Google Books API)
- **Payment Gateway**: Online payment for late fees (e.g., Stripe, PayPal)
- **Search Engine**: Enhanced search capabilities (e.g., Elasticsearch, Algolia)

### Development Dependencies
- Go modules for dependency management
- Testing frameworks (standard library `testing` package)
- Linting tools (`golangci-lint`)
- CI/CD pipeline (GitHub Actions recommended)

### Configuration Management
- Environment variables for sensitive configuration (Supabase URL, API keys)
- Configuration files for application settings
- Separate configurations for development, staging, and production environments
