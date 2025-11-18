package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/infrastructure/logger"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/repository"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/usecase"
	"github.com/gorilla/mux"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userUseCase usecase.UserUseCase
	authUseCase usecase.AuthUseCase
	logger      *logger.Logger
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userUseCase usecase.UserUseCase, authUseCase usecase.AuthUseCase, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		authUseCase: authUseCase,
		logger:      logger,
	}
}

// GetUser handles getting a user by ID
// GET /api/v1/users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if userID == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "User ID is required")
		return
	}

	// Get user
	user, err := h.userUseCase.GetUserByID(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("Failed to get user: %v", err)
		ErrorResponse(w, http.StatusNotFound, ErrCodeNotFound, "User not found")
		return
	}

	// Convert to response
	userResp := &UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		FullName:       user.FullName,
		Role:           string(user.Role),
		Status:         string(user.Status),
		BorrowingLimit: user.BorrowingLimit,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}

	SuccessResponse(w, http.StatusOK, userResp)
}

// ListUsers handles listing users with pagination (Admin only)
// GET /api/v1/users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	params := repository.UserListParams{
		Page:     page,
		PageSize: pageSize,
	}

	// List users
	users, total, err := h.userUseCase.ListUsers(r.Context(), params)
	if err != nil {
		h.logger.Errorf("Failed to list users: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to list users")
		return
	}

	// Convert to response
	userResponses := make([]*UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = &UserResponse{
			ID:             user.ID,
			Email:          user.Email,
			FullName:       user.FullName,
			Role:           string(user.Role),
			Status:         string(user.Status),
			BorrowingLimit: user.BorrowingLimit,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
		}
	}

	// Create pagination metadata
	meta := &MetaData{
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: CalculateTotalPages(total, pageSize),
	}

	SuccessResponseWithMeta(w, http.StatusOK, userResponses, meta)
}

// UpdateUserProfile handles updating the current user's profile
// PUT /api/v1/users/profile
func (h *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, ErrCodeUnauthorized, "User not authenticated")
		return
	}

	var req UpdateUserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body")
		return
	}

	// Get existing user
	user, err := h.userUseCase.GetUserByID(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("Failed to get user: %v", err)
		ErrorResponse(w, http.StatusNotFound, ErrCodeNotFound, "User not found")
		return
	}

	// Update user fields
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	// Update user
	if err := h.userUseCase.UpdateUser(r.Context(), user); err != nil {
		if err.Error() == "user with this email already exists" {
			ErrorResponse(w, http.StatusConflict, ErrCodeConflict, "Email already in use")
			return
		}
		h.logger.Errorf("Failed to update user: %v", err)
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}

	// Convert to response
	userResp := &UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		FullName:       user.FullName,
		Role:           string(user.Role),
		Status:         string(user.Status),
		BorrowingLimit: user.BorrowingLimit,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}

	SuccessResponse(w, http.StatusOK, userResp)
}

// UpdateUser handles updating a user (Admin only)
// PUT /api/v1/users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if userID == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "User ID is required")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body")
		return
	}

	// Get existing user
	user, err := h.userUseCase.GetUserByID(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("Failed to get user: %v", err)
		ErrorResponse(w, http.StatusNotFound, ErrCodeNotFound, "User not found")
		return
	}

	// Update user fields
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = entity.UserRole(req.Role)
	}
	if req.Status != "" {
		user.Status = entity.UserStatus(req.Status)
	}
	if req.BorrowingLimit >= 0 {
		user.BorrowingLimit = req.BorrowingLimit
	}

	// Update user
	if err := h.userUseCase.UpdateUser(r.Context(), user); err != nil {
		if err.Error() == "user with this email already exists" {
			ErrorResponse(w, http.StatusConflict, ErrCodeConflict, "Email already in use")
			return
		}
		h.logger.Errorf("Failed to update user: %v", err)
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}

	// Convert to response
	userResp := &UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		FullName:       user.FullName,
		Role:           string(user.Role),
		Status:         string(user.Status),
		BorrowingLimit: user.BorrowingLimit,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}

	SuccessResponse(w, http.StatusOK, userResp)
}

// DeleteUser handles deleting a user (Admin only)
// DELETE /api/v1/users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if userID == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "User ID is required")
		return
	}

	// Delete user
	if err := h.userUseCase.DeleteUser(r.Context(), userID); err != nil {
		if err.Error() == "cannot delete user with active borrow records" {
			ErrorResponse(w, http.StatusConflict, ErrCodeConflict, "Cannot delete user with active borrow records")
			return
		}
		h.logger.Errorf("Failed to delete user: %v", err)
		ErrorResponse(w, http.StatusNotFound, ErrCodeNotFound, "User not found")
		return
	}

	SuccessResponse(w, http.StatusNoContent, nil)
}

// SuspendUser handles suspending a user (Admin only)
// POST /api/v1/users/{id}/suspend
func (h *UserHandler) SuspendUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if userID == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "User ID is required")
		return
	}

	// Suspend user
	if err := h.userUseCase.SuspendUser(r.Context(), userID); err != nil {
		h.logger.Errorf("Failed to suspend user: %v", err)
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, err.Error())
		return
	}

	SuccessResponse(w, http.StatusOK, map[string]string{"message": "User suspended successfully"})
}

// ActivateUser handles activating a user (Admin only)
// POST /api/v1/users/{id}/activate
func (h *UserHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if userID == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "User ID is required")
		return
	}

	// Activate user
	if err := h.userUseCase.ActivateUser(r.Context(), userID); err != nil {
		h.logger.Errorf("Failed to activate user: %v", err)
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, err.Error())
		return
	}

	SuccessResponse(w, http.StatusOK, map[string]string{"message": "User activated successfully"})
}