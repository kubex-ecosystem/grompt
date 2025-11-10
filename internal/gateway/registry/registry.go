// Package registry provides provider registration and resolution functionality.
package registry

import (
	"context"
	"fmt"
	"os"

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
func (r *Registry) GetConfig() *types.Config {
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

// /v1/chat/completions — SSE endpoints
