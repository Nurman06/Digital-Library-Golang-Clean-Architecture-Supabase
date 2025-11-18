package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/infrastructure/logger"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/usecase"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authUseCase usecase.AuthUseCase
	userUseCase usecase.UserUseCase
	logger      *logger.Logger
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authUseCase usecase.AuthUseCase, userUseCase usecase.UserUseCase, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		userUseCase: userUseCase,
		logger:      logger,
	}
}

// Register handles user registration
// POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.Email == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, "Email is required")
		return
	}
	if req.Password == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, "Password is required")
		return
	}
	if len(req.Password) < 8 {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, "Password must be at least 8 characters")
		return
	}
	if req.FullName == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, "Full name is required")
		return
	}

	// Create user entity
	user := &entity.User{
		Email:          req.Email,
		FullName:       req.FullName,
		Role:           entity.RoleMember, // Default role
		Status:         entity.StatusActive,
		BorrowingLimit: entity.GetDefaultBorrowingLimit(entity.RoleMember),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Register user
	// Note: In production, this should integrate with Supabase Auth for password handling
	if err := h.userUseCase.RegisterUser(r.Context(), user); err != nil {
		if err.Error() == "user with this email already exists" {
			ErrorResponse(w, http.StatusConflict, ErrCodeConflict, "Email already registered")
			return
		}
		h.logger.Errorf("Failed to register user: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to register user")
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

	SuccessResponse(w, http.StatusCreated, userResp)
}

// Login handles user login
// POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.Email == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, "Email is required")
		return
	}
	if req.Password == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, "Password is required")
		return
	}

	// Authenticate user
	// Note: In production, this should integrate with Supabase Auth for password verification
	user, err := h.authUseCase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Errorf("Login failed: %v", err)
		ErrorResponse(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Invalid email or password")
		return
	}

	// Check user status
	if user.IsSuspended() {
		ErrorResponse(w, http.StatusForbidden, ErrCodeForbidden, "Account is suspended")
		return
	}
	if user.IsExpired() {
		ErrorResponse(w, http.StatusForbidden, ErrCodeForbidden, "Account has expired")
		return
	}

	// Generate JWT token
	// Note: In production, this should use Supabase Auth to generate tokens
	token := "placeholder_jwt_token"

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

	loginResp := &LoginResponse{
		User:  userResp,
		Token: token,
	}

	SuccessResponse(w, http.StatusOK, loginResp)
}

// GetProfile handles getting the current user's profile
// GET /api/v1/auth/profile
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, ErrCodeUnauthorized, "User not authenticated")
		return
	}

	// Get user
	user, err := h.userUseCase.GetUserByID(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("Failed to get user profile: %v", err)
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