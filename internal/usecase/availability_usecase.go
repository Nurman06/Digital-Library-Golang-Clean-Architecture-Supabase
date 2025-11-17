package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/repository"
)

// AvailabilityUseCase defines the interface for availability tracking business logic
type AvailabilityUseCase interface {
	// GetBookAvailability retrieves availability information for a specific book
	GetBookAvailability(ctx context.Context, bookID string) (*BookAvailability, error)

	// GetBookCopyAvailability retrieves availability information for a specific book copy
	GetBookCopyAvailability(ctx context.Context, bookCopyID string) (*BookCopyAvailability, error)

	// CheckBookAvailability checks if a book has any available copies
	CheckBookAvailability(ctx context.Context, bookID string) (bool, error)

	// GetAvailableBooks retrieves all books that have at least one available copy
	GetAvailableBooks(ctx context.Context, params repository.ListParams) ([]*BookWithAvailability, int64, error)

	// GetBookCopiesByStatus retrieves book copies by their status
	GetBookCopiesByStatus(ctx context.Context, status entity.CopyStatus, params repository.ListParams) ([]*entity.BookCopy, int64, error)

	// GetExpectedReturnDate retrieves the expected return date for a borrowed book copy
	GetExpectedReturnDate(ctx context.Context, bookCopyID string) (*time.Time, error)

	// GetBooksWithLowAvailability retrieves books with low availability (below threshold)
	GetBooksWithLowAvailability(ctx context.Context, threshold int) ([]*BookAvailability, error)

	// GetAvailabilityStatistics retrieves overall availability statistics
	GetAvailabilityStatistics(ctx context.Context) (*AvailabilityStatistics, error)

	// ReserveBook creates a reservation for a book
	ReserveBook(ctx context.Context, userID, bookID string) (*entity.Reservation, error)

	// CancelReservation cancels a book reservation
	CancelReservation(ctx context.Context, reservationID string) error

	// GetUserReservations retrieves all reservations for a user
	GetUserReservations(ctx context.Context, userID string) ([]*entity.Reservation, error)

	// GetBookReservationQueue retrieves the reservation queue for a book
	GetBookReservationQueue(ctx context.Context, bookID string) ([]*entity.Reservation, error)

	// ProcessExpiredReservations processes and expires old reservations
	ProcessExpiredReservations(ctx context.Context) error
}

// BookAvailability represents the availability information for a book
type BookAvailability struct {
	BookID           string
	BookTitle        string
	TotalCopies      int64
	AvailableCopies  int64
	BorrowedCopies   int64
	ReservedCopies   int64
	DamagedCopies    int64
	LostCopies       int64
	NextAvailableDate *time.Time
	ReservationCount int64
}

// BookCopyAvailability represents the availability information for a book copy
type BookCopyAvailability struct {
	BookCopy         *entity.BookCopy
	IsAvailable      bool
	CurrentBorrower  *entity.User
	ExpectedReturn   *time.Time
	ReservationCount int64
}

// BookWithAvailability represents a book with its availability information
type BookWithAvailability struct {
	Book             *entity.Book
	AvailableCopies  int64
	TotalCopies      int64
}

// AvailabilityStatistics represents overall availability statistics
type AvailabilityStatistics struct {
	TotalBooks          int64
	TotalCopies         int64
	AvailableCopies     int64
	BorrowedCopies      int64
	ReservedCopies      int64
	DamagedCopies       int64
	LostCopies          int64
	AvailabilityRate    float64
	TotalReservations   int64
	PendingReservations int64
}

// availabilityUseCase implements the AvailabilityUseCase interface
type availabilityUseCase struct {
	bookRepo         repository.BookRepository
	bookCopyRepo     repository.BookCopyRepository
	borrowRecordRepo repository.BorrowRecordRepository
	reservationRepo  repository.ReservationRepository
	userRepo         repository.UserRepository
}

// NewAvailabilityUseCase creates a new instance of AvailabilityUseCase
func NewAvailabilityUseCase(
	bookRepo repository.BookRepository,
	bookCopyRepo repository.BookCopyRepository,
	borrowRecordRepo repository.BorrowRecordRepository,
	reservationRepo repository.ReservationRepository,
	userRepo repository.UserRepository,
) AvailabilityUseCase {
	return &availabilityUseCase{
		bookRepo:         bookRepo,
		bookCopyRepo:     bookCopyRepo,
		borrowRecordRepo: borrowRecordRepo,
		reservationRepo:  reservationRepo,
		userRepo:         userRepo,
	}
}

