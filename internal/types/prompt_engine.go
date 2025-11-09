package types

import (
	"github.com/kubex-ecosystem/grompt/internal/interfaces"
)

// PromptEngine exposes legacy prompt-processing capabilities.
type PromptEngine interface {
	ProcessPrompt(tmpl string, vars map[string]interface{}) (*interfaces.Result, error)
	GetProviders() []interfaces.Provider
	GetHistory() []interfaces.Result
	SaveToHistory(prompt, response string) error
	BatchProcess(prompts []string, vars map[string]interface{}) ([]interfaces.Result, error)
}

// NewPromptEngine returns a legacy-compatible engine backed by the new gateway stack.
func NewPromptEngine(cfg interfaces.IConfig) PromptEngine {
	// impl, ok := cfg.(*itypes.Config)
	// if !ok {
	// 	// Attempt to wrap arbitrary implementations with sane defaults.
	// 	wrapped := itypes.NewConfig()
	// 	wrapped.port = cfg.GetPort()
	// 	for _, name := range legacyProviders {
	// 		if key := cfg.GetAPIKey(name); key != "" {
	// 			wrapped.apiKeys[name] = key
	// 		}
	// 		if endpoint := cfg.GetAPIEndpoint(name); endpoint != "" {
	// 			wrapped.endpoints[name] = endpoint
	// 		}
	// 	}
	// 	impl = wrapped
	// }

	// engine := &promptEngine{
	// 	cfg:       impl,
	// 	history:   newHistoryStore(impl.historyLimit),
	// 	providers: make(map[string]interfaces.Provider),
	// }

	// if err := engine.bootstrap(); err != nil {
	// 	// Surface bootstrap errors through a stub engine that will report the failure at runtime.
	// 	engine.bootstrapErr = err
	// }

	// impl.attachEngine(engine)
	// return engine
	return nil
}

// ---------- Internal engine implementation ----------

// type promptEngine struct {
// 	cfg          *itypes.Config
// 	registry     *registry.Registry
// 	middleware   *middleware.ProductionMiddleware
// 	history      *interfaces.IHistoryManager
// 	providers    map[string]interfaces.Provider
// 	providersMu  sync.RWMutex
// 	bootstrapErr error
// }

// func (pe *promptEngine) bootstrap() error {
// 	reg, err := registry.FromConfig(pe.cfg.engine.registry.GetConfig())
// 	if err != nil {
// 		return err
// 	}

// 	pm := middleware.NewProductionMiddleware(middleware.DefaultProductionConfig())

// 	for _, name := range reg.ListProviders() {
// 		pm.RegisterProvider(name)
// 	}

// 	pe.registry = reg
// 	pe.middleware = pm
// 	return nil
// }

// func (pe *promptEngine) ensureBootstrap() error {
// 	if pe.bootstrapErr != nil {
// 		return pe.bootstrapErr
// 	}
// 	if pe.registry == nil {
// 		return errors.New("prompt engine not initialized")
// 	}
// 	return nil
// }

// func (pe *promptEngine) ProcessPrompt(tmpl string, vars map[string]interface{}) (*interfaces.Result, error) {
// 	if err := pe.ensureBootstrap(); err != nil {
// 		return nil, err
// 	}

// 	promptText, err := executeTemplate(tmpl, vars)
// 	if err != nil {
// 		return nil, err
// 	}

// 	providerName := pe.cfg.defaultProvider
// 	if providerOverride, ok := vars["provider"].(string); ok && providerOverride != "" {
// 		providerName = providerOverride
// 	}

// 	if providerName == "" {
// 		names := pe.registry.ListProviders()
// 		if len(names) == 0 {
// 			return nil, errors.New("no providers configured")
// 		}
// 		providerName = names[0]
// 	}

// 	result, err := pe.invokeProvider(context.Background(), providerName, promptText, vars)
// 	if err != nil {
// 		return nil, err
// 	}

// 	pe.history.add(*result)
// 	return result, nil
// }

// func (pe *promptEngine) BatchProcess(prompts []string, vars map[string]interface{}) ([]interfaces.Result, error) {
// 	if err := pe.ensureBootstrap(); err != nil {
// 		return nil, err
// 	}

// 	results := make([]interfaces.Result, len(prompts))
// 	var wg sync.WaitGroup
// 	errCh := make(chan error, len(prompts))

// 	for idx, tmpl := range prompts {
// 		wg.Add(1)
// 		go func(i int, template string) {
// 			defer wg.Done()

// 			res, err := pe.ProcessPrompt(template, vars)
// 			if err != nil {
// 				errCh <- fmt.Errorf("prompt %d: %w", i, err)
// 				return
// 			}
// 			results[i] = *res
// 		}(idx, tmpl)
// 	}

// 	wg.Wait()
// 	close(errCh)

// 	if len(errCh) > 0 {
// 		var builder strings.Builder
// 		for err := range errCh {
// 			if builder.Len() > 0 {
// 				builder.WriteString("; ")
// 			}
// 			builder.WriteString(err.Error())
// 		}
// 		return results, errors.New(builder.String())
// 	}

