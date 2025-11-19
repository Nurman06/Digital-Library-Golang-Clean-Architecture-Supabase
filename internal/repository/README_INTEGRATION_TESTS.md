# Integration Tests for Repository Layer

This directory contains integration tests for repository implementations that test against a real PostgreSQL database.

## Prerequisites

- Docker and Docker Compose installed
- Go 1.21 or later

## Setup Test Database

### 1. Start Test Database

```bash
# Start the test database container
docker-compose -f docker-compose.test.yml up -d

# Wait for database to be ready
docker-compose -f docker-compose.test.yml ps
```

### 2. Verify Database Connection

```bash
# Connect to test database
docker exec -it library-test-db psql -U postgres -d library_test

# Or check if it's ready
docker-compose -f docker-compose.test.yml exec test-db pg_isready -U postgres
```

## Running Integration Tests

### Run All Integration Tests

```bash
# Run integration tests with the integration build tag
go test -v -tags=integration ./internal/repository/...

# Or using make
make test-integration
```

### Run Specific Test

```bash
# Run specific test file
go test -v -tags=integration ./internal/repository/ -run TestPostgresBookRepository

# Run specific test function
go test -v -tags=integration ./internal/repository/ -run TestPostgresBookRepository_Create
```

### Skip Integration Tests (Run Only Unit Tests)

```bash
# Integration tests are skipped by default without the tag
go test -v ./internal/repository/...

# Or explicitly skip with -short flag
go test -v -short ./internal/repository/...
```

## Test Database Configuration

Integration tests use environment variables for database configuration:

```bash
# Default values (matches docker-compose.test.yml)
TEST_DB_HOST=localhost
TEST_DB_PORT=5433
TEST_DB_USER=postgres
TEST_DB_PASSWORD=postgres
TEST_DB_NAME=library_test
TEST_DB_SSL_MODE=disable
```

You can override these by setting environment variables:

```bash
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5433
go test -v -tags=integration ./internal/repository/...
```

Or create a `.env.test` file (already provided).

## Test Structure

### Test Helper (`testhelper/db_helper.go`)

Provides utilities for integration testing:

- `SetupTestDB()` - Creates test database connection
- `RunMigrations()` - Runs database migrations
- `CleanupTables()` - Truncates all tables for clean state
- `DropTables()` - Drops all tables (cleanup)
- `BeginTx()` / `Rollback()` - Transaction support

### Integration Test Files

- `postgres_book_repository_integration_test.go` - Book repository tests
- `postgres_user_repository_integration_test.go` - User repository tests
- `postgres_borrow_record_repository_integration_test.go` - Borrow record repository tests (to be added)
- `postgres_book_copy_repository_integration_test.go` - Book copy repository tests (to be added)

## Test Coverage

Integration tests cover:

### BookRepository
- ✅ Create (with duplicate ISBN validation)
- ✅ GetByID
- ✅ GetByISBN
- ✅ Update
- ✅ Delete (soft delete)
- ✅ List (with pagination and sorting)
- ✅ Search (by query, author, category, ISBN, year range)
- ✅ ExistsByISBN
- ✅ Count

### UserRepository
- ✅ Create (with duplicate email validation)
- ✅ GetByID
- ✅ GetByEmail
- ✅ Update
- ✅ Delete
- ✅ List (with pagination)
- ✅ Search (by name and email)
- ✅ GetByRole
- ✅ GetByStatus
- ✅ UpdateStatus
- ✅ UpdateBorrowingLimit
- ✅ ExistsByEmail
- ✅ Count

## Cleanup

### Stop Test Database

```bash
# Stop the test database
docker-compose -f docker-compose.test.yml down

# Stop and remove volumes (clean slate)
docker-compose -f docker-compose.test.yml down -v
```

## CI/CD Integration

To run integration tests in CI/CD:

```yaml
# Example GitHub Actions workflow
- name: Start test database
  run: docker-compose -f docker-compose.test.yml up -d

- name: Wait for database
  run: |
    timeout 30 bash -c 'until docker-compose -f docker-compose.test.yml exec -T test-db pg_isready -U postgres; do sleep 1; done'

- name: Run integration tests
  run: go test -v -tags=integration ./internal/repository/...

- name: Stop test database
  run: docker-compose -f docker-compose.test.yml down
```

## Best Practices

1. **Isolation**: Each test should clean up after itself
2. **Idempotency**: Tests should be repeatable
3. **No Side Effects**: Tests should not affect each other
4. **Fast Cleanup**: Use `CleanupTables()` instead of recreating database
5. **Transaction Rollback**: Use transactions for tests that don't need to commit

## Troubleshooting

### Database Connection Failed

```bash
# Check if database is running
docker-compose -f docker-compose.test.yml ps

# Check database logs
docker-compose -f docker-compose.test.yml logs test-db

# Restart database
docker-compose -f docker-compose.test.yml restart test-db
```

### Port Already in Use

If port 5433 is already in use, modify `docker-compose.test.yml`:

```yaml
ports:
  - "5434:5432"  # Use different port
```

And update environment variable:

```bash
export TEST_DB_PORT=5434
```

### Migration Errors

If migrations fail, check:

1. Database is running and accessible
2. Migration SQL syntax is correct
3. Tables don't already exist (use `CleanupTables()` or `DropTables()`)

## Notes

- Integration tests use the `// +build integration` build tag
- Tests are skipped when running `go test` without `-tags=integration`
- Tests are also skipped when running with `-short` flag
- Test database runs on port 5433 (different from production port 5432)