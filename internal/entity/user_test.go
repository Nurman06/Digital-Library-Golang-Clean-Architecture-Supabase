package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid user",
			user: &User{
				Email:          "john.doe@example.com",
				FullName:       "John Doe",
				Role:           RoleMember,
				Status:         StatusActive,
				BorrowingLimit: 5,
			},
			wantErr: false,
		},
		{
			name: "missing email",
			user: &User{
				Email:          "",
				FullName:       "John Doe",
				Role:           RoleMember,
				Status:         StatusActive,
				BorrowingLimit: 5,
			},
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name: "invalid email format",
			user: &User{
				Email:          "invalid-email",
				FullName:       "John Doe",
				Role:           RoleMember,
				Status:         StatusActive,
				BorrowingLimit: 5,
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "missing full name",
			user: &User{
				Email:          "john.doe@example.com",
				FullName:       "",
				Role:           RoleMember,
				Status:         StatusActive,
				BorrowingLimit: 5,
			},
			wantErr: true,
			errMsg:  "full name is required",
		},
		{
			name: "missing role",
			user: &User{
				Email:          "john.doe@example.com",
				FullName:       "John Doe",
				Role:           "",
				Status:         StatusActive,
				BorrowingLimit: 5,
			},
			wantErr: true,
			errMsg:  "role is required",
		},
		{
			name: "invalid role",
			user: &User{
				Email:          "john.doe@example.com",
				FullName:       "John Doe",
				Role:           "invalid",
				Status:         StatusActive,
				BorrowingLimit: 5,
			},
			wantErr: true,
			errMsg:  "invalid user role",
		},
		{
			name: "missing status",
			user: &User{
				Email:          "john.doe@example.com",
				FullName:       "John Doe",
				Role:           RoleMember,
				Status:         "",
				BorrowingLimit: 5,
			},
			wantErr: true,
			errMsg:  "status is required",
		},
		{
			name: "invalid status",
			user: &User{
				Email:          "john.doe@example.com",
				FullName:       "John Doe",
				Role:           RoleMember,
				Status:         "invalid",
				BorrowingLimit: 5,
			},
			wantErr: true,
			errMsg:  "invalid user status",
		},
		{
			name: "negative borrowing limit",
			user: &User{
				Email:          "john.doe@example.com",
				FullName:       "John Doe",
				Role:           RoleMember,
				Status:         StatusActive,
				BorrowingLimit: -1,
			},
			wantErr: true,
			errMsg:  "borrowing limit cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUser_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   UserStatus
		expected bool
	}{
		{
			name:     "active user",
			status:   StatusActive,
			expected: true,
		},
		{
			name:     "suspended user",
			status:   StatusSuspended,
			expected: false,
		},
		{
			name:     "expired user",
			status:   StatusExpired,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Status: tt.status}
			assert.Equal(t, tt.expected, user.IsActive())
		})
	}
}

func TestUser_CanBorrow(t *testing.T) {
	tests := []struct {
		name     string
		status   UserStatus
		expected bool
	}{
		{
			name:     "active user can borrow",
			status:   StatusActive,
			expected: true,
		},
		{
			name:     "suspended user cannot borrow",
			status:   StatusSuspended,
			expected: false,
		},
		{
			name:     "expired user cannot borrow",
			status:   StatusExpired,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Status: tt.status}
			assert.Equal(t, tt.expected, user.CanBorrow())
		})
	}
}

func TestGetDefaultBorrowingLimit(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected int
	}{
		{
			name:     "admin has unlimited",
			role:     RoleAdmin,
			expected: 0,
		},
		{
			name:     "librarian has 10",
			role:     RoleLibrarian,
			expected: 10,
		},
		{
			name:     "member has 5",
			role:     RoleMember,
			expected: 5,
		},
		{
			name:     "invalid role defaults to 5",
			role:     "invalid",
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit := GetDefaultBorrowingLimit(tt.role)
			assert.Equal(t, tt.expected, limit)
		})
	}
}

