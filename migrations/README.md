# Database Migrations

This directory contains SQL migration scripts for the Digital Library Management System database schema.

## Migration Files

- `001_create_tables.sql` - Initial schema creation (up migration)
- `001_create_tables_down.sql` - Rollback script for initial schema (down migration)
- `002_setup_rls_policies.sql` - Row Level Security policies setup (up migration)
- `002_setup_rls_policies_down.sql` - Rollback script for RLS policies (down migration)

## Database Schema Overview

The schema includes the following tables:

1. **books** - Stores book catalog information (title, author, ISBN, category, etc.)
2. **book_copies** - Tracks individual physical/digital copies of books
3. **users** - Stores user profiles extending Supabase Auth
4. **borrow_records** - Records all borrowing transactions with history
5. **reservations** - Manages book reservation queue system

### Row Level Security (RLS) Policies

Migration 002 implements comprehensive RLS policies for data security:

- **Books**: Public read access, admin/librarian write access
- **Book Copies**: Public read access, admin/librarian management
- **Users**: Self-profile access, admin full access
- **Borrow Records**: Users see own records, staff see all
- **Reservations**: Users manage own reservations, staff manage all

## Running Migrations

### Using Supabase CLI

```bash
# Apply migration
supabase db push

# Or apply specific migration
supabase migration up

# Rollback migration
supabase migration down
```

### Using psql (Direct PostgreSQL)

```bash
# Apply migration
psql -h <host> -U <user> -d <database> -f migrations/001_create_tables.sql

# Rollback migration
psql -h <host> -U <user> -d <database> -f migrations/001_create_tables_down.sql
```

### Using golang-migrate

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Apply all migrations
migrate -path migrations -database "postgresql://user:pass@host:port/dbname?sslmode=disable" up

# Rollback one migration
migrate -path migrations -database "postgresql://user:pass@host:port/dbname?sslmode=disable" down 1

# Check migration version
migrate -path migrations -database "postgresql://user:pass@host:port/dbname?sslmode=disable" version
```

## Schema Features

### UUID Generation
- Uses PostgreSQL's `gen_random_uuid()` function from the `pgcrypto` extension
- All primary keys are UUIDs for distributed system compatibility

### Indexes
- Full-text search indexes on book titles and authors (GIN indexes)
- Performance indexes on foreign keys and frequently queried columns
- Partial indexes for active records (e.g., active borrow records)

### Constraints
- Foreign key constraints with appropriate CASCADE/RESTRICT policies
- CHECK constraints for enum-like fields (status, role, etc.)
- UNIQUE constraints on business keys (ISBN, email, copy_number)

### Soft Deletes
- Books table includes `deleted_at` column for soft delete functionality
- Allows maintaining historical data while hiding deleted records

### Timestamps
- All tables include `created_at` and `updated_at` timestamps
- Automatic `updated_at` updates via database triggers

## Migration Best Practices

1. **Never modify existing migrations** - Always create new migration files for changes
2. **Test migrations** - Always test both up and down migrations in development first
3. **Backup before production** - Always backup production database before applying migrations
4. **Version control** - Keep all migration files in version control
5. **Sequential numbering** - Use sequential numbering for migration files (001, 002, 003, etc.)

## Rollback Strategy

Each migration includes a corresponding rollback script (`*_down.sql`) that:
- Drops tables in reverse order to respect foreign key constraints
- Removes indexes before dropping tables
- Drops triggers before dropping functions
- Preserves extension (pgcrypto) as it may be used by other schemas

## Data Integrity

The schema enforces data integrity through:
- Foreign key constraints prevent orphaned records
- CHECK constraints validate enum values
- NOT NULL constraints ensure required fields
- UNIQUE constraints prevent duplicates
- Triggers maintain timestamp accuracy

## Performance Considerations

- GIN indexes for full-text search on books
- B-tree indexes on foreign keys for JOIN performance
- Partial indexes on status fields for active record queries
- Proper index selection based on query patterns from design document

## Future Migrations

When adding new migrations:
1. Create `00X_description.sql` for the up migration
2. Create `00X_description_down.sql` for the rollback
3. Update this README with the new migration details
4. Test both up and down migrations thoroughly
5. Document any breaking changes or special considerations