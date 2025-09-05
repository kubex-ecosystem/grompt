// Package providers defines interfaces for AI providers.
package providers

import (
	"github.com/rafa-mori/grompt/internal/core/provider"
	"github.com/rafa-mori/grompt/internal/types"
	"github.com/rafa-mori/logz"
)

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

func NewProvider(name, apiKey, version string, cfg provider.IConfig) Provider {
	return &types.ProviderImpl{
		VName:    name,
		VVersion: version,
		VAPI:     cfg.GetAPIConfig(name),
		VConfig:  cfg,
	}
}

// Initialize creates and returns all available providers
func Initialize(
	bindAddr,
	port,
	openAIKey,
	deepSeekKey,
	ollamaEndpoint,
	claudeKey,
	geminiKey,
	chatGPTKey string,
	logger logz.Logger,
) []Provider {

	if bindAddr == "" &&
		port == "" &&
		openAIKey == "" &&
		deepSeekKey == "" &&
		ollamaEndpoint == "" &&
		claudeKey == "" &&
		geminiKey == "" &&
		chatGPTKey == "" {
		return []Provider{}
	}

	var cfg = provider.NewConfig(
		bindAddr,
		"8080",
		openAIKey,
		deepSeekKey,
		ollamaEndpoint,
		claudeKey,
		geminiKey,
		chatGPTKey,
		nil,
	)

	cfg.Logger = logger

	var providers []Provider
	if claudeKey != "" {
		providers = append(providers, NewProvider("claude", claudeKey, "v1", cfg))
	}
	if openAIKey != "" {
		providers = append(providers, NewProvider("openai", openAIKey, "v1", cfg))
	}
	if deepSeekKey != "" {
		providers = append(providers, NewProvider("deepseek", deepSeekKey, "v1", cfg))
	}
	if ollamaEndpoint != "" {
		providers = append(providers, NewProvider("ollama", ollamaEndpoint, "v1", cfg))
	}
	if geminiKey != "" {
		providers = append(providers, NewProvider("gemini", geminiKey, "v1", cfg))
	}
	if chatGPTKey != "" {
		providers = append(providers, NewProvider("chatgpt", chatGPTKey, "v1", cfg))
	}

	return providers
}
