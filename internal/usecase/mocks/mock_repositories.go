package mocks

import (
	"context"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/repository"
)

// MockBookRepository is a mock implementation of BookRepository
type MockBookRepository struct {
	CreateFunc            func(ctx context.Context, book *entity.Book) error
	GetByIDFunc           func(ctx context.Context, id string) (*entity.Book, error)
	GetByISBNFunc         func(ctx context.Context, isbn string) (*entity.Book, error)
	UpdateFunc            func(ctx context.Context, book *entity.Book) error
	DeleteFunc            func(ctx context.Context, id string) error
	RestoreFunc           func(ctx context.Context, id string) error
	ListFunc              func(ctx context.Context, params repository.ListParams) ([]*entity.Book, int64, error)
	SearchFunc            func(ctx context.Context, params repository.SearchParams) ([]*entity.Book, int64, error)
	GetByCategoryFunc     func(ctx context.Context, category string, params repository.ListParams) ([]*entity.Book, int64, error)
	GetByAuthorFunc       func(ctx context.Context, author string, params repository.ListParams) ([]*entity.Book, int64, error)
	ExistsByISBNFunc      func(ctx context.Context, isbn string) (bool, error)
	CountFunc             func(ctx context.Context) (int64, error)
	GetRecentlyAddedFunc  func(ctx context.Context, limit int) ([]*entity.Book, error)
}

func (m *MockBookRepository) Create(ctx context.Context, book *entity.Book) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, book)
	}
	return nil
}

func (m *MockBookRepository) GetByID(ctx context.Context, id string) (*entity.Book, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockBookRepository) GetByISBN(ctx context.Context, isbn string) (*entity.Book, error) {
	if m.GetByISBNFunc != nil {
		return m.GetByISBNFunc(ctx, isbn)
	}
	return nil, nil
}

func (m *MockBookRepository) Update(ctx context.Context, book *entity.Book) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, book)
	}
	return nil
}

func (m *MockBookRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockBookRepository) HardDelete(ctx context.Context, id string) error {
	return nil
}

func (m *MockBookRepository) Restore(ctx context.Context, id string) error {
	if m.RestoreFunc != nil {
		return m.RestoreFunc(ctx, id)
	}
	return nil
}

func (m *MockBookRepository) List(ctx context.Context, params repository.ListParams) ([]*entity.Book, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, params)
	}
	return nil, 0, nil
}

func (m *MockBookRepository) Search(ctx context.Context, params repository.SearchParams) ([]*entity.Book, int64, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(ctx, params)
	}
	return nil, 0, nil
}

func (m *MockBookRepository) GetByCategory(ctx context.Context, category string, params repository.ListParams) ([]*entity.Book, int64, error) {
	if m.GetByCategoryFunc != nil {
		return m.GetByCategoryFunc(ctx, category, params)
	}
	return nil, 0, nil
}

func (m *MockBookRepository) GetByAuthor(ctx context.Context, author string, params repository.ListParams) ([]*entity.Book, int64, error) {
	if m.GetByAuthorFunc != nil {
		return m.GetByAuthorFunc(ctx, author, params)
	}
	return nil, 0, nil
}

func (m *MockBookRepository) ExistsByISBN(ctx context.Context, isbn string) (bool, error) {
	if m.ExistsByISBNFunc != nil {
		return m.ExistsByISBNFunc(ctx, isbn)
	}
	return false, nil
}

func (m *MockBookRepository) Count(ctx context.Context) (int64, error) {
	if m.CountFunc != nil {
		return m.CountFunc(ctx)
	}
	return 0, nil
}

func (m *MockBookRepository) GetRecentlyAdded(ctx context.Context, limit int) ([]*entity.Book, error) {
	if m.GetRecentlyAddedFunc != nil {
		return m.GetRecentlyAddedFunc(ctx, limit)
	}
	return nil, nil
}

