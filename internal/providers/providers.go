// Package providers provides concrete implementations of AI providers.
package providers

import (
	"github.com/rafa-mori/grompt/factory/providers"
	"github.com/rafa-mori/grompt/internal/types"
)

// concreteProvider wraps the types.IAPIConfig to implement providers.Provider
type concreteProvider struct {
	name   string
	api    types.IAPIConfig
	config types.IConfig
}

// Name returns the provider name
func (cp *concreteProvider) Name() string {
	return cp.name
}

// Execute sends a prompt to the provider and returns the response
func (cp *concreteProvider) Execute(prompt string) (string, error) {
	return cp.api.Complete(prompt, 2048, "") // Default max tokens
}

// IsAvailable checks if the provider is configured and ready
func (cp *concreteProvider) IsAvailable() bool {
	return cp.api.IsAvailable()
}

// GetCapabilities returns provider-specific capabilities
func (cp *concreteProvider) GetCapabilities() providers.Capabilities {
	models := cp.api.GetCommonModels()
	return providers.Capabilities{
		MaxTokens:         getMaxTokensForProvider(cp.name),
		SupportsBatch:     true,
		SupportsStreaming: false, // For now, streaming is not implemented
		Models:            models,
		Pricing:           getPricingForProvider(cp.name),
	}
}

// getMaxTokensForProvider returns max tokens for each provider
func getMaxTokensForProvider(providerName string) int {
	switch providerName {
	case "openai":
		return 4096
	case "claude":
		return 8192
	case "gemini":
		return 8192
	case "deepseek":
		return 4096
	case "ollama":
		return 2048
	default:
		return 2048
	}
}

// getPricingForProvider returns pricing information for each provider
func getPricingForProvider(providerName string) *providers.Pricing {
	switch providerName {
	case "openai":
		return &providers.Pricing{
			InputCostPer1K:  0.0015,
			OutputCostPer1K: 0.002,
			Currency:        "USD",
		}
	case "claude":
		return &providers.Pricing{
			InputCostPer1K:  0.003,
			OutputCostPer1K: 0.015,
			Currency:        "USD",
		}
	case "gemini":
		return &providers.Pricing{
			InputCostPer1K:  0.000125,
			OutputCostPer1K: 0.000375,
			Currency:        "USD",
		}
	case "deepseek":
		return &providers.Pricing{
			InputCostPer1K:  0.00014,
			OutputCostPer1K: 0.00028,
			Currency:        "USD",
		}
	case "ollama":
		return &providers.Pricing{
			InputCostPer1K:  0.0,
			OutputCostPer1K: 0.0,
			Currency:        "USD", // Free local model
		}
	default:
		return nil
	}
}

// NewProviders creates all available providers based on configuration
func NewProviders(config types.IConfig) []providers.Provider {
	var activeProviders []providers.Provider

	// List of all supported providers
	providerConfigs := []struct {
		name string
		key  string
	}{
		{"openai", "openai"},
		{"claude", "claude"},
		{"gemini", "gemini"},
		{"deepseek", "deepseek"},
		{"ollama", "ollama"},
	}

	for _, pc := range providerConfigs {
		// Check if provider is configured
		if config.GetAPIKey(pc.key) != "" {
			api := config.GetAPIConfig(pc.name)
			if api != nil && api.IsAvailable() {
				provider := &concreteProvider{
					name:   pc.name,
					api:    api,
					config: config,
				}
				activeProviders = append(activeProviders, provider)
			}
		}
	}

	return activeProviders
}

// Individual provider constructors for engine initialization

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string) providers.Provider {
	api := types.NewOpenAIAPI(apiKey)
	return &concreteProvider{
		name: "openai",
		api:  api,
	}
}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider(apiKey string) providers.Provider {
	api := types.NewClaudeAPI(apiKey)
	return &concreteProvider{
		name: "claude",
		api:  api,
	}
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider(apiKey string) providers.Provider {
	api := types.NewGeminiAPI(apiKey)
	return &concreteProvider{
		name: "gemini",
		api:  api,
	}
}

// NewDeepSeekProvider creates a new DeepSeek provider
func NewDeepSeekProvider(apiKey string) providers.Provider {
	api := types.NewDeepSeekAPI(apiKey)
	return &concreteProvider{
		name: "deepseek",
		api:  api,
	}
}

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider() providers.Provider {
	api := types.NewOllamaAPI("http://localhost:11434")
	return &concreteProvider{
		name: "ollama",
		api:  api,
	}
}
