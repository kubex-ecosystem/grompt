package health

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"
)

// GeminiProber implementa health checking para Google Gemini API
type GeminiProber struct {
	BaseURL string
	Client  *http.Client
}

// NewGeminiProber cria uma nova instância do prober Gemini
func NewGeminiProber() *GeminiProber {
	return &GeminiProber{
		BaseURL: "https://generativelanguage.googleapis.com/v1beta",
		Client: &http.Client{
			Timeout: 8 * time.Second, // Gemini pode ser um pouco mais lento
		},
	}
}

// Name retorna o nome do provider
func (p *GeminiProber) Name() string {
	return "gemini"
}

// HealthyHint retorna TTL sugerido por tier
func (p *GeminiProber) HealthyHint(t Tier) time.Duration {
	switch t {
	case Tier1Key:
		return 30 * time.Minute // Key validation: cache por 30min
	case Tier2Handshake:
		return 15 * time.Minute // Handshake: cache por 15min (Gemini mais lento)
	case Tier3Real:
		return 1 * time.Hour // Real request: cache por 1h
	default:
		return 15 * time.Minute
	}
}

// Check executa health check baseado no tier
func (p *GeminiProber) Check(t Tier) (ProbeResult, error) {
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
// Lista models com a API key (GET /models barato)
func (p *GeminiProber) tier1() (ProbeResult, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return ProbeResult{
			Provider: p.Name(),
			Tier:     Tier1Key,
			Status:   StatusDown,
			Details:  "GEMINI_API_KEY ausente",
		}, nil
	}

	// GET /models com API key no query param
	url := p.BaseURL + "/models?key=" + apiKey
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ProbeResult{
			Provider: p.Name(),
			Tier:     Tier1Key,
			Status:   StatusDown,
			Details:  "falha ao criar request: " + err.Error(),
		}, err
	}

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
	details := "API key válida"

	switch resp.StatusCode {
	case 200:
		status = StatusOK
		details = "API key válida, models disponíveis"
	case 400:
		status = StatusDown
		details = "API key inválida ou mal formada"
	case 403:
		status = StatusDown
		details = "API key sem permissão"
	case 429:
		status = StatusDegraded
		details = "rate limit atingido"
	case 500, 502, 503, 504:
		status = StatusDegraded
		details = "servidor com problemas"
	default:
		status = StatusSuspect
		details = "resposta inesperada"
	}

	return ProbeResult{
		Provider: p.Name(),
		Tier:     Tier1Key,
		Status:   status,
		Details:  details,
		HTTPCode: resp.StatusCode,
	}, nil
}

// tier2 - Handshake mínimo (sem tokens)
// HEAD ou OPTIONS em generateContent endpoint
func (p *GeminiProber) tier2() (ProbeResult, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return ProbeResult{
			Provider: p.Name(),
			Tier:     Tier2Handshake,
			Status:   StatusDown,
			Details:  "GEMINI_API_KEY ausente",
		}, nil
	}

	// Tenta HEAD no endpoint generateContent (sem body)
	url := p.BaseURL + "/models/gemini-1.5-flash:generateContent?key=" + apiKey
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return ProbeResult{
			Provider: p.Name(),
			Tier:     Tier2Handshake,
			Status:   StatusDown,
			Details:  "falha ao criar request: " + err.Error(),
		}, err
	}

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
	details := "endpoint responsivo"

	// Para Gemini, 400 em HEAD pode ser normal (need body), mas >= 500 é problema
	switch {
	case resp.StatusCode >= 500:
		status = StatusDegraded
		details = "servidor instável"
	case resp.StatusCode == 429:
		status = StatusDegraded
		details = "rate limit atingido"
	case resp.StatusCode == 403:
		status = StatusDown
		details = "API key sem permissão"
	case resp.StatusCode == 200 || resp.StatusCode == 400:
		// 200: ok, 400: normal p/ HEAD sem body no Gemini
		status = StatusOK
		details = "endpoint responsivo"
	default:
		status = StatusSuspect
		details = "resposta inesperada"
	}

	return ProbeResult{
		Provider: p.Name(),
		Tier:     Tier2Handshake,
		Status:   status,
		Details:  details,
		HTTPCode: resp.StatusCode,
	}, nil
}

// tier3 - Micro-request controlada (pode gastar migalha)
// generateContent com texto mínimo
func (p *GeminiProber) tier3() (ProbeResult, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return ProbeResult{
			Provider: p.Name(),
			Tier:     Tier3Real,
			Status:   StatusDown,
			Details:  "GEMINI_API_KEY ausente",
		}, nil
	}

	// Micro-request para Gemini: minimal prompt, sem stream
	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": "hi"},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"maxOutputTokens": 1,
			"temperature":     0.0,
		},
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

	url := p.BaseURL + "/models/gemini-1.5-flash:generateContent?key=" + apiKey
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return ProbeResult{
			Provider: p.Name(),
			Tier:     Tier3Real,
			Status:   StatusDown,
			Details:  "falha ao criar request: " + err.Error(),
		}, err
	}

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
	case resp.StatusCode == 400:
		// Precisamos ler o body pra saber se é problema de quota/key
		var respBody map[string]interface{}
		if json.NewDecoder(resp.Body).Decode(&respBody) == nil {
			if errData, ok := respBody["error"].(map[string]interface{}); ok {
				if msg, ok := errData["message"].(string); ok {
					if strings.Contains(strings.ToLower(msg), "quota") {
						status = StatusDegraded
						details = "quota esgotada"
					} else if strings.Contains(strings.ToLower(msg), "key") {
						status = StatusDown
						details = "API key inválida"
					} else {
						status = StatusSuspect
						details = "bad request: " + msg
					}
				}
			}
		} else {
			status = StatusSuspect
			details = "bad request (não conseguiu ler detalhes)"
		}
	case resp.StatusCode == 403:
		status = StatusDown
		details = "API key sem permissão"
	case resp.StatusCode != 200:
		status = StatusSuspect
		details = "resposta inesperada para generateContent"
	}

	return ProbeResult{
		Provider: p.Name(),
		Tier:     Tier3Real,
		Status:   status,
		Details:  details,
		HTTPCode: resp.StatusCode,
	}, nil
}
