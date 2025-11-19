package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/repository"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/usecase/mocks"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/usecase/testdata"
)

func TestUserUseCase_RegisterUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		user          *entity.User
		setupMocks    func(*mocks.MockUserRepository)
		wantErr       bool
		errorContains string
	}{
		{
			name: "successful user registration",
			user: &entity.User{
				Email:    "newuser@example.com",
				FullName: "New User",
				Role:     entity.RoleMember,
				Status:   entity.StatusActive,
			},
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.ExistsByEmailFunc = func(ctx context.Context, email string) (bool, error) {
					return false, nil
				}
				ur.CreateFunc = func(ctx context.Context, user *entity.User) error {
					user.ID = "new-user-id"
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "validation error - missing email",
			user: &entity.User{
				FullName: "New User",
				Role:     entity.RoleMember,
			},
			setupMocks:    func(ur *mocks.MockUserRepository) {},
			wantErr:       true,
			errorContains: "validation failed",
		},
		{
			name: "duplicate email error",
			user: &entity.User{
				Email:    "existing@example.com",
				FullName: "New User",
				Role:     entity.RoleMember,
				Status:   entity.StatusActive,
			},
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.ExistsByEmailFunc = func(ctx context.Context, email string) (bool, error) {
					return true, nil
				}
			},
			wantErr:       true,
			errorContains: "email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mocks.MockUserRepository{}
			borrowRecordRepo := &mocks.MockBorrowRecordRepository{}
			tt.setupMocks(userRepo)

			uc := NewUserUseCase(userRepo, borrowRecordRepo)
			err := uc.RegisterUser(ctx, tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("RegisterUser() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr {
				if tt.user.Status != entity.StatusActive {
					t.Errorf("Status should be active, got %v", tt.user.Status)
				}
				if tt.user.BorrowingLimit == 0 {
					t.Error("BorrowingLimit should be set to default")
				}
			}
		})
	}
}

func TestUserUseCase_GetUserByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        string
		setupMocks    func(*mocks.MockUserRepository)
		wantErr       bool
		errorContains string
	}{
		{
			name:   "successful retrieval",
			userID: "user-123",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
			},
			wantErr: false,
		},
		{
			name:          "empty user ID",
			userID:        "",
			setupMocks:    func(ur *mocks.MockUserRepository) {},
			wantErr:       true,
			errorContains: "user ID is required",
		},
		{
			name:   "user not found",
			userID: "nonexistent",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr:       true,
			errorContains: "failed to get user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mocks.MockUserRepository{}
			borrowRecordRepo := &mocks.MockBorrowRecordRepository{}
			tt.setupMocks(userRepo)

			uc := NewUserUseCase(userRepo, borrowRecordRepo)
			user, err := uc.GetUserByID(ctx, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("GetUserByID() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr && user == nil {
				t.Error("GetUserByID() should return a user")
			}
		})
	}
}

func TestUserUseCase_DeleteUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        string
		setupMocks    func(*mocks.MockUserRepository, *mocks.MockBorrowRecordRepository)
		wantErr       bool
		errorContains string
	}{
		{
			name:   "successful deletion",
			userID: "user-123",
			setupMocks: func(ur *mocks.MockUserRepository, brr *mocks.MockBorrowRecordRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
				brr.GetActiveByUserIDFunc = func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
					return []*entity.BorrowRecord{}, nil
				}
				brr.GetOverdueByUserIDFunc = func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
					return []*entity.BorrowRecord{}, nil
				}
				brr.CalculateTotalLateFeesFunc = func(ctx context.Context, userID string) (float64, error) {
					return 0, nil
				}
				ur.DeleteFunc = func(ctx context.Context, id string) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:          "empty user ID",
			userID:        "",
			setupMocks:    func(ur *mocks.MockUserRepository, brr *mocks.MockBorrowRecordRepository) {},
			wantErr:       true,
			errorContains: "user ID is required",
		},
		{
			name:   "cannot delete user with active borrows",
			userID: "user-123",
			setupMocks: func(ur *mocks.MockUserRepository, brr *mocks.MockBorrowRecordRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
				brr.GetActiveByUserIDFunc = func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
					return []*entity.BorrowRecord{
						testdata.CreateTestBorrowRecord("rec-1", userID, "copy-1"),
					}, nil
				}
			},
			wantErr:       true,
			errorContains: "cannot delete user with active borrows",
		},
		{
			name:   "cannot delete user with unpaid late fees",
			userID: "user-123",
			setupMocks: func(ur *mocks.MockUserRepository, brr *mocks.MockBorrowRecordRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
				brr.GetActiveByUserIDFunc = func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
					return []*entity.BorrowRecord{}, nil
				}
				brr.GetOverdueByUserIDFunc = func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
					return []*entity.BorrowRecord{}, nil
				}
				brr.CalculateTotalLateFeesFunc = func(ctx context.Context, userID string) (float64, error) {
					return 5.50, nil
				}
			},
			wantErr:       true,
			errorContains: "cannot delete user with unpaid late fees",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mocks.MockUserRepository{}
			borrowRecordRepo := &mocks.MockBorrowRecordRepository{}
			tt.setupMocks(userRepo, borrowRecordRepo)

			uc := NewUserUseCase(userRepo, borrowRecordRepo)
			err := uc.DeleteUser(ctx, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("DeleteUser() error = %v, should contain %v", err, tt.errorContains)
				}
			}
		})
	}
}

