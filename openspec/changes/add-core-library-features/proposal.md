# Change: Add Core Digital Library Features

## Why

This change establishes the foundational capabilities for the Digital Library Management System. The system currently has no implemented features, and this proposal defines the complete set of core functionalities required for a functional library system: book catalog management, user authentication and authorization, borrowing/lending operations, borrowing history tracking, search and discovery features, and real-time availability tracking.

## What Changes

This proposal adds six new capabilities to the system:

- **Book Catalog Management**: Complete CRUD operations for books with metadata (title, author, ISBN, category, publication year, etc.), supporting multiple copies per book
- **User Authentication & Authorization**: Secure user registration, login, JWT-based authentication, and role-based access control (Admin, Librarian, Member roles)
- **Borrowing/Lending Operations**: Book checkout and return processes with due date management, late fee calculation, and renewal functionality
- **Borrowing History**: Complete audit trail of all borrowing transactions with timestamps and status tracking
- **Search & Discovery**: Advanced search functionality allowing users to find books by title, author, ISBN, category, and other criteria
- **Availability Tracking**: Real-time tracking of book availability status (available, borrowed, reserved, lost) with copy-level granularity

## Impact

- **Affected specs**: Creates 6 new capability specifications
  - `specs/book-catalog/spec.md`
  - `specs/user-auth/spec.md`
  - `specs/borrowing-operations/spec.md`
  - `specs/borrowing-history/spec.md`
  - `specs/search-discovery/spec.md`
  - `specs/availability-tracking/spec.md`

- **Affected code**: This is a greenfield implementation that will create:
  - Entity models for Book, User, BorrowRecord, and related domain objects
  - Use case/service layer for business logic orchestration
  - Repository interfaces and implementations for data access
  - HTTP handlers for RESTful API endpoints
  - Supabase client integration for database and authentication
  - Database schema and migrations
  - Configuration management for environment-specific settings

- **External dependencies**: 
  - Supabase Go client library for database and authentication
  - PostgreSQL database via Supabase
  - Standard Go libraries for HTTP server and routing

- **Breaking changes**: None (initial implementation)