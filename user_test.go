package gobase

import (
	"testing"
)

// TestUser_SetPassword tests password hashing
func TestUser_SetPassword(t *testing.T) {
	user := &User{}

	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{
			name:        "Valid password",
			password:    "StrongPass123",
			expectError: false,
		},
		{
			name:        "Short password",
			password:    "123",
			expectError: true,
		},
		{
			name:        "Empty password",
			password:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.SetPassword(tt.password)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else if !tt.expectError {
				// Verify password was hashed
				if user.PasswordHash == "" {
					t.Error("Expected password hash to be set")
				}
				if user.PasswordHash == tt.password {
					t.Error("Password should be hashed, not stored in plain text")
				}
			}
		})
	}
}

// TestUser_CheckPassword tests password verification
func TestUser_CheckPassword(t *testing.T) {
	user := &User{}
	password := "TestPassword123"

	// Set password
	err := user.SetPassword(password)
	if err != nil {
		t.Fatalf("Failed to set password: %v", err)
	}

	// Test correct password
	if !user.CheckPassword(password) {
		t.Error("CheckPassword should return true for correct password")
	}

	// Test incorrect password
	if user.CheckPassword("WrongPassword") {
		t.Error("CheckPassword should return false for incorrect password")
	}
}

// TestUser_IsValidEmail tests email validation
func TestUser_IsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "Valid email",
			email:    "test@example.com",
			expected: true,
		},
		{
			name:     "Valid email with subdomain",
			email:    "user@mail.example.com",
			expected: true,
		},
		{
			name:     "Invalid email without @",
			email:    "testexample.com",
			expected: false,
		},
		{
			name:     "Invalid email without domain",
			email:    "test@",
			expected: false,
		},
		{
			name:     "Empty email",
			email:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Email: tt.email}
			result := user.IsValidEmail()

			if result != tt.expected {
				t.Errorf("Expected %v, got %v for email: %s", tt.expected, result, tt.email)
			}
		})
	}
}

// TestUser_IsValidUsername tests username validation
func TestUser_IsValidUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		expected bool
	}{
		{
			name:     "Valid username",
			username: "testuser",
			expected: true,
		},
		{
			name:     "Valid username with numbers",
			username: "user123",
			expected: true,
		},
		{
			name:     "Valid username with underscore",
			username: "test_user",
			expected: true,
		},
		{
			name:     "Too short username",
			username: "ab",
			expected: false,
		},
		{
			name:     "Username with special characters",
			username: "test-user",
			expected: false,
		},
		{
			name:     "Empty username",
			username: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Username: tt.username}
			result := user.IsValidUsername()

			if result != tt.expected {
				t.Errorf("Expected %v, got %v for username: %s", tt.expected, result, tt.username)
			}
		})
	}
}

// TestUser_HasRole tests role checking
func TestUser_HasRole(t *testing.T) {
	user := &User{Role: RoleAdmin}

	if !user.HasRole(RoleAdmin) {
		t.Error("User should have admin role")
	}

	if user.HasRole(RoleUser) {
		t.Error("User should not have user role")
	}
}

// TestUser_IsAdminUser tests admin user checking
func TestUser_IsAdminUser(t *testing.T) {
	tests := []struct {
		name        string
		user        *User
		expected    bool
		description string
	}{
		{
			name:        "Superuser role",
			user:        &User{Role: RoleSuperuser},
			expected:    true,
			description: "User with superuser role should be admin",
		},
		{
			name:        "Admin role",
			user:        &User{Role: RoleAdmin},
			expected:    true,
			description: "User with admin role should be admin",
		},
		{
			name:        "IsSuperuser flag",
			user:        &User{Role: RoleUser, IsSuperuser: true},
			expected:    true,
			description: "User with IsSuperuser flag should be admin",
		},
		{
			name:        "Regular user",
			user:        &User{Role: RoleUser},
			expected:    false,
			description: "Regular user should not be admin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.IsAdminUser()
			if result != tt.expected {
				t.Errorf("%s: Expected %v, got %v", tt.description, tt.expected, result)
			}
		})
	}
}

