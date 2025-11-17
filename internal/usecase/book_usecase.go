package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/repository"
)

// BookUseCase defines the interface for book business logic operations
type BookUseCase interface {
	// CreateBook creates a new book in the catalog
	CreateBook(ctx context.Context, book *entity.Book) error

	// GetBookByID retrieves a book by its ID
	GetBookByID(ctx context.Context, id string) (*entity.Book, error)

	// GetBookByISBN retrieves a book by its ISBN
	GetBookByISBN(ctx context.Context, isbn string) (*entity.Book, error)

	// UpdateBook updates an existing book
	UpdateBook(ctx context.Context, book *entity.Book) error

	// DeleteBook performs a soft delete on a book
	DeleteBook(ctx context.Context, id string) error

	// RestoreBook restores a soft-deleted book
	RestoreBook(ctx context.Context, id string) error

	// ListBooks retrieves a paginated list of books
	ListBooks(ctx context.Context, params repository.ListParams) ([]*entity.Book, int64, error)

	// SearchBooks searches for books based on search criteria
	SearchBooks(ctx context.Context, params repository.SearchParams) ([]*entity.Book, int64, error)

	// GetBooksByCategory retrieves books by category
	GetBooksByCategory(ctx context.Context, category string, params repository.ListParams) ([]*entity.Book, int64, error)

	// GetBooksByAuthor retrieves books by author
	GetBooksByAuthor(ctx context.Context, author string, params repository.ListParams) ([]*entity.Book, int64, error)

	// GetRecentlyAddedBooks retrieves recently added books
	GetRecentlyAddedBooks(ctx context.Context, limit int) ([]*entity.Book, error)

	// GetBookCount returns the total number of books
	GetBookCount(ctx context.Context) (int64, error)

	// CreateBookCopy creates a new copy of a book
	CreateBookCopy(ctx context.Context, copy *entity.BookCopy) error

	// GetBookCopyByID retrieves a book copy by its ID
	GetBookCopyByID(ctx context.Context, id string) (*entity.BookCopy, error)

	// GetBookCopies retrieves all copies of a specific book
	GetBookCopies(ctx context.Context, bookID string) ([]*entity.BookCopy, error)

	// GetAvailableBookCopies retrieves available copies of a specific book
	GetAvailableBookCopies(ctx context.Context, bookID string) ([]*entity.BookCopy, error)

	// UpdateBookCopy updates an existing book copy
	UpdateBookCopy(ctx context.Context, copy *entity.BookCopy) error

	// UpdateBookCopyStatus updates the status of a book copy
	UpdateBookCopyStatus(ctx context.Context, id string, status entity.CopyStatus) error

	// DeleteBookCopy deletes a book copy
	DeleteBookCopy(ctx context.Context, id string) error

	// GetBookAvailabilityCount gets the count of available copies for a book
	GetBookAvailabilityCount(ctx context.Context, bookID string) (int64, error)
}

// bookUseCase implements the BookUseCase interface
type bookUseCase struct {
	bookRepo     repository.BookRepository
	bookCopyRepo repository.BookCopyRepository
}

// NewBookUseCase creates a new instance of BookUseCase
func NewBookUseCase(bookRepo repository.BookRepository, bookCopyRepo repository.BookCopyRepository) BookUseCase {
	return &bookUseCase{
		bookRepo:     bookRepo,
		bookCopyRepo: bookCopyRepo,
	}
}

// CreateBook creates a new book in the catalog
func (uc *bookUseCase) CreateBook(ctx context.Context, book *entity.Book) error {
	// Validate the book entity
	if err := book.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if book with same ISBN already exists
	exists, err := uc.bookRepo.ExistsByISBN(ctx, book.ISBN)
	if err != nil {
		return fmt.Errorf("failed to check ISBN existence: %w", err)
	}
	if exists {
		return errors.New("book with this ISBN already exists")
	}

	// Set timestamps
	now := time.Now()
	book.CreatedAt = now
	book.UpdatedAt = now

	// Create the book
	if err := uc.bookRepo.Create(ctx, book); err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}

	return nil
}

// GetBookByID retrieves a book by its ID
func (uc *bookUseCase) GetBookByID(ctx context.Context, id string) (*entity.Book, error) {
	if id == "" {
		return nil, errors.New("book ID is required")
	}

	book, err := uc.bookRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	return book, nil
}

