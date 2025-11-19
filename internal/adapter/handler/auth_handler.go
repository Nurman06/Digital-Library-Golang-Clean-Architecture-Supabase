package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/infrastructure/auth"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/infrastructure/logger"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/usecase"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authUseCase     usecase.AuthUseCase
	userUseCase     usecase.UserUseCase
	jwtService      *auth.JWTService
	passwordService *auth.PasswordService
	logger          *logger.Logger
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	authUseCase usecase.AuthUseCase,
	userUseCase usecase.UserUseCase,
	jwtService *auth.JWTService,
	passwordService *auth.PasswordService,
	logger *logger.Logger,
) *AuthHandler {
	return &AuthHandler{
		authUseCase:     authUseCase,
		userUseCase:     userUseCase,
		jwtService:      jwtService,
		passwordService: passwordService,
		logger:          logger,
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
	if req.FullName == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, "Full name is required")
		return
	}

	// Validate password strength
	if err := h.passwordService.ValidatePasswordStrength(req.Password); err != nil {
		ErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}

	// Hash password
	hashedPassword, err := h.passwordService.HashPassword(req.Password)
	if err != nil {
		h.logger.Errorf("Failed to hash password: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to process password")
		return
	}

	// Create user entity
	user := &entity.User{
		Email:          req.Email,
		PasswordHash:   hashedPassword,
		FullName:       req.FullName,
		Role:           entity.RoleMember, // Default role
		Status:         entity.StatusActive,
		BorrowingLimit: entity.GetDefaultBorrowingLimit(entity.RoleMember),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Register user
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

	// Get user by email
	user, err := h.authUseCase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Errorf("Login failed: %v", err)
		ErrorResponse(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Invalid email or password")
		return
	}

	// Verify password
	if err := h.passwordService.ComparePassword(user.PasswordHash, req.Password); err != nil {
		h.logger.Errorf("Password verification failed: %v", err)
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

	// Generate JWT access token
	accessToken, err := h.jwtService.GenerateToken(user)
	if err != nil {
		h.logger.Errorf("Failed to generate access token: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to generate token")
		return
	}

	// Generate refresh token
	refreshToken, err := h.jwtService.GenerateRefreshToken(user)
	if err != nil {
		h.logger.Errorf("Failed to generate refresh token: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to generate token")
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
		Token:        accessToken,
		RefreshToken: refreshToken,
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

	// Validate refresh token and extract user ID
	claims, err := h.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		h.logger.Errorf("Invalid refresh token: %v", err)
		ErrorResponse(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Invalid or expired refresh token")
		return
	}

	// Verify it's a refresh token
	if claims.Issuer != "digital-library-refresh" {
		ErrorResponse(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Not a refresh token")
		return
	}

	// Get user
	user, err := h.userUseCase.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		h.logger.Errorf("Failed to get user: %v", err)
		ErrorResponse(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Invalid refresh token")
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

	// Generate new access token
	accessToken, err := h.jwtService.GenerateToken(user)
	if err != nil {
		h.logger.Errorf("Failed to generate access token: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to generate token")
		return
	}

	// Generate new refresh token
	newRefreshToken, err := h.jwtService.GenerateRefreshToken(user)
	if err != nil {
		h.logger.Errorf("Failed to generate refresh token: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to generate token")
		return
	}

	resp := &RefreshTokenResponse{
		Token:        accessToken,
		RefreshToken: newRefreshToken,
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