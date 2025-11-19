package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/infrastructure/logger"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/infrastructure/supabase"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/usecase"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authUseCase usecase.AuthUseCase
	userUseCase usecase.UserUseCase
	authService *supabase.AuthService
	logger      *logger.Logger
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	authUseCase usecase.AuthUseCase,
	userUseCase usecase.UserUseCase,
	authService *supabase.AuthService,
	logger *logger.Logger,
) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		userUseCase: userUseCase,
		authService: authService,
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

	// Register user with Supabase Auth
	authResp, err := h.authService.SignUp(r.Context(), supabase.SignUpRequest{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Role:     entity.RoleMember, // Default role
	})
	if err != nil {
		h.logger.Errorf("Failed to register user with Supabase Auth: %v", err)
		if err.Error() == "user with this email already exists" ||
		   err.Error() == "User already registered" {
			ErrorResponse(w, http.StatusConflict, ErrCodeConflict, "Email already registered")
			return
		}
		ErrorResponse(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to register user")
		return
	}

	// Create user record in our database
	user := &entity.User{
		ID:             authResp.User.ID, // Use Supabase Auth user ID
		Email:          authResp.User.Email,
		FullName:       req.FullName,
		Role:           entity.RoleMember,
		Status:         entity.StatusActive,
		BorrowingLimit: entity.GetDefaultBorrowingLimit(entity.RoleMember),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Save user to database
	if err := h.userUseCase.RegisterUser(r.Context(), user); err != nil {
		h.logger.Errorf("Failed to save user to database: %v", err)
		// User is created in Supabase Auth but not in our DB
		// This is acceptable as auth.users is the source of truth
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

	loginResp := &LoginResponse{
		User:         userResp,
		Token:        authResp.AccessToken,
		RefreshToken: authResp.Session.RefreshToken,
	}

	SuccessResponse(w, http.StatusCreated, loginResp)
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

	// Authenticate with Supabase Auth
	authResp, err := h.authService.SignIn(r.Context(), supabase.SignInRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		h.logger.Errorf("Login failed: %v", err)
		ErrorResponse(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Invalid email or password")
		return
	}

	// Get user from database to check status and get role
	user, err := h.userUseCase.GetUserByID(r.Context(), authResp.User.ID)
	if err != nil {
		h.logger.Errorf("Failed to get user: %v", err)
		// User authenticated but not in our DB, this shouldn't happen
		ErrorResponse(w, http.StatusInternalServerError, ErrCodeInternal, "User record not found")
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
		User:         userResp,
		Token:        authResp.AccessToken,
		RefreshToken: authResp.Session.RefreshToken,
	}

	SuccessResponse(w, http.StatusOK, loginResp)
}

// RefreshToken handles token refresh
// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.RefreshToken == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, "Refresh token is required")
		return
	}

	// Refresh token with Supabase Auth
	authResp, err := h.authService.RefreshToken(r.Context(), supabase.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		h.logger.Errorf("Failed to refresh token: %v", err)
		ErrorResponse(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Invalid or expired refresh token")
		return
	}

	resp := &RefreshTokenResponse{
		Token:        authResp.AccessToken,
		RefreshToken: authResp.RefreshToken,
	}

	SuccessResponse(w, http.StatusOK, resp)
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

// ChangePassword handles password change for authenticated user
// POST /api/v1/auth/change-password
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.CurrentPassword == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, "Current password is required")
		return
	}
	if req.NewPassword == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, "New password is required")
		return
	}
	if len(req.NewPassword) < 8 {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, "New password must be at least 8 characters")
		return
	}

	// Get user ID from context
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, ErrCodeUnauthorized, "User not authenticated")
		return
	}

	// Get user to verify current password
	user, err := h.userUseCase.GetUserByID(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("Failed to get user: %v", err)
		ErrorResponse(w, http.StatusNotFound, ErrCodeNotFound, "User not found")
		return
	}

	// Verify current password by attempting to sign in
	_, err = h.authService.SignIn(r.Context(), supabase.SignInRequest{
		Email:    user.Email,
		Password: req.CurrentPassword,
	})
	if err != nil {
		h.logger.Errorf("Current password verification failed: %v", err)
		ErrorResponse(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Current password is incorrect")
		return
	}

	// Get access token from Authorization header
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		ErrorResponse(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Invalid authorization header")
		return
	}
	accessToken := parts[1]

	// Update password via Supabase Auth
	err = h.authService.UpdatePassword(r.Context(), supabase.UpdatePasswordRequest{
		AccessToken: accessToken,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		h.logger.Errorf("Failed to update password: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to update password")
		return
	}

	SuccessResponse(w, http.StatusOK, map[string]string{
		"message": "Password updated successfully",
	})
}