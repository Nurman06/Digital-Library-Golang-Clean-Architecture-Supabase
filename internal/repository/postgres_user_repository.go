package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
)

// PostgresUserRepository implements UserRepository interface using PostgreSQL
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new PostgresUserRepository
func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

// Create creates a new user in the repository
func (r *PostgresUserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, full_name, role, status, borrowing_limit, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.Role,
		user.Status,
		user.BorrowingLimit,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by their ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, status, borrowing_limit, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.Status,
		&user.BorrowingLimit,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by their email address
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, status, borrowing_limit, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.Status,
		&user.BorrowingLimit,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Update updates an existing user
func (r *PostgresUserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET email = $1, full_name = $2, role = $3, status = $4, borrowing_limit = $5, updated_at = $6
		WHERE id = $7
	`

	user.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.FullName,
		user.Role,
		user.Status,
		user.BorrowingLimit,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete deletes a user from the repository
func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// List retrieves a paginated list of users
func (r *PostgresUserRepository) List(ctx context.Context, params UserListParams) ([]*entity.User, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM users`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
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

	// Query users
	query := fmt.Sprintf(`
		SELECT id, email, full_name, role, status, borrowing_limit, created_at, updated_at
		FROM users
		ORDER BY %s
		LIMIT $1 OFFSET $2
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	users := make([]*entity.User, 0)
	for rows.Next() {
		user := &entity.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FullName,
			&user.Role,
			&user.Status,
			&user.BorrowingLimit,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating users: %w", err)
	}

	return users, total, nil
}

// GetByRole retrieves users by role with pagination
func (r *PostgresUserRepository) GetByRole(ctx context.Context, role entity.UserRole, params UserListParams) ([]*entity.User, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM users WHERE role = $1`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, role).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users by role: %w", err)
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

	// Query users
	query := fmt.Sprintf(`
		SELECT id, email, full_name, role, status, borrowing_limit, created_at, updated_at
		FROM users
		WHERE role = $1
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, role, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users by role: %w", err)
	}
	defer rows.Close()

	users := make([]*entity.User, 0)
	for rows.Next() {
		user := &entity.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FullName,
			&user.Role,
			&user.Status,
			&user.BorrowingLimit,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating users: %w", err)
	}

	return users, total, nil
}

// GetByStatus retrieves users by status with pagination
func (r *PostgresUserRepository) GetByStatus(ctx context.Context, status entity.UserStatus, params UserListParams) ([]*entity.User, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM users WHERE status = $1`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, status).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users by status: %w", err)
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

	// Query users
	query := fmt.Sprintf(`
		SELECT id, email, full_name, role, status, borrowing_limit, created_at, updated_at
		FROM users
		WHERE status = $1
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, status, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users by status: %w", err)
	}
	defer rows.Close()

	users := make([]*entity.User, 0)
	for rows.Next() {
		user := &entity.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FullName,
			&user.Role,
			&user.Status,
			&user.BorrowingLimit,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating users: %w", err)
	}

	return users, total, nil
}

// ExistsByEmail checks if a user with the given email exists
func (r *PostgresUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

// UpdateStatus updates the status of a user
func (r *PostgresUserRepository) UpdateStatus(ctx context.Context, id string, status entity.UserStatus) error {
	query := `
		UPDATE users
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, status, now, id)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdateBorrowingLimit updates the borrowing limit of a user
func (r *PostgresUserRepository) UpdateBorrowingLimit(ctx context.Context, id string, limit int) error {
	query := `
		UPDATE users
		SET borrowing_limit = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, limit, now, id)
	if err != nil {
		return fmt.Errorf("failed to update borrowing limit: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Count returns the total number of users
func (r *PostgresUserRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// CountByRole returns the number of users by role
func (r *PostgresUserRepository) CountByRole(ctx context.Context, role entity.UserRole) (int64, error) {
	query := `SELECT COUNT(*) FROM users WHERE role = $1`

	var count int64
	err := r.db.QueryRowContext(ctx, query, role).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users by role: %w", err)
	}

	return count, nil
}

// CountByStatus returns the number of users by status
func (r *PostgresUserRepository) CountByStatus(ctx context.Context, status entity.UserStatus) (int64, error) {
	query := `SELECT COUNT(*) FROM users WHERE status = $1`

	var count int64
	err := r.db.QueryRowContext(ctx, query, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users by status: %w", err)
	}

	return count, nil
}

// GetActiveMembers retrieves all active members
func (r *PostgresUserRepository) GetActiveMembers(ctx context.Context, params UserListParams) ([]*entity.User, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM users WHERE role = $1 AND status = $2`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, entity.RoleMember, entity.StatusActive).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count active members: %w", err)
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

	// Query users
	query := fmt.Sprintf(`
		SELECT id, email, full_name, role, status, borrowing_limit, created_at, updated_at
		FROM users
		WHERE role = $1 AND status = $2
		ORDER BY %s
		LIMIT $3 OFFSET $4
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, entity.RoleMember, entity.StatusActive, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get active members: %w", err)
	}
	defer rows.Close()

	users := make([]*entity.User, 0)
	for rows.Next() {
		user := &entity.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FullName,
			&user.Role,
			&user.Status,
			&user.BorrowingLimit,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating users: %w", err)
	}

	return users, total, nil
}

// Search searches for users by name or email
func (r *PostgresUserRepository) Search(ctx context.Context, query string, params UserListParams) ([]*entity.User, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM users WHERE full_name ILIKE $1 OR email ILIKE $1`
	searchPattern := "%" + query + "%"
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, searchPattern).Scan(&total)
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

	// Query users
	searchQuery := fmt.Sprintf(`
		SELECT id, email, full_name, role, status, borrowing_limit, created_at, updated_at
		FROM users
		WHERE full_name ILIKE $1 OR email ILIKE $1
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, searchQuery, searchPattern, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	users := make([]*entity.User, 0)
	for rows.Next() {
		user := &entity.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FullName,
			&user.Role,
			&user.Status,
			&user.BorrowingLimit,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating search results: %w", err)
	}

	return users, total, nil
}