// Package providers defines interfaces for AI providers.
package providers

import "github.com/rafa-mori/grompt/internal/types"

// Provider represents an AI provider interface
type Provider interface {
	// Name returns the provider name (e.g., "openai", "claude", "ollama")
	Name() string

	// Version returns the provider version
	Version() string

	// Execute sends a prompt to the provider and returns the response
	Execute(prompt string) (string, error)

	// IsAvailable checks if the provider is configured and ready
	IsAvailable() bool

	// GetCapabilities returns provider-specific capabilities
	GetCapabilities() *types.Capabilities
}

type Capabilities = types.Capabilities

func NewProvider(name, apiKey, version string, cfg types.IConfig) Provider {
	return &types.ProviderImpl{
		VName:    name,
		VVersion: version,
		VAPI:     cfg.GetAPIConfig(name),
		VConfig:  cfg,
	}
}

// Initialize creates and returns all available providers
func Initialize(claudeKey, openaiKey, deepseekKey, ollamaEndpoint string) []Provider {
	if claudeKey == "" && openaiKey == "" && deepseekKey == "" && ollamaEndpoint == "" {
		return []Provider{}
	}
	var cfg types.IConfig = types.NewConfig("8080", openaiKey, deepseekKey, ollamaEndpoint, claudeKey, "")
	var providers []Provider
	if claudeKey != "" {
		providers = append(providers, NewProvider("claude", claudeKey, "v1", cfg))
	}
	if openaiKey != "" {
		providers = append(providers, NewProvider("openai", openaiKey, "v1", cfg))
	}
	if deepseekKey != "" {
		providers = append(providers, NewProvider("deepseek", deepseekKey, "v1", cfg))
	}
	if ollamaEndpoint != "" {
		providers = append(providers, NewProvider("ollama", ollamaEndpoint, "v1", cfg))
	}
	return providers
}
