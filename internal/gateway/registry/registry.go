// Package registry provides provider registration and resolution functionality.
package registry

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kubex-ecosystem/grompt/internal/interfaces"
	"github.com/kubex-ecosystem/grompt/internal/providers"
	"github.com/kubex-ecosystem/grompt/internal/types"
	"gopkg.in/yaml.v3"
)

// Registry manages provider registration and resolution
type Registry struct {
	cfg       *types.Config
	providers map[string]interfaces.Provider
}

// Load creates a new registry from a YAML configuration file
func Load(path string) (*Registry, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var cfg types.Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return FromConfig(&cfg)
}

// FromConfig builds a registry from an in-memory configuration structure.
func FromConfig(cfg *types.Config) (*Registry, error) {
	r := &Registry{
		cfg:       cfg,
		providers: make(map[string]interfaces.Provider),
	}

	if err := r.initializeProviders(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Registry) initializeProviders() error {
	for name, pc := range r.cfg.Providers {
		switch pc.Type() {
		case "openai":
			key := os.Getenv(pc.KeyEnv())
			if key == "" {
				fmt.Printf("Warning: Skipping OpenAI provider '%s' - no API key found in %s\n", name, pc.KeyEnv())
				continue
			}
			p := providers.NewOpenAIProvider(key)
			// if err != nil {
			// 	return fmt.Errorf("failed to create OpenAI provider %s: %w", name, err)
			// }
			r.providers[name] = p
		case "gemini":
			key := os.Getenv(pc.KeyEnv())
			if key == "" {
				fmt.Printf("Warning: Skipping Gemini provider '%s' - no API key found in %s\n", name, pc.KeyEnv())
				continue
			}
			p := providers.NewGeminiProvider(key)
			// if err != nil {
			// 	return fmt.Errorf("failed to create Gemini provider %s: %w", name, err)
			// }
			r.providers[name] = p
		case "anthropic":
			key := os.Getenv(pc.KeyEnv())
			if key == "" {
				fmt.Printf("Warning: Skipping Anthropic provider '%s' - no API key found in %s\n", name, pc.KeyEnv())
				continue
			}
			p := providers.NewClaudeProvider(key)
			// if err != nil {
			// 	return fmt.Errorf("failed to create Anthropic provider %s: %w", name, err)
			// }
			r.providers[name] = p
		case "openrouter":
			fmt.Printf("Warning: Skipping OpenRouter provider '%s' - not yet implemented\n", name)
		case "ollama":
			fmt.Printf("Warning: Skipping Ollama provider '%s' - not yet implemented\n", name)
		default:
			fmt.Printf("Warning: Skipping provider '%s' - unknown type '%s'\n", name, pc.Type())
		}
	}

	return nil
}

// FromRuntimeConfig constrói um registry utilizando a configuração já carregada
// pelo servidor principal, evitando depender de arquivos YAML separados.
func FromRuntimeConfig(cfg interfaces.IConfig) (*Registry, error) {
	if cfg == nil {
		return nil, fmt.Errorf("runtime config is nil")
	}

	providerMap := make(map[string]interfaces.Provider)

	if key := cfg.GetAPIKey("openai"); key != "" {
		providerMap["openai"] = providers.NewOpenAIProvider(key)
	}
	if key := cfg.GetAPIKey("claude"); key != "" {
		providerMap["claude"] = providers.NewClaudeProvider(key)
	}
	if key := cfg.GetAPIKey("gemini"); key != "" {
		providerMap["gemini"] = providers.NewGeminiProvider(key)
	}
	if key := cfg.GetAPIKey("deepseek"); key != "" {
		providerMap["deepseek"] = providers.NewDeepSeekProvider(key)
	}
	if key := cfg.GetAPIKey("chatgpt"); key != "" {
		providerMap["chatgpt"] = providers.NewChatGPTProvider(key)
	}
	if endpoint := cfg.GetAPIEndpoint("ollama"); endpoint != "" {
		providerMap["ollama"] = &types.ProviderImpl{
			VName: "ollama",
			VAPI:  types.NewOllamaAPI(endpoint),
		}
	}

	if len(providerMap) == 0 {
		return nil, fmt.Errorf("no providers available in runtime config")
	}

	runtimeCfg := &types.Config{
		Providers: providerMap,
	}

	return &Registry{
		cfg:       runtimeCfg,
		providers: providerMap,
	}, nil
}

// Resolve returns a provider by name
func (r *Registry) Resolve(name string) providers.Provider {
	return r.providers[name]
}

// ListProviders returns all available provider names
func (r *Registry) ListProviders() []string {
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}

// GetConfig returns the provider configuration
func (r *Registry) GetConfig() interfaces.IConfig {
	return r.cfg
}

func (r *Registry) ResolveProvider(name string) interfaces.Provider { return r.providers[name] }

func (r *Registry) Config() *types.Config { return r.cfg } // <- usado por /v1/providers

func (r *Registry) Chat(ctx context.Context, req interfaces.ChatRequest) (<-chan interfaces.ChatChunk, error) {
	p := r.ResolveProvider(req.Provider)
	if p == nil {
		return nil, fmt.Errorf("provider '%s' not found", req.Provider)
	}
	return p.Chat(ctx, req)
}
func (r *Registry) Notify(ctx context.Context, event interfaces.NotificationEvent) error {
	p := r.ResolveProvider(event.Type)
	if p == nil {
		return fmt.Errorf("provider '%s' not found", event.Type)
	}
	return p.Notify(ctx, event)
}
func (r *Registry) AddProvider(provider interfaces.Provider) error {
	if _, exists := r.providers[provider.Name()]; exists {
		return fmt.Errorf("provider '%s' already exists", provider.Name())
	}
	r.providers[provider.Name()] = provider
	return nil
}
func (r *Registry) BatchProcess(ctx context.Context, requests []string, options map[string]interface{}) ([]interfaces.Result, error) {
	results := make([]interfaces.Result, 0, len(requests))
	for _, req := range requests {
		p := r.ResolveProvider(req)
		if p == nil {
			return nil, fmt.Errorf("provider '%s' not found", req)
		}
		chunks, err := p.Chat(ctx, interfaces.ChatRequest{
			Provider: req,
			Messages: options["messages"].([]interfaces.Message),
		})
		if err != nil {
			return nil, fmt.Errorf("error processing request with provider '%s': %w", req, err)
		}
		var responseChunks []interfaces.ChatChunk
		for chunk := range chunks {
			responseChunks = append(responseChunks, chunk)
		}
		result := interfaces.Result{
			Provider: req,
			Response: strings.Join(extractContents(responseChunks), ""),
		}
		results = append(results, result)
	}
	return results, nil
}

func (r *Registry) GetCapabilities(ctx context.Context) *interfaces.Capabilities {
	combined := &interfaces.Capabilities{
		MaxTokens:         0,
		SupportsBatch:     false,
		SupportsStreaming: false,
		Models: map[string]any{},
	}

	for _, provider := range r.providers {
		caps := provider.GetCapabilities(ctx)
		if caps == nil {
			continue
		}
		if caps.MaxTokens > combined.MaxTokens {
			combined.MaxTokens = caps.MaxTokens
		}
		if caps.SupportsBatch {
			combined.SupportsBatch = true
		}
		if caps.SupportsStreaming {
			combined.SupportsStreaming = true
		}
		for modelName, model := range caps.Models {
			combined.Models[modelName] = model
		}
		// Note: Pricing aggregation logic can be more complex based on requirements
	}

	return combined
}

func (r *Registry) Close() error {
	// Close all providers if they implement io.Closer
	for name, provider := range r.providers {
		if closer, ok := provider.(interface{ Close() error }); ok {
			if err := closer.Close(); err != nil {
				fmt.Printf("Warning: failed to close provider '%s': %v\n", name, err)
			}
		}
	}
	return nil
}

func (r *Registry) GetHistory() []interfaces.Result {
	// Dummy implementation - replace with actual history retrieval logic
	return []interfaces.Result{}
}

func (r *Registry) GetProviders() []interfaces.Provider {
	providersList := make([]interfaces.Provider, 0, len(r.providers))
	for _, provider := range r.providers {
		providersList = append(providersList, provider)
	}
	return providersList
}

func (r *Registry) GetRegistry() interfaces.Provider {
	return r
}

func (r *Registry) InvokeProvider(ctx context.Context, service string, method string, params map[string]interface{}) (*interfaces.Result, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *Registry) ProcessPrompt(ctx context.Context, prompt string, options map[string]interface{}) (*interfaces.Result, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *Registry) SaveToHistory(ctx context.Context, key string, value string) error {
	return nil
}

func (r *Registry) Execute(ctx context.Context, command string, args map[string]any) (*interfaces.Result, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *Registry) IsAvailable() bool {
	return len(r.providers) > 0
}

func (r *Registry) KeyEnv() string {
	return ""
}

func (r *Registry) Type() string {
	return "registry"
}

func (r *Registry) Name() string {
	return "registry"
}

func (r *Registry) Version() string {
	return "1.0.0"
}

func extractContents(chunks []interfaces.ChatChunk) []string {
	contents := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		contents = append(contents, chunk.Content)
	}
	return contents
}
