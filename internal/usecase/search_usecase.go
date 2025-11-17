package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/repository"
)

// SearchUseCase defines the interface for search business logic operations
type SearchUseCase interface {
	// SearchBooks performs a comprehensive search across books
	SearchBooks(ctx context.Context, params SearchBooksParams) (*SearchResult, error)

	// SearchByTitle searches books by title
	SearchByTitle(ctx context.Context, title string, params PaginationParams) (*SearchResult, error)

	// SearchByAuthor searches books by author
	SearchByAuthor(ctx context.Context, author string, params PaginationParams) (*SearchResult, error)

	// SearchByISBN searches books by ISBN
	SearchByISBN(ctx context.Context, isbn string) (*entity.Book, error)

	// SearchByCategory searches books by category
	SearchByCategory(ctx context.Context, category string, params PaginationParams) (*SearchResult, error)

	// AdvancedSearch performs an advanced search with multiple criteria
	AdvancedSearch(ctx context.Context, params AdvancedSearchParams) (*SearchResult, error)

	// GetAvailableCategories retrieves all unique book categories
	GetAvailableCategories(ctx context.Context) ([]string, error)

	// GetPopularBooks retrieves popular books based on borrowing statistics
	GetPopularBooks(ctx context.Context, limit int) ([]*entity.Book, error)

	// SuggestBooks suggests books based on user preferences or history
	SuggestBooks(ctx context.Context, userID string, limit int) ([]*entity.Book, error)
}

// SearchBooksParams defines parameters for general book search
type SearchBooksParams struct {
	Query      string
	Page       int
	PageSize   int
	SortBy     string
	SortOrder  string
	OnlyAvailable bool
}

// PaginationParams defines pagination parameters
type PaginationParams struct {
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
}

// AdvancedSearchParams defines parameters for advanced search
type AdvancedSearchParams struct {
	Title           string
	Author          string
	ISBN            string
	Category        string
	PublicationYear int
	YearFrom        int
	YearTo          int
	OnlyAvailable   bool
	Page            int
	PageSize        int
	SortBy          string
	SortOrder       string
}

