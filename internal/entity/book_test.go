package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				Title:           "Clean Code",
				Author:          "Robert C. Martin",
				ISBN:            "9780132350884",
				Category:        "Programming",
				PublicationYear: 2008,
				Description:     "A handbook of agile software craftsmanship",
			},
			wantErr: false,
		},
		{
			name: "valid book with ISBN-10",
			book: &Book{
				Title:           "The Pragmatic Programmer",
				Author:          "Andrew Hunt",
				ISBN:            "020161622X",
				Category:        "Programming",
				PublicationYear: 1999,
			},
			wantErr: false,
		},
		{
			name: "missing title",
			book: &Book{
				Title:           "",
				Author:          "Robert C. Martin",
				ISBN:            "9780132350884",
				PublicationYear: 2008,
			},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name: "missing author",
			book: &Book{
				Title:           "Clean Code",
				Author:          "",
				ISBN:            "9780132350884",
				PublicationYear: 2008,
			},
			wantErr: true,
			errMsg:  "author is required",
		},
		{
			name: "missing ISBN",
			book: &Book{
				Title:           "Clean Code",
				Author:          "Robert C. Martin",
				ISBN:            "",
				PublicationYear: 2008,
			},
			wantErr: true,
			errMsg:  "ISBN is required",
		},
		{
			name: "invalid ISBN format",
			book: &Book{
				Title:           "Clean Code",
				Author:          "Robert C. Martin",
				ISBN:            "invalid-isbn",
				PublicationYear: 2008,
			},
			wantErr: true,
			errMsg:  "invalid ISBN format",
		},
		{
			name: "publication year in future",
			book: &Book{
				Title:           "Clean Code",
				Author:          "Robert C. Martin",
				ISBN:            "9780132350884",
				PublicationYear: time.Now().Year() + 1,
			},
			wantErr: true,
			errMsg:  "publication year cannot be in the future",
		},
		{
			name: "invalid publication year",
			book: &Book{
				Title:           "Clean Code",
				Author:          "Robert C. Martin",
				ISBN:            "9780132350884",
				PublicationYear: 999,
			},
			wantErr: true,
			errMsg:  "publication year must be a valid year",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.book.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBook_IsDeleted(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name     string
		book     *Book
		expected bool
	}{
		{
			name: "not deleted",
			book: &Book{
				DeletedAt: nil,
			},
			expected: false,
		},
		{
			name: "deleted",
			book: &Book{
				DeletedAt: &now,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.book.IsDeleted()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBook_SoftDelete(t *testing.T) {
	book := &Book{
		ID:        "1",
		Title:     "Test Book",
		DeletedAt: nil,
	}

	assert.Nil(t, book.DeletedAt)
	
	book.SoftDelete()
	
	assert.NotNil(t, book.DeletedAt)
	assert.True(t, book.IsDeleted())
}

func TestBook_Restore(t *testing.T) {
	now := time.Now()
	book := &Book{
		ID:        "1",
		Title:     "Test Book",
		DeletedAt: &now,
	}

	assert.True(t, book.IsDeleted())
	
	book.Restore()
	
	assert.Nil(t, book.DeletedAt)
	assert.False(t, book.IsDeleted())
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
				BookID:     "book-1",
				CopyNumber: "COPY-001",
				Status:     CopyStatusAvailable,
				Location:   "Shelf A1",
			},
			wantErr: false,
		},
		{
			name: "missing book ID",
			copy: &BookCopy{
				BookID:     "",
				CopyNumber: "COPY-001",
				Status:     CopyStatusAvailable,
			},
			wantErr: true,
			errMsg:  "book ID is required",
		},
		{
			name: "missing copy number",
			copy: &BookCopy{
				BookID:     "book-1",
				CopyNumber: "",
				Status:     CopyStatusAvailable,
			},
			wantErr: true,
			errMsg:  "copy number is required",
		},
		{
			name: "missing status",
			copy: &BookCopy{
				BookID:     "book-1",
				CopyNumber: "COPY-001",
				Status:     "",
			},
			wantErr: true,
			errMsg:  "status is required",
		},
		{
			name: "invalid status",
			copy: &BookCopy{
				BookID:     "book-1",
				CopyNumber: "COPY-001",
				Status:     "invalid",
			},
			wantErr: true,
			errMsg:  "invalid copy status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.copy.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBookCopy_StatusChecks(t *testing.T) {
	tests := []struct {
		name               string
		status             CopyStatus
		expectedAvailable  bool
		expectedBorrowed   bool
		expectedReserved   bool
		expectedDamaged    bool
		expectedLost       bool
		expectedCanBorrow  bool
	}{
		{
			name:               "available",
			status:             CopyStatusAvailable,
			expectedAvailable:  true,
			expectedBorrowed:   false,
			expectedReserved:   false,
			expectedDamaged:    false,
			expectedLost:       false,
			expectedCanBorrow:  true,
		},
		{
			name:               "borrowed",
			status:             CopyStatusBorrowed,
			expectedAvailable:  false,
			expectedBorrowed:   true,
			expectedReserved:   false,
			expectedDamaged:    false,
			expectedLost:       false,
			expectedCanBorrow:  false,
		},
		{
			name:               "reserved",
			status:             CopyStatusReserved,
			expectedAvailable:  false,
			expectedBorrowed:   false,
			expectedReserved:   true,
			expectedDamaged:    false,
			expectedLost:       false,
			expectedCanBorrow:  false,
		},
		{
			name:               "damaged",
			status:             CopyStatusDamaged,
			expectedAvailable:  false,
			expectedBorrowed:   false,
			expectedReserved:   false,
			expectedDamaged:    true,
			expectedLost:       false,
			expectedCanBorrow:  false,
		},
		{
			name:               "lost",
			status:             CopyStatusLost,
			expectedAvailable:  false,
			expectedBorrowed:   false,
			expectedReserved:   false,
			expectedDamaged:    false,
			expectedLost:       true,
			expectedCanBorrow:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copy := &BookCopy{Status: tt.status}
			
			assert.Equal(t, tt.expectedAvailable, copy.IsAvailable())
			assert.Equal(t, tt.expectedBorrowed, copy.IsBorrowed())
			assert.Equal(t, tt.expectedReserved, copy.IsReserved())
			assert.Equal(t, tt.expectedDamaged, copy.IsDamaged())
			assert.Equal(t, tt.expectedLost, copy.IsLost())
			assert.Equal(t, tt.expectedCanBorrow, copy.CanBeBorrowed())
		})
	}
}

func TestBookCopy_MarkAsBorrowed(t *testing.T) {
	t.Run("success - from available", func(t *testing.T) {
		copy := &BookCopy{
			ID:     "1",
			Status: CopyStatusAvailable,
		}

		err := copy.MarkAsBorrowed()
		
		require.NoError(t, err)
		assert.Equal(t, CopyStatusBorrowed, copy.Status)
	})

	t.Run("error - not available", func(t *testing.T) {
		copy := &BookCopy{
			ID:     "1",
			Status: CopyStatusBorrowed,
		}

		err := copy.MarkAsBorrowed()
		
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not available for borrowing")
	})
}

func TestBookCopy_MarkAsAvailable(t *testing.T) {
	t.Run("success - from borrowed", func(t *testing.T) {
		copy := &BookCopy{
			ID:     "1",
			Status: CopyStatusBorrowed,
		}

		err := copy.MarkAsAvailable()
		
		require.NoError(t, err)
		assert.Equal(t, CopyStatusAvailable, copy.Status)
	})

	t.Run("error - from lost", func(t *testing.T) {
		copy := &BookCopy{
			ID:     "1",
			Status: CopyStatusLost,
		}

		err := copy.MarkAsAvailable()
		
		require.Error(t, err)
		assert.Contains(t, err.Error(), "lost book copy cannot be marked as available")
	})
}

func TestBookCopy_MarkAsReserved(t *testing.T) {
	t.Run("success - from available", func(t *testing.T) {
		copy := &BookCopy{
			ID:     "1",
			Status: CopyStatusAvailable,
		}

		err := copy.MarkAsReserved()
		
		require.NoError(t, err)
		assert.Equal(t, CopyStatusReserved, copy.Status)
	})

	t.Run("error - not available", func(t *testing.T) {
		copy := &BookCopy{
			ID:     "1",
			Status: CopyStatusBorrowed,
		}

		err := copy.MarkAsReserved()
		
		require.Error(t, err)
		assert.Contains(t, err.Error(), "only available book copies can be reserved")
	})
}

func TestBookCopy_MarkAsDamaged(t *testing.T) {
	copy := &BookCopy{
		ID:     "1",
		Status: CopyStatusAvailable,
	}

	err := copy.MarkAsDamaged()
	
	require.NoError(t, err)
	assert.Equal(t, CopyStatusDamaged, copy.Status)
}

func TestBookCopy_MarkAsLost(t *testing.T) {
	copy := &BookCopy{
		ID:     "1",
		Status: CopyStatusBorrowed,
	}

	err := copy.MarkAsLost()
	
	require.NoError(t, err)
	assert.Equal(t, CopyStatusLost, copy.Status)
}