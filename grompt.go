// Package grompt provides a simple and flexible way to create interactive command-line prompts in Go.
package grompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Prompt represents a command-line prompt with a message and an optional default value.
type Prompt struct {
	Message      string
	DefaultValue string
}

// NewPrompt creates a new Prompt instance with the given message and default value.
func NewPrompt(message, defaultValue string) *Prompt {
	return &Prompt{
		Message:      message,
		DefaultValue: defaultValue,
	}
}

// Show displays the prompt to the user and captures their input.
// If the user provides no input, the default value is returned.
func (p *Prompt) Show() string {
	reader := bufio.NewReader(os.Stdin)
	if p.DefaultValue != "" {
		fmt.Printf("%s [%s]: ", p.Message, p.DefaultValue)
	} else {
		fmt.Printf("%s: ", p.Message)
	}
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return p.DefaultValue
	}
	return input
}

// Example usage:
// func main() {
//     prompt := NewPrompt("Enter your name", "Guest")
//     name := prompt.Show()
//     fmt.Printf("Hello, %s!\n", name)
// }
