// Package providers defines interfaces for AI providers.
package providers

import "github.com/rafa-mori/grompt/internal/types"

// Provider represents an AI provider interface
type Provider interface {
	// Name returns the provider name (e.g., "openai", "claude", "ollama")
	Name() string

	// Execute sends a prompt to the provider and returns the response
	Execute(prompt string) (string, error)

	// IsAvailable checks if the provider is configured and ready
	IsAvailable() bool

	// GetCapabilities returns provider-specific capabilities
	GetCapabilities() *types.Capabilities
}

func NewProvider(name, apiKey string, cfg types.IConfig) Provider {
	return &types.ProviderImpl{
		VName:   name,
		VAPI:    cfg.GetAPIConfig(name),
		VConfig: cfg,
	}
}

// Initialize creates and returns all available providers
func Initialize(claudeKey, openaiKey, deepseekKey, ollamaEndpoint string) []Provider {
	// This function signature is kept for backwards compatibility
	// The actual implementation is done in internal/engine which calls internal/providers
	// For now, return empty slice - the real initialization happens in the engine
	return []Provider{}
}
