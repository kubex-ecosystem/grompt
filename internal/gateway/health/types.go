package health

import "time"

// Tier representa o nível de verificação de saúde
type Tier int

const (
	Tier1Key       Tier = iota + 1 // Validação de key (sem tokens)
	Tier2Handshake                 // Handshake mínimo (OPTIONS/HEAD)
	Tier3Real                      // Micro-request controlada (pode gastar migalha)
)

// Status representa o estado de saúde do provider
type Status string

const (
	StatusOK       Status = "ok"       // Funcionando perfeitamente
	StatusSuspect  Status = "suspect"  // Comportamento suspeito
	StatusDegraded Status = "degraded" // Funcionando mas com problemas
	StatusDown     Status = "down"     // Indisponível
)

// ProbeResult contém o resultado de uma verificação de saúde
type ProbeResult struct {
	Provider     string    `json:"provider"`
	Tier         Tier      `json:"tier"`
	Status       Status    `json:"status"`
	Details      string    `json:"details,omitempty"`
	CheckedAt    time.Time `json:"checked_at"`
	TTLSeconds   int       `json:"ttl_seconds"`
	LatencyMs    int64     `json:"latency_ms"`
	RateLimitRem *int      `json:"rate_limit_remaining,omitempty"`
	HTTPCode     int       `json:"http_code,omitempty"`
}

// Prober é a interface que cada provider deve implementar
type Prober interface {
	Name() string
	// Check executa o tier solicitado.
	// Deve ser idempotente e barato quando possível.
	Check(t Tier) (ProbeResult, error)
	// HealthyHint retorna um TTL sugerido por tier.
	HealthyHint(t Tier) time.Duration
}
