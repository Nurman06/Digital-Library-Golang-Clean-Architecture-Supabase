package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/usecase/mocks"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/usecase/testdata"
)

func TestBorrowingUseCase_CheckoutBook(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        string
		bookID        string
		setupMocks    func(*mocks.MockBorrowRecordRepository, *mocks.MockBookCopyRepository, *mocks.MockUserRepository, *mocks.MockBookRepository)
		wantErr       bool
		errorContains string
	}{
		{
			name:   "successful checkout",
			userID: "user-123",
			bookID: "book-123",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository, bcr *mocks.MockBookCopyRepository, ur *mocks.MockUserRepository, br *mocks.MockBookRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
				brr.CountActiveByUserIDFunc = func(ctx context.Context, userID string) (int64, error) {
					return 2, nil
				}
				brr.GetOverdueByUserIDFunc = func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
					return []*entity.BorrowRecord{}, nil
				}
				brr.CalculateTotalLateFeesFunc = func(ctx context.Context, userID string) (float64, error) {
					return 0, nil
				}
				br.GetByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return testdata.CreateTestBook(id), nil
				}
				bcr.GetAvailableCopiesFunc = func(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
					return []*entity.BookCopy{testdata.CreateTestBookCopy("copy-1", bookID)}, nil
				}
				bcr.UpdateStatusFunc = func(ctx context.Context, id string, status entity.CopyStatus) error {
					return nil
				}
				brr.CreateFunc = func(ctx context.Context, record *entity.BorrowRecord) error {
					record.ID = "new-record-id"
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:   "empty user ID",
			userID: "",
			bookID: "book-123",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository, bcr *mocks.MockBookCopyRepository, ur *mocks.MockUserRepository, br *mocks.MockBookRepository) {
			},
			wantErr:       true,
			errorContains: "user ID is required",
		},
		{
			name:   "empty book ID",
			userID: "user-123",
			bookID: "",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository, bcr *mocks.MockBookCopyRepository, ur *mocks.MockUserRepository, br *mocks.MockBookRepository) {
			},
			wantErr:       true,
			errorContains: "book ID is required",
		},
		{
			name:   "user suspended",
			userID: "user-123",
			bookID: "book-123",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository, bcr *mocks.MockBookCopyRepository, ur *mocks.MockUserRepository, br *mocks.MockBookRepository) {
				user := testdata.CreateTestUser("user-123", entity.RoleMember)
				user.Status = entity.StatusSuspended
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return user, nil
				}
			},
			wantErr:       true,
			errorContains: "not allowed to borrow",
		},
		{
			name:   "borrowing limit reached",
			userID: "user-123",
			bookID: "book-123",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository, bcr *mocks.MockBookCopyRepository, ur *mocks.MockUserRepository, br *mocks.MockBookRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
				brr.CountActiveByUserIDFunc = func(ctx context.Context, userID string) (int64, error) {
					return 5, nil // Member limit is 5
				}
			},
			wantErr:       true,
			errorContains: "reached borrowing limit",
		},
		{
			name:   "user has overdue books",
			userID: "user-123",
			bookID: "book-123",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository, bcr *mocks.MockBookCopyRepository, ur *mocks.MockUserRepository, br *mocks.MockBookRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
				brr.CountActiveByUserIDFunc = func(ctx context.Context, userID string) (int64, error) {
					return 2, nil
				}
				brr.GetOverdueByUserIDFunc = func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
					return []*entity.BorrowRecord{
						testdata.CreateTestBorrowRecord("rec-1", userID, "copy-1"),
					}, nil
				}
			},
			wantErr:       true,
			errorContains: "overdue books",
		},
		{
			name:   "user has unpaid late fees",
			userID: "user-123",
			bookID: "book-123",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository, bcr *mocks.MockBookCopyRepository, ur *mocks.MockUserRepository, br *mocks.MockBookRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
				brr.CountActiveByUserIDFunc = func(ctx context.Context, userID string) (int64, error) {
					return 2, nil
				}
				brr.GetOverdueByUserIDFunc = func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
					return []*entity.BorrowRecord{}, nil
				}
				brr.CalculateTotalLateFeesFunc = func(ctx context.Context, userID string) (float64, error) {
					return 5.50, nil
				}
			},
			wantErr:       true,
			errorContains: "unpaid late fees",
		},
		{
			name:   "book is deleted",
			userID: "user-123",
			bookID: "book-123",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository, bcr *mocks.MockBookCopyRepository, ur *mocks.MockUserRepository, br *mocks.MockBookRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
				brr.CountActiveByUserIDFunc = func(ctx context.Context, userID string) (int64, error) {
					return 2, nil
				}
				brr.GetOverdueByUserIDFunc = func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
					return []*entity.BorrowRecord{}, nil
				}
				brr.CalculateTotalLateFeesFunc = func(ctx context.Context, userID string) (float64, error) {
					return 0, nil
				}
				br.GetByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					book := testdata.CreateTestBook(id)
					now := time.Now()
					book.DeletedAt = &now
					return book, nil
				}
			},
			wantErr:       true,
			errorContains: "not available",
		},
		{
			name:   "no available copies",
			userID: "user-123",
			bookID: "book-123",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository, bcr *mocks.MockBookCopyRepository, ur *mocks.MockUserRepository, br *mocks.MockBookRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
				brr.CountActiveByUserIDFunc = func(ctx context.Context, userID string) (int64, error) {
					return 2, nil
				}
				brr.GetOverdueByUserIDFunc = func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
					return []*entity.BorrowRecord{}, nil
				}
				brr.CalculateTotalLateFeesFunc = func(ctx context.Context, userID string) (float64, error) {
					return 0, nil
				}
				br.GetByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return testdata.CreateTestBook(id), nil
				}
				bcr.GetAvailableCopiesFunc = func(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
					return []*entity.BookCopy{}, nil
				}
			},
			wantErr:       true,
			errorContains: "no available copies",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			borrowRecordRepo := &mocks.MockBorrowRecordRepository{}
			bookCopyRepo := &mocks.MockBookCopyRepository{}
			userRepo := &mocks.MockUserRepository{}
			bookRepo := &mocks.MockBookRepository{}
			tt.setupMocks(borrowRecordRepo, bookCopyRepo, userRepo, bookRepo)

			uc := NewBorrowingUseCase(borrowRecordRepo, bookCopyRepo, userRepo, bookRepo)
			record, err := uc.CheckoutBook(ctx, tt.userID, tt.bookID)

			if (err != nil) != tt.wantErr {
				t.Errorf("CheckoutBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("CheckoutBook() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr {
				if record == nil {
					t.Error("CheckoutBook() should return a borrow record")
				}
				if record.Status != entity.BorrowStatusActive {
					t.Errorf("BorrowRecord status = %v, want %v", record.Status, entity.BorrowStatusActive)
				}
				if record.DueDate.Before(time.Now()) {
					t.Error("DueDate should be in the future")
				}
			}
		})
	}
}

