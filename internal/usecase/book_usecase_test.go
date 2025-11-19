package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/repository"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/usecase/mocks"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/usecase/testdata"
)

func TestBookUseCase_CreateBook(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		book          *entity.Book
		setupMocks    func(*mocks.MockBookRepository, *mocks.MockBookCopyRepository)
		wantErr       bool
		errorContains string
	}{
		{
			name: "successful book creation",
			book: testdata.CreateTestBook(""),
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {
				br.ExistsByISBNFunc = func(ctx context.Context, isbn string) (bool, error) {
					return false, nil
				}
				br.CreateFunc = func(ctx context.Context, book *entity.Book) error {
					book.ID = "new-book-id"
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "validation error - missing title",
			book: &entity.Book{
				Author:          "Test Author",
				ISBN:            "9780134190440",
				PublicationYear: 2020,
			},
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {},
			wantErr:    true,
			errorContains: "validation failed",
		},
		{
			name: "duplicate ISBN error",
			book: testdata.CreateTestBook(""),
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {
				br.ExistsByISBNFunc = func(ctx context.Context, isbn string) (bool, error) {
					return true, nil
				}
			},
			wantErr:       true,
			errorContains: "ISBN already exists",
		},
		{
			name: "repository error on ISBN check",
			book: testdata.CreateTestBook(""),
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {
				br.ExistsByISBNFunc = func(ctx context.Context, isbn string) (bool, error) {
					return false, errors.New("database error")
				}
			},
			wantErr:       true,
			errorContains: "failed to check ISBN existence",
		},
		{
			name: "repository error on create",
			book: testdata.CreateTestBook(""),
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {
				br.ExistsByISBNFunc = func(ctx context.Context, isbn string) (bool, error) {
					return false, nil
				}
				br.CreateFunc = func(ctx context.Context, book *entity.Book) error {
					return errors.New("database error")
				}
			},
			wantErr:       true,
			errorContains: "failed to create book",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookRepo := &mocks.MockBookRepository{}
			bookCopyRepo := &mocks.MockBookCopyRepository{}
			tt.setupMocks(bookRepo, bookCopyRepo)

			uc := NewBookUseCase(bookRepo, bookCopyRepo)
			err := uc.CreateBook(ctx, tt.book)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("CreateBook() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr {
				if tt.book.CreatedAt.IsZero() {
					t.Error("CreatedAt should be set")
				}
				if tt.book.UpdatedAt.IsZero() {
					t.Error("UpdatedAt should be set")
				}
			}
		})
	}
}

func TestBookUseCase_GetBookByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		bookID        string
		setupMocks    func(*mocks.MockBookRepository)
		wantErr       bool
		errorContains string
	}{
		{
			name:   "successful retrieval",
			bookID: "book-123",
			setupMocks: func(br *mocks.MockBookRepository) {
				br.GetByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return testdata.CreateTestBook(id), nil
				}
			},
			wantErr: false,
		},
		{
			name:          "empty book ID",
			bookID:        "",
			setupMocks:    func(br *mocks.MockBookRepository) {},
			wantErr:       true,
			errorContains: "book ID is required",
		},
		{
			name:   "book not found",
			bookID: "nonexistent",
			setupMocks: func(br *mocks.MockBookRepository) {
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
			tt.setupMocks(bookRepo)

			uc := NewBookUseCase(bookRepo, bookCopyRepo)
			book, err := uc.GetBookByID(ctx, tt.bookID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetBookByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("GetBookByID() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr && book == nil {
				t.Error("GetBookByID() should return a book")
			}
		})
	}
}

func TestBookUseCase_UpdateBook(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		book          *entity.Book
		setupMocks    func(*mocks.MockBookRepository)
		wantErr       bool
		errorContains string
	}{
		{
			name: "successful update",
			book: testdata.CreateTestBook("book-123"),
			setupMocks: func(br *mocks.MockBookRepository) {
				br.GetByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return testdata.CreateTestBook(id), nil
				}
				br.ExistsByISBNFunc = func(ctx context.Context, isbn string) (bool, error) {
					return false, nil
				}
				br.UpdateFunc = func(ctx context.Context, book *entity.Book) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "validation error",
			book: &entity.Book{
				ID:              "book-123",
				Author:          "Test Author",
				ISBN:            "9780134190440",
				PublicationYear: 2020,
			},
			setupMocks: func(br *mocks.MockBookRepository) {},
			wantErr:    true,
			errorContains: "validation failed",
		},
		{
			name: "book not found",
			book: testdata.CreateTestBook("book-123"),
			setupMocks: func(br *mocks.MockBookRepository) {
				br.GetByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr:       true,
			errorContains: "failed to get existing book",
		},
		{
			name: "duplicate ISBN on update",
			book: func() *entity.Book {
				book := testdata.CreateTestBook("book-123")
				book.ISBN = "9999999999999"
				return book
			}(),
			setupMocks: func(br *mocks.MockBookRepository) {
				br.GetByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return testdata.CreateTestBook(id), nil
				}
				br.ExistsByISBNFunc = func(ctx context.Context, isbn string) (bool, error) {
					return true, nil
				}
			},
			wantErr:       true,
			errorContains: "ISBN already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookRepo := &mocks.MockBookRepository{}
			bookCopyRepo := &mocks.MockBookCopyRepository{}
			tt.setupMocks(bookRepo)

			uc := NewBookUseCase(bookRepo, bookCopyRepo)
			err := uc.UpdateBook(ctx, tt.book)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("UpdateBook() error = %v, should contain %v", err, tt.errorContains)
				}
			}
		})
	}
}

