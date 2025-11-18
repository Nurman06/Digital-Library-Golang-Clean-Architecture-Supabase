package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/infrastructure/logger"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/repository"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/usecase"
	"github.com/gorilla/mux"
)

// BookHandler handles book-related HTTP requests
type BookHandler struct {
	bookUseCase usecase.BookUseCase
	authUseCase usecase.AuthUseCase
	logger      *logger.Logger
}

// NewBookHandler creates a new BookHandler
func NewBookHandler(bookUseCase usecase.BookUseCase, authUseCase usecase.AuthUseCase, logger *logger.Logger) *BookHandler {
	return &BookHandler{
		bookUseCase: bookUseCase,
		authUseCase: authUseCase,
		logger:      logger,
	}
}

// CreateBook handles book creation
// POST /api/v1/books
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body")
		return
	}

	// Create book entity
	book := &entity.Book{
		Title:           req.Title,
		Author:          req.Author,
		ISBN:            req.ISBN,
		Category:        req.Category,
		PublicationYear: req.PublicationYear,
		Description:     req.Description,
	}

	// Create book
	if err := h.bookUseCase.CreateBook(r.Context(), book); err != nil {
		if err.Error() == "book with this ISBN already exists" {
			ErrorResponse(w, http.StatusConflict, ErrCodeConflict, "Book with this ISBN already exists")
			return
		}
		h.logger.Errorf("Failed to create book: %v", err)
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}

	// Convert to response
	bookResp := &BookResponse{
		ID:              book.ID,
		Title:           book.Title,
		Author:          book.Author,
		ISBN:            book.ISBN,
		Category:        book.Category,
		PublicationYear: book.PublicationYear,
		Description:     book.Description,
		CreatedAt:       book.CreatedAt,
		UpdatedAt:       book.UpdatedAt,
	}

	SuccessResponse(w, http.StatusCreated, bookResp)
}

// GetBook handles getting a book by ID
// GET /api/v1/books/{id}
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]

	if bookID == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "Book ID is required")
		return
	}

	// Get book
	book, err := h.bookUseCase.GetBookByID(r.Context(), bookID)
	if err != nil {
		h.logger.Errorf("Failed to get book: %v", err)
		ErrorResponse(w, http.StatusNotFound, ErrCodeNotFound, "Book not found")
		return
	}

	// Get availability count
	availableCount, err := h.bookUseCase.GetBookAvailabilityCount(r.Context(), bookID)
	if err != nil {
		h.logger.Errorf("Failed to get availability count: %v", err)
		availableCount = 0
	}

	// Get total copies
	copies, err := h.bookUseCase.GetBookCopies(r.Context(), bookID)
	if err != nil {
		h.logger.Errorf("Failed to get book copies: %v", err)
	}

	// Convert to response with availability
	bookResp := &BookWithAvailabilityResponse{
		BookResponse: BookResponse{
			ID:              book.ID,
			Title:           book.Title,
			Author:          book.Author,
			ISBN:            book.ISBN,
			Category:        book.Category,
			PublicationYear: book.PublicationYear,
			Description:     book.Description,
			CreatedAt:       book.CreatedAt,
			UpdatedAt:       book.UpdatedAt,
			DeletedAt:       book.DeletedAt,
		},
		TotalCopies:     len(copies),
		AvailableCopies: int(availableCount),
	}

	SuccessResponse(w, http.StatusOK, bookResp)
}

// ListBooks handles listing books with pagination
// GET /api/v1/books
func (h *BookHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	sortBy := r.URL.Query().Get("sort_by")
	sortDir := r.URL.Query().Get("sort_dir")

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortDir == "" {
		sortDir = "desc"
	}

	params := repository.ListParams{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortDir,
	}

	// List books
	books, total, err := h.bookUseCase.ListBooks(r.Context(), params)
	if err != nil {
		h.logger.Errorf("Failed to list books: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to list books")
		return
	}

	// Convert to response
	bookResponses := make([]*BookResponse, len(books))
	for i, book := range books {
		bookResponses[i] = &BookResponse{
			ID:              book.ID,
			Title:           book.Title,
			Author:          book.Author,
			ISBN:            book.ISBN,
			Category:        book.Category,
			PublicationYear: book.PublicationYear,
			Description:     book.Description,
			CreatedAt:       book.CreatedAt,
			UpdatedAt:       book.UpdatedAt,
			DeletedAt:       book.DeletedAt,
		}
	}

	// Create pagination metadata
	meta := &MetaData{
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: CalculateTotalPages(total, pageSize),
	}

	SuccessResponseWithMeta(w, http.StatusOK, bookResponses, meta)
}