func TestUserUseCase_SuspendUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        string
		setupMocks    func(*mocks.MockUserRepository)
		wantErr       bool
		errorContains string
	}{
		{
			name:   "successful suspension",
			userID: "user-123",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
				ur.UpdateStatusFunc = func(ctx context.Context, id string, status entity.UserStatus) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:   "already suspended user",
			userID: "user-123",
			setupMocks: func(ur *mocks.MockUserRepository) {
				user := testdata.CreateTestUser("user-123", entity.RoleMember)
				user.Status = entity.StatusSuspended
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return user, nil
				}
			},
			wantErr:       true,
			errorContains: "failed to suspend user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mocks.MockUserRepository{}
			borrowRecordRepo := &mocks.MockBorrowRecordRepository{}
			tt.setupMocks(userRepo)

			uc := NewUserUseCase(userRepo, borrowRecordRepo)
			err := uc.SuspendUser(ctx, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("SuspendUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("SuspendUser() error = %v, should contain %v", err, tt.errorContains)
				}
			}
		})
	}
}

func TestUserUseCase_UpdateBorrowingLimit(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        string
		newLimit      int
		setupMocks    func(*mocks.MockUserRepository)
		wantErr       bool
		errorContains string
	}{
		{
			name:     "successful limit update",
			userID:   "user-123",
			newLimit: 10,
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
				ur.UpdateBorrowingLimitFunc = func(ctx context.Context, id string, limit int) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:          "negative limit",
			userID:        "user-123",
			newLimit:      -1,
			setupMocks:    func(ur *mocks.MockUserRepository) {},
			wantErr:       true,
			errorContains: "borrowing limit cannot be negative",
		},
		{
			name:   "user not found",
			userID: "nonexistent",
			newLimit: 10,
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr:       true,
			errorContains: "failed to get user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mocks.MockUserRepository{}
			borrowRecordRepo := &mocks.MockBorrowRecordRepository{}
			tt.setupMocks(userRepo)

			uc := NewUserUseCase(userRepo, borrowRecordRepo)
			err := uc.UpdateBorrowingLimit(ctx, tt.userID, tt.newLimit)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateBorrowingLimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("UpdateBorrowingLimit() error = %v, should contain %v", err, tt.errorContains)
				}
			}
		})
	}
}

func TestUserUseCase_ListUsers(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		params     repository.UserListParams
		setupMocks func(*mocks.MockUserRepository)
		wantCount  int
		wantTotal  int64
		wantErr    bool
	}{
		{
			name: "successful list with default params",
			params: repository.UserListParams{
				Page:     0,
				PageSize: 0,
			},
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.ListFunc = func(ctx context.Context, params repository.UserListParams) ([]*entity.User, int64, error) {
					users := []*entity.User{
						testdata.CreateTestUser("user-1", entity.RoleMember),
						testdata.CreateTestUser("user-2", entity.RoleMember),
					}
					return users, 2, nil
				}
			},
			wantCount: 2,
			wantTotal: 2,
			wantErr:   false,
		},
		{
			name: "repository error",
			params: repository.UserListParams{
				Page:     1,
				PageSize: 20,
			},
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.ListFunc = func(ctx context.Context, params repository.UserListParams) ([]*entity.User, int64, error) {
					return nil, 0, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mocks.MockUserRepository{}
			borrowRecordRepo := &mocks.MockBorrowRecordRepository{}
			tt.setupMocks(userRepo)

			uc := NewUserUseCase(userRepo, borrowRecordRepo)
			users, total, err := uc.ListUsers(ctx, tt.params)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(users) != tt.wantCount {
					t.Errorf("ListUsers() returned %d users, want %d", len(users), tt.wantCount)
				}
				if total != tt.wantTotal {
					t.Errorf("ListUsers() total = %d, want %d", total, tt.wantTotal)
				}
			}
		})
	}
}

func TestUserUseCase_SearchUsers(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		query         string
		params        repository.UserListParams
		setupMocks    func(*mocks.MockUserRepository)
		wantCount     int
		wantErr       bool
		errorContains string
	}{
		{
			name:  "successful search",
			query: "test",
			params: repository.UserListParams{
				Page:     1,
				PageSize: 20,
			},
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.SearchFunc = func(ctx context.Context, query string, params repository.UserListParams) ([]*entity.User, int64, error) {
					users := []*entity.User{
						testdata.CreateTestUser("user-1", entity.RoleMember),
					}
					return users, 1, nil
				}
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:  "empty query",
			query: "",
			params: repository.UserListParams{
				Page:     1,
				PageSize: 20,
			},
			setupMocks:    func(ur *mocks.MockUserRepository) {},
			wantErr:       true,
			errorContains: "search query is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mocks.MockUserRepository{}
			borrowRecordRepo := &mocks.MockBorrowRecordRepository{}
			tt.setupMocks(userRepo)

			uc := NewUserUseCase(userRepo, borrowRecordRepo)
			users, _, err := uc.SearchUsers(ctx, tt.query, tt.params)

			if (err != nil) != tt.wantErr {
				t.Errorf("SearchUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("SearchUsers() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr && len(users) != tt.wantCount {
				t.Errorf("SearchUsers() returned %d users, want %d", len(users), tt.wantCount)
			}
		})
	}
}