func TestBorrowingUseCase_ReturnBook(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name             string
		borrowRecordID   string
		setupMocks       func(*mocks.MockBorrowRecordRepository, *mocks.MockBookCopyRepository)
		wantErr          bool
		errorContains    string
	}{
		{
			name:           "successful return",
			borrowRecordID: "record-123",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository, bcr *mocks.MockBookCopyRepository) {
				brr.GetByIDFunc = func(ctx context.Context, id string) (*entity.BorrowRecord, error) {
					return testdata.CreateTestBorrowRecord(id, "user-1", "copy-1"), nil
				}
				bcr.GetByIDFunc = func(ctx context.Context, id string) (*entity.BookCopy, error) {
					copy := testdata.CreateTestBookCopy(id, "book-1")
					copy.Status = entity.CopyStatusBorrowed
					return copy, nil
				}
				brr.UpdateFunc = func(ctx context.Context, record *entity.BorrowRecord) error {
					return nil
				}
				bcr.UpdateStatusFunc = func(ctx context.Context, id string, status entity.CopyStatus) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:           "empty borrow record ID",
			borrowRecordID: "",
			setupMocks:     func(brr *mocks.MockBorrowRecordRepository, bcr *mocks.MockBookCopyRepository) {},
			wantErr:        true,
			errorContains:  "borrow record ID is required",
		},
		{
			name:           "borrow record not found",
			borrowRecordID: "nonexistent",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository, bcr *mocks.MockBookCopyRepository) {
				brr.GetByIDFunc = func(ctx context.Context, id string) (*entity.BorrowRecord, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr:       true,
			errorContains: "failed to get borrow record",
		},
		{
			name:           "already returned",
			borrowRecordID: "record-123",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository, bcr *mocks.MockBookCopyRepository) {
				record := testdata.CreateTestBorrowRecord("record-123", "user-1", "copy-1")
				now := time.Now()
				record.ReturnDate = &now
				record.Status = entity.BorrowStatusReturned
				brr.GetByIDFunc = func(ctx context.Context, id string) (*entity.BorrowRecord, error) {
					return record, nil
				}
			},
			wantErr:       true,
			errorContains: "already been returned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			borrowRecordRepo := &mocks.MockBorrowRecordRepository{}
			bookCopyRepo := &mocks.MockBookCopyRepository{}
			userRepo := &mocks.MockUserRepository{}
			bookRepo := &mocks.MockBookRepository{}
			tt.setupMocks(borrowRecordRepo, bookCopyRepo)

			uc := NewBorrowingUseCase(borrowRecordRepo, bookCopyRepo, userRepo, bookRepo)
			err := uc.ReturnBook(ctx, tt.borrowRecordID)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReturnBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("ReturnBook() error = %v, should contain %v", err, tt.errorContains)
				}
			}
		})
	}
}

