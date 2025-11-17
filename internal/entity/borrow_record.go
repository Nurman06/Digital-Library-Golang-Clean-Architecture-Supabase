package entity

import (
	"errors"
	"time"
)

// BorrowRecord represents a borrowing transaction
type BorrowRecord struct {
	ID           string              `json:"id"`
	UserID       string              `json:"user_id"`
	BookCopyID   string              `json:"book_copy_id"`
	CheckoutDate time.Time           `json:"checkout_date"`
	DueDate      time.Time           `json:"due_date"`
	ReturnDate   *time.Time          `json:"return_date,omitempty"`
	RenewalCount int                 `json:"renewal_count"`
	LateFee      float64             `json:"late_fee"`
	Status       BorrowRecordStatus  `json:"status"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

// BorrowRecordStatus represents the status of a borrow record
type BorrowRecordStatus string

const (
	BorrowStatusActive   BorrowRecordStatus = "active"
	BorrowStatusReturned BorrowRecordStatus = "returned"
	BorrowStatusOverdue  BorrowRecordStatus = "overdue"
)

const (
	// DefaultBorrowingPeriod is 14 days
	DefaultBorrowingPeriod = 14 * 24 * time.Hour
	// MaxRenewalCount is the maximum number of renewals allowed
	MaxRenewalCount = 2
	// LateFeePerDay is the late fee charged per day
	LateFeePerDay = 0.50
	// MaxLateFee is the maximum late fee that can be charged
	MaxLateFee = 10.00
)

// Validate validates the borrow record entity
func (br *BorrowRecord) Validate() error {
	if br.UserID == "" {
		return errors.New("user ID is required")
	}
	if br.BookCopyID == "" {
		return errors.New("book copy ID is required")
	}
	if br.CheckoutDate.IsZero() {
		return errors.New("checkout date is required")
	}
	if br.DueDate.IsZero() {
		return errors.New("due date is required")
	}
	if br.DueDate.Before(br.CheckoutDate) {
		return errors.New("due date must be after checkout date")
	}
	if br.Status == "" {
		return errors.New("status is required")
	}
	if !isValidBorrowStatus(br.Status) {
		return errors.New("invalid borrow status")
	}
	if br.RenewalCount < 0 {
		return errors.New("renewal count cannot be negative")
	}
	if br.LateFee < 0 {
		return errors.New("late fee cannot be negative")
	}
	return nil
}

// isValidBorrowStatus checks if the borrow status is valid
func isValidBorrowStatus(status BorrowRecordStatus) bool {
	switch status {
	case BorrowStatusActive, BorrowStatusReturned, BorrowStatusOverdue:
		return true
	}
	return false
}

// IsOverdue checks if the borrow record is overdue
func (br *BorrowRecord) IsOverdue() bool {
	return time.Now().After(br.DueDate) && br.Status == BorrowStatusActive
}

// CanRenew checks if the borrow record can be renewed
func (br *BorrowRecord) CanRenew() bool {
	return br.Status == BorrowStatusActive && 
		   br.RenewalCount < MaxRenewalCount && 
		   !br.IsOverdue()
}

// CalculateLateFee calculates the late fee based on overdue days
func (br *BorrowRecord) CalculateLateFee() float64 {
	if br.ReturnDate == nil || br.ReturnDate.Before(br.DueDate) || br.ReturnDate.Equal(br.DueDate) {
		return 0
	}
	
	overdueDays := int(br.ReturnDate.Sub(br.DueDate).Hours() / 24)
	if overdueDays <= 0 {
		return 0
	}
	
	lateFee := float64(overdueDays) * LateFeePerDay
	if lateFee > MaxLateFee {
		return MaxLateFee
	}
	return lateFee
}

// Reservation represents a book reservation
type Reservation struct {
	ID              string            `json:"id"`
	UserID          string            `json:"user_id"`
	BookID          string            `json:"book_id"`
	ReservationDate time.Time         `json:"reservation_date"`
	ExpiryDate      time.Time         `json:"expiry_date"`
	Status          ReservationStatus `json:"status"`
	QueuePosition   int               `json:"queue_position"`
}

// ReservationStatus represents the status of a reservation
type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "pending"
	ReservationStatusFulfilled ReservationStatus = "fulfilled"
	ReservationStatusExpired   ReservationStatus = "expired"
	ReservationStatusCancelled ReservationStatus = "cancelled"
)