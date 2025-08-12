package gobase

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Role constants for user roles
const (
	RoleUser      = "user"
	RoleAdmin     = "admin"
	RoleSuperuser = "superuser"
)

// User represents the default user model with authentication and authorization capabilities
type User struct {
	BaseModel
	Username     string     `gorm:"uniqueIndex;size:150;not null" json:"username"`
	Email        string     `gorm:"uniqueIndex;size:254;not null" json:"email"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"` // Don't expose in JSON
	FirstName    string     `gorm:"size:100" json:"first_name"`
	LastName     string     `gorm:"size:100" json:"last_name"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	IsStaff      bool       `gorm:"default:false" json:"is_staff"`
	IsSuperuser  bool       `gorm:"default:false" json:"is_superuser"`
	Role         string     `gorm:"size:50;default:'user'" json:"role"`
	LastLogin    *time.Time `json:"last_login"`
}

// SetPassword hashes and sets the user's password
func (u *User) SetPassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword verifies if the provided password matches the user's password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// IsValidEmail validates the email format
func (u *User) IsValidEmail() bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(u.Email)
}

// IsValidUsername validates the username format
func (u *User) IsValidUsername() bool {
	if len(u.Username) < 3 || len(u.Username) > 150 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(u.Username)
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Username
	}
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// HasRole checks if the user has a specific role
func (u *User) HasRole(role string) bool {
	return u.Role == role
}

// IsAdminUser checks if the user is an admin or superuser
func (u *User) IsAdminUser() bool {
	return u.Role == RoleAdmin || u.Role == RoleSuperuser || u.IsSuperuser
}

// BeforeCreate is a GORM hook that validates the user before creation
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Validate required fields
	if u.Username == "" {
		return errors.New("username is required")
	}

	if u.Email == "" {
		return errors.New("email is required")
	}

	if u.PasswordHash == "" {
		return errors.New("password is required")
	}

	// Validate email format
	if !u.IsValidEmail() {
		return errors.New("invalid email format")
	}

	// Validate username format
	if !u.IsValidUsername() {
		return errors.New("invalid username format (3-150 chars, alphanumeric and underscore only)")
	}

	// Set default role if not provided
	if u.Role == "" {
		u.Role = RoleUser
	}

	// Validate role
	validRoles := []string{RoleUser, RoleAdmin, RoleSuperuser}
	isValidRole := false
	for _, role := range validRoles {
		if u.Role == role {
			isValidRole = true
			break
		}
	}
	if !isValidRole {
		return fmt.Errorf("invalid role: %s. Valid roles: %v", u.Role, validRoles)
	}

	// Set superuser flags based on role
	if u.Role == RoleSuperuser {
		u.IsSuperuser = true
		u.IsStaff = true
	} else if u.Role == RoleAdmin {
		u.IsStaff = true
	}

	return nil
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
