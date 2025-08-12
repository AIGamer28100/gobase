package gobase

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// CreateSuperuser creates a new superuser with the specified credentials
func CreateSuperuser(accessor *Accessor, username, email, password string) error {
	// Validate inputs
	if err := validateSuperuserInputs(username, email, password); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if user already exists
	existingUser := &User{}
	err := accessor.Get(existingUser, username)
	if err == nil {
		return fmt.Errorf("user with username '%s' already exists", username)
	}

	// Check if email already exists
	var users []User
	err = accessor.Filter(&users, map[string]interface{}{"email": email})
	if err != nil {
		return fmt.Errorf("failed to check email uniqueness: %w", err)
	}
	if len(users) > 0 {
		return fmt.Errorf("user with email '%s' already exists", email)
	}

	// Create new superuser
	superuser := &User{
		Username:    username,
		Email:       email,
		Role:        RoleSuperuser,
		IsActive:    true,
		IsStaff:     true,
		IsSuperuser: true,
	}

	// Set password
	if err := superuser.SetPassword(password); err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	// Save to database
	if err := accessor.Create(superuser); err != nil {
		return fmt.Errorf("failed to create superuser: %w", err)
	}

	return nil
}

// validateSuperuserInputs validates the inputs for superuser creation
func validateSuperuserInputs(username, email, password string) error {
	// Validate username
	if strings.TrimSpace(username) == "" {
		return errors.New("username cannot be empty")
	}
	if len(username) < 3 || len(username) > 150 {
		return errors.New("username must be between 3 and 150 characters")
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(username) {
		return errors.New("username can only contain letters, numbers, and underscores")
	}

	// Validate email
	if strings.TrimSpace(email) == "" {
		return errors.New("email cannot be empty")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	// Validate password
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if !hasUpperCase(password) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLowerCase(password) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasDigit(password) {
		return errors.New("password must contain at least one digit")
	}

	return nil
}

// hasUpperCase checks if string contains uppercase letters
func hasUpperCase(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

// hasLowerCase checks if string contains lowercase letters
func hasLowerCase(s string) bool {
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			return true
		}
	}
	return false
}

// hasDigit checks if string contains digits
func hasDigit(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}
