// Package providers supplies legacy helper constructors for compatibility with the refactored engine.
package providers

import (
    "github.com/kubex-ecosystem/grompt"
    logz "github.com/kubex-ecosystem/logz"
)

type Provider = grompt.Provider

type Capabilities = grompt.Capabilities

type Pricing = grompt.Pricing

// Initialize mirrors the legacy helper, wiring a prompt engine and returning the available providers.
func Initialize(
    bindAddr string,
    port string,
    openAIKey string,
    deepSeekKey string,
    ollamaEndpoint string,
    claudeKey string,
    geminiKey string,
    chatgptKey string,
    logger logz.Logger,
) []Provider {
    cfg := grompt.NewConfig(
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

    engine := grompt.NewPromptEngine(cfg)
    return engine.GetProviders()
}

// NewProvider lazily ensures the named provider is available and returns its adapter.
func NewProvider(name, apiKey, version string, cfg grompt.Config) Provider {
    if cfg == nil {
        return nil
    }

    if apiKey != "" {
        _ = cfg.SetAPIKey(name, apiKey)
    }

    engine := grompt.NewPromptEngine(cfg)
    for _, provider := range engine.GetProviders() {
        if provider.Name() == name {
            return provider
        }
    }
    return nil
}