// SearchResult represents the result of a search operation
type SearchResult struct {
	Books      []*entity.Book
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// searchUseCase implements the SearchUseCase interface
type searchUseCase struct {
	bookRepo         repository.BookRepository
	bookCopyRepo     repository.BookCopyRepository
	borrowRecordRepo repository.BorrowRecordRepository
}

// NewSearchUseCase creates a new instance of SearchUseCase
func NewSearchUseCase(
	bookRepo repository.BookRepository,
	bookCopyRepo repository.BookCopyRepository,
	borrowRecordRepo repository.BorrowRecordRepository,
) SearchUseCase {
	return &searchUseCase{
		bookRepo:         bookRepo,
		bookCopyRepo:     bookCopyRepo,
		borrowRecordRepo: borrowRecordRepo,
	}
}

// SearchBooks performs a comprehensive search across books
func (uc *searchUseCase) SearchBooks(ctx context.Context, params SearchBooksParams) (*SearchResult, error) {
	if params.Query == "" {
		return nil, errors.New("search query is required")
	}

	// Set default pagination values
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}

	// Build search parameters
	searchParams := repository.SearchParams{
		Query:     params.Query,
		Page:      params.Page,
		PageSize:  params.PageSize,
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
	}

	// Perform search
	books, total, err := uc.bookRepo.Search(ctx, searchParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search books: %w", err)
	}

	// Filter by availability if requested
	if params.OnlyAvailable {
		books, total = uc.filterAvailableBooks(ctx, books)
	}

	// Calculate total pages
	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &SearchResult{
		Books:      books,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

// SearchByTitle searches books by title
func (uc *searchUseCase) SearchByTitle(ctx context.Context, title string, params PaginationParams) (*SearchResult, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	// Set default pagination values
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}

	// Build search parameters
	searchParams := repository.SearchParams{
		Title:     title,
		Page:      params.Page,
		PageSize:  params.PageSize,
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
	}

	// Perform search
	books, total, err := uc.bookRepo.Search(ctx, searchParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search books by title: %w", err)
	}

	// Calculate total pages
	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &SearchResult{
		Books:      books,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

// SearchByAuthor searches books by author
func (uc *searchUseCase) SearchByAuthor(ctx context.Context, author string, params PaginationParams) (*SearchResult, error) {
	if author == "" {
		return nil, errors.New("author is required")
	}

	// Set default pagination values
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}

	// Build list parameters
	listParams := repository.ListParams{
		Page:      params.Page,
		PageSize:  params.PageSize,
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
	}

	// Perform search
	books, total, err := uc.bookRepo.GetByAuthor(ctx, author, listParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search books by author: %w", err)
	}

	// Calculate total pages
	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &SearchResult{
		Books:      books,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

// SearchByISBN searches books by ISBN
func (uc *searchUseCase) SearchByISBN(ctx context.Context, isbn string) (*entity.Book, error) {
	if isbn == "" {
		return nil, errors.New("ISBN is required")
	}

	// Clean ISBN (remove dashes and spaces)
	isbn = strings.ReplaceAll(isbn, "-", "")
	isbn = strings.ReplaceAll(isbn, " ", "")

	book, err := uc.bookRepo.GetByISBN(ctx, isbn)
	if err != nil {
		return nil, fmt.Errorf("failed to search book by ISBN: %w", err)
	}

	return book, nil
}

// SearchByCategory searches books by category
func (uc *searchUseCase) SearchByCategory(ctx context.Context, category string, params PaginationParams) (*SearchResult, error) {
	if category == "" {
		return nil, errors.New("category is required")
	}

	// Set default pagination values
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}

	// Build list parameters
	listParams := repository.ListParams{
		Page:      params.Page,
		PageSize:  params.PageSize,
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
	}

	// Perform search
	books, total, err := uc.bookRepo.GetByCategory(ctx, category, listParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search books by category: %w", err)
	}

	// Calculate total pages
	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &SearchResult{
		Books:      books,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

// AdvancedSearch performs an advanced search with multiple criteria
func (uc *searchUseCase) AdvancedSearch(ctx context.Context, params AdvancedSearchParams) (*SearchResult, error) {
	// Set default pagination values
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}

	// Build search parameters
	searchParams := repository.SearchParams{
		Title:           params.Title,
		Author:          params.Author,
		ISBN:            params.ISBN,
		Category:        params.Category,
		PublicationYear: params.PublicationYear,
		YearFrom:        params.YearFrom,
		YearTo:          params.YearTo,
		Page:            params.Page,
		PageSize:        params.PageSize,
		SortBy:          params.SortBy,
		SortOrder:       params.SortOrder,
	}

	// Perform search
	books, total, err := uc.bookRepo.Search(ctx, searchParams)
	if err != nil {
		return nil, fmt.Errorf("failed to perform advanced search: %w", err)
	}

	// Filter by availability if requested
	if params.OnlyAvailable {
		books, total = uc.filterAvailableBooks(ctx, books)
	}

	// Calculate total pages
	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &SearchResult{
		Books:      books,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetAvailableCategories retrieves all unique book categories
func (uc *searchUseCase) GetAvailableCategories(ctx context.Context) ([]string, error) {
	// Get all books
	listParams := repository.ListParams{
		Page:     1,
		PageSize: 1000, // Large page size to get all categories
	}

	books, _, err := uc.bookRepo.List(ctx, listParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get books: %w", err)
	}

	// Extract unique categories
	categoryMap := make(map[string]bool)
	for _, book := range books {
		if book.Category != "" {
			categoryMap[book.Category] = true
		}
	}

	// Convert map to slice
	categories := make([]string, 0, len(categoryMap))
	for category := range categoryMap {
		categories = append(categories, category)
	}

	return categories, nil
}

// GetPopularBooks retrieves popular books based on borrowing statistics
func (uc *searchUseCase) GetPopularBooks(ctx context.Context, limit int) ([]*entity.Book, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	// Get most borrowed books statistics
	stats, err := uc.borrowRecordRepo.GetMostBorrowedBooks(ctx, limit, time.Time{}, time.Time{})
	if err != nil {
		return nil, fmt.Errorf("failed to get popular books: %w", err)
	}

	// Fetch book details
	books := make([]*entity.Book, 0, len(stats))
	for _, stat := range stats {
		book, err := uc.bookRepo.GetByID(ctx, stat.BookID)
		if err != nil {
			// Skip books that can't be found
			continue
		}
		books = append(books, book)
	}

	return books, nil
}

// SuggestBooks suggests books based on user preferences or history
func (uc *searchUseCase) SuggestBooks(ctx context.Context, userID string, limit int) ([]*entity.Book, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	// Get user's borrowing history
	historyParams := repository.BorrowRecordListParams{
		Page:     1,
		PageSize: 20, // Get recent borrows
		SortBy:   "checkout_date",
		SortOrder: "desc",
	}

	borrowHistory, _, err := uc.borrowRecordRepo.GetUserBorrowingHistory(ctx, userID, historyParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get user borrowing history: %w", err)
	}

	// Extract categories from borrowed books
	categoryMap := make(map[string]int)
	for _, record := range borrowHistory {
		// Get book copy
		bookCopy, err := uc.bookCopyRepo.GetByID(ctx, record.BookCopyID)
		if err != nil {
			continue
		}

		// Get book
		book, err := uc.bookRepo.GetByID(ctx, bookCopy.BookID)
		if err != nil {
			continue
		}

		// Count category frequency
		if book.Category != "" {
			categoryMap[book.Category]++
		}
	}

	// Find most frequent category
	var mostFrequentCategory string
	maxCount := 0
	for category, count := range categoryMap {
		if count > maxCount {
			maxCount = count
			mostFrequentCategory = category
		}
	}

	// If no category found, return recently added books
	if mostFrequentCategory == "" {
		return uc.bookRepo.GetRecentlyAdded(ctx, limit)
	}

	// Get books from the most frequent category
	listParams := repository.ListParams{
		Page:      1,
		PageSize:  limit,
		SortBy:    "created_at",
		SortOrder: "desc",
	}

	books, _, err := uc.bookRepo.GetByCategory(ctx, mostFrequentCategory, listParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested books: %w", err)
	}

	return books, nil
}

// filterAvailableBooks filters books that have available copies
func (uc *searchUseCase) filterAvailableBooks(ctx context.Context, books []*entity.Book) ([]*entity.Book, int64) {
	availableBooks := make([]*entity.Book, 0)

	for _, book := range books {
		// Check if book has available copies
		count, err := uc.bookCopyRepo.CountAvailableByBookID(ctx, book.ID)
		if err != nil {
			continue
		}

		if count > 0 {
			availableBooks = append(availableBooks, book)
		}
	}

	return availableBooks, int64(len(availableBooks))
}