// MockBookCopyRepository is a mock implementation of BookCopyRepository
type MockBookCopyRepository struct {
	CreateFunc                 func(ctx context.Context, copy *entity.BookCopy) error
	GetByIDFunc                func(ctx context.Context, id string) (*entity.BookCopy, error)
	GetByCopyNumberFunc        func(ctx context.Context, copyNumber string) (*entity.BookCopy, error)
	UpdateFunc                 func(ctx context.Context, copy *entity.BookCopy) error
	DeleteFunc                 func(ctx context.Context, id string) error
	GetByBookIDFunc            func(ctx context.Context, bookID string) ([]*entity.BookCopy, error)
	GetAvailableCopiesFunc     func(ctx context.Context, bookID string) ([]*entity.BookCopy, error)
	GetByStatusFunc            func(ctx context.Context, status entity.CopyStatus, params repository.ListParams) ([]*entity.BookCopy, int64, error)
	CountByBookIDFunc          func(ctx context.Context, bookID string) (int64, error)
	CountAvailableByBookIDFunc func(ctx context.Context, bookID string) (int64, error)
	UpdateStatusFunc           func(ctx context.Context, id string, status entity.CopyStatus) error
	ExistsByCopyNumberFunc     func(ctx context.Context, copyNumber string) (bool, error)
	GetByLocationFunc          func(ctx context.Context, location string, params repository.ListParams) ([]*entity.BookCopy, int64, error)
	ListFunc                   func(ctx context.Context, params repository.ListParams) ([]*entity.BookCopy, int64, error)
}

func (m *MockBookCopyRepository) Create(ctx context.Context, copy *entity.BookCopy) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, copy)
	}
	return nil
}

func (m *MockBookCopyRepository) GetByID(ctx context.Context, id string) (*entity.BookCopy, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockBookCopyRepository) GetByCopyNumber(ctx context.Context, copyNumber string) (*entity.BookCopy, error) {
	if m.GetByCopyNumberFunc != nil {
		return m.GetByCopyNumberFunc(ctx, copyNumber)
	}
	return nil, nil
}

func (m *MockBookCopyRepository) Update(ctx context.Context, copy *entity.BookCopy) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, copy)
	}
	return nil
}

func (m *MockBookCopyRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockBookCopyRepository) GetByBookID(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
	if m.GetByBookIDFunc != nil {
		return m.GetByBookIDFunc(ctx, bookID)
	}
	return nil, nil
}

func (m *MockBookCopyRepository) GetAvailableCopies(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
	if m.GetAvailableCopiesFunc != nil {
		return m.GetAvailableCopiesFunc(ctx, bookID)
	}
	return nil, nil
}

func (m *MockBookCopyRepository) GetByStatus(ctx context.Context, status entity.CopyStatus, params repository.ListParams) ([]*entity.BookCopy, int64, error) {
	if m.GetByStatusFunc != nil {
		return m.GetByStatusFunc(ctx, status, params)
	}
	return nil, 0, nil
}

func (m *MockBookCopyRepository) CountByBookID(ctx context.Context, bookID string) (int64, error) {
	if m.CountByBookIDFunc != nil {
		return m.CountByBookIDFunc(ctx, bookID)
	}
	return 0, nil
}

func (m *MockBookCopyRepository) CountAvailableByBookID(ctx context.Context, bookID string) (int64, error) {
	if m.CountAvailableByBookIDFunc != nil {
		return m.CountAvailableByBookIDFunc(ctx, bookID)
	}
	return 0, nil
}

func (m *MockBookCopyRepository) UpdateStatus(ctx context.Context, id string, status entity.CopyStatus) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, id, status)
	}
	return nil
}

func (m *MockBookCopyRepository) ExistsByCopyNumber(ctx context.Context, copyNumber string) (bool, error) {
	if m.ExistsByCopyNumberFunc != nil {
		return m.ExistsByCopyNumberFunc(ctx, copyNumber)
	}
	return false, nil
}

