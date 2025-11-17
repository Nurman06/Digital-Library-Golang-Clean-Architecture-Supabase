package entity

import (
	"errors"
	"time"
)

// BorrowRecord represents a borrowing transaction
type BorrowRecord struct {
	ID           string             `json:"id"`
	UserID       string             `json:"user_id"`
	BookCopyID   string             `json:"book_copy_id"`
	CheckoutDate time.Time          `json:"checkout_date"`
	DueDate      time.Time          `json:"due_date"`
	ReturnDate   *time.Time         `json:"return_date,omitempty"`
	RenewalCount int                `json:"renewal_count"`
	LateFee      float64            `json:"late_fee"`
	Status       BorrowRecordStatus `json:"status"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
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

// CalculateCurrentLateFee calculates the current late fee if book is overdue
func (br *BorrowRecord) CalculateCurrentLateFee() float64 {
	if br.Status != BorrowStatusActive || !br.IsOverdue() {
		return 0
	}

	overdueDays := int(time.Since(br.DueDate).Hours() / 24)
	if overdueDays <= 0 {
		return 0
	}

	lateFee := float64(overdueDays) * LateFeePerDay
	if lateFee > MaxLateFee {
		return MaxLateFee
	}
	return lateFee
}

// GetOverdueDays returns the number of days overdue (0 if not overdue)
func (br *BorrowRecord) GetOverdueDays() int {
	if !br.IsOverdue() {
		return 0
	}

	days := int(time.Since(br.DueDate).Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}

// IsActive checks if the borrow record is currently active
func (br *BorrowRecord) IsActive() bool {
	return br.Status == BorrowStatusActive
}

// IsReturned checks if the borrow record has been returned
func (br *BorrowRecord) IsReturned() bool {
	return br.Status == BorrowStatusReturned
}

// Renew extends the due date for the borrow record
func (br *BorrowRecord) Renew() error {
	if !br.CanRenew() {
		return errors.New("book cannot be renewed")
	}

	br.DueDate = br.DueDate.Add(DefaultBorrowingPeriod)
	br.RenewalCount++
	br.UpdatedAt = time.Now()

	return nil
}

// Return processes the return of a borrowed book
func (br *BorrowRecord) Return() error {
	if br.Status == BorrowStatusReturned {
		return errors.New("book has already been returned")
	}

	now := time.Now()
	br.ReturnDate = &now
	br.Status = BorrowStatusReturned
	br.LateFee = br.CalculateLateFee()
	br.UpdatedAt = now

	return nil
}

// UpdateStatus updates the status of the borrow record based on current time
func (br *BorrowRecord) UpdateStatus() {
	if br.Status == BorrowStatusActive && br.IsOverdue() {
		br.Status = BorrowStatusOverdue
		br.UpdatedAt = time.Now()
	}
}

// GetDaysRemaining returns the number of days remaining until due date (negative if overdue)
func (br *BorrowRecord) GetDaysRemaining() int {
	if br.Status != BorrowStatusActive {
		return 0
	}

	duration := time.Until(br.DueDate)
	return int(duration.Hours() / 24)
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

const (
	// DefaultReservationHoldPeriod is 3 days
	DefaultReservationHoldPeriod = 3 * 24 * time.Hour
)

// Validate validates the reservation entity
func (r *Reservation) Validate() error {
	if r.UserID == "" {
		return errors.New("user ID is required")
	}
	if r.BookID == "" {
		return errors.New("book ID is required")
	}
	if r.ReservationDate.IsZero() {
		return errors.New("reservation date is required")
	}
	if r.ExpiryDate.IsZero() {
		return errors.New("expiry date is required")
	}
	if r.ExpiryDate.Before(r.ReservationDate) {
		return errors.New("expiry date must be after reservation date")
	}
	if r.Status == "" {
		return errors.New("status is required")
	}
	if !isValidReservationStatus(r.Status) {
		return errors.New("invalid reservation status")
	}
	if r.QueuePosition < 1 {
		return errors.New("queue position must be at least 1")
	}
	return nil
}

// isValidReservationStatus checks if the reservation status is valid
func isValidReservationStatus(status ReservationStatus) bool {
	switch status {
	case ReservationStatusPending, ReservationStatusFulfilled, ReservationStatusExpired, ReservationStatusCancelled:
		return true
	}
	return false
}

// IsExpired checks if the reservation has expired
func (r *Reservation) IsExpired() bool {
	return time.Now().After(r.ExpiryDate) && r.Status == ReservationStatusPending
}

// IsPending checks if the reservation is still pending
func (r *Reservation) IsPending() bool {
	return r.Status == ReservationStatusPending && !r.IsExpired()
}

// CanBeFulfilled checks if the reservation can be fulfilled
func (r *Reservation) CanBeFulfilled() bool {
	return r.Status == ReservationStatusPending && !r.IsExpired()
}

// Cancel cancels the reservation
func (r *Reservation) Cancel() error {
	if r.Status != ReservationStatusPending {
		return errors.New("only pending reservations can be cancelled")
	}
	r.Status = ReservationStatusCancelled
	return nil
}

// Fulfill marks the reservation as fulfilled
func (r *Reservation) Fulfill() error {
	if !r.CanBeFulfilled() {
		return errors.New("reservation cannot be fulfilled")
	}
	r.Status = ReservationStatusFulfilled
	return nil
}

// Expire marks the reservation as expired
func (r *Reservation) Expire() error {
	if r.Status != ReservationStatusPending {
		return errors.New("only pending reservations can be expired")
	}
	r.Status = ReservationStatusExpired
	return nil
}
