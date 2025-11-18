package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
)

// PostgresBookRepository implements BookRepository interface using PostgreSQL
type PostgresBookRepository struct {
	db *sql.DB
}

// NewPostgresBookRepository creates a new PostgresBookRepository
func NewPostgresBookRepository(db *sql.DB) *PostgresBookRepository {
	return &PostgresBookRepository{
		db: db,
	}
}

// Create creates a new book in the repository
func (r *PostgresBookRepository) Create(ctx context.Context, book *entity.Book) error {
	query := `
		INSERT INTO books (id, title, author, isbn, category, publication_year, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	now := time.Now()
	book.CreatedAt = now
	book.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		book.ID,
		book.Title,
		book.Author,
		book.ISBN,
		book.Category,
		book.PublicationYear,
		book.Description,
		book.CreatedAt,
		book.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}

	return nil
}

// GetByID retrieves a book by its ID
func (r *PostgresBookRepository) GetByID(ctx context.Context, id string) (*entity.Book, error) {
	query := `
		SELECT id, title, author, isbn, category, publication_year, description, created_at, updated_at, deleted_at
		FROM books
		WHERE id = $1 AND deleted_at IS NULL
	`

	book := &entity.Book{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.Category,
		&book.PublicationYear,
		&book.Description,
		&book.CreatedAt,
		&book.UpdatedAt,
		&book.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("book not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	return book, nil
}

// GetByISBN retrieves a book by its ISBN
func (r *PostgresBookRepository) GetByISBN(ctx context.Context, isbn string) (*entity.Book, error) {
	query := `
		SELECT id, title, author, isbn, category, publication_year, description, created_at, updated_at, deleted_at
		FROM books
		WHERE isbn = $1 AND deleted_at IS NULL
	`

	book := &entity.Book{}
	err := r.db.QueryRowContext(ctx, query, isbn).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.Category,
		&book.PublicationYear,
		&book.Description,
		&book.CreatedAt,
		&book.UpdatedAt,
		&book.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("book not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	return book, nil
}

// Update updates an existing book
func (r *PostgresBookRepository) Update(ctx context.Context, book *entity.Book) error {
	query := `
		UPDATE books
		SET title = $1, author = $2, isbn = $3, category = $4, publication_year = $5, 
		    description = $6, updated_at = $7
		WHERE id = $8 AND deleted_at IS NULL
	`

	book.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		book.Title,
		book.Author,
		book.ISBN,
		book.Category,
		book.PublicationYear,
		book.Description,
		book.UpdatedAt,
		book.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}

// Delete performs a soft delete on a book
func (r *PostgresBookRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE books
		SET deleted_at = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}

// HardDelete permanently deletes a book from the repository
func (r *PostgresBookRepository) HardDelete(ctx context.Context, id string) error {
	query := `DELETE FROM books WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}

// Restore restores a soft-deleted book
func (r *PostgresBookRepository) Restore(ctx context.Context, id string) error {
	query := `
		UPDATE books
		SET deleted_at = NULL, updated_at = $1
		WHERE id = $2 AND deleted_at IS NOT NULL
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("failed to restore book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book not found or not deleted")
	}

	return nil
}

// List retrieves a paginated list of books
func (r *PostgresBookRepository) List(ctx context.Context, params ListParams) ([]*entity.Book, int64, error) {
	// Build WHERE clause
	whereClause := "WHERE deleted_at IS NULL"
	if params.IncludeDeleted {
		whereClause = ""
	}

	// Build ORDER BY clause
	orderBy := "created_at DESC"
	if params.SortBy != "" {
		order := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", params.SortBy, order)
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM books %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count books: %w", err)
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Query books
	query := fmt.Sprintf(`
		SELECT id, title, author, isbn, category, publication_year, description, created_at, updated_at, deleted_at
		FROM books
		%s
		ORDER BY %s
		LIMIT $1 OFFSET $2
	`, whereClause, orderBy)

	rows, err := r.db.QueryContext(ctx, query, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list books: %w", err)
	}
	defer rows.Close()

	books := make([]*entity.Book, 0)
	for rows.Next() {
		book := &entity.Book{}
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.Category,
			&book.PublicationYear,
			&book.Description,
			&book.CreatedAt,
			&book.UpdatedAt,
			&book.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan book: %w", err)
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating books: %w", err)
	}

	return books, total, nil
}

// Search searches for books based on search criteria
func (r *PostgresBookRepository) Search(ctx context.Context, params SearchParams) ([]*entity.Book, int64, error) {
	whereConditions := []string{"deleted_at IS NULL"}
	args := []interface{}{}
	argIndex := 1

	// Full-text search query
	if params.Query != "" {
		whereConditions = append(whereConditions, fmt.Sprintf(
			"(to_tsvector('english', title) @@ plainto_tsquery('english', $%d) OR to_tsvector('english', author) @@ plainto_tsquery('english', $%d))",
			argIndex, argIndex,
		))
		args = append(args, params.Query)
		argIndex++
	}

	// Title search
	if params.Title != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("title ILIKE $%d", argIndex))
		args = append(args, "%"+params.Title+"%")
		argIndex++
	}

	// Author search
	if params.Author != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("author ILIKE $%d", argIndex))
		args = append(args, "%"+params.Author+"%")
		argIndex++
	}

	// ISBN search
	if params.ISBN != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("isbn = $%d", argIndex))
		args = append(args, params.ISBN)
		argIndex++
	}

	// Category filter
	if params.Category != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, params.Category)
		argIndex++
	}

	// Publication year filter
	if params.PublicationYear > 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("publication_year = $%d", argIndex))
		args = append(args, params.PublicationYear)
		argIndex++
	}

	// Publication year range
	if params.YearFrom > 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("publication_year >= $%d", argIndex))
		args = append(args, params.YearFrom)
		argIndex++
	}
	if params.YearTo > 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("publication_year <= $%d", argIndex))
		args = append(args, params.YearTo)
		argIndex++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM books WHERE %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "created_at DESC"
	if params.SortBy != "" {
		order := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", params.SortBy, order)
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Query books
	query := fmt.Sprintf(`
		SELECT id, title, author, isbn, category, publication_year, description, created_at, updated_at, deleted_at
		FROM books
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, argIndex, argIndex+1)

	args = append(args, params.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search books: %w", err)
	}
	defer rows.Close()

	books := make([]*entity.Book, 0)
	for rows.Next() {
		book := &entity.Book{}
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.Category,
			&book.PublicationYear,
			&book.Description,
			&book.CreatedAt,
			&book.UpdatedAt,
			&book.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan book: %w", err)
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating search results: %w", err)
	}

	return books, total, nil
}

