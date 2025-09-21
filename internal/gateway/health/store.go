package health

import (
	"sync"
	"time"
)

// memo armazena um resultado com expiração
type memo struct {
	res ProbeResult
	exp time.Time
}

// Store é um cache thread-safe para resultados de health check
type Store struct {
	mu   sync.RWMutex
	data map[string]memo // key: provider|tier
}

// NewStore cria uma nova instância do store
func NewStore() *Store {
	return &Store{
		data: make(map[string]memo),
	}
}

// key gera uma chave única para provider + tier
func key(provider string, tier Tier) string {
	return provider + "|" + string(rune('0'+int(tier)))
}

// Get recupera um resultado do cache se ainda válido
func (s *Store) Get(provider string, tier Tier) (ProbeResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	m, ok := s.data[key(provider, tier)]
	if !ok || time.Now().After(m.exp) {
		return ProbeResult{}, false
	}
	return m.res, true
}

// Set armazena um resultado no cache com TTL
func (s *Store) Set(provider string, tier Tier, res ProbeResult, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key(provider, tier)] = memo{
		res: res,
		exp: time.Now().Add(ttl),
	}
}

// Clear remove uma entrada específica do cache
func (s *Store) Clear(provider string, tier Tier) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key(provider, tier))
}

// ClearProvider remove todas as entradas de um provider
func (s *Store) ClearProvider(provider string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k := range s.data {
		if len(k) > len(provider) && k[:len(provider)] == provider && k[len(provider)] == '|' {
			delete(s.data, k)
		}
	}
}

// Stats retorna estatísticas do cache
func (s *Store) Stats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := len(s.data)
	expired := 0
	now := time.Now()

	for _, m := range s.data {
		if now.After(m.exp) {
			expired++
		}
	}

	return map[string]interface{}{
		"total_entries": total,
		"expired":       expired,
		"active":        total - expired,
	}
}
