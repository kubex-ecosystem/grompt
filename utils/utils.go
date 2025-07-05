// Package utils provides utility functions for the application.
// It provides a function to get environment variables with a default value.
// It is used to retrieve configuration values from the environment, allowing for flexible application settings.
package utils

import "os"

func GetEnvOr(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
