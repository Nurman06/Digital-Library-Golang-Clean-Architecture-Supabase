package supabase

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Client represents a Supabase client
type Client struct {
	DB  *sql.DB
	URL string
}

// Config holds Supabase client configuration
type Config struct {
	URL             string
	AnonKey         string
	ServiceKey      string
	DatabaseDSN     string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// NewClient creates a new Supabase client
func NewClient(config Config) (*Client, error) {
	// Connect to PostgreSQL database via Supabase
	db, err := sql.Open("postgres", config.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Client{
		DB:  db,
		URL: config.URL,
	}, nil
}

// Close closes the database connection
func (c *Client) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

// Health checks the health of the database connection
func (c *Client) Health() error {
	return c.DB.Ping()
}