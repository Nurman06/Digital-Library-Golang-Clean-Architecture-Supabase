// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/repository/testhelper"
	"github.com/google/uuid"
)

func TestPostgresUserRepository_Create_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresUserRepository(testDB.DB)
	ctx := context.Background()

	tests := []struct {
		name    string
		user    *entity.User
		wantErr bool
	}{
		{
			name: "create valid user",
			user: &entity.User{
				ID:             uuid.New().String(),
				Email:          "test@example.com",
				FullName:       "Test User",
				Role:           entity.RoleMember,
				Status:         entity.StatusActive,
				BorrowingLimit: 3,
			},
			wantErr: false,
		},
		{
			name: "create user with duplicate email should fail",
			user: &entity.User{
				ID:             uuid.New().String(),
				Email:          "test@example.com", // Same email
				FullName:       "Another User",
				Role:           entity.RoleMember,
				Status:         entity.StatusActive,
				BorrowingLimit: 3,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify user was created
				retrieved, err := repo.GetByID(ctx, tt.user.ID)
				if err != nil {
					t.Errorf("Failed to retrieve created user: %v", err)
					return
				}

				if retrieved.Email != tt.user.Email {
					t.Errorf("Email mismatch: got %v, want %v", retrieved.Email, tt.user.Email)
				}
				if retrieved.FullName != tt.user.FullName {
					t.Errorf("FullName mismatch: got %v, want %v", retrieved.FullName, tt.user.FullName)
				}
			}
		})
	}
}

func TestPostgresUserRepository_GetByID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresUserRepository(testDB.DB)
	ctx := context.Background()

	// Create a test user
	userID := uuid.New().String()
	user := &entity.User{
		ID:             userID,
		Email:          "test@example.com",
		FullName:       "Test User",
		Role:           entity.RoleMember,
		Status:         entity.StatusActive,
		BorrowingLimit: 3,
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test GetByID
	retrieved, err := repo.GetByID(ctx, userID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
		return
	}

	if retrieved.ID != userID {
		t.Errorf("ID mismatch: got %v, want %v", retrieved.ID, userID)
	}

	// Test with non-existent ID
	_, err = repo.GetByID(ctx, uuid.New().String())
	if err == nil {
		t.Error("GetByID() should return error for non-existent user")
	}
}

func TestPostgresUserRepository_GetByEmail_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresUserRepository(testDB.DB)
	ctx := context.Background()

	// Create a test user
	email := "test@example.com"
	user := &entity.User{
		ID:             uuid.New().String(),
		Email:          email,
		FullName:       "Test User",
		Role:           entity.RoleMember,
		Status:         entity.StatusActive,
		BorrowingLimit: 3,
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test GetByEmail
	retrieved, err := repo.GetByEmail(ctx, email)
	if err != nil {
		t.Errorf("GetByEmail() error = %v", err)
		return
	}

	if retrieved.Email != email {
		t.Errorf("Email mismatch: got %v, want %v", retrieved.Email, email)
	}

	// Test with non-existent email
	_, err = repo.GetByEmail(ctx, "nonexistent@example.com")
	if err == nil {
		t.Error("GetByEmail() should return error for non-existent email")
	}
}

func TestPostgresUserRepository_Update_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresUserRepository(testDB.DB)
	ctx := context.Background()

	// Create a test user
	userID := uuid.New().String()
	user := &entity.User{
		ID:             userID,
		Email:          "original@example.com",
		FullName:       "Original Name",
		Role:           entity.RoleMember,
		Status:         entity.StatusActive,
		BorrowingLimit: 3,
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Update the user
	user.FullName = "Updated Name"
	user.BorrowingLimit = 5
	user.Role = entity.RoleLibrarian

	err = repo.Update(ctx, user)
	if err != nil {
		t.Errorf("Update() error = %v", err)
		return
	}

	// Verify update
	retrieved, err := repo.GetByID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated user: %v", err)
	}

	if retrieved.FullName != "Updated Name" {
		t.Errorf("FullName not updated: got %v, want %v", retrieved.FullName, "Updated Name")
	}
	if retrieved.BorrowingLimit != 5 {
		t.Errorf("BorrowingLimit not updated: got %v, want %v", retrieved.BorrowingLimit, 5)
	}
	if retrieved.Role != entity.RoleLibrarian {
		t.Errorf("Role not updated: got %v, want %v", retrieved.Role, entity.RoleLibrarian)
	}
}

func TestPostgresUserRepository_Delete_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresUserRepository(testDB.DB)
	ctx := context.Background()

	// Create a test user
	userID := uuid.New().String()
	user := &entity.User{
		ID:             userID,
		Email:          "delete@example.com",
		FullName:       "User To Delete",
		Role:           entity.RoleMember,
		Status:         entity.StatusActive,
		BorrowingLimit: 3,
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Delete the user
	err = repo.Delete(ctx, userID)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
		return
	}

	// Verify user is deleted
	_, err = repo.GetByID(ctx, userID)
	if err == nil {
		t.Error("GetByID() should return error for deleted user")
	}
}

