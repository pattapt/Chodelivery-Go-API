package config

import (
	"regexp"
)

func IsValidEmail(email string) bool {
	// Use a regular expression to validate email format
	// This is a basic example and may not cover all edge cases
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`
	matched, err := regexp.MatchString(pattern, email)
	return matched && err == nil
}
