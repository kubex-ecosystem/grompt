// Package transport sets up HTTP routes and handlers for the Grompt Gateway,
// including merged Repository Intelligence endpoints.
package transport

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/config"
	"github.com/kubex-ecosystem/grompt/internal/gateway/health"
	"github.com/kubex-ecosystem/grompt/internal/gateway/middleware"
	"github.com/kubex-ecosystem/grompt/internal/gateway/registry"
	"github.com/kubex-ecosystem/grompt/internal/handlers/lookatni"
	"github.com/kubex-ecosystem/grompt/internal/scorecard"
	providers "github.com/kubex-ecosystem/grompt/internal/types"
	"github.com/kubex-ecosystem/grompt/internal/web"
	"github.com/kubex-ecosystem/grompt/internal/webhook"
)

// httpHandlers holds the HTTP route handlers
type httpHandlers struct {
	registry             *registry.Registry
	productionMiddleware *middleware.ProductionMiddleware
	engine               *scorecard.Engine    // Repository Intelligence engine
	lookAtniHandler      *lookatni.Handler    // LookAtni integration
	webhookHandler       *webhook.HTTPHandler // Meta-recursive webhook handler
	healthEngine         *health.Engine       // AI Provider health monitoring
	healthRegistry       *health.ProberRegistry
	healthScheduler      *health.Scheduler // Background health checks
}

// WireHTTP sets up HTTP routes
func WireHTTP(mux *http.ServeMux, reg *registry.Registry, prodMiddleware *middleware.ProductionMiddleware) {
	// Initialize Grompt V1 handlers with GoBE proxy support
	gobeBaseURL := getEnv("GOBE_BASE_URL", "")
	gromptV1Handlers := NewGromptV1Handlers(reg, prodMiddleware, gobeBaseURL)

	// Wire Grompt V1 specific routes
	WireGromptV1Routes(mux, gromptV1Handlers)
	// Initialize LookAtni handler
	workDir := "./lookatni_workspace" // TODO: Make configurable
	lookAtniHandler := lookatni.NewHandler(workDir)

	// Initialize webhook handler (mock for now - TODO: implement real actors)
	webhookHandler := webhook.NewHTTPHandler(nil) // TODO: Initialize with real handler

	// Initialize AI Provider Health Monitoring
	healthStore := health.NewStore()
	healthRegistry := health.NewProberRegistry()

	// Registra probers no registry local
	healthRegistry.Register(health.NewGroqProber())
	healthRegistry.Register(health.NewGeminiProber())

	// Cria engine com probers do registry
	groqProber := health.NewGroqProber()
	geminiProber := health.NewGeminiProber()
	healthEngine := health.NewEngine(healthStore, groqProber, geminiProber)

	// Initialize Background Health Scheduler - ARQUITETURA QUE NÃO SE SABOTA! 🔥
	schedulerConfig := health.DefaultSchedulerConfig()
	schedulerConfig.LogVerbose = false // Production mode
	healthScheduler := health.NewScheduler(healthEngine, healthRegistry, schedulerConfig)

	// Start scheduler in background
	if err := healthScheduler.Start(); err != nil {
		log.Printf("⚠️  Failed to start health scheduler: %v", err)
	}

	h := &httpHandlers{
		registry:             reg,
		productionMiddleware: prodMiddleware,
		engine:               nil, // TODO: Initialize scorecard engine with real clients
		lookAtniHandler:      lookAtniHandler,
		webhookHandler:       webhookHandler,
		healthEngine:         healthEngine,
		healthRegistry:       healthRegistry,
		healthScheduler:      healthScheduler,
	}

	// Web Interface - Frontend embarcado! 🚀
	webHandler, err := web.NewHandler()
	if err != nil {
		log.Printf("⚠️  Failed to initialize web interface: %v", err)
	} else {
		// Register web interface on /app/* and root
		mux.Handle("/app/", http.StripPrefix("/app", webHandler))
		// Root path serves the frontend (but with lower priority than API endpoints)
		mux.Handle("/", webHandler)
		log.Println("✅ Web interface enabled at /app/ and /")
	}

	// API endpoints (higher priority routes)
	mux.HandleFunc("/healthz", h.healthCheck)
	mux.HandleFunc("/chat", h.chatSSE)
	mux.HandleFunc("/providers", h.listProviders)
	mux.HandleFunc("/advise", h.handleAdvise)
	mux.HandleFunc("/status", h.productionStatus)

	// Repository Intelligence endpoints - MERGE POINT! 🚀
	mux.HandleFunc("/api/v1/scorecard", h.handleRepositoryScorecard)
	mux.HandleFunc("/api/v1/scorecard/advice", h.handleScorecardAdvice)
	mux.HandleFunc("/api/v1/metrics/ai", h.handleAIMetrics)
	mux.HandleFunc("/api/v1/health", h.handleRepositoryHealth)

	// AI Provider Health Monitoring - ARQUITETURA QUE NÃO SE SABOTA! 🔥
	health.RegisterRoutes(mux, h.healthEngine, h.healthRegistry)
	mux.HandleFunc("/health/scheduler/stats", h.handleSchedulerStats)
	mux.HandleFunc("/health/scheduler/force", h.handleSchedulerForce)

	// LookAtni Integration endpoints - CODE NAVIGATION! 🔍
	mux.HandleFunc("/api/v1/lookatni/extract", h.lookAtniHandler.HandleExtractProject)
	mux.HandleFunc("/api/v1/lookatni/archive", h.lookAtniHandler.HandleCreateArchive)
	mux.HandleFunc("/api/v1/lookatni/download/", h.lookAtniHandler.HandleDownloadArchive)
	mux.HandleFunc("/api/v1/lookatni/projects", h.lookAtniHandler.HandleListExtractedProjects)
	mux.HandleFunc("/api/v1/lookatni/projects/", h.lookAtniHandler.HandleProjectFragments)

	// Meta-Recursive Webhook endpoints - INSANIDADE RACIONAL! 🔄
	mux.HandleFunc("/v1/webhooks", h.webhookHandler.HandleWebhook)
	mux.HandleFunc("/v1/webhooks/health", h.webhookHandler.HealthCheck)

	log.Println("✅ LookAtni integration enabled - Code extraction and navigation ready!")
	log.Println("🔄 Meta-recursive webhook system enabled")
	log.Println("🔥 AI Provider Health Monitoring enabled")
}

