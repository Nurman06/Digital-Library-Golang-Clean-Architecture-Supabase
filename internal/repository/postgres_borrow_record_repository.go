package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
)

// PostgresBorrowRecordRepository implements BorrowRecordRepository interface using PostgreSQL
type PostgresBorrowRecordRepository struct {
	db *sql.DB
}

// NewPostgresBorrowRecordRepository creates a new PostgresBorrowRecordRepository
func NewPostgresBorrowRecordRepository(db *sql.DB) *PostgresBorrowRecordRepository {
	return &PostgresBorrowRecordRepository{
		db: db,
	}
}

// Create creates a new borrow record in the repository
func (r *PostgresBorrowRecordRepository) Create(ctx context.Context, record *entity.BorrowRecord) error {
	query := `
		INSERT INTO borrow_records (id, user_id, book_copy_id, checkout_date, due_date, return_date, renewal_count, late_fee, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	now := time.Now()
	record.CreatedAt = now
	record.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		record.ID,
		record.UserID,
		record.BookCopyID,
		record.CheckoutDate,
		record.DueDate,
		record.ReturnDate,
		record.RenewalCount,
		record.LateFee,
		record.Status,
		record.CreatedAt,
		record.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create borrow record: %w", err)
	}

	return nil
}

// GetByID retrieves a borrow record by its ID
func (r *PostgresBorrowRecordRepository) GetByID(ctx context.Context, id string) (*entity.BorrowRecord, error) {
	query := `
		SELECT id, user_id, book_copy_id, checkout_date, due_date, return_date, renewal_count, late_fee, status, created_at, updated_at
		FROM borrow_records
		WHERE id = $1
	`

	record := &entity.BorrowRecord{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&record.ID,
		&record.UserID,
		&record.BookCopyID,
		&record.CheckoutDate,
		&record.DueDate,
		&record.ReturnDate,
		&record.RenewalCount,
		&record.LateFee,
		&record.Status,
		&record.CreatedAt,
		&record.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("borrow record not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get borrow record: %w", err)
	}

	return record, nil
}

// Update updates an existing borrow record
func (r *PostgresBorrowRecordRepository) Update(ctx context.Context, record *entity.BorrowRecord) error {
	query := `
		UPDATE borrow_records
		SET user_id = $1, book_copy_id = $2, checkout_date = $3, due_date = $4, return_date = $5, 
		    renewal_count = $6, late_fee = $7, status = $8, updated_at = $9
		WHERE id = $10
	`

	record.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		record.UserID,
		record.BookCopyID,
		record.CheckoutDate,
		record.DueDate,
		record.ReturnDate,
		record.RenewalCount,
		record.LateFee,
		record.Status,
		record.UpdatedAt,
		record.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update borrow record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("borrow record not found")
	}

	return nil
}

// Delete deletes a borrow record from the repository
func (r *PostgresBorrowRecordRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM borrow_records WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete borrow record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("borrow record not found")
	}

	return nil
}

// GetByUserID retrieves all borrow records for a specific user
func (r *PostgresBorrowRecordRepository) GetByUserID(ctx context.Context, userID string, params BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM borrow_records WHERE user_id = $1`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count borrow records: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "checkout_date DESC"
	if params.SortBy != "" {
		order := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", params.SortBy, order)
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Query borrow records
	query := fmt.Sprintf(`
		SELECT id, user_id, book_copy_id, checkout_date, due_date, return_date, renewal_count, late_fee, status, created_at, updated_at
		FROM borrow_records
		WHERE user_id = $1
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, userID, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get borrow records: %w", err)
	}
	defer rows.Close()

	records := make([]*entity.BorrowRecord, 0)
	for rows.Next() {
		record := &entity.BorrowRecord{}
		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.BookCopyID,
			&record.CheckoutDate,
			&record.DueDate,
			&record.ReturnDate,
			&record.RenewalCount,
			&record.LateFee,
			&record.Status,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan borrow record: %w", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating borrow records: %w", err)
	}

	return records, total, nil
}

// GetByBookCopyID retrieves all borrow records for a specific book copy
func (r *PostgresBorrowRecordRepository) GetByBookCopyID(ctx context.Context, bookCopyID string, params BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM borrow_records WHERE book_copy_id = $1`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, bookCopyID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count borrow records: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "checkout_date DESC"
	if params.SortBy != "" {
		order := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", params.SortBy, order)
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Query borrow records
	query := fmt.Sprintf(`
		SELECT id, user_id, book_copy_id, checkout_date, due_date, return_date, renewal_count, late_fee, status, created_at, updated_at
		FROM borrow_records
		WHERE book_copy_id = $1
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, bookCopyID, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get borrow records: %w", err)
	}
	defer rows.Close()

	records := make([]*entity.BorrowRecord, 0)
	for rows.Next() {
		record := &entity.BorrowRecord{}
		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.BookCopyID,
			&record.CheckoutDate,
			&record.DueDate,
			&record.ReturnDate,
			&record.RenewalCount,
			&record.LateFee,
			&record.Status,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan borrow record: %w", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating borrow records: %w", err)
	}

	return records, total, nil
}

// GetByStatus retrieves borrow records by status
func (r *PostgresBorrowRecordRepository) GetByStatus(ctx context.Context, status entity.BorrowRecordStatus, params BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM borrow_records WHERE status = $1`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, status).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count borrow records by status: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "checkout_date DESC"
	if params.SortBy != "" {
		order := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", params.SortBy, order)
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Query borrow records
	query := fmt.Sprintf(`
		SELECT id, user_id, book_copy_id, checkout_date, due_date, return_date, renewal_count, late_fee, status, created_at, updated_at
		FROM borrow_records
		WHERE status = $1
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, status, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get borrow records by status: %w", err)
	}
	defer rows.Close()

	records := make([]*entity.BorrowRecord, 0)
	for rows.Next() {
		record := &entity.BorrowRecord{}
		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.BookCopyID,
			&record.CheckoutDate,
			&record.DueDate,
			&record.ReturnDate,
			&record.RenewalCount,
			&record.LateFee,
			&record.Status,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan borrow record: %w", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating borrow records: %w", err)
	}

	return records, total, nil
}

// GetActiveByUserID retrieves active borrow records for a specific user
func (r *PostgresBorrowRecordRepository) GetActiveByUserID(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
	query := `
		SELECT id, user_id, book_copy_id, checkout_date, due_date, return_date, renewal_count, late_fee, status, created_at, updated_at
		FROM borrow_records
		WHERE user_id = $1 AND status = $2
		ORDER BY checkout_date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, entity.BorrowStatusActive)
	if err != nil {
		return nil, fmt.Errorf("failed to get active borrow records: %w", err)
	}
	defer rows.Close()

	records := make([]*entity.BorrowRecord, 0)
	for rows.Next() {
		record := &entity.BorrowRecord{}
		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.BookCopyID,
			&record.CheckoutDate,
			&record.DueDate,
			&record.ReturnDate,
			&record.RenewalCount,
			&record.LateFee,
			&record.Status,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan borrow record: %w", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating borrow records: %w", err)
	}

	return records, nil
}

// GetOverdueRecords retrieves all overdue borrow records
func (r *PostgresBorrowRecordRepository) GetOverdueRecords(ctx context.Context, params BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM borrow_records WHERE status = $1 AND due_date < $2`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, entity.BorrowStatusActive, time.Now()).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count overdue records: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "due_date ASC"
	if params.SortBy != "" {
		order := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", params.SortBy, order)
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Query borrow records
	query := fmt.Sprintf(`
		SELECT id, user_id, book_copy_id, checkout_date, due_date, return_date, renewal_count, late_fee, status, created_at, updated_at
		FROM borrow_records
		WHERE status = $1 AND due_date < $2
		ORDER BY %s
		LIMIT $3 OFFSET $4
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, entity.BorrowStatusActive, time.Now(), params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get overdue records: %w", err)
	}
	defer rows.Close()

	records := make([]*entity.BorrowRecord, 0)
	for rows.Next() {
		record := &entity.BorrowRecord{}
		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.BookCopyID,
			&record.CheckoutDate,
			&record.DueDate,
			&record.ReturnDate,
			&record.RenewalCount,
			&record.LateFee,
			&record.Status,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan borrow record: %w", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating borrow records: %w", err)
	}

	return records, total, nil
}

// GetOverdueByUserID retrieves overdue borrow records for a specific user
func (r *PostgresBorrowRecordRepository) GetOverdueByUserID(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
	query := `
		SELECT id, user_id, book_copy_id, checkout_date, due_date, return_date, renewal_count, late_fee, status, created_at, updated_at
		FROM borrow_records
		WHERE user_id = $1 AND status = $2 AND due_date < $3
		ORDER BY due_date ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, entity.BorrowStatusActive, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue records: %w", err)
	}
	defer rows.Close()

	records := make([]*entity.BorrowRecord, 0)
	for rows.Next() {
		record := &entity.BorrowRecord{}
		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.BookCopyID,
			&record.CheckoutDate,
			&record.DueDate,
			&record.ReturnDate,
			&record.RenewalCount,
			&record.LateFee,
			&record.Status,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan borrow record: %w", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating borrow records: %w", err)
	}

	return records, nil
}

// GetDueSoon retrieves borrow records that are due within the specified days
func (r *PostgresBorrowRecordRepository) GetDueSoon(ctx context.Context, days int, params BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	dueDate := time.Now().AddDate(0, 0, days)

	// Count total
	countQuery := `SELECT COUNT(*) FROM borrow_records WHERE status = $1 AND due_date <= $2 AND due_date >= $3`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, entity.BorrowStatusActive, dueDate, time.Now()).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count due soon records: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "due_date ASC"
	if params.SortBy != "" {
		order := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", params.SortBy, order)
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Query borrow records
	query := fmt.Sprintf(`
		SELECT id, user_id, book_copy_id, checkout_date, due_date, return_date, renewal_count, late_fee, status, created_at, updated_at
		FROM borrow_records
		WHERE status = $1 AND due_date <= $2 AND due_date >= $3
		ORDER BY %s
		LIMIT $4 OFFSET $5
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, entity.BorrowStatusActive, dueDate, time.Now(), params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get due soon records: %w", err)
	}
	defer rows.Close()

	records := make([]*entity.BorrowRecord, 0)
	for rows.Next() {
		record := &entity.BorrowRecord{}
		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.BookCopyID,
			&record.CheckoutDate,
			&record.DueDate,
			&record.ReturnDate,
			&record.RenewalCount,
			&record.LateFee,
			&record.Status,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan borrow record: %w", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating borrow records: %w", err)
	}

	return records, total, nil
}

