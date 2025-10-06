// Package providers supplies legacy helper constructors for compatibility with the refactored engine.
package providers

import (
	eng "github.com/kubex-ecosystem/grompt/internal/engine"
	tp "github.com/kubex-ecosystem/grompt/internal/types"
	logz "github.com/kubex-ecosystem/logz"
)

type Provider = tp.Provider

type Capabilities = tp.Capabilities

type Pricing = tp.Pricing

// Initialize mirrors the legacy helper, wiring a prompt engine and returning the available providers.
func Initialize(
	port string,
	bindAddr string,
	openAIKey string,
	deepSeekKey string,
	ollamaEndpoint string,
	claudeKey string,
	geminiKey string,
	chatgptKey string,
	logger logz.Logger,
) []Provider {
	cfg := tp.NewConfig(
		bindAddr,
		port,
		openAIKey,
		deepSeekKey,
		ollamaEndpoint,
		claudeKey,
		geminiKey,
		chatgptKey,
		logger,
	)

	e := eng.NewEngine(cfg)
	providers := []Provider{}
	providers = append(providers, e.GetProviders()...)

	return providers
}

// NewProvider lazily ensures the named provider is available and returns its adapter.
func NewProvider(name, apiKey, version string, cfg *tp.Config) Provider {
	if cfg == nil {
		return nil
	}

	if apiKey != "" {
		_ = cfg.SetAPIKey(name, apiKey)
	}

	engine := eng.NewEngine(cfg)
	for _, provider := range engine.GetProviders() {
		if provider.Name() == name {
			return provider
		}
	}
	return nil
}
