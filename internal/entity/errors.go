package entity

import (
	"fmt"
)

// DomainError represents a domain-level error with structured information
type DomainError struct {
	Code    ErrorCode
	Message string
	Details map[string]any
	Err     error
}

// ErrorCode represents standardized error codes for domain errors
type ErrorCode string

const (
	// Validation errors
	ErrCodeValidation    ErrorCode = "VALIDATION_ERROR"
	ErrCodeInvalidISBN   ErrorCode = "INVALID_ISBN"
	ErrCodeInvalidEmail  ErrorCode = "INVALID_EMAIL"
	ErrCodeInvalidRole   ErrorCode = "INVALID_ROLE"
	ErrCodeInvalidStatus ErrorCode = "INVALID_STATUS"
	ErrCodeInvalidDate   ErrorCode = "INVALID_DATE"

	// Business logic errors
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeConflict     ErrorCode = "CONFLICT"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"

	// Book-specific errors
	ErrCodeBookNotAvailable     ErrorCode = "BOOK_NOT_AVAILABLE"
	ErrCodeDuplicateISBN        ErrorCode = "DUPLICATE_ISBN"
	ErrCodeBookHasActiveBorrows ErrorCode = "BOOK_HAS_ACTIVE_BORROWS"

	// Borrowing-specific errors
	ErrCodeBorrowingLimitReached ErrorCode = "BORROWING_LIMIT_REACHED"
	ErrCodeOverdueBooks          ErrorCode = "OVERDUE_BOOKS"
	ErrCodeUnpaidLateFees        ErrorCode = "UNPAID_LATE_FEES"
	ErrCodeCannotRenew           ErrorCode = "CANNOT_RENEW"
	ErrCodeAlreadyReturned       ErrorCode = "ALREADY_RETURNED"
	ErrCodeNotBorrowed           ErrorCode = "NOT_BORROWED"

	// User-specific errors
	ErrCodeUserSuspended        ErrorCode = "USER_SUSPENDED"
	ErrCodeUserExpired          ErrorCode = "USER_EXPIRED"
	ErrCodeDuplicateEmail       ErrorCode = "DUPLICATE_EMAIL"
	ErrCodeUserHasActiveBorrows ErrorCode = "USER_HAS_ACTIVE_BORROWS"

	// Reservation-specific errors
	ErrCodeReservationExpired  ErrorCode = "RESERVATION_EXPIRED"
	ErrCodeReservationNotFound ErrorCode = "RESERVATION_NOT_FOUND"

	// Internal errors
	ErrCodeInternal ErrorCode = "INTERNAL_ERROR"
)

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError creates a new domain error
func NewDomainError(code ErrorCode, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: make(map[string]any),
	}
}

// NewDomainErrorWithDetails creates a new domain error with details
func NewDomainErrorWithDetails(code ErrorCode, message string, details map[string]any) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// WrapError wraps an existing error as a domain error
func WrapError(code ErrorCode, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: make(map[string]any),
		Err:     err,
	}
}

// AddDetail adds a detail to the error
func (e *DomainError) AddDetail(key string, value any) *DomainError {
	if e.Details == nil {
		e.Details = make(map[string]any)
	}
	e.Details[key] = value
	return e
}

// Predefined error constructors for common scenarios

// NewValidationError creates a validation error
func NewValidationError(message string) *DomainError {
	return NewDomainError(ErrCodeValidation, message)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) *DomainError {
	return NewDomainError(ErrCodeNotFound, fmt.Sprintf("%s not found", resource))
}

// NewConflictError creates a conflict error
func NewConflictError(message string) *DomainError {
	return NewDomainError(ErrCodeConflict, message)
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *DomainError {
	return NewDomainError(ErrCodeUnauthorized, message)
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string) *DomainError {
	return NewDomainError(ErrCodeForbidden, message)
}

// NewBookNotAvailableError creates a book not available error
func NewBookNotAvailableError(bookID string) *DomainError {
	return NewDomainErrorWithDetails(
		ErrCodeBookNotAvailable,
		"Book is not available for borrowing",
		map[string]any{"book_id": bookID},
	)
}

// NewBorrowingLimitReachedError creates a borrowing limit reached error
func NewBorrowingLimitReachedError(userID string, limit int) *DomainError {
	return NewDomainErrorWithDetails(
		ErrCodeBorrowingLimitReached,
		fmt.Sprintf("Borrowing limit of %d books reached", limit),
		map[string]any{
			"user_id": userID,
			"limit":   limit,
		},
	)
}

// NewOverdueBooksError creates an overdue books error
func NewOverdueBooksError(userID string, overdueCount int) *DomainError {
	return NewDomainErrorWithDetails(
		ErrCodeOverdueBooks,
		fmt.Sprintf("User has %d overdue books", overdueCount),
		map[string]any{
			"user_id":       userID,
			"overdue_count": overdueCount,
		},
	)
}

// NewCannotRenewError creates a cannot renew error
func NewCannotRenewError(reason string) *DomainError {
	return NewDomainErrorWithDetails(
		ErrCodeCannotRenew,
		"Book cannot be renewed",
		map[string]any{"reason": reason},
	)
}

// NewUserSuspendedError creates a user suspended error
func NewUserSuspendedError(userID string) *DomainError {
	return NewDomainErrorWithDetails(
		ErrCodeUserSuspended,
		"User account is suspended",
		map[string]any{"user_id": userID},
	)
}

// NewDuplicateISBNError creates a duplicate ISBN error
func NewDuplicateISBNError(isbn string) *DomainError {
	return NewDomainErrorWithDetails(
		ErrCodeDuplicateISBN,
		"A book with this ISBN already exists",
		map[string]any{"isbn": isbn},
	)
}

// NewDuplicateEmailError creates a duplicate email error
func NewDuplicateEmailError(email string) *DomainError {
	return NewDomainErrorWithDetails(
		ErrCodeDuplicateEmail,
		"A user with this email already exists",
		map[string]any{"email": email},
	)
}
