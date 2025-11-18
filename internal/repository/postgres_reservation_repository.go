package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
)

// PostgresReservationRepository implements ReservationRepository interface using PostgreSQL
type PostgresReservationRepository struct {
	db *sql.DB
}

// NewPostgresReservationRepository creates a new PostgresReservationRepository
func NewPostgresReservationRepository(db *sql.DB) *PostgresReservationRepository {
	return &PostgresReservationRepository{
		db: db,
	}
}

// Create creates a new reservation in the repository
func (r *PostgresReservationRepository) Create(ctx context.Context, reservation *entity.Reservation) error {
	query := `
		INSERT INTO reservations (id, user_id, book_id, reservation_date, expiry_date, status, queue_position)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		reservation.ID,
		reservation.UserID,
		reservation.BookID,
		reservation.ReservationDate,
		reservation.ExpiryDate,
		reservation.Status,
		reservation.QueuePosition,
	)

	if err != nil {
		return fmt.Errorf("failed to create reservation: %w", err)
	}

	return nil
}

// GetByID retrieves a reservation by its ID
func (r *PostgresReservationRepository) GetByID(ctx context.Context, id string) (*entity.Reservation, error) {
	query := `
		SELECT id, user_id, book_id, reservation_date, expiry_date, status, queue_position
		FROM reservations
		WHERE id = $1
	`

	reservation := &entity.Reservation{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&reservation.ID,
		&reservation.UserID,
		&reservation.BookID,
		&reservation.ReservationDate,
		&reservation.ExpiryDate,
		&reservation.Status,
		&reservation.QueuePosition,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("reservation not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation: %w", err)
	}

	return reservation, nil
}

// Update updates an existing reservation
func (r *PostgresReservationRepository) Update(ctx context.Context, reservation *entity.Reservation) error {
	query := `
		UPDATE reservations
		SET user_id = $1, book_id = $2, reservation_date = $3, expiry_date = $4, status = $5, queue_position = $6
		WHERE id = $7
	`

	result, err := r.db.ExecContext(ctx, query,
		reservation.UserID,
		reservation.BookID,
		reservation.ReservationDate,
		reservation.ExpiryDate,
		reservation.Status,
		reservation.QueuePosition,
		reservation.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update reservation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("reservation not found")
	}

	return nil
}

// Delete deletes a reservation from the repository
func (r *PostgresReservationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM reservations WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete reservation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("reservation not found")
	}

	return nil
}

// GetByUserID retrieves all reservations for a specific user
func (r *PostgresReservationRepository) GetByUserID(ctx context.Context, userID string, params ReservationListParams) ([]*entity.Reservation, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM reservations WHERE user_id = $1`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reservations: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "reservation_date DESC"
	if params.SortBy != "" {
		order := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", params.SortBy, order)
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Query reservations
	query := fmt.Sprintf(`
		SELECT id, user_id, book_id, reservation_date, expiry_date, status, queue_position
		FROM reservations
		WHERE user_id = $1
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, userID, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get reservations: %w", err)
	}
	defer rows.Close()

	reservations := make([]*entity.Reservation, 0)
	for rows.Next() {
		reservation := &entity.Reservation{}
		err := rows.Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.BookID,
			&reservation.ReservationDate,
			&reservation.ExpiryDate,
			&reservation.Status,
			&reservation.QueuePosition,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan reservation: %w", err)
		}
		reservations = append(reservations, reservation)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating reservations: %w", err)
	}

	return reservations, total, nil
}

// GetByBookID retrieves all reservations for a specific book
func (r *PostgresReservationRepository) GetByBookID(ctx context.Context, bookID string, params ReservationListParams) ([]*entity.Reservation, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM reservations WHERE book_id = $1`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, bookID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reservations: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "queue_position ASC"
	if params.SortBy != "" {
		order := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", params.SortBy, order)
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Query reservations
	query := fmt.Sprintf(`
		SELECT id, user_id, book_id, reservation_date, expiry_date, status, queue_position
		FROM reservations
		WHERE book_id = $1
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, bookID, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get reservations: %w", err)
	}
	defer rows.Close()

	reservations := make([]*entity.Reservation, 0)
	for rows.Next() {
		reservation := &entity.Reservation{}
		err := rows.Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.BookID,
			&reservation.ReservationDate,
			&reservation.ExpiryDate,
			&reservation.Status,
			&reservation.QueuePosition,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan reservation: %w", err)
		}
		reservations = append(reservations, reservation)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating reservations: %w", err)
	}

	return reservations, total, nil
}

// GetByStatus retrieves reservations by status
func (r *PostgresReservationRepository) GetByStatus(ctx context.Context, status entity.ReservationStatus, params ReservationListParams) ([]*entity.Reservation, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM reservations WHERE status = $1`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, status).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reservations by status: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "reservation_date DESC"
	if params.SortBy != "" {
		order := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", params.SortBy, order)
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Query reservations
	query := fmt.Sprintf(`
		SELECT id, user_id, book_id, reservation_date, expiry_date, status, queue_position
		FROM reservations
		WHERE status = $1
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, status, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get reservations by status: %w", err)
	}
	defer rows.Close()

	reservations := make([]*entity.Reservation, 0)
	for rows.Next() {
		reservation := &entity.Reservation{}
		err := rows.Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.BookID,
			&reservation.ReservationDate,
			&reservation.ExpiryDate,
			&reservation.Status,
			&reservation.QueuePosition,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan reservation: %w", err)
		}
		reservations = append(reservations, reservation)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating reservations: %w", err)
	}

	return reservations, total, nil
}

// GetPendingByUserID retrieves pending reservations for a specific user
func (r *PostgresReservationRepository) GetPendingByUserID(ctx context.Context, userID string) ([]*entity.Reservation, error) {
	query := `
		SELECT id, user_id, book_id, reservation_date, expiry_date, status, queue_position
		FROM reservations
		WHERE user_id = $1 AND status = $2
		ORDER BY reservation_date ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, entity.ReservationStatusPending)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending reservations: %w", err)
	}
	defer rows.Close()

	reservations := make([]*entity.Reservation, 0)
	for rows.Next() {
		reservation := &entity.Reservation{}
		err := rows.Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.BookID,
			&reservation.ReservationDate,
			&reservation.ExpiryDate,
			&reservation.Status,
			&reservation.QueuePosition,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reservation: %w", err)
		}
		reservations = append(reservations, reservation)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reservations: %w", err)
	}

	return reservations, nil
}

// GetPendingByBookID retrieves pending reservations for a specific book ordered by queue position
func (r *PostgresReservationRepository) GetPendingByBookID(ctx context.Context, bookID string) ([]*entity.Reservation, error) {
	query := `
		SELECT id, user_id, book_id, reservation_date, expiry_date, status, queue_position
		FROM reservations
		WHERE book_id = $1 AND status = $2
		ORDER BY queue_position ASC
	`

	rows, err := r.db.QueryContext(ctx, query, bookID, entity.ReservationStatusPending)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending reservations: %w", err)
	}
	defer rows.Close()

	reservations := make([]*entity.Reservation, 0)
	for rows.Next() {
		reservation := &entity.Reservation{}
		err := rows.Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.BookID,
			&reservation.ReservationDate,
			&reservation.ExpiryDate,
			&reservation.Status,
			&reservation.QueuePosition,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reservation: %w", err)
		}
		reservations = append(reservations, reservation)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reservations: %w", err)
	}

	return reservations, nil
}

