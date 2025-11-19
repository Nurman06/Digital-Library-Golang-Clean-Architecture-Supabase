package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockBorrowRecordRepository is a mock implementation of BorrowRecordRepository
type MockBorrowRecordRepository struct {
	mock.Mock
}

func (m *MockBorrowRecordRepository) Create(ctx context.Context, record *entity.BorrowRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockBorrowRecordRepository) GetByID(ctx context.Context, id string) (*entity.BorrowRecord, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BorrowRecord), args.Error(1)
}

func (m *MockBorrowRecordRepository) Update(ctx context.Context, record *entity.BorrowRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockBorrowRecordRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBorrowRecordRepository) GetByUserID(ctx context.Context, userID string, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	args := m.Called(ctx, userID, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.BorrowRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockBorrowRecordRepository) GetActiveByUserID(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.BorrowRecord), args.Error(1)
}

func (m *MockBorrowRecordRepository) GetOverdueByUserID(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.BorrowRecord), args.Error(1)
}

func (m *MockBorrowRecordRepository) CountActiveByUserID(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBorrowRecordRepository) CountOverdueByUserID(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBorrowRecordRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	args := m.Called(ctx, startDate, endDate, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.BorrowRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockBorrowRecordRepository) GetBookBorrowingHistory(ctx context.Context, bookCopyID string, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	args := m.Called(ctx, bookCopyID, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.BorrowRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockBorrowRecordRepository) MarkAsReturned(ctx context.Context, id string, returnDate time.Time) error {
	args := m.Called(ctx, id, returnDate)
	return args.Error(0)
}

func (m *MockBorrowRecordRepository) HasActiveBookCopyBorrow(ctx context.Context, bookCopyID string) (bool, error) {
	args := m.Called(ctx, bookCopyID)
	return args.Bool(0), args.Error(1)
}

func (m *MockBorrowRecordRepository) GetActiveByBookCopyID(ctx context.Context, bookCopyID string) (*entity.BorrowRecord, error) {
	args := m.Called(ctx, bookCopyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BorrowRecord), args.Error(1)
}

func (m *MockBorrowRecordRepository) GetByBookCopyID(ctx context.Context, bookCopyID string, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	args := m.Called(ctx, bookCopyID, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.BorrowRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockBorrowRecordRepository) List(ctx context.Context, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.BorrowRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockBorrowRecordRepository) GetUserBorrowingHistory(ctx context.Context, userID string, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	args := m.Called(ctx, userID, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.BorrowRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockBorrowRecordRepository) GetByStatus(ctx context.Context, status entity.BorrowRecordStatus, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	args := m.Called(ctx, status, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.BorrowRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockBorrowRecordRepository) GetOverdueRecords(ctx context.Context, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.BorrowRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockBorrowRecordRepository) GetDueSoon(ctx context.Context, days int, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	args := m.Called(ctx, days, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.BorrowRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockBorrowRecordRepository) UpdateStatus(ctx context.Context, id string, status entity.BorrowRecordStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockBorrowRecordRepository) CalculateTotalLateFees(ctx context.Context, userID string) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockBorrowRecordRepository) GetUserBorrowingStatistics(ctx context.Context, userID string) (*repository.UserBorrowStatistics, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.UserBorrowStatistics), args.Error(1)
}

func (m *MockBorrowRecordRepository) GetMostBorrowedBooks(ctx context.Context, limit int, startDate, endDate time.Time) ([]repository.BorrowStatistic, error) {
	args := m.Called(ctx, limit, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.BorrowStatistic), args.Error(1)
}

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, params repository.UserListParams) ([]*entity.User, int64, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) UpdateStatus(ctx context.Context, id string, status entity.UserStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockUserRepository) GetByRole(ctx context.Context, role entity.UserRole, params repository.UserListParams) ([]*entity.User, int64, error) {
	args := m.Called(ctx, role, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) GetByStatus(ctx context.Context, status entity.UserStatus, params repository.UserListParams) ([]*entity.User, int64, error) {
	args := m.Called(ctx, status, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) UpdateBorrowingLimit(ctx context.Context, id string, limit int) error {
	args := m.Called(ctx, id, limit)
	return args.Error(0)
}

func (m *MockUserRepository) CountByRole(ctx context.Context, role entity.UserRole) (int64, error) {
	args := m.Called(ctx, role)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) CountByStatus(ctx context.Context, status entity.UserStatus) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) GetActiveMembers(ctx context.Context, params repository.UserListParams) ([]*entity.User, int64, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) Search(ctx context.Context, query string, params repository.UserListParams) ([]*entity.User, int64, error) {
	args := m.Called(ctx, query, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.User), args.Get(1).(int64), args.Error(2)
}

// Test CheckoutBook
func TestBorrowingUseCase_CheckoutBook(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		user := &entity.User{
			ID:             "user-1",
			Status:         entity.StatusActive,
			Role:           entity.RoleMember,
			BorrowingLimit: 5,
		}

		book := &entity.Book{
			ID:        "book-1",
			Title:     "Test Book",
			DeletedAt: nil,
		}

		bookCopy := &entity.BookCopy{
			ID:     "copy-1",
			BookID: "book-1",
			Status: entity.CopyStatusAvailable,
		}

		mockUserRepo.On("GetByID", ctx, "user-1").Return(user, nil)
		mockBorrowRepo.On("CountActiveByUserID", ctx, "user-1").Return(int64(0), nil)
		mockBorrowRepo.On("GetOverdueByUserID", ctx, "user-1").Return([]*entity.BorrowRecord{}, nil)
		mockBorrowRepo.On("CalculateTotalLateFees", ctx, "user-1").Return(0.0, nil)
		mockBookRepo.On("GetByID", ctx, "book-1").Return(book, nil)
		mockCopyRepo.On("GetAvailableCopies", ctx, "book-1").Return([]*entity.BookCopy{bookCopy}, nil)
		mockCopyRepo.On("UpdateStatus", ctx, "copy-1", entity.CopyStatusBorrowed).Return(nil)
		mockBorrowRepo.On("Create", ctx, mock.AnythingOfType("*entity.BorrowRecord")).Return(nil)

		record, err := uc.CheckoutBook(ctx, "user-1", "book-1")

		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, "user-1", record.UserID)
		assert.Equal(t, "copy-1", record.BookCopyID)
		assert.Equal(t, entity.BorrowStatusActive, record.Status)
		mockUserRepo.AssertExpectations(t)
		mockBorrowRepo.AssertExpectations(t)
		mockCopyRepo.AssertExpectations(t)
		mockBookRepo.AssertExpectations(t)
	})

	t.Run("error - user suspended", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		user := &entity.User{
			ID:     "user-1",
			Status: entity.StatusSuspended,
		}

		mockUserRepo.On("GetByID", ctx, "user-1").Return(user, nil)

		record, err := uc.CheckoutBook(ctx, "user-1", "book-1")

		assert.Error(t, err)
		assert.Nil(t, record)
		assert.Contains(t, err.Error(), "not allowed to borrow")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error - borrowing limit reached", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		user := &entity.User{
			ID:             "user-1",
			Status:         entity.StatusActive,
			Role:           entity.RoleMember,
			BorrowingLimit: 5,
		}

		mockUserRepo.On("GetByID", ctx, "user-1").Return(user, nil)
		mockBorrowRepo.On("CountActiveByUserID", ctx, "user-1").Return(int64(5), nil)

		record, err := uc.CheckoutBook(ctx, "user-1", "book-1")

		assert.Error(t, err)
		assert.Nil(t, record)
		assert.Contains(t, err.Error(), "borrowing limit")
		mockUserRepo.AssertExpectations(t)
		mockBorrowRepo.AssertExpectations(t)
	})

	t.Run("error - has overdue books", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		user := &entity.User{
			ID:             "user-1",
			Status:         entity.StatusActive,
			Role:           entity.RoleMember,
			BorrowingLimit: 5,
		}

		overdueRecord := &entity.BorrowRecord{
			ID:     "record-1",
			Status: entity.BorrowStatusOverdue,
		}

		mockUserRepo.On("GetByID", ctx, "user-1").Return(user, nil)
		mockBorrowRepo.On("CountActiveByUserID", ctx, "user-1").Return(int64(0), nil)
		mockBorrowRepo.On("GetOverdueByUserID", ctx, "user-1").Return([]*entity.BorrowRecord{overdueRecord}, nil)

		record, err := uc.CheckoutBook(ctx, "user-1", "book-1")

		assert.Error(t, err)
		assert.Nil(t, record)
		assert.Contains(t, err.Error(), "overdue books")
		mockUserRepo.AssertExpectations(t)
		mockBorrowRepo.AssertExpectations(t)
	})

	t.Run("error - no available copies", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		user := &entity.User{
			ID:             "user-1",
			Status:         entity.StatusActive,
			Role:           entity.RoleMember,
			BorrowingLimit: 5,
		}

		book := &entity.Book{
			ID:        "book-1",
			DeletedAt: nil,
		}

		mockUserRepo.On("GetByID", ctx, "user-1").Return(user, nil)
		mockBorrowRepo.On("CountActiveByUserID", ctx, "user-1").Return(int64(0), nil)
		mockBorrowRepo.On("GetOverdueByUserID", ctx, "user-1").Return([]*entity.BorrowRecord{}, nil)
		mockBorrowRepo.On("CalculateTotalLateFees", ctx, "user-1").Return(0.0, nil)
		mockBookRepo.On("GetByID", ctx, "book-1").Return(book, nil)
		mockCopyRepo.On("GetAvailableCopies", ctx, "book-1").Return([]*entity.BookCopy{}, nil)

		record, err := uc.CheckoutBook(ctx, "user-1", "book-1")

		assert.Error(t, err)
		assert.Nil(t, record)
		assert.Contains(t, err.Error(), "no available copies")
		mockUserRepo.AssertExpectations(t)
		mockBorrowRepo.AssertExpectations(t)
		mockCopyRepo.AssertExpectations(t)
		mockBookRepo.AssertExpectations(t)
	})
}

