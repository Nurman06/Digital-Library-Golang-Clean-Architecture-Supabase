package testdata

import (
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
)

// CreateTestBook creates a test book with default values
func CreateTestBook(id string) *entity.Book {
	return &entity.Book{
		ID:              id,
		Title:           "Test Book",
		Author:          "Test Author",
		ISBN:            "9780134190440",
		Category:        "Programming",
		PublicationYear: 2020,
		Description:     "Test Description",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		DeletedAt:       nil,
	}
}

// CreateTestBookCopy creates a test book copy with default values
func CreateTestBookCopy(id, bookID string) *entity.BookCopy {
	return &entity.BookCopy{
		ID:         id,
		BookID:     bookID,
		CopyNumber: "COPY-001",
		Status:     entity.CopyStatusAvailable,
		Location:   "Shelf A1",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// CreateTestUser creates a test user with default values
func CreateTestUser(id string, role entity.UserRole) *entity.User {
	return &entity.User{
		ID:             id,
		Email:          "test@example.com",
		FullName:       "Test User",
		Role:           role,
		Status:         entity.StatusActive,
		BorrowingLimit: entity.GetDefaultBorrowingLimit(role),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// CreateTestBorrowRecord creates a test borrow record with default values
func CreateTestBorrowRecord(id, userID, bookCopyID string) *entity.BorrowRecord {
	now := time.Now()
	return &entity.BorrowRecord{
		ID:           id,
		UserID:       userID,
		BookCopyID:   bookCopyID,
		CheckoutDate: now,
		DueDate:      now.Add(entity.DefaultBorrowingPeriod),
		ReturnDate:   nil,
		RenewalCount: 0,
		LateFee:      0,
		Status:       entity.BorrowStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// CreateTestReservation creates a test reservation with default values
func CreateTestReservation(id, userID, bookID string) *entity.Reservation {
	now := time.Now()
	return &entity.Reservation{
		ID:              id,
		UserID:          userID,
		BookID:          bookID,
		ReservationDate: now,
		ExpiryDate:      now.Add(entity.DefaultReservationHoldPeriod),
		Status:          entity.ReservationStatusPending,
		QueuePosition:   1,
	}
}

// CreateTestBooks creates multiple test books
func CreateTestBooks(count int) []*entity.Book {
	books := make([]*entity.Book, count)
	for i := 0; i < count; i++ {
		books[i] = &entity.Book{
			ID:              string(rune(i + 1)),
			Title:           "Test Book " + string(rune(i+1)),
			Author:          "Test Author " + string(rune(i+1)),
			ISBN:            "978013419044" + string(rune(i)),
			Category:        "Programming",
			PublicationYear: 2020 + i,
			Description:     "Test Description",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
	}
	return books
}

// CreateTestBookCopies creates multiple test book copies
func CreateTestBookCopies(bookID string, count int) []*entity.BookCopy {
	copies := make([]*entity.BookCopy, count)
	for i := 0; i < count; i++ {
		copies[i] = &entity.BookCopy{
			ID:         string(rune(i + 1)),
			BookID:     bookID,
			CopyNumber: "COPY-00" + string(rune(i+1)),
			Status:     entity.CopyStatusAvailable,
			Location:   "Shelf A" + string(rune(i+1)),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
	}
	return copies
}