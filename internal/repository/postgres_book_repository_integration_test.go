// +build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/repository/testhelper"
	"github.com/google/uuid"
)

func TestPostgresBookRepository_Create_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)

	repo := NewPostgresBookRepository(testDB.DB)
	ctx := context.Background()

	tests := []struct {
		name    string
		book    *entity.Book
		wantErr bool
	}{
		{
			name: "create valid book",
			book: &entity.Book{
				ID:              uuid.New().String(),
				Title:           "Clean Code",
				Author:          "Robert C. Martin",
				ISBN:            "9780132350884",
				Category:        "Programming",
				PublicationYear: 2008,
				Description:     "A Handbook of Agile Software Craftsmanship",
			},
			wantErr: false,
		},
		{
			name: "create book with duplicate ISBN should fail",
			book: &entity.Book{
				ID:              uuid.New().String(),
				Title:           "Another Book",
				Author:          "Another Author",
				ISBN:            "9780132350884", // Same ISBN as above
				Category:        "Programming",
				PublicationYear: 2020,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.book)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify book was created
				retrieved, err := repo.GetByID(ctx, tt.book.ID)
				if err != nil {
					t.Errorf("Failed to retrieve created book: %v", err)
					return
				}

				if retrieved.Title != tt.book.Title {
					t.Errorf("Title mismatch: got %v, want %v", retrieved.Title, tt.book.Title)
				}
				if retrieved.ISBN != tt.book.ISBN {
					t.Errorf("ISBN mismatch: got %v, want %v", retrieved.ISBN, tt.book.ISBN)
				}
			}
		})
	}
}

func TestPostgresBookRepository_GetByID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresBookRepository(testDB.DB)
	ctx := context.Background()

	// Create a test book
	bookID := uuid.New().String()
	book := &entity.Book{
		ID:              bookID,
		Title:           "Test Book",
		Author:          "Test Author",
		ISBN:            "9781234567890",
		Category:        "Test",
		PublicationYear: 2023,
		Description:     "Test Description",
	}

	err := repo.Create(ctx, book)
	if err != nil {
		t.Fatalf("Failed to create test book: %v", err)
	}

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "get existing book",
			id:      bookID,
			wantErr: false,
		},
		{
			name:    "get non-existent book",
			id:      uuid.New().String(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retrieved, err := repo.GetByID(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && retrieved == nil {
				t.Error("GetByID() returned nil book")
			}

			if !tt.wantErr && retrieved.ID != tt.id {
				t.Errorf("ID mismatch: got %v, want %v", retrieved.ID, tt.id)
			}
		})
	}
}

func TestPostgresBookRepository_Update_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresBookRepository(testDB.DB)
	ctx := context.Background()

	// Create a test book
	bookID := uuid.New().String()
	book := &entity.Book{
		ID:              bookID,
		Title:           "Original Title",
		Author:          "Original Author",
		ISBN:            "9781234567890",
		Category:        "Original",
		PublicationYear: 2023,
		Description:     "Original Description",
	}

	err := repo.Create(ctx, book)
	if err != nil {
		t.Fatalf("Failed to create test book: %v", err)
	}

	// Update the book
	book.Title = "Updated Title"
	book.Description = "Updated Description"
	book.PublicationYear = 2024

	err = repo.Update(ctx, book)
	if err != nil {
		t.Errorf("Update() error = %v", err)
		return
	}

	// Verify update
	retrieved, err := repo.GetByID(ctx, bookID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated book: %v", err)
	}

	if retrieved.Title != "Updated Title" {
		t.Errorf("Title not updated: got %v, want %v", retrieved.Title, "Updated Title")
	}
	if retrieved.Description != "Updated Description" {
		t.Errorf("Description not updated: got %v, want %v", retrieved.Description, "Updated Description")
	}
	if retrieved.PublicationYear != 2024 {
		t.Errorf("PublicationYear not updated: got %v, want %v", retrieved.PublicationYear, 2024)
	}
}

