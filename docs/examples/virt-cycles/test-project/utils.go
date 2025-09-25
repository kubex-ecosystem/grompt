// Package main provides utility functions - needs better organization
package main

import "strings"

// stringUtils - should be separate package or better organized
func stringUtils(s string) string {
	// No validation
	return strings.ToUpper(s)
}

// validateEmail - poor implementation
func validateEmail(email string) bool {
	// Very basic validation - should use proper regex
	return strings.Contains(email, "@")
}