// List retrieves a paginated list of borrow records
func (r *PostgresBorrowRecordRepository) List(ctx context.Context, params BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM borrow_records`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count borrow records: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "checkout_date DESC"
	if params.SortBy != "" {
		order := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", params.SortBy, order)
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Query borrow records
	query := fmt.Sprintf(`
		SELECT id, user_id, book_copy_id, checkout_date, due_date, return_date, renewal_count, late_fee, status, created_at, updated_at
		FROM borrow_records
		ORDER BY %s
		LIMIT $1 OFFSET $2
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list borrow records: %w", err)
	}
	defer rows.Close()

	records := make([]*entity.BorrowRecord, 0)
	for rows.Next() {
		record := &entity.BorrowRecord{}
		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.BookCopyID,
			&record.CheckoutDate,
			&record.DueDate,
			&record.ReturnDate,
			&record.RenewalCount,
			&record.LateFee,
			&record.Status,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan borrow record: %w", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating borrow records: %w", err)
	}

	return records, total, nil
}

// CountActiveByUserID counts active borrow records for a specific user
func (r *PostgresBorrowRecordRepository) CountActiveByUserID(ctx context.Context, userID string) (int64, error) {
	query := `SELECT COUNT(*) FROM borrow_records WHERE user_id = $1 AND status = $2`

	var count int64
	err := r.db.QueryRowContext(ctx, query, userID, entity.BorrowStatusActive).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count active borrow records: %w", err)
	}

	return count, nil
}