func (m *MockBookCopyRepository) GetByLocation(ctx context.Context, location string, params repository.ListParams) ([]*entity.BookCopy, int64, error) {
	if m.GetByLocationFunc != nil {
		return m.GetByLocationFunc(ctx, location, params)
	}
	return nil, 0, nil
}

func (m *MockBookCopyRepository) List(ctx context.Context, params repository.ListParams) ([]*entity.BookCopy, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, params)
	}
	return nil, 0, nil
}

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	CreateFunc  func(ctx context.Context, user *entity.User) error
	GetByIDFunc func(ctx context.Context, id string) (*entity.User, error)
	UpdateFunc  func(ctx context.Context, user *entity.User) error
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *MockUserRepository) List(ctx context.Context, params repository.UserListParams) ([]*entity.User, int64, error) {
	return nil, 0, nil
}

func (m *MockUserRepository) GetByRole(ctx context.Context, role entity.UserRole, params repository.UserListParams) ([]*entity.User, int64, error) {
	return nil, 0, nil
}

func (m *MockUserRepository) GetByStatus(ctx context.Context, status entity.UserStatus, params repository.UserListParams) ([]*entity.User, int64, error) {
	return nil, 0, nil
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *MockUserRepository) CountByRole(ctx context.Context, role entity.UserRole) (int64, error) {
	return 0, nil
}

func (m *MockUserRepository) CountByStatus(ctx context.Context, status entity.UserStatus) (int64, error) {
	return 0, nil
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, nil
}

func (m *MockUserRepository) UpdateStatus(ctx context.Context, id string, status entity.UserStatus) error {
	return nil
}

func (m *MockUserRepository) UpdateBorrowingLimit(ctx context.Context, id string, limit int) error {
	return nil
}

func (m *MockUserRepository) GetActiveMembers(ctx context.Context, params repository.UserListParams) ([]*entity.User, int64, error) {
	return nil, 0, nil
}

func (m *MockUserRepository) Search(ctx context.Context, query string, params repository.UserListParams) ([]*entity.User, int64, error) {
	return nil, 0, nil
}

// MockBorrowRecordRepository is a mock implementation of BorrowRecordRepository
type MockBorrowRecordRepository struct {
	CreateFunc                    func(ctx context.Context, record *entity.BorrowRecord) error
	GetByIDFunc                   func(ctx context.Context, id string) (*entity.BorrowRecord, error)
	UpdateFunc                    func(ctx context.Context, record *entity.BorrowRecord) error
	GetActiveByUserIDFunc         func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error)
	GetOverdueByUserIDFunc        func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error)
	CountActiveByUserIDFunc       func(ctx context.Context, userID string) (int64, error)
	CalculateTotalLateFeesFunc    func(ctx context.Context, userID string) (float64, error)
	GetUserBorrowingHistoryFunc   func(ctx context.Context, userID string, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error)
	GetOverdueRecordsFunc         func(ctx context.Context, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error)
	GetDueSoonFunc                func(ctx context.Context, days int, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error)
	GetByStatusFunc               func(ctx context.Context, status entity.BorrowRecordStatus, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error)
	UpdateStatusFunc              func(ctx context.Context, id string, status entity.BorrowRecordStatus) error
	GetUserBorrowingStatisticsFunc func(ctx context.Context, userID string) (*repository.UserBorrowStatistics, error)
	GetMostBorrowedBooksFunc      func(ctx context.Context, limit int, startDate, endDate interface{}) ([]repository.BorrowStatistic, error)
}

func (m *MockBorrowRecordRepository) Create(ctx context.Context, record *entity.BorrowRecord) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, record)
	}
	return nil
}

func (m *MockBorrowRecordRepository) GetByID(ctx context.Context, id string) (*entity.BorrowRecord, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockBorrowRecordRepository) Update(ctx context.Context, record *entity.BorrowRecord) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, record)
	}
	return nil
}

