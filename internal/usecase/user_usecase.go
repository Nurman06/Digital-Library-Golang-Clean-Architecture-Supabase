package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/repository"
)

// UserUseCase defines the interface for user business logic operations
type UserUseCase interface {
	// RegisterUser registers a new user in the system
	RegisterUser(ctx context.Context, user *entity.User) error

	// GetUserByID retrieves a user by their ID
	GetUserByID(ctx context.Context, id string) (*entity.User, error)

	// GetUserByEmail retrieves a user by their email
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, user *entity.User) error

	// DeleteUser deletes a user from the system
	DeleteUser(ctx context.Context, id string) error

	// ListUsers retrieves a paginated list of users
	ListUsers(ctx context.Context, params repository.UserListParams) ([]*entity.User, int64, error)

	// GetUsersByRole retrieves users by role
	GetUsersByRole(ctx context.Context, role entity.UserRole, params repository.UserListParams) ([]*entity.User, int64, error)

	// GetUsersByStatus retrieves users by status
	GetUsersByStatus(ctx context.Context, status entity.UserStatus, params repository.UserListParams) ([]*entity.User, int64, error)

	// SearchUsers searches for users by name or email
	SearchUsers(ctx context.Context, query string, params repository.UserListParams) ([]*entity.User, int64, error)

	// SuspendUser suspends a user account
	SuspendUser(ctx context.Context, id string) error

	// ActivateUser activates a user account
	ActivateUser(ctx context.Context, id string) error

	// ExpireUser marks a user account as expired
	ExpireUser(ctx context.Context, id string) error

	// UpdateBorrowingLimit updates a user's borrowing limit
	UpdateBorrowingLimit(ctx context.Context, id string, limit int) error

	// GetUserCount returns the total number of users
	GetUserCount(ctx context.Context) (int64, error)

	// GetUserCountByRole returns the number of users by role
	GetUserCountByRole(ctx context.Context, role entity.UserRole) (int64, error)

	// GetActiveMembers retrieves all active members
	GetActiveMembers(ctx context.Context, params repository.UserListParams) ([]*entity.User, int64, error)
}

// userUseCase implements the UserUseCase interface
type userUseCase struct {
	userRepo         repository.UserRepository
	borrowRecordRepo repository.BorrowRecordRepository
}

// NewUserUseCase creates a new instance of UserUseCase
func NewUserUseCase(userRepo repository.UserRepository, borrowRecordRepo repository.BorrowRecordRepository) UserUseCase {
	return &userUseCase{
		userRepo:         userRepo,
		borrowRecordRepo: borrowRecordRepo,
	}
}

// RegisterUser registers a new user in the system
func (uc *userUseCase) RegisterUser(ctx context.Context, user *entity.User) error {
	// Validate the user entity
	if err := user.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if user with same email already exists
	exists, err := uc.userRepo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return errors.New("user with this email already exists")
	}

	// Set default borrowing limit based on role if not provided
	if user.BorrowingLimit == 0 {
		user.BorrowingLimit = entity.GetDefaultBorrowingLimit(user.Role)
	}

	// Set default status to active if not provided
	if user.Status == "" {
		user.Status = entity.StatusActive
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Create the user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user by their ID
func (uc *userUseCase) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	if id == "" {
		return nil, errors.New("user ID is required")
	}

	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email
func (uc *userUseCase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// UpdateUser updates an existing user
func (uc *userUseCase) UpdateUser(ctx context.Context, user *entity.User) error {
	// Validate the user entity
	if err := user.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if user exists
	existingUser, err := uc.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing user: %w", err)
	}

	// If email is being changed, check if new email already exists
	if existingUser.Email != user.Email {
		exists, err := uc.userRepo.ExistsByEmail(ctx, user.Email)
		if err != nil {
			return fmt.Errorf("failed to check email existence: %w", err)
		}
		if exists {
			return errors.New("user with this email already exists")
		}
	}

	// Update timestamp
	user.UpdatedAt = time.Now()

	// Update the user
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser deletes a user from the system
func (uc *userUseCase) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("user ID is required")
	}

	// Check if user exists
	_, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user has any active borrows
	activeBorrows, err := uc.borrowRecordRepo.GetActiveByUserID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get active borrows: %w", err)
	}
	if len(activeBorrows) > 0 {
		return errors.New("cannot delete user with active borrows")
	}

	// Check if user has any overdue borrows
	overdueBorrows, err := uc.borrowRecordRepo.GetOverdueByUserID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get overdue borrows: %w", err)
	}
	if len(overdueBorrows) > 0 {
		return errors.New("cannot delete user with overdue borrows")
	}

	// Check if user has unpaid late fees
	totalLateFees, err := uc.borrowRecordRepo.CalculateTotalLateFees(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to calculate late fees: %w", err)
	}
	if totalLateFees > 0 {
		return errors.New("cannot delete user with unpaid late fees")
	}

	// Delete the user
	if err := uc.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers retrieves a paginated list of users
