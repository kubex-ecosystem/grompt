package types

import (
	"context"
	"errors"

	"github.com/kubex-ecosystem/grompt/internal/interfaces"
)

// ---------- Provider adapter ----------

type providerAdapter struct {
	name    string
	version string
	type_   string
	keyEnv  string
	legacy  interfaces.LegacyProvider
	engine  interfaces.IEngine
}

func (p *providerAdapter) Name() string {
	return p.name
}

func (p *providerAdapter) Type() string {
	if p.type_ == "" {
		return "text"
	}
	return p.type_
}

func (p *providerAdapter) KeyEnv() string {
	return p.keyEnv
}

func (p *providerAdapter) Version() string {
	if p.version == "" {
		return "v1"
	}
	return p.version
}

func (p *providerAdapter) Execute(ctx context.Context, template string, vars map[string]any) (*interfaces.Result, error) {
	if p.engine == nil {
		return nil, errors.New("grompt provider adapter not bound to engine")
	}
	return p.engine.ProcessPrompt(ctx, template, vars)
}

func (p *providerAdapter) IsAvailable() bool {
	if p.engine == nil || p.engine.GetRegistry() == nil {
		return false
	}
	prov := p.engine.Resolve(p.name)
	if prov == nil {
		return false
	}
	return prov.IsAvailable()
}

func (p *providerAdapter) GetCapabilities(ctx context.Context) *interfaces.Capabilities {
	models, _ := p.engine.GetConfig().GetAPIConfig(p.name).ListModels()
	caps := defaultCapabilities(p.name, models[p.name].(string))
	if prov := p.engine.Resolve(p.name); prov != nil {
		if !prov.IsAvailable() {
			return nil
		}
	}
	return caps
}

func (p *providerAdapter) Chat(ctx context.Context, req interfaces.ChatRequest) (<-chan interfaces.ChatChunk, error) {
	if p.legacy == nil {
		return nil, errors.New("legacy provider not set")
	}
	return p.legacy.Chat(ctx, req)
}

func (p *providerAdapter) Notify(ctx context.Context, event interfaces.NotificationEvent) error {
	if p.legacy == nil {
		return errors.New("legacy provider not set")
	}
	return p.legacy.Notify(ctx, event)
}

// NewLegacyProviderAdapter creates a new Provider adapter from a LegacyProvider
func NewLegacyProviderAdapter(name, version, type_, keyEnv string, legacy interfaces.LegacyProvider, engine interfaces.IEngine) interfaces.Provider {
	return &providerAdapter{
		name:    name,
		version: version,
		type_:   type_,
		keyEnv:  keyEnv,
		legacy:  legacy,
		engine:  engine,
	}
}
