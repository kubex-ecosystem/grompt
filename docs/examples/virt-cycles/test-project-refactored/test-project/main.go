package main

import (
	"errors"
	"fmt"
	"os"
)

// User represents a user in the system
type User struct {
	ID    int    // Changed to ID to follow Go naming conventions (exported field)
	Name  string // Changed to Name to follow Go naming conventions (exported field)
	Email string // Changed to Email to follow Go naming conventions (exported field)
}

// ErrInvalidID is a custom error for invalid user IDs.
var ErrInvalidID = errors.New("invalid user ID")

// GetUserInfo returns user information.  Returns an error if the ID is invalid.
func GetUserInfo(id int) (User, error) {
	// Improved error handling: return an error instead of printing and exiting.
	if id < 0 {
		return User{}, ErrInvalidID // Return a zero-value User and the custom error
	}

	// Hardcoded data - consider using a proper data source (e.g., database, file).
	// This is just for demonstration purposes.
	user := User{
		ID:    id,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	return user, nil // Return nil error on success
}

// PrintUser prints user information. Includes validation and better formatting.
func PrintUser(user User) {
	// Validate input to prevent unexpected behavior.
	if user.ID <= 0 {
		fmt.Println("Invalid user ID")
		return
	}
	if user.Name == "" {
		fmt.Println("Invalid user name")
		return
	}
	if user.Email == "" {
		fmt.Println("Invalid user email")
		return
	}

	// Improved formatting for better readability.
	fmt.Printf("User ID: %d
", user.ID)
	fmt.Printf("Name: %s
", user.Name)
	fmt.Printf("Email: %s
", user.Email)
}

// isValidUser checks if a user is valid based on ID, Name, and Email.
func isValidUser(user User) bool {
	// Consolidated validation logic for better readability and maintainability.
	return user.ID > 0 && user.Name != "" && user.Email != ""
}

func main() {
	// Improved error handling: check for errors and handle them appropriately.
	user, err := GetUserInfo(1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting user info: %v
", err) // Print error to stderr
		os.Exit(1)                                                     // Exit with an error code
		return
	}

	PrintUser(user)

	// Simplified logic using the isValidUser function.
	if isValidUser(user) {
		fmt.Println("User is valid")
	} else {
		fmt.Println("User is not valid") // More generic message as the validation is now consolidated
	}
}