func (uc *userUseCase) ListUsers(ctx context.Context, params repository.UserListParams) ([]*entity.User, int64, error) {
	// Set default pagination values if not provided
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100 // Max page size
	}

	users, total, err := uc.userRepo.List(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// GetUsersByRole retrieves users by role
func (uc *userUseCase) GetUsersByRole(ctx context.Context, role entity.UserRole, params repository.UserListParams) ([]*entity.User, int64, error) {
	// Set default pagination values if not provided
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100 // Max page size
	}

	users, total, err := uc.userRepo.GetByRole(ctx, role, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users by role: %w", err)
	}

	return users, total, nil
}

// GetUsersByStatus retrieves users by status
func (uc *userUseCase) GetUsersByStatus(ctx context.Context, status entity.UserStatus, params repository.UserListParams) ([]*entity.User, int64, error) {
	// Set default pagination values if not provided
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100 // Max page size
	}

	users, total, err := uc.userRepo.GetByStatus(ctx, status, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users by status: %w", err)
	}

	return users, total, nil
}

// SearchUsers searches for users by name or email
func (uc *userUseCase) SearchUsers(ctx context.Context, query string, params repository.UserListParams) ([]*entity.User, int64, error) {
	if query == "" {
		return nil, 0, errors.New("search query is required")
	}

	// Set default pagination values if not provided
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100 // Max page size
	}

	users, total, err := uc.userRepo.Search(ctx, query, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}

	return users, total, nil
}

// SuspendUser suspends a user account
func (uc *userUseCase) SuspendUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("user ID is required")
	}

	// Get the user
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Suspend the user
	if err := user.Suspend(); err != nil {
		return fmt.Errorf("failed to suspend user: %w", err)
	}

	// Update the user status in repository
	if err := uc.userRepo.UpdateStatus(ctx, id, entity.StatusSuspended); err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	return nil
}

// ActivateUser activates a user account
func (uc *userUseCase) ActivateUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("user ID is required")
	}

	// Get the user
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Activate the user
	if err := user.Activate(); err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}

	// Update the user status in repository
	if err := uc.userRepo.UpdateStatus(ctx, id, entity.StatusActive); err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	return nil
}

// ExpireUser marks a user account as expired
func (uc *userUseCase) ExpireUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("user ID is required")
	}

	// Get the user
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Expire the user
	if err := user.Expire(); err != nil {
		return fmt.Errorf("failed to expire user: %w", err)
	}

	// Update the user status in repository
	if err := uc.userRepo.UpdateStatus(ctx, id, entity.StatusExpired); err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	return nil
}

// UpdateBorrowingLimit updates a user's borrowing limit
func (uc *userUseCase) UpdateBorrowingLimit(ctx context.Context, id string, limit int) error {
	if id == "" {
		return errors.New("user ID is required")
	}

	if limit < 0 {
		return errors.New("borrowing limit cannot be negative")
	}

	// Get the user
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Update the borrowing limit
	if err := user.UpdateBorrowingLimit(limit); err != nil {
		return fmt.Errorf("failed to update borrowing limit: %w", err)
	}

	// Update the borrowing limit in repository
	if err := uc.userRepo.UpdateBorrowingLimit(ctx, id, limit); err != nil {
		return fmt.Errorf("failed to update borrowing limit in repository: %w", err)
	}

	return nil
}

// GetUserCount returns the total number of users
func (uc *userUseCase) GetUserCount(ctx context.Context) (int64, error) {
	count, err := uc.userRepo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// GetUserCountByRole returns the number of users by role
func (uc *userUseCase) GetUserCountByRole(ctx context.Context, role entity.UserRole) (int64, error) {
	count, err := uc.userRepo.CountByRole(ctx, role)
	if err != nil {
		return 0, fmt.Errorf("failed to count users by role: %w", err)
	}

	return count, nil
}

// GetActiveMembers retrieves all active members
func (uc *userUseCase) GetActiveMembers(ctx context.Context, params repository.UserListParams) ([]*entity.User, int64, error) {
	// Set default pagination values if not provided
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100 // Max page size
	}

	users, total, err := uc.userRepo.GetActiveMembers(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get active members: %w", err)
	}

	return users, total, nil
}