func TestBookUseCase_DeleteBook(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		bookID        string
		setupMocks    func(*mocks.MockBookRepository, *mocks.MockBookCopyRepository)
		wantErr       bool
		errorContains string
	}{
		{
			name:   "successful deletion",
			bookID: "book-123",
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {
				br.GetByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return testdata.CreateTestBook(id), nil
				}
				bcr.GetByBookIDFunc = func(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
					copies := testdata.CreateTestBookCopies(bookID, 2)
					copies[0].Status = entity.CopyStatusAvailable
					copies[1].Status = entity.CopyStatusDamaged
					return copies, nil
				}
				br.DeleteFunc = func(ctx context.Context, id string) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:          "empty book ID",
			bookID:        "",
			setupMocks:    func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {},
			wantErr:       true,
			errorContains: "book ID is required",
		},
		{
			name:   "book not found",
			bookID: "nonexistent",
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {
				br.GetByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr:       true,
			errorContains: "failed to get book",
		},
		{
			name:   "already deleted book",
			bookID: "book-123",
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {
				book := testdata.CreateTestBook("book-123")
				now := time.Now()
				book.DeletedAt = &now
				br.GetByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return book, nil
				}
			},
			wantErr:       true,
			errorContains: "already deleted",
		},
		{
			name:   "cannot delete book with active borrows",
			bookID: "book-123",
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {
				br.GetByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return testdata.CreateTestBook(id), nil
				}
				bcr.GetByBookIDFunc = func(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
					copies := testdata.CreateTestBookCopies(bookID, 1)
					copies[0].Status = entity.CopyStatusBorrowed
					return copies, nil
				}
			},
			wantErr:       true,
			errorContains: "cannot delete book with active borrows",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookRepo := &mocks.MockBookRepository{}
			bookCopyRepo := &mocks.MockBookCopyRepository{}
			tt.setupMocks(bookRepo, bookCopyRepo)

			uc := NewBookUseCase(bookRepo, bookCopyRepo)
			err := uc.DeleteBook(ctx, tt.bookID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("DeleteBook() error = %v, should contain %v", err, tt.errorContains)
				}
			}
		})
	}
}

func TestBookUseCase_ListBooks(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		params     repository.ListParams
		setupMocks func(*mocks.MockBookRepository)
		wantCount  int
		wantTotal  int64
		wantErr    bool
	}{
		{
			name: "successful list with default params",
			params: repository.ListParams{
				Page:     0,
				PageSize: 0,
			},
			setupMocks: func(br *mocks.MockBookRepository) {
				br.ListFunc = func(ctx context.Context, params repository.ListParams) ([]*entity.Book, int64, error) {
					// Should set defaults: Page=1, PageSize=20
					if params.Page != 1 || params.PageSize != 20 {
						t.Errorf("Expected defaults Page=1, PageSize=20, got Page=%d, PageSize=%d", params.Page, params.PageSize)
					}
					books := testdata.CreateTestBooks(5)
					return books, 5, nil
				}
			},
			wantCount: 5,
			wantTotal: 5,
			wantErr:   false,
		},
		{
			name: "successful list with custom params",
			params: repository.ListParams{
				Page:     2,
				PageSize: 10,
			},
			setupMocks: func(br *mocks.MockBookRepository) {
				br.ListFunc = func(ctx context.Context, params repository.ListParams) ([]*entity.Book, int64, error) {
					books := testdata.CreateTestBooks(10)
					return books, 50, nil
				}
			},
			wantCount: 10,
			wantTotal: 50,
			wantErr:   false,
		},
		{
			name: "page size exceeds maximum",
			params: repository.ListParams{
				Page:     1,
				PageSize: 150,
			},
			setupMocks: func(br *mocks.MockBookRepository) {
				br.ListFunc = func(ctx context.Context, params repository.ListParams) ([]*entity.Book, int64, error) {
					// Should cap at 100
					if params.PageSize != 100 {
						t.Errorf("Expected PageSize=100, got PageSize=%d", params.PageSize)
					}
					return testdata.CreateTestBooks(10), 10, nil
				}
			},
			wantCount: 10,
			wantTotal: 10,
			wantErr:   false,
		},
		{
			name: "repository error",
			params: repository.ListParams{
				Page:     1,
				PageSize: 20,
			},
			setupMocks: func(br *mocks.MockBookRepository) {
				br.ListFunc = func(ctx context.Context, params repository.ListParams) ([]*entity.Book, int64, error) {
					return nil, 0, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookRepo := &mocks.MockBookRepository{}
			bookCopyRepo := &mocks.MockBookCopyRepository{}
			tt.setupMocks(bookRepo)

			uc := NewBookUseCase(bookRepo, bookCopyRepo)
			books, total, err := uc.ListBooks(ctx, tt.params)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListBooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(books) != tt.wantCount {
					t.Errorf("ListBooks() returned %d books, want %d", len(books), tt.wantCount)
				}
				if total != tt.wantTotal {
					t.Errorf("ListBooks() total = %d, want %d", total, tt.wantTotal)
				}
			}
		})
	}
}

