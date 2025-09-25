package health

import (
	"context"
	"log"
	"sync"
	"time"
)

// Scheduler gerencia health checks automáticos em background
type Scheduler struct {
	engine    *Engine
	registry  *ProberRegistry
	intervals map[Tier]time.Duration
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	running   bool
	mu        sync.RWMutex
}

// SchedulerConfig configuração do scheduler
type SchedulerConfig struct {
	// Intervalos por tier (quanto tempo entre checks automáticos)
	Tier1Interval time.Duration // Default: 15min (key validation)
	Tier2Interval time.Duration // Default: 5min (handshake check)
	Tier3Interval time.Duration // Default: 30min (real request - mais caro)

	// Configurações gerais
	EnableStaggering bool          // Espalha checks no tempo para evitar picos
	StaggerWindow    time.Duration // Janela para espalhar (default: 2min)
	LogVerbose       bool          // Log detalhado
}

// DefaultSchedulerConfig retorna configuração padrão do scheduler
func DefaultSchedulerConfig() SchedulerConfig {
	return SchedulerConfig{
		Tier1Interval:    15 * time.Minute,
		Tier2Interval:    5 * time.Minute,
		Tier3Interval:    30 * time.Minute,
		EnableStaggering: true,
		StaggerWindow:    2 * time.Minute,
		LogVerbose:       false,
	}
}

// NewScheduler cria novo scheduler
func NewScheduler(engine *Engine, registry *ProberRegistry, config SchedulerConfig) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	intervals := map[Tier]time.Duration{
		Tier1Key:       config.Tier1Interval,
		Tier2Handshake: config.Tier2Interval,
		Tier3Real:      config.Tier3Interval,
	}

	return &Scheduler{
		engine:    engine,
		registry:  registry,
		intervals: intervals,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Start inicia o scheduler
func (s *Scheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return nil // Já está rodando
	}

	s.running = true

	// Inicia goroutines para cada tier
	for tier, interval := range s.intervals {
		s.wg.Add(1)
		go s.runTierScheduler(tier, interval)
	}

	log.Printf("[HealthScheduler] Started with intervals: T1=%v, T2=%v, T3=%v",
		s.intervals[Tier1Key], s.intervals[Tier2Handshake], s.intervals[Tier3Real])

	return nil
}

// Stop para o scheduler
func (s *Scheduler) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil // Já parado
	}

	s.cancel()
	s.wg.Wait()
	s.running = false

	log.Printf("[HealthScheduler] Stopped")
	return nil
}

// IsRunning verifica se o scheduler está ativo
func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// runTierScheduler executa checks periódicos para um tier específico
func (s *Scheduler) runTierScheduler(tier Tier, interval time.Duration) {
	defer s.wg.Done()

	log.Printf("[HealthScheduler] Tier %d scheduler started (interval: %v)", int(tier), interval)

	// Primeira execução com stagger para espalhar carga
	stagger := s.calculateStagger(tier)

	select {
	case <-time.After(stagger):
		s.runChecksForTier(tier)
	case <-s.ctx.Done():
		return
	}

	// Loop principal
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.runChecksForTier(tier)
		case <-s.ctx.Done():
			log.Printf("[HealthScheduler] Tier %d scheduler stopped", int(tier))
			return
		}
	}
}

// runChecksForTier executa health checks para todos os providers em um tier
func (s *Scheduler) runChecksForTier(tier Tier) {
	providers := s.registry.List()

	log.Printf("[HealthScheduler] Running Tier %d checks for %d providers", int(tier), len(providers))

	for _, provider := range providers {
		// Executa check em background para não bloquear scheduler
		go func(providerName string) {
			result, err := s.engine.Check(providerName, tier, false) // usa cache se disponível

			if err != nil {
				log.Printf("[HealthScheduler] Tier %d check failed for %s: %v", int(tier), providerName, err)
				return
			}

			// Log apenas para problemas ou verbose mode
			if result.Status != StatusOK {
				log.Printf("[HealthScheduler] Tier %d %s: %s - %s",
					int(tier), providerName, result.Status, result.Details)
			}
		}(provider)
	}
}

// calculateStagger calcula delay inicial para espalhar checks
func (s *Scheduler) calculateStagger(tier Tier) time.Duration {
	// Espalha checks baseado no tier para evitar todos executarem simultaneamente
	tierOffset := time.Duration(int(tier)) * 30 * time.Second
	return tierOffset
}

// GetStats retorna estatísticas do scheduler
func (s *Scheduler) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"running":     s.running,
		"providers":   len(s.registry.List()),
		"tier1_every": s.intervals[Tier1Key].String(),
		"tier2_every": s.intervals[Tier2Handshake].String(),
		"tier3_every": s.intervals[Tier3Real].String(),
	}
}

// ForceCheck força execução imediata de todos os tiers para todos os providers
func (s *Scheduler) ForceCheck() {
	providers := s.registry.List()
	log.Printf("[HealthScheduler] Force check for %d providers", len(providers))

	for _, provider := range providers {
		for tier := Tier1Key; tier <= Tier3Real; tier++ {
			go func(providerName string, t Tier) {
				result, err := s.engine.Check(providerName, t, true) // força novo check
				if err != nil {
					log.Printf("[HealthScheduler] Force check T%d %s failed: %v", int(t), providerName, err)
				} else {
					log.Printf("[HealthScheduler] Force check T%d %s: %s", int(t), providerName, result.Status)
				}
			}(provider, tier)
		}
	}
}