// UpdateBook handles updating a book
// PUT /api/v1/books/{id}
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]

	if bookID == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "Book ID is required")
		return
	}

	var req UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body")
		return
	}

	// Get existing book
	book, err := h.bookUseCase.GetBookByID(r.Context(), bookID)
	if err != nil {
		h.logger.Errorf("Failed to get book: %v", err)
		ErrorResponse(w, http.StatusNotFound, ErrCodeNotFound, "Book not found")
		return
	}

	// Update book fields
	book.Title = req.Title
	book.Author = req.Author
	book.ISBN = req.ISBN
	book.Category = req.Category
	book.PublicationYear = req.PublicationYear
	book.Description = req.Description
	book.UpdatedAt = time.Now()

	// Update book
	if err := h.bookUseCase.UpdateBook(r.Context(), book); err != nil {
		if err.Error() == "book with this ISBN already exists" {
			ErrorResponse(w, http.StatusConflict, ErrCodeConflict, "Book with this ISBN already exists")
			return
		}
		h.logger.Errorf("Failed to update book: %v", err)
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}

	// Convert to response
	bookResp := &BookResponse{
		ID:              book.ID,
		Title:           book.Title,
		Author:          book.Author,
		ISBN:            book.ISBN,
		Category:        book.Category,
		PublicationYear: book.PublicationYear,
		Description:     book.Description,
		CreatedAt:       book.CreatedAt,
		UpdatedAt:       book.UpdatedAt,
	}

	SuccessResponse(w, http.StatusOK, bookResp)
}

// DeleteBook handles deleting a book
// DELETE /api/v1/books/{id}
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]

	if bookID == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "Book ID is required")
		return
	}

	// Delete book
	if err := h.bookUseCase.DeleteBook(r.Context(), bookID); err != nil {
		if err.Error() == "cannot delete book with active borrows" {
			ErrorResponse(w, http.StatusConflict, ErrCodeConflict, "Cannot delete book with active borrows")
			return
		}
		h.logger.Errorf("Failed to delete book: %v", err)
		ErrorResponse(w, http.StatusNotFound, ErrCodeNotFound, "Book not found")
		return
	}

	SuccessResponse(w, http.StatusNoContent, nil)
}

// SearchBooks handles book search
// GET /api/v1/books/search
func (h *BookHandler) SearchBooks(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	author := r.URL.Query().Get("author")
	isbn := r.URL.Query().Get("isbn")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	sortBy := r.URL.Query().Get("sort_by")
	sortDir := r.URL.Query().Get("sort_dir")

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortDir == "" {
		sortDir = "desc"
	}

	params := repository.SearchParams{
		Query:    query,
		Category: category,
		Author:   author,
		ISBN:     isbn,
		Page:     page,
		PageSize: pageSize,
		SortBy:   sortBy,
		SortOrder: sortDir,
	}

	// Search books
	books, total, err := h.bookUseCase.SearchBooks(r.Context(), params)
	if err != nil {
		h.logger.Errorf("Failed to search books: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to search books")
		return
	}

	// Convert to response
	bookResponses := make([]*BookResponse, len(books))
	for i, book := range books {
		bookResponses[i] = &BookResponse{
			ID:              book.ID,
			Title:           book.Title,
			Author:          book.Author,
			ISBN:            book.ISBN,
			Category:        book.Category,
			PublicationYear: book.PublicationYear,
			Description:     book.Description,
			CreatedAt:       book.CreatedAt,
			UpdatedAt:       book.UpdatedAt,
			DeletedAt:       book.DeletedAt,
		}
	}

	// Create pagination metadata
	meta := &MetaData{
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: CalculateTotalPages(total, pageSize),
	}

	SuccessResponseWithMeta(w, http.StatusOK, bookResponses, meta)
}