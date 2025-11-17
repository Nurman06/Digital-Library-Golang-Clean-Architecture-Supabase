package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/repository"
)

// BorrowingUseCase defines the interface for borrowing business logic operations
type BorrowingUseCase interface {
	// CheckoutBook processes a book checkout request
	CheckoutBook(ctx context.Context, userID, bookID string) (*entity.BorrowRecord, error)

	// ReturnBook processes a book return
	ReturnBook(ctx context.Context, borrowRecordID string) error

	// RenewBook renews a borrowed book
	RenewBook(ctx context.Context, borrowRecordID string) error

	// GetBorrowRecord retrieves a borrow record by ID
	GetBorrowRecord(ctx context.Context, id string) (*entity.BorrowRecord, error)

	// GetUserBorrowingHistory retrieves a user's borrowing history
	GetUserBorrowingHistory(ctx context.Context, userID string, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error)

	// GetActiveBorrows retrieves active borrows for a user
	GetActiveBorrows(ctx context.Context, userID string) ([]*entity.BorrowRecord, error)

	// GetOverdueBorrows retrieves overdue borrows for a user
	GetOverdueBorrows(ctx context.Context, userID string) ([]*entity.BorrowRecord, error)

	// GetAllOverdueRecords retrieves all overdue borrow records
	GetAllOverdueRecords(ctx context.Context, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error)

	// GetDueSoonRecords retrieves records that are due soon
	GetDueSoonRecords(ctx context.Context, days int, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error)

	// GetUserBorrowingStatistics retrieves borrowing statistics for a user
	GetUserBorrowingStatistics(ctx context.Context, userID string) (*repository.UserBorrowStatistics, error)

	// CalculateUserLateFees calculates total unpaid late fees for a user
	CalculateUserLateFees(ctx context.Context, userID string) (float64, error)

	// CanUserBorrowMore checks if a user can borrow more books
	CanUserBorrowMore(ctx context.Context, userID string) (bool, error)

	// GetMostBorrowedBooks retrieves statistics of most borrowed books
	GetMostBorrowedBooks(ctx context.Context, limit int, startDate, endDate time.Time) ([]repository.BorrowStatistic, error)

	// UpdateOverdueStatuses updates the status of overdue records
	UpdateOverdueStatuses(ctx context.Context) error
}

// borrowingUseCase implements the BorrowingUseCase interface
type borrowingUseCase struct {
	borrowRecordRepo repository.BorrowRecordRepository
	bookCopyRepo     repository.BookCopyRepository
	userRepo         repository.UserRepository
	bookRepo         repository.BookRepository
}

// NewBorrowingUseCase creates a new instance of BorrowingUseCase
func NewBorrowingUseCase(
	borrowRecordRepo repository.BorrowRecordRepository,
	bookCopyRepo repository.BookCopyRepository,
	userRepo repository.UserRepository,
	bookRepo repository.BookRepository,
) BorrowingUseCase {
	return &borrowingUseCase{
		borrowRecordRepo: borrowRecordRepo,
		bookCopyRepo:     bookCopyRepo,
		userRepo:         userRepo,
		bookRepo:         bookRepo,
	}
}

