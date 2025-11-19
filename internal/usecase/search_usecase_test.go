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

func TestSearchUseCase_SearchBooks(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		params        SearchBooksParams
		setupMocks    func(*mocks.MockBookRepository, *mocks.MockBookCopyRepository)
		wantCount     int
		wantTotal     int64
		wantErr       bool
		errorContains string
	}{
		{
			name: "successful search by query",
			params: SearchBooksParams{
				Query:    "Go Programming",
				Page:     1,
				PageSize: 20,
			},
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {
				br.SearchFunc = func(ctx context.Context, params repository.SearchParams) ([]*entity.Book, int64, error) {
					books := testdata.CreateTestBooks(3)
					return books, 3, nil
				}
			},
			wantCount: 3,
			wantTotal: 3,
			wantErr:   false,
		},
		{
			name: "empty query error",
			params: SearchBooksParams{
				Query:    "",
				Page:     1,
				PageSize: 20,
			},
			setupMocks:    func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {},
			wantErr:       true,
			errorContains: "search query is required",
		},
		{
			name: "repository error",
			params: SearchBooksParams{
				Query:    "Test",
				Page:     1,
				PageSize: 20,
			},
			setupMocks: func(br *mocks.MockBookRepository, bcr *mocks.MockBookCopyRepository) {
				br.SearchFunc = func(ctx context.Context, params repository.SearchParams) ([]*entity.Book, int64, error) {
					return nil, 0, errors.New("database error")
				}
			},
			wantErr:       true,
			errorContains: "failed to search books",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookRepo := &mocks.MockBookRepository{}
			bookCopyRepo := &mocks.MockBookCopyRepository{}
			borrowRecordRepo := &mocks.MockBorrowRecordRepository{}
			tt.setupMocks(bookRepo, bookCopyRepo)

			uc := NewSearchUseCase(bookRepo, bookCopyRepo, borrowRecordRepo)
			result, err := uc.SearchBooks(ctx, tt.params)

			if (err != nil) != tt.wantErr {
				t.Errorf("SearchBooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("SearchBooks() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr {
				if len(result.Books) != tt.wantCount {
					t.Errorf("SearchBooks() returned %d books, want %d", len(result.Books), tt.wantCount)
				}
				if result.Total != tt.wantTotal {
					t.Errorf("SearchBooks() total = %d, want %d", result.Total, tt.wantTotal)
				}
			}
		})
	}
}

func TestSearchUseCase_SearchByAuthor(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		author        string
		params        PaginationParams
		setupMocks    func(*mocks.MockBookRepository)
		wantCount     int
		wantTotal     int64
		wantErr       bool
		errorContains string
	}{
		{
			name:   "successful search by author",
			author: "Robert Martin",
			params: PaginationParams{
				Page:     1,
				PageSize: 20,
			},
			setupMocks: func(br *mocks.MockBookRepository) {
				br.GetByAuthorFunc = func(ctx context.Context, author string, params repository.ListParams) ([]*entity.Book, int64, error) {
					books := testdata.CreateTestBooks(3)
					return books, 3, nil
				}
			},
			wantCount: 3,
			wantTotal: 3,
			wantErr:   false,
		},
		{
			name:   "empty author",
			author: "",
			params: PaginationParams{
				Page:     1,
				PageSize: 20,
			},
			setupMocks:    func(br *mocks.MockBookRepository) {},
			wantErr:       true,
			errorContains: "author is required",
		},
		{
			name:   "repository error",
			author: "Test Author",
			params: PaginationParams{
				Page:     1,
				PageSize: 20,
			},
			setupMocks: func(br *mocks.MockBookRepository) {
				br.GetByAuthorFunc = func(ctx context.Context, author string, params repository.ListParams) ([]*entity.Book, int64, error) {
					return nil, 0, errors.New("database error")
				}
			},
			wantErr:       true,
			errorContains: "failed to search books by author",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookRepo := &mocks.MockBookRepository{}
			bookCopyRepo := &mocks.MockBookCopyRepository{}
			borrowRecordRepo := &mocks.MockBorrowRecordRepository{}
			tt.setupMocks(bookRepo)

			uc := NewSearchUseCase(bookRepo, bookCopyRepo, borrowRecordRepo)
			result, err := uc.SearchByAuthor(ctx, tt.author, tt.params)

			if (err != nil) != tt.wantErr {
				t.Errorf("SearchByAuthor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("SearchByAuthor() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr {
				if len(result.Books) != tt.wantCount {
					t.Errorf("SearchByAuthor() returned %d books, want %d", len(result.Books), tt.wantCount)
				}
				if result.Total != tt.wantTotal {
					t.Errorf("SearchByAuthor() total = %d, want %d", result.Total, tt.wantTotal)
				}
			}
		})
	}
}

