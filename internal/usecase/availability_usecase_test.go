package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/usecase/mocks"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/usecase/testdata"
)

func TestAvailabilityUseCase_GetBookAvailability(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		bookID        string
		setupMocks    func(*mocks.MockBookRepository, *mocks.MockBookCopyRepository, *mocks.MockReservationRepository)
		wantErr       bool
		errorContains string
		wantAvailable int64
		wantTotal     int64
	}{
		{
			name:   "successful availability check",
			bookID: "book-123",
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository, rr *mocks.MockReservationRepository) {
				br.GetByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return testdata.CreateTestBook(id), nil
				}
				bcr.GetByBookIDFunc = func(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
					copies := []*entity.BookCopy{
						testdata.CreateTestBookCopy("copy-1", bookID),
						testdata.CreateTestBookCopy("copy-2", bookID),
						testdata.CreateTestBookCopy("copy-3", bookID),
					}
					copies[0].Status = entity.CopyStatusAvailable
					copies[1].Status = entity.CopyStatusAvailable
					copies[2].Status = entity.CopyStatusBorrowed
					return copies, nil
				}
				rr.CountPendingByBookIDFunc = func(ctx context.Context, bookID string) (int64, error) {
					return 0, nil
				}
			},
			wantErr:       false,
			wantAvailable: 2,
			wantTotal:     3,
		},
		{
			name:          "empty book ID",
			bookID:        "",
			setupMocks:    func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository, rr *mocks.MockReservationRepository) {},
			wantErr:       true,
			errorContains: "book ID is required",
		},
		{
			name:   "book not found",
			bookID: "nonexistent",
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository, rr *mocks.MockReservationRepository) {
				br.GetByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr:       true,
			errorContains: "failed to get book",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookRepo := &mocks.MockBookRepository{}
			bookCopyRepo := &mocks.MockBookCopyRepository{}
			borrowRecordRepo := &mocks.MockBorrowRecordRepository{}
			reservationRepo := &mocks.MockReservationRepository{}
			userRepo := &mocks.MockUserRepository{}
			tt.setupMocks(bookRepo, bookCopyRepo, reservationRepo)

			uc := NewAvailabilityUseCase(bookRepo, bookCopyRepo, borrowRecordRepo, reservationRepo, userRepo)
			availability, err := uc.GetBookAvailability(ctx, tt.bookID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetBookAvailability() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("GetBookAvailability() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr {
				if availability.AvailableCopies != tt.wantAvailable {
					t.Errorf("AvailableCopies = %d, want %d", availability.AvailableCopies, tt.wantAvailable)
				}
				if availability.TotalCopies != tt.wantTotal {
					t.Errorf("TotalCopies = %d, want %d", availability.TotalCopies, tt.wantTotal)
				}
			}
		})
	}
}

func TestAvailabilityUseCase_CheckBookAvailability(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		bookID        string
		setupMocks    func(*mocks.MockBookCopyRepository)
		want          bool
		wantErr       bool
		errorContains string
	}{
		{
			name:   "book is available",
			bookID: "book-123",
			setupMocks: func(bcr *mocks.MockBookCopyRepository) {
				bcr.CountAvailableByBookIDFunc = func(ctx context.Context, bookID string) (int64, error) {
					return 2, nil
				}
			},
			want:    true,
			wantErr: false,
		},
		{
			name:   "book is not available",
			bookID: "book-123",
			setupMocks: func(bcr *mocks.MockBookCopyRepository) {
				bcr.CountAvailableByBookIDFunc = func(ctx context.Context, bookID string) (int64, error) {
					return 0, nil
				}
			},
			want:    false,
			wantErr: false,
		},
		{
			name:          "empty book ID",
			bookID:        "",
			setupMocks:    func(bcr *mocks.MockBookCopyRepository) {},
			wantErr:       true,
			errorContains: "book ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookRepo := &mocks.MockBookRepository{}
			bookCopyRepo := &mocks.MockBookCopyRepository{}
			borrowRecordRepo := &mocks.MockBorrowRecordRepository{}
			reservationRepo := &mocks.MockReservationRepository{}
			userRepo := &mocks.MockUserRepository{}
			tt.setupMocks(bookCopyRepo)

			uc := NewAvailabilityUseCase(bookRepo, bookCopyRepo, borrowRecordRepo, reservationRepo, userRepo)
			available, err := uc.CheckBookAvailability(ctx, tt.bookID)

			if (err != nil) != tt.wantErr {
				t.Errorf("CheckBookAvailability() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("CheckBookAvailability() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr && available != tt.want {
				t.Errorf("CheckBookAvailability() = %v, want %v", available, tt.want)
			}
		})
	}
}