// 	return results, nil
// }

// func (pe *promptEngine) GetProviders() []interfaces.Provider {
// 	if err := pe.ensureBootstrap(); err != nil {
// 		return nil
// 	}

// 	pe.providersMu.RLock()
// 	cached := make([]interfaces.Provider, 0, len(pe.providers))
// 	for _, p := range pe.providers {
// 		cached = append(cached, p)
// 	}
// 	pe.providersMu.RUnlock()

// 	if len(cached) > 0 {
// 		sort.Slice(cached, func(i, j int) bool {
// 			return cached[i].Name() < cached[j].Name()
// 		})
// 		return cached
// 	}

// 	names := pe.registry.ListProviders()
// 	providers := make(map[string]interfaces.Provider)
// 	for _, name := range names {
// 		if prov := pe.registry.Resolve(name); prov != nil {
// 			providers[name] = prov
// 		}
// 	}

// 	pe.providersMu.Lock()
// 	for _, p := range providers {
// 		pe.providers[p.Name()] = p
// 	}
// 	pe.providersMu.Unlock()

// 	sortedNames := make([]string, 0, len(providers))
// 	for name := range providers {
// 		sortedNames = append(sortedNames, name)
// 	}
// 	sort.Strings(sortedNames)

// 	sortedProviders := make([]interfaces.Provider, 0, len(providers))
// 	for _, name := range sortedNames {
// 		sortedProviders = append(sortedProviders, providers[name])
// 	}
// 	return sortedProviders
// }

// func (pe *promptEngine) GetHistory() []interfaces.Result {
// 	return pe.history.snapshot()
// }

// func (pe *promptEngine) SaveToHistory(prompt, response string) error {
// 	pe.history.add(interfaces.Result{
// 		ID:        uuid.NewString(),
// 		Prompt:    prompt,
// 		Response:  response,
// 		Provider:  "manual",
// 		Timestamp: time.Now(),
// 	})
// 	return nil
// }

// // executeRaw allows provider adapters and API configs to run direct prompts without templating.
// func (pe *promptEngine) executeRaw(providerName, prompt string) (string, error) {
// 	res, err := pe.invokeProvider(context.Background(), providerName, prompt, map[string]interface{}{})
// 	if err != nil {
// 		return "", err
// 	}
// 	return res.Response, nil
// }

// func (pe *promptEngine) invokeProvider(ctx context.Context, providerName, prompt string, vars map[string]interface{}) (*interfaces.Result, error) {
// 	if err := pe.ensureBootstrap(); err != nil {
// 		return nil, err
// 	}

// 	provider := pe.registry.Resolve(providerName)
// 	if provider == nil {
// 		return nil, fmt.Errorf("provider '%s' not found", providerName)
// 	}

// 	if ok := provider.IsAvailable(); !ok {
// 		return nil, fmt.Errorf("provider '%s' is not available", providerName)
// 	}

// 	model := pe.cfg.defaultModels[providerName]
// 	if override, ok := vars["model"].(string); ok && override != "" {
// 		model = override
// 	}

// 	req := interfaces.ChatRequest{
// 		Provider: providerName,
// 		Model:    model,
// 		Messages: []interfaces.Message{{Role: "user", Content: prompt}},
// 		Temp:     pe.cfg.defaultTemperature,
// 		Stream:   false,
// 		Meta:     map[string]any{},
// 	}

// 	var responseContent bytes.Buffer
// 	var usage *interfaces.Usage

// 	operation := func() error {
// 		ctxWithTimeout, cancel := context.WithTimeout(ctx, pe.cfg.requestTimeout())
// 		defer cancel()

// 		stream, err := provider.Chat(ctxWithTimeout, req)
// 		if err != nil {
// 			return err
// 		}

// 		for chunk := range stream {
// 			if chunk.Error != "" {
// 				return errors.New(chunk.Error)
// 			}
// 			if chunk.Content != "" {
// 				responseContent.WriteString(chunk.Content)
// 			}
// 			if chunk.Done && chunk.Usage != nil {
// 				usage = chunk.Usage
// 			}
// 		}

// 		return nil
// 	}

// 	if pe.middleware != nil {
// 		if err := pe.middleware.WrapProvider(providerName, operation); err != nil {
// 			return nil, err
// 		}
// 	} else if err := operation(); err != nil {
// 		return nil, err
// 	}

// 	metadata := map[string]interface{}{}
// 	if usage != nil {
// 		metadata["usage"] = usage
// 	}

// 	result := &interfaces.Result{
// 		ID:        uuid.NewString(),
// 		Prompt:    prompt,
// 		Response:  responseContent.String(),
// 		Provider:  providerName,
// 		Variables: vars,
// 		Metadata:  metadata,
// 		Timestamp: time.Now(),
// 	}

// 	return result, nil
// }
