// Package advise implements the advice generation for Repository Intelligence.
package advise

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/gateway/registry"
	"github.com/kubex-ecosystem/grompt/internal/interfaces"
)

type Handler struct{ reg *registry.Registry }

func New(reg *registry.Registry) *Handler { return &Handler{reg: reg} }

type adviseReq struct {
	Mode        string         `json:"mode"`
	Provider    string         `json:"provider"`
	Model       string         `json:"model"`
	Scorecard   map[string]any `json:"scorecard"`
	Hotspots    []string       `json:"hotspots"`
	Temperature float32        `json:"temperature"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	if mode == "" {
		http.Error(w, "mode required: exec|code|ops|community", http.StatusBadRequest)
		return
	}

	var in adviseReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p := h.reg.Resolve(in.Provider)
	if p == nil {
		http.Error(w, "bad provider", http.StatusBadRequest)
		return
	}

	sys := systemPrompt(in.Mode)
	user := userPrompt(in.Scorecard, in.Hotspots)

	headers := map[string]string{
		"x-external-api-key": r.Header.Get("x-external-api-key"),
		"x-tenant-id":        r.Header.Get("x-tenant-id"),
		"x-user-id":          r.Header.Get("x-user-id"),
	}

	ch, err := p.Chat(r.Context(), interfaces.ChatRequest{
		Provider: in.Provider,
		Model:    in.Model,
		Temp:     in.Temperature,
		Stream:   true,
		Messages: []interfaces.Message{
			{Role: "system", Content: sys},
			{Role: "user", Content: user},
		},
		Meta:    map[string]any{},
		Headers: headers,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	flusher, _ := w.(http.Flusher)

	enc := func(v any) []byte { b, _ := json.Marshal(v); return b }
	start := time.Now()
	for c := range ch {
		if c.Content != "" {
			w.Write([]byte("data: "))
			w.Write(enc(map[string]any{"content": c.Content}))
			w.Write([]byte("\n\n"))
			if flusher != nil {
				flusher.Flush()
			}
			if c.Done {
				w.Write([]byte("data: "))
				w.Write(enc(map[string]any{"done": true, "usage": c.Usage, "mode": mode}))
				w.Write([]byte("nn"))
				flusher.Flush()
			}
		}
		_ = start

	}
}

// systemPrompt retorna o prompt de sistema apropriado para o modo requerido.
// Inclui casos para exec|code|ops|community e um fallback genérico.
func systemPrompt(mode string) string {
	switch mode {
	case "exec":
		return "You are an expert in repository execution and operational guidance. Provide clear, concise instructions focusing on runtime behavior, command usage, and how to reproduce issues."
	case "code":
		return "You are a code reviewer and refactorer. Provide focused suggestions to improve code quality, architecture, maintainability, and testing. Highlight concrete code changes when appropriate."
	case "ops":
		return "You are an infrastructure and operations expert. Suggest improvements for deployment, CI/CD, monitoring, reliability, and security relevant to this repository."
	case "community":
		return "You are a community and contributor experience advisor. Recommend improvements for documentation, contributing guidelines, issue templates, and onboarding for new contributors."
	default:
		return "You are a general repository advisor. Provide concise, actionable recommendations to improve the repository across code, operations, and contributor experience."
	}
}

func userPrompt(scorecard map[string]any, hotspots []string) string {
	scorecardStr, _ := json.MarshalIndent(scorecard, "", "  ")
	hotspotsStr, _ := json.MarshalIndent(hotspots, "", "  ")

	return `Here are the scorecard results for a software repository:
` + "```json\n" + string(scorecardStr) + "\n```\n" + `
Here are some identified hotspots in the repository that may need attention:
` + "```json\n" + string(hotspotsStr) + "\n```\n" + `
Based on the above scorecard results and hotspots, please provide specific, actionable advice to improve the repository. Focus on practical steps that can be taken to address any issues or weaknesses identified. Be concise and prioritize the most impactful recommendations.`
}
