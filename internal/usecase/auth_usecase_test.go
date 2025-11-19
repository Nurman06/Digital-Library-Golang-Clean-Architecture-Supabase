package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/entity"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/usecase/mocks"
	"github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase/internal/usecase/testdata"
)

func TestAuthUseCase_Login(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		email         string
		password      string
		setupMocks    func(*mocks.MockUserRepository)
		wantErr       bool
		errorContains string
	}{
		{
			name:     "successful login",
			email:    "test@example.com",
			password: "password123",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByEmailFunc = func(ctx context.Context, email string) (*entity.User, error) {
					return testdata.CreateTestUser("user-1", entity.RoleMember), nil
				}
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
			},
			wantErr: false,
		},
		{
			name:          "empty email",
			email:         "",
			password:      "password123",
			setupMocks:    func(ur *mocks.MockUserRepository) {},
			wantErr:       true,
			errorContains: "email is required",
		},
		{
			name:          "empty password",
			email:         "test@example.com",
			password:      "",
			setupMocks:    func(ur *mocks.MockUserRepository) {},
			wantErr:       true,
			errorContains: "password is required",
		},
		{
			name:     "user not found",
			email:    "nonexistent@example.com",
			password: "password123",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByEmailFunc = func(ctx context.Context, email string) (*entity.User, error) {
					return nil, errors.New("user not found")
				}
			},
			wantErr:       true,
			errorContains: "authentication failed",
		},
		{
			name:     "suspended user",
			email:    "suspended@example.com",
			password: "password123",
			setupMocks: func(ur *mocks.MockUserRepository) {
				user := testdata.CreateTestUser("user-1", entity.RoleMember)
				user.Status = entity.StatusSuspended
				ur.GetByEmailFunc = func(ctx context.Context, email string) (*entity.User, error) {
					return user, nil
				}
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return user, nil
				}
			},
			wantErr:       true,
			errorContains: "authentication failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mocks.MockUserRepository{}
			tt.setupMocks(userRepo)

			uc := NewAuthUseCase(userRepo)
			user, err := uc.Login(ctx, tt.email, tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("Login() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr && user == nil {
				t.Error("Login() should return a user")
			}
		})
	}
}

func TestAuthUseCase_ValidateUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        string
		setupMocks    func(*mocks.MockUserRepository)
		wantErr       bool
		errorContains string
	}{
		{
			name:   "valid active user",
			userID: "user-123",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
			},
			wantErr: false,
		},
		{
			name:          "empty user ID",
			userID:        "",
			setupMocks:    func(ur *mocks.MockUserRepository) {},
			wantErr:       true,
			errorContains: "user ID is required",
		},
		{
			name:   "user not found",
			userID: "nonexistent",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return nil, errors.New("user not found")
				}
			},
			wantErr:       true,
			errorContains: "user validation failed",
		},
		{
			name:   "inactive user",
			userID: "user-123",
			setupMocks: func(ur *mocks.MockUserRepository) {
				user := testdata.CreateTestUser("user-123", entity.RoleMember)
				user.Status = entity.StatusSuspended
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return user, nil
				}
			},
			wantErr:       true,
			errorContains: "user account is not active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mocks.MockUserRepository{}
			tt.setupMocks(userRepo)

			uc := NewAuthUseCase(userRepo)
			err := uc.ValidateUser(ctx, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("ValidateUser() error = %v, should contain %v", err, tt.errorContains)
				}
			}
		})
	}
}

func TestAuthUseCase_CheckPermission(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        string
		operation     string
		setupMocks    func(*mocks.MockUserRepository)
		want          bool
		wantErr       bool
		errorContains string
	}{
		{
			name:      "admin can manage users",
			userID:    "admin-1",
			operation: "manage_users",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleAdmin), nil
				}
			},
			want:    true,
			wantErr: false,
		},
		{
			name:      "member cannot manage users",
			userID:    "member-1",
			operation: "manage_users",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
			},
			want:    false,
			wantErr: false,
		},
		{
			name:          "empty user ID",
			userID:        "",
			operation:     "manage_users",
			setupMocks:    func(ur *mocks.MockUserRepository) {},
			wantErr:       true,
			errorContains: "user ID is required",
		},
		{
			name:          "empty operation",
			userID:        "user-1",
			operation:     "",
			setupMocks:    func(ur *mocks.MockUserRepository) {},
			wantErr:       true,
			errorContains: "operation is required",
		},
		{
			name:      "inactive user cannot have permissions",
			userID:    "user-1",
			operation: "view_books",
			setupMocks: func(ur *mocks.MockUserRepository) {
				user := testdata.CreateTestUser("user-1", entity.RoleMember)
				user.Status = entity.StatusSuspended
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return user, nil
				}
			},
			wantErr:       true,
			errorContains: "user account is not active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mocks.MockUserRepository{}
			tt.setupMocks(userRepo)

			uc := NewAuthUseCase(userRepo)
			hasPermission, err := uc.CheckPermission(ctx, tt.userID, tt.operation)

			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPermission() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("CheckPermission() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr && hasPermission != tt.want {
				t.Errorf("CheckPermission() = %v, want %v", hasPermission, tt.want)
			}
		})
	}
}

