package testhelper

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// TestDB wraps a database connection for testing
type TestDB struct {
	DB *sql.DB
	tx *sql.Tx
}

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) *TestDB {
	t.Helper()

	// Get test database configuration from environment
	host := getEnv("TEST_DB_HOST", "localhost")
	port := getEnv("TEST_DB_PORT", "5433") // Different port for test DB
	user := getEnv("TEST_DB_USER", "postgres")
	password := getEnv("TEST_DB_PASSWORD", "postgres")
	dbname := getEnv("TEST_DB_NAME", "library_test")
	sslmode := getEnv("TEST_DB_SSL_MODE", "disable")

	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &TestDB{DB: db}
}

// BeginTx starts a transaction for the test
func (tdb *TestDB) BeginTx(t *testing.T) {
	t.Helper()

	tx, err := tdb.DB.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	tdb.tx = tx
}

// Rollback rolls back the transaction
func (tdb *TestDB) Rollback(t *testing.T) {
	t.Helper()

	if tdb.tx != nil {
		if err := tdb.tx.Rollback(); err != nil {
			t.Logf("Failed to rollback transaction: %v", err)
		}
		tdb.tx = nil
	}
}

// Close closes the database connection
func (tdb *TestDB) Close() {
	if tdb.DB != nil {
		tdb.DB.Close()
	}
}

// CleanupTables truncates all tables for a clean test state
func (tdb *TestDB) CleanupTables(t *testing.T) {
	t.Helper()

	tables := []string{
		"reservations",
		"borrow_records",
		"book_copies",
		"books",
		"users",
	}

	for _, table := range tables {
		_, err := tdb.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Logf("Warning: Failed to truncate table %s: %v", table, err)
		}
	}
}

// RunMigrations runs database migrations for testing
func (tdb *TestDB) RunMigrations(t *testing.T) {
	t.Helper()

	// Create tables
	migrations := []string{
		// Books table
		`CREATE TABLE IF NOT EXISTS books (
			id VARCHAR(255) PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			author VARCHAR(255) NOT NULL,
			isbn VARCHAR(20) NOT NULL UNIQUE,
			category VARCHAR(100) NOT NULL,
			publication_year INTEGER NOT NULL,
			description TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)`,

		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(255) PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			full_name VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL,
			status VARCHAR(50) NOT NULL,
			borrowing_limit INTEGER NOT NULL DEFAULT 3,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,

		// Book copies table
		`CREATE TABLE IF NOT EXISTS book_copies (
			id VARCHAR(255) PRIMARY KEY,
			book_id VARCHAR(255) NOT NULL REFERENCES books(id) ON DELETE CASCADE,
			copy_number INTEGER NOT NULL,
			status VARCHAR(50) NOT NULL,
			location VARCHAR(255),
			condition VARCHAR(50),
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(book_id, copy_number)
		)`,

		// Borrow records table
		`CREATE TABLE IF NOT EXISTS borrow_records (
			id VARCHAR(255) PRIMARY KEY,
			user_id VARCHAR(255) NOT NULL REFERENCES users(id),
			book_copy_id VARCHAR(255) NOT NULL REFERENCES book_copies(id),
			checkout_date TIMESTAMP NOT NULL,
			due_date TIMESTAMP NOT NULL,
			return_date TIMESTAMP,
			renewal_count INTEGER NOT NULL DEFAULT 0,
			late_fee DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
			status VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,

		// Reservations table
		`CREATE TABLE IF NOT EXISTS reservations (
			id VARCHAR(255) PRIMARY KEY,
			user_id VARCHAR(255) NOT NULL REFERENCES users(id),
			book_id VARCHAR(255) NOT NULL REFERENCES books(id),
			reservation_date TIMESTAMP NOT NULL,
			expiry_date TIMESTAMP NOT NULL,
			status VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,

		// Create indexes
		`CREATE INDEX IF NOT EXISTS idx_books_isbn ON books(isbn)`,
		`CREATE INDEX IF NOT EXISTS idx_books_category ON books(category)`,
		`CREATE INDEX IF NOT EXISTS idx_books_author ON books(author)`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		`CREATE INDEX IF NOT EXISTS idx_book_copies_book_id ON book_copies(book_id)`,
		`CREATE INDEX IF NOT EXISTS idx_book_copies_status ON book_copies(status)`,
		`CREATE INDEX IF NOT EXISTS idx_borrow_records_user_id ON borrow_records(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_borrow_records_status ON borrow_records(status)`,
		`CREATE INDEX IF NOT EXISTS idx_reservations_user_id ON reservations(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_reservations_book_id ON reservations(book_id)`,
	}

	for i, migration := range migrations {
		if _, err := tdb.DB.Exec(migration); err != nil {
			t.Fatalf("Failed to run migration %d: %v\nMigration: %s", i+1, err, migration)
		}
	}
}

// DropTables drops all tables (for cleanup)
func (tdb *TestDB) DropTables(t *testing.T) {
	t.Helper()

	tables := []string{
		"reservations",
		"borrow_records",
		"book_copies",
		"books",
		"users",
	}

	for _, table := range tables {
		_, err := tdb.DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
		if err != nil {
			t.Logf("Warning: Failed to drop table %s: %v", table, err)
		}
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}