// TestCreateSuperuser tests superuser creation functionality
func TestCreateSuperuser(t *testing.T) {
	connection := setupTestDB(t)
	accessor := NewAccessor(connection)

	// Register User model
	RegisterModel(&User{})

	err := accessor.Migrate(&User{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	tests := []struct {
		name        string
		username    string
		email       string
		password    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid superuser",
			username:    "admin",
			email:       "admin@example.com",
			password:    "SuperSecret123",
			expectError: false,
		},
		{
			name:        "Invalid username",
			username:    "a",
			email:       "admin@example.com",
			password:    "SuperSecret123",
			expectError: true,
			errorMsg:    "username must be between 3 and 150 characters",
		},
		{
			name:        "Invalid email",
			username:    "admin2",
			email:       "invalid-email",
			password:    "SuperSecret123",
			expectError: true,
			errorMsg:    "invalid email format",
		},
		{
			name:        "Weak password",
			username:    "admin3",
			email:       "admin3@example.com",
			password:    "weak",
			expectError: true,
			errorMsg:    "password must be at least 8 characters long",
		},
		{
			name:        "Password without uppercase",
			username:    "admin4",
			email:       "admin4@example.com",
			password:    "nouppercase123",
			expectError: true,
			errorMsg:    "password must contain at least one uppercase letter",
		},
		{
			name:        "Password without lowercase",
			username:    "admin5",
			email:       "admin5@example.com",
			password:    "NOLOWERCASE123",
			expectError: true,
			errorMsg:    "password must contain at least one lowercase letter",
		},
		{
			name:        "Password without digit",
			username:    "admin6",
			email:       "admin6@example.com",
			password:    "NoDigitPassword",
			expectError: true,
			errorMsg:    "password must contain at least one digit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateSuperuser(accessor, tt.username, tt.email, tt.password)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != "validation failed: "+tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else {
					// Verify superuser was created correctly
					users := []*User{}
					conditions := map[string]interface{}{
						"username": tt.username,
					}
					err = accessor.Filter(&users, conditions)
					if err != nil {
						t.Errorf("Failed to retrieve created superuser: %v", err)
					} else if len(users) == 0 {
						t.Error("No superuser found with the given username")
					} else {
						user := users[0]
						if user.Role != RoleSuperuser {
							t.Errorf("Expected role %s, got %s", RoleSuperuser, user.Role)
						}
						if !user.IsSuperuser {
							t.Error("Expected IsSuperuser to be true")
						}
						if !user.IsStaff {
							t.Error("Expected IsStaff to be true")
						}
						if !user.IsActive {
							t.Error("Expected IsActive to be true")
						}
						if user.Email != tt.email {
							t.Errorf("Expected email %s, got %s", tt.email, user.Email)
						}
					}
				}
			}
		})
	}
}

// TestCreateSuperuser_DuplicateUser tests creating superuser with existing username/email
func TestCreateSuperuser_DuplicateUser(t *testing.T) {
	connection := setupTestDB(t)
	accessor := NewAccessor(connection)

	// Register User model
	RegisterModel(&User{})

	err := accessor.Migrate(&User{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Create first superuser
	err = CreateSuperuser(accessor, "admin", "admin@example.com", "SuperSecret123")
	if err != nil {
		t.Fatalf("Failed to create first superuser: %v", err)
	}

	// Try to create another with same username
	err = CreateSuperuser(accessor, "admin", "different@example.com", "SuperSecret123")
	if err == nil {
		t.Error("Expected error when creating superuser with duplicate username")
	}

	// Try to create another with same email
	err = CreateSuperuser(accessor, "different", "admin@example.com", "SuperSecret123")
	if err == nil {
		t.Error("Expected error when creating superuser with duplicate email")
	}
}

// TestUserModelValidation tests User model validation during creation
func TestUserModelValidation(t *testing.T) {
	connection := setupTestDB(t)
	accessor := NewAccessor(connection)

	// Register User model
	RegisterModel(&User{})

	err := accessor.Migrate(&User{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	tests := []struct {
		name        string
		user        *User
		expectError bool
		setupFunc   func(*User)
	}{
		{
			name: "Valid user",
			user: &User{
				Username: "testuser",
				Email:    "test@example.com",
			},
			expectError: false,
			setupFunc: func(u *User) {
				u.SetPassword("TestPassword123")
			},
		},
		{
			name: "Missing username",
			user: &User{
				Email: "test@example.com",
			},
			expectError: true,
			setupFunc: func(u *User) {
				u.SetPassword("TestPassword123")
			},
		},
		{
			name: "Missing email",
			user: &User{
				Username: "testuser",
			},
			expectError: true,
			setupFunc: func(u *User) {
				u.SetPassword("TestPassword123")
			},
		},
		{
			name: "Missing password",
			user: &User{
				Username: "testuser",
				Email:    "test@example.com",
			},
			expectError: true,
			setupFunc: func(u *User) {
				// Don't set password
			},
		},
		{
			name: "Invalid email format",
			user: &User{
				Username: "testuser",
				Email:    "invalid-email",
			},
			expectError: true,
			setupFunc: func(u *User) {
				u.SetPassword("TestPassword123")
			},
		},
		{
			name: "Invalid username format",
			user: &User{
				Username: "test-user", // Contains hyphen
				Email:    "test@example.com",
			},
			expectError: true,
			setupFunc: func(u *User) {
				u.SetPassword("TestPassword123")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup user
			if tt.setupFunc != nil {
				tt.setupFunc(tt.user)
			}

			err := accessor.Create(tt.user)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else {
					// Verify default role was set
					if tt.user.Role == "" {
						t.Error("Expected default role to be set")
					}
				}
			}
		})
	}
}