// healthCheck provides a simple health endpoint
func (h *httpHandlers) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "grompt-gw",
	})
}

// listProviders returns available providers with health status
func (h *httpHandlers) listProviders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	providerNames := h.registry.ListProviders()
	config := h.registry.GetConfig()

	// Get health status for all providers from health monitor
	healthStatuses := make(map[string]interface{})
	if h.productionMiddleware != nil {
		healthMonitor := h.productionMiddleware.GetHealthMonitor()
		if healthMonitor != nil {
			allHealth := healthMonitor.GetAllHealth()
			for providerName, health := range allHealth {
				healthStatuses[providerName] = map[string]interface{}{
					"status":        health.Status.String(),
					"last_check":    health.LastCheck,
					"response_time": health.ResponseTime.String(),
					"uptime":        health.Uptime,
					"error":         health.ErrorMsg,
				}
			}
		}
	}

	// Enrich providers with health information
	enrichedProviders := make([]map[string]interface{}, 0, len(providerNames))
	for _, providerName := range providerNames {
		enrichedProvider := map[string]interface{}{
			"name":      providerName,
			"type":      providerName, // TODO: Get actual type from registry
			"available": true,         // Default to true
		}

		// Add health status if available
		if healthStatus, exists := healthStatuses[providerName]; exists {
			enrichedProvider["health"] = healthStatus
			// Update availability based on health status
			if healthMap, ok := healthStatus.(map[string]interface{}); ok {
				if status, ok := healthMap["status"].(string); ok {
					enrichedProvider["available"] = (status == "healthy" || status == "degraded")
				}
			}
		} else {
			// No health data available
			enrichedProvider["health"] = map[string]interface{}{
				"status":     "unknown",
				"last_check": nil,
				"uptime":     100.0,
				"error":      "",
			}
		}

		enrichedProviders = append(enrichedProviders, enrichedProvider)
	}

	response := map[string]interface{}{
		"providers": enrichedProviders,
		"config":    config.Providers,
		"timestamp": "2024-01-01T00:00:00Z", // TODO: Use real timestamp
		"service":   "grompt-gateway",
		"version":   "v1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// chatSSE handles chat completion with Server-Sent Events
func (h *httpHandlers) chatSSE(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req providers.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Provider == "" {
		http.Error(w, "Provider is required", http.StatusBadRequest)
		return
	}

	provider := h.registry.Resolve(req.Provider)
	if provider == nil {
		http.Error(w, fmt.Sprintf("Provider '%s' not found", req.Provider), http.StatusBadRequest)
		return
	}

	// Check if provider is available
	if ok := provider.Available(); !ok {
		http.Error(w, "Provider unavailable", http.StatusServiceUnavailable)
		return
	}

	// Handle BYOK (Bring Your Own Key)
	if externalKey := r.Header.Get("x-external-api-key"); externalKey != "" {
		// TODO: Implement secure BYOK handling
		// For now, we'll pass it through meta
		if req.Meta == nil {
			req.Meta = make(map[string]interface{})
		}
		req.Meta["external_api_key"] = externalKey
	}

	// Set default temperature if not provided
	if req.Temp == 0 {
		req.Temp = 0.7
	}

	// Force streaming for SSE
	req.Stream = true

	// Start chat completion
	ch, err := provider.Chat(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Chat request failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Create SSE coalescer to improve streaming UX
	coalescer := NewSSECoalescer(func(content string) {
		data, _ := json.Marshal(map[string]interface{}{
			"content": content,
			"done":    false,
		})
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	})
	defer coalescer.Close()

	// Stream the response with coalescence
	for chunk := range ch {
		if chunk.Error != "" {
			// Flush any pending content before error
			coalescer.Close()

			// Send error event
			data, _ := json.Marshal(map[string]interface{}{
				"error": chunk.Error,
				"done":  true,
			})
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
			return
		}

		if chunk.Content != "" {
			// Add to coalescer instead of immediate flush
			coalescer.AddChunk(chunk.Content)
		}

		if chunk.Done {
			// Flush any remaining content before final chunk
			coalescer.Close()

			// Send final chunk with usage info
			data, _ := json.Marshal(map[string]interface{}{
				"done":  true,
				"usage": chunk.Usage,
			})
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()

			// Log usage for monitoring
			if chunk.Usage != nil {
				log.Printf("Usage: provider=%s model=%s tokens=%d latency=%dms cost=$%.6f",
					chunk.Usage.Provider, chunk.Usage.Model, chunk.Usage.Tokens,
					chunk.Usage.Ms, chunk.Usage.CostUSD)
			}
			break
		}
	}
}

// productionStatus returns comprehensive status including middleware metrics
func (h *httpHandlers) productionStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := map[string]interface{}{
		"service":   "grompt-gw",
		"status":    "healthy",
		"providers": h.registry.ListProviders(),
	}

	// Add production middleware status if available
	if h.productionMiddleware != nil {
		status["production_features"] = h.productionMiddleware.GetStatus()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// handleAdvise handles POST /v1/advise for AI-powered analysis advice
func (h *httpHandlers) handleAdvise(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check BF1_MODE restrictions
	bf1Config := config.GetBF1Config()
	if bf1Config.Enabled {
		// In BF1 mode, add guard-rails and limitations
		w.Header().Set("X-BF1-Mode", "true")
		w.Header().Set("X-BF1-WIP-Cap", fmt.Sprintf("%d", bf1Config.WIPCap))
	}

	// Parse request body
	var req struct {
		Mode     string                 `json:"mode"`     // "exec" or "code"
		Provider string                 `json:"provider"` // optional, defaults to first available
		Context  map[string]interface{} `json:"context"`  // repository, hotspots, scorecard
		Options  map[string]interface{} `json:"options"`  // timeout_sec, temperature
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	// Validate mode
	if req.Mode != "exec" && req.Mode != "code" {
		http.Error(w, "Mode must be 'exec' or 'code'", http.StatusBadRequest)
		return
	}

	// Set SSE headers for streaming response
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Generate mock response based on mode
	if req.Mode == "exec" {
		// Simulate exec mode response
		execResponse := map[string]interface{}{
			"summary": map[string]interface{}{
				"grade":               "B+",
				"chi":                 72.5,
				"lead_time_p95_hours": 24.0,
				"deploys_per_week":    3.2,
			},
			"top_focus": []map[string]interface{}{
				{
					"title":      "Reduce Lead Time",
					"why":        "Long deployment cycles impacting delivery velocity",
					"kpi":        "lead_time_p95_hours",
					"target":     "< 12 hours",
					"confidence": 0.85,
				},
			},
			"quick_wins": []map[string]interface{}{
				{
					"action":        "Implement automated testing",
					"effort":        "M",
					"expected_gain": "20% faster deployments",
				},
			},
			"risks": []map[string]interface{}{
				{
					"risk":       "Technical debt accumulation",
					"mitigation": "Schedule dedicated refactoring sprints",
				},
			},
			"call_to_action": "Focus on deployment automation and testing coverage",
		}

		// Send as SSE
		fmt.Fprintf(w, "data: %s\n\n", mustMarshal(execResponse))

	} else { // code mode
		// Simulate code mode response
		codeResponse := map[string]interface{}{
			"chi_now": 72.5,
			"drivers": []map[string]interface{}{
				{
					"metric": "cyclomatic_complexity",
					"value":  15.2,
					"impact": "high",
				},
			},
			"refactor_plan": []map[string]interface{}{
				{
					"step":    1,
					"theme":   "Simplify complex functions",
					"actions": []string{"Break down large functions", "Extract common utilities"},
					"kpi":     "cyclomatic_complexity",
					"target":  "< 10",
				},
			},
			"guardrails": []string{"Maintain test coverage > 80%", "No functions > 50 lines"},
			"milestones": []map[string]interface{}{
				{
					"in_days": 14,
					"goal":    "Reduce complexity by 30%",
				},
			},
		}

		// Send as SSE
		fmt.Fprintf(w, "data: %s\n\n", mustMarshal(codeResponse))
	}

	// Send completion event
	fmt.Fprintf(w, "data: {\"done\": true}\n\n")

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

// Helper function to marshal JSON (panic on error for simplicity)
func mustMarshal(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}

// Repository Intelligence Handlers - MERGED! 🚀

// handleRepositoryScorecard handles GET /api/v1/scorecard
func (h *httpHandlers) handleRepositoryScorecard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Implement with real scorecard engine
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Schema-Version", "scorecard@1.0.0")
	w.Header().Set("X-Server-Version", "grompt-v1.0.0")

	// Placeholder response
	placeholder := map[string]interface{}{
		"status":  "not_implemented",
		"message": "Repository Intelligence API under development",
		"endpoints": []string{
			"/api/v1/scorecard",
			"/api/v1/scorecard/advice",
			"/api/v1/metrics/ai",
		},
	}
	json.NewEncoder(w).Encode(placeholder)
}

// handleScorecardAdvice handles POST /api/v1/scorecard/advice
func (h *httpHandlers) handleScorecardAdvice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Implement advice generation using existing advise system
	w.Header().Set("Content-Type", "application/json")
	placeholder := map[string]interface{}{
		"status":  "not_implemented",
		"message": "Will integrate with existing /v1/advise system",
	}
	json.NewEncoder(w).Encode(placeholder)
}

// handleAIMetrics handles GET /api/v1/metrics/ai
func (h *httpHandlers) handleAIMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Implement AI metrics calculation
	w.Header().Set("Content-Type", "application/json")
	placeholder := map[string]interface{}{
		"status":  "not_implemented",
		"message": "AI Metrics (HIR/AAC/TPH) calculation under development",
	}
	json.NewEncoder(w).Encode(placeholder)
}

// handleRepositoryHealth handles GET /api/v1/health
func (h *httpHandlers) handleRepositoryHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	health := map[string]interface{}{
		"status":  "healthy",
		"service": "repository-intelligence",
		"components": map[string]string{
			"scorecard_engine": "not_initialized",
			"dora_calculator":  "not_initialized",
			"chi_calculator":   "not_initialized",
			"ai_metrics":       "not_initialized",
		},
		"version": "grompt-v1.0.0",
	}
	json.NewEncoder(w).Encode(health)
}

// handleSchedulerStats returns health scheduler statistics
func (h *httpHandlers) handleSchedulerStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	stats := h.healthScheduler.GetStats()
	json.NewEncoder(w).Encode(stats)
}

// handleSchedulerForce forces immediate health checks for all providers
func (h *httpHandlers) handleSchedulerForce(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.healthScheduler.ForceCheck()

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":  "triggered",
		"message": "Force health check initiated for all providers",
		"time":    time.Now(),
	}
	json.NewEncoder(w).Encode(response)
}

// getEnv returns environment variable value or default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
