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

// Mock BookUseCase for AvailabilityHandler
type MockBookUseCaseForAvailability struct {
	GetBookByIDFunc             func(ctx context.Context, id string) (*entity.Book, error)
	GetBookCopiesFunc           func(ctx context.Context, bookID string) ([]*entity.BookCopy, error)
	GetAvailableBookCopiesFunc  func(ctx context.Context, bookID string) ([]*entity.BookCopy, error)
}

func (m *MockBookUseCaseForAvailability) GetBookByID(ctx context.Context, id string) (*entity.Book, error) {
	if m.GetBookByIDFunc != nil {
		return m.GetBookByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockBookUseCaseForAvailability) GetBookCopies(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
	if m.GetBookCopiesFunc != nil {
		return m.GetBookCopiesFunc(ctx, bookID)
	}
	return nil, nil
}

func (m *MockBookUseCaseForAvailability) GetAvailableBookCopies(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
	if m.GetAvailableBookCopiesFunc != nil {
		return m.GetAvailableBookCopiesFunc(ctx, bookID)
	}
	return nil, nil
}

// Implement other required methods
func (m *MockBookUseCaseForAvailability) CreateBook(ctx context.Context, book *entity.Book) error {
	return nil
}

func (m *MockBookUseCaseForAvailability) UpdateBook(ctx context.Context, book *entity.Book) error {
	return nil
}

func (m *MockBookUseCaseForAvailability) DeleteBook(ctx context.Context, id string) error {
	return nil
}

func (m *MockBookUseCaseForAvailability) ListBooks(ctx context.Context, params repository.ListParams) ([]*entity.Book, int64, error) {
	return nil, 0, nil
}

func (m *MockBookUseCaseForAvailability) SearchBooks(ctx context.Context, params repository.SearchParams) ([]*entity.Book, int64, error) {
	return nil, 0, nil
}

func (m *MockBookUseCaseForAvailability) GetBookAvailabilityCount(ctx context.Context, bookID string) (int64, error) {
	return 0, nil
}

func (m *MockBookUseCaseForAvailability) CreateBookCopy(ctx context.Context, copy *entity.BookCopy) error {
	return nil
}

func (m *MockBookUseCaseForAvailability) GetBookByISBN(ctx context.Context, isbn string) (*entity.Book, error) {
	return nil, nil
}

func (m *MockBookUseCaseForAvailability) RestoreBook(ctx context.Context, id string) error {
	return nil
}

func (m *MockBookUseCaseForAvailability) GetBooksByCategory(ctx context.Context, category string, params repository.ListParams) ([]*entity.Book, int64, error) {
	return nil, 0, nil
}

func (m *MockBookUseCaseForAvailability) GetBooksByAuthor(ctx context.Context, author string, params repository.ListParams) ([]*entity.Book, int64, error) {
	return nil, 0, nil
}

func (m *MockBookUseCaseForAvailability) GetRecentlyAddedBooks(ctx context.Context, limit int) ([]*entity.Book, error) {
	return nil, nil
}

func (m *MockBookUseCaseForAvailability) GetBookCount(ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *MockBookUseCaseForAvailability) GetBookCopyByID(ctx context.Context, id string) (*entity.BookCopy, error) {
	return nil, nil
}

func (m *MockBookUseCaseForAvailability) UpdateBookCopy(ctx context.Context, copy *entity.BookCopy) error {
	return nil
}

func (m *MockBookUseCaseForAvailability) UpdateBookCopyStatus(ctx context.Context, id string, status entity.CopyStatus) error {
	return nil
}

func (m *MockBookUseCaseForAvailability) DeleteBookCopy(ctx context.Context, id string) error {
	return nil
}

func TestAvailabilityHandler_GetBookAvailability(t *testing.T) {
	tests := []struct {
		name           string
		bookID         string
		setupMocks     func(*MockBookUseCaseForAvailability)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "successful availability retrieval",
			bookID: "book-123",
			setupMocks: func(m *MockBookUseCaseForAvailability) {
				m.GetBookByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return &entity.Book{
						ID:              id,
						Title:           "Test Book",
						Author:          "Test Author",
						ISBN:            "9780134190440",
						Category:        "Programming",
						PublicationYear: 2020,
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					}, nil
				}
				m.GetBookCopiesFunc = func(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
					return []*entity.BookCopy{
						{
							ID:         "copy-1",
							BookID:     bookID,
							CopyNumber: "COPY-001",
							Status:     entity.CopyStatusAvailable,
							Location:   "Shelf A1",
							CreatedAt:  time.Now(),
							UpdatedAt:  time.Now(),
						},
						{
							ID:         "copy-2",
							BookID:     bookID,
							CopyNumber: "COPY-002",
							Status:     entity.CopyStatusBorrowed,
							Location:   "Shelf A1",
							CreatedAt:  time.Now(),
							UpdatedAt:  time.Now(),
						},
						{
							ID:         "copy-3",
							BookID:     bookID,
							CopyNumber: "COPY-003",
							Status:     entity.CopyStatusReserved,
							Location:   "Shelf A1",
							CreatedAt:  time.Now(),
							UpdatedAt:  time.Now(),
						},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:   "book not found",
			bookID: "nonexistent",
			setupMocks: func(m *MockBookUseCaseForAvailability) {
				m.GetBookByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return nil, errors.New("book not found")
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:   "failed to get copies",
			bookID: "book-123",
			setupMocks: func(m *MockBookUseCaseForAvailability) {
				m.GetBookByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return &entity.Book{
						ID:    id,
						Title: "Test Book",
					}, nil
				}
				m.GetBookCopiesFunc = func(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
					return nil, errors.New("database error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockBookUseCase := &MockBookUseCaseForAvailability{}
			tt.setupMocks(mockBookUseCase)

			handler := NewAvailabilityHandler(mockBookUseCase, logger.New())

			req := httptest.NewRequest(http.MethodGet, "/api/v1/books/"+tt.bookID+"/availability", nil)
			w := httptest.NewRecorder()

			// Set mux vars
			req = mux.SetURLVars(req, map[string]string{"id": tt.bookID})

			// Execute
			handler.GetBookAvailability(w, req)

			// Assert
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

			// Verify availability counts for successful response
			if !tt.expectedError {
				data, ok := response.Data.(map[string]interface{})
				if !ok {
					t.Error("Expected data to be a map")
					return
				}

				if totalCopies, ok := data["total_copies"].(float64); !ok || totalCopies != 3 {
					t.Errorf("Expected total_copies to be 3, got %v", data["total_copies"])
				}
				if availableCopies, ok := data["available_copies"].(float64); !ok || availableCopies != 1 {
					t.Errorf("Expected available_copies to be 1, got %v", data["available_copies"])
				}
				if borrowedCopies, ok := data["borrowed_copies"].(float64); !ok || borrowedCopies != 1 {
					t.Errorf("Expected borrowed_copies to be 1, got %v", data["borrowed_copies"])
				}
				if reservedCopies, ok := data["reserved_copies"].(float64); !ok || reservedCopies != 1 {
					t.Errorf("Expected reserved_copies to be 1, got %v", data["reserved_copies"])
				}
			}
		})
	}
}

func TestAvailabilityHandler_GetAvailableBookCopies(t *testing.T) {
	tests := []struct {
		name           string
		bookID         string
		setupMocks     func(*MockBookUseCaseForAvailability)
		expectedStatus int
		expectedError  bool
		expectedCount  int
	}{
		{
			name:   "successful available copies retrieval",
			bookID: "book-123",
			setupMocks: func(m *MockBookUseCaseForAvailability) {
				m.GetAvailableBookCopiesFunc = func(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
					return []*entity.BookCopy{
						{
							ID:         "copy-1",
							BookID:     bookID,
							CopyNumber: "COPY-001",
							Status:     entity.CopyStatusAvailable,
							Location:   "Shelf A1",
							CreatedAt:  time.Now(),
							UpdatedAt:  time.Now(),
						},
						{
							ID:         "copy-2",
							BookID:     bookID,
							CopyNumber: "COPY-002",
							Status:     entity.CopyStatusAvailable,
							Location:   "Shelf A2",
							CreatedAt:  time.Now(),
							UpdatedAt:  time.Now(),
						},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedCount:  2,
		},
		{
			name:   "book not found",
			bookID: "nonexistent",
			setupMocks: func(m *MockBookUseCaseForAvailability) {
				m.GetAvailableBookCopiesFunc = func(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
					return nil, errors.New("book not found")
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:   "no available copies",
			bookID: "book-123",
			setupMocks: func(m *MockBookUseCaseForAvailability) {
				m.GetAvailableBookCopiesFunc = func(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
					return []*entity.BookCopy{}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockBookUseCase := &MockBookUseCaseForAvailability{}
			tt.setupMocks(mockBookUseCase)

			handler := NewAvailabilityHandler(mockBookUseCase, logger.New())

			req := httptest.NewRequest(http.MethodGet, "/api/v1/books/"+tt.bookID+"/copies/available", nil)
			w := httptest.NewRecorder()

			// Set mux vars
			req = mux.SetURLVars(req, map[string]string{"id": tt.bookID})

			// Execute
			handler.GetAvailableBookCopies(w, req)

			// Assert
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

			// Verify count for successful response
			if !tt.expectedError {
				data, ok := response.Data.([]interface{})
				if !ok {
					t.Error("Expected data to be an array")
					return
				}

				if len(data) != tt.expectedCount {
					t.Errorf("Expected %d copies, got %d", tt.expectedCount, len(data))
				}
			}
		})
	}
}

func TestAvailabilityHandler_GetBookCopies(t *testing.T) {
	tests := []struct {
		name           string
		bookID         string
		setupMocks     func(*MockBookUseCaseForAvailability)
		expectedStatus int
		expectedError  bool
		expectedCount  int
	}{
		{
			name:   "successful copies retrieval",
			bookID: "book-123",
			setupMocks: func(m *MockBookUseCaseForAvailability) {
				m.GetBookCopiesFunc = func(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
					return []*entity.BookCopy{
						{
							ID:         "copy-1",
							BookID:     bookID,
							CopyNumber: "COPY-001",
							Status:     entity.CopyStatusAvailable,
							Location:   "Shelf A1",
							CreatedAt:  time.Now(),
							UpdatedAt:  time.Now(),
						},
						{
							ID:         "copy-2",
							BookID:     bookID,
							CopyNumber: "COPY-002",
							Status:     entity.CopyStatusBorrowed,
							Location:   "Shelf A2",
							CreatedAt:  time.Now(),
							UpdatedAt:  time.Now(),
						},
						{
							ID:         "copy-3",
							BookID:     bookID,
							CopyNumber: "COPY-003",
							Status:     entity.CopyStatusDamaged,
							Location:   "Shelf A3",
							CreatedAt:  time.Now(),
							UpdatedAt:  time.Now(),
						},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedCount:  3,
		},
		{
			name:   "book not found",
			bookID: "nonexistent",
			setupMocks: func(m *MockBookUseCaseForAvailability) {
				m.GetBookCopiesFunc = func(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
					return nil, errors.New("book not found")
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:   "book with no copies",
			bookID: "book-123",
			setupMocks: func(m *MockBookUseCaseForAvailability) {
				m.GetBookCopiesFunc = func(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
					return []*entity.BookCopy{}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockBookUseCase := &MockBookUseCaseForAvailability{}
			tt.setupMocks(mockBookUseCase)

			handler := NewAvailabilityHandler(mockBookUseCase, logger.New())

			req := httptest.NewRequest(http.MethodGet, "/api/v1/books/"+tt.bookID+"/copies", nil)
			w := httptest.NewRecorder()

			// Set mux vars
			req = mux.SetURLVars(req, map[string]string{"id": tt.bookID})

			// Execute
			handler.GetBookCopies(w, req)

			// Assert
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

			// Verify count for successful response
			if !tt.expectedError {
				data, ok := response.Data.([]interface{})
				if !ok {
					t.Error("Expected data to be an array")
					return
				}

				if len(data) != tt.expectedCount {
					t.Errorf("Expected %d copies, got %d", tt.expectedCount, len(data))
				}
			}
		})
	}
}