func TestUser_HasPermission(t *testing.T) {
	tests := []struct {
		name      string
		role      UserRole
		operation string
		expected  bool
	}{
		// Admin tests
		{
			name:      "admin can manage users",
			role:      RoleAdmin,
			operation: "manage_users",
			expected:  true,
		},
		{
			name:      "admin can view books",
			role:      RoleAdmin,
			operation: "view_books",
			expected:  true,
		},
		// Librarian tests
		{
			name:      "librarian cannot manage users",
			role:      RoleLibrarian,
			operation: "manage_users",
			expected:  false,
		},
		{
			name:      "librarian can view books",
			role:      RoleLibrarian,
			operation: "view_books",
			expected:  true,
		},
		// Member tests
		{
			name:      "member can view books",
			role:      RoleMember,
			operation: "view_books",
			expected:  true,
		},
		{
			name:      "member can borrow books",
			role:      RoleMember,
			operation: "borrow_books",
			expected:  true,
		},
		{
			name:      "member can view own profile",
			role:      RoleMember,
			operation: "view_own_profile",
			expected:  true,
		},
		{
			name:      "member cannot manage users",
			role:      RoleMember,
			operation: "manage_users",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			assert.Equal(t, tt.expected, user.HasPermission(tt.operation))
		})
	}
}

func TestUser_RoleChecks(t *testing.T) {
	tests := []struct {
		name            string
		role            UserRole
		expectedAdmin   bool
		expectedLibrary bool
		expectedMember  bool
	}{
		{
			name:            "admin role",
			role:            RoleAdmin,
			expectedAdmin:   true,
			expectedLibrary: false,
			expectedMember:  false,
		},
		{
			name:            "librarian role",
			role:            RoleLibrarian,
			expectedAdmin:   false,
			expectedLibrary: true,
			expectedMember:  false,
		},
		{
			name:            "member role",
			role:            RoleMember,
			expectedAdmin:   false,
			expectedLibrary: false,
			expectedMember:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			assert.Equal(t, tt.expectedAdmin, user.IsAdmin())
			assert.Equal(t, tt.expectedLibrary, user.IsLibrarian())
			assert.Equal(t, tt.expectedMember, user.IsMember())
		})
	}
}

func TestUser_CanManageBooks(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{
			name:     "admin can manage books",
			role:     RoleAdmin,
			expected: true,
		},
		{
			name:     "librarian can manage books",
			role:     RoleLibrarian,
			expected: true,
		},
		{
			name:     "member cannot manage books",
			role:     RoleMember,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			assert.Equal(t, tt.expected, user.CanManageBooks())
		})
	}
}

func TestUser_CanManageUsers(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{
			name:     "admin can manage users",
			role:     RoleAdmin,
			expected: true,
		},
		{
			name:     "librarian cannot manage users",
			role:     RoleLibrarian,
			expected: false,
		},
		{
			name:     "member cannot manage users",
			role:     RoleMember,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			assert.Equal(t, tt.expected, user.CanManageUsers())
		})
	}
}

func TestUser_CanDeleteBooks(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{
			name:     "admin can delete books",
			role:     RoleAdmin,
			expected: true,
		},
		{
			name:     "librarian cannot delete books",
			role:     RoleLibrarian,
			expected: false,
		},
		{
			name:     "member cannot delete books",
			role:     RoleMember,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			assert.Equal(t, tt.expected, user.CanDeleteBooks())
		})
	}
}

func TestUser_Suspend(t *testing.T) {
	t.Run("success - suspend active user", func(t *testing.T) {
		user := &User{
			ID:     "1",
			Status: StatusActive,
		}

		err := user.Suspend()
		
		require.NoError(t, err)
		assert.Equal(t, StatusSuspended, user.Status)
	})

	t.Run("error - already suspended", func(t *testing.T) {
		user := &User{
			ID:     "1",
			Status: StatusSuspended,
		}

		err := user.Suspend()
		
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already suspended")
	})
}

func TestUser_Activate(t *testing.T) {
	t.Run("success - activate suspended user", func(t *testing.T) {
		user := &User{
			ID:     "1",
			Status: StatusSuspended,
		}

		err := user.Activate()
		
		require.NoError(t, err)
		assert.Equal(t, StatusActive, user.Status)
	})

	t.Run("error - already active", func(t *testing.T) {
		user := &User{
			ID:     "1",
			Status: StatusActive,
		}

		err := user.Activate()
		
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already active")
	})
}