// GetExpiredReservations retrieves all expired reservations
func (r *PostgresReservationRepository) GetExpiredReservations(ctx context.Context) ([]*entity.Reservation, error) {
	query := `
		SELECT id, user_id, book_id, reservation_date, expiry_date, status, queue_position
		FROM reservations
		WHERE status = $1 AND expiry_date < $2
		ORDER BY expiry_date ASC
	`

	rows, err := r.db.QueryContext(ctx, query, entity.ReservationStatusPending, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to get expired reservations: %w", err)
	}
	defer rows.Close()

	reservations := make([]*entity.Reservation, 0)
	for rows.Next() {
		reservation := &entity.Reservation{}
		err := rows.Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.BookID,
			&reservation.ReservationDate,
			&reservation.ExpiryDate,
			&reservation.Status,
			&reservation.QueuePosition,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reservation: %w", err)
		}
		reservations = append(reservations, reservation)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reservations: %w", err)
	}

	return reservations, nil
}

// GetNextInQueue retrieves the next reservation in queue for a book
func (r *PostgresReservationRepository) GetNextInQueue(ctx context.Context, bookID string) (*entity.Reservation, error) {
	query := `
		SELECT id, user_id, book_id, reservation_date, expiry_date, status, queue_position
		FROM reservations
		WHERE book_id = $1 AND status = $2
		ORDER BY queue_position ASC
		LIMIT 1
	`

	reservation := &entity.Reservation{}
	err := r.db.QueryRowContext(ctx, query, bookID, entity.ReservationStatusPending).Scan(
		&reservation.ID,
		&reservation.UserID,
		&reservation.BookID,
		&reservation.ReservationDate,
		&reservation.ExpiryDate,
		&reservation.Status,
		&reservation.QueuePosition,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no pending reservations found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get next reservation in queue: %w", err)
	}

	return reservation, nil
}

