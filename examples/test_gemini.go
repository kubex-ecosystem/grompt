package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rafa-mori/grompt"
)

func main() {
	// Test Gemini provider with user's API key
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required")
	}

	// Create config with Gemini API key
	config := grompt.DefaultConfig()
	config.SetAPIKey("gemini", apiKey)

	// Create engine
	engine := grompt.NewPromptEngine(config)

	// Get available providers
	providers := engine.GetProviders()
	fmt.Printf("Available providers: %d\n", len(providers))

	for _, provider := range providers {
		fmt.Printf("- %s\n", provider.Name())
	}

	// Test prompt processing
	result, err := engine.ProcessPrompt("Hello, what is the capital of Brazil?", nil)
	if err != nil {
		log.Fatalf("Failed to process prompt: %v", err)
	}

	fmt.Printf("\nPrompt: %s\n", result.Prompt)
	fmt.Printf("Response: %s\n", result.Response)
	fmt.Printf("Provider: %s\n", result.Provider)
}
