package health

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"
)

// GroqProber implementa health checking para Groq API
type GroqProber struct {
	BaseURL string
	Client  *http.Client
}

// NewGroqProber cria uma nova instância do prober Groq
func NewGroqProber() *GroqProber {
	return &GroqProber{
		BaseURL: "https://api.groq.com/openai/v1",
		Client: &http.Client{
			Timeout: 6 * time.Second,
		},
	}
}

// Name retorna o nome do provider
func (p *GroqProber) Name() string {
	return "groq"
}

// HealthyHint retorna TTL sugerido por tier
func (p *GroqProber) HealthyHint(t Tier) time.Duration {
	switch t {
	case Tier1Key:
		return 30 * time.Minute // Key validation: cache por 30min
	case Tier2Handshake:
		return 10 * time.Minute // Handshake: cache por 10min
	case Tier3Real:
		return 1 * time.Hour // Real request: cache por 1h
	default:
		return 10 * time.Minute
	}
}

// Check executa health check baseado no tier
func (p *GroqProber) Check(t Tier) (ProbeResult, error) {
	switch t {
	case Tier1Key:
		return p.tier1()
	case Tier2Handshake:
		return p.tier2()
	case Tier3Real:
		return p.tier3()
	default:
		return ProbeResult{
			Provider: p.Name(),
			Tier:     t,
			Status:   StatusDown,
		}, errors.New("tier inválido")
	}
}

// tier1 - Validação de chave (sem tokens)
// GET /models?limit=1 é mais confiável que HEAD (evita falsos positivos)
func (p *GroqProber) tier1() (ProbeResult, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return ProbeResult{
			Provider: p.Name(),
			Tier:     Tier1Key,
			Status:   StatusDown,
			Details:  "GROQ_API_KEY ausente",
		}, nil
	}

	// GET em /models?limit=1 não gasta tokens (lista de modelos) e quase sempre responde 200
	req, err := http.NewRequest("GET", p.BaseURL+"/models?limit=1", nil)
	if err != nil {
		return ProbeResult{
			Provider: p.Name(),
			Tier:     Tier1Key,
			Status:   StatusDown,
			Details:  "falha ao criar request: " + err.Error(),
		}, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := p.Client.Do(req)
	if err != nil {
		return ProbeResult{
			Provider: p.Name(),
			Tier:     Tier1Key,
			Status:   StatusSuspect,
			Details:  "falha na conexão: " + err.Error(),
		}, err
	}
	defer resp.Body.Close()

	status := StatusOK
	details := "models ok"

	switch resp.StatusCode {
	case 401, 403:
		status = StatusDown
		details = "chave inválida ou sem permissão"
	case 404, 405:
		status = StatusSuspect
		details = "rota/método inesperado (proxy/CDN?)"
	case 429:
		status = StatusDegraded
		details = "rate limit atingido"
	default:
		if resp.StatusCode >= 500 {
			status = StatusDegraded
			details = "upstream 5xx"
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			status = StatusOK
			details = "models ok"
		} else {
			status = StatusSuspect
			details = "resposta inesperada"
		}
	}

	// Captura rate limit info se disponível
	var rateLimitRem *int
	if remain := resp.Header.Get("X-RateLimit-Remaining"); remain != "" {
		if val, err := strconv.Atoi(remain); err == nil {
			rateLimitRem = &val
		}
	}

	return ProbeResult{
		Provider:     p.Name(),
		Tier:         Tier1Key,
		Status:       status,
		Details:      details,
		HTTPCode:     resp.StatusCode,
		RateLimitRem: rateLimitRem,
		CheckedAt:    time.Now(),
		TTLSeconds:   1800,
	}, nil
}

// tier2 - Handshake mínimo (sem tokens)
// OPTIONS com tolerância a 405 (rota existe, método bloqueado)
func (p *GroqProber) tier2() (ProbeResult, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return ProbeResult{
			Provider: p.Name(),
			Tier:     Tier2Handshake,
			Status:   StatusDown,
			Details:  "GROQ_API_KEY ausente",
		}, nil
	}

	req, err := http.NewRequest("OPTIONS", p.BaseURL+"/chat/completions", nil)
	if err != nil {
		return ProbeResult{
			Provider: p.Name(),
			Tier:     Tier2Handshake,
			Status:   StatusDown,
			Details:  "falha ao criar request: " + err.Error(),
		}, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := p.Client.Do(req)
	if err != nil {
		return ProbeResult{
			Provider: p.Name(),
			Tier:     Tier2Handshake,
			Status:   StatusSuspect,
			Details:  "falha na conexão: " + err.Error(),
		}, err
	}
	defer resp.Body.Close()

	status := StatusOK
	details := "handshake ok"

	switch resp.StatusCode {
	case 200, 204:
		status = StatusOK
		details = "handshake ok"
	case 405:
		// Rota existe, método bloqueado -> ok
		status = StatusOK
		details = "rota existe (405)"
	case 404:
		status = StatusSuspect
		details = "rota não encontrada (proxy/CDN?)"
	case 429:
		status = StatusDegraded
		details = "rate limit atingido"
	default:
		if resp.StatusCode >= 500 {
			status = StatusDegraded
			details = "upstream 5xx"
		} else {
			status = StatusSuspect
			details = "resposta inesperada"
		}
	}

	return ProbeResult{
		Provider:   p.Name(),
		Tier:       Tier2Handshake,
		Status:     status,
		Details:    details,
		HTTPCode:   resp.StatusCode,
		CheckedAt:  time.Now(),
		TTLSeconds: 600,
	}, nil
}

// tier3 - Micro-request controlada (pode gastar migalha)
// max_tokens ultra baixo + prompt curtíssimo; não stream; cache 1h
func (p *GroqProber) tier3() (ProbeResult, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return ProbeResult{
			Provider: p.Name(),
			Tier:     Tier3Real,
			Status:   StatusDown,
			Details:  "GROQ_API_KEY ausente",
		}, nil
	}

	// Micro-request: minimal payload, single token
	body := map[string]interface{}{
		"model":       "llama-3.1-8b-instant",
		"messages":    []map[string]string{{"role": "user", "content": "ping"}},
		"temperature": 0.0,
		"max_tokens":  1,
		"stream":      false,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return ProbeResult{
			Provider: p.Name(),
			Tier:     Tier3Real,
			Status:   StatusDown,
			Details:  "falha ao marshalar request: " + err.Error(),
		}, err
	}

	req, err := http.NewRequest("POST", p.BaseURL+"/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return ProbeResult{
			Provider: p.Name(),
			Tier:     Tier3Real,
			Status:   StatusDown,
			Details:  "falha ao criar request: " + err.Error(),
		}, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.Client.Do(req)
	if err != nil {
		return ProbeResult{
			Provider: p.Name(),
			Tier:     Tier3Real,
			Status:   StatusSuspect,
			Details:  "falha na conexão: " + err.Error(),
		}, err
	}
	defer resp.Body.Close()

	status := StatusOK
	details := "API completamente funcional"

	switch {
	case resp.StatusCode >= 500:
		status = StatusDegraded
		details = "servidor instável"
	case resp.StatusCode == 429:
		status = StatusDegraded
		details = "rate limit atingido"
	case resp.StatusCode == 401 || resp.StatusCode == 403:
		status = StatusDown
		details = "API key inválida"
	case resp.StatusCode != 200:
		status = StatusSuspect
		details = "resposta inesperada para completion"
	}

	return ProbeResult{
		Provider: p.Name(),
		Tier:     Tier3Real,
		Status:   status,
		Details:  details,
		HTTPCode: resp.StatusCode,
	}, nil
}
