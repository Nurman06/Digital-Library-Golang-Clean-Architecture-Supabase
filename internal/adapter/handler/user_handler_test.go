package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/infrastructure/logger"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/repository"
	"github.com/gorilla/mux"
)

func TestUserHandler_GetUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupMocks     func(*MockUserUseCaseForHandler)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "successful user retrieval",
			userID: "user-123",
			setupMocks: func(m *MockUserUseCaseForHandler) {
				m.GetUserByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return &entity.User{
						ID:             id,
						Email:          "user@example.com",
						FullName:       "Test User",
						Role:           entity.RoleMember,
						Status:         entity.StatusActive,
						BorrowingLimit: 5,
						CreatedAt:      time.Now(),
						UpdatedAt:      time.Now(),
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:   "user not found",
			userID: "nonexistent",
			setupMocks: func(m *MockUserUseCaseForHandler) {
				m.GetUserByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return nil, errors.New("user not found")
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserUseCase := &MockUserUseCaseForHandler{}
			mockAuthUseCase := &MockAuthUseCase{}
			tt.setupMocks(mockUserUseCase)

			handler := NewUserHandler(mockUserUseCase, mockAuthUseCase, logger.New())

			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+tt.userID, nil)
			w := httptest.NewRecorder()
			req = mux.SetURLVars(req, map[string]string{"id": tt.userID})

			handler.GetUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response Response
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if tt.expectedError && response.Success {
				t.Error("Expected error response, got success")
			}
			if !tt.expectedError && !response.Success {
				t.Error("Expected success response, got error")
			}
		})
	}
}

func TestUserHandler_ListUsers(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		setupMocks     func(*MockUserUseCaseForHandler)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "successful user listing",
			queryParams: map[string]string{
				"page":      "1",
				"page_size": "20",
			},
			setupMocks: func(m *MockUserUseCaseForHandler) {
				m.ListUsersFunc = func(ctx context.Context, params repository.UserListParams) ([]*entity.User, int64, error) {
					users := []*entity.User{
						{
							ID:             "user-1",
							Email:          "user1@example.com",
							FullName:       "User One",
							Role:           entity.RoleMember,
							Status:         entity.StatusActive,
							BorrowingLimit: 5,
							CreatedAt:      time.Now(),
							UpdatedAt:      time.Now(),
						},
					}
					return users, 1, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserUseCase := &MockUserUseCaseForHandler{}
			mockAuthUseCase := &MockAuthUseCase{}
			tt.setupMocks(mockUserUseCase)

			handler := NewUserHandler(mockUserUseCase, mockAuthUseCase, logger.New())

			url := "/api/v1/users"
			if len(tt.queryParams) > 0 {
				url += "?"
				first := true
				for k, v := range tt.queryParams {
					if !first {
						url += "&"
					}
					url += k + "=" + v
					first = false
				}
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			handler.ListUsers(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupMocks     func(*MockUserUseCaseForHandler)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "successful user deletion",
			userID: "user-123",
			setupMocks: func(m *MockUserUseCaseForHandler) {
				m.DeleteUserFunc = func(ctx context.Context, id string) error {
					return nil
				}
			},
			expectedStatus: http.StatusNoContent,
			expectedError:  false,
		},
		{
			name:   "cannot delete user with active borrows",
			userID: "user-123",
			setupMocks: func(m *MockUserUseCaseForHandler) {
				m.DeleteUserFunc = func(ctx context.Context, id string) error {
					return errors.New("cannot delete user with active borrow records")
				}
			},
			expectedStatus: http.StatusConflict,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserUseCase := &MockUserUseCaseForHandler{}
			mockAuthUseCase := &MockAuthUseCase{}
			tt.setupMocks(mockUserUseCase)

			handler := NewUserHandler(mockUserUseCase, mockAuthUseCase, logger.New())

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+tt.userID, nil)
			w := httptest.NewRecorder()
			req = mux.SetURLVars(req, map[string]string{"id": tt.userID})

			handler.DeleteUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestUserHandler_SuspendUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupMocks     func(*MockUserUseCaseForHandler)
		expectedStatus int
	}{
		{
			name:   "successful suspension",
			userID: "user-123",
			setupMocks: func(m *MockUserUseCaseForHandler) {
				m.SuspendUserFunc = func(ctx context.Context, id string) error {
					return nil
				}
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserUseCase := &MockUserUseCaseForHandler{}
			mockAuthUseCase := &MockAuthUseCase{}
			tt.setupMocks(mockUserUseCase)

			handler := NewUserHandler(mockUserUseCase, mockAuthUseCase, logger.New())

			req := httptest.NewRequest(http.MethodPost, "/api/v1/users/"+tt.userID+"/suspend", nil)
			w := httptest.NewRecorder()
			req = mux.SetURLVars(req, map[string]string{"id": tt.userID})

			handler.SuspendUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}