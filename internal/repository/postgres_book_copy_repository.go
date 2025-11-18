package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
)

// PostgresBookCopyRepository implements BookCopyRepository interface using PostgreSQL
type PostgresBookCopyRepository struct {
	db *sql.DB
}

// NewPostgresBookCopyRepository creates a new PostgresBookCopyRepository
func NewPostgresBookCopyRepository(db *sql.DB) *PostgresBookCopyRepository {
	return &PostgresBookCopyRepository{
		db: db,
	}
}

// Create creates a new book copy in the repository
func (r *PostgresBookCopyRepository) Create(ctx context.Context, copy *entity.BookCopy) error {
	query := `
		INSERT INTO book_copies (id, book_id, copy_number, status, location, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()
	copy.CreatedAt = now
	copy.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		copy.ID,
		copy.BookID,
		copy.CopyNumber,
		copy.Status,
		copy.Location,
		copy.CreatedAt,
		copy.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create book copy: %w", err)
	}

	return nil
}

// GetByID retrieves a book copy by its ID
func (r *PostgresBookCopyRepository) GetByID(ctx context.Context, id string) (*entity.BookCopy, error) {
	query := `
		SELECT id, book_id, copy_number, status, location, created_at, updated_at
		FROM book_copies
		WHERE id = $1
	`

	copy := &entity.BookCopy{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&copy.ID,
		&copy.BookID,
		&copy.CopyNumber,
		&copy.Status,
		&copy.Location,
		&copy.CreatedAt,
		&copy.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("book copy not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get book copy: %w", err)
	}

	return copy, nil
}

// GetByCopyNumber retrieves a book copy by its copy number
func (r *PostgresBookCopyRepository) GetByCopyNumber(ctx context.Context, copyNumber string) (*entity.BookCopy, error) {
	query := `
		SELECT id, book_id, copy_number, status, location, created_at, updated_at
		FROM book_copies
		WHERE copy_number = $1
	`

	copy := &entity.BookCopy{}
	err := r.db.QueryRowContext(ctx, query, copyNumber).Scan(
		&copy.ID,
		&copy.BookID,
		&copy.CopyNumber,
		&copy.Status,
		&copy.Location,
		&copy.CreatedAt,
		&copy.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("book copy not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get book copy: %w", err)
	}

	return copy, nil
}

// Update updates an existing book copy
func (r *PostgresBookCopyRepository) Update(ctx context.Context, copy *entity.BookCopy) error {
	query := `
		UPDATE book_copies
		SET book_id = $1, copy_number = $2, status = $3, location = $4, updated_at = $5
		WHERE id = $6
	`

	copy.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		copy.BookID,
		copy.CopyNumber,
		copy.Status,
		copy.Location,
		copy.UpdatedAt,
		copy.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update book copy: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book copy not found")
	}

	return nil
}

// Delete deletes a book copy from the repository
func (r *PostgresBookCopyRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM book_copies WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete book copy: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book copy not found")
	}

	return nil
}

// GetByBookID retrieves all copies of a specific book
func (r *PostgresBookCopyRepository) GetByBookID(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
	query := `
		SELECT id, book_id, copy_number, status, location, created_at, updated_at
		FROM book_copies
		WHERE book_id = $1
		ORDER BY copy_number ASC
	`

	rows, err := r.db.QueryContext(ctx, query, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book copies: %w", err)
	}
	defer rows.Close()

	copies := make([]*entity.BookCopy, 0)
	for rows.Next() {
		copy := &entity.BookCopy{}
		err := rows.Scan(
			&copy.ID,
			&copy.BookID,
			&copy.CopyNumber,
			&copy.Status,
			&copy.Location,
			&copy.CreatedAt,
			&copy.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan book copy: %w", err)
		}
		copies = append(copies, copy)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating book copies: %w", err)
	}

	return copies, nil
}

// GetAvailableCopies retrieves available copies of a specific book
func (r *PostgresBookCopyRepository) GetAvailableCopies(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
	query := `
		SELECT id, book_id, copy_number, status, location, created_at, updated_at
		FROM book_copies
		WHERE book_id = $1 AND status = $2
		ORDER BY copy_number ASC
	`

	rows, err := r.db.QueryContext(ctx, query, bookID, entity.CopyStatusAvailable)
	if err != nil {
		return nil, fmt.Errorf("failed to get available book copies: %w", err)
	}
	defer rows.Close()

	copies := make([]*entity.BookCopy, 0)
	for rows.Next() {
		copy := &entity.BookCopy{}
		err := rows.Scan(
			&copy.ID,
			&copy.BookID,
			&copy.CopyNumber,
			&copy.Status,
			&copy.Location,
			&copy.CreatedAt,
			&copy.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan book copy: %w", err)
		}
		copies = append(copies, copy)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating book copies: %w", err)
	}

	return copies, nil
}

// GetByStatus retrieves book copies by status
func (r *PostgresBookCopyRepository) GetByStatus(ctx context.Context, status entity.CopyStatus, params ListParams) ([]*entity.BookCopy, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM book_copies WHERE status = $1`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, status).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count book copies by status: %w", err)
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

	// Query book copies
	query := fmt.Sprintf(`
		SELECT id, book_id, copy_number, status, location, created_at, updated_at
		FROM book_copies
		WHERE status = $1
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, status, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get book copies by status: %w", err)
	}
	defer rows.Close()

	copies := make([]*entity.BookCopy, 0)
	for rows.Next() {
		copy := &entity.BookCopy{}
		err := rows.Scan(
			&copy.ID,
			&copy.BookID,
			&copy.CopyNumber,
			&copy.Status,
			&copy.Location,
			&copy.CreatedAt,
			&copy.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan book copy: %w", err)
		}
		copies = append(copies, copy)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating book copies: %w", err)
	}

	return copies, total, nil
}

// CountByBookID counts the total number of copies for a book
func (r *PostgresBookCopyRepository) CountByBookID(ctx context.Context, bookID string) (int64, error) {
	query := `SELECT COUNT(*) FROM book_copies WHERE book_id = $1`

	var count int64
	err := r.db.QueryRowContext(ctx, query, bookID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count book copies: %w", err)
	}

	return count, nil
}

// CountAvailableByBookID counts available copies for a book
func (r *PostgresBookCopyRepository) CountAvailableByBookID(ctx context.Context, bookID string) (int64, error) {
	query := `SELECT COUNT(*) FROM book_copies WHERE book_id = $1 AND status = $2`

	var count int64
	err := r.db.QueryRowContext(ctx, query, bookID, entity.CopyStatusAvailable).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count available book copies: %w", err)
	}

	return count, nil
}

// UpdateStatus updates the status of a book copy
func (r *PostgresBookCopyRepository) UpdateStatus(ctx context.Context, id string, status entity.CopyStatus) error {
	query := `
		UPDATE book_copies
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, status, now, id)
	if err != nil {
		return fmt.Errorf("failed to update book copy status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book copy not found")
	}

	return nil
}

