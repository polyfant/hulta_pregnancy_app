package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/microcosm-cc/bluemonday"
)

type Sanitizer struct {
	policy *bluemonday.Policy
}

func NewSanitizer() *Sanitizer {
	policy := bluemonday.StrictPolicy()
	policy.AllowStandardURLs()
	return &Sanitizer{policy: policy}
}

func (s *Sanitizer) Sanitize(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	// Remove non-printable characters
	input = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, input)

	// HTML sanitization
	input = s.policy.Sanitize(input)

	// Remove potential injection patterns
	sqlInjectionRegex := regexp.MustCompile(`(?i)(--|\b(SELECT|INSERT|UPDATE|DELETE|DROP|UNION|ALTER)\b)`)
	if sqlInjectionRegex.MatchString(input) {
		return ""
	}

	return input
}

func (s *Sanitizer) ValidateEmail(email string) (string, error) {
	cleanEmail := s.Sanitize(email)
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	
	if !emailRegex.MatchString(strings.ToLower(cleanEmail)) {
		return "", fmt.Errorf("invalid email format")
	}
	
	return cleanEmail, nil
}

func (s *Sanitizer) ValidateID(id string) (uint, error) {
	cleanID := s.Sanitize(id)
	idRegex := regexp.MustCompile(`^\d+$`)
	
	if !idRegex.MatchString(cleanID) {
		return 0, fmt.Errorf("invalid ID format")
	}

	var parsedID uint
	_, err := fmt.Sscanf(cleanID, "%d", &parsedID)
	if err != nil {
		return 0, fmt.Errorf("failed to parse ID")
	}
	
	return parsedID, nil
}

func (s *Sanitizer) ValidatePassword(password string) error {
	// Minimum 8 characters, at least one uppercase, one lowercase, one number
	passwordRegex := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$`)
	
	if !passwordRegex.MatchString(password) {
		return fmt.Errorf("password does not meet complexity requirements")
	}
	
	return nil
}
