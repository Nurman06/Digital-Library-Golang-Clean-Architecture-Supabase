# Integration Test Setup Guide

## Overview

Integration tests untuk repository layer telah dibuat dan siap digunakan. Tests ini memerlukan PostgreSQL database yang berjalan.

## Files Created

1. **Test Infrastructure**
   - `testhelper/db_helper.go` - Database helper utilities
   - `docker-compose.test.yml` - Test database container configuration
   - `.env.test` - Test database environment variables

2. **Integration Test Files**
   - `postgres_book_repository_integration_test.go` - 604 lines, 6 test functions
   - `postgres_user_repository_integration_test.go` - 479 lines, 9 test functions

3. **Documentation**
   - `README_INTEGRATION_TESTS.md` - Complete guide for running tests

## Quick Start

### 1. Start Test Database

```bash
# Using docker compose v2
docker compose -f docker-compose.test.yml up -d

# Or using docker-compose v1
docker-compose -f docker-compose.test.yml up -d

# Wait for database to be ready
sleep 5
```

### 2. Run Integration Tests

```bash
# Run all integration tests
make test-integration

# Or manually
go test -v -tags=integration ./internal/repository/...

# Run specific test
go test -v -tags=integration ./internal/repository/ -run TestPostgresBookRepository_Create
```

### 3. Stop Test Database

```bash
# Stop database
make test-db-down

# Or manually
docker compose -f docker-compose.test.yml down

# Clean everything (including volumes)
docker compose -f docker-compose.test.yml down -v
```

## Test Coverage

### BookRepository Integration Tests (6 tests)
- ✅ Create (including duplicate ISBN validation)
- ✅ GetByID (existing and non-existent)
- ✅ Update (verify all fields updated)
- ✅ Delete (soft delete verification)
- ✅ List (pagination and sorting)
- ✅ Search (query, author, category, ISBN, year range)
- ✅ GetByISBN
- ✅ ExistsByISBN
- ✅ Count

### UserRepository Integration Tests (9 tests)
- ✅ Create (including duplicate email validation)
- ✅ GetByID (existing and non-existent)
- ✅ GetByEmail (existing and non-existent)
- ✅ Update (verify all fields updated)
- ✅ Delete (hard delete verification)
- ✅ List (pagination)
- ✅ Search (by name and email)
- ✅ GetByRole (filter by user role)
- ✅ GetByStatus (filter by user status)
- ✅ UpdateStatus
- ✅ UpdateBorrowingLimit
- ✅ ExistsByEmail
- ✅ Count

## Environment Variables

Tests use these environment variables (with defaults):

```bash
TEST_DB_HOST=localhost      # Database host
TEST_DB_PORT=5433          # Database port (different from production)
TEST_DB_USER=postgres      # Database user
TEST_DB_PASSWORD=postgres  # Database password
TEST_DB_NAME=library_test  # Database name
TEST_DB_SSL_MODE=disable   # SSL mode
```

## Build Tags

Integration tests use the `integration` build tag:

```go
// +build integration
```

This allows:
- Skip integration tests in normal `go test` runs
- Run only integration tests with `-tags=integration`
- Skip with `-short` flag

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Integration Tests

on: [push, pull_request]

jobs:
  integration-test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: library_test
        ports:
          - 5433:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run integration tests
        env:
          TEST_DB_HOST: localhost
          TEST_DB_PORT: 5433
          TEST_DB_USER: postgres
          TEST_DB_PASSWORD: postgres
          TEST_DB_NAME: library_test
        run: go test -v -tags=integration ./internal/repository/...
```

## Troubleshooting

### Cannot Connect to Database

1. Check if database is running:
   ```bash
   docker ps | grep library-test-db
   ```

2. Check database logs:
   ```bash
   docker logs library-test-db
   ```

3. Test connection manually:
   ```bash
   docker exec -it library-test-db psql -U postgres -d library_test
   ```

### Port Conflict

If port 5433 is already in use:

1. Change port in `docker-compose.test.yml`:
   ```yaml
   ports:
     - "5434:5432"
   ```

2. Update environment variable:
   ```bash
   export TEST_DB_PORT=5434
   ```

### Tests Fail with "table already exists"

The test helper automatically creates tables. If you see this error:

1. Clean the database:
   ```bash
   make test-db-clean
   make test-db-up
   ```

2. Or drop tables manually:
   ```bash
   docker exec -it library-test-db psql -U postgres -d library_test -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
   ```

## Next Steps

To complete repository integration tests, add:

1. `postgres_borrow_record_repository_integration_test.go`
2. `postgres_book_copy_repository_integration_test.go`
3. `postgres_reservation_repository_integration_test.go`

These can be added in future PRs following the same pattern.

## Notes

- ⚠️ Integration tests require Docker to be installed and running
- ⚠️ Tests are NOT run in CI/CD by default (requires setup)
- ✅ Tests are isolated and can run in parallel
- ✅ Each test cleans up after itself
- ✅ Tests use transactions where possible for speed