// CountOverdueByUserID counts overdue borrow records for a specific user
func (r *PostgresBorrowRecordRepository) CountOverdueByUserID(ctx context.Context, userID string) (int64, error) {
	query := `SELECT COUNT(*) FROM borrow_records WHERE user_id = $1 AND status = $2 AND due_date < $3`

	var count int64
	err := r.db.QueryRowContext(ctx, query, userID, entity.BorrowStatusActive, time.Now()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count overdue borrow records: %w", err)
	}

	return count, nil
}

// GetByDateRange retrieves borrow records within a date range
func (r *PostgresBorrowRecordRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time, params BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM borrow_records WHERE checkout_date >= $1 AND checkout_date <= $2`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, startDate, endDate).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count borrow records by date range: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "checkout_date DESC"
	if params.SortBy != "" {
		order := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", params.SortBy, order)
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Query borrow records
	query := fmt.Sprintf(`
		SELECT id, user_id, book_copy_id, checkout_date, due_date, return_date, renewal_count, late_fee, status, created_at, updated_at
		FROM borrow_records
		WHERE checkout_date >= $1 AND checkout_date <= $2
		ORDER BY %s
		LIMIT $3 OFFSET $4
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get borrow records by date range: %w", err)
	}
	defer rows.Close()

	records := make([]*entity.BorrowRecord, 0)
	for rows.Next() {
		record := &entity.BorrowRecord{}
		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.BookCopyID,
			&record.CheckoutDate,
			&record.DueDate,
			&record.ReturnDate,
			&record.RenewalCount,
			&record.LateFee,
			&record.Status,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan borrow record: %w", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating borrow records: %w", err)
	}

	return records, total, nil
}

