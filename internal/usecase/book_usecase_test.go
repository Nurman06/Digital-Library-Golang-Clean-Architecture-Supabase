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

// MockBookRepository is a mock implementation of BookRepository
type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) Create(ctx context.Context, book *entity.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *MockBookRepository) GetByID(ctx context.Context, id string) (*entity.Book, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Book), args.Error(1)
}

func (m *MockBookRepository) GetByISBN(ctx context.Context, isbn string) (*entity.Book, error) {
	args := m.Called(ctx, isbn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Book), args.Error(1)
}

func (m *MockBookRepository) Update(ctx context.Context, book *entity.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *MockBookRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBookRepository) HardDelete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBookRepository) Restore(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBookRepository) List(ctx context.Context, params repository.ListParams) ([]*entity.Book, int64, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookRepository) Search(ctx context.Context, params repository.SearchParams) ([]*entity.Book, int64, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookRepository) GetByCategory(ctx context.Context, category string, params repository.ListParams) ([]*entity.Book, int64, error) {
	args := m.Called(ctx, category, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookRepository) GetByAuthor(ctx context.Context, author string, params repository.ListParams) ([]*entity.Book, int64, error) {
	args := m.Called(ctx, author, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookRepository) GetRecentlyAdded(ctx context.Context, limit int) ([]*entity.Book, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Book), args.Error(1)
}

func (m *MockBookRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) ExistsByISBN(ctx context.Context, isbn string) (bool, error) {
	args := m.Called(ctx, isbn)
	return args.Bool(0), args.Error(1)
}

// MockBookCopyRepository is a mock implementation of BookCopyRepository
type MockBookCopyRepository struct {
	mock.Mock
}

func (m *MockBookCopyRepository) Create(ctx context.Context, copy *entity.BookCopy) error {
	args := m.Called(ctx, copy)
	return args.Error(0)
}

func (m *MockBookCopyRepository) GetByID(ctx context.Context, id string) (*entity.BookCopy, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BookCopy), args.Error(1)
}

func (m *MockBookCopyRepository) GetByBookID(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.BookCopy), args.Error(1)
}

func (m *MockBookCopyRepository) GetAvailableCopies(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.BookCopy), args.Error(1)
}

func (m *MockBookCopyRepository) Update(ctx context.Context, copy *entity.BookCopy) error {
	args := m.Called(ctx, copy)
	return args.Error(0)
}

func (m *MockBookCopyRepository) UpdateStatus(ctx context.Context, id string, status entity.CopyStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockBookCopyRepository) GetByCopyNumber(ctx context.Context, copyNumber string) (*entity.BookCopy, error) {
	args := m.Called(ctx, copyNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BookCopy), args.Error(1)
}

func (m *MockBookCopyRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBookCopyRepository) GetByStatus(ctx context.Context, status entity.CopyStatus, params repository.ListParams) ([]*entity.BookCopy, int64, error) {
	args := m.Called(ctx, status, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.BookCopy), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookCopyRepository) CountByBookID(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookCopyRepository) GetByLocation(ctx context.Context, location string, params repository.ListParams) ([]*entity.BookCopy, int64, error) {
	args := m.Called(ctx, location, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.BookCopy), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookCopyRepository) List(ctx context.Context, params repository.ListParams) ([]*entity.BookCopy, int64, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.BookCopy), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookCopyRepository) CountAvailableByBookID(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookCopyRepository) ExistsByCopyNumber(ctx context.Context, copyNumber string) (bool, error) {
	args := m.Called(ctx, copyNumber)
	return args.Bool(0), args.Error(1)
}

// Test CreateBook
func TestBookUseCase_CreateBook(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		book := &entity.Book{
			Title:           "Clean Code",
			Author:          "Robert C. Martin",
			ISBN:            "9780132350884",
			Category:        "Programming",
			PublicationYear: 2008,
		}

		mockBookRepo.On("ExistsByISBN", ctx, book.ISBN).Return(false, nil)
		mockBookRepo.On("Create", ctx, book).Return(nil)

		err := uc.CreateBook(ctx, book)

		assert.NoError(t, err)
		assert.NotZero(t, book.CreatedAt)
		assert.NotZero(t, book.UpdatedAt)
		mockBookRepo.AssertExpectations(t)
	})

	t.Run("validation error - missing title", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		book := &entity.Book{
			Title:  "",
			Author: "Robert C. Martin",
			ISBN:   "9780132350884",
		}

		err := uc.CreateBook(ctx, book)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("duplicate ISBN error", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		book := &entity.Book{
			Title:           "Clean Code",
			Author:          "Robert C. Martin",
			ISBN:            "9780132350884",
			PublicationYear: 2008,
		}

		mockBookRepo.On("ExistsByISBN", ctx, book.ISBN).Return(true, nil)

		err := uc.CreateBook(ctx, book)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
		mockBookRepo.AssertExpectations(t)
	})
}