func (m *MockBorrowRecordRepository) GetActiveByUserID(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
	if m.GetActiveByUserIDFunc != nil {
		return m.GetActiveByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockBorrowRecordRepository) GetOverdueByUserID(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
	if m.GetOverdueByUserIDFunc != nil {
		return m.GetOverdueByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockBorrowRecordRepository) CountActiveByUserID(ctx context.Context, userID string) (int64, error) {
	if m.CountActiveByUserIDFunc != nil {
		return m.CountActiveByUserIDFunc(ctx, userID)
	}
	return 0, nil
}

func (m *MockBorrowRecordRepository) CalculateTotalLateFees(ctx context.Context, userID string) (float64, error) {
	if m.CalculateTotalLateFeesFunc != nil {
		return m.CalculateTotalLateFeesFunc(ctx, userID)
	}
	return 0, nil
}

func (m *MockBorrowRecordRepository) GetUserBorrowingHistory(ctx context.Context, userID string, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	if m.GetUserBorrowingHistoryFunc != nil {
		return m.GetUserBorrowingHistoryFunc(ctx, userID, params)
	}
	return nil, 0, nil
}

func (m *MockBorrowRecordRepository) GetOverdueRecords(ctx context.Context, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	if m.GetOverdueRecordsFunc != nil {
		return m.GetOverdueRecordsFunc(ctx, params)
	}
	return nil, 0, nil
}

func (m *MockBorrowRecordRepository) GetDueSoon(ctx context.Context, days int, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	if m.GetDueSoonFunc != nil {
		return m.GetDueSoonFunc(ctx, days, params)
	}
	return nil, 0, nil
}

func (m *MockBorrowRecordRepository) GetByStatus(ctx context.Context, status entity.BorrowRecordStatus, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	if m.GetByStatusFunc != nil {
		return m.GetByStatusFunc(ctx, status, params)
	}
	return nil, 0, nil
}

func (m *MockBorrowRecordRepository) UpdateStatus(ctx context.Context, id string, status entity.BorrowRecordStatus) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, id, status)
	}
	return nil
}

func (m *MockBorrowRecordRepository) GetUserBorrowingStatistics(ctx context.Context, userID string) (*repository.UserBorrowStatistics, error) {
	if m.GetUserBorrowingStatisticsFunc != nil {
		return m.GetUserBorrowingStatisticsFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockBorrowRecordRepository) GetMostBorrowedBooks(ctx context.Context, limit int, startDate, endDate time.Time) ([]repository.BorrowStatistic, error) {
	return nil, nil
}

func (m *MockBorrowRecordRepository) CountOverdueByUserID(ctx context.Context, userID string) (int64, error) {
	return 0, nil
}

func (m *MockBorrowRecordRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *MockBorrowRecordRepository) List(ctx context.Context, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	return nil, 0, nil
}

func (m *MockBorrowRecordRepository) GetByUserID(ctx context.Context, userID string, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	return nil, 0, nil
}

func (m *MockBorrowRecordRepository) GetByBookCopyID(ctx context.Context, bookCopyID string, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	return nil, 0, nil
}

func (m *MockBorrowRecordRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	return nil, 0, nil
}

func (m *MockBorrowRecordRepository) GetBookBorrowingHistory(ctx context.Context, bookCopyID string, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	return nil, 0, nil
}

func (m *MockBorrowRecordRepository) MarkAsReturned(ctx context.Context, id string, returnDate time.Time) error {
	return nil
}

func (m *MockBorrowRecordRepository) HasActiveBookCopyBorrow(ctx context.Context, bookCopyID string) (bool, error) {
	return false, nil
}

func (m *MockBorrowRecordRepository) GetActiveByBookCopyID(ctx context.Context, bookCopyID string) (*entity.BorrowRecord, error) {
	return nil, nil
}