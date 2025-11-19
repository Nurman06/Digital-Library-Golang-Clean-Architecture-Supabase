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

// Mock BookUseCase
type MockBookUseCase struct {
	CreateBookFunc               func(ctx context.Context, book *entity.Book) error
	GetBookByIDFunc              func(ctx context.Context, id string) (*entity.Book, error)
	UpdateBookFunc               func(ctx context.Context, book *entity.Book) error
	DeleteBookFunc               func(ctx context.Context, id string) error
	ListBooksFunc                func(ctx context.Context, params repository.ListParams) ([]*entity.Book, int64, error)
	SearchBooksFunc              func(ctx context.Context, params repository.SearchParams) ([]*entity.Book, int64, error)
	GetBookAvailabilityCountFunc func(ctx context.Context, bookID string) (int64, error)
	GetBookCopiesFunc            func(ctx context.Context, bookID string) ([]*entity.BookCopy, error)
	CreateBookCopyFunc           func(ctx context.Context, copy *entity.BookCopy) error
	DeleteBookCopyFunc           func(ctx context.Context, id string) error
}

func (m *MockBookUseCase) CreateBook(ctx context.Context, book *entity.Book) error {
	if m.CreateBookFunc != nil {
		return m.CreateBookFunc(ctx, book)
	}
	return nil
}