// Test GetBookByID
func TestBookUseCase_GetBookByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		expectedBook := &entity.Book{
			ID:     "book-1",
			Title:  "Clean Code",
			Author: "Robert C. Martin",
			ISBN:   "9780132350884",
		}

		mockBookRepo.On("GetByID", ctx, "book-1").Return(expectedBook, nil)

		book, err := uc.GetBookByID(ctx, "book-1")

		assert.NoError(t, err)
		assert.Equal(t, expectedBook, book)
		mockBookRepo.AssertExpectations(t)
	})

	t.Run("empty ID error", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		book, err := uc.GetBookByID(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, book)
		assert.Contains(t, err.Error(), "required")
	})

	t.Run("not found error", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		mockBookRepo.On("GetByID", ctx, "nonexistent").Return(nil, errors.New("not found"))

		book, err := uc.GetBookByID(ctx, "nonexistent")

		assert.Error(t, err)
		assert.Nil(t, book)
		mockBookRepo.AssertExpectations(t)
	})
}

// Test UpdateBook
func TestBookUseCase_UpdateBook(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		existingBook := &entity.Book{
			ID:              "book-1",
			Title:           "Clean Code",
			Author:          "Robert C. Martin",
			ISBN:            "9780132350884",
			PublicationYear: 2008,
		}

		updatedBook := &entity.Book{
			ID:              "book-1",
			Title:           "Clean Code - Updated",
			Author:          "Robert C. Martin",
			ISBN:            "9780132350884",
			PublicationYear: 2008,
		}

		mockBookRepo.On("GetByID", ctx, "book-1").Return(existingBook, nil)
		mockBookRepo.On("Update", ctx, updatedBook).Return(nil)

		err := uc.UpdateBook(ctx, updatedBook)

		assert.NoError(t, err)
		mockBookRepo.AssertExpectations(t)
	})

	t.Run("ISBN change - duplicate error", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		existingBook := &entity.Book{
			ID:              "book-1",
			ISBN:            "9780132350884",
			Title:           "Clean Code",
			Author:          "Robert C. Martin",
			PublicationYear: 2008,
		}

		updatedBook := &entity.Book{
			ID:              "book-1",
			ISBN:            "9781234567890", // Different ISBN
			Title:           "Clean Code",
			Author:          "Robert C. Martin",
			PublicationYear: 2008,
		}

		mockBookRepo.On("GetByID", ctx, "book-1").Return(existingBook, nil)
		mockBookRepo.On("ExistsByISBN", ctx, "9781234567890").Return(true, nil)

		err := uc.UpdateBook(ctx, updatedBook)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
		mockBookRepo.AssertExpectations(t)
	})
}

// Test DeleteBook
func TestBookUseCase_DeleteBook(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		book := &entity.Book{
			ID:        "book-1",
			DeletedAt: nil,
		}

		copies := []*entity.BookCopy{
			{ID: "copy-1", Status: entity.CopyStatusAvailable},
		}

		mockBookRepo.On("GetByID", ctx, "book-1").Return(book, nil)
		mockCopyRepo.On("GetByBookID", ctx, "book-1").Return(copies, nil)
		mockBookRepo.On("Delete", ctx, "book-1").Return(nil)

		err := uc.DeleteBook(ctx, "book-1")

		assert.NoError(t, err)
		mockBookRepo.AssertExpectations(t)
		mockCopyRepo.AssertExpectations(t)
	})

	t.Run("error - has active borrows", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		book := &entity.Book{
			ID:        "book-1",
			DeletedAt: nil,
		}

		copies := []*entity.BookCopy{
			{ID: "copy-1", Status: entity.CopyStatusBorrowed},
		}

		mockBookRepo.On("GetByID", ctx, "book-1").Return(book, nil)
		mockCopyRepo.On("GetByBookID", ctx, "book-1").Return(copies, nil)

		err := uc.DeleteBook(ctx, "book-1")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "active borrows")
		mockBookRepo.AssertExpectations(t)
		mockCopyRepo.AssertExpectations(t)
	})
}

// Test CreateBookCopy
func TestBookUseCase_CreateBookCopy(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		book := &entity.Book{ID: "book-1"}
		copy := &entity.BookCopy{
			BookID:     "book-1",
			CopyNumber: "COPY-001",
			Status:     entity.CopyStatusAvailable,
		}

		mockBookRepo.On("GetByID", ctx, "book-1").Return(book, nil)
		mockCopyRepo.On("ExistsByCopyNumber", ctx, "COPY-001").Return(false, nil)
		mockCopyRepo.On("Create", ctx, copy).Return(nil)

		err := uc.CreateBookCopy(ctx, copy)

		assert.NoError(t, err)
		assert.NotZero(t, copy.CreatedAt)
		mockBookRepo.AssertExpectations(t)
		mockCopyRepo.AssertExpectations(t)
	})

	t.Run("duplicate copy number error", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		book := &entity.Book{ID: "book-1"}
		copy := &entity.BookCopy{
			BookID:     "book-1",
			CopyNumber: "COPY-001",
			Status:     entity.CopyStatusAvailable,
		}

		mockBookRepo.On("GetByID", ctx, "book-1").Return(book, nil)
		mockCopyRepo.On("ExistsByCopyNumber", ctx, "COPY-001").Return(true, nil)

		err := uc.CreateBookCopy(ctx, copy)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
		mockBookRepo.AssertExpectations(t)
		mockCopyRepo.AssertExpectations(t)
	})
}

