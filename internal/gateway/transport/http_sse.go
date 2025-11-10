// Package transport implements HTTP transport handlers for server-sent events (SSE) in the Grompt application.
package transport

import (
	"encoding/json"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/kubex-ecosystem/grompt/internal/advise"
	"github.com/kubex-ecosystem/grompt/internal/gateway/middleware"
	"github.com/kubex-ecosystem/grompt/internal/gateway/registry"
	"github.com/kubex-ecosystem/grompt/internal/interfaces"
	"github.com/kubex-ecosystem/grompt/internal/scorecard"
)

type httpHandlersSSE struct {
	reg        *registry.Registry
	middleware *middleware.ProductionMiddleware
	engine     *scorecard.Engine // Add scorecard engine
}

func WireHTTPSSE(router gin.IRouter, reg *registry.Registry, mw *middleware.ProductionMiddleware) {
	hh := &httpHandlersSSE{reg: reg, middleware: mw, engine: nil} // TODO: Initialize engine when ready
	router.GET("/healthz", hh.healthz)

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

	if hm := h.getHealthMonitor(); hm != nil {
		if health, ok := hm.GetHealth(in.Provider); ok && health.Status == middleware.HealthUnhealthy {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":    "provider temporarily unavailable",
				"provider": in.Provider,
				"status":   health.Status.String(),
			})
			return
		}
	}
	headers := map[string]string{
		"x-external-api-key": c.GetHeader("x-external-api-key"),
		"x-tenant-id":        c.GetHeader("x-tenant-id"),
		"x-user-id":          c.GetHeader("x-user-id"),
	}
	w := c.Writer
	fl, ok := w.(http.Flusher)
	if !ok {
		c.String(http.StatusInternalServerError, "streaming unsupported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)

	var wroteAny atomic.Bool
	send := func(payload map[string]any) {
		if len(payload) == 0 {
			return
		}
		payload["provider"] = in.Provider
		b, _ := json.Marshal(payload)
		_, _ = w.Write([]byte("data: "))
		_, _ = w.Write(b)
		_, _ = w.Write([]byte("\n\n"))
		fl.Flush()
		wroteAny.Store(true)
	}

	streamErr := h.wrapWithMiddleware(in.Provider, func() error {
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
			return err
		}

		coalescer := NewSSECoalescer(func(content string) {
			send(map[string]any{"content": content})
		})
		defer coalescer.Close()

		for msg := range ch {
			if msg.Content != "" {
				coalescer.AddChunk(msg.Content)
			}

			if msg.Done {
				coalescer.Close()
				payload := map[string]any{"done": true}
				if msg.Usage != nil {
					payload["usage"] = msg.Usage
				}
				send(payload)
			}
		}
		return nil
	})

	if streamErr != nil {
		if !wroteAny.Load() {
			c.Status(http.StatusBadGateway)
		}
		send(map[string]any{"error": streamErr.Error(), "done": true})
	}
}

// /v1/providers — lista nomes e tipos carregados (pra pintar “verde” no dropdown)

func (h *httpHandlersSSE) providers(c *gin.Context) {
	cfg := h.reg.Config() // adicione um getter simples no registry
	hm := h.getHealthMonitor()
	type item struct {
		Name   string                  `json:"name"`
		Type   string                  `json:"type"`
		Health *middleware.HealthCheck `json:"health,omitempty"`
	}
	out := []item{}
	for name, pc := range cfg.Providers {
		var health *middleware.HealthCheck
		if hm != nil {
			if chk, ok := hm.GetHealth(name); ok {
				health = chk
			}
		}
		out = append(out, item{Name: name, Type: pc.Type(), Health: health})
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

func (h *httpHandlersSSE) healthz(c *gin.Context) {
	resp := gin.H{"status": "ok"}
	if hm := h.getHealthMonitor(); hm != nil {
		resp["providers"] = hm.GetAllHealth()
	}
	c.JSON(http.StatusOK, resp)
}

func (h *httpHandlersSSE) getHealthMonitor() *middleware.HealthMonitor {
	if h == nil || h.middleware == nil {
		return nil
	}
	return h.middleware.GetHealthMonitor()
}

func (h *httpHandlersSSE) wrapWithMiddleware(provider string, fn func() error) error {
	if h.middleware == nil {
		return fn()
	}
	return h.middleware.WrapProvider(provider, fn)
}
