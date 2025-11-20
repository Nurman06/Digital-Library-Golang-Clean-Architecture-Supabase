package handler

import (
	"bytes"
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

func TestBorrowingHandler_CheckoutBook(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    interface{}
		setupMocks     func(*MockBorrowingUseCase)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "successful checkout",
			userID: "user-123",
			requestBody: CheckoutRequest{
				BookID: "book-123",
			},
			setupMocks: func(m *MockBorrowingUseCase) {
				m.CheckoutBookFunc = func(ctx context.Context, userID, bookID string) (*entity.BorrowRecord, error) {
					return &entity.BorrowRecord{
						ID:           "record-123",
						UserID:       userID,
						BookCopyID:   "copy-123",
						CheckoutDate: time.Now(),
						DueDate:      time.Now().Add(14 * 24 * time.Hour),
						RenewalCount: 0,
						LateFee:      0,
						Status:       entity.BorrowStatusActive,
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
					}, nil
				}
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:           "user not authenticated",
			userID:         "",
			requestBody:    CheckoutRequest{BookID: "book-123"},
			setupMocks:     func(m *MockBorrowingUseCase) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:   "no available copies",
			userID: "user-123",
			requestBody: CheckoutRequest{
				BookID: "book-123",
			},
			setupMocks: func(m *MockBorrowingUseCase) {
				m.CheckoutBookFunc = func(ctx context.Context, userID, bookID string) (*entity.BorrowRecord, error) {
					return nil, errors.New("no available copies of this book")
				}
			},
			expectedStatus: http.StatusConflict,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBorrowingUseCase := &MockBorrowingUseCase{}
			mockBookUseCase := &MockBookUseCase{}
			mockUserUseCase := &MockUserUseCaseForHandler{}
			tt.setupMocks(mockBorrowingUseCase)

			handler := NewBorrowingHandler(mockBorrowingUseCase, mockBookUseCase, mockUserUseCase, logger.New())

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/borrowing/checkout", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			if tt.userID != "" {
				ctx := SetUserIDInContext(req.Context(), tt.userID)
				req = req.WithContext(ctx)
			}

			handler.CheckoutBook(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestBorrowingHandler_ReturnBook(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func(*MockBorrowingUseCase)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "successful return",
			requestBody: ReturnBookRequest{
				BorrowRecordID: "record-123",
			},
			setupMocks: func(m *MockBorrowingUseCase) {
				m.ReturnBookFunc = func(ctx context.Context, borrowRecordID string) error {
					return nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "already returned",
			requestBody: ReturnBookRequest{
				BorrowRecordID: "record-123",
			},
			setupMocks: func(m *MockBorrowingUseCase) {
				m.ReturnBookFunc = func(ctx context.Context, borrowRecordID string) error {
					return errors.New("book has already been returned")
				}
			},
			expectedStatus: http.StatusConflict,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBorrowingUseCase := &MockBorrowingUseCase{}
			mockBookUseCase := &MockBookUseCase{}
			mockUserUseCase := &MockUserUseCaseForHandler{}
			tt.setupMocks(mockBorrowingUseCase)

			handler := NewBorrowingHandler(mockBorrowingUseCase, mockBookUseCase, mockUserUseCase, logger.New())

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/borrowing/return", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.ReturnBook(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestBorrowingHandler_GetBorrowingHistory(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupMocks     func(*MockBorrowingUseCase)
		expectedStatus int
	}{
		{
			name:   "successful history retrieval",
			userID: "user-123",
			setupMocks: func(m *MockBorrowingUseCase) {
				m.GetUserBorrowingHistoryFunc = func(ctx context.Context, userID string, params repository.BorrowRecordListParams) ([]*entity.BorrowRecord, int64, error) {
					records := []*entity.BorrowRecord{
						{
							ID:           "record-1",
							UserID:       userID,
							BookCopyID:   "copy-1",
							CheckoutDate: time.Now(),
							DueDate:      time.Now().Add(14 * 24 * time.Hour),
							Status:       entity.BorrowStatusActive,
							CreatedAt:    time.Now(),
							UpdatedAt:    time.Now(),
						},
					}
					return records, 1, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBorrowingUseCase := &MockBorrowingUseCase{}
			mockBookUseCase := &MockBookUseCase{}
			mockUserUseCase := &MockUserUseCaseForHandler{}
			tt.setupMocks(mockBorrowingUseCase)

			handler := NewBorrowingHandler(mockBorrowingUseCase, mockBookUseCase, mockUserUseCase, logger.New())

			req := httptest.NewRequest(http.MethodGet, "/api/v1/borrowing/history", nil)
			w := httptest.NewRecorder()

			ctx := SetUserIDInContext(req.Context(), tt.userID)
			req = req.WithContext(ctx)

			handler.GetBorrowingHistory(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestBorrowingHandler_GetBorrowRecord(t *testing.T) {
	tests := []struct {
		name           string
		recordID       string
		setupMocks     func(*MockBorrowingUseCase)
		expectedStatus int
	}{
		{
			name:     "successful record retrieval",
			recordID: "record-123",
			setupMocks: func(m *MockBorrowingUseCase) {
				m.GetBorrowRecordFunc = func(ctx context.Context, id string) (*entity.BorrowRecord, error) {
					return &entity.BorrowRecord{
						ID:           id,
						UserID:       "user-123",
						BookCopyID:   "copy-123",
						CheckoutDate: time.Now(),
						DueDate:      time.Now().Add(14 * 24 * time.Hour),
						Status:       entity.BorrowStatusActive,
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBorrowingUseCase := &MockBorrowingUseCase{}
			mockBookUseCase := &MockBookUseCase{}
			mockUserUseCase := &MockUserUseCaseForHandler{}
			tt.setupMocks(mockBorrowingUseCase)

			handler := NewBorrowingHandler(mockBorrowingUseCase, mockBookUseCase, mockUserUseCase, logger.New())

			req := httptest.NewRequest(http.MethodGet, "/api/v1/borrowing/"+tt.recordID, nil)
			w := httptest.NewRecorder()
			req = mux.SetURLVars(req, map[string]string{"id": tt.recordID})

			handler.GetBorrowRecord(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}