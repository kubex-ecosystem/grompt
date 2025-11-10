// Package providers supplies legacy helper constructors for compatibility with the refactored engine.
package providers

import (
	"time"

	eng "github.com/kubex-ecosystem/grompt/internal/engine"
	i "github.com/kubex-ecosystem/grompt/internal/interfaces"
	tp "github.com/kubex-ecosystem/grompt/internal/types"
	logz "github.com/kubex-ecosystem/logz/logger"
)

type Provider = i.Provider

type Capabilities = i.Capabilities

type Pricing = i.Pricing

// Initialize mirrors the legacy helper, wiring a prompt engine and returning the available providers.
func Initialize(
	name               string,
	debug              bool,
	logger             logz.Logger,
	bindAddr           string,
	port               string,
	tempDir            string,
	logFile           string,
	envFile           string,
	configFile        string,
	cwd               string,
	openAIKey         string,
	claudeKey        string,
	geminiKey        string,
	deepSeekKey      string,
	chatGPTKey      string,
	ollamaEndpoint   string,
	apiKeys            map[string]string,
	endpoints          map[string]string,
	defaultModels      map[string]string,
	providerTypes      map[string]string,
	defaultProvider    string,
	defaultTemperature float32,
	historyLimit       int,
	timeout            time.Duration,
	providerConfigPath string,
) []Provider {
	cfg := tp.NewConfig(
		name,
		debug,
		logger,
		bindAddr,
		port,
		tempDir,
		logFile,
		envFile,
		configFile,
		cwd,
		openAIKey,
		claudeKey,
		geminiKey,
		deepSeekKey,
		chatGPTKey,
		ollamaEndpoint,
		apiKeys,
		endpoints,
		defaultModels,
		providerTypes,
		defaultProvider,
		defaultTemperature,
		historyLimit,
		timeout,
		providerConfigPath,
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
