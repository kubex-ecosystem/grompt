// Package providers defines interfaces for AI providers.
package providers

// Provider represents an AI provider interface
type Provider interface {
	// Name returns the provider name (e.g., "openai", "claude", "ollama")
	Name() string

	// Execute sends a prompt to the provider and returns the response
	Execute(prompt string) (string, error)

	// IsAvailable checks if the provider is configured and ready
	IsAvailable() bool

	// GetCapabilities returns provider-specific capabilities
	GetCapabilities() Capabilities
}

// Capabilities describes what a provider can do
type Capabilities struct {
	MaxTokens         int      `json:"max_tokens"`
	SupportsBatch     bool     `json:"supports_batch"`
	SupportsStreaming bool     `json:"supports_streaming"`
	Models            []string `json:"models"`
	Pricing           *Pricing `json:"pricing,omitempty"`
}

// Pricing information for the provider
type Pricing struct {
	InputCostPer1K  float64 `json:"input_cost_per_1k"`
	OutputCostPer1K float64 `json:"output_cost_per_1k"`
	Currency        string  `json:"currency"`
}

// Initialize creates and returns all available providers
func Initialize(claudeKey, openaiKey, deepseekKey, ollamaEndpoint string) []Provider {
	var providers []Provider

	// This will be implemented in internal/providers with concrete implementations
	// For now, return empty slice to avoid compilation errors
	return providers
}