func TestPostgresUserRepository_List_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresUserRepository(testDB.DB)
	ctx := context.Background()

	// Create test users
	for i := 1; i <= 5; i++ {
		user := &entity.User{
			ID:             uuid.New().String(),
			Email:          "user" + string(rune('0'+i)) + "@example.com",
			FullName:       "User " + string(rune('A'+i-1)),
			Role:           entity.RoleMember,
			Status:         entity.StatusActive,
			BorrowingLimit: 3,
		}
		if err := repo.Create(ctx, user); err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	tests := []struct {
		name      string
		params    UserListParams
		wantCount int
		wantTotal int64
	}{
		{
			name: "list first page",
			params: UserListParams{
				Page:     1,
				PageSize: 3,
			},
			wantCount: 3,
			wantTotal: 5,
		},
		{
			name: "list second page",
			params: UserListParams{
				Page:     2,
				PageSize: 3,
			},
			wantCount: 2,
			wantTotal: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, total, err := repo.List(ctx, tt.params)
			if err != nil {
				t.Errorf("List() error = %v", err)
				return
			}

			if len(users) != tt.wantCount {
				t.Errorf("List() returned %d users, want %d", len(users), tt.wantCount)
			}

			if total != tt.wantTotal {
				t.Errorf("List() total = %d, want %d", total, tt.wantTotal)
			}
		})
	}
}

func TestPostgresUserRepository_Search_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresUserRepository(testDB.DB)
	ctx := context.Background()

	// Create test users
	users := []struct {
		email    string
		fullName string
	}{
		{"john.doe@example.com", "John Doe"},
		{"jane.smith@example.com", "Jane Smith"},
		{"john.smith@example.com", "John Smith"},
	}

	for _, u := range users {
		user := &entity.User{
			ID:             uuid.New().String(),
			Email:          u.email,
			FullName:       u.fullName,
			Role:           entity.RoleMember,
			Status:         entity.StatusActive,
			BorrowingLimit: 3,
		}
		if err := repo.Create(ctx, user); err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	tests := []struct {
		name      string
		query     string
		wantCount int
	}{
		{
			name:      "search by name 'John'",
			query:     "John",
			wantCount: 2,
		},
		{
			name:      "search by email 'smith'",
			query:     "smith",
			wantCount: 2,
		},
		{
			name:      "search with no results",
			query:     "nonexistent",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := UserListParams{
				Page:     1,
				PageSize: 10,
			}

			users, total, err := repo.Search(ctx, tt.query, params)
			if err != nil {
				t.Errorf("Search() error = %v", err)
				return
			}

			if len(users) != tt.wantCount {
				t.Errorf("Search() returned %d users, want %d", len(users), tt.wantCount)
			}

			if total != int64(tt.wantCount) {
				t.Errorf("Search() total = %d, want %d", total, tt.wantCount)
			}
		})
	}
}

func TestPostgresUserRepository_UpdateStatus_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresUserRepository(testDB.DB)
	ctx := context.Background()

	// Create a test user
	userID := uuid.New().String()
	user := &entity.User{
		ID:             userID,
		Email:          "test@example.com",
		FullName:       "Test User",
		Role:           entity.RoleMember,
		Status:         entity.StatusActive,
		BorrowingLimit: 3,
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Update status
	err = repo.UpdateStatus(ctx, userID, entity.StatusSuspended)
	if err != nil {
		t.Errorf("UpdateStatus() error = %v", err)
		return
	}

	// Verify status update
	retrieved, err := repo.GetByID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated user: %v", err)
	}

	if retrieved.Status != entity.StatusSuspended {
		t.Errorf("Status not updated: got %v, want %v", retrieved.Status, entity.StatusSuspended)
	}
}

func TestPostgresUserRepository_UpdateBorrowingLimit_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresUserRepository(testDB.DB)
	ctx := context.Background()

	// Create a test user
	userID := uuid.New().String()
	user := &entity.User{
		ID:             userID,
		Email:          "test@example.com",
		FullName:       "Test User",
		Role:           entity.RoleMember,
		Status:         entity.StatusActive,
		BorrowingLimit: 3,
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Update borrowing limit
	newLimit := 10
	err = repo.UpdateBorrowingLimit(ctx, userID, newLimit)
	if err != nil {
		t.Errorf("UpdateBorrowingLimit() error = %v", err)
		return
	}

	// Verify limit update
	retrieved, err := repo.GetByID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated user: %v", err)
	}

	if retrieved.BorrowingLimit != newLimit {
		t.Errorf("BorrowingLimit not updated: got %v, want %v", retrieved.BorrowingLimit, newLimit)
	}
}

