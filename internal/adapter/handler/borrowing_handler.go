package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/infrastructure/logger"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/repository"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/usecase"
	"github.com/gorilla/mux"
)

// BorrowingHandler handles borrowing-related HTTP requests
type BorrowingHandler struct {
	borrowingUseCase usecase.BorrowingUseCase
	bookUseCase      usecase.BookUseCase
	userUseCase      usecase.UserUseCase
	logger           *logger.Logger
}

// NewBorrowingHandler creates a new BorrowingHandler
func NewBorrowingHandler(
	borrowingUseCase usecase.BorrowingUseCase,
	bookUseCase usecase.BookUseCase,
	userUseCase usecase.UserUseCase,
	logger *logger.Logger,
) *BorrowingHandler {
	return &BorrowingHandler{
		borrowingUseCase: borrowingUseCase,
		bookUseCase:      bookUseCase,
		userUseCase:      userUseCase,
		logger:           logger,
	}
}

// CheckoutBook handles book checkout
// POST /api/v1/borrowing/checkout
func (h *BorrowingHandler) CheckoutBook(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, ErrCodeUnauthorized, "User not authenticated")
		return
	}

	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body")
		return
	}

	if req.BookID == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, "Book ID is required")
		return
	}

	// Checkout book
	borrowRecord, err := h.borrowingUseCase.CheckoutBook(r.Context(), userID, req.BookID)
	if err != nil {
		h.logger.Errorf("Failed to checkout book: %v", err)
		
		// Handle specific errors
		switch err.Error() {
		case "user is not allowed to borrow books":
			ErrorResponse(w, http.StatusForbidden, ErrCodeForbidden, "User is not allowed to borrow books")
		case "no available copies of this book":
			ErrorResponse(w, http.StatusConflict, ErrCodeUnavailable, "No available copies of this book")
		case "user has overdue books and cannot borrow more":
			ErrorResponse(w, http.StatusForbidden, ErrCodeForbidden, "User has overdue books and cannot borrow more")
		default:
			if len(err.Error()) > 30 && err.Error()[:30] == "user has reached borrowing lim" {
				ErrorResponse(w, http.StatusConflict, ErrCodeConflict, err.Error())
			} else if len(err.Error()) > 25 && err.Error()[:25] == "user has unpaid late fees" {
				ErrorResponse(w, http.StatusForbidden, ErrCodeForbidden, err.Error())
			} else {
				ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, err.Error())
			}
		}
		return
	}

	// Convert to response
	borrowResp := &BorrowRecordResponse{
		ID:           borrowRecord.ID,
		UserID:       borrowRecord.UserID,
		BookCopyID:   borrowRecord.BookCopyID,
		CheckoutDate: borrowRecord.CheckoutDate,
		DueDate:      borrowRecord.DueDate,
		ReturnDate:   borrowRecord.ReturnDate,
		RenewalCount: borrowRecord.RenewalCount,
		LateFee:      borrowRecord.LateFee,
		Status:       string(borrowRecord.Status),
		CreatedAt:    borrowRecord.CreatedAt,
		UpdatedAt:    borrowRecord.UpdatedAt,
	}

	SuccessResponse(w, http.StatusCreated, borrowResp)
}

// ReturnBook handles book return
// POST /api/v1/borrowing/return
func (h *BorrowingHandler) ReturnBook(w http.ResponseWriter, r *http.Request) {
	var req ReturnBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body")
		return
	}

	if req.BorrowRecordID == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, "Borrow record ID is required")
		return
	}

	// Return book
	if err := h.borrowingUseCase.ReturnBook(r.Context(), req.BorrowRecordID); err != nil {
		h.logger.Errorf("Failed to return book: %v", err)
		
		if err.Error() == "book has already been returned" {
			ErrorResponse(w, http.StatusConflict, ErrCodeConflict, "Book has already been returned")
		} else {
			ErrorResponse(w, http.StatusNotFound, ErrCodeNotFound, "Borrow record not found")
		}
		return
	}

	SuccessResponse(w, http.StatusOK, map[string]string{"message": "Book returned successfully"})
}

// RenewBook handles book renewal
// POST /api/v1/borrowing/renew
func (h *BorrowingHandler) RenewBook(w http.ResponseWriter, r *http.Request) {
	var req RenewBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body")
		return
	}

	if req.BorrowRecordID == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, "Borrow record ID is required")
		return
	}

	// Renew book
	if err := h.borrowingUseCase.RenewBook(r.Context(), req.BorrowRecordID); err != nil {
		h.logger.Errorf("Failed to renew book: %v", err)
		
		switch err.Error() {
		case "cannot renew overdue book":
			ErrorResponse(w, http.StatusForbidden, ErrCodeForbidden, "Cannot renew overdue book")
		case "book cannot be renewed":
			ErrorResponse(w, http.StatusConflict, ErrCodeConflict, "Book cannot be renewed")
		default:
			if len(err.Error()) > 20 && err.Error()[:20] == "maximum renewal limi" {
				ErrorResponse(w, http.StatusConflict, ErrCodeConflict, err.Error())
			} else {
				ErrorResponse(w, http.StatusNotFound, ErrCodeNotFound, "Borrow record not found")
			}
		}
		return
	}

	SuccessResponse(w, http.StatusOK, map[string]string{"message": "Book renewed successfully"})
}

