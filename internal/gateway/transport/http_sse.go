package transport

import (
	"encoding/json"
	"net/http"

	"github.com/kubex-ecosystem/grompt/internal/advise"
	"github.com/kubex-ecosystem/grompt/internal/gateway/registry"
	"github.com/kubex-ecosystem/grompt/internal/interfaces"
	"github.com/kubex-ecosystem/grompt/internal/scorecard"
)

type httpHandlersSSE struct {
	reg    *registry.Registry
	engine *scorecard.Engine // Add scorecard engine
}

// WireHTTP sets up HTTP routes
// func WireHTTP(mux *http.ServeMux, reg *registry.Registry) {
// 	h := &httpHandlers{registry: reg}

// 	mux.HandleFunc("/healthz", h.healthCheck)
// 	mux.HandleFunc("/v1/chat", h.chatSSE)
// 	mux.HandleFunc("/v1/providers", h.listProviders)
// }

func WireHTTPSSE(mux *http.ServeMux, reg *registry.Registry) {
	hh := &httpHandlersSSE{reg: reg, engine: nil} // TODO: Initialize engine when ready
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	mux.HandleFunc("/v1/chat", hh.chatSSE)
	mux.HandleFunc("/v1/session", hh.session)
	mux.HandleFunc("/v1/providers", hh.providers) // status simples
	mux.HandleFunc("/v1/auth/login", hh.authLoginPassthrough)
	mux.HandleFunc("/v1/state/export", hh.stateExport)
	mux.HandleFunc("/v1/state/import", hh.stateImport)
	mux.Handle("/v1/advise", advise.New(reg))

	// Repository Intelligence APIs (to be implemented)
	// mux.HandleFunc("/api/v1/scorecard", hh.handleScorecard)
	// mux.HandleFunc("/api/v1/scorecard/advice", hh.handleScorecardAdvice)
	// mux.HandleFunc("/api/v1/metrics/ai", hh.handleAIMetrics)
	// mux.HandleFunc("/api/v1/health", hh.handleHealthRI)
}

type chatReq struct {
	Provider string              `json:"provider"`
	Model    string              `json:"model"`
	Messages []interfaces.Message `json:"messages"`
	Temp     float32             `json:"temperature"`
	Stream   bool                `json:"stream"`
	Meta     map[string]any      `json:"meta"`
}

func (h *httpHandlersSSE) chatSSE(w http.ResponseWriter, r *http.Request) {
	var in chatReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p := h.reg.Resolve(in.Provider)
	if p == nil {
		http.Error(w, "bad provider", http.StatusBadRequest)
		return
	}
	headers := map[string]string{
		"x-external-api-key": r.Header.Get("x-external-api-key"),
		"x-tenant-id":        r.Header.Get("x-tenant-id"),
		"x-user-id":          r.Header.Get("x-user-id"),
	}
	ch, err := p.Chat(r.Context(), interfaces.ChatRequest{
		Provider: in.Provider,
		Model:    in.Model,
		Messages: in.Messages,
		Temp:     in.Temp,
		Stream:   in.Stream,
		Meta:     in.Meta,
		Headers:  headers,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	fl, _ := w.(http.Flusher)

	enc := func(v any) []byte { b, _ := json.Marshal(v); return b }
	for c := range ch {
		payload := map[string]any{}
		if c.Content != "" {
			payload["content"] = c.Content
		}
		// if c.ToolCall != nil {
		// 	payload["toolCall"] = c.ToolCall
		// }
		if c.Done {
			payload["done"] = true
			if c.Usage != nil {
				payload["usage"] = c.Usage
			}
		}

		if len(payload) == 0 {
			continue
		}
		w.Write([]byte("data: "))
		w.Write(enc(payload))
		w.Write([]byte("\n\n"))
		fl.Flush()
	}
}

// /v1/providers — lista nomes e tipos carregados (pra pintar “verde” no dropdown)
func (h *httpHandlersSSE) providers(w http.ResponseWriter, r *http.Request) {
	cfg := h.reg.Config() // adicione um getter simples no registry
	type item struct{ Name, Type string }
	out := []item{}
	for name, pc := range cfg.Providers {
		out = append(out, item{Name: name, Type: pc.Type()})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"providers": out})
}

func (h *httpHandlersSSE) session(w http.ResponseWriter, r *http.Request) {
	// Implementar lógica para gerenciar sessões
}

func (h *httpHandlersSSE) authLoginPassthrough(w http.ResponseWriter, r *http.Request) {
	// Implementar lógica para login via passthrough
}

func (h *httpHandlersSSE) stateExport(w http.ResponseWriter, r *http.Request) {
	// Implementar lógica para exportar estado
}

func (h *httpHandlersSSE) stateImport(w http.ResponseWriter, r *http.Request) {
	// Implementar lógica para importar estado
}
