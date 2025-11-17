package repository

import (
	"context"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
)

// BorrowRecordRepository defines the interface for borrow record data access operations
type BorrowRecordRepository interface {
	// Create creates a new borrow record in the repository
	Create(ctx context.Context, record *entity.BorrowRecord) error

	// GetByID retrieves a borrow record by its ID
	GetByID(ctx context.Context, id string) (*entity.BorrowRecord, error)

	// Update updates an existing borrow record
	Update(ctx context.Context, record *entity.BorrowRecord) error

	// Delete deletes a borrow record from the repository
	Delete(ctx context.Context, id string) error

	// GetByUserID retrieves all borrow records for a specific user
	GetByUserID(ctx context.Context, userID string, params BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error)

	// GetByBookCopyID retrieves all borrow records for a specific book copy
	GetByBookCopyID(ctx context.Context, bookCopyID string, params BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error)

	// GetByStatus retrieves borrow records by status
	GetByStatus(ctx context.Context, status entity.BorrowRecordStatus, params BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error)

	// GetActiveByUserID retrieves active borrow records for a specific user
	GetActiveByUserID(ctx context.Context, userID string) ([]*entity.BorrowRecord, error)

	// GetOverdueRecords retrieves all overdue borrow records
	GetOverdueRecords(ctx context.Context, params BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error)

	// GetOverdueByUserID retrieves overdue borrow records for a specific user
	GetOverdueByUserID(ctx context.Context, userID string) ([]*entity.BorrowRecord, error)

	// GetDueSoon retrieves borrow records that are due within the specified days
	GetDueSoon(ctx context.Context, days int, params BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error)

	// List retrieves a paginated list of borrow records
	List(ctx context.Context, params BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error)

	// CountActiveByUserID counts active borrow records for a specific user
	CountActiveByUserID(ctx context.Context, userID string) (int64, error)

	// CountOverdueByUserID counts overdue borrow records for a specific user
	CountOverdueByUserID(ctx context.Context, userID string) (int64, error)

	// GetByDateRange retrieves borrow records within a date range
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, params BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error)

	// GetUserBorrowingHistory retrieves complete borrowing history for a user
	GetUserBorrowingHistory(ctx context.Context, userID string, params BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error)

	// GetBookBorrowingHistory retrieves complete borrowing history for a book copy
	GetBookBorrowingHistory(ctx context.Context, bookCopyID string, params BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error)

	// UpdateStatus updates the status of a borrow record
	UpdateStatus(ctx context.Context, id string, status entity.BorrowRecordStatus) error

	// MarkAsReturned marks a borrow record as returned
	MarkAsReturned(ctx context.Context, id string, returnDate time.Time) error

	// CalculateTotalLateFees calculates total unpaid late fees for a user
	CalculateTotalLateFees(ctx context.Context, userID string) (float64, error)

	// GetMostBorrowedBooks retrieves statistics of most borrowed books
	GetMostBorrowedBooks(ctx context.Context, limit int, startDate, endDate time.Time) ([]BorrowStatistic, error)

	// GetUserBorrowingStatistics retrieves borrowing statistics for a user
	GetUserBorrowingStatistics(ctx context.Context, userID string) (*UserBorrowStatistics, error)

	// HasActiveBookCopyBorrow checks if a book copy has an active borrow record
	HasActiveBookCopyBorrow(ctx context.Context, bookCopyID string) (bool, error)

	// GetActiveByBookCopyID retrieves active borrow record for a specific book copy
	GetActiveByBookCopyID(ctx context.Context, bookCopyID string) (*entity.BorrowRecord, error)
}

// BorrowRecordListParams defines parameters for listing borrow records
type BorrowRecordListParams struct {
	Page      int
	PageSize  int
	SortBy    string // e.g., "checkout_date", "due_date", "return_date"
	SortOrder string // "asc" or "desc"
}

// BorrowStatistic represents borrowing statistics for a book
type BorrowStatistic struct {
	BookID      string
	BookTitle   string
	BorrowCount int64
	TotalDays   int64
}

// UserBorrowStatistics represents borrowing statistics for a user
type UserBorrowStatistics struct {
	TotalBorrows    int64
	ActiveBorrows   int64
	OverdueBorrows  int64
	ReturnedBorrows int64
	TotalLateFees   float64
	AverageDays     float64
}

// ReservationRepository defines the interface for reservation data access operations
type ReservationRepository interface {
	// Create creates a new reservation in the repository
	Create(ctx context.Context, reservation *entity.Reservation) error

	// GetByID retrieves a reservation by its ID
	GetByID(ctx context.Context, id string) (*entity.Reservation, error)

	// Update updates an existing reservation
	Update(ctx context.Context, reservation *entity.Reservation) error

	// Delete deletes a reservation from the repository
	Delete(ctx context.Context, id string) error

	// GetByUserID retrieves all reservations for a specific user
	GetByUserID(ctx context.Context, userID string, params ReservationListParams) ([]*entity.Reservation, int64, error)

	// GetByBookID retrieves all reservations for a specific book
	GetByBookID(ctx context.Context, bookID string, params ReservationListParams) ([]*entity.Reservation, int64, error)

	// GetByStatus retrieves reservations by status
	GetByStatus(ctx context.Context, status entity.ReservationStatus, params ReservationListParams) ([]*entity.Reservation, int64, error)

	// GetPendingByUserID retrieves pending reservations for a specific user
	GetPendingByUserID(ctx context.Context, userID string) ([]*entity.Reservation, error)

	// GetPendingByBookID retrieves pending reservations for a specific book ordered by queue position
	GetPendingByBookID(ctx context.Context, bookID string) ([]*entity.Reservation, error)

	// GetExpiredReservations retrieves all expired reservations
	GetExpiredReservations(ctx context.Context) ([]*entity.Reservation, error)

	// GetNextInQueue retrieves the next reservation in queue for a book
	GetNextInQueue(ctx context.Context, bookID string) (*entity.Reservation, error)

	// UpdateStatus updates the status of a reservation
	UpdateStatus(ctx context.Context, id string, status entity.ReservationStatus) error

	// UpdateQueuePosition updates the queue position of a reservation
	UpdateQueuePosition(ctx context.Context, id string, position int) error

	// CountPendingByBookID counts pending reservations for a specific book
	CountPendingByBookID(ctx context.Context, bookID string) (int64, error)

	// CountPendingByUserID counts pending reservations for a specific user
	CountPendingByUserID(ctx context.Context, userID string) (int64, error)

	// List retrieves a paginated list of reservations
	List(ctx context.Context, params ReservationListParams) ([]*entity.Reservation, int64, error)

	// HasActiveReservation checks if a user has an active reservation for a book
	HasActiveReservation(ctx context.Context, userID, bookID string) (bool, error)
}

// ReservationListParams defines parameters for listing reservations
type ReservationListParams struct {
	Page      int
	PageSize  int
	SortBy    string // e.g., "reservation_date", "queue_position", "expiry_date"
	SortOrder string // "asc" or "desc"
}