func (m *MockBookUseCase) GetBookByID(ctx context.Context, id string) (*entity.Book, error) {
	if m.GetBookByIDFunc != nil {
		return m.GetBookByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockBookUseCase) UpdateBook(ctx context.Context, book *entity.Book) error {
	if m.UpdateBookFunc != nil {
		return m.UpdateBookFunc(ctx, book)
	}
	return nil
}

func (m *MockBookUseCase) DeleteBook(ctx context.Context, id string) error {
	if m.DeleteBookFunc != nil {
		return m.DeleteBookFunc(ctx, id)
	}
	return nil
}

func (m *MockBookUseCase) ListBooks(ctx context.Context, params repository.ListParams) ([]*entity.Book, int64, error) {
	if m.ListBooksFunc != nil {
		return m.ListBooksFunc(ctx, params)
	}
	return nil, 0, nil
}

func (m *MockBookUseCase) SearchBooks(ctx context.Context, params repository.SearchParams) ([]*entity.Book, int64, error) {
	if m.SearchBooksFunc != nil {
		return m.SearchBooksFunc(ctx, params)
	}
	return nil, 0, nil
}

func (m *MockBookUseCase) GetBookAvailabilityCount(ctx context.Context, bookID string) (int64, error) {
	if m.GetBookAvailabilityCountFunc != nil {
		return m.GetBookAvailabilityCountFunc(ctx, bookID)
	}
	return 0, nil
}

func (m *MockBookUseCase) GetBookCopies(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
	if m.GetBookCopiesFunc != nil {
		return m.GetBookCopiesFunc(ctx, bookID)
	}
	return nil, nil
}

func (m *MockBookUseCase) CreateBookCopy(ctx context.Context, copy *entity.BookCopy) error {
	if m.CreateBookCopyFunc != nil {
		return m.CreateBookCopyFunc(ctx, copy)
	}
	return nil
}

func (m *MockBookUseCase) GetBookByISBN(ctx context.Context, isbn string) (*entity.Book, error) {
	return nil, nil
}

func (m *MockBookUseCase) RestoreBook(ctx context.Context, id string) error {
	return nil
}

func (m *MockBookUseCase) GetBooksByCategory(ctx context.Context, category string, params repository.ListParams) ([]*entity.Book, int64, error) {
	return nil, 0, nil
}

func (m *MockBookUseCase) GetBooksByAuthor(ctx context.Context, author string, params repository.ListParams) ([]*entity.Book, int64, error) {
	return nil, 0, nil
}

func (m *MockBookUseCase) GetRecentlyAddedBooks(ctx context.Context, limit int) ([]*entity.Book, error) {
	return nil, nil
}

func (m *MockBookUseCase) GetBookCount(ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *MockBookUseCase) GetBookCopyByID(ctx context.Context, id string) (*entity.BookCopy, error) {
	return nil, nil
}

func (m *MockBookUseCase) GetAvailableBookCopies(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
	return nil, nil
}

func (m *MockBookUseCase) UpdateBookCopy(ctx context.Context, copy *entity.BookCopy) error {
	return nil
}

func (m *MockBookUseCase) UpdateBookCopyStatus(ctx context.Context, id string, status entity.CopyStatus) error {
	return nil
}

func (m *MockBookUseCase) DeleteBookCopy(ctx context.Context, id string) error {
	if m.DeleteBookCopyFunc != nil {
		return m.DeleteBookCopyFunc(ctx, id)
	}
	return nil
}

// Mock AuthUseCase
type MockAuthUseCase struct{}

func (m *MockAuthUseCase) Login(ctx context.Context, email, password string) (*entity.User, error) {
	return nil, nil
}

func (m *MockAuthUseCase) ValidateUser(ctx context.Context, userID string) error {
	return nil
}

func (m *MockAuthUseCase) CheckPermission(ctx context.Context, userID string, operation string) (bool, error) {
	return false, nil
}

func (m *MockAuthUseCase) ValidateUserStatus(ctx context.Context, userID string) error {
	return nil
}

func (m *MockAuthUseCase) GetUserRole(ctx context.Context, userID string) (entity.UserRole, error) {
	return "", nil
}

func (m *MockAuthUseCase) CanUserBorrow(ctx context.Context, userID string) (bool, error) {
	return false, nil
}

func (m *MockAuthUseCase) CanUserManageBooks(ctx context.Context, userID string) (bool, error) {
	return false, nil
}

func (m *MockAuthUseCase) CanUserManageUsers(ctx context.Context, userID string) (bool, error) {
	return false, nil
}

func TestBookHandler_CreateBook(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockBookUseCase)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "successful book creation",
			requestBody: CreateBookRequest{
				Title:           "Test Book",
				Author:          "Test Author",
				ISBN:            "9780134190440",
				Category:        "Programming",
				PublicationYear: 2020,
				Description:     "Test Description",
			},
			setupMock: func(m *MockBookUseCase) {
				m.CreateBookFunc = func(ctx context.Context, book *entity.Book) error {
					book.ID = "test-id"
					book.CreatedAt = time.Now()
					book.UpdatedAt = time.Now()
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			setupMock:      func(m *MockBookUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "duplicate ISBN error",
			requestBody: CreateBookRequest{
				Title:           "Test Book",
				Author:          "Test Author",
				ISBN:            "9780134190440",
				Category:        "Programming",
				PublicationYear: 2020,
			},
			setupMock: func(m *MockBookUseCase) {
				m.CreateBookFunc = func(ctx context.Context, book *entity.Book) error {
					return errors.New("book with this ISBN already exists")
				}
			},
			expectedStatus: http.StatusConflict,
			expectedError:  true,
		},
		{
			name: "validation error",
			requestBody: CreateBookRequest{
				Title:           "",
				Author:          "Test Author",
				ISBN:            "9780134190440",
				Category:        "Programming",
				PublicationYear: 2020,
			},
			setupMock: func(m *MockBookUseCase) {
				m.CreateBookFunc = func(ctx context.Context, book *entity.Book) error {
					return errors.New("title is required")
				}
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUseCase := &MockBookUseCase{}
			mockAuthUseCase := &MockAuthUseCase{}
			tt.setupMock(mockUseCase)

			handler := NewBookHandler(mockUseCase, mockAuthUseCase, logger.New())

			// Create request
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

			req := httptest.NewRequest(http.MethodPost, "/api/v1/books", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute
			handler.CreateBook(w, req)

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
		})
	}
}

func TestBookHandler_GetBook(t *testing.T) {
	tests := []struct {
		name           string
		bookID         string
		setupMock      func(*MockBookUseCase)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "successful book retrieval",
			bookID: "test-id",
			setupMock: func(m *MockBookUseCase) {
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
				m.GetBookAvailabilityCountFunc = func(ctx context.Context, bookID string) (int64, error) {
					return 3, nil
				}
				m.GetBookCopiesFunc = func(ctx context.Context, bookID string) ([]*entity.BookCopy, error) {
					return []*entity.BookCopy{
						{ID: "copy-1", BookID: bookID, Status: entity.CopyStatusAvailable},
						{ID: "copy-2", BookID: bookID, Status: entity.CopyStatusAvailable},
						{ID: "copy-3", BookID: bookID, Status: entity.CopyStatusBorrowed},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:   "book not found",
			bookID: "nonexistent",
			setupMock: func(m *MockBookUseCase) {
				m.GetBookByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return nil, errors.New("book not found")
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUseCase := &MockBookUseCase{}
			mockAuthUseCase := &MockAuthUseCase{}
			tt.setupMock(mockUseCase)

			handler := NewBookHandler(mockUseCase, mockAuthUseCase, logger.New())

			// Create request with mux vars
			req := httptest.NewRequest(http.MethodGet, "/api/v1/books/"+tt.bookID, nil)
			w := httptest.NewRecorder()

			// Set mux vars
			req = mux.SetURLVars(req, map[string]string{"id": tt.bookID})

			// Execute
			handler.GetBook(w, req)

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
		})
	}
}

func TestBookHandler_ListBooks(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		setupMock      func(*MockBookUseCase)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "successful book listing with default params",
			queryParams: map[string]string{},
			setupMock: func(m *MockBookUseCase) {
				m.ListBooksFunc = func(ctx context.Context, params repository.ListParams) ([]*entity.Book, int64, error) {
					books := []*entity.Book{
						{
							ID:              "book-1",
							Title:           "Test Book 1",
							Author:          "Author 1",
							ISBN:            "1234567890",
							Category:        "Programming",
							PublicationYear: 2020,
							CreatedAt:       time.Now(),
							UpdatedAt:       time.Now(),
						},
						{
							ID:              "book-2",
							Title:           "Test Book 2",
							Author:          "Author 2",
							ISBN:            "0987654321",
							Category:        "Programming",
							PublicationYear: 2021,
							CreatedAt:       time.Now(),
							UpdatedAt:       time.Now(),
						},
					}
					return books, 2, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "successful book listing with custom params",
			queryParams: map[string]string{
				"page":      "2",
				"page_size": "10",
				"sort_by":   "title",
				"sort_dir":  "asc",
			},
			setupMock: func(m *MockBookUseCase) {
				m.ListBooksFunc = func(ctx context.Context, params repository.ListParams) ([]*entity.Book, int64, error) {
					// Verify params
					if params.Page != 2 || params.PageSize != 10 || params.SortBy != "title" || params.SortOrder != "asc" {
						return nil, 0, errors.New("invalid params")
					}
					return []*entity.Book{}, 0, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:        "repository error",
			queryParams: map[string]string{},
			setupMock: func(m *MockBookUseCase) {
				m.ListBooksFunc = func(ctx context.Context, params repository.ListParams) ([]*entity.Book, int64, error) {
					return nil, 0, errors.New("database error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUseCase := &MockBookUseCase{}
			mockAuthUseCase := &MockAuthUseCase{}
			tt.setupMock(mockUseCase)

			handler := NewBookHandler(mockUseCase, mockAuthUseCase, logger.New())

			// Build URL with query params
			url := "/api/v1/books"
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

			// Execute
			handler.ListBooks(w, req)

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

			// Check metadata for successful responses
			if !tt.expectedError && response.Meta == nil {
				t.Error("Expected metadata in successful response")
			}
		})
	}
}

func TestBookHandler_UpdateBook(t *testing.T) {
	tests := []struct {
		name           string
		bookID         string
		requestBody    interface{}
		setupMock      func(*MockBookUseCase)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "successful book update",
			bookID: "test-id",
			requestBody: UpdateBookRequest{
				Title:           "Updated Book",
				Author:          "Updated Author",
				ISBN:            "9780134190440",
				Category:        "Programming",
				PublicationYear: 2021,
				Description:     "Updated Description",
			},
			setupMock: func(m *MockBookUseCase) {
				m.GetBookByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return &entity.Book{
						ID:              id,
						Title:           "Old Book",
						Author:          "Old Author",
						ISBN:            "9780134190440",
						Category:        "Programming",
						PublicationYear: 2020,
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					}, nil
				}
				m.UpdateBookFunc = func(ctx context.Context, book *entity.Book) error {
					return nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:   "book not found",
			bookID: "nonexistent",
			requestBody: UpdateBookRequest{
				Title:           "Updated Book",
				Author:          "Updated Author",
				ISBN:            "9780134190440",
				Category:        "Programming",
				PublicationYear: 2021,
			},
			setupMock: func(m *MockBookUseCase) {
				m.GetBookByIDFunc = func(ctx context.Context, id string) (*entity.Book, error) {
					return nil, errors.New("book not found")
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:           "invalid request body",
			bookID:         "test-id",
			requestBody:    "invalid json",
			setupMock:      func(m *MockBookUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUseCase := &MockBookUseCase{}
			mockAuthUseCase := &MockAuthUseCase{}
			tt.setupMock(mockUseCase)

			handler := NewBookHandler(mockUseCase, mockAuthUseCase, logger.New())

			// Create request
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

			req := httptest.NewRequest(http.MethodPut, "/api/v1/books/"+tt.bookID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Set mux vars
			req = mux.SetURLVars(req, map[string]string{"id": tt.bookID})

			// Execute
			handler.UpdateBook(w, req)

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
		})
	}
}

func TestBookHandler_DeleteBook(t *testing.T) {
	tests := []struct {
		name           string
		bookID         string
		setupMock      func(*MockBookUseCase)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "successful book deletion",
			bookID: "test-id",
			setupMock: func(m *MockBookUseCase) {
				m.DeleteBookFunc = func(ctx context.Context, id string) error {
					return nil
				}
			},
			expectedStatus: http.StatusNoContent,
			expectedError:  false,
		},
		{
			name:   "book not found",
			bookID: "nonexistent",
			setupMock: func(m *MockBookUseCase) {
				m.DeleteBookFunc = func(ctx context.Context, id string) error {
					return errors.New("book not found")
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:   "cannot delete book with active borrows",
			bookID: "test-id",
			setupMock: func(m *MockBookUseCase) {
				m.DeleteBookFunc = func(ctx context.Context, id string) error {
					return errors.New("cannot delete book with active borrows")
				}
			},
			expectedStatus: http.StatusConflict,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUseCase := &MockBookUseCase{}
			mockAuthUseCase := &MockAuthUseCase{}
			tt.setupMock(mockUseCase)

			handler := NewBookHandler(mockUseCase, mockAuthUseCase, logger.New())

			// Create request
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/books/"+tt.bookID, nil)
			w := httptest.NewRecorder()

			// Set mux vars
			req = mux.SetURLVars(req, map[string]string{"id": tt.bookID})

			// Execute
			handler.DeleteBook(w, req)

			// Assert
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// For 204 No Content, there's no body to decode
			if tt.expectedStatus != http.StatusNoContent {
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
			}
		})
	}
}