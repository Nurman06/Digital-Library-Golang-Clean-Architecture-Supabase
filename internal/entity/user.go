package entity

import (
	"errors"
	"regexp"
	"time"
)

// User represents a library user
type User struct {
	ID             string     `json:"id"`
	Email          string     `json:"email"`
	FullName       string     `json:"full_name"`
	Role           UserRole   `json:"role"`
	Status         UserStatus `json:"status"`
	BorrowingLimit int        `json:"borrowing_limit"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// UserRole represents the role of a user
type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleLibrarian UserRole = "librarian"
	RoleMember    UserRole = "member"
)

// UserStatus represents the status of a user account
type UserStatus string

const (
	StatusActive    UserStatus = "active"
	StatusSuspended UserStatus = "suspended"
	StatusExpired   UserStatus = "expired"
)

var (
	emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// Validate validates the user entity
func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(u.Email) {
		return errors.New("invalid email format")
	}
	if u.FullName == "" {
		return errors.New("full name is required")
	}
	if u.Role == "" {
		return errors.New("role is required")
	}
	if !isValidRole(u.Role) {
		return errors.New("invalid user role")
	}
	if u.Status == "" {
		return errors.New("status is required")
	}
	if !isValidStatus(u.Status) {
		return errors.New("invalid user status")
	}
	if u.BorrowingLimit < 0 {
		return errors.New("borrowing limit cannot be negative")
	}
	return nil
}

// isValidEmail checks if the email is in valid format
func isValidEmail(email string) bool {
	return emailPattern.MatchString(email)
}

// isValidRole checks if the user role is valid
func isValidRole(role UserRole) bool {
	switch role {
	case RoleAdmin, RoleLibrarian, RoleMember:
		return true
	}
	return false
}

// isValidStatus checks if the user status is valid
func isValidStatus(status UserStatus) bool {
	switch status {
	case StatusActive, StatusSuspended, StatusExpired:
		return true
	}
	return false
}

// IsActive checks if the user account is active
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// CanBorrow checks if the user can borrow books
func (u *User) CanBorrow() bool {
	return u.Status == StatusActive
}

// GetDefaultBorrowingLimit returns the default borrowing limit based on role
func GetDefaultBorrowingLimit(role UserRole) int {
	switch role {
	case RoleAdmin:
		return 0 // unlimited
	case RoleLibrarian:
		return 10
	case RoleMember:
		return 5
	default:
		return 5
	}
}

// HasPermission checks if the user has permission for an operation
func (u *User) HasPermission(operation string) bool {
	switch u.Role {
	case RoleAdmin:
		return true // Admin has all permissions
	case RoleLibrarian:
		// Librarian can manage books and borrowing but not users
		return operation != "manage_users"
	case RoleMember:
		// Member can only view and borrow
		return operation == "view_books" || operation == "borrow_books" || operation == "view_own_profile"
	}
	return false
}