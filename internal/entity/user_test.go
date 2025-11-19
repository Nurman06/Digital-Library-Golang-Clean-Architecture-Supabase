package entity

import (
	"testing"
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
				Email:          "test@example.com",
				FullName:       "Test User",
				Role:           RoleMember,
				Status:         StatusActive,
				BorrowingLimit: 5,
			},
			wantErr: false,
		},
		{
			name: "missing email",
			user: &User{
				FullName:       "Test User",
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
				FullName:       "Test User",
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
				Email:          "test@example.com",
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
				Email:          "test@example.com",
				FullName:       "Test User",
				Status:         StatusActive,
				BorrowingLimit: 5,
			},
			wantErr: true,
			errMsg:  "role is required",
		},
		{
			name: "invalid role",
			user: &User{
				Email:          "test@example.com",
				FullName:       "Test User",
				Role:           "invalid-role",
				Status:         StatusActive,
				BorrowingLimit: 5,
			},
			wantErr: true,
			errMsg:  "invalid user role",
		},
		{
			name: "missing status",
			user: &User{
				Email:          "test@example.com",
				FullName:       "Test User",
				Role:           RoleMember,
				BorrowingLimit: 5,
			},
			wantErr: true,
			errMsg:  "status is required",
		},
		{
			name: "invalid status",
			user: &User{
				Email:          "test@example.com",
				FullName:       "Test User",
				Role:           RoleMember,
				Status:         "invalid-status",
				BorrowingLimit: 5,
			},
			wantErr: true,
			errMsg:  "invalid user status",
		},
		{
			name: "negative borrowing limit",
			user: &User{
				Email:          "test@example.com",
				FullName:       "Test User",
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
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("User.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUser_StatusChecks(t *testing.T) {
	tests := []struct {
		name   string
		status UserStatus
		checks map[string]bool
	}{
		{
			name:   "active user",
			status: StatusActive,
			checks: map[string]bool{
				"IsActive":    true,
				"IsSuspended": false,
				"IsExpired":   false,
				"CanBorrow":   true,
			},
		},
		{
			name:   "suspended user",
			status: StatusSuspended,
			checks: map[string]bool{
				"IsActive":    false,
				"IsSuspended": true,
				"IsExpired":   false,
				"CanBorrow":   false,
			},
		},
		{
			name:   "expired user",
			status: StatusExpired,
			checks: map[string]bool{
				"IsActive":    false,
				"IsSuspended": false,
				"IsExpired":   true,
				"CanBorrow":   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Status: tt.status}

			if got := user.IsActive(); got != tt.checks["IsActive"] {
				t.Errorf("IsActive() = %v, want %v", got, tt.checks["IsActive"])
			}
			if got := user.IsSuspended(); got != tt.checks["IsSuspended"] {
				t.Errorf("IsSuspended() = %v, want %v", got, tt.checks["IsSuspended"])
			}
			if got := user.IsExpired(); got != tt.checks["IsExpired"] {
				t.Errorf("IsExpired() = %v, want %v", got, tt.checks["IsExpired"])
			}
			if got := user.CanBorrow(); got != tt.checks["CanBorrow"] {
				t.Errorf("CanBorrow() = %v, want %v", got, tt.checks["CanBorrow"])
			}
		})
	}
}

func TestUser_RoleChecks(t *testing.T) {
	tests := []struct {
		name   string
		role   UserRole
		checks map[string]bool
	}{
		{
			name: "admin user",
			role: RoleAdmin,
			checks: map[string]bool{
				"IsAdmin":        true,
				"IsLibrarian":    false,
				"IsMember":       false,
				"CanManageBooks": true,
				"CanManageUsers": true,
				"CanDeleteBooks": true,
			},
		},
		{
			name: "librarian user",
			role: RoleLibrarian,
			checks: map[string]bool{
				"IsAdmin":        false,
				"IsLibrarian":    true,
				"IsMember":       false,
				"CanManageBooks": true,
				"CanManageUsers": false,
				"CanDeleteBooks": false,
			},
		},
		{
			name: "member user",
			role: RoleMember,
			checks: map[string]bool{
				"IsAdmin":        false,
				"IsLibrarian":    false,
				"IsMember":       true,
				"CanManageBooks": false,
				"CanManageUsers": false,
				"CanDeleteBooks": false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}

			if got := user.IsAdmin(); got != tt.checks["IsAdmin"] {
				t.Errorf("IsAdmin() = %v, want %v", got, tt.checks["IsAdmin"])
			}
			if got := user.IsLibrarian(); got != tt.checks["IsLibrarian"] {
				t.Errorf("IsLibrarian() = %v, want %v", got, tt.checks["IsLibrarian"])
			}
			if got := user.IsMember(); got != tt.checks["IsMember"] {
				t.Errorf("IsMember() = %v, want %v", got, tt.checks["IsMember"])
			}
			if got := user.CanManageBooks(); got != tt.checks["CanManageBooks"] {
				t.Errorf("CanManageBooks() = %v, want %v", got, tt.checks["CanManageBooks"])
			}
			if got := user.CanManageUsers(); got != tt.checks["CanManageUsers"] {
				t.Errorf("CanManageUsers() = %v, want %v", got, tt.checks["CanManageUsers"])
			}
			if got := user.CanDeleteBooks(); got != tt.checks["CanDeleteBooks"] {
				t.Errorf("CanDeleteBooks() = %v, want %v", got, tt.checks["CanDeleteBooks"])
			}
		})
	}
}