// GetBookByISBN retrieves a book by its ISBN
func (uc *bookUseCase) GetBookByISBN(ctx context.Context, isbn string) (*entity.Book, error) {
	if isbn == "" {
		return nil, errors.New("ISBN is required")
	}

	book, err := uc.bookRepo.GetByISBN(ctx, isbn)
	if err != nil {
		return nil, fmt.Errorf("failed to get book by ISBN: %w", err)
	}

	return book, nil
}

// UpdateBook updates an existing book
func (uc *bookUseCase) UpdateBook(ctx context.Context, book *entity.Book) error {
	// Validate the book entity
	if err := book.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if book exists
	existingBook, err := uc.bookRepo.GetByID(ctx, book.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing book: %w", err)
	}

	// If ISBN is being changed, check if new ISBN already exists
	if existingBook.ISBN != book.ISBN {
		exists, err := uc.bookRepo.ExistsByISBN(ctx, book.ISBN)
		if err != nil {
			return fmt.Errorf("failed to check ISBN existence: %w", err)
		}
		if exists {
			return errors.New("book with this ISBN already exists")
		}
	}

	// Update timestamp
	book.UpdatedAt = time.Now()

	// Update the book
	if err := uc.bookRepo.Update(ctx, book); err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}

	return nil
}

// DeleteBook performs a soft delete on a book
func (uc *bookUseCase) DeleteBook(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("book ID is required")
	}

	// Check if book exists
	book, err := uc.bookRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get book: %w", err)
	}

	// Check if book is already deleted
	if book.IsDeleted() {
		return errors.New("book is already deleted")
	}

	// Check if book has any active borrows
	copies, err := uc.bookCopyRepo.GetByBookID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get book copies: %w", err)
	}

	for _, copy := range copies {
		if copy.IsBorrowed() {
			return errors.New("cannot delete book with active borrows")
		}
	}

	// Perform soft delete
	if err := uc.bookRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	return nil
}

// RestoreBook restores a soft-deleted book
func (uc *bookUseCase) RestoreBook(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("book ID is required")
	}

	// Check if book exists
	book, err := uc.bookRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get book: %w", err)
	}

	// Check if book is deleted
	if !book.IsDeleted() {
		return errors.New("book is not deleted")
	}

	// Restore the book
	if err := uc.bookRepo.Restore(ctx, id); err != nil {
		return fmt.Errorf("failed to restore book: %w", err)
	}

	return nil
}

// ListBooks retrieves a paginated list of books
func (uc *bookUseCase) ListBooks(ctx context.Context, params repository.ListParams) ([]*entity.Book, int64, error) {
	// Set default pagination values if not provided
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100 // Max page size
	}

	books, total, err := uc.bookRepo.List(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list books: %w", err)
	}

	return books, total, nil
}

// SearchBooks searches for books based on search criteria
func (uc *bookUseCase) SearchBooks(ctx context.Context, params repository.SearchParams) ([]*entity.Book, int64, error) {
	// Set default pagination values if not provided
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100 // Max page size
	}

	books, total, err := uc.bookRepo.Search(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search books: %w", err)
	}

	return books, total, nil
}

// GetBooksByCategory retrieves books by category
func (uc *bookUseCase) GetBooksByCategory(ctx context.Context, category string, params repository.ListParams) ([]*entity.Book, int64, error) {
	if category == "" {
		return nil, 0, errors.New("category is required")
	}

	// Set default pagination values if not provided
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100 // Max page size
	}

	books, total, err := uc.bookRepo.GetByCategory(ctx, category, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get books by category: %w", err)
	}

	return books, total, nil
}

// GetBooksByAuthor retrieves books by author
func (uc *bookUseCase) GetBooksByAuthor(ctx context.Context, author string, params repository.ListParams) ([]*entity.Book, int64, error) {
	if author == "" {
		return nil, 0, errors.New("author is required")
	}

	// Set default pagination values if not provided
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100 // Max page size
	}

	books, total, err := uc.bookRepo.GetByAuthor(ctx, author, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get books by author: %w", err)
	}

	return books, total, nil
}

// GetRecentlyAddedBooks retrieves recently added books
func (uc *bookUseCase) GetRecentlyAddedBooks(ctx context.Context, limit int) ([]*entity.Book, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 50 {
		limit = 50 // Max limit
	}

	books, err := uc.bookRepo.GetRecentlyAdded(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recently added books: %w", err)
	}

	return books, nil
}

// GetBookCount returns the total number of books
func (uc *bookUseCase) GetBookCount(ctx context.Context) (int64, error) {
	count, err := uc.bookRepo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count books: %w", err)
	}

	return count, nil
}

