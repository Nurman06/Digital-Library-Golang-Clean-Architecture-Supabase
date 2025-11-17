package entity

import (
	"errors"
	"regexp"
	"time"
)

// Book represents a book in the library catalog
type Book struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	ISBN            string    `json:"isbn"`
	Category        string    `json:"category"`
	PublicationYear int       `json:"publication_year"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

// BookCopy represents a physical or digital copy of a book
type BookCopy struct {
	ID         string     `json:"id"`
	BookID     string     `json:"book_id"`
	CopyNumber string     `json:"copy_number"`
	Status     CopyStatus `json:"status"`
	Location   string     `json:"location"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// CopyStatus represents the status of a book copy
type CopyStatus string

const (
	CopyStatusAvailable CopyStatus = "available"
	CopyStatusBorrowed  CopyStatus = "borrowed"
	CopyStatusReserved  CopyStatus = "reserved"
	CopyStatusLost      CopyStatus = "lost"
	CopyStatusDamaged   CopyStatus = "damaged"
)

var (
	// ISBN validation patterns
	isbn10Pattern = regexp.MustCompile(`^\d{9}[\dX]$`)
	isbn13Pattern = regexp.MustCompile(`^\d{13}$`)
)

// Validate validates the book entity
func (b *Book) Validate() error {
	if b.Title == "" {
		return errors.New("title is required")
	}
	if b.Author == "" {
		return errors.New("author is required")
	}
	if b.ISBN == "" {
		return errors.New("ISBN is required")
	}
	if !isValidISBN(b.ISBN) {
		return errors.New("invalid ISBN format")
	}
	if b.PublicationYear > time.Now().Year() {
		return errors.New("publication year cannot be in the future")
	}
	if b.PublicationYear < 1000 {
		return errors.New("publication year must be a valid year")
	}
	return nil
}

// isValidISBN checks if the ISBN is in valid format (ISBN-10 or ISBN-13)
func isValidISBN(isbn string) bool {
	return isbn10Pattern.MatchString(isbn) || isbn13Pattern.MatchString(isbn)
}

// Validate validates the book copy entity
func (bc *BookCopy) Validate() error {
	if bc.BookID == "" {
		return errors.New("book ID is required")
	}
	if bc.CopyNumber == "" {
		return errors.New("copy number is required")
	}
	if bc.Status == "" {
		return errors.New("status is required")
	}
	if !isValidCopyStatus(bc.Status) {
		return errors.New("invalid copy status")
	}
	return nil
}

// isValidCopyStatus checks if the copy status is valid
func isValidCopyStatus(status CopyStatus) bool {
	switch status {
	case CopyStatusAvailable, CopyStatusBorrowed, CopyStatusReserved, CopyStatusLost, CopyStatusDamaged:
		return true
	}
	return false
}

// IsAvailable checks if the book copy is available for borrowing
func (bc *BookCopy) IsAvailable() bool {
	return bc.Status == CopyStatusAvailable
}