// GetBookAvailability retrieves availability information for a specific book
func (uc *availabilityUseCase) GetBookAvailability(ctx context.Context, bookID string) (*BookAvailability, error) {
	if bookID == "" {
		return nil, errors.New("book ID is required")
	}

	// Get book
	book, err := uc.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	// Get all copies
	copies, err := uc.bookCopyRepo.GetByBookID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book copies: %w", err)
	}

	// Count copies by status
	var availableCount, borrowedCount, reservedCount, damagedCount, lostCount int64
	for _, copy := range copies {
		switch copy.Status {
		case entity.CopyStatusAvailable:
			availableCount++
		case entity.CopyStatusBorrowed:
			borrowedCount++
		case entity.CopyStatusReserved:
			reservedCount++
		case entity.CopyStatusDamaged:
			damagedCount++
		case entity.CopyStatusLost:
			lostCount++
		}
	}

	// Get reservation count
	reservationCount, err := uc.reservationRepo.CountPendingByBookID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to count reservations: %w", err)
	}

	// Find next available date (earliest return date for borrowed copies)
	var nextAvailableDate *time.Time
	if availableCount == 0 && borrowedCount > 0 {
		// Get all active borrow records for this book's copies
		for _, copy := range copies {
			if copy.IsBorrowed() {
				borrowRecord, err := uc.borrowRecordRepo.GetActiveByBookCopyID(ctx, copy.ID)
				if err != nil {
					continue
				}
				if nextAvailableDate == nil || borrowRecord.DueDate.Before(*nextAvailableDate) {
					nextAvailableDate = &borrowRecord.DueDate
				}
			}
		}
	}

	return &BookAvailability{
		BookID:            bookID,
		BookTitle:         book.Title,
		TotalCopies:       int64(len(copies)),
		AvailableCopies:   availableCount,
		BorrowedCopies:    borrowedCount,
		ReservedCopies:    reservedCount,
		DamagedCopies:     damagedCount,
		LostCopies:        lostCount,
		NextAvailableDate: nextAvailableDate,
		ReservationCount:  reservationCount,
	}, nil
}

// GetBookCopyAvailability retrieves availability information for a specific book copy
func (uc *availabilityUseCase) GetBookCopyAvailability(ctx context.Context, bookCopyID string) (*BookCopyAvailability, error) {
	if bookCopyID == "" {
		return nil, errors.New("book copy ID is required")
	}

	// Get book copy
	bookCopy, err := uc.bookCopyRepo.GetByID(ctx, bookCopyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book copy: %w", err)
	}

	availability := &BookCopyAvailability{
		BookCopy:    bookCopy,
		IsAvailable: bookCopy.IsAvailable(),
	}

	// If borrowed, get borrower and expected return date
	if bookCopy.IsBorrowed() {
		borrowRecord, err := uc.borrowRecordRepo.GetActiveByBookCopyID(ctx, bookCopyID)
		if err == nil {
			availability.ExpectedReturn = &borrowRecord.DueDate

			// Get borrower information
			user, err := uc.userRepo.GetByID(ctx, borrowRecord.UserID)
			if err == nil {
				availability.CurrentBorrower = user
			}
		}
	}

	// Get reservation count for the book
	reservationCount, err := uc.reservationRepo.CountPendingByBookID(ctx, bookCopy.BookID)
	if err == nil {
		availability.ReservationCount = reservationCount
	}

	return availability, nil
}

// CheckBookAvailability checks if a book has any available copies
func (uc *availabilityUseCase) CheckBookAvailability(ctx context.Context, bookID string) (bool, error) {
	if bookID == "" {
		return false, errors.New("book ID is required")
	}

	count, err := uc.bookCopyRepo.CountAvailableByBookID(ctx, bookID)
	if err != nil {
		return false, fmt.Errorf("failed to count available copies: %w", err)
	}

	return count > 0, nil
}