// Test ReturnBook
func TestBorrowingUseCase_ReturnBook(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		borrowRecord := &entity.BorrowRecord{
			ID:           "record-1",
			BookCopyID:   "copy-1",
			Status:       entity.BorrowStatusActive,
			CheckoutDate: now.Add(-7 * 24 * time.Hour),
			DueDate:      now.Add(7 * 24 * time.Hour),
		}

		bookCopy := &entity.BookCopy{
			ID:     "copy-1",
			Status: entity.CopyStatusBorrowed,
		}

		mockBorrowRepo.On("GetByID", ctx, "record-1").Return(borrowRecord, nil)
		mockCopyRepo.On("GetByID", ctx, "copy-1").Return(bookCopy, nil)
		mockBorrowRepo.On("Update", ctx, borrowRecord).Return(nil)
		mockCopyRepo.On("UpdateStatus", ctx, "copy-1", entity.CopyStatusAvailable).Return(nil)

		err := uc.ReturnBook(ctx, "record-1")

		assert.NoError(t, err)
		mockBorrowRepo.AssertExpectations(t)
		mockCopyRepo.AssertExpectations(t)
	})

	t.Run("error - already returned", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		returnDate := now
		borrowRecord := &entity.BorrowRecord{
			ID:         "record-1",
			Status:     entity.BorrowStatusReturned,
			ReturnDate: &returnDate,
		}

		mockBorrowRepo.On("GetByID", ctx, "record-1").Return(borrowRecord, nil)

		err := uc.ReturnBook(ctx, "record-1")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already been returned")
		mockBorrowRepo.AssertExpectations(t)
	})

	t.Run("error - empty record ID", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		err := uc.ReturnBook(ctx, "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required")
	})
}