func TestUser_HasPermission(t *testing.T) {
	tests := []struct {
		name      string
		role      UserRole
		operation string
		want      bool
	}{
		// Admin permissions
		{
			name:      "admin can manage users",
			role:      RoleAdmin,
			operation: "manage_users",
			want:      true,
		},
		{
			name:      "admin can view books",
			role:      RoleAdmin,
			operation: "view_books",
			want:      true,
		},
		{
			name:      "admin can borrow books",
			role:      RoleAdmin,
			operation: "borrow_books",
			want:      true,
		},
		// Librarian permissions
		{
			name:      "librarian cannot manage users",
			role:      RoleLibrarian,
			operation: "manage_users",
			want:      false,
		},
		{
			name:      "librarian can view books",
			role:      RoleLibrarian,
			operation: "view_books",
			want:      true,
		},
		{
			name:      "librarian can borrow books",
			role:      RoleLibrarian,
			operation: "borrow_books",
			want:      true,
		},
		// Member permissions
		{
			name:      "member cannot manage users",
			role:      RoleMember,
			operation: "manage_users",
			want:      false,
		},
		{
			name:      "member can view books",
			role:      RoleMember,
			operation: "view_books",
			want:      true,
		},
		{
			name:      "member can borrow books",
			role:      RoleMember,
			operation: "borrow_books",
			want:      true,
		},
		{
			name:      "member can view own profile",
			role:      RoleMember,
			operation: "view_own_profile",
			want:      true,
		},
		{
			name:      "member cannot perform other operations",
			role:      RoleMember,
			operation: "delete_books",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			if got := user.HasPermission(tt.operation); got != tt.want {
				t.Errorf("HasPermission(%v) = %v, want %v", tt.operation, got, tt.want)
			}
		})
	}
}

func TestGetDefaultBorrowingLimit(t *testing.T) {
	tests := []struct {
		name string
		role UserRole
		want int
	}{
		{
			name: "admin unlimited",
			role: RoleAdmin,
			want: 0,
		},
		{
			name: "librarian limit 10",
			role: RoleLibrarian,
			want: 10,
		},
		{
			name: "member limit 5",
			role: RoleMember,
			want: 5,
		},
		{
			name: "invalid role defaults to 5",
			role: "invalid",
			want: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDefaultBorrowingLimit(tt.role); got != tt.want {
				t.Errorf("GetDefaultBorrowingLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_Suspend(t *testing.T) {
	tests := []struct {
		name    string
		status  UserStatus
		wantErr bool
	}{
		{
			name:    "active user can be suspended",
			status:  StatusActive,
			wantErr: false,
		},
		{
			name:    "expired user can be suspended",
			status:  StatusExpired,
			wantErr: false,
		},
		{
			name:    "already suspended user",
			status:  StatusSuspended,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Status: tt.status}
			err := user.Suspend()

			if (err != nil) != tt.wantErr {
				t.Errorf("Suspend() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && user.Status != StatusSuspended {
				t.Errorf("Status should be suspended, got %v", user.Status)
			}
		})
	}
}

func TestUser_Activate(t *testing.T) {
	tests := []struct {
		name    string
		status  UserStatus
		wantErr bool
	}{
		{
			name:    "suspended user can be activated",
			status:  StatusSuspended,
			wantErr: false,
		},
		{
			name:    "expired user can be activated",
			status:  StatusExpired,
			wantErr: false,
		},
		{
			name:    "already active user",
			status:  StatusActive,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Status: tt.status}
			err := user.Activate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Activate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && user.Status != StatusActive {
				t.Errorf("Status should be active, got %v", user.Status)
			}
		})
	}
}

func TestUser_Expire(t *testing.T) {
	tests := []struct {
		name    string
		status  UserStatus
		wantErr bool
	}{
		{
			name:    "active user can be expired",
			status:  StatusActive,
			wantErr: false,
		},
		{
			name:    "suspended user can be expired",
			status:  StatusSuspended,
			wantErr: false,
		},
		{
			name:    "already expired user",
			status:  StatusExpired,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Status: tt.status}
			err := user.Expire()

			if (err != nil) != tt.wantErr {
				t.Errorf("Expire() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && user.Status != StatusExpired {
				t.Errorf("Status should be expired, got %v", user.Status)
			}
		})
	}
}