// CheckoutBook processes a book checkout request
func (uc *borrowingUseCase) CheckoutBook(ctx context.Context, userID, bookID string) (*entity.BorrowRecord, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if bookID == "" {
		return nil, errors.New("book ID is required")
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Validate user can borrow
	if !user.CanBorrow() {
		return nil, errors.New("user is not allowed to borrow books")
	}

	// Check if user has reached borrowing limit
	activeBorrowCount, err := uc.borrowRecordRepo.CountActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count active borrows: %w", err)
	}

	if !user.CanBorrowMoreBooks(int(activeBorrowCount)) {
		return nil, fmt.Errorf("user has reached borrowing limit of %d books", user.BorrowingLimit)
	}

	// Check if user has overdue books
	overdueBorrows, err := uc.borrowRecordRepo.GetOverdueByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check overdue borrows: %w", err)
	}
	if len(overdueBorrows) > 0 {
		return nil, errors.New("user has overdue books and cannot borrow more")
	}

	// Check if user has unpaid late fees
	totalLateFees, err := uc.borrowRecordRepo.CalculateTotalLateFees(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate late fees: %w", err)
	}
	if totalLateFees > 0 {
		return nil, fmt.Errorf("user has unpaid late fees of $%.2f", totalLateFees)
	}

	// Get book
	book, err := uc.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	// Check if book is deleted
	if book.IsDeleted() {
		return nil, errors.New("book is not available")
	}

	// Find an available copy
	availableCopies, err := uc.bookCopyRepo.GetAvailableCopies(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get available copies: %w", err)
	}

	if len(availableCopies) == 0 {
		return nil, errors.New("no available copies of this book")
	}

	// Use the first available copy
	bookCopy := availableCopies[0]

	// Mark book copy as borrowed
	if err := bookCopy.MarkAsBorrowed(); err != nil {
		return nil, fmt.Errorf("failed to mark book copy as borrowed: %w", err)
	}

	// Update book copy status
	if err := uc.bookCopyRepo.UpdateStatus(ctx, bookCopy.ID, entity.CopyStatusBorrowed); err != nil {
		return nil, fmt.Errorf("failed to update book copy status: %w", err)
	}

	// Create borrow record
	now := time.Now()
	borrowRecord := &entity.BorrowRecord{
		UserID:       userID,
		BookCopyID:   bookCopy.ID,
		CheckoutDate: now,
		DueDate:      now.Add(entity.DefaultBorrowingPeriod),
		RenewalCount: 0,
		LateFee:      0,
		Status:       entity.BorrowStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Validate borrow record
	if err := borrowRecord.Validate(); err != nil {
		// Rollback: Mark book copy as available
		_ = uc.bookCopyRepo.UpdateStatus(ctx, bookCopy.ID, entity.CopyStatusAvailable)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create the borrow record
	if err := uc.borrowRecordRepo.Create(ctx, borrowRecord); err != nil {
		// Rollback: Mark book copy as available
		_ = uc.bookCopyRepo.UpdateStatus(ctx, bookCopy.ID, entity.CopyStatusAvailable)
		return nil, fmt.Errorf("failed to create borrow record: %w", err)
	}

	return borrowRecord, nil
}

// ReturnBook processes a book return
func (uc *borrowingUseCase) ReturnBook(ctx context.Context, borrowRecordID string) error {
	if borrowRecordID == "" {
		return errors.New("borrow record ID is required")
	}

	// Get borrow record
	borrowRecord, err := uc.borrowRecordRepo.GetByID(ctx, borrowRecordID)
	if err != nil {
		return fmt.Errorf("failed to get borrow record: %w", err)
	}

	// Check if book is already returned
	if borrowRecord.IsReturned() {
		return errors.New("book has already been returned")
	}

	// Get book copy
	bookCopy, err := uc.bookCopyRepo.GetByID(ctx, borrowRecord.BookCopyID)
	if err != nil {
		return fmt.Errorf("failed to get book copy: %w", err)
	}

	// Process return in entity
	if err := borrowRecord.Return(); err != nil {
		return fmt.Errorf("failed to process return: %w", err)
	}

	// Mark book copy as available
	if err := bookCopy.MarkAsAvailable(); err != nil {
		return fmt.Errorf("failed to mark book copy as available: %w", err)
	}

	// Update borrow record
	if err := uc.borrowRecordRepo.Update(ctx, borrowRecord); err != nil {
		return fmt.Errorf("failed to update borrow record: %w", err)
	}

	// Update book copy status
	if err := uc.bookCopyRepo.UpdateStatus(ctx, bookCopy.ID, entity.CopyStatusAvailable); err != nil {
		return fmt.Errorf("failed to update book copy status: %w", err)
	}

	return nil
}

// RenewBook renews a borrowed book
func (uc *borrowingUseCase) RenewBook(ctx context.Context, borrowRecordID string) error {
	if borrowRecordID == "" {
		return errors.New("borrow record ID is required")
	}

	// Get borrow record
	borrowRecord, err := uc.borrowRecordRepo.GetByID(ctx, borrowRecordID)
	if err != nil {
		return fmt.Errorf("failed to get borrow record: %w", err)
	}

	// Check if book can be renewed
	if !borrowRecord.CanRenew() {
		if borrowRecord.IsOverdue() {
			return errors.New("cannot renew overdue book")
		}
		if borrowRecord.RenewalCount >= entity.MaxRenewalCount {
			return fmt.Errorf("maximum renewal limit of %d reached", entity.MaxRenewalCount)
		}
		return errors.New("book cannot be renewed")
	}

	// Renew the book
	if err := borrowRecord.Renew(); err != nil {
		return fmt.Errorf("failed to renew book: %w", err)
	}

	// Update borrow record
	if err := uc.borrowRecordRepo.Update(ctx, borrowRecord); err != nil {
		return fmt.Errorf("failed to update borrow record: %w", err)
	}

	return nil
}

// GetBorrowRecord retrieves a borrow record by ID
func (uc *borrowingUseCase) GetBorrowRecord(ctx context.Context, id string) (*entity.BorrowRecord, error) {
	if id == "" {
		return nil, errors.New("borrow record ID is required")
	}

	borrowRecord, err := uc.borrowRecordRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get borrow record: %w", err)
	}

	return borrowRecord, nil
}