// Test RenewBook
func TestBorrowingUseCase_RenewBook(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		borrowRecord := &entity.BorrowRecord{
			ID:           "record-1",
			Status:       entity.BorrowStatusActive,
			RenewalCount: 0,
			CheckoutDate: now.Add(-7 * 24 * time.Hour),
			DueDate:      now.Add(7 * 24 * time.Hour),
		}

		mockBorrowRepo.On("GetByID", ctx, "record-1").Return(borrowRecord, nil)
		mockBorrowRepo.On("Update", ctx, borrowRecord).Return(nil)

		err := uc.RenewBook(ctx, "record-1")

		assert.NoError(t, err)
		mockBorrowRepo.AssertExpectations(t)
	})

	t.Run("error - overdue", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		borrowRecord := &entity.BorrowRecord{
			ID:           "record-1",
			Status:       entity.BorrowStatusActive,
			RenewalCount: 0,
			CheckoutDate: now.Add(-20 * 24 * time.Hour),
			DueDate:      now.Add(-5 * 24 * time.Hour),
		}

		mockBorrowRepo.On("GetByID", ctx, "record-1").Return(borrowRecord, nil)

		err := uc.RenewBook(ctx, "record-1")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "overdue")
		mockBorrowRepo.AssertExpectations(t)
	})

	t.Run("error - max renewal reached", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		borrowRecord := &entity.BorrowRecord{
			ID:           "record-1",
			Status:       entity.BorrowStatusActive,
			RenewalCount: entity.MaxRenewalCount,
			CheckoutDate: now.Add(-7 * 24 * time.Hour),
			DueDate:      now.Add(7 * 24 * time.Hour),
		}

		mockBorrowRepo.On("GetByID", ctx, "record-1").Return(borrowRecord, nil)

		err := uc.RenewBook(ctx, "record-1")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "renewal limit")
		mockBorrowRepo.AssertExpectations(t)
	})
}