func TestUser_Expire(t *testing.T) {
	t.Run("success - expire active user", func(t *testing.T) {
		user := &User{
			ID:     "1",
			Status: StatusActive,
		}

		err := user.Expire()
		
		require.NoError(t, err)
		assert.Equal(t, StatusExpired, user.Status)
	})

	t.Run("error - already expired", func(t *testing.T) {
		user := &User{
			ID:     "1",
			Status: StatusExpired,
		}

		err := user.Expire()
		
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already expired")
	})
}

func TestUser_UpdateBorrowingLimit(t *testing.T) {
	t.Run("success - update limit", func(t *testing.T) {
		user := &User{
			ID:             "1",
			BorrowingLimit: 5,
		}

		err := user.UpdateBorrowingLimit(10)
		
		require.NoError(t, err)
		assert.Equal(t, 10, user.BorrowingLimit)
	})

	t.Run("error - negative limit", func(t *testing.T) {
		user := &User{
			ID:             "1",
			BorrowingLimit: 5,
		}

		err := user.UpdateBorrowingLimit(-1)
		
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be negative")
		assert.Equal(t, 5, user.BorrowingLimit) // Should not change
	})
}

func TestUser_HasUnlimitedBorrowing(t *testing.T) {
	tests := []struct {
		name           string
		role           UserRole
		borrowingLimit int
		expected       bool
	}{
		{
			name:           "admin with 0 limit has unlimited",
			role:           RoleAdmin,
			borrowingLimit: 0,
			expected:       true,
		},
		{
			name:           "admin with positive limit does not have unlimited",
			role:           RoleAdmin,
			borrowingLimit: 10,
			expected:       false,
		},
		{
			name:           "non-admin with 0 limit does not have unlimited",
			role:           RoleMember,
			borrowingLimit: 0,
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Role:           tt.role,
				BorrowingLimit: tt.borrowingLimit,
			}
			assert.Equal(t, tt.expected, user.HasUnlimitedBorrowing())
		})
	}
}

func TestUser_CanBorrowMoreBooks(t *testing.T) {
	tests := []struct {
		name              string
		status            UserStatus
		role              UserRole
		borrowingLimit    int
		currentBorrowCount int
		expected          bool
	}{
		{
			name:              "active user under limit can borrow",
			status:            StatusActive,
			role:              RoleMember,
			borrowingLimit:    5,
			currentBorrowCount: 3,
			expected:          true,
		},
		{
			name:              "active user at limit cannot borrow",
			status:            StatusActive,
			role:              RoleMember,
			borrowingLimit:    5,
			currentBorrowCount: 5,
			expected:          false,
		},
		{
			name:              "suspended user cannot borrow",
			status:            StatusSuspended,
			role:              RoleMember,
			borrowingLimit:    5,
			currentBorrowCount: 0,
			expected:          false,
		},
		{
			name:              "admin with unlimited can always borrow",
			status:            StatusActive,
			role:              RoleAdmin,
			borrowingLimit:    0,
			currentBorrowCount: 100,
			expected:          true,
		},
		{
			name:              "active user over limit cannot borrow",
			status:            StatusActive,
			role:              RoleMember,
			borrowingLimit:    5,
			currentBorrowCount: 6,
			expected:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Status:         tt.status,
				Role:           tt.role,
				BorrowingLimit: tt.borrowingLimit,
			}
			assert.Equal(t, tt.expected, user.CanBorrowMoreBooks(tt.currentBorrowCount))
		})
	}
}

func TestUser_IsSuspended(t *testing.T) {
	tests := []struct {
		name     string
		status   UserStatus
		expected bool
	}{
		{
			name:     "suspended user",
			status:   StatusSuspended,
			expected: true,
		},
		{
			name:     "active user",
			status:   StatusActive,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Status: tt.status}
			assert.Equal(t, tt.expected, user.IsSuspended())
		})
	}
}

func TestUser_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		status   UserStatus
		expected bool
	}{
		{
			name:     "expired user",
			status:   StatusExpired,
			expected: true,
		},
		{
			name:     "active user",
			status:   StatusActive,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Status: tt.status}
			assert.Equal(t, tt.expected, user.IsExpired())
		})
	}
}