func TestSearchUseCase_SearchByCategory(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		category      string
		params        PaginationParams
		setupMocks    func(*mocks.MockBookRepository)
		wantCount     int
		wantTotal     int64
		wantErr       bool
		errorContains string
	}{
		{
			name:     "successful search by category",
			category: "Programming",
			params: PaginationParams{
				Page:     1,
				PageSize: 20,
			},
			setupMocks: func(br *mocks.MockBookRepository) {
				br.GetByCategoryFunc = func(ctx context.Context, category string, params repository.ListParams) ([]*entity.Book, int64, error) {
					books := testdata.CreateTestBooks(5)
					return books, 5, nil
				}
			},
			wantCount: 5,
			wantTotal: 5,
			wantErr:   false,
		},
		{
			name:     "empty category",
			category: "",
			params: PaginationParams{
				Page:     1,
				PageSize: 20,
			},
			setupMocks:    func(br *mocks.MockBookRepository) {},
			wantErr:       true,
			errorContains: "category is required",
		},
		{
			name:     "repository error",
			category: "Programming",
			params: PaginationParams{
				Page:     1,
				PageSize: 20,
			},
			setupMocks: func(br *mocks.MockBookRepository) {
				br.GetByCategoryFunc = func(ctx context.Context, category string, params repository.ListParams) ([]*entity.Book, int64, error) {
					return nil, 0, errors.New("database error")
				}
			},
			wantErr:       true,
			errorContains: "failed to search books by category",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookRepo := &mocks.MockBookRepository{}
			bookCopyRepo := &mocks.MockBookCopyRepository{}
			borrowRecordRepo := &mocks.MockBorrowRecordRepository{}
			tt.setupMocks(bookRepo)

			uc := NewSearchUseCase(bookRepo, bookCopyRepo, borrowRecordRepo)
			result, err := uc.SearchByCategory(ctx, tt.category, tt.params)

			if (err != nil) != tt.wantErr {
				t.Errorf("SearchByCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("SearchByCategory() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr {
				if len(result.Books) != tt.wantCount {
					t.Errorf("SearchByCategory() returned %d books, want %d", len(result.Books), tt.wantCount)
				}
				if result.Total != tt.wantTotal {
					t.Errorf("SearchByCategory() total = %d, want %d", result.Total, tt.wantTotal)
				}
			}
		})
	}
}

func TestSearchUseCase_SearchByISBN(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		isbn          string
		setupMocks    func(*mocks.MockBookRepository)
		wantBook      bool
		wantErr       bool
		errorContains string
	}{
		{
			name: "successful search by ISBN",
			isbn: "9780134190440",
			setupMocks: func(br *mocks.MockBookRepository) {
				br.GetByISBNFunc = func(ctx context.Context, isbn string) (*entity.Book, error) {
					return testdata.CreateTestBook("book-1"), nil
				}
			},
			wantBook: true,
			wantErr:  false,
		},
		{
			name:          "empty ISBN",
			isbn:          "",
			setupMocks:    func(br *mocks.MockBookRepository) {},
			wantErr:       true,
			errorContains: "ISBN is required",
		},
		{
			name: "ISBN with dashes",
			isbn: "978-0-13-419044-0",
			setupMocks: func(br *mocks.MockBookRepository) {
				br.GetByISBNFunc = func(ctx context.Context, isbn string) (*entity.Book, error) {
					// Should receive cleaned ISBN
					if isbn != "9780134190440" {
						t.Errorf("Expected cleaned ISBN, got %s", isbn)
					}
					return testdata.CreateTestBook("book-1"), nil
				}
			},
			wantBook: true,
			wantErr:  false,
		},
		{
			name: "repository error",
			isbn: "9780134190440",
			setupMocks: func(br *mocks.MockBookRepository) {
				br.GetByISBNFunc = func(ctx context.Context, isbn string) (*entity.Book, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr:       true,
			errorContains: "failed to search book by ISBN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookRepo := &mocks.MockBookRepository{}
			bookCopyRepo := &mocks.MockBookCopyRepository{}
			borrowRecordRepo := &mocks.MockBorrowRecordRepository{}
			tt.setupMocks(bookRepo)

			uc := NewSearchUseCase(bookRepo, bookCopyRepo, borrowRecordRepo)
			book, err := uc.SearchByISBN(ctx, tt.isbn)

			if (err != nil) != tt.wantErr {
				t.Errorf("SearchByISBN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("SearchByISBN() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr && tt.wantBook && book == nil {
				t.Error("SearchByISBN() should return a book")
			}
		})
	}
}