func TestAuthUseCase_CanUserBorrow(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        string
		setupMocks    func(*mocks.MockUserRepository)
		want          bool
		wantErr       bool
		errorContains string
	}{
		{
			name:   "active user can borrow",
			userID: "user-123",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
			},
			want:    true,
			wantErr: false,
		},
		{
			name:   "suspended user cannot borrow",
			userID: "user-123",
			setupMocks: func(ur *mocks.MockUserRepository) {
				user := testdata.CreateTestUser("user-123", entity.RoleMember)
				user.Status = entity.StatusSuspended
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return user, nil
				}
			},
			want:    false,
			wantErr: false,
		},
		{
			name:          "empty user ID",
			userID:        "",
			setupMocks:    func(ur *mocks.MockUserRepository) {},
			wantErr:       true,
			errorContains: "user ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mocks.MockUserRepository{}
			tt.setupMocks(userRepo)

			uc := NewAuthUseCase(userRepo)
			canBorrow, err := uc.CanUserBorrow(ctx, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("CanUserBorrow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("CanUserBorrow() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr && canBorrow != tt.want {
				t.Errorf("CanUserBorrow() = %v, want %v", canBorrow, tt.want)
			}
		})
	}
}

func TestAuthUseCase_CanUserManageBooks(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		userID     string
		setupMocks func(*mocks.MockUserRepository)
		want       bool
		wantErr    bool
	}{
		{
			name:   "admin can manage books",
			userID: "admin-1",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleAdmin), nil
				}
			},
			want:    true,
			wantErr: false,
		},
		{
			name:   "librarian can manage books",
			userID: "librarian-1",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleLibrarian), nil
				}
			},
			want:    true,
			wantErr: false,
		},
		{
			name:   "member cannot manage books",
			userID: "member-1",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
			},
			want:    false,
			wantErr: false,
		},
		{
			name:   "inactive admin cannot manage books",
			userID: "admin-1",
			setupMocks: func(ur *mocks.MockUserRepository) {
				user := testdata.CreateTestUser("admin-1", entity.RoleAdmin)
				user.Status = entity.StatusSuspended
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return user, nil
				}
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mocks.MockUserRepository{}
			tt.setupMocks(userRepo)

			uc := NewAuthUseCase(userRepo)
			canManage, err := uc.CanUserManageBooks(ctx, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("CanUserManageBooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && canManage != tt.want {
				t.Errorf("CanUserManageBooks() = %v, want %v", canManage, tt.want)
			}
		})
	}
}

func TestAuthUseCase_CanUserManageUsers(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		userID     string
		setupMocks func(*mocks.MockUserRepository)
		want       bool
		wantErr    bool
	}{
		{
			name:   "admin can manage users",
			userID: "admin-1",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleAdmin), nil
				}
			},
			want:    true,
			wantErr: false,
		},
		{
			name:   "librarian cannot manage users",
			userID: "librarian-1",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleLibrarian), nil
				}
			},
			want:    false,
			wantErr: false,
		},
		{
			name:   "member cannot manage users",
			userID: "member-1",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mocks.MockUserRepository{}
			tt.setupMocks(userRepo)

			uc := NewAuthUseCase(userRepo)
			canManage, err := uc.CanUserManageUsers(ctx, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("CanUserManageUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && canManage != tt.want {
				t.Errorf("CanUserManageUsers() = %v, want %v", canManage, tt.want)
			}
		})
	}
}

func TestAuthUseCase_GetUserRole(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        string
		setupMocks    func(*mocks.MockUserRepository)
		wantRole      entity.UserRole
		wantErr       bool
		errorContains string
	}{
		{
			name:   "get admin role",
			userID: "admin-1",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleAdmin), nil
				}
			},
			wantRole: entity.RoleAdmin,
			wantErr:  false,
		},
		{
			name:   "get member role",
			userID: "member-1",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
			},
			wantRole: entity.RoleMember,
			wantErr:  false,
		},
		{
			name:          "empty user ID",
			userID:        "",
			setupMocks:    func(ur *mocks.MockUserRepository) {},
			wantErr:       true,
			errorContains: "user ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mocks.MockUserRepository{}
			tt.setupMocks(userRepo)

			uc := NewAuthUseCase(userRepo)
			role, err := uc.GetUserRole(ctx, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("GetUserRole() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr && role != tt.wantRole {
				t.Errorf("GetUserRole() = %v, want %v", role, tt.wantRole)
			}
		})
	}
}

func TestAuthUseCase_ValidateUserStatus(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        string
		setupMocks    func(*mocks.MockUserRepository)
		wantErr       bool
		errorContains string
	}{
		{
			name:   "active user is valid",
			userID: "user-123",
			setupMocks: func(ur *mocks.MockUserRepository) {
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return testdata.CreateTestUser(id, entity.RoleMember), nil
				}
			},
			wantErr: false,
		},
		{
			name:   "suspended user is invalid",
			userID: "user-123",
			setupMocks: func(ur *mocks.MockUserRepository) {
				user := testdata.CreateTestUser("user-123", entity.RoleMember)
				user.Status = entity.StatusSuspended
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return user, nil
				}
			},
			wantErr:       true,
			errorContains: "user account is suspended",
		},
		{
			name:   "expired user is invalid",
			userID: "user-123",
			setupMocks: func(ur *mocks.MockUserRepository) {
				user := testdata.CreateTestUser("user-123", entity.RoleMember)
				user.Status = entity.StatusExpired
				ur.GetByIDFunc = func(ctx context.Context, id string) (*entity.User, error) {
					return user, nil
				}
			},
			wantErr:       true,
			errorContains: "user account has expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mocks.MockUserRepository{}
			tt.setupMocks(userRepo)

			uc := NewAuthUseCase(userRepo)
			err := uc.ValidateUserStatus(ctx, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUserStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil || !contains(err.Error(), tt.errorContains) {
					t.Errorf("ValidateUserStatus() error = %v, should contain %v", err, tt.errorContains)
				}
			}
		})
	}
}