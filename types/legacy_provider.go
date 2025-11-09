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
	engine  *promptEngine
}

func (p *providerAdapter) Name() string {
	return p.name
}

func (p *providerAdapter) Version() string {
	if p.version == "" {
		return "v1"
	}
	return p.version
}

func (p *providerAdapter) Execute(ctx context.Context, prompt string) (string, error) {
	if p.engine == nil {
		return "", errors.New("grompt provider adapter not bound to engine")
	}
	return p.engine.executeRaw(p.name, prompt)
}

func (p *providerAdapter) Available() bool {
	if p.engine == nil || p.engine.registry == nil {
		return false
	}
	prov := p.engine.registry.Resolve(p.name)
	if prov == nil {
		return false
	}
	return prov.IsAvailable()
}

func (p *providerAdapter) GetCapabilities(ctx context.Context) *interfaces.Capabilities {
	caps := defaultCapabilities(p.name, p.engine.cfg.defaultModels[p.name])
	if prov := p.engine.registry.Resolve(p.name); prov != nil {
		if !prov.IsAvailable() {
			return nil
		}
	}
	return caps
}

func (p *providerAdapter) Chat(ctx context.Context, req interfaces.ChatRequest) (<-chan interfaces.ChatChunk, error) {
	return nil, errors.New("not implemented")
}

func (p *providerAdapter) Complete(ctx context.Context, prompt string, options map[string]interface{}) (*interfaces.Result, error) {
	return &interfaces.Result{}, errors.New("not implemented")
}