// GetAvailableBooks retrieves all books that have at least one available copy
func (uc *availabilityUseCase) GetAvailableBooks(ctx context.Context, params repository.ListParams) ([]*BookWithAvailability, int64, error) {
	// Set default pagination values
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}

	// Get all books
	books, _, err := uc.bookRepo.List(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list books: %w", err)
	}

	// Filter books with available copies
	availableBooks := make([]*BookWithAvailability, 0)
	for _, book := range books {
		availableCount, err := uc.bookCopyRepo.CountAvailableByBookID(ctx, book.ID)
		if err != nil {
			continue
		}

		if availableCount > 0 {
			totalCount, err := uc.bookCopyRepo.CountByBookID(ctx, book.ID)
			if err != nil {
				totalCount = 0
			}

			availableBooks = append(availableBooks, &BookWithAvailability{
				Book:            book,
				AvailableCopies: availableCount,
				TotalCopies:     totalCount,
			})
		}
	}

	return availableBooks, int64(len(availableBooks)), nil
}

// GetBookCopiesByStatus retrieves book copies by their status
func (uc *availabilityUseCase) GetBookCopiesByStatus(ctx context.Context, status entity.CopyStatus, params repository.ListParams) ([]*entity.BookCopy, int64, error) {
	// Set default pagination values
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}

	copies, total, err := uc.bookCopyRepo.GetByStatus(ctx, status, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get book copies by status: %w", err)
	}

	return copies, total, nil
}

// GetExpectedReturnDate retrieves the expected return date for a borrowed book copy
func (uc *availabilityUseCase) GetExpectedReturnDate(ctx context.Context, bookCopyID string) (*time.Time, error) {
	if bookCopyID == "" {
		return nil, errors.New("book copy ID is required")
	}

	// Get book copy
	bookCopy, err := uc.bookCopyRepo.GetByID(ctx, bookCopyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book copy: %w", err)
	}

	// Check if borrowed
	if !bookCopy.IsBorrowed() {
		return nil, errors.New("book copy is not currently borrowed")
	}

	// Get active borrow record
	borrowRecord, err := uc.borrowRecordRepo.GetActiveByBookCopyID(ctx, bookCopyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get borrow record: %w", err)
	}

	return &borrowRecord.DueDate, nil
}

// GetBooksWithLowAvailability retrieves books with low availability (below threshold)
func (uc *availabilityUseCase) GetBooksWithLowAvailability(ctx context.Context, threshold int) ([]*BookAvailability, error) {
	if threshold <= 0 {
		threshold = 1 // Default threshold
	}

	// Get all books
	params := repository.ListParams{
		Page:     1,
		PageSize: 1000, // Large page size
	}

	books, _, err := uc.bookRepo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list books: %w", err)
	}

	// Filter books with low availability
	lowAvailabilityBooks := make([]*BookAvailability, 0)
	for _, book := range books {
		availability, err := uc.GetBookAvailability(ctx, book.ID)
		if err != nil {
			continue
		}

		if availability.AvailableCopies < int64(threshold) && availability.TotalCopies > 0 {
			lowAvailabilityBooks = append(lowAvailabilityBooks, availability)
		}
	}

	return lowAvailabilityBooks, nil
}

// GetAvailabilityStatistics retrieves overall availability statistics
func (uc *availabilityUseCase) GetAvailabilityStatistics(ctx context.Context) (*AvailabilityStatistics, error) {
	// Count total books
	totalBooks, err := uc.bookRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count books: %w", err)
	}

	// Get all book copies to count by status
	params := repository.ListParams{
		Page:     1,
		PageSize: 10000, // Large page size to get all
	}

	allCopies, _, err := uc.bookCopyRepo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list book copies: %w", err)
	}

	// Count copies by status
	var availableCount, borrowedCount, reservedCount, damagedCount, lostCount int64
	for _, copy := range allCopies {
		switch copy.Status {
		case entity.CopyStatusAvailable:
			availableCount++
		case entity.CopyStatusBorrowed:
			borrowedCount++
		case entity.CopyStatusReserved:
			reservedCount++
		case entity.CopyStatusDamaged:
			damagedCount++
		case entity.CopyStatusLost:
			lostCount++
		}
	}

	totalCopies := int64(len(allCopies))

	// Calculate availability rate
	var availabilityRate float64
	if totalCopies > 0 {
		availabilityRate = float64(availableCount) / float64(totalCopies) * 100
	}

	// Count reservations
	reservationParams := repository.ReservationListParams{
		Page:     1,
		PageSize: 10000,
	}

	allReservations, totalReservations, err := uc.reservationRepo.List(ctx, reservationParams)
	if err != nil {
		return nil, fmt.Errorf("failed to count reservations: %w", err)
	}

	// Count pending reservations
	var pendingReservations int64
	for _, reservation := range allReservations {
		if reservation.IsPending() {
			pendingReservations++
		}
	}

	return &AvailabilityStatistics{
		TotalBooks:          totalBooks,
		TotalCopies:         totalCopies,
		AvailableCopies:     availableCount,
		BorrowedCopies:      borrowedCount,
		ReservedCopies:      reservedCount,
		DamagedCopies:       damagedCount,
		LostCopies:          lostCount,
		AvailabilityRate:    availabilityRate,
		TotalReservations:   totalReservations,
		PendingReservations: pendingReservations,
	}, nil
}