// GetBorrowingHistory handles getting user's borrowing history
// GET /api/v1/borrowing/history
func (h *BorrowingHandler) GetBorrowingHistory(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, ErrCodeUnauthorized, "User not authenticated")
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	params := repository.BorrowRecordListParams{
		Page:     page,
		PageSize: pageSize,
	}

	// Get borrowing history
	records, total, err := h.borrowingUseCase.GetUserBorrowingHistory(r.Context(), userID, params)
	if err != nil {
		h.logger.Errorf("Failed to get borrowing history: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to get borrowing history")
		return
	}

	// Convert to response
	recordResponses := make([]*BorrowRecordResponse, len(records))
	for i, record := range records {
		recordResponses[i] = &BorrowRecordResponse{
			ID:           record.ID,
			UserID:       record.UserID,
			BookCopyID:   record.BookCopyID,
			CheckoutDate: record.CheckoutDate,
			DueDate:      record.DueDate,
			ReturnDate:   record.ReturnDate,
			RenewalCount: record.RenewalCount,
			LateFee:      record.LateFee,
			Status:       string(record.Status),
			CreatedAt:    record.CreatedAt,
			UpdatedAt:    record.UpdatedAt,
		}
	}

	// Create pagination metadata
	meta := &MetaData{
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: CalculateTotalPages(total, pageSize),
	}

	SuccessResponseWithMeta(w, http.StatusOK, recordResponses, meta)
}

// GetActiveBorrows handles getting user's active borrows
// GET /api/v1/borrowing/active
func (h *BorrowingHandler) GetActiveBorrows(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, ErrCodeUnauthorized, "User not authenticated")
		return
	}

	// Get active borrows
	records, err := h.borrowingUseCase.GetActiveBorrows(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("Failed to get active borrows: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to get active borrows")
		return
	}

	// Convert to response
	recordResponses := make([]*BorrowRecordResponse, len(records))
	for i, record := range records {
		recordResponses[i] = &BorrowRecordResponse{
			ID:           record.ID,
			UserID:       record.UserID,
			BookCopyID:   record.BookCopyID,
			CheckoutDate: record.CheckoutDate,
			DueDate:      record.DueDate,
			ReturnDate:   record.ReturnDate,
			RenewalCount: record.RenewalCount,
			LateFee:      record.LateFee,
			Status:       string(record.Status),
			CreatedAt:    record.CreatedAt,
			UpdatedAt:    record.UpdatedAt,
		}
	}

	SuccessResponse(w, http.StatusOK, recordResponses)
}

// GetOverdueBorrows handles getting user's overdue borrows
// GET /api/v1/borrowing/overdue
func (h *BorrowingHandler) GetOverdueBorrows(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, ErrCodeUnauthorized, "User not authenticated")
		return
	}

	// Get overdue borrows
	records, err := h.borrowingUseCase.GetOverdueBorrows(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("Failed to get overdue borrows: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to get overdue borrows")
		return
	}

	// Convert to response
	recordResponses := make([]*BorrowRecordResponse, len(records))
	for i, record := range records {
		recordResponses[i] = &BorrowRecordResponse{
			ID:           record.ID,
			UserID:       record.UserID,
			BookCopyID:   record.BookCopyID,
			CheckoutDate: record.CheckoutDate,
			DueDate:      record.DueDate,
			ReturnDate:   record.ReturnDate,
			RenewalCount: record.RenewalCount,
			LateFee:      record.LateFee,
			Status:       string(record.Status),
			CreatedAt:    record.CreatedAt,
			UpdatedAt:    record.UpdatedAt,
		}
	}

	SuccessResponse(w, http.StatusOK, recordResponses)
}

// GetBorrowRecord handles getting a specific borrow record
// GET /api/v1/borrowing/{id}
func (h *BorrowingHandler) GetBorrowRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recordID := vars["id"]

	if recordID == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "Borrow record ID is required")
		return
	}

	// Get borrow record
	record, err := h.borrowingUseCase.GetBorrowRecord(r.Context(), recordID)
	if err != nil {
		h.logger.Errorf("Failed to get borrow record: %v", err)
		ErrorResponse(w, http.StatusNotFound, ErrCodeNotFound, "Borrow record not found")
		return
	}

	// Convert to response
	recordResp := &BorrowRecordResponse{
		ID:           record.ID,
		UserID:       record.UserID,
		BookCopyID:   record.BookCopyID,
		CheckoutDate: record.CheckoutDate,
		DueDate:      record.DueDate,
		ReturnDate:   record.ReturnDate,
		RenewalCount: record.RenewalCount,
		LateFee:      record.LateFee,
		Status:       string(record.Status),
		CreatedAt:    record.CreatedAt,
		UpdatedAt:    record.UpdatedAt,
	}

	SuccessResponse(w, http.StatusOK, recordResp)
}