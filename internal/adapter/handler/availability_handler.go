package handler

import (
	"net/http"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/infrastructure/logger"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/usecase"
	"github.com/gorilla/mux"
)

// AvailabilityHandler handles availability-related HTTP requests
type AvailabilityHandler struct {
	bookUseCase usecase.BookUseCase
	logger      *logger.Logger
}

// NewAvailabilityHandler creates a new AvailabilityHandler
func NewAvailabilityHandler(bookUseCase usecase.BookUseCase, logger *logger.Logger) *AvailabilityHandler {
	return &AvailabilityHandler{
		bookUseCase: bookUseCase,
		logger:      logger,
	}
}

// GetBookAvailability handles getting book availability information
// GET /api/v1/books/{id}/availability
func (h *AvailabilityHandler) GetBookAvailability(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]

	if bookID == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "Book ID is required")
		return
	}

	// Check if book exists
	book, err := h.bookUseCase.GetBookByID(r.Context(), bookID)
	if err != nil {
		h.logger.Errorf("Failed to get book: %v", err)
		ErrorResponse(w, http.StatusNotFound, ErrCodeNotFound, "Book not found")
		return
	}

	// Get all copies
	copies, err := h.bookUseCase.GetBookCopies(r.Context(), bookID)
	if err != nil {
		h.logger.Errorf("Failed to get book copies: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to get book copies")
		return
	}

	// Count copies by status
	var availableCount, borrowedCount, reservedCount int
	copyResponses := make([]*BookCopyResponse, len(copies))
	
	for i, copy := range copies {
		switch copy.Status {
		case entity.CopyStatusAvailable:
			availableCount++
		case entity.CopyStatusBorrowed:
			borrowedCount++
		case entity.CopyStatusReserved:
			reservedCount++
		}

		copyResponses[i] = &BookCopyResponse{
			ID:         copy.ID,
			BookID:     copy.BookID,
			CopyNumber: copy.CopyNumber,
			Status:     string(copy.Status),
			Location:   copy.Location,
			CreatedAt:  copy.CreatedAt,
			UpdatedAt:  copy.UpdatedAt,
		}
	}

	// Create availability response
	availabilityResp := &AvailabilityResponse{
		BookID:          book.ID,
		TotalCopies:     len(copies),
		AvailableCopies: availableCount,
		BorrowedCopies:  borrowedCount,
		ReservedCopies:  reservedCount,
		Copies:          copyResponses,
	}

	SuccessResponse(w, http.StatusOK, availabilityResp)
}

// GetAvailableBookCopies handles getting available copies of a book
// GET /api/v1/books/{id}/copies/available
func (h *AvailabilityHandler) GetAvailableBookCopies(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]

	if bookID == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "Book ID is required")
		return
	}

	// Get available copies
	copies, err := h.bookUseCase.GetAvailableBookCopies(r.Context(), bookID)
	if err != nil {
		h.logger.Errorf("Failed to get available book copies: %v", err)
		ErrorResponse(w, http.StatusNotFound, ErrCodeNotFound, "Book not found")
		return
	}

	// Convert to response
	copyResponses := make([]*BookCopyResponse, len(copies))
	for i, copy := range copies {
		copyResponses[i] = &BookCopyResponse{
			ID:         copy.ID,
			BookID:     copy.BookID,
			CopyNumber: copy.CopyNumber,
			Status:     string(copy.Status),
			Location:   copy.Location,
			CreatedAt:  copy.CreatedAt,
			UpdatedAt:  copy.UpdatedAt,
		}
	}

	SuccessResponse(w, http.StatusOK, copyResponses)
}

// GetBookCopies handles getting all copies of a book
// GET /api/v1/books/{id}/copies
func (h *AvailabilityHandler) GetBookCopies(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]

	if bookID == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "Book ID is required")
		return
	}

	// Get all copies
	copies, err := h.bookUseCase.GetBookCopies(r.Context(), bookID)
	if err != nil {
		h.logger.Errorf("Failed to get book copies: %v", err)
		ErrorResponse(w, http.StatusNotFound, ErrCodeNotFound, "Book not found")
		return
	}

	// Convert to response
	copyResponses := make([]*BookCopyResponse, len(copies))
	for i, copy := range copies {
		copyResponses[i] = &BookCopyResponse{
			ID:         copy.ID,
			BookID:     copy.BookID,
			CopyNumber: copy.CopyNumber,
			Status:     string(copy.Status),
			Location:   copy.Location,
			CreatedAt:  copy.CreatedAt,
			UpdatedAt:  copy.UpdatedAt,
		}
	}

	SuccessResponse(w, http.StatusOK, copyResponses)
}