// Test ListBooks
func TestBookUseCase_ListBooks(t *testing.T) {
	ctx := context.Background()

	t.Run("success with default pagination", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		expectedBooks := []*entity.Book{
			{ID: "book-1", Title: "Book 1"},
			{ID: "book-2", Title: "Book 2"},
		}

		params := repository.ListParams{Page: 0, PageSize: 0}
		expectedParams := repository.ListParams{Page: 1, PageSize: 20}

		mockBookRepo.On("List", ctx, expectedParams).Return(expectedBooks, int64(2), nil)

		books, total, err := uc.ListBooks(ctx, params)

		assert.NoError(t, err)
		assert.Equal(t, expectedBooks, books)
		assert.Equal(t, int64(2), total)
		mockBookRepo.AssertExpectations(t)
	})

	t.Run("success with custom pagination", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		expectedBooks := []*entity.Book{
			{ID: "book-1", Title: "Book 1"},
		}

		params := repository.ListParams{Page: 2, PageSize: 10}

		mockBookRepo.On("List", ctx, params).Return(expectedBooks, int64(15), nil)

		books, total, err := uc.ListBooks(ctx, params)

		assert.NoError(t, err)
		assert.Equal(t, expectedBooks, books)
		assert.Equal(t, int64(15), total)
		mockBookRepo.AssertExpectations(t)
	})
}

// Test GetBookAvailabilityCount
func TestBookUseCase_GetBookAvailabilityCount(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		book := &entity.Book{ID: "book-1"}

		mockBookRepo.On("GetByID", ctx, "book-1").Return(book, nil)
		mockCopyRepo.On("CountAvailableByBookID", ctx, "book-1").Return(int64(3), nil)

		count, err := uc.GetBookAvailabilityCount(ctx, "book-1")

		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)
		mockBookRepo.AssertExpectations(t)
		mockCopyRepo.AssertExpectations(t)
	})

	t.Run("empty book ID error", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		count, err := uc.GetBookAvailabilityCount(ctx, "")

		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
		assert.Contains(t, err.Error(), "required")
	})
}

// Test DeleteBookCopy
func TestBookUseCase_DeleteBookCopy(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		copy := &entity.BookCopy{
			ID:     "copy-1",
			Status: entity.CopyStatusAvailable,
		}

		mockCopyRepo.On("GetByID", ctx, "copy-1").Return(copy, nil)
		mockCopyRepo.On("Delete", ctx, "copy-1").Return(nil)

		err := uc.DeleteBookCopy(ctx, "copy-1")

		assert.NoError(t, err)
		mockCopyRepo.AssertExpectations(t)
	})

	t.Run("error - copy is borrowed", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		copy := &entity.BookCopy{
			ID:     "copy-1",
			Status: entity.CopyStatusBorrowed,
		}

		mockCopyRepo.On("GetByID", ctx, "copy-1").Return(copy, nil)

		err := uc.DeleteBookCopy(ctx, "copy-1")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "currently borrowed")
		mockCopyRepo.AssertExpectations(t)
	})
}

// Test RestoreBook
func TestBookUseCase_RestoreBook(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		deletedAt := time.Now()
		book := &entity.Book{
			ID:        "book-1",
			DeletedAt: &deletedAt,
		}

		mockBookRepo.On("GetByID", ctx, "book-1").Return(book, nil)
		mockBookRepo.On("Restore", ctx, "book-1").Return(nil)

		err := uc.RestoreBook(ctx, "book-1")

		assert.NoError(t, err)
		mockBookRepo.AssertExpectations(t)
	})

	t.Run("error - book not deleted", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		book := &entity.Book{
			ID:        "book-1",
			DeletedAt: nil,
		}

		mockBookRepo.On("GetByID", ctx, "book-1").Return(book, nil)

		err := uc.RestoreBook(ctx, "book-1")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not deleted")
		mockBookRepo.AssertExpectations(t)
	})
}

// Test GetBookCount
func TestBookUseCase_GetBookCount(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		mockBookRepo.On("Count", ctx).Return(int64(42), nil)

		count, err := uc.GetBookCount(ctx)

		assert.NoError(t, err)
		assert.Equal(t, int64(42), count)
		mockBookRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockBookRepo := new(MockBookRepository)
		mockCopyRepo := new(MockBookCopyRepository)
		uc := NewBookUseCase(mockBookRepo, mockCopyRepo)

		mockBookRepo.On("Count", ctx).Return(int64(0), errors.New("database error"))

		count, err := uc.GetBookCount(ctx)

		require.Error(t, err)
		assert.Equal(t, int64(0), count)
		mockBookRepo.AssertExpectations(t)
	})
}