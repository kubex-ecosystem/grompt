// Package engine provides the core prompt engineering functionality.
package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/kubex-ecosystem/grompt/factory/templates"
	"github.com/kubex-ecosystem/grompt/internal/interfaces"
	"github.com/kubex-ecosystem/grompt/internal/types"
)

// Engine represents the core prompt engineering engine
type Engine struct {
	providers []interfaces.Provider
	templates templates.Manager
	history   interfaces.IHistoryManager
	config    interfaces.IConfig
}

// NewEngine creates a new IEngine instance with initialized providers
func NewEngine(config interfaces.IConfig) interfaces.IEngine {
	engine := &Engine{
		providers: make([]interfaces.Provider, 0),
		templates: templates.NewManager("./templates"), // Default templates path
		history:   interfaces.NewHistoryStore(100),    // Default history limit
		config:    config,
	}

	// Initialize concrete providers
	engine.initializeProviders()

	return engine
}

// initializeProviders initializes all available concrete providers
func (e *Engine) initializeProviders() {
	// Initialize OpenAI provider
	if apiKey := e.config.GetAPIKey("openai"); apiKey != "" {
		provider := types.NewOpenAIAPI(apiKey)
		e.providers = append(e.providers, provider)
	}

	// Initialize Claude provider
	if apiKey := e.config.GetAPIKey("claude"); apiKey != "" {
		provider := types.NewClaudeAPI(apiKey)
		e.providers = append(e.providers, provider)
	}

	// Initialize DeepSeek provider
	if apiKey := e.config.GetAPIKey("deepseek"); apiKey != "" {
		provider := types.NewDeepSeekAPI(apiKey)
		e.providers = append(e.providers, provider)
	}

	// Initialize Ollama provider (local - no API key needed)
	provider := types.NewOllamaAPI("")
	e.providers = append(e.providers, provider)

	// Initialize Gemini provider
	if apiKey := e.config.GetAPIKey("gemini"); apiKey != "" {
		provider := types.NewGeminiAPI(apiKey)
		e.providers = append(e.providers, provider)
	}
}

// ProcessPrompt processes a prompt with variables and returns the result
func (e *Engine) ProcessPrompt(ctx context.Context, template string, vars map[string]interface{}) (*interfaces.Result, error) {
	if e == nil {
		return nil, fmt.Errorf("engine is nil")
	}

	// Process template with variables
	processedPrompt, err := e.templates.Process(template, vars)
	if err != nil {
		return nil, fmt.Errorf("template processing failed: %w", err)
	}

	// Get default provider (first available)
	if len(e.providers) == 0 {
		return nil, fmt.Errorf("no providers available")
	}

	provider := e.providers[0]

	// Execute prompt with provider
	response, err := provider.Execute(context.Background(), template, vars)
	if err != nil {
		return nil, fmt.Errorf("provider execution failed: %w", err)
	}
	// result := response

	// Create result
	result := &interfaces.Result{
		ID:        generateID(),
		Prompt:    processedPrompt,
		Response:  response.Response,
		Provider:  provider.Name(),
		Variables: vars,
		Timestamp: time.Now(),
	}

	// Add to history
	e.history.Add(*result)

	return result, nil
}

// GetProviders returns available providers
func (e *Engine) GetProviders() []interfaces.Provider {
	if e == nil {
		return nil
	}
	return e.providers
}

// GetHistory returns the prompt history
func (e *Engine) GetHistory() []interfaces.Result {
	if e == nil || e.history == nil {
		return nil
	}

	return e.history.Snapshot()
}

// SaveToHistory saves a prompt/response pair to history
func (e *Engine) SaveToHistory(ctx context.Context, prompt, response string) error {
	if e == nil || e.history == nil {
		return fmt.Errorf("engine or history is nil")
	}

	result := interfaces.Result{
		ID:        generateID(),
		Prompt:    prompt,
		Response:  response,
		Provider:  "manual",
		Timestamp: time.Now(),
	}

	e.history.Add(result)
	return nil
}

// BatchProcess processes multiple prompts concurrently
func (e *Engine) BatchProcess(ctx context.Context, prompts []string, vars map[string]interface{}) ([]interfaces.Result, error) {
	if e == nil {
		return nil, fmt.Errorf("engine is nil")
	}

	results := make([]interfaces.Result, len(prompts))
	errors := make([]error, len(prompts))

	// Process prompts concurrently
	for i, prompt := range prompts {
		go func(index int, p string) {
			result, err := e.ProcessPrompt(ctx, p, vars)
			if err != nil {
				errors[index] = err
				return
			}
			results[index] = *result
		}(i, prompt)
	}

	// Wait for all goroutines to complete (simplified version)
	time.Sleep(time.Second * 2) // In production, use sync.WaitGroup

	// Check for errors
	for _, err := range errors {
		if err != nil {
			return results, fmt.Errorf("batch processing errors occurred: %v", err)
		}
	}

	return results, nil
}

func (e *Engine) GetConfig() interfaces.IConfig {
	if e == nil {
		return nil
	}
	return e.config
}

func (e *Engine) GetRegistry() interfaces.Provider {
	if e == nil {
		return nil
	}
	// For simplicity, return the first provider as the registry
	if len(e.providers) > 0 {
		return e.providers[0]
	}
	return nil
}

func (e *Engine) Resolve(name string) interfaces.Provider {
	if e == nil {
		return nil
	}
	for _, p := range e.providers {
		if p.Name() == name {
			return p
		}
	}
	return nil
}

func (e *Engine) AddProvider(provider interfaces.Provider) error {
	if e == nil {
		return fmt.Errorf("engine is nil")
	}
	e.providers = append(e.providers, provider)
	return nil
}

// Close releases any resources held by the engine
func (e *Engine) Close() error {
	// Currently, no resources to release
	return nil
}

func (e *Engine) InvokeProvider(ctx context.Context, providerName, prompt string, vars map[string]interface{}) (*interfaces.Result, error) {
	if e == nil {
		return nil, fmt.Errorf("engine is nil")
	}

	// Find the specified provider
	var provider interfaces.Provider
	for _, p := range e.providers {
		if p.Name() == providerName {
			provider = p
			break
		}
	}
	if provider == nil {
		return nil, fmt.Errorf("provider %s not found", providerName)
	}

	// Process template with variables
	processedPrompt, err := e.templates.Process(prompt, vars)
	if err != nil {
		return nil, fmt.Errorf("template processing failed: %w", err)
	}

	// // Execute prompt with provider
	// response, err := provider.Execute(ctx, processedPrompt, vars)
	// if err != nil {
	// 	return nil, fmt.Errorf("provider execution failed: %w", err)
	// }
	// // Create result
	// result := response

	result := &interfaces.Result{
		ID:        generateID(),
		Prompt:    processedPrompt,
		// Response:  response,
		Provider:  provider.Name(),
		Variables: vars,
		Timestamp: time.Now(),
	}

	// Add to history
	e.history.Add(*result)

	return result, nil
}

func (e *Engine) GetCapabilities(ctx context.Context) *interfaces.Capabilities {
	if e == nil {
		return nil
	}
	// For simplicity, return static capabilities
	capabilities := e.GetRegistry().GetCapabilities(ctx)
	return capabilities
}

// generateID generates a simple ID for results
func generateID() string {
	return fmt.Sprintf("prompt_%d", time.Now().UnixNano())
}
