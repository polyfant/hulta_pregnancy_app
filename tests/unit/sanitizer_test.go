package unit

import (
	"testing"

	"github.com/polyfant/hulta_pregnancy_app/internal/validation"
)

func TestSanitizer(t *testing.T) {
	sanitizer := validation.NewSanitizer()

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Remove SQL Injection",
			input:    "SELECT * FROM users; DROP TABLE users;",
			expected: "",
		},
		{
			name:     "Trim Whitespace",
			input:    "  test  ",
			expected: "test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := sanitizer.Sanitize(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestEmailValidation(t *testing.T) {
	sanitizer := validation.NewSanitizer()

	testCases := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "Valid Email",
			email:    "test@example.com",
			expected: true,
		},
		{
			name:     "Invalid Email",
			email:    "invalid-email",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := sanitizer.ValidateEmail(tc.email)
			hasError := err != nil
			
			if hasError == tc.expected {
				t.Errorf("Expected valid=%v, got error=%v", tc.expected, err)
			}
		})
	}
}
