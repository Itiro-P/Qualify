package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
	Code    string
}

// PasswordStrength validates password strength and returns errors
func ValidatePasswordStrength(password string) []ValidationError {
	var errors []ValidationError

	if len(password) < 8 {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password must be at least 8 characters long",
			Code:    "PASSWORD_TOO_SHORT",
		})
	}

	if !containsUpperCase(password) {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password must contain at least one uppercase letter",
			Code:    "PASSWORD_MISSING_UPPERCASE",
		})
	}

	if !containsLowerCase(password) {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password must contain at least one lowercase letter",
			Code:    "PASSWORD_MISSING_LOWERCASE",
		})
	}

	if !containsDigit(password) {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password must contain at least one digit",
			Code:    "PASSWORD_MISSING_DIGIT",
		})
	}

	if !containsSpecialChar(password) {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password must contain at least one special character (!@#$%^&*)",
			Code:    "PASSWORD_MISSING_SPECIAL_CHAR",
		})
	}

	return errors
}

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	if !re.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	if len(email) > 255 {
		return fmt.Errorf("email is too long (max 255 characters)")
	}

	return nil
}

// SanitizeEmail removes whitespace and converts to lowercase
func SanitizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

// SanitizeInput removes leading/trailing whitespace
func SanitizeInput(input string) string {
	return strings.TrimSpace(input)
}

// ValidateName validates user name
func ValidateName(name string) error {
	name = SanitizeInput(name)

	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters long")
	}

	if len(name) > 255 {
		return fmt.Errorf("name is too long (max 255 characters)")
	}

	return nil
}

// ValidateRefreshToken validates refresh token format
func ValidateRefreshToken(token string) error {
	if len(token) == 0 {
		return fmt.Errorf("refresh token is required")
	}

	if len(token) > 500 {
		return fmt.Errorf("invalid refresh token format")
	}

	return nil
}

// Helper functions
func containsUpperCase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func containsLowerCase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

func containsDigit(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func containsSpecialChar(s string) bool {
	specialChars := "!@#$%^&*()-_=+[]{}|;:,.<>?"
	for _, r := range specialChars {
		if strings.ContainsRune(s, r) {
			return true
		}
	}
	return false
}

// ValidatePasswordMatch ensures both passwords match
func ValidatePasswordMatch(password1, password2 string) error {
	if password1 != password2 {
		return fmt.Errorf("passwords do not match")
	}
	return nil
}

// ValidateBuildValidationErrorMap builds a map of validation errors
func BuildValidationErrorMap(errors []ValidationError) map[string]string {
	errorMap := make(map[string]string)
	for _, err := range errors {
		errorMap[err.Field] = err.Message
	}
	return errorMap
}
