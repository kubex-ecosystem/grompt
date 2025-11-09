// Package types defines interfaces and types for AI providers
package types

import (
	"context"
	"fmt"

	"github.com/kubex-ecosystem/grompt/internal/interfaces"
)

// ProviderImpl wraps the types.IAPIConfig to implement providers.Provider
type ProviderImpl struct {
	VName    string
	VVersion string
	VKeyEnv  string
	VType    string
	VAPI     interfaces.IAPIConfig
	VConfig  interfaces.IConfig
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
func (cp *ProviderImpl) Execute(ctx context.Context, prompt string, vars map[string]any) (*interfaces.Result, error) {
	if cp == nil || cp.VAPI == nil {
		return nil, fmt.Errorf("provider is not available")
	}
	response, err := cp.VAPI.Complete(prompt, 2048, "") // Default max tokens
	if err != nil {
		return nil, err
	}
	return &interfaces.Result{Response: response}, nil
}

// IsAvailable checks if the provider is configured and ready
func (cp *ProviderImpl) IsAvailable() bool {
	if cp == nil || cp.VAPI == nil {
		return false
	}
	return cp.VAPI.IsAvailable()
}

func (cp *ProviderImpl) KeyEnv() string {
	if cp == nil {
		return ""
	}
	return cp.VKeyEnv
}

func (cp *ProviderImpl) Type() string {
	if cp == nil || cp.VType == "" {
		return "text"
	}
	return cp.VType
}

// GetCapabilities returns provider-specific capabilities
func (cp *ProviderImpl) GetCapabilities(ctx context.Context) *interfaces.Capabilities {
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
	return &interfaces.Capabilities{
		MaxTokens:         getMaxTokensForProvider(cp.VName),
		SupportsBatch:     true,
		SupportsStreaming: false, // For now, streaming is not implemented
		Models:            models,
		Pricing:           getPricingForProvider(cp.VName),
	}
}

// Chat sends a chat request to the provider and returns a stream of responses
func (cp *ProviderImpl) Chat(ctx context.Context, req interfaces.ChatRequest) (<-chan interfaces.ChatChunk, error) {
	if cp == nil || cp.VAPI == nil {
		return nil, fmt.Errorf("provider is not available")
	}
	if _, ok := cp.VAPI.(interfaces.Provider); !ok {
		return nil, fmt.Errorf("provider does not support chat")
	}
	chatProvider := cp.VAPI.(interfaces.Provider)
	return chatProvider.Chat(ctx, req)
}

// Notify sends a notification event to the provider if supported
func (cp *ProviderImpl) Notify(ctx context.Context, event interfaces.NotificationEvent) error {
	if cp == nil || cp.VAPI == nil {
		return fmt.Errorf("provider is not available")
	}
	if notifier, ok := cp.VAPI.(interfaces.Notifier); ok {
		return notifier.Notify(ctx, event)
	}
	return fmt.Errorf("provider does not support notifications")
}


func getMaxTokensForProvider(name string) int {
	switch name {
	case "openai":
		return 8192
	case "claude":
		return 9000
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

func getPricingForProvider(name string) *interfaces.Pricing {
	switch name {
	case "openai":
		return &interfaces.Pricing{
			InputCostPer1K: 		 0.03,
			OutputCostPer1K: 		 0.06,
			Currency: 				"USD",
		}
	case "claude":
		return &interfaces.Pricing{
			InputCostPer1K: 		 0.015,
			OutputCostPer1K: 		 0.03,
			Currency: 				"USD",
		}
	case "gemini":
		return &interfaces.Pricing{
			InputCostPer1K: 		 0.02,
			OutputCostPer1K: 		 0.04,
			Currency: 				"USD",
		}
	case "deepseek":
		return &interfaces.Pricing{
			InputCostPer1K: 		 0.01,
			OutputCostPer1K: 		 0.02,
			Currency: 				"USD",
		}
	case "ollama":
		return &interfaces.Pricing{
			InputCostPer1K: 		 0.0,
			OutputCostPer1K: 		 0.0,
			Currency: 				"USD",
		}
	default:
		return nil
	}
}