// Test GetBorrowRecord
func TestBorrowingUseCase_GetBorrowRecord(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		expectedRecord := &entity.BorrowRecord{
			ID:     "record-1",
			UserID: "user-1",
		}

		mockBorrowRepo.On("GetByID", ctx, "record-1").Return(expectedRecord, nil)

		record, err := uc.GetBorrowRecord(ctx, "record-1")

		assert.NoError(t, err)
		assert.Equal(t, expectedRecord, record)
		mockBorrowRepo.AssertExpectations(t)
	})

	t.Run("error - empty ID", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		record, err := uc.GetBorrowRecord(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, record)
		assert.Contains(t, err.Error(), "required")
	})
}

// Test CanUserBorrowMore
func TestBorrowingUseCase_CanUserBorrowMore(t *testing.T) {
	ctx := context.Background()

	t.Run("can borrow", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		user := &entity.User{
			ID:             "user-1",
			Status:         entity.StatusActive,
			Role:           entity.RoleMember,
			BorrowingLimit: 5,
		}

		mockUserRepo.On("GetByID", ctx, "user-1").Return(user, nil)
		mockBorrowRepo.On("CountActiveByUserID", ctx, "user-1").Return(int64(2), nil)
		mockBorrowRepo.On("GetOverdueByUserID", ctx, "user-1").Return([]*entity.BorrowRecord{}, nil)
		mockBorrowRepo.On("CalculateTotalLateFees", ctx, "user-1").Return(0.0, nil)

		canBorrow, err := uc.CanUserBorrowMore(ctx, "user-1")

		assert.NoError(t, err)
		assert.True(t, canBorrow)
		mockUserRepo.AssertExpectations(t)
		mockBorrowRepo.AssertExpectations(t)
	})

	t.Run("cannot borrow - suspended", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		user := &entity.User{
			ID:     "user-1",
			Status: entity.StatusSuspended,
		}

		mockUserRepo.On("GetByID", ctx, "user-1").Return(user, nil)

		canBorrow, err := uc.CanUserBorrowMore(ctx, "user-1")

		assert.NoError(t, err)
		assert.False(t, canBorrow)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("cannot borrow - has overdue", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		user := &entity.User{
			ID:             "user-1",
			Status:         entity.StatusActive,
			Role:           entity.RoleMember,
			BorrowingLimit: 5,
		}

		overdueRecord := &entity.BorrowRecord{ID: "record-1"}

		mockUserRepo.On("GetByID", ctx, "user-1").Return(user, nil)
		mockBorrowRepo.On("CountActiveByUserID", ctx, "user-1").Return(int64(2), nil)
		mockBorrowRepo.On("GetOverdueByUserID", ctx, "user-1").Return([]*entity.BorrowRecord{overdueRecord}, nil)

		canBorrow, err := uc.CanUserBorrowMore(ctx, "user-1")

		assert.NoError(t, err)
		assert.False(t, canBorrow)
		mockUserRepo.AssertExpectations(t)
		mockBorrowRepo.AssertExpectations(t)
	})
}