// GetByCategory retrieves books by category with pagination
func (r *PostgresBookRepository) GetByCategory(ctx context.Context, category string, params ListParams) ([]*entity.Book, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM books WHERE category = $1 AND deleted_at IS NULL`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, category).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count books by category: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "created_at DESC"
	if params.SortBy != "" {
		order := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", params.SortBy, order)
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Query books
	query := fmt.Sprintf(`
		SELECT id, title, author, isbn, category, publication_year, description, created_at, updated_at, deleted_at
		FROM books
		WHERE category = $1 AND deleted_at IS NULL
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, category, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get books by category: %w", err)
	}
	defer rows.Close()

	books := make([]*entity.Book, 0)
	for rows.Next() {
		book := &entity.Book{}
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.Category,
			&book.PublicationYear,
			&book.Description,
			&book.CreatedAt,
			&book.UpdatedAt,
			&book.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan book: %w", err)
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating books: %w", err)
	}

	return books, total, nil
}

// GetByAuthor retrieves books by author with pagination
func (r *PostgresBookRepository) GetByAuthor(ctx context.Context, author string, params ListParams) ([]*entity.Book, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM books WHERE author ILIKE $1 AND deleted_at IS NULL`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, "%"+author+"%").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count books by author: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "created_at DESC"
	if params.SortBy != "" {
		order := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", params.SortBy, order)
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Query books
	query := fmt.Sprintf(`
		SELECT id, title, author, isbn, category, publication_year, description, created_at, updated_at, deleted_at
		FROM books
		WHERE author ILIKE $1 AND deleted_at IS NULL
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, "%"+author+"%", params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get books by author: %w", err)
	}
	defer rows.Close()

	books := make([]*entity.Book, 0)
	for rows.Next() {
		book := &entity.Book{}
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.Category,
			&book.PublicationYear,
			&book.Description,
			&book.CreatedAt,
			&book.UpdatedAt,
			&book.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan book: %w", err)
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating books: %w", err)
	}

	return books, total, nil
}

// ExistsByISBN checks if a book with the given ISBN exists
func (r *PostgresBookRepository) ExistsByISBN(ctx context.Context, isbn string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM books WHERE isbn = $1 AND deleted_at IS NULL)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, isbn).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check book existence: %w", err)
	}

	return exists, nil
}

// Count returns the total number of books (excluding soft-deleted)
func (r *PostgresBookRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM books WHERE deleted_at IS NULL`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count books: %w", err)
	}

	return count, nil
}

// GetRecentlyAdded retrieves recently added books
func (r *PostgresBookRepository) GetRecentlyAdded(ctx context.Context, limit int) ([]*entity.Book, error) {
	query := `
		SELECT id, title, author, isbn, category, publication_year, description, created_at, updated_at, deleted_at
		FROM books
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recently added books: %w", err)
	}
	defer rows.Close()

	books := make([]*entity.Book, 0)
	for rows.Next() {
		book := &entity.Book{}
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.Category,
			&book.PublicationYear,
			&book.Description,
			&book.CreatedAt,
			&book.UpdatedAt,
			&book.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan book: %w", err)
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating books: %w", err)
	}

	return books, nil
}