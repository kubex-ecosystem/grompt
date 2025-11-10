// Package transport implements HTTP transport handlers for server-sent events (SSE) in the Grompt application.
package transport

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kubex-ecosystem/grompt/internal/advise"
	"github.com/kubex-ecosystem/grompt/internal/gateway/registry"
	"github.com/kubex-ecosystem/grompt/internal/interfaces"
	"github.com/kubex-ecosystem/grompt/internal/scorecard"
)

type httpHandlersSSE struct {
	reg    *registry.Registry
	engine *scorecard.Engine // Add scorecard engine
}

func WireHTTPSSE(router gin.IRouter, reg *registry.Registry) {
	hh := &httpHandlersSSE{reg: reg, engine: nil} // TODO: Initialize engine when ready
	router.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	v1 := router.Group("/v1")
	v1.Any("/chat", hh.chatSSE)
	v1.Any("/session", hh.session)
	v1.Any("/providers", hh.providers) // status simples
	v1.Any("/auth/login", hh.authLoginPassthrough)
	v1.Any("/state/export", hh.stateExport)
	v1.Any("/state/import", hh.stateImport)
	v1.Any("/advise", gin.WrapH(advise.New(reg)))

	// Repository Intelligence APIs (to be implemented)
	// v1.Any("/scorecard", hh.handleScorecard)
	// v1.Any("/scorecard/advice", hh.handleScorecardAdvice)
	// v1.Any("/metrics/ai", hh.handleAIMetrics)
	// v1.Any("/health", hh.handleHealthRI)
}

type chatReq struct {
	Provider string               `json:"provider"`
	Model    string               `json:"model"`
	Messages []interfaces.Message `json:"messages"`
	Temp     float32              `json:"temperature"`
	Stream   bool                 `json:"stream"`
	Meta     map[string]any       `json:"meta"`
}

func (h *httpHandlersSSE) chatSSE(c *gin.Context) {
	var in chatReq
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := h.reg.Resolve(in.Provider)
	if p == nil {
		c.String(http.StatusBadRequest, "bad provider")
		return
	}
	headers := map[string]string{
		"x-external-api-key": c.GetHeader("x-external-api-key"),
		"x-tenant-id":        c.GetHeader("x-tenant-id"),
		"x-user-id":          c.GetHeader("x-user-id"),
	}
	ch, err := p.Chat(c.Request.Context(), interfaces.ChatRequest{
		Provider: in.Provider,
		Model:    in.Model,
		Messages: in.Messages,
		Temp:     in.Temp,
		Stream:   in.Stream,
		Meta:     in.Meta,
		Headers:  headers,
	})
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}

	w := c.Writer
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	c.Status(http.StatusOK)
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

func (h *httpHandlersSSE) providers(c *gin.Context) {
	cfg := h.reg.Config() // adicione um getter simples no registry
	type item struct{ Name, Type string }
	out := []item{}
	for name, pc := range cfg.Providers {
		out = append(out, item{Name: name, Type: pc.Type()})
	}
	c.JSON(http.StatusOK, gin.H{"providers": out})
}

func (h *httpHandlersSSE) session(c *gin.Context) {
	// Implementar lógica para gerenciar sessões
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (h *httpHandlersSSE) authLoginPassthrough(c *gin.Context) {
	// Implementar lógica para login via passthrough
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (h *httpHandlersSSE) stateExport(c *gin.Context) {
	// Implementar lógica para exportar estado
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (h *httpHandlersSSE) stateImport(c *gin.Context) {
	// Implementar lógica para importar estado
	c.AbortWithStatus(http.StatusNotImplemented)
}
