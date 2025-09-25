package health

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// HealthHandler implementa endpoints HTTP para health checking
type HealthHandler struct {
	engine   *Engine
	registry *ProberRegistry
}

// NewHealthHandler cria um novo handler
func NewHealthHandler(engine *Engine, registry *ProberRegistry) *HealthHandler {
	return &HealthHandler{
		engine:   engine,
		registry: registry,
	}
}

// HealthResponse representa a resposta de health check
type HealthResponse struct {
	Provider  string    `json:"provider"`
	Status    string    `json:"status"`
	Tier      int       `json:"tier"`
	Details   string    `json:"details,omitempty"`
	HTTPCode  int       `json:"http_code,omitempty"`
	Latency   int64     `json:"latency_ms,omitempty"`
	TTL       int64     `json:"ttl_seconds,omitempty"`
	RateLimit *int      `json:"rate_limit_remaining,omitempty"`
	CheckedAt time.Time `json:"checked_at"`
}

// OverallHealthResponse representa status geral de múltiplos providers
type OverallHealthResponse struct {
	Status    string                    `json:"status"`
	Providers map[string]HealthResponse `json:"providers"`
	Summary   map[string]int            `json:"summary"` // count por status
	CheckedAt time.Time                 `json:"checked_at"`
}

// ServeHTTP implementa http.Handler para health checks
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse do path: /health, /health/{provider}, /health/{provider}/{tier}
	path := strings.TrimPrefix(r.URL.Path, "/health")
	path = strings.Trim(path, "/")

	if path == "" {
		// Health geral de todos os providers
		h.handleOverallHealth(w, r)
		return
	}

	parts := strings.Split(path, "/")
	provider := parts[0]

	// Verifica se o provider existe
	if _, exists := h.registry.Get(provider); !exists {
		http.Error(w, `{"error":"provider not found"}`, http.StatusNotFound)
		return
	}

	tier := Tier1Key // default
	if len(parts) > 1 {
		if t, err := strconv.Atoi(parts[1]); err == nil && t >= 1 && t <= 3 {
			tier = Tier(t)
		} else {
			http.Error(w, `{"error":"invalid tier, use 1, 2, or 3"}`, http.StatusBadRequest)
			return
		}
	}

	// Executa health check específico
	h.handleProviderHealth(w, r, provider, tier)
}

// handleOverallHealth retorna status de todos os providers
func (h *HealthHandler) handleOverallHealth(w http.ResponseWriter, r *http.Request) {
	providers := h.registry.List()
	results := make(map[string]HealthResponse)
	summary := map[string]int{
		"ok":       0,
		"suspect":  0,
		"degraded": 0,
		"down":     0,
	}

	// Query params para controle
	escalate := r.URL.Query().Get("escalate") == "true"
	force := r.URL.Query().Get("force") == "true" || r.URL.Query().Get("force") == "1"
	tierParam := r.URL.Query().Get("tier")

	tier := Tier1Key
	if tierParam != "" {
		if t, err := strconv.Atoi(tierParam); err == nil && t >= 1 && t <= 3 {
			tier = Tier(t)
		}
	}

	for _, providerName := range providers {
		var result ProbeResult
		var err error

		if escalate {
			result, err = h.engine.CheckWithEscalation(providerName, Tier3Real, force)
		} else {
			result, err = h.engine.Check(providerName, tier, force)
		}

		if err != nil {
			// Se deu erro, marca como down
			result = ProbeResult{
				Provider: providerName,
				Tier:     tier,
				Status:   StatusDown,
				Details:  "health check failed: " + err.Error(),
			}
		}

		response := h.probeResultToResponse(result)
		results[providerName] = response
		summary[strings.ToLower(string(result.Status))]++
	}

	// Status geral: ok se todos ok, senão degraded se algum ok, senão down
	overallStatus := "down"
	statusCode := http.StatusServiceUnavailable // 503

	if summary["ok"] == len(providers) {
		overallStatus = "ok"
		statusCode = http.StatusOK // 200
	} else if summary["ok"] > 0 {
		overallStatus = "degraded"
		statusCode = http.StatusPartialContent // 206
	}

	// Headers úteis para LB/observabilidade
	w.Header().Set("X-Health", overallStatus)
	w.Header().Set("X-Health-Counts", fmt.Sprintf("ok=%d,suspect=%d,degraded=%d,down=%d",
		summary["ok"], summary["suspect"], summary["degraded"], summary["down"]))

	// Retry-After com mínimo TTL dos estados ruins
	if minTTL := h.calculateMinBadTTL(results); minTTL > 0 {
		w.Header().Set("Retry-After", strconv.Itoa(minTTL))
	}

	overall := OverallHealthResponse{
		Status:    overallStatus,
		Providers: results,
		Summary:   summary,
		CheckedAt: time.Now(),
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(overall)
}

// handleProviderHealth retorna status de um provider específico
func (h *HealthHandler) handleProviderHealth(w http.ResponseWriter, r *http.Request, provider string, tier Tier) {
	escalate := r.URL.Query().Get("escalate") == "true"
	force := r.URL.Query().Get("force") == "true" || r.URL.Query().Get("force") == "1"

	var result ProbeResult
	var err error

	if escalate {
		result, err = h.engine.CheckWithEscalation(provider, Tier3Real, force)
	} else {
		result, err = h.engine.Check(provider, tier, force)
	}

	if err != nil {
		http.Error(w, `{"error":"health check failed: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := h.probeResultToResponse(result)

	// Headers úteis para single provider
	w.Header().Set("X-Health", strings.ToLower(string(result.Status)))
	if result.TTLSeconds > 0 && result.Status != StatusOK {
		w.Header().Set("Retry-After", strconv.Itoa(result.TTLSeconds))
	}

	// Status code baseado no resultado
	statusCode := http.StatusOK
	switch result.Status {
	case StatusSuspect, StatusDegraded:
		statusCode = http.StatusPartialContent // 206
	case StatusDown:
		statusCode = http.StatusServiceUnavailable // 503
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// probeResultToResponse converte ProbeResult para HealthResponse
func (h *HealthHandler) probeResultToResponse(result ProbeResult) HealthResponse {
	response := HealthResponse{
		Provider:  result.Provider,
		Status:    strings.ToLower(string(result.Status)),
		Tier:      int(result.Tier),
		Details:   result.Details,
		HTTPCode:  result.HTTPCode,
		CheckedAt: result.CheckedAt,
	}

	if result.LatencyMs > 0 {
		response.Latency = result.LatencyMs
	}

	if result.TTLSeconds > 0 {
		response.TTL = int64(result.TTLSeconds)
	}

	if result.RateLimitRem != nil {
		response.RateLimit = result.RateLimitRem
	}

	return response
}

// calculateMinBadTTL calcula o menor TTL dos providers com status ruim
func (h *HealthHandler) calculateMinBadTTL(results map[string]HealthResponse) int {
	minTTL := int64(999999)
	foundBad := false

	for _, result := range results {
		if result.Status != "ok" && result.TTL > 0 {
			if result.TTL < minTTL {
				minTTL = result.TTL
				foundBad = true
			}
		}
	}

	if !foundBad {
		return 0
	}
	return int(minTTL)
}

// RegisterRoutes registra as rotas de health check em um mux
func RegisterRoutes(mux *http.ServeMux, engine *Engine, registry *ProberRegistry) {
	handler := NewHealthHandler(engine, registry)
	mux.Handle("/health", handler)
	mux.Handle("/health/", handler)
}