func TestBorrowingUseCase_RenewBook(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name             string
		borrowRecordID   string
		setupMocks       func(*mocks.MockBorrowRecordRepository)
		wantErr          bool
		errorContains    string
	}{
		{
			name:           "successful renewal",
			borrowRecordID: "record-123",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository) {
				record := testdata.CreateTestBorrowRecord("record-123", "user-1", "copy-1")
				record.RenewalCount = 0
				brr.GetByIDFunc = func(ctx context.Context, id string) (*entity.BorrowRecord, error) {
					return record, nil
				}
				brr.UpdateFunc = func(ctx context.Context, record *entity.BorrowRecord) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:           "empty borrow record ID",
			borrowRecordID: "",
			setupMocks:     func(brr *mocks.MockBorrowRecordRepository) {},
			wantErr:        true,
			errorContains:  "borrow record ID is required",
		},
		{
			name:           "borrow record not found",
			borrowRecordID: "nonexistent",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository) {
				brr.GetByIDFunc = func(ctx context.Context, id string) (*entity.BorrowRecord, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr:       true,
			errorContains: "failed to get borrow record",
		},
		{
			name:           "renewal limit reached",
			borrowRecordID: "record-123",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository) {
				record := testdata.CreateTestBorrowRecord("record-123", "user-1", "copy-1")
				record.RenewalCount = entity.MaxRenewalCount
				brr.GetByIDFunc = func(ctx context.Context, id string) (*entity.BorrowRecord, error) {
					return record, nil
				}
			},
			wantErr:       true,
			errorContains: "maximum renewal limit",
		},
		{
			name:           "book is overdue",
			borrowRecordID: "record-123",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository) {
				record := testdata.CreateTestBorrowRecord("record-123", "user-1", "copy-1")
				record.DueDate = time.Now().Add(-24 * time.Hour) // Overdue
				record.RenewalCount = 0
				brr.GetByIDFunc = func(ctx context.Context, id string) (*entity.BorrowRecord, error) {
					return record, nil
				}
			},
			wantErr:       true,
			errorContains: "cannot renew overdue book",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			borrowRecordRepo := &mocks.MockBorrowRecordRepository{}
			bookCopyRepo := &mocks.MockBookCopyRepository{}
			userRepo := &mocks.MockUserRepository{}
			bookRepo := &mocks.MockBookRepository{}
			tt.setupMocks(borrowRecordRepo)

			uc := NewBorrowingUseCase(borrowRecordRepo, bookCopyRepo, userRepo, bookRepo)
			err := uc.RenewBook(ctx, tt.borrowRecordID)

			if (err != nil) != tt.wantErr {
				t.Errorf("RenewBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("RenewBook() error = %v, should contain %v", err, tt.errorContains)
				}
			}
		})
	}
}

