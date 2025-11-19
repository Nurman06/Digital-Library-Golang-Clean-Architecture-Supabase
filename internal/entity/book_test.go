package entity

import (
	"testing"
	"time"
)

func TestBook_Validate(t *testing.T) {
	tests := []struct {
		name    string
		book    *Book
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid book",
			book: &Book{
				Title:           "The Go Programming Language",
				Author:          "Alan Donovan",
				ISBN:            "9780134190440",
				Category:        "Programming",
				PublicationYear: 2015,
			},
			wantErr: false,
		},
		{
			name: "valid book with ISBN-10",
			book: &Book{
				Title:           "Clean Code",
				Author:          "Robert Martin",
				ISBN:            "0132350882",
				Category:        "Programming",
				PublicationYear: 2008,
			},
			wantErr: false,
		},
		{
			name: "missing title",
			book: &Book{
				Author:          "Alan Donovan",
				ISBN:            "9780134190440",
				PublicationYear: 2015,
			},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name: "missing author",
			book: &Book{
				Title:           "The Go Programming Language",
				ISBN:            "9780134190440",
				PublicationYear: 2015,
			},
			wantErr: true,
			errMsg:  "author is required",
		},
		{
			name: "missing ISBN",
			book: &Book{
				Title:           "The Go Programming Language",
				Author:          "Alan Donovan",
				PublicationYear: 2015,
			},
			wantErr: true,
			errMsg:  "ISBN is required",
		},
		{
			name: "invalid ISBN format",
			book: &Book{
				Title:           "The Go Programming Language",
				Author:          "Alan Donovan",
				ISBN:            "invalid-isbn",
				PublicationYear: 2015,
			},
			wantErr: true,
			errMsg:  "invalid ISBN format",
		},
		{
			name: "publication year in future",
			book: &Book{
				Title:           "The Go Programming Language",
				Author:          "Alan Donovan",
				ISBN:            "9780134190440",
				PublicationYear: time.Now().Year() + 1,
			},
			wantErr: true,
			errMsg:  "publication year cannot be in the future",
		},
		{
			name: "invalid publication year",
			book: &Book{
				Title:           "The Go Programming Language",
				Author:          "Alan Donovan",
				ISBN:            "9780134190440",
				PublicationYear: 999,
			},
			wantErr: true,
			errMsg:  "publication year must be a valid year",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.book.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Book.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Book.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestBook_IsDeleted(t *testing.T) {
	tests := []struct {
		name string
		book *Book
		want bool
	}{
		{
			name: "not deleted",
			book: &Book{
				DeletedAt: nil,
			},
			want: false,
		},
		{
			name: "deleted",
			book: &Book{
				DeletedAt: &time.Time{},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.book.IsDeleted(); got != tt.want {
				t.Errorf("Book.IsDeleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBook_SoftDelete(t *testing.T) {
	book := &Book{
		ID:        "test-id",
		Title:     "Test Book",
		Author:    "Test Author",
		ISBN:      "9780134190440",
		DeletedAt: nil,
	}

	if book.IsDeleted() {
		t.Error("Book should not be deleted initially")
	}

	book.SoftDelete()

	if !book.IsDeleted() {
		t.Error("Book should be deleted after SoftDelete()")
	}

	if book.DeletedAt == nil {
		t.Error("DeletedAt should not be nil after SoftDelete()")
	}
}

func TestBook_Restore(t *testing.T) {
	now := time.Now()
	book := &Book{
		ID:        "test-id",
		Title:     "Test Book",
		Author:    "Test Author",
		ISBN:      "9780134190440",
		DeletedAt: &now,
	}

	if !book.IsDeleted() {
		t.Error("Book should be deleted initially")
	}

	book.Restore()

	if book.IsDeleted() {
		t.Error("Book should not be deleted after Restore()")
	}

	if book.DeletedAt != nil {
		t.Error("DeletedAt should be nil after Restore()")
	}
}

func TestBookCopy_Validate(t *testing.T) {
	tests := []struct {
		name    string
		copy    *BookCopy
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid book copy",
			copy: &BookCopy{
				BookID:     "book-id",
				CopyNumber: "COPY-001",
				Status:     CopyStatusAvailable,
			},
			wantErr: false,
		},
		{
			name: "missing book ID",
			copy: &BookCopy{
				CopyNumber: "COPY-001",
				Status:     CopyStatusAvailable,
			},
			wantErr: true,
			errMsg:  "book ID is required",
		},
		{
			name: "missing copy number",
			copy: &BookCopy{
				BookID: "book-id",
				Status: CopyStatusAvailable,
			},
			wantErr: true,
			errMsg:  "copy number is required",
		},
		{
			name: "missing status",
			copy: &BookCopy{
				BookID:     "book-id",
				CopyNumber: "COPY-001",
			},
			wantErr: true,
			errMsg:  "status is required",
		},
		{
			name: "invalid status",
			copy: &BookCopy{
				BookID:     "book-id",
				CopyNumber: "COPY-001",
				Status:     "invalid-status",
			},
			wantErr: true,
			errMsg:  "invalid copy status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.copy.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BookCopy.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("BookCopy.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestBookCopy_StatusChecks(t *testing.T) {
	tests := []struct {
		name   string
		status CopyStatus
		checks map[string]bool
	}{
		{
			name:   "available copy",
			status: CopyStatusAvailable,
			checks: map[string]bool{
				"IsAvailable":   true,
				"IsBorrowed":    false,
				"IsReserved":    false,
				"IsDamaged":     false,
				"IsLost":        false,
				"CanBeBorrowed": true,
			},
		},
		{
			name:   "borrowed copy",
			status: CopyStatusBorrowed,
			checks: map[string]bool{
				"IsAvailable":   false,
				"IsBorrowed":    true,
				"IsReserved":    false,
				"IsDamaged":     false,
				"IsLost":        false,
				"CanBeBorrowed": false,
			},
		},
		{
			name:   "reserved copy",
			status: CopyStatusReserved,
			checks: map[string]bool{
				"IsAvailable":   false,
				"IsBorrowed":    false,
				"IsReserved":    true,
				"IsDamaged":     false,
				"IsLost":        false,
				"CanBeBorrowed": false,
			},
		},
		{
			name:   "damaged copy",
			status: CopyStatusDamaged,
			checks: map[string]bool{
				"IsAvailable":   false,
				"IsBorrowed":    false,
				"IsReserved":    false,
				"IsDamaged":     true,
				"IsLost":        false,
				"CanBeBorrowed": false,
			},
		},
		{
			name:   "lost copy",
			status: CopyStatusLost,
			checks: map[string]bool{
				"IsAvailable":   false,
				"IsBorrowed":    false,
				"IsReserved":    false,
				"IsDamaged":     false,
				"IsLost":        true,
				"CanBeBorrowed": false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copy := &BookCopy{Status: tt.status}

			if got := copy.IsAvailable(); got != tt.checks["IsAvailable"] {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.checks["IsAvailable"])
			}
			if got := copy.IsBorrowed(); got != tt.checks["IsBorrowed"] {
				t.Errorf("IsBorrowed() = %v, want %v", got, tt.checks["IsBorrowed"])
			}
			if got := copy.IsReserved(); got != tt.checks["IsReserved"] {
				t.Errorf("IsReserved() = %v, want %v", got, tt.checks["IsReserved"])
			}
			if got := copy.IsDamaged(); got != tt.checks["IsDamaged"] {
				t.Errorf("IsDamaged() = %v, want %v", got, tt.checks["IsDamaged"])
			}
			if got := copy.IsLost(); got != tt.checks["IsLost"] {
				t.Errorf("IsLost() = %v, want %v", got, tt.checks["IsLost"])
			}
			if got := copy.CanBeBorrowed(); got != tt.checks["CanBeBorrowed"] {
				t.Errorf("CanBeBorrowed() = %v, want %v", got, tt.checks["CanBeBorrowed"])
			}
		})
	}
}

func TestBookCopy_MarkAsBorrowed(t *testing.T) {
	tests := []struct {
		name    string
		status  CopyStatus
		wantErr bool
	}{
		{
			name:    "available copy can be borrowed",
			status:  CopyStatusAvailable,
			wantErr: false,
		},
		{
			name:    "borrowed copy cannot be borrowed",
			status:  CopyStatusBorrowed,
			wantErr: true,
		},
		{
			name:    "reserved copy cannot be borrowed",
			status:  CopyStatusReserved,
			wantErr: true,
		},
		{
			name:    "damaged copy cannot be borrowed",
			status:  CopyStatusDamaged,
			wantErr: true,
		},
		{
			name:    "lost copy cannot be borrowed",
			status:  CopyStatusLost,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copy := &BookCopy{
				BookID:     "book-id",
				CopyNumber: "COPY-001",
				Status:     tt.status,
			}

			err := copy.MarkAsBorrowed()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarkAsBorrowed() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && copy.Status != CopyStatusBorrowed {
				t.Errorf("Status should be borrowed, got %v", copy.Status)
			}
		})
	}
}

func TestBookCopy_MarkAsAvailable(t *testing.T) {
	tests := []struct {
		name    string
		status  CopyStatus
		wantErr bool
	}{
		{
			name:    "borrowed copy can be marked available",
			status:  CopyStatusBorrowed,
			wantErr: false,
		},
		{
			name:    "reserved copy can be marked available",
			status:  CopyStatusReserved,
			wantErr: false,
		},
		{
			name:    "damaged copy can be marked available",
			status:  CopyStatusDamaged,
			wantErr: false,
		},
		{
			name:    "lost copy cannot be marked available",
			status:  CopyStatusLost,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copy := &BookCopy{
				BookID:     "book-id",
				CopyNumber: "COPY-001",
				Status:     tt.status,
			}

			err := copy.MarkAsAvailable()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarkAsAvailable() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && copy.Status != CopyStatusAvailable {
				t.Errorf("Status should be available, got %v", copy.Status)
			}
		})
	}
}

func TestBookCopy_MarkAsReserved(t *testing.T) {
	tests := []struct {
		name    string
		status  CopyStatus
		wantErr bool
	}{
		{
			name:    "available copy can be reserved",
			status:  CopyStatusAvailable,
			wantErr: false,
		},
		{
			name:    "borrowed copy cannot be reserved",
			status:  CopyStatusBorrowed,
			wantErr: true,
		},
		{
			name:    "damaged copy cannot be reserved",
			status:  CopyStatusDamaged,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copy := &BookCopy{
				BookID:     "book-id",
				CopyNumber: "COPY-001",
				Status:     tt.status,
			}

			err := copy.MarkAsReserved()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarkAsReserved() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && copy.Status != CopyStatusReserved {
				t.Errorf("Status should be reserved, got %v", copy.Status)
			}
		})
	}
}

func TestBookCopy_MarkAsDamaged(t *testing.T) {
	copy := &BookCopy{
		BookID:     "book-id",
		CopyNumber: "COPY-001",
		Status:     CopyStatusAvailable,
	}

	err := copy.MarkAsDamaged()
	if err != nil {
		t.Errorf("MarkAsDamaged() unexpected error = %v", err)
	}

	if copy.Status != CopyStatusDamaged {
		t.Errorf("Status should be damaged, got %v", copy.Status)
	}
}

func TestBookCopy_MarkAsLost(t *testing.T) {
	copy := &BookCopy{
		BookID:     "book-id",
		CopyNumber: "COPY-001",
		Status:     CopyStatusBorrowed,
	}

	err := copy.MarkAsLost()
	if err != nil {
		t.Errorf("MarkAsLost() unexpected error = %v", err)
	}

	if copy.Status != CopyStatusLost {
		t.Errorf("Status should be lost, got %v", copy.Status)
	}
}