func TestUser_UpdateBorrowingLimit(t *testing.T) {
	tests := []struct {
		name     string
		newLimit int
		wantErr  bool
	}{
		{
			name:     "valid limit",
			newLimit: 10,
			wantErr:  false,
		},
		{
			name:     "zero limit",
			newLimit: 0,
			wantErr:  false,
		},
		{
			name:     "negative limit",
			newLimit: -1,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{BorrowingLimit: 5}
			err := user.UpdateBorrowingLimit(tt.newLimit)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateBorrowingLimit() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && user.BorrowingLimit != tt.newLimit {
				t.Errorf("BorrowingLimit = %v, want %v", user.BorrowingLimit, tt.newLimit)
			}
		})
	}
}

func TestUser_HasUnlimitedBorrowing(t *testing.T) {
	tests := []struct {
		name           string
		role           UserRole
		borrowingLimit int
		want           bool
	}{
		{
			name:           "admin with zero limit",
			role:           RoleAdmin,
			borrowingLimit: 0,
			want:           true,
		},
		{
			name:           "admin with non-zero limit",
			role:           RoleAdmin,
			borrowingLimit: 10,
			want:           false,
		},
		{
			name:           "librarian with zero limit",
			role:           RoleLibrarian,
			borrowingLimit: 0,
			want:           false,
		},
		{
			name:           "member with zero limit",
			role:           RoleMember,
			borrowingLimit: 0,
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Role:           tt.role,
				BorrowingLimit: tt.borrowingLimit,
			}
			if got := user.HasUnlimitedBorrowing(); got != tt.want {
				t.Errorf("HasUnlimitedBorrowing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_CanBorrowMoreBooks(t *testing.T) {
	tests := []struct {
		name              string
		role              UserRole
		status            UserStatus
		borrowingLimit    int
		currentBorrowCount int
		want              bool
	}{
		{
			name:              "member under limit",
			role:              RoleMember,
			status:            StatusActive,
			borrowingLimit:    5,
			currentBorrowCount: 3,
			want:              true,
		},
		{
			name:              "member at limit",
			role:              RoleMember,
			status:            StatusActive,
			borrowingLimit:    5,
			currentBorrowCount: 5,
			want:              false,
		},
		{
			name:              "member over limit",
			role:              RoleMember,
			status:            StatusActive,
			borrowingLimit:    5,
			currentBorrowCount: 6,
			want:              false,
		},
		{
			name:              "admin unlimited",
			role:              RoleAdmin,
			status:            StatusActive,
			borrowingLimit:    0,
			currentBorrowCount: 100,
			want:              true,
		},
		{
			name:              "suspended user",
			role:              RoleMember,
			status:            StatusSuspended,
			borrowingLimit:    5,
			currentBorrowCount: 0,
			want:              false,
		},
		{
			name:              "expired user",
			role:              RoleMember,
			status:            StatusExpired,
			borrowingLimit:    5,
			currentBorrowCount: 0,
			want:              false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Role:           tt.role,
				Status:         tt.status,
				BorrowingLimit: tt.borrowingLimit,
			}
			if got := user.CanBorrowMoreBooks(tt.currentBorrowCount); got != tt.want {
				t.Errorf("CanBorrowMoreBooks() = %v, want %v", got, tt.want)
			}
		})
	}
}