// ExistsByCopyNumber checks if a copy with the given copy number exists
func (r *PostgresBookCopyRepository) ExistsByCopyNumber(ctx context.Context, copyNumber string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM book_copies WHERE copy_number = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, copyNumber).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check book copy existence: %w", err)
	}

	return exists, nil
}

// GetByLocation retrieves book copies by location
func (r *PostgresBookCopyRepository) GetByLocation(ctx context.Context, location string, params ListParams) ([]*entity.BookCopy, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM book_copies WHERE location = $1`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, location).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count book copies by location: %w", err)
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

	// Query book copies
	query := fmt.Sprintf(`
		SELECT id, book_id, copy_number, status, location, created_at, updated_at
		FROM book_copies
		WHERE location = $1
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, location, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get book copies by location: %w", err)
	}
	defer rows.Close()

	copies := make([]*entity.BookCopy, 0)
	for rows.Next() {
		copy := &entity.BookCopy{}
		err := rows.Scan(
			&copy.ID,
			&copy.BookID,
			&copy.CopyNumber,
			&copy.Status,
			&copy.Location,
			&copy.CreatedAt,
			&copy.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan book copy: %w", err)
		}
		copies = append(copies, copy)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating book copies: %w", err)
	}

	return copies, total, nil
}

// List retrieves a paginated list of book copies
func (r *PostgresBookCopyRepository) List(ctx context.Context, params ListParams) ([]*entity.BookCopy, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM book_copies`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count book copies: %w", err)
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

	// Query book copies
	query := fmt.Sprintf(`
		SELECT id, book_id, copy_number, status, location, created_at, updated_at
		FROM book_copies
		ORDER BY %s
		LIMIT $1 OFFSET $2
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list book copies: %w", err)
	}
	defer rows.Close()

	copies := make([]*entity.BookCopy, 0)
	for rows.Next() {
		copy := &entity.BookCopy{}
		err := rows.Scan(
			&copy.ID,
			&copy.BookID,
			&copy.CopyNumber,
			&copy.Status,
			&copy.Location,
			&copy.CreatedAt,
			&copy.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan book copy: %w", err)
		}
		copies = append(copies, copy)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating book copies: %w", err)
	}

	return copies, total, nil
}