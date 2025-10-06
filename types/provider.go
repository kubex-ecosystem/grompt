package types

import "errors"

// Provider is the legacy provider interface exposed to consumers.
type Provider interface {
	Name() string
	Version() string
	Execute(prompt string) (string, error)
	IsAvailable() bool
	GetCapabilities() *Capabilities
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

func (p *providerAdapter) Execute(prompt string) (string, error) {
	if p.engine == nil {
		return "", errors.New("grompt provider adapter not bound to engine")
	}
	return p.engine.executeRaw(p.name, prompt)
}

func (p *providerAdapter) IsAvailable() bool {
	if p.engine == nil || p.engine.registry == nil {
		return false
	}
	prov := p.engine.registry.Resolve(p.name)
	if prov == nil {
		return false
	}
	return prov.Available() == nil
}

func (p *providerAdapter) GetCapabilities() *Capabilities {
	caps := defaultCapabilities(p.name, p.engine.cfg.defaultModels[p.name])
	if prov := p.engine.registry.Resolve(p.name); prov != nil {
		if prov.Available() != nil {
			return nil
		}
	}
	return caps
}