// Test GetActiveBorrows
func TestBorrowingUseCase_GetActiveBorrows(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		user := &entity.User{ID: "user-1"}
		expectedRecords := []*entity.BorrowRecord{
			{ID: "record-1", Status: entity.BorrowStatusActive},
		}

		mockUserRepo.On("GetByID", ctx, "user-1").Return(user, nil)
		mockBorrowRepo.On("GetActiveByUserID", ctx, "user-1").Return(expectedRecords, nil)

		records, err := uc.GetActiveBorrows(ctx, "user-1")

		assert.NoError(t, err)
		assert.Equal(t, expectedRecords, records)
		mockUserRepo.AssertExpectations(t)
		mockBorrowRepo.AssertExpectations(t)
	})

	t.Run("error - empty user ID", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		records, err := uc.GetActiveBorrows(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, records)
		assert.Contains(t, err.Error(), "required")
	})
}

// Test CalculateUserLateFees
func TestBorrowingUseCase_CalculateUserLateFees(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		user := &entity.User{ID: "user-1"}

		mockUserRepo.On("GetByID", ctx, "user-1").Return(user, nil)
		mockBorrowRepo.On("CalculateTotalLateFees", ctx, "user-1").Return(5.50, nil)

		fees, err := uc.CalculateUserLateFees(ctx, "user-1")

		assert.NoError(t, err)
		assert.Equal(t, 5.50, fees)
		mockUserRepo.AssertExpectations(t)
		mockBorrowRepo.AssertExpectations(t)
	})

	t.Run("error - empty user ID", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		fees, err := uc.CalculateUserLateFees(ctx, "")

		assert.Error(t, err)
		assert.Equal(t, 0.0, fees)
		assert.Contains(t, err.Error(), "required")
	})
}

// Test GetMostBorrowedBooks
func TestBorrowingUseCase_GetMostBorrowedBooks(t *testing.T) {
	ctx := context.Background()
	startDate := time.Now().Add(-30 * 24 * time.Hour)
	endDate := time.Now()

	t.Run("success", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		expectedStats := []repository.BorrowStatistic{
			{BookID: "book-1", BorrowCount: 10},
			{BookID: "book-2", BorrowCount: 5},
		}

		mockBorrowRepo.On("GetMostBorrowedBooks", ctx, 10, startDate, endDate).Return(expectedStats, nil)

		stats, err := uc.GetMostBorrowedBooks(ctx, 10, startDate, endDate)

		assert.NoError(t, err)
		assert.Equal(t, expectedStats, stats)
		mockBorrowRepo.AssertExpectations(t)
	})

	t.Run("default limit", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		mockBorrowRepo.On("GetMostBorrowedBooks", ctx, 10, startDate, endDate).Return([]repository.BorrowStatistic{}, nil)

		stats, err := uc.GetMostBorrowedBooks(ctx, 0, startDate, endDate)

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		mockBorrowRepo.AssertExpectations(t)
	})
}

// Test UpdateOverdueStatuses
func TestBorrowingUseCase_UpdateOverdueStatuses(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		activeRecords := []*entity.BorrowRecord{
			{
				ID:      "record-1",
				Status:  entity.BorrowStatusActive,
				DueDate: now.Add(-5 * 24 * time.Hour), // Overdue
			},
			{
				ID:      "record-2",
				Status:  entity.BorrowStatusActive,
				DueDate: now.Add(5 * 24 * time.Hour), // Not overdue
			},
		}

		params := repository.BorrowRecordListParams{
			Page:     1,
			PageSize: 1000,
		}

		mockBorrowRepo.On("GetByStatus", ctx, entity.BorrowStatusActive, params).Return(activeRecords, int64(2), nil)
		mockBorrowRepo.On("UpdateStatus", ctx, "record-1", entity.BorrowStatusOverdue).Return(nil)

		err := uc.UpdateOverdueStatuses(ctx)

		assert.NoError(t, err)
		mockBorrowRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockBorrowRepo := new(MockBorrowRecordRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		uc := NewBorrowingUseCase(mockBorrowRepo, mockCopyRepo, mockUserRepo, mockBookRepo)

		params := repository.BorrowRecordListParams{
			Page:     1,
			PageSize: 1000,
		}

		mockBorrowRepo.On("GetByStatus", ctx, entity.BorrowStatusActive, params).Return(nil, int64(0), errors.New("database error"))

		err := uc.UpdateOverdueStatuses(ctx)

		require.Error(t, err)
		mockBorrowRepo.AssertExpectations(t)
	})
}