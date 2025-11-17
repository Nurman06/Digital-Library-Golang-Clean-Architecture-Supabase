package repository

import (
	"context"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
)

// UserRepository defines the interface for user data access operations
type UserRepository interface {
	// Create creates a new user in the repository
	Create(ctx context.Context, user *entity.User) error

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id string) (*entity.User, error)

	// GetByEmail retrieves a user by their email address
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entity.User) error

	// Delete deletes a user from the repository
	Delete(ctx context.Context, id string) error

	// List retrieves a paginated list of users
	List(ctx context.Context, params UserListParams) ([]*entity.User, int64, error)

	// GetByRole retrieves users by role with pagination
	GetByRole(ctx context.Context, role entity.UserRole, params UserListParams) ([]*entity.User, int64, error)

	// GetByStatus retrieves users by status with pagination
	GetByStatus(ctx context.Context, status entity.UserStatus, params UserListParams) ([]*entity.User, int64, error)

	// ExistsByEmail checks if a user with the given email exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// UpdateStatus updates the status of a user
	UpdateStatus(ctx context.Context, id string, status entity.UserStatus) error

	// UpdateBorrowingLimit updates the borrowing limit of a user
	UpdateBorrowingLimit(ctx context.Context, id string, limit int) error

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)

	// CountByRole returns the number of users by role
	CountByRole(ctx context.Context, role entity.UserRole) (int64, error)

	// CountByStatus returns the number of users by status
	CountByStatus(ctx context.Context, status entity.UserStatus) (int64, error)

	// GetActiveMembers retrieves all active members
	GetActiveMembers(ctx context.Context, params UserListParams) ([]*entity.User, int64, error)

	// Search searches for users by name or email
	Search(ctx context.Context, query string, params UserListParams) ([]*entity.User, int64, error)
}

// UserListParams defines parameters for listing users
type UserListParams struct {
	Page      int
	PageSize  int
	SortBy    string // e.g., "full_name", "email", "created_at"
	SortOrder string // "asc" or "desc"
}