func TestBookUseCase_CreateBookCopy(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		copy          *entity.BookCopy
		setupMocks    func(*mocks.MockBookRepository, *mocks.MockBookCopyRepository)
		wantErr       bool
		errorContains string
	}{
		{
			name: "successful book copy creation",
			copy: testdata.CreateTestBookCopy("", "book-123"),
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {
				br.GetByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return testdata.CreateTestBook(id), nil
				}
				bcr.ExistsByCopyNumberFunc = func(ctx context.Context, copyNumber string) (bool, error) {
					return false, nil
				}
				bcr.CreateFunc = func(ctx context.Context, copy *entity.BookCopy) error {
					copy.ID = "new-copy-id"
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "validation error",
			copy: &entity.BookCopy{
				CopyNumber: "COPY-001",
				Status:     entity.CopyStatusAvailable,
			},
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {},
			wantErr:    true,
			errorContains: "validation failed",
		},
		{
			name: "book not found",
			copy: testdata.CreateTestBookCopy("", "nonexistent"),
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {
				br.GetByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr:       true,
			errorContains: "failed to get book",
		},
		{
			name: "duplicate copy number",
			copy: testdata.CreateTestBookCopy("", "book-123"),
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {
				br.GetByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return testdata.CreateTestBook(id), nil
				}
				bcr.ExistsByCopyNumberFunc = func(ctx context.Context, copyNumber string) (bool, error) {
					return true, nil
				}
			},
			wantErr:       true,
			errorContains: "copy number already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookRepo := &mocks.MockBookRepository{}
			bookCopyRepo := &mocks.MockBookCopyRepository{}
			tt.setupMocks(bookRepo, bookCopyRepo)

			uc := NewBookUseCase(bookRepo, bookCopyRepo)
			err := uc.CreateBookCopy(ctx, tt.copy)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBookCopy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("CreateBookCopy() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr {
				if tt.copy.CreatedAt.IsZero() {
					t.Error("CreatedAt should be set")
				}
				if tt.copy.UpdatedAt.IsZero() {
					t.Error("UpdatedAt should be set")
				}
			}
		})
	}
}

func TestBookUseCase_GetBookAvailabilityCount(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		bookID        string
		setupMocks    func(*mocks.MockBookRepository, *mocks.MockBookCopyRepository)
		wantCount     int64
		wantErr       bool
		errorContains string
	}{
		{
			name:   "successful count",
			bookID: "book-123",
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {
				br.GetByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return testdata.CreateTestBook(id), nil
				}
				bcr.CountAvailableByBookIDFunc = func(ctx context.Context, bookID string) (int64, error) {
					return 3, nil
				}
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:          "empty book ID",
			bookID:        "",
			setupMocks:    func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {},
			wantErr:       true,
			errorContains: "book ID is required",
		},
		{
			name:   "book not found",
			bookID: "nonexistent",
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {
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
			tt.setupMocks(bookRepo, bookCopyRepo)

			uc := NewBookUseCase(bookRepo, bookCopyRepo)
			count, err := uc.GetBookAvailabilityCount(ctx, tt.bookID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetBookAvailabilityCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("GetBookAvailabilityCount() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr && count != tt.wantCount {
				t.Errorf("GetBookAvailabilityCount() = %d, want %d", count, tt.wantCount)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}