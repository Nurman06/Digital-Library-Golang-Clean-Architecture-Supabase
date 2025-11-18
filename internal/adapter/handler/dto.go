package handler

import "time"

// Auth DTOs

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	FullName       string    `json:"full_name"`
	Role           string    `json:"role"`
	Status         string    `json:"status"`
	BorrowingLimit int       `json:"borrowing_limit"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Book DTOs

// CreateBookRequest represents a request to create a book
type CreateBookRequest struct {
	Title           string `json:"title"`
	Author          string `json:"author"`
	ISBN            string `json:"isbn"`
	Category        string `json:"category"`
	PublicationYear int    `json:"publication_year"`
	Description     string `json:"description"`
}

// UpdateBookRequest represents a request to update a book
type UpdateBookRequest struct {
	Title           string `json:"title"`
	Author          string `json:"author"`
	ISBN            string `json:"isbn"`
	Category        string `json:"category"`
	PublicationYear int    `json:"publication_year"`
	Description     string `json:"description"`
}

// BookResponse represents a book in API responses
type BookResponse struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Author          string     `json:"author"`
	ISBN            string     `json:"isbn"`
	Category        string     `json:"category"`
	PublicationYear int        `json:"publication_year"`
	Description     string     `json:"description"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

// BookWithAvailabilityResponse represents a book with availability info
type BookWithAvailabilityResponse struct {
	BookResponse
	TotalCopies     int `json:"total_copies"`
	AvailableCopies int `json:"available_copies"`
}

// BookCopy DTOs

// CreateBookCopyRequest represents a request to create a book copy
type CreateBookCopyRequest struct {
	BookID     string `json:"book_id"`
	CopyNumber string `json:"copy_number"`
	Location   string `json:"location"`
}

// UpdateBookCopyRequest represents a request to update a book copy
type UpdateBookCopyRequest struct {
	CopyNumber string `json:"copy_number"`
	Status     string `json:"status"`
	Location   string `json:"location"`
}

// BookCopyResponse represents a book copy in API responses
type BookCopyResponse struct {
	ID         string    `json:"id"`
	BookID     string    `json:"book_id"`
	CopyNumber string    `json:"copy_number"`
	Status     string    `json:"status"`
	Location   string    `json:"location"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Borrowing DTOs

// CheckoutRequest represents a request to checkout a book
type CheckoutRequest struct {
	BookID string `json:"book_id"`
}

// ReturnBookRequest represents a request to return a book
type ReturnBookRequest struct {
	BorrowRecordID string `json:"borrow_record_id"`
}

// RenewBookRequest represents a request to renew a book
type RenewBookRequest struct {
	BorrowRecordID string `json:"borrow_record_id"`
}

// BorrowRecordResponse represents a borrow record in API responses
type BorrowRecordResponse struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	BookCopyID   string     `json:"book_copy_id"`
	CheckoutDate time.Time  `json:"checkout_date"`
	DueDate      time.Time  `json:"due_date"`
	ReturnDate   *time.Time `json:"return_date,omitempty"`
	RenewalCount int        `json:"renewal_count"`
	LateFee      float64    `json:"late_fee"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// BorrowRecordWithDetailsResponse represents a borrow record with book and user details
type BorrowRecordWithDetailsResponse struct {
	BorrowRecordResponse
	Book     *BookResponse     `json:"book,omitempty"`
	BookCopy *BookCopyResponse `json:"book_copy,omitempty"`
	User     *UserResponse     `json:"user,omitempty"`
}

// Search DTOs

// SearchBooksRequest represents a book search request
type SearchBooksRequest struct {
	Query    string `json:"query"`
	Category string `json:"category"`
	Author   string `json:"author"`
	ISBN     string `json:"isbn"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	SortBy   string `json:"sort_by"`
	SortDir  string `json:"sort_dir"`
}

// User Profile DTOs

// UpdateUserProfileRequest represents a request to update user profile
type UpdateUserProfileRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

// UpdateUserRequest represents a request to update a user (admin only)
type UpdateUserRequest struct {
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	Status         string `json:"status"`
	BorrowingLimit int    `json:"borrowing_limit"`
}

// Availability DTOs

// AvailabilityResponse represents availability information for a book
type AvailabilityResponse struct {
	BookID          string                  `json:"book_id"`
	TotalCopies     int                     `json:"total_copies"`
	AvailableCopies int                     `json:"available_copies"`
	BorrowedCopies  int                     `json:"borrowed_copies"`
	ReservedCopies  int                     `json:"reserved_copies"`
	Copies          []*BookCopyResponse     `json:"copies"`
	NextAvailable   *time.Time              `json:"next_available,omitempty"`
}