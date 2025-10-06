package types

import (
	"context"
	"errors"
)

// Provider is the legacy provider interface exposed to consumers.
type Provider interface {
	Name() string
	Version() string
	Execute(ctx context.Context, prompt string) (string, error)
	Available() bool
	GetCapabilities(ctx context.Context) *Capabilities
}

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
	return prov.Available()
}

func (p *providerAdapter) GetCapabilities(ctx context.Context) *Capabilities {
	caps := defaultCapabilities(p.name, p.engine.cfg.defaultModels[p.name])
	if prov := p.engine.registry.Resolve(p.name); prov != nil {
		if !prov.Available() {
			return nil
		}
	}
	return caps
}