func TestBorrowingUseCase_GetActiveBorrows(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        string
		setupMocks    func(*mocks.MockBorrowRecordRepository, *mocks.MockUserRepository)
		wantCount     int
		wantErr       bool
		errorContains string
	}{
		{
			name:   "successful retrieval",
			userID: "user-123",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository, ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
				brr.GetActiveByUserIDFunc = func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
					return []*entity.BorrowRecord{
						testdata.CreateTestBorrowRecord("rec-1", userID, "copy-1"),
						testdata.CreateTestBorrowRecord("rec-2", userID, "copy-2"),
					}, nil
				}
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:          "empty user ID",
			userID:        "",
			setupMocks:    func(brr *mocks.MockBorrowRecordRepository, ur *mocks.MockUserRepository) {},
			wantErr:       true,
			errorContains: "user ID is required",
		},
		{
			name:   "user not found",
			userID: "nonexistent",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository, ur *mocks.MockUserRepository) {
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
			borrowRecordRepo := &mocks.MockBorrowRecordRepository{}
			bookCopyRepo := &mocks.MockBookCopyRepository{}
			userRepo := &mocks.MockUserRepository{}
			bookRepo := &mocks.MockBookRepository{}
			tt.setupMocks(borrowRecordRepo, userRepo)

			uc := NewBorrowingUseCase(borrowRecordRepo, bookCopyRepo, userRepo, bookRepo)
			records, err := uc.GetActiveBorrows(ctx, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetActiveBorrows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("GetActiveBorrows() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr && len(records) != tt.wantCount {
				t.Errorf("GetActiveBorrows() returned %d records, want %d", len(records), tt.wantCount)
			}
		})
	}
}

func TestBorrowingUseCase_CanUserBorrowMore(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        string
		setupMocks    func(*mocks.MockBorrowRecordRepository, *mocks.MockUserRepository)
		want          bool
		wantErr       bool
		errorContains string
	}{
		{
			name:   "user can borrow more",
			userID: "user-123",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository, ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
				brr.CountActiveByUserIDFunc = func(ctx context.Context, userID string) (int64, error) {
					return 2, nil
				}
				brr.GetOverdueByUserIDFunc = func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
					return []*entity.BorrowRecord{}, nil
				}
				brr.CalculateTotalLateFeesFunc = func(ctx context.Context, userID string) (float64, error) {
					return 0, nil
				}
			},
			want:    true,
			wantErr: false,
		},
		{
			name:   "user at borrowing limit",
			userID: "user-123",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository, ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
				brr.CountActiveByUserIDFunc = func(ctx context.Context, userID string) (int64, error) {
					return 5, nil // Member limit is 5
				}
				brr.GetOverdueByUserIDFunc = func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
					return []*entity.BorrowRecord{}, nil
				}
				brr.CalculateTotalLateFeesFunc = func(ctx context.Context, userID string) (float64, error) {
					return 0, nil
				}
			},
			want:    false,
			wantErr: false,
		},
		{
			name:   "user has overdue books",
			userID: "user-123",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository, ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
				brr.CountActiveByUserIDFunc = func(ctx context.Context, userID string) (int64, error) {
					return 2, nil
				}
				brr.GetOverdueByUserIDFunc = func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
					return []*entity.BorrowRecord{
						testdata.CreateTestBorrowRecord("rec-1", userID, "copy-1"),
					}, nil
				}
			},
			want:    false,
			wantErr: false,
		},
		{
			name:   "admin with unlimited borrowing",
			userID: "admin-123",
			setupMocks: func(brr *mocks.MockBorrowRecordRepository, ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					user := testdata.CreateTestUser(id, entity.RoleAdmin)
					user.BorrowingLimit = 0 // Unlimited
					return user, nil
				}
				brr.CountActiveByUserIDFunc = func(ctx context.Context, userID string) (int64, error) {
					return 100, nil
				}
				brr.GetOverdueByUserIDFunc = func(ctx context.Context, userID string) ([]*entity.BorrowRecord, error) {
					return []*entity.BorrowRecord{}, nil
				}
				brr.CalculateTotalLateFeesFunc = func(ctx context.Context, userID string) (float64, error) {
					return 0, nil
				}
			},
			want:    true,
			wantErr: false,
		},
		{
			name:          "empty user ID",
			userID:        "",
			setupMocks:    func(brr *mocks.MockBorrowRecordRepository, ur *mocks.MockUserRepository) {},
			wantErr:       true,
			errorContains: "user ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			borrowRecordRepo := &mocks.MockBorrowRecordRepository{}
			bookCopyRepo := &mocks.MockBookCopyRepository{}
			userRepo := &mocks.MockUserRepository{}
			bookRepo := &mocks.MockBookRepository{}
			tt.setupMocks(borrowRecordRepo, userRepo)

			uc := NewBorrowingUseCase(borrowRecordRepo, bookCopyRepo, userRepo, bookRepo)
			canBorrow, err := uc.CanUserBorrowMore(ctx, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("CanUserBorrowMore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("CanUserBorrowMore() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr && canBorrow != tt.want {
				t.Errorf("CanUserBorrowMore() = %v, want %v", canBorrow, tt.want)
			}
		})
	}
}