package health

import (
	"errors"
	"math/rand"
	"time"
)

// Engine orquestra os health checks com escalação automática
type Engine struct {
	store   *Store
	probers map[string]Prober
}

// NewEngine cria uma nova instância do engine
func NewEngine(store *Store, probers ...Prober) *Engine {
	m := make(map[string]Prober)
	for _, p := range probers {
		m[p.Name()] = p
	}
	return &Engine{
		store:   store,
		probers: m,
	}
}

// Providers retorna lista de providers disponíveis
func (e *Engine) Providers() []string {
	out := make([]string, 0, len(e.probers))
	for k := range e.probers {
		out = append(out, k)
	}
	return out
}

// Check executa health check com escalação automática
func (e *Engine) Check(provider string, tier Tier, force bool) (ProbeResult, error) {
	p, ok := e.probers[provider]
	if !ok {
		return ProbeResult{}, errors.New("provider desconhecido: " + provider)
	}

	// Verifica cache primeiro (se não for force)
	if !force {
		if res, ok := e.store.Get(provider, tier); ok {
			return res, nil
		}
	}

	// Executa o check
	start := time.Now()
	res, err := p.Check(tier)
	res.LatencyMs = time.Since(start).Milliseconds()
	res.CheckedAt = time.Now()
	res.Provider = provider
	res.Tier = tier

	if err != nil {
		if res.Status == "" {
			res.Status = StatusDown
		}
		res.Details = err.Error()
	}

	// TTL sugerido + jitter leve (para espalhar carga)
	ttl := p.HealthyHint(tier)
	if ttl > 0 {
		// Adiciona jitter de ±400ms para TTLs > 1s
		if ttl > time.Second {
			jitter := time.Duration(rand.Intn(800)-400) * time.Millisecond
			ttl = ttl + jitter
		}
		res.TTLSeconds = int(ttl.Seconds())
		e.store.Set(provider, tier, res, ttl)
	}

	return res, nil
}

// CheckWithEscalation executa health check com escalação automática entre tiers
func (e *Engine) CheckWithEscalation(provider string, maxTier Tier, force bool) (ProbeResult, error) {
	// Começa pelo Tier 1 e escala conforme necessário
	for tier := Tier1Key; tier <= maxTier; tier++ {
		res, err := e.Check(provider, tier, force)

		// Se OK ou erro fatal, retorna
		if res.Status == StatusOK || err != nil {
			return res, err
		}

		// Se suspect/degraded/down e ainda há tiers, continua
		if res.Status != StatusOK && tier < maxTier {
			continue
		}

		// Último tier testado, retorna resultado
		return res, err
	}

	return ProbeResult{}, errors.New("nenhum tier disponível")
}

// GetProviderStatus retorna o status geral de um provider baseado em todos os tiers
func (e *Engine) GetProviderStatus(provider string) (*ProviderStatus, error) {
	if _, ok := e.probers[provider]; !ok {
		return nil, errors.New("provider desconhecido: " + provider)
	}

	status := &ProviderStatus{
		Provider:  provider,
		Overall:   StatusDown,
		Tiers:     make(map[Tier]ProbeResult),
		UpdatedAt: time.Now(),
	}

	// Verifica todos os tiers no cache
	bestStatus := StatusDown
	for tier := Tier1Key; tier <= Tier3Real; tier++ {
		if res, ok := e.store.Get(provider, tier); ok {
			status.Tiers[tier] = res

			// Define status geral baseado no melhor tier
			if res.Status == StatusOK {
				bestStatus = StatusOK
			} else if res.Status == StatusDegraded && bestStatus != StatusOK {
				bestStatus = StatusDegraded
			} else if res.Status == StatusSuspect && bestStatus == StatusDown {
				bestStatus = StatusSuspect
			}
		}
	}

	status.Overall = bestStatus
	return status, nil
}

// ClearCache limpa o cache de um provider específico
func (e *Engine) ClearCache(provider string) {
	e.store.ClearProvider(provider)
}

// Stats retorna estatísticas do engine
func (e *Engine) Stats() map[string]interface{} {
	stats := e.store.Stats()
	stats["providers_count"] = len(e.probers)
	stats["providers"] = e.Providers()
	return stats
}

// ProviderStatus representa o status completo de um provider
type ProviderStatus struct {
	Provider  string               `json:"provider"`
	Overall   Status               `json:"overall"`
	Tiers     map[Tier]ProbeResult `json:"tiers"`
	UpdatedAt time.Time            `json:"updated_at"`
}
