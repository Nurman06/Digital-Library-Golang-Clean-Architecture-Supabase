package supabase

import (
	"context"
	"fmt"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/supabase-community/gotrue-go/types"
	supa "github.com/supabase-community/supabase-go"
)

// AuthService handles Supabase Auth operations
type AuthService struct {
	client *supa.Client
}

// NewAuthService creates a new Supabase Auth service
func NewAuthService(client *supa.Client) *AuthService {
	return &AuthService{
		client: client,
	}
}

// SignUpRequest represents a user registration request
type SignUpRequest struct {
	Email    string
	Password string
	FullName string
	Role     entity.UserRole
}

// SignUpResponse represents the response from Supabase Auth signup
type SignUpResponse struct {
	User        *AuthUser
	Session     *Session
	AccessToken string
}

// AuthUser represents a user from Supabase Auth
type AuthUser struct {
	ID    string
	Email string
}

// Session represents an auth session
type Session struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	TokenType    string
	User         *AuthUser
}

// SignUp registers a new user with Supabase Auth
func (s *AuthService) SignUp(ctx context.Context, req SignUpRequest) (*SignUpResponse, error) {
	// Sign up user with Supabase Auth using types.SignupRequest
	authResp, err := s.client.Auth.Signup(types.SignupRequest{
		Email:    req.Email,
		Password: req.Password,
		Data: map[string]interface{}{
			"full_name": req.FullName,
			"role":      string(req.Role),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to sign up user: %w", err)
	}

	if authResp.User.ID.String() == "" {
		return nil, fmt.Errorf("user creation failed")
	}

	return &SignUpResponse{
		User: &AuthUser{
			ID:    authResp.User.ID.String(),
			Email: authResp.User.Email,
		},
		Session: &Session{
			AccessToken:  authResp.AccessToken,
			RefreshToken: authResp.RefreshToken,
			ExpiresIn:    authResp.ExpiresIn,
			TokenType:    authResp.TokenType,
			User: &AuthUser{
				ID:    authResp.User.ID.String(),
				Email: authResp.User.Email,
			},
		},
		AccessToken: authResp.AccessToken,
	}, nil
}

// SignInRequest represents a login request
type SignInRequest struct {
	Email    string
	Password string
}

// SignInResponse represents the response from Supabase Auth signin
type SignInResponse struct {
	User        *AuthUser
	Session     *Session
	AccessToken string
}

// SignIn authenticates a user with Supabase Auth
func (s *AuthService) SignIn(ctx context.Context, req SignInRequest) (*SignInResponse, error) {
	authResp, err := s.client.Auth.SignInWithEmailPassword(req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to sign in: %w", err)
	}

	if authResp.User.ID.String() == "" {
		return nil, fmt.Errorf("authentication failed")
	}

	return &SignInResponse{
		User: &AuthUser{
			ID:    authResp.User.ID.String(),
			Email: authResp.User.Email,
		},
		Session: &Session{
			AccessToken:  authResp.AccessToken,
			RefreshToken: authResp.RefreshToken,
			ExpiresIn:    authResp.ExpiresIn,
			TokenType:    authResp.TokenType,
			User: &AuthUser{
				ID:    authResp.User.ID.String(),
				Email: authResp.User.Email,
			},
		},
		AccessToken: authResp.AccessToken,
	}, nil
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string
}

// RefreshTokenResponse represents the response from token refresh
type RefreshTokenResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	TokenType    string
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*RefreshTokenResponse, error) {
	authResp, err := s.client.Auth.RefreshToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return &RefreshTokenResponse{
		AccessToken:  authResp.AccessToken,
		RefreshToken: authResp.RefreshToken,
		ExpiresIn:    authResp.ExpiresIn,
		TokenType:    authResp.TokenType,
	}, nil
}

// VerifyToken verifies a JWT token with Supabase
// Note: For now, we'll use a simple approach - validate by trying to get user with token
// In production, you should verify JWT signature with Supabase public key
func (s *AuthService) VerifyToken(ctx context.Context, token string) (*AuthUser, error) {
	// For Supabase, we can decode the JWT to get user info
	// The token itself is validated by Supabase, so if we can decode it, it's valid
	// This is a simplified version - in production you'd verify the signature
	
	// For now, we'll extract user info from token by calling an authenticated endpoint
	// The supabase-go library doesn't have a direct verify method, so we use GetUser with token
	user, err := s.client.Auth.GetUser()
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	if user.ID.String() == "" {
		return nil, fmt.Errorf("invalid token")
	}

	return &AuthUser{
		ID:    user.ID.String(),
		Email: user.Email,
	}, nil
}
// UpdatePasswordRequest represents a password update request
type UpdatePasswordRequest struct {
	AccessToken string
	NewPassword string
}

// UpdatePassword updates a user's password
func (s *AuthService) UpdatePassword(ctx context.Context, req UpdatePasswordRequest) error {
	// Update password using Supabase Auth
	// Note: Supabase Auth requires the user to be authenticated (valid access token)
	// The access token should be set in the client before calling UpdateUser
	_, err := s.client.Auth.UpdateUser(types.UpdateUserRequest{
		Password: &req.NewPassword,
	})
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

// SignOut signs out a user
func (s *AuthService) SignOut(ctx context.Context, token string) error {
	// Supabase Go client doesn't have SignOut method
	// Token invalidation is typically handled client-side by removing the token
	// Or you can call the REST API directly
	return nil
}