// ReserveBook creates a reservation for a book
func (uc *availabilityUseCase) ReserveBook(ctx context.Context, userID, bookID string) (*entity.Reservation, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if bookID == "" {
		return nil, errors.New("book ID is required")
	}

	// Check if user exists and is active
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if !user.CanBorrow() {
		return nil, errors.New("user is not allowed to make reservations")
	}

	// Check if book exists
	_, err = uc.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	// Check if user already has an active reservation for this book
	hasReservation, err := uc.reservationRepo.HasActiveReservation(ctx, userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing reservation: %w", err)
	}
	if hasReservation {
		return nil, errors.New("user already has an active reservation for this book")
	}

	// Get current queue position
	pendingReservations, err := uc.reservationRepo.GetPendingByBookID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending reservations: %w", err)
	}
	queuePosition := len(pendingReservations) + 1

	// Create reservation
	now := time.Now()
	reservation := &entity.Reservation{
		UserID:          userID,
		BookID:          bookID,
		ReservationDate: now,
		ExpiryDate:      now.Add(entity.DefaultReservationHoldPeriod),
		Status:          entity.ReservationStatusPending,
		QueuePosition:   queuePosition,
	}

	// Validate reservation
	if err := reservation.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create the reservation
	if err := uc.reservationRepo.Create(ctx, reservation); err != nil {
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}

	return reservation, nil
}

// CancelReservation cancels a book reservation
func (uc *availabilityUseCase) CancelReservation(ctx context.Context, reservationID string) error {
	if reservationID == "" {
		return errors.New("reservation ID is required")
	}

	// Get reservation
	reservation, err := uc.reservationRepo.GetByID(ctx, reservationID)
	if err != nil {
		return fmt.Errorf("failed to get reservation: %w", err)
	}

	// Cancel the reservation
	if err := reservation.Cancel(); err != nil {
		return fmt.Errorf("failed to cancel reservation: %w", err)
	}

	// Update reservation status
	if err := uc.reservationRepo.UpdateStatus(ctx, reservationID, entity.ReservationStatusCancelled); err != nil {
		return fmt.Errorf("failed to update reservation status: %w", err)
	}

	return nil
}

// GetUserReservations retrieves all reservations for a user
func (uc *availabilityUseCase) GetUserReservations(ctx context.Context, userID string) ([]*entity.Reservation, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Check if user exists
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	reservations, err := uc.reservationRepo.GetPendingByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user reservations: %w", err)
	}

	return reservations, nil
}

// GetBookReservationQueue retrieves the reservation queue for a book
func (uc *availabilityUseCase) GetBookReservationQueue(ctx context.Context, bookID string) ([]*entity.Reservation, error) {
	if bookID == "" {
		return nil, errors.New("book ID is required")
	}

	// Check if book exists
	_, err := uc.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	reservations, err := uc.reservationRepo.GetPendingByBookID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation queue: %w", err)
	}

	return reservations, nil
}

// ProcessExpiredReservations processes and expires old reservations
func (uc *availabilityUseCase) ProcessExpiredReservations(ctx context.Context) error {
	// Get all expired reservations
	expiredReservations, err := uc.reservationRepo.GetExpiredReservations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get expired reservations: %w", err)
	}

	// Expire each reservation
	for _, reservation := range expiredReservations {
		if err := reservation.Expire(); err != nil {
			continue
		}

		if err := uc.reservationRepo.UpdateStatus(ctx, reservation.ID, entity.ReservationStatusExpired); err != nil {
			continue
		}
	}

	return nil
}