// GetUserBorrowingHistory retrieves complete borrowing history for a user
func (r *PostgresBorrowRecordRepository) GetUserBorrowingHistory(ctx context.Context, userID string, params BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	return r.GetByUserID(ctx, userID, params)
}

// GetBookBorrowingHistory retrieves complete borrowing history for a book copy
func (r *PostgresBorrowRecordRepository) GetBookBorrowingHistory(ctx context.Context, bookCopyID string, params BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	return r.GetByBookCopyID(ctx, bookCopyID, params)
}

// UpdateStatus updates the status of a borrow record
func (r *PostgresBorrowRecordRepository) UpdateStatus(ctx context.Context, id string, status entity.BorrowRecordStatus) error {
	query := `
		UPDATE borrow_records
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, status, now, id)
	if err != nil {
		return fmt.Errorf("failed to update borrow record status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("borrow record not found")
	}

	return nil
}

// MarkAsReturned marks a borrow record as returned
func (r *PostgresBorrowRecordRepository) MarkAsReturned(ctx context.Context, id string, returnDate time.Time) error {
	query := `
		UPDATE borrow_records
		SET return_date = $1, status = $2, updated_at = $3
		WHERE id = $4
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, returnDate, entity.BorrowStatusReturned, now, id)
	if err != nil {
		return fmt.Errorf("failed to mark borrow record as returned: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("borrow record not found")
	}

	return nil
}

// CalculateTotalLateFees calculates total unpaid late fees for a user
func (r *PostgresBorrowRecordRepository) CalculateTotalLateFees(ctx context.Context, userID string) (float64, error) {
	query := `SELECT COALESCE(SUM(late_fee), 0) FROM borrow_records WHERE user_id = $1 AND status != $2`

	var totalFees float64
	err := r.db.QueryRowContext(ctx, query, userID, entity.BorrowStatusReturned).Scan(&totalFees)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total late fees: %w", err)
	}

	return totalFees, nil
}