// UpdateStatus updates the status of a reservation
func (r *PostgresReservationRepository) UpdateStatus(ctx context.Context, id string, status entity.ReservationStatus) error {
	query := `
		UPDATE reservations
		SET status = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update reservation status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("reservation not found")
	}

	return nil
}

// UpdateQueuePosition updates the queue position of a reservation
func (r *PostgresReservationRepository) UpdateQueuePosition(ctx context.Context, id string, position int) error {
	query := `
		UPDATE reservations
		SET queue_position = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, position, id)
	if err != nil {
		return fmt.Errorf("failed to update queue position: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("reservation not found")
	}

	return nil
}

// CountPendingByBookID counts pending reservations for a specific book
func (r *PostgresReservationRepository) CountPendingByBookID(ctx context.Context, bookID string) (int64, error) {
	query := `SELECT COUNT(*) FROM reservations WHERE book_id = $1 AND status = $2`

	var count int64
	err := r.db.QueryRowContext(ctx, query, bookID, entity.ReservationStatusPending).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count pending reservations: %w", err)
	}

	return count, nil
}

// CountPendingByUserID counts pending reservations for a specific user
func (r *PostgresReservationRepository) CountPendingByUserID(ctx context.Context, userID string) (int64, error) {
	query := `SELECT COUNT(*) FROM reservations WHERE user_id = $1 AND status = $2`

	var count int64
	err := r.db.QueryRowContext(ctx, query, userID, entity.ReservationStatusPending).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count pending reservations: %w", err)
	}

	return count, nil
}

// List retrieves a paginated list of reservations
func (r *PostgresReservationRepository) List(ctx context.Context, params ReservationListParams) ([]*entity.Reservation, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM reservations`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reservations: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "reservation_date DESC"
	if params.SortBy != "" {
		order := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", params.SortBy, order)
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Query reservations
	query := fmt.Sprintf(`
		SELECT id, user_id, book_id, reservation_date, expiry_date, status, queue_position
		FROM reservations
		ORDER BY %s
		LIMIT $1 OFFSET $2
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list reservations: %w", err)
	}
	defer rows.Close()

	reservations := make([]*entity.Reservation, 0)
	for rows.Next() {
		reservation := &entity.Reservation{}
		err := rows.Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.BookID,
			&reservation.ReservationDate,
			&reservation.ExpiryDate,
			&reservation.Status,
			&reservation.QueuePosition,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan reservation: %w", err)
		}
		reservations = append(reservations, reservation)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating reservations: %w", err)
	}

	return reservations, total, nil
}

// HasActiveReservation checks if a user has an active reservation for a book
func (r *PostgresReservationRepository) HasActiveReservation(ctx context.Context, userID, bookID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM reservations WHERE user_id = $1 AND book_id = $2 AND status = $3)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID, bookID, entity.ReservationStatusPending).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check active reservation: %w", err)
	}

	return exists, nil
}