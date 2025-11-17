package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/repository"
)

// AuthUseCase defines the interface for authentication business logic operations
type AuthUseCase interface {
	// Login authenticates a user and returns user information
	Login(ctx context.Context, email, password string) (*entity.User, error)

	// ValidateUser validates if a user can perform actions
	ValidateUser(ctx context.Context, userID string) error

	// CheckPermission checks if a user has permission for an operation
	CheckPermission(ctx context.Context, userID string, operation string) (bool, error)

	// ValidateUserStatus validates if a user's status allows them to perform actions
	ValidateUserStatus(ctx context.Context, userID string) error

	// GetUserRole retrieves a user's role
	GetUserRole(ctx context.Context, userID string) (entity.UserRole, error)

	// CanUserBorrow checks if a user can borrow books
	CanUserBorrow(ctx context.Context, userID string) (bool, error)

	// CanUserManageBooks checks if a user can manage books
	CanUserManageBooks(ctx context.Context, userID string) (bool, error)

	// CanUserManageUsers checks if a user can manage other users
	CanUserManageUsers(ctx context.Context, userID string) (bool, error)
}

// authUseCase implements the AuthUseCase interface
type authUseCase struct {
	userRepo repository.UserRepository
}

// NewAuthUseCase creates a new instance of AuthUseCase
func NewAuthUseCase(userRepo repository.UserRepository) AuthUseCase {
	return &authUseCase{
		userRepo: userRepo,
	}
}

// Login authenticates a user and returns user information
// Note: Actual password verification is handled by Supabase Auth
// This method is for retrieving user information after Supabase authentication
func (uc *authUseCase) Login(ctx context.Context, email, password string) (*entity.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Validate user status
	if err := uc.ValidateUserStatus(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return user, nil
}

// ValidateUser validates if a user can perform actions
func (uc *authUseCase) ValidateUser(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user validation failed: %w", err)
	}

	// Check if user is active
	if !user.IsActive() {
		return errors.New("user account is not active")
	}

	return nil
}

// CheckPermission checks if a user has permission for an operation
func (uc *authUseCase) CheckPermission(ctx context.Context, userID string, operation string) (bool, error) {
	if userID == "" {
		return false, errors.New("user ID is required")
	}
	if operation == "" {
		return false, errors.New("operation is required")
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is active
	if !user.IsActive() {
		return false, errors.New("user account is not active")
	}

	// Check permission
	hasPermission := user.HasPermission(operation)
	return hasPermission, nil
}

// ValidateUserStatus validates if a user's status allows them to perform actions
func (uc *authUseCase) ValidateUserStatus(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check user status
	switch {
	case user.IsSuspended():
		return errors.New("user account is suspended")
	case user.IsExpired():
		return errors.New("user account has expired")
	case !user.IsActive():
		return errors.New("user account is not active")
	}

	return nil
}

// GetUserRole retrieves a user's role
func (uc *authUseCase) GetUserRole(ctx context.Context, userID string) (entity.UserRole, error) {
	if userID == "" {
		return "", errors.New("user ID is required")
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	return user.Role, nil
}

// CanUserBorrow checks if a user can borrow books
func (uc *authUseCase) CanUserBorrow(ctx context.Context, userID string) (bool, error) {
	if userID == "" {
		return false, errors.New("user ID is required")
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user can borrow
	canBorrow := user.CanBorrow()
	return canBorrow, nil
}

// CanUserManageBooks checks if a user can manage books
func (uc *authUseCase) CanUserManageBooks(ctx context.Context, userID string) (bool, error) {
	if userID == "" {
		return false, errors.New("user ID is required")
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is active
	if !user.IsActive() {
		return false, nil
	}

	// Check if user can manage books
	canManage := user.CanManageBooks()
	return canManage, nil
}

// CanUserManageUsers checks if a user can manage other users
func (uc *authUseCase) CanUserManageUsers(ctx context.Context, userID string) (bool, error) {
	if userID == "" {
		return false, errors.New("user ID is required")
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is active
	if !user.IsActive() {
		return false, nil
	}

	// Check if user can manage users
	canManage := user.CanManageUsers()
	return canManage, nil
}

// TokenService defines the interface for JWT token operations
// This will be used in conjunction with Supabase Auth
type TokenService interface {
	// GenerateToken generates a JWT token for a user
	GenerateToken(user *entity.User, expiresIn time.Duration) (string, error)

	// ValidateToken validates a JWT token and returns the user ID
	ValidateToken(token string) (string, error)

	// RefreshToken refreshes an expired token
	RefreshToken(refreshToken string) (string, error)

	// RevokeToken revokes a token
	RevokeToken(token string) error
}