// CreateBookCopy creates a new copy of a book
func (uc *bookUseCase) CreateBookCopy(ctx context.Context, copy *entity.BookCopy) error {
	// Validate the book copy entity
	if err := copy.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if book exists
	_, err := uc.bookRepo.GetByID(ctx, copy.BookID)
	if err != nil {
		return fmt.Errorf("failed to get book: %w", err)
	}

	// Check if copy number already exists
	exists, err := uc.bookCopyRepo.ExistsByCopyNumber(ctx, copy.CopyNumber)
	if err != nil {
		return fmt.Errorf("failed to check copy number existence: %w", err)
	}
	if exists {
		return errors.New("book copy with this copy number already exists")
	}

	// Set timestamps
	now := time.Now()
	copy.CreatedAt = now
	copy.UpdatedAt = now

	// Create the book copy
	if err := uc.bookCopyRepo.Create(ctx, copy); err != nil {
		return fmt.Errorf("failed to create book copy: %w", err)
	}

	return nil
}

// GetBookCopyByID retrieves a book copy by its ID
func (uc *bookUseCase) GetBookCopyByID(ctx context.Context, id string) (*entity.BookCopy, error) {
	if id == "" {
		return nil, errors.New("book copy ID is required")
	}

	copy, err := uc.bookCopyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get book copy: %w", err)
	}

	return copy, nil
}

// GetBookCopies retrieves all copies of a specific book
func (uc *bookUseCase) GetBookCopies(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
	if bookID == "" {
		return nil, errors.New("book ID is required")
	}

	// Check if book exists
	_, err := uc.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	copies, err := uc.bookCopyRepo.GetByBookID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book copies: %w", err)
	}

	return copies, nil
}

// GetAvailableBookCopies retrieves available copies of a specific book
func (uc *bookUseCase) GetAvailableBookCopies(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
	if bookID == "" {
		return nil, errors.New("book ID is required")
	}

	// Check if book exists
	_, err := uc.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	copies, err := uc.bookCopyRepo.GetAvailableCopies(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get available book copies: %w", err)
	}

	return copies, nil
}

// UpdateBookCopy updates an existing book copy
func (uc *bookUseCase) UpdateBookCopy(ctx context.Context, copy *entity.BookCopy) error {
	// Validate the book copy entity
	if err := copy.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if book copy exists
	existingCopy, err := uc.bookCopyRepo.GetByID(ctx, copy.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing book copy: %w", err)
	}

	// If copy number is being changed, check if new copy number already exists
	if existingCopy.CopyNumber != copy.CopyNumber {
		exists, err := uc.bookCopyRepo.ExistsByCopyNumber(ctx, copy.CopyNumber)
		if err != nil {
			return fmt.Errorf("failed to check copy number existence: %w", err)
		}
		if exists {
			return errors.New("book copy with this copy number already exists")
		}
	}

	// Update timestamp
	copy.UpdatedAt = time.Now()

	// Update the book copy
	if err := uc.bookCopyRepo.Update(ctx, copy); err != nil {
		return fmt.Errorf("failed to update book copy: %w", err)
	}

	return nil
}

// UpdateBookCopyStatus updates the status of a book copy
func (uc *bookUseCase) UpdateBookCopyStatus(ctx context.Context, id string, status entity.CopyStatus) error {
	if id == "" {
		return errors.New("book copy ID is required")
	}

	// Check if book copy exists
	_, err := uc.bookCopyRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get book copy: %w", err)
	}

	// Update the status
	if err := uc.bookCopyRepo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update book copy status: %w", err)
	}

	return nil
}

// DeleteBookCopy deletes a book copy
func (uc *bookUseCase) DeleteBookCopy(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("book copy ID is required")
	}

	// Check if book copy exists
	copy, err := uc.bookCopyRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get book copy: %w", err)
	}

	// Check if book copy is currently borrowed
	if copy.IsBorrowed() {
		return errors.New("cannot delete book copy that is currently borrowed")
	}

	// Delete the book copy
	if err := uc.bookCopyRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete book copy: %w", err)
	}

	return nil
}

// GetBookAvailabilityCount gets the count of available copies for a book
func (uc *bookUseCase) GetBookAvailabilityCount(ctx context.Context, bookID string) (int64, error) {
	if bookID == "" {
		return 0, errors.New("book ID is required")
	}

	// Check if book exists
	_, err := uc.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return 0, fmt.Errorf("failed to get book: %w", err)
	}

	count, err := uc.bookCopyRepo.CountAvailableByBookID(ctx, bookID)
	if err != nil {
		return 0, fmt.Errorf("failed to count available book copies: %w", err)
	}

	return count, nil
}