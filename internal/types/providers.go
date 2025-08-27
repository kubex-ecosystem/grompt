package types

import "fmt"

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

// ProviderImpl wraps the types.IAPIConfig to implement providers.Provider
type ProviderImpl struct {
	VName    string
	VVersion string
	VAPI     IAPIConfig
	VConfig  IConfig
}

// Name returns the provider name
func (cp *ProviderImpl) Name() string {
	return cp.VName
}

// Version returns the provider version
func (cp *ProviderImpl) Version() string {
	return cp.VVersion
}

// Execute sends a prompt to the provider and returns the response
func (cp *ProviderImpl) Execute(prompt string) (string, error) {
	if cp == nil || cp.VAPI == nil {
		return "", fmt.Errorf("provider is not available")
	}
	return cp.VAPI.Complete(prompt, 2048, "") // Default max tokens
}

// IsAvailable checks if the provider is configured and ready
func (cp *ProviderImpl) IsAvailable() bool {
	if cp == nil || cp.VAPI == nil {
		return false
	}
	return cp.VAPI.IsAvailable()
}

// GetCapabilities returns provider-specific capabilities
func (cp *ProviderImpl) GetCapabilities() *Capabilities {
	if cp == nil {
		return nil
	}
	if cp.VAPI == nil {
		if cp.VConfig != nil {
			switch cp.VName {
			case "openai":
				cp.VAPI = cp.VConfig.GetAPIConfig("openai")
			case "claude":
				cp.VAPI = cp.VConfig.GetAPIConfig("claude")
			case "gemini":
				cp.VAPI = cp.VConfig.GetAPIConfig("gemini")
			case "deepseek":
				cp.VAPI = cp.VConfig.GetAPIConfig("deepseek")
			case "ollama":
				cp.VAPI = cp.VConfig.GetAPIConfig("ollama")
			default:
				return nil // No API config available for this provider
			}
		} else {
			return nil // No API config available
		}
	}
	models, err := cp.VAPI.ListModels()
	if err != nil {
		return nil
	}
	return &Capabilities{
		MaxTokens:         getMaxTokensForProvider(cp.VName),
		SupportsBatch:     true,
		SupportsStreaming: false, // For now, streaming is not implemented
		Models:            models,
		Pricing:           getPricingForProvider(cp.VName),
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
func getPricingForProvider(providerName string) *Pricing {
	switch providerName {
	case "openai":
		return &Pricing{
			InputCostPer1K:  0.0015,
			OutputCostPer1K: 0.002,
			Currency:        "USD",
		}
	case "claude":
		return &Pricing{
			InputCostPer1K:  0.003,
			OutputCostPer1K: 0.015,
			Currency:        "USD",
		}
	case "gemini":
		return &Pricing{
			InputCostPer1K:  0.000125,
			OutputCostPer1K: 0.000375,
			Currency:        "USD",
		}
	case "deepseek":
		return &Pricing{
			InputCostPer1K:  0.00014,
			OutputCostPer1K: 0.00028,
			Currency:        "USD",
		}
	case "ollama":
		return &Pricing{
			InputCostPer1K:  0.0,
			OutputCostPer1K: 0.0,
			Currency:        "USD", // Free local model
		}
	default:
		return nil
	}
}

// NewProviders creates all available providers based on configuration
func NewProviders(config IConfig) []*ProviderImpl {
	var activeProviders []*ProviderImpl

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
		{"chatgpt", "chatgpt"},
	}

	for _, pc := range providerConfigs {
		// Check if provider is configured
		if config.GetAPIKey(pc.key) != "" {
			api := config.GetAPIConfig(pc.name)
			if api != nil && api.IsAvailable() {
				provider := &ProviderImpl{
					VName:   pc.name,
					VAPI:    api,
					VConfig: config,
				}
				activeProviders = append(activeProviders, provider)
			}
		}
	}

	return activeProviders
}

// Individual provider constructors for engine initialization

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string) *ProviderImpl {
	api := NewOpenAIAPI(apiKey)
	return &ProviderImpl{
		VName: "openai",
		VAPI:  api,
	}
}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider(apiKey string) *ProviderImpl {
	api := NewClaudeAPI(apiKey)
	return &ProviderImpl{
		VName: "claude",
		VAPI:  api,
	}
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider(apiKey string) *ProviderImpl {
	api := NewGeminiAPI(apiKey)
	return &ProviderImpl{
		VName: "gemini",
		VAPI:  api,
	}
}

// NewDeepSeekProvider creates a new DeepSeek provider
func NewDeepSeekProvider(apiKey string) *ProviderImpl {
	api := NewDeepSeekAPI(apiKey)
	return &ProviderImpl{
		VName: "deepseek",
		VAPI:  api,
	}
}

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider() *ProviderImpl {
	api := NewOllamaAPI("http://localhost:11434")
	return &ProviderImpl{
		VName: "ollama",
		VAPI:  api,
	}
}

func NewChatGPTProvider(apiKey string) *ProviderImpl {
	api := NewChatGPTAPI(apiKey)
	return &ProviderImpl{
		VName: "chatgpt",
		VAPI:  api,
	}
}

func NewProvider(name, apiKey string, cfg IConfig) *ProviderImpl {
	return &ProviderImpl{
		VName:   name,
		VAPI:    cfg.GetAPIConfig(name),
		VConfig: cfg,
	}
}
