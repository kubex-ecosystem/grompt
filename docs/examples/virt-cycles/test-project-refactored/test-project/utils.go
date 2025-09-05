package main

import (
	"strings"
	"regexp"
)

// StringUtils provides utility functions for string manipulation.
// Consider moving this to a separate package for better organization if it grows.

// StringToUpper converts a string to uppercase.
func StringToUpper(s string) string {
	// Input validation is not strictly necessary here, as ToUpper handles empty strings fine.
	return strings.ToUpper(s)
}

// IsValidEmail validates an email address using a regular expression.
func IsValidEmail(email string) bool {
	// More robust email validation using regular expression.
	// This regex is a simplified version and might not cover all valid email formats.
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
```