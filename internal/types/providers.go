// Package types defines interfaces and types for AI providers
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

// Chat sends a chat request to the provider and returns a stream of responses
func (cp *ProviderImpl) Chat(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error) {
	if cp == nil || cp.VAPI == nil {
		return nil, fmt.Errorf("provider is not available")
	}
	if _, ok := cp.VAPI.(Provider); !ok {
		return nil, fmt.Errorf("provider does not support chat")
	}
	chatProvider := cp.VAPI.(Provider)
	return chatProvider.Chat(ctx, req)
}

// Notify sends a notification event to the provider if supported
func (cp *ProviderImpl) Notify(ctx context.Context, event NotificationEvent) error {
	if cp == nil || cp.VAPI == nil {
		return fmt.Errorf("provider is not available")
	}
	if notifier, ok := cp.VAPI.(Notifier); ok {
		return notifier.Notify(ctx, event)
	}
	return fmt.Errorf("provider does not support notifications")
}

// Provider interface defines the contract for AI providers
type Provider interface {
	Name() string
	Chat(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error)
	Available() error
	Notify(ctx context.Context, event NotificationEvent) error
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Headers  map[string]string `json:"-"`
	Provider string            `json:"provider"`
	Model    string            `json:"model"`
	Messages []Message         `json:"messages"`
	Temp     float32           `json:"temperature"`
	Stream   bool              `json:"stream"`
	Meta     map[string]any    `json:"meta"`
}

// Message represents a single chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Usage represents token usage and cost information
type Usage struct {
	Completion int     `json:"completion_tokens"`
	Prompt     int     `json:"prompt_tokens"`
	Tokens     int     `json:"tokens"`
	Ms         int64   `json:"latency_ms"`
	CostUSD    float64 `json:"cost_usd"`
	Provider   string  `json:"provider"`
	Model      string  `json:"model"`
}

// ChatChunk represents a streaming response chunk
type ChatChunk struct {
	Content  string    `json:"content,omitempty"`
	Done     bool      `json:"done"`
	Usage    *Usage    `json:"usage,omitempty"`
	Error    string    `json:"error,omitempty"`
	ToolCall *ToolCall `json:"toolCall,omitempty"`
}