func TestPostgresUserRepository_ExistsByEmail_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresUserRepository(testDB.DB)
	ctx := context.Background()

	// Create a test user
	email := "test@example.com"
	user := &entity.User{
		ID:             uuid.New().String(),
		Email:          email,
		FullName:       "Test User",
		Role:           entity.RoleMember,
		Status:         entity.StatusActive,
		BorrowingLimit: 3,
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test ExistsByEmail
	exists, err := repo.ExistsByEmail(ctx, email)
	if err != nil {
		t.Errorf("ExistsByEmail() error = %v", err)
		return
	}

	if !exists {
		t.Error("ExistsByEmail() should return true for existing email")
	}

	// Test with non-existent email
	exists, err = repo.ExistsByEmail(ctx, "nonexistent@example.com")
	if err != nil {
		t.Errorf("ExistsByEmail() error = %v", err)
		return
	}

	if exists {
		t.Error("ExistsByEmail() should return false for non-existent email")
	}
}

func TestPostgresUserRepository_Count_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresUserRepository(testDB.DB)
	ctx := context.Background()

	// Create test users
	for i := 1; i <= 3; i++ {
		user := &entity.User{
			ID:             uuid.New().String(),
			Email:          "user" + string(rune('0'+i)) + "@example.com",
			FullName:       "User " + string(rune('A'+i-1)),
			Role:           entity.RoleMember,
			Status:         entity.StatusActive,
			BorrowingLimit: 3,
		}
		if err := repo.Create(ctx, user); err != nil {
			t.Fatalf("Failed to create test user: %v", err)
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

func TestPostgresUserRepository_GetByRole_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresUserRepository(testDB.DB)
	ctx := context.Background()

	// Create test users with different roles
	roles := []entity.UserRole{
		entity.RoleMember,
		entity.RoleMember,
		entity.RoleLibrarian,
		entity.RoleAdmin,
	}

	for i, role := range roles {
		user := &entity.User{
			ID:             uuid.New().String(),
			Email:          "user" + string(rune('0'+i)) + "@example.com",
			FullName:       "User " + string(rune('A'+i)),
			Role:           role,
			Status:         entity.StatusActive,
			BorrowingLimit: 3,
		}
		if err := repo.Create(ctx, user); err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	// Test GetByRole
	params := UserListParams{
		Page:     1,
		PageSize: 10,
	}

	members, total, err := repo.GetByRole(ctx, entity.RoleMember, params)
	if err != nil {
		t.Errorf("GetByRole() error = %v", err)
		return
	}

	if len(members) != 2 {
		t.Errorf("GetByRole(Member) returned %d users, want %d", len(members), 2)
	}
	if total != 2 {
		t.Errorf("GetByRole(Member) total = %d, want %d", total, 2)
	}

	// Verify all returned users have the correct role
	for _, user := range members {
		if user.Role != entity.RoleMember {
			t.Errorf("Expected role %v, got %v", entity.RoleMember, user.Role)
		}
	}
}

func TestPostgresUserRepository_GetByStatus_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelper.SetupTestDB(t)
	defer testDB.Close()

	testDB.RunMigrations(t)
	defer testDB.DropTables(t)
	testDB.CleanupTables(t)

	repo := NewPostgresUserRepository(testDB.DB)
	ctx := context.Background()

	// Create test users with different statuses
	statuses := []entity.UserStatus{
		entity.StatusActive,
		entity.StatusActive,
		entity.StatusSuspended,
	}

	for i, status := range statuses {
		user := &entity.User{
			ID:             uuid.New().String(),
			Email:          "user" + string(rune('0'+i)) + "@example.com",
			FullName:       "User " + string(rune('A'+i)),
			Role:           entity.RoleMember,
			Status:         status,
			BorrowingLimit: 3,
		}
		if err := repo.Create(ctx, user); err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	// Test GetByStatus
	params := UserListParams{
		Page:     1,
		PageSize: 10,
	}

	activeUsers, total, err := repo.GetByStatus(ctx, entity.StatusActive, params)
	if err != nil {
		t.Errorf("GetByStatus() error = %v", err)
		return
	}

	if len(activeUsers) != 2 {
		t.Errorf("GetByStatus(Active) returned %d users, want %d", len(activeUsers), 2)
	}
	if total != 2 {
		t.Errorf("GetByStatus(Active) total = %d, want %d", total, 2)
	}

	// Verify all returned users have the correct status
	for _, user := range activeUsers {
		if user.Status != entity.StatusActive {
			t.Errorf("Expected status %v, got %v", entity.StatusActive, user.Status)
		}
	}
}