// GetUserBorrowingHistory retrieves a user's borrowing history
func (uc *borrowingUseCase) GetUserBorrowingHistory(ctx context.Context, userID string, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	if userID == "" {
		return nil, 0, errors.New("user ID is required")
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

	// Check if user exists
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user: %w", err)
	}

	records, total, err := uc.borrowRecordRepo.GetUserBorrowingHistory(ctx, userID, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get borrowing history: %w", err)
	}

	return records, total, nil
}

// GetActiveBorrows retrieves active borrows for a user
func (uc *borrowingUseCase) GetActiveBorrows(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Check if user exists
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	records, err := uc.borrowRecordRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active borrows: %w", err)
	}

	return records, nil
}

// GetOverdueBorrows retrieves overdue borrows for a user
func (uc *borrowingUseCase) GetOverdueBorrows(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Check if user exists
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	records, err := uc.borrowRecordRepo.GetOverdueByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue borrows: %w", err)
	}

	return records, nil
}

// GetAllOverdueRecords retrieves all overdue borrow records
func (uc *borrowingUseCase) GetAllOverdueRecords(ctx context.Context, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
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

	records, total, err := uc.borrowRecordRepo.GetOverdueRecords(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get overdue records: %w", err)
	}

	return records, total, nil
}

// GetDueSoonRecords retrieves records that are due soon
func (uc *borrowingUseCase) GetDueSoonRecords(ctx context.Context, days int, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	if days <= 0 {
		days = 3 // Default to 3 days
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

	records, total, err := uc.borrowRecordRepo.GetDueSoon(ctx, days, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get due soon records: %w", err)
	}

	return records, total, nil
}

// GetUserBorrowingStatistics retrieves borrowing statistics for a user
func (uc *borrowingUseCase) GetUserBorrowingStatistics(ctx context.Context, userID string) (*repository.UserBorrowStatistics, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Check if user exists
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	stats, err := uc.borrowRecordRepo.GetUserBorrowingStatistics(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get borrowing statistics: %w", err)
	}

	return stats, nil
}

// CalculateUserLateFees calculates total unpaid late fees for a user
func (uc *borrowingUseCase) CalculateUserLateFees(ctx context.Context, userID string) (float64, error) {
	if userID == "" {
		return 0, errors.New("user ID is required")
	}

	// Check if user exists
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user: %w", err)
	}

	totalFees, err := uc.borrowRecordRepo.CalculateTotalLateFees(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate late fees: %w", err)
	}

	return totalFees, nil
}

// CanUserBorrowMore checks if a user can borrow more books
func (uc *borrowingUseCase) CanUserBorrowMore(ctx context.Context, userID string) (bool, error) {
	if userID == "" {
		return false, errors.New("user ID is required")
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user can borrow
	if !user.CanBorrow() {
		return false, nil
	}

	// Count active borrows
	activeBorrowCount, err := uc.borrowRecordRepo.CountActiveByUserID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to count active borrows: %w", err)
	}

	// Check borrowing limit
	canBorrowMore := user.CanBorrowMoreBooks(int(activeBorrowCount))

	// Check for overdue books
	overdueBorrows, err := uc.borrowRecordRepo.GetOverdueByUserID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check overdue borrows: %w", err)
	}
	if len(overdueBorrows) > 0 {
		return false, nil
	}

	// Check for unpaid late fees
	totalLateFees, err := uc.borrowRecordRepo.CalculateTotalLateFees(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to calculate late fees: %w", err)
	}
	if totalLateFees > 0 {
		return false, nil
	}

	return canBorrowMore, nil
}

// GetMostBorrowedBooks retrieves statistics of most borrowed books
func (uc *borrowingUseCase) GetMostBorrowedBooks(ctx context.Context, limit int, startDate, endDate time.Time) ([]repository.BorrowStatistic, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	stats, err := uc.borrowRecordRepo.GetMostBorrowedBooks(ctx, limit, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get most borrowed books: %w", err)
	}

	return stats, nil
}

// UpdateOverdueStatuses updates the status of overdue records
func (uc *borrowingUseCase) UpdateOverdueStatuses(ctx context.Context) error {
	// Get all active borrow records
	params := repository.BorrowRecordListParams{
		Page:     1,
		PageSize: 1000, // Process in batches
	}

	records, _, err := uc.borrowRecordRepo.GetByStatus(ctx, entity.BorrowStatusActive, params)
	if err != nil {
		return fmt.Errorf("failed to get active records: %w", err)
	}

	// Update status for overdue records
	for _, record := range records {
		if record.IsOverdue() {
			record.UpdateStatus()
			if err := uc.borrowRecordRepo.UpdateStatus(ctx, record.ID, entity.BorrowStatusOverdue); err != nil {
				// Log error but continue processing
				continue
			}
		}
	}

	return nil
}