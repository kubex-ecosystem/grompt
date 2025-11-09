package interfaces

import "context"

type IEngine interface {
	// ProcessPrompt processes a prompt with variables and returns the result
	ProcessPrompt(ctx context.Context, template string, vars map[string]interface{}) (*Result, error)

	// GetProviders returns available providers
	GetProviders() []Provider

	// GetHistory returns the prompt history
	GetHistory() []Result

	// SaveToHistory saves a prompt/response pair to history
	SaveToHistory(ctx context.Context, prompt, response string) error

	// BatchProcess processes multiple prompts concurrently
	BatchProcess(ctx context.Context, prompts []string, vars map[string]interface{}) ([]Result, error)

	// InvokeProvider invokes a specific provider with a prompt and variables
	InvokeProvider(ctx context.Context, providerName, prompt string, vars map[string]interface{}) (*Result, error)

	// Close releases any resources held by the engine
	Close() error

	// GetCapabilities returns the capabilities of the engine
	GetCapabilities(ctx context.Context) *Capabilities

	// registry holds the provider registry
	GetRegistry() Provider

	// GetConfig returns the engine configuration
	GetConfig() IConfig

	// Resolve finds a provider by name
	Resolve(name string) Provider

	// AddProvider adds a new provider to the engine
	AddProvider(provider Provider) error
}
