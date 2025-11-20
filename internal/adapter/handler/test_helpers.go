package handler

import (
	"context"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/repository"
)

// SetUserIDInContext is a helper function for tests to set user ID in context
func SetUserIDInContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// MockUserUseCaseForHandler is a mock implementation of UserUseCase for handler tests
type MockUserUseCaseForHandler struct {
	GetUserByIDFunc          func(ctx context.Context, id string) (*entity.User, error)
	ListUsersFunc            func(ctx context.Context, params repository.UserListParams) ([]*entity.User, int64, error)
	UpdateUserFunc           func(ctx context.Context, user *entity.User) error
	DeleteUserFunc           func(ctx context.Context, id string) error
	SuspendUserFunc          func(ctx context.Context, id string) error
	ActivateUserFunc         func(ctx context.Context, id string) error
	UpdateBorrowingLimitFunc func(ctx context.Context, id string, limit int) error
}

func (m *MockUserUseCaseForHandler) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockUserUseCaseForHandler) ListUsers(ctx context.Context, params repository.UserListParams) ([]*entity.User, int64, error) {
	if m.ListUsersFunc != nil {
		return m.ListUsersFunc(ctx, params)
	}
	return nil, 0, nil
}

func (m *MockUserUseCaseForHandler) UpdateUser(ctx context.Context, user *entity.User) error {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(ctx, user)
	}
	return nil
}

func (m *MockUserUseCaseForHandler) DeleteUser(ctx context.Context, id string) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(ctx, id)
	}
	return nil
}

func (m *MockUserUseCaseForHandler) SuspendUser(ctx context.Context, id string) error {
	if m.SuspendUserFunc != nil {
		return m.SuspendUserFunc(ctx, id)
	}
	return nil
}

func (m *MockUserUseCaseForHandler) ActivateUser(ctx context.Context, id string) error {
	if m.ActivateUserFunc != nil {
		return m.ActivateUserFunc(ctx, id)
	}
	return nil
}

func (m *MockUserUseCaseForHandler) UpdateBorrowingLimit(ctx context.Context, id string, limit int) error {
	if m.UpdateBorrowingLimitFunc != nil {
		return m.UpdateBorrowingLimitFunc(ctx, id, limit)
	}
	return nil
}

func (m *MockUserUseCaseForHandler) RegisterUser(ctx context.Context, user *entity.User) error {
	return nil
}

func (m *MockUserUseCaseForHandler) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, nil
}

func (m *MockUserUseCaseForHandler) GetUsersByRole(ctx context.Context, role entity.UserRole, params repository.UserListParams) ([]*entity.User, int64, error) {
	return nil, 0, nil
}

func (m *MockUserUseCaseForHandler) GetUsersByStatus(ctx context.Context, status entity.UserStatus, params repository.UserListParams) ([]*entity.User, int64, error) {
	return nil, 0, nil
}

func (m *MockUserUseCaseForHandler) SearchUsers(ctx context.Context, query string, params repository.UserListParams) ([]*entity.User, int64, error) {
	return nil, 0, nil
}

func (m *MockUserUseCaseForHandler) ExpireUser(ctx context.Context, id string) error {
	return nil
}

func (m *MockUserUseCaseForHandler) GetUserCount(ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *MockUserUseCaseForHandler) GetUserCountByRole(ctx context.Context, role entity.UserRole) (int64, error) {
	return 0, nil
}

func (m *MockUserUseCaseForHandler) GetActiveMembers(ctx context.Context, params repository.UserListParams) ([]*entity.User, int64, error) {
	return nil, 0, nil
}

// MockBorrowingUseCase is a mock implementation of BorrowingUseCase
type MockBorrowingUseCase struct {
	CheckoutBookFunc            func(ctx context.Context, userID, bookID string) (*entity.BorrowRecord, error)
	ReturnBookFunc              func(ctx context.Context, borrowRecordID string) error
	RenewBookFunc               func(ctx context.Context, borrowRecordID string) error
	GetUserBorrowingHistoryFunc func(ctx context.Context, userID string, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error)
	GetActiveBorrowsFunc        func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error)
	GetOverdueBorrowsFunc       func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error)
	GetBorrowRecordFunc         func(ctx context.Context, id string) (*entity.BorrowRecord, error)
}

func (m *MockBorrowingUseCase) CheckoutBook(ctx context.Context, userID, bookID string) (*entity.BorrowRecord, error) {
	if m.CheckoutBookFunc != nil {
		return m.CheckoutBookFunc(ctx, userID, bookID)
	}
	return nil, nil
}

func (m *MockBorrowingUseCase) ReturnBook(ctx context.Context, borrowRecordID string) error {
	if m.ReturnBookFunc != nil {
		return m.ReturnBookFunc(ctx, borrowRecordID)
	}
	return nil
}

func (m *MockBorrowingUseCase) RenewBook(ctx context.Context, borrowRecordID string) error {
	if m.RenewBookFunc != nil {
		return m.RenewBookFunc(ctx, borrowRecordID)
	}
	return nil
}

func (m *MockBorrowingUseCase) GetUserBorrowingHistory(ctx context.Context, userID string, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	if m.GetUserBorrowingHistoryFunc != nil {
		return m.GetUserBorrowingHistoryFunc(ctx, userID, params)
	}
	return nil, 0, nil
}

func (m *MockBorrowingUseCase) GetActiveBorrows(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
	if m.GetActiveBorrowsFunc != nil {
		return m.GetActiveBorrowsFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockBorrowingUseCase) GetOverdueBorrows(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
	if m.GetOverdueBorrowsFunc != nil {
		return m.GetOverdueBorrowsFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockBorrowingUseCase) GetBorrowRecord(ctx context.Context, id string) (*entity.BorrowRecord, error) {
	if m.GetBorrowRecordFunc != nil {
		return m.GetBorrowRecordFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockBorrowingUseCase) CanUserBorrowMore(ctx context.Context, userID string) (bool, error) {
	return true, nil
}

func (m *MockBorrowingUseCase) GetAllOverdueRecords(ctx context.Context, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	return nil, 0, nil
}

func (m *MockBorrowingUseCase) GetDueSoonRecords(ctx context.Context, days int, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
	return nil, 0, nil
}

func (m *MockBorrowingUseCase) GetUserBorrowingStatistics(ctx context.Context, userID string) (*repository.UserBorrowStatistics, error) {
	return nil, nil
}

func (m *MockBorrowingUseCase) CalculateUserLateFees(ctx context.Context, userID string) (float64, error) {
	return 0, nil
}

func (m *MockBorrowingUseCase) GetMostBorrowedBooks(ctx context.Context, limit int, startDate, endDate time.Time) ([]repository.BorrowStatistic, error) {
	return nil, nil
}

func (m *MockBorrowingUseCase) UpdateOverdueStatuses(ctx context.Context) error {
	return nil
}