func TestPostgresBookRepository_Delete_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresBookRepository(testDB.DB)
	ctx := context.Background()

	// Create a test book
	bookID := uuid.New().String()
	book := &entity.Book{
		ID:              bookID,
		Title:           "Book to Delete",
		Author:          "Test Author",
		ISBN:            "9781234567890",
		Category:        "Test",
		PublicationYear: 2023,
	}

	err := repo.Create(ctx, book)
	if err != nil {
		t.Fatalf("Failed to create test book: %v", err)
	}

	// Soft delete the book
	err = repo.Delete(ctx, bookID)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
		return
	}

	// Verify book is soft deleted (should not be found)
	_, err = repo.GetByID(ctx, bookID)
	if err == nil {
		t.Error("GetByID() should return error for soft-deleted book")
	}
}

func TestPostgresBookRepository_List_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresBookRepository(testDB.DB)
	ctx := context.Background()

	// Create multiple test books
	for i := 1; i <= 5; i++ {
		book := &entity.Book{
			ID:              uuid.New().String(),
			Title:           "Book " + string(rune('A'+i-1)),
			Author:          "Author " + string(rune('A'+i-1)),
			ISBN:            "978123456789" + string(rune('0'+i)),
			Category:        "Test",
			PublicationYear: 2020 + i,
		}
		if err := repo.Create(ctx, book); err != nil {
			t.Fatalf("Failed to create test book: %v", err)
		}
		time.Sleep(10 * time.Millisecond) // Ensure different created_at timestamps
	}

	tests := []struct {
		name      string
		params    ListParams
		wantCount int
		wantTotal int64
	}{
		{
			name: "list first page",
			params: ListParams{
				Page:     1,
				PageSize: 3,
			},
			wantCount: 3,
			wantTotal: 5,
		},
		{
			name: "list second page",
			params: ListParams{
				Page:     2,
				PageSize: 3,
			},
			wantCount: 2,
			wantTotal: 5,
		},
		{
			name: "list with sorting by title asc",
			params: ListParams{
				Page:      1,
				PageSize:  10,
				SortBy:    "title",
				SortOrder: "ASC",
			},
			wantCount: 5,
			wantTotal: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			books, total, err := repo.List(ctx, tt.params)
			if err != nil {
				t.Errorf("List() error = %v", err)
				return
			}

			if len(books) != tt.wantCount {
				t.Errorf("List() returned %d books, want %d", len(books), tt.wantCount)
			}

			if total != tt.wantTotal {
				t.Errorf("List() total = %d, want %d", total, tt.wantTotal)
			}

			// Verify sorting if specified
			if tt.params.SortBy == "title" && tt.params.SortOrder == "ASC" && len(books) > 1 {
				for i := 0; i < len(books)-1; i++ {
					if books[i].Title > books[i+1].Title {
						t.Errorf("Books not sorted by title ascending")
						break
					}
				}
			}
		})
	}
}

func TestPostgresBookRepository_Search_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresBookRepository(testDB.DB)
	ctx := context.Background()

	// Create test books
	testBooks := []struct {
		title    string
		author   string
		isbn     string
		category string
		year     int
	}{
		{"Clean Code", "Robert Martin", "9780132350884", "Programming", 2008},
		{"The Clean Coder", "Robert Martin", "9780137081073", "Programming", 2011},
		{"Design Patterns", "Gang of Four", "9780201633610", "Programming", 1994},
		{"Introduction to Algorithms", "CLRS", "9780262033848", "Computer Science", 2009},
	}

	for _, tb := range testBooks {
		book := &entity.Book{
			ID:              uuid.New().String(),
			Title:           tb.title,
			Author:          tb.author,
			ISBN:            tb.isbn,
			Category:        tb.category,
			PublicationYear: tb.year,
		}
		if err := repo.Create(ctx, book); err != nil {
			t.Fatalf("Failed to create test book: %v", err)
		}
	}

	tests := []struct {
		name      string
		params    SearchParams
		wantCount int
	}{
		{
			name: "search by query 'Clean'",
			params: SearchParams{
				Query:    "Clean",
				Page:     1,
				PageSize: 10,
			},
			wantCount: 2,
		},
		{
			name: "search by author 'Martin'",
			params: SearchParams{
				Author:   "Martin",
				Page:     1,
				PageSize: 10,
			},
			wantCount: 2,
		},
		{
			name: "search by category 'Programming'",
			params: SearchParams{
				Category: "Programming",
				Page:     1,
				PageSize: 10,
			},
			wantCount: 3,
		},
		{
			name: "search by ISBN",
			params: SearchParams{
				ISBN:     "9780132350884",
				Page:     1,
				PageSize: 10,
			},
			wantCount: 1,
		},
		{
			name: "search by year range",
			params: SearchParams{
				YearFrom: 2000,
				YearTo:   2010,
				Page:     1,
				PageSize: 10,
			},
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			books, total, err := repo.Search(ctx, tt.params)
			if err != nil {
				t.Errorf("Search() error = %v", err)
				return
			}

			if len(books) != tt.wantCount {
				t.Errorf("Search() returned %d books, want %d", len(books), tt.wantCount)
			}

			if total != int64(tt.wantCount) {
				t.Errorf("Search() total = %d, want %d", total, tt.wantCount)
			}
		})
	}
}

