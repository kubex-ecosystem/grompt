// Package providers provides concrete implementations of AI providers.
package providers

import (
	"github.com/kubex-ecosystem/gemx/grompt/internal/types"
)

// Provider defines the interface for AI providers.
type Provider interface {
	Name() string
	Version() string
	Execute(prompt string) (string, error)
	IsAvailable() bool
	GetCapabilities() *types.Capabilities
}

// Individual provider constructors for engine initialization

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string) Provider {
	api := types.NewOpenAIAPI(apiKey)
	return &types.ProviderImpl{
		VName: "openai",
		VAPI:  api,
	}
}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider(apiKey string) Provider {
	api := types.NewClaudeAPI(apiKey)
	return &types.ProviderImpl{
		VName: "claude",
		VAPI:  api,
	}
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider(apiKey string) Provider {
	api := types.NewGeminiAPI(apiKey)
	return &types.ProviderImpl{
		VName: "gemini",
		VAPI:  api,
	}
}

// NewDeepSeekProvider creates a new DeepSeek provider
func NewDeepSeekProvider(apiKey string) Provider {
	api := types.NewDeepSeekAPI(apiKey)
	return &types.ProviderImpl{
		VName: "deepseek",
		VAPI:  api,
	}
}

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider() Provider {
	api := types.NewOllamaAPI("http://localhost:11434")
	return &types.ProviderImpl{
		VName: "ollama",
		VAPI:  api,
	}
}

func NewChatGPTProvider(apiKey string) Provider {
	api := types.NewChatGPTAPI(apiKey)
	return &types.ProviderImpl{
		VName: "chatgpt",
		VAPI:  api,
	}
}

func NewProvider(name, apiKey, version string, cfg types.IConfig) Provider {
	return &types.ProviderImpl{
		VName:    name,
		VVersion: version,
		VAPI:     cfg.GetAPIConfig(name),
		VConfig:  cfg,
	}
}
