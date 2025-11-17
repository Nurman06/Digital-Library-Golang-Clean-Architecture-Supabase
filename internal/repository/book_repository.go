package repository

import (
	"context"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
)

// BookRepository defines the interface for book data access operations
type BookRepository interface {
	// Create creates a new book in the repository
	Create(ctx context.Context, book *entity.Book) error

	// GetByID retrieves a book by its ID
	GetByID(ctx context.Context, id string) (*entity.Book, error)

	// GetByISBN retrieves a book by its ISBN
	GetByISBN(ctx context.Context, isbn string) (*entity.Book, error)

	// Update updates an existing book
	Update(ctx context.Context, book *entity.Book) error

	// Delete performs a soft delete on a book
	Delete(ctx context.Context, id string) error

	// HardDelete permanently deletes a book from the repository
	HardDelete(ctx context.Context, id string) error

	// Restore restores a soft-deleted book
	Restore(ctx context.Context, id string) error

	// List retrieves a paginated list of books
	List(ctx context.Context, params ListParams) ([]*entity.Book, int64, error)

	// Search searches for books based on search criteria
	Search(ctx context.Context, params SearchParams) ([]*entity.Book, int64, error)

	// GetByCategory retrieves books by category with pagination
	GetByCategory(ctx context.Context, category string, params ListParams) ([]*entity.Book, int64, error)

	// GetByAuthor retrieves books by author with pagination
	GetByAuthor(ctx context.Context, author string, params ListParams) ([]*entity.Book, int64, error)

	// ExistsByISBN checks if a book with the given ISBN exists
	ExistsByISBN(ctx context.Context, isbn string) (bool, error)

	// Count returns the total number of books (excluding soft-deleted)
	Count(ctx context.Context) (int64, error)

	// GetRecentlyAdded retrieves recently added books
	GetRecentlyAdded(ctx context.Context, limit int) ([]*entity.Book, error)
}

// ListParams defines parameters for listing books
type ListParams struct {
	Page           int
	PageSize       int
	SortBy         string // e.g., "title", "created_at", "publication_year"
	SortOrder      string // "asc" or "desc"
	IncludeDeleted bool
}

// SearchParams defines parameters for searching books
type SearchParams struct {
	Query           string // Full-text search query
	Title           string // Search by title
	Author          string // Search by author
	ISBN            string // Search by ISBN
	Category        string // Filter by category
	PublicationYear int    // Filter by publication year
	YearFrom        int    // Filter by publication year range (from)
	YearTo          int    // Filter by publication year range (to)
	Page            int
	PageSize        int
	SortBy          string
	SortOrder       string
}

// BookCopyRepository defines the interface for book copy data access operations
type BookCopyRepository interface {
	// Create creates a new book copy in the repository
	Create(ctx context.Context, copy *entity.BookCopy) error

	// GetByID retrieves a book copy by its ID
	GetByID(ctx context.Context, id string) (*entity.BookCopy, error)

	// GetByCopyNumber retrieves a book copy by its copy number
	GetByCopyNumber(ctx context.Context, copyNumber string) (*entity.BookCopy, error)

	// Update updates an existing book copy
	Update(ctx context.Context, copy *entity.BookCopy) error

	// Delete deletes a book copy from the repository
	Delete(ctx context.Context, id string) error

	// GetByBookID retrieves all copies of a specific book
	GetByBookID(ctx context.Context, bookID string) ([]*entity.BookCopy, error)

	// GetAvailableCopies retrieves available copies of a specific book
	GetAvailableCopies(ctx context.Context, bookID string) ([]*entity.BookCopy, error)

	// GetByStatus retrieves book copies by status
	GetByStatus(ctx context.Context, status entity.CopyStatus, params ListParams) ([]*entity.BookCopy, int64, error)

	// CountByBookID counts the total number of copies for a book
	CountByBookID(ctx context.Context, bookID string) (int64, error)

	// CountAvailableByBookID counts available copies for a book
	CountAvailableByBookID(ctx context.Context, bookID string) (int64, error)

	// UpdateStatus updates the status of a book copy
	UpdateStatus(ctx context.Context, id string, status entity.CopyStatus) error

	// ExistsByCopyNumber checks if a copy with the given copy number exists
	ExistsByCopyNumber(ctx context.Context, copyNumber string) (bool, error)

	// GetByLocation retrieves book copies by location
	GetByLocation(ctx context.Context, location string, params ListParams) ([]*entity.BookCopy, int64, error)

	// List retrieves a paginated list of book copies
	List(ctx context.Context, params ListParams) ([]*entity.BookCopy, int64, error)
}
