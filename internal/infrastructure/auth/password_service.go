package auth

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// PasswordService handles password hashing and validation
type PasswordService struct {
	cost int // bcrypt cost factor
}

// NewPasswordService creates a new password service
func NewPasswordService() *PasswordService {
	return &PasswordService{
		cost: bcrypt.DefaultCost, // Cost of 10
	}
}

// HashPassword hashes a plain text password using bcrypt
func (s *PasswordService) HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	// Validate password requirements
	if err := s.ValidatePasswordStrength(password); err != nil {
		return "", err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// ComparePassword compares a plain text password with a hashed password
func (s *PasswordService) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ValidatePasswordStrength validates password strength requirements
func (s *PasswordService) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Check for at least one letter
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	if !hasLetter {
		return errors.New("password must contain at least one letter")
	}

	// Check for at least one number
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}

	return nil
}

// IsPasswordValid checks if a password meets the strength requirements
func (s *PasswordService) IsPasswordValid(password string) bool {
	return s.ValidatePasswordStrength(password) == nil
}