// GetMostBorrowedBooks retrieves statistics of most borrowed books
func (r *PostgresBorrowRecordRepository) GetMostBorrowedBooks(ctx context.Context, limit int, startDate, endDate time.Time) ([]BorrowStatistic, error) {
	query := `
		SELECT 
			bc.book_id,
			b.title,
			COUNT(*) as borrow_count,
			COALESCE(SUM(EXTRACT(EPOCH FROM (COALESCE(br.return_date, NOW()) - br.checkout_date)) / 86400), 0) as total_days
		FROM borrow_records br
		JOIN book_copies bc ON br.book_copy_id = bc.id
		JOIN books b ON bc.book_id = b.id
		WHERE br.checkout_date >= $1 AND br.checkout_date <= $2
		GROUP BY bc.book_id, b.title
		ORDER BY borrow_count DESC
		LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get most borrowed books: %w", err)
	}
	defer rows.Close()

	statistics := make([]BorrowStatistic, 0)
	for rows.Next() {
		stat := BorrowStatistic{}
		err := rows.Scan(
			&stat.BookID,
			&stat.BookTitle,
			&stat.BorrowCount,
			&stat.TotalDays,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan borrow statistic: %w", err)
		}
		statistics = append(statistics, stat)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating borrow statistics: %w", err)
	}

	return statistics, nil
}

// GetUserBorrowingStatistics retrieves borrowing statistics for a user
func (r *PostgresBorrowRecordRepository) GetUserBorrowingStatistics(ctx context.Context, userID string) (*UserBorrowStatistics, error) {
	query := `
		SELECT 
			COUNT(*) as total_borrows,
			COUNT(*) FILTER (WHERE status = $2) as active_borrows,
			COUNT(*) FILTER (WHERE status = $2 AND due_date < $3) as overdue_borrows,
			COUNT(*) FILTER (WHERE status = $4) as returned_borrows,
			COALESCE(SUM(late_fee), 0) as total_late_fees,
			COALESCE(AVG(EXTRACT(EPOCH FROM (COALESCE(return_date, NOW()) - checkout_date)) / 86400), 0) as average_days
		FROM borrow_records
		WHERE user_id = $1
	`

	stats := &UserBorrowStatistics{}
	err := r.db.QueryRowContext(ctx, query, userID, entity.BorrowStatusActive, time.Now(), entity.BorrowStatusReturned).Scan(
		&stats.TotalBorrows,
		&stats.ActiveBorrows,
		&stats.OverdueBorrows,
		&stats.ReturnedBorrows,
		&stats.TotalLateFees,
		&stats.AverageDays,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user borrowing statistics: %w", err)
	}

	return stats, nil
}

// HasActiveBookCopyBorrow checks if a book copy has an active borrow record
func (r *PostgresBorrowRecordRepository) HasActiveBookCopyBorrow(ctx context.Context, bookCopyID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM borrow_records WHERE book_copy_id = $1 AND status = $2)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, bookCopyID, entity.BorrowStatusActive).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check active borrow: %w", err)
	}

	return exists, nil
}

// GetActiveByBookCopyID retrieves active borrow record for a specific book copy
func (r *PostgresBorrowRecordRepository) GetActiveByBookCopyID(ctx context.Context, bookCopyID string) (*entity.BorrowRecord, error) {
	query := `
		SELECT id, user_id, book_copy_id, checkout_date, due_date, return_date, renewal_count, late_fee, status, created_at, updated_at
		FROM borrow_records
		WHERE book_copy_id = $1 AND status = $2
		LIMIT 1
	`

	record := &entity.BorrowRecord{}
	err := r.db.QueryRowContext(ctx, query, bookCopyID, entity.BorrowStatusActive).Scan(
		&record.ID,
		&record.UserID,
		&record.BookCopyID,
		&record.CheckoutDate,
		&record.DueDate,
		&record.ReturnDate,
		&record.RenewalCount,
		&record.LateFee,
		&record.Status,
		&record.CreatedAt,
		&record.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no active borrow record found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get active borrow record: %w", err)
	}

	return record, nil
}