func TestPostgresBookRepository_GetByISBN_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresBookRepository(testDB.DB)
	ctx := context.Background()

	// Create a test book
	isbn := "9780132350884"
	book := &entity.Book{
		ID:              uuid.New().String(),
		Title:           "Clean Code",
		Author:          "Robert Martin",
		ISBN:            isbn,
		Category:        "Programming",
		PublicationYear: 2008,
	}

	err := repo.Create(ctx, book)
	if err != nil {
		t.Fatalf("Failed to create test book: %v", err)
	}

	// Test GetByISBN
	retrieved, err := repo.GetByISBN(ctx, isbn)
	if err != nil {
		t.Errorf("GetByISBN() error = %v", err)
		return
	}

	if retrieved.ISBN != isbn {
		t.Errorf("ISBN mismatch: got %v, want %v", retrieved.ISBN, isbn)
	}
	if retrieved.Title != book.Title {
		t.Errorf("Title mismatch: got %v, want %v", retrieved.Title, book.Title)
	}

	// Test with non-existent ISBN
	_, err = repo.GetByISBN(ctx, "9999999999999")
	if err == nil {
		t.Error("GetByISBN() should return error for non-existent ISBN")
	}
}

func TestPostgresBookRepository_ExistsByISBN_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresBookRepository(testDB.DB)
	ctx := context.Background()

	// Create a test book
	isbn := "9780132350884"
	book := &entity.Book{
		ID:              uuid.New().String(),
		Title:           "Clean Code",
		Author:          "Robert Martin",
		ISBN:            isbn,
		Category:        "Programming",
		PublicationYear: 2008,
	}

	err := repo.Create(ctx, book)
	if err != nil {
		t.Fatalf("Failed to create test book: %v", err)
	}

	// Test ExistsByISBN
	exists, err := repo.ExistsByISBN(ctx, isbn)
	if err != nil {
		t.Errorf("ExistsByISBN() error = %v", err)
		return
	}

	if !exists {
		t.Error("ExistsByISBN() should return true for existing ISBN")
	}

	// Test with non-existent ISBN
	exists, err = repo.ExistsByISBN(ctx, "9999999999999")
	if err != nil {
		t.Errorf("ExistsByISBN() error = %v", err)
		return
	}

	if exists {
		t.Error("ExistsByISBN() should return false for non-existent ISBN")
	}
}

func TestPostgresBookRepository_Count_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresBookRepository(testDB.DB)
	ctx := context.Background()

	// Create test books
	for i := 1; i <= 3; i++ {
		book := &entity.Book{
			ID:              uuid.New().String(),
			Title:           "Book " + string(rune('A'+i-1)),
			Author:          "Author " + string(rune('A'+i-1)),
			ISBN:            "978123456789" + string(rune('0'+i)),
			Category:        "Test",
			PublicationYear: 2023,
		}
		if err := repo.Create(ctx, book); err != nil {
			t.Fatalf("Failed to create test book: %v", err)
		}
	}

	// Test Count
	count, err := repo.Count(ctx)
	if err != nil {
		t.Errorf("Count() error = %v", err)
		return
	}

	if count != 3 {
		t.Errorf("Count() = %d, want %d", count, 3)
	}
}