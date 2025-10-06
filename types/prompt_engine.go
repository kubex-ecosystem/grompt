package types

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kubex-ecosystem/grompt/internal/gateway/middleware"
	"github.com/kubex-ecosystem/grompt/internal/gateway/registry"
	providers "github.com/kubex-ecosystem/grompt/internal/types"
)

// PromptEngine exposes legacy prompt-processing capabilities.
type PromptEngine interface {
	ProcessPrompt(tmpl string, vars map[string]interface{}) (*Result, error)
	GetProviders() []Provider
	GetHistory() []Result
	SaveToHistory(prompt, response string) error
	BatchProcess(prompts []string, vars map[string]interface{}) ([]Result, error)
}

// NewPromptEngine returns a legacy-compatible engine backed by the new gateway stack.
func NewPromptEngine(cfg Config) PromptEngine {
	impl, ok := cfg.(*configImpl)
	if !ok {
		// Attempt to wrap arbitrary implementations with sane defaults.
		wrapped := newConfig()
		wrapped.port = cfg.GetPort()
		for _, name := range legacyProviders {
			if key := cfg.GetAPIKey(name); key != "" {
				wrapped.apiKeys[name] = key
			}
			if endpoint := cfg.GetAPIEndpoint(name); endpoint != "" {
				wrapped.endpoints[name] = endpoint
			}
		}
		impl = wrapped
	}

	engine := &promptEngine{
		cfg:       impl,
		history:   newHistoryStore(impl.historyLimit),
		providers: map[string]Provider{},
	}

	if err := engine.bootstrap(); err != nil {
		// Surface bootstrap errors through a stub engine that will report the failure at runtime.
		engine.bootstrapErr = err
	}

	impl.attachEngine(engine)
	return engine
}

// ---------- Internal engine implementation ----------

type promptEngine struct {
	cfg          *configImpl
	registry     *registry.Registry
	middleware   *middleware.ProductionMiddleware
	history      *historyStore
	providers    map[string]Provider
	providersMu  sync.RWMutex
	bootstrapErr error
}

func (pe *promptEngine) bootstrap() error {
	reg, err := registry.FromConfig(pe.cfg.registryConfig())
	if err != nil {
		return err
	}

	pm := middleware.NewProductionMiddleware(middleware.DefaultProductionConfig())

	for _, name := range reg.ListProviders() {
		pm.RegisterProvider(name)
	}

	pe.registry = reg
	pe.middleware = pm
	return nil
}

func (pe *promptEngine) ensureBootstrap() error {
	if pe.bootstrapErr != nil {
		return pe.bootstrapErr
	}
	if pe.registry == nil {
		return errors.New("prompt engine not initialized")
	}
	return nil
}

func (pe *promptEngine) ProcessPrompt(tmpl string, vars map[string]interface{}) (*Result, error) {
	if err := pe.ensureBootstrap(); err != nil {
		return nil, err
	}

	promptText, err := executeTemplate(tmpl, vars)
	if err != nil {
		return nil, err
	}

	providerName := pe.cfg.defaultProvider
	if providerOverride, ok := vars["provider"].(string); ok && providerOverride != "" {
		providerName = providerOverride
	}

	if providerName == "" {
		names := pe.registry.ListProviders()
		if len(names) == 0 {
			return nil, errors.New("no providers configured")
		}
		providerName = names[0]
	}

	result, err := pe.invokeProvider(context.Background(), providerName, promptText, vars)
	if err != nil {
		return nil, err
	}

	pe.history.add(*result)
	return result, nil
}

func (pe *promptEngine) BatchProcess(prompts []string, vars map[string]interface{}) ([]Result, error) {
	if err := pe.ensureBootstrap(); err != nil {
		return nil, err
	}

	results := make([]Result, len(prompts))
	var wg sync.WaitGroup
	errCh := make(chan error, len(prompts))

	for idx, tmpl := range prompts {
		wg.Add(1)
		go func(i int, template string) {
			defer wg.Done()

			res, err := pe.ProcessPrompt(template, vars)
			if err != nil {
				errCh <- fmt.Errorf("prompt %d: %w", i, err)
				return
			}
			results[i] = *res
		}(idx, tmpl)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		var builder strings.Builder
		for err := range errCh {
			if builder.Len() > 0 {
				builder.WriteString("; ")
			}
			builder.WriteString(err.Error())
		}
		return results, errors.New(builder.String())
	}

	return results, nil
}

func (pe *promptEngine) GetProviders() []Provider {
	if err := pe.ensureBootstrap(); err != nil {
		return nil
	}

	pe.providersMu.RLock()
	cached := make([]Provider, 0, len(pe.providers))
	for _, p := range pe.providers {
		cached = append(cached, p)
	}
	pe.providersMu.RUnlock()

	if len(cached) > 0 {
		sort.Slice(cached, func(i, j int) bool {
			return cached[i].Name() < cached[j].Name()
		})
		return cached
	}

	names := pe.registry.ListProviders()
	providers := make([]Provider, 0, len(names))
	for _, name := range names {
		providers = append(providers, &providerAdapter{
			name:    name,
			engine:  pe,
			version: "v1",
		})
	}

	pe.providersMu.Lock()
	for _, p := range providers {
		pe.providers[p.Name()] = p
	}
	pe.providersMu.Unlock()

	sort.Slice(providers, func(i, j int) bool { return providers[i].Name() < providers[j].Name() })
	return providers
}

func (pe *promptEngine) GetHistory() []Result {
	return pe.history.snapshot()
}

func (pe *promptEngine) SaveToHistory(prompt, response string) error {
	pe.history.add(Result{
		ID:        uuid.NewString(),
		Prompt:    prompt,
		Response:  response,
		Provider:  "manual",
		Timestamp: time.Now(),
	})
	return nil
}

// executeRaw allows provider adapters and API configs to run direct prompts without templating.
func (pe *promptEngine) executeRaw(providerName, prompt string) (string, error) {
	res, err := pe.invokeProvider(context.Background(), providerName, prompt, map[string]interface{}{})
	if err != nil {
		return "", err
	}
	return res.Response, nil
}

func (pe *promptEngine) invokeProvider(ctx context.Context, providerName, prompt string, vars map[string]interface{}) (*Result, error) {
	if err := pe.ensureBootstrap(); err != nil {
		return nil, err
	}

	provider := pe.registry.Resolve(providerName)
	if provider == nil {
		return nil, fmt.Errorf("provider '%s' not found", providerName)
	}

	if err := provider.Available(); err != nil {
		return nil, err
	}

	model := pe.cfg.defaultModels[providerName]
	if override, ok := vars["model"].(string); ok && override != "" {
		model = override
	}

	req := providers.ChatRequest{
		Provider: providerName,
		Model:    model,
		Messages: []providers.Message{{Role: "user", Content: prompt}},
		Temp:     pe.cfg.defaultTemperature,
		Stream:   false,
		Meta:     map[string]any{},
	}

	var responseContent bytes.Buffer
	var usage *providers.Usage

	operation := func() error {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, pe.cfg.requestTimeout())
		defer cancel()

		stream, err := provider.Chat(ctxWithTimeout, req)
		if err != nil {
			return err
		}

		for chunk := range stream {
			if chunk.Error != "" {
				return errors.New(chunk.Error)
			}
			if chunk.Content != "" {
				responseContent.WriteString(chunk.Content)
			}
			if chunk.Done && chunk.Usage != nil {
				usage = chunk.Usage
			}
		}

		return nil
	}

	if pe.middleware != nil {
		if err := pe.middleware.WrapProvider(providerName, operation); err != nil {
			return nil, err
		}
	} else if err := operation(); err != nil {
		return nil, err
	}

	metadata := map[string]interface{}{}
	if usage != nil {
		metadata["usage"] = usage
	}

	result := &Result{
		ID:        uuid.NewString(),
		Prompt:    prompt,
		Response:  responseContent.String(),
		Provider:  providerName,
		Variables: vars,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}

	return result, nil
}
