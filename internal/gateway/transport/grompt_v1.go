// Package transport implements the Grompt V1 API routes
// Following the Analyzer structure pattern with Multi-Provider integration
package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/gateway/middleware"
	"github.com/kubex-ecosystem/grompt/internal/gateway/registry"
	providers "github.com/kubex-ecosystem/grompt/internal/types"
)

// GromptV1Handlers holds the V1 API handlers for Grompt
type GromptV1Handlers struct {
	registry             *registry.Registry
	productionMiddleware *middleware.ProductionMiddleware
	gobeBaseURL          string
}

// NewGromptV1Handlers creates a new set of V1 handlers
func NewGromptV1Handlers(reg *registry.Registry, prodMiddleware *middleware.ProductionMiddleware, gobeBaseURL string) *GromptV1Handlers {
	return &GromptV1Handlers{
		registry:             reg,
		productionMiddleware: prodMiddleware,
		gobeBaseURL:          gobeBaseURL,
	}
}

// WireGromptV1Routes sets up the Grompt V1 API routes
func WireGromptV1Routes(mux *http.ServeMux, handlers *GromptV1Handlers) {
	// Core Grompt V1 routes with middleware
	mux.HandleFunc("/v1/generate", handlers.withTimeout(handlers.withConcurrencyLimit(handlers.generatePrompt)))
	mux.HandleFunc("/v1/generate/stream", handlers.withTimeout(handlers.withConcurrencyLimit(handlers.generatePromptStream)))
	mux.HandleFunc("/v1/providers", handlers.listProviders)
	mux.HandleFunc("/v1/health", handlers.healthCheck)

	// GoBE proxy route - delegates auth, storage, billing
	mux.HandleFunc("/v1/proxy/", handlers.withTimeout(handlers.proxyToGoBE))

	log.Println("✅ Grompt V1 API routes wired successfully")
}

// GenerateRequest represents a prompt generation request
type GenerateRequest struct {
	Provider    string                 `json:"provider"`
	Model       string                 `json:"model,omitempty"`
	Ideas       []string               `json:"ideas"`
	Purpose     string                 `json:"purpose,omitempty"` // "code", "creative", "analysis"
	Context     map[string]interface{} `json:"context,omitempty"`
	Temperature float32                `json:"temperature,omitempty"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Stream      bool                   `json:"stream,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
}

// GenerateResponse represents the response for prompt generation
type GenerateResponse struct {
	ID        string                 `json:"id"`
	Object    string                 `json:"object"`
	CreatedAt int64                  `json:"created_at"`
	Provider  string                 `json:"provider"`
	Model     string                 `json:"model"`
	Prompt    string                 `json:"prompt"`
	Ideas     []string               `json:"ideas"`
	Purpose   string                 `json:"purpose"`
	Usage     *providers.Usage       `json:"usage,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// generatePrompt handles POST /v1/generate - synchronous prompt generation
func (h *GromptV1Handlers) generatePrompt(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Enhanced validation
	if err := h.validateGenerateRequest(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Sanitize and set defaults
	h.sanitizeAndSetDefaults(&req)

	// Resolve provider
	provider := h.registry.Resolve(req.Provider)
	if provider == nil {
		http.Error(w, fmt.Sprintf("Provider '%s' not found", req.Provider), http.StatusNotFound)
		return
	}

	// Check provider availability
	if ok := provider.Available(); !ok {
		http.Error(w, "Provider unavailable", http.StatusServiceUnavailable)
		return
	}

	// Build prompt engineering message
	systemPrompt := h.buildPromptEngineeringSystem(req.Purpose)
	userPrompt := h.buildPromptFromIdeas(req.Ideas, req.Context)

	// Prepare chat request
	chatReq := providers.ChatRequest{
		Provider: req.Provider,
		Model:    req.Model,
		Messages: []providers.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temp:   req.Temperature,
		Stream: false, // Synchronous response
		Meta:   req.Meta,
		Headers: map[string]string{
			"x-external-api-key": r.Header.Get("x-external-api-key"),
			"x-tenant-id":        r.Header.Get("x-tenant-id"),
			"x-user-id":          r.Header.Get("x-user-id"),
		},
	}

	// Execute chat completion
	ch, err := provider.Chat(r.Context(), chatReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Chat request failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Collect full response
	var fullContent strings.Builder
	var usage *providers.Usage

	for chunk := range ch {
		if chunk.Error != "" {
			http.Error(w, fmt.Sprintf("Generation failed: %s", chunk.Error), http.StatusInternalServerError)
			return
		}

		if chunk.Content != "" {
			fullContent.WriteString(chunk.Content)
		}

		if chunk.Done && chunk.Usage != nil {
			usage = chunk.Usage
		}
	}

	// Record metrics for observability
	duration := time.Since(startTime)
	tokens := 0
	cost := 0.0
	if usage != nil {
		tokens = usage.Tokens
		cost = usage.CostUSD
	}
	h.recordMetrics("/v1/generate", req.Provider, req.Model, duration, tokens, cost, nil)

	// Build response
	response := GenerateResponse{
		ID:        fmt.Sprintf("gen_%d", time.Now().Unix()),
		Object:    "prompt.generation",
		CreatedAt: time.Now().Unix(),
		Provider:  req.Provider,
		Model:     req.Model,
		Prompt:    fullContent.String(),
		Ideas:     req.Ideas,
		Purpose:   req.Purpose,
		Usage:     usage,
		Metadata: map[string]interface{}{
			"temperature": req.Temperature,
			"max_tokens":  req.MaxTokens,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// generatePromptStream handles GET /v1/generate/stream - SSE streaming generation
func (h *GromptV1Handlers) generatePromptStream(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters for streaming request
	provider := r.URL.Query().Get("provider")
	ideas := r.URL.Query()["ideas"]
	purpose := r.URL.Query().Get("purpose")
	model := r.URL.Query().Get("model")
	temperature := r.URL.Query().Get("temperature")

	if provider == "" {
		http.Error(w, "Provider parameter is required", http.StatusBadRequest)
		return
	}
	if len(ideas) == 0 {
		http.Error(w, "At least one idea parameter is required", http.StatusBadRequest)
		return
	}

	// Resolve provider
	prov := h.registry.Resolve(provider)
	if prov == nil {
		http.Error(w, fmt.Sprintf("Provider '%s' not found", provider), http.StatusNotFound)
		return
	}

	// Check provider availability
	if ok := prov.Available(); !ok {
		http.Error(w, "Provider unavailable", http.StatusServiceUnavailable)
		return
	}

	// Set default model if not provided
	if model == "" {
		cfg := h.registry.GetConfig()
		if providerConfig, exists := cfg.Providers[provider]; exists {
			model = providerConfig.DefaultModel
		}
	}

	// Parse temperature parameter
	temp := float32(0.7) // default
	if temperature != "" {
		if parsedTemp, err := strconv.ParseFloat(temperature, 32); err == nil {
			if parsedTemp >= 0 && parsedTemp <= 2.0 {
				temp = float32(parsedTemp)
			}
		}
	}

	// Build prompts
	systemPrompt := h.buildPromptEngineeringSystem(purpose)
	userPrompt := h.buildPromptFromIdeas(ideas, nil)

	// Prepare chat request for streaming
	chatReq := providers.ChatRequest{
		Provider: provider,
		Model:    model,
		Messages: []providers.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temp:   temp,
		Stream: true,
		Headers: map[string]string{
			"x-external-api-key": r.Header.Get("x-external-api-key"),
			"x-tenant-id":        r.Header.Get("x-tenant-id"),
			"x-user-id":          r.Header.Get("x-user-id"),
		},
	}

	// Start chat completion
	ch, err := prov.Chat(r.Context(), chatReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Chat request failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "content-type, authorization, x-external-api-key, x-tenant-id, x-user-id")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Send initial event
	fmt.Fprintf(w, "data: %s\n\n", mustMarshalJSON(map[string]interface{}{
		"event":    "generation.started",
		"provider": provider,
		"model":    model,
		"ideas":    ideas,
	}))
	flusher.Flush()

	// Create SSE coalescer for smooth streaming UX
	coalescer := NewSSECoalescer(func(content string) {
		data := mustMarshalJSON(map[string]interface{}{
			"event":   "generation.chunk",
			"content": content,
		})
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	})
	defer coalescer.Close()

	// Stream response chunks with coalescence
	for chunk := range ch {
		if chunk.Error != "" {
			// Flush pending content before error
			coalescer.Close()

			fmt.Fprintf(w, "data: %s\n\n", mustMarshalJSON(map[string]interface{}{
				"event": "generation.error",
				"error": chunk.Error,
			}))
			flusher.Flush()
			return
		}

		if chunk.Content != "" {
			// Add to coalescer instead of immediate flush
			coalescer.AddChunk(chunk.Content)
		}

		if chunk.Done {
			// Flush any remaining content before completion
			coalescer.Close()

			fmt.Fprintf(w, "data: %s\n\n", mustMarshalJSON(map[string]interface{}{
				"event": "generation.complete",
				"usage": chunk.Usage,
			}))
			flusher.Flush()

			// Record metrics for streaming
			duration := time.Since(startTime)
			tokens := 0
			cost := 0.0
			if chunk.Usage != nil {
				tokens = chunk.Usage.Tokens
				cost = chunk.Usage.CostUSD
			}
			h.recordMetrics("/v1/generate/stream", provider, model, duration, tokens, cost, nil)
			break
		}
	}
}

// listProviders handles GET /v1/providers
func (h *GromptV1Handlers) listProviders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	providerNames := h.registry.ListProviders()
	config := h.registry.GetConfig()

	providers := make([]map[string]interface{}, 0, len(providerNames))
	for _, name := range providerNames {
		providerInfo := map[string]interface{}{
			"name":      name,
			"available": true,
		}

		// Add configuration info
		if providerConfig, exists := config.Providers[name]; exists {
			providerInfo["type"] = providerConfig.Type
			providerInfo["default_model"] = providerConfig.DefaultModel
		}

		// Check availability
		if provider := h.registry.Resolve(name); provider != nil {
			if ok := provider.Available(); !ok {
				providerInfo["available"] = false
				providerInfo["error"] = "Provider unavailable"
			}
		}

		providers = append(providers, providerInfo)
	}

	response := map[string]interface{}{
		"object":    "list",
		"data":      providers,
		"has_more":  false,
		"timestamp": time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// healthCheck handles GET /v1/health - intelligent health check
func (h *GromptV1Handlers) healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	health := map[string]interface{}{
		"status":       "healthy",
		"service":      "grompt-v1",
		"timestamp":    time.Now().Unix(),
		"version":      "1.0.0",
		"dependencies": make(map[string]interface{}),
	}

	overallHealthy := true

	// Check provider connectivity
	providerNames := h.registry.ListProviders()
	providerHealth := make(map[string]interface{})

	for _, name := range providerNames {
		provider := h.registry.Resolve(name)
		if provider == nil {
			providerHealth[name] = map[string]interface{}{
				"status": "unavailable",
				"error":  "provider not found",
			}
			overallHealthy = false
			continue
		}

		if ok := provider.Available(); !ok {
			providerHealth[name] = map[string]interface{}{
				"status": "unhealthy",
				"error":  "Provider unavailable",
			}
			overallHealthy = false
		} else {
			providerHealth[name] = map[string]interface{}{
				"status": "healthy",
			}
		}
	}
	health["dependencies"].(map[string]interface{})["providers"] = providerHealth

	// Check GoBE proxy connectivity
	gobeHealth := h.checkGobeHealth()
	health["dependencies"].(map[string]interface{})["gobe_proxy"] = gobeHealth
	if gobeHealth["status"] != "healthy" {
		overallHealthy = false
	}

	// Set overall status
	if !overallHealthy {
		health["status"] = "degraded"
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// proxyToGoBE handles POST /v1/proxy/* - delegates to GoBE for auth, storage, billing
func (h *GromptV1Handlers) proxyToGoBE(w http.ResponseWriter, r *http.Request) {
	if h.gobeBaseURL == "" {
		http.Error(w, "GoBE proxy not configured", http.StatusServiceUnavailable)
		return
	}

	// Extract path after /v1/proxy/
	targetPath := strings.TrimPrefix(r.URL.Path, "/v1/proxy")
	if targetPath == "" {
		targetPath = "/"
	}

	// Build target URL
	targetURL, err := url.Parse(h.gobeBaseURL + targetPath)
	if err != nil {
		http.Error(w, "Invalid GoBE URL", http.StatusInternalServerError)
		return
	}

	// Copy query parameters
	targetURL.RawQuery = r.URL.RawQuery

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Modify request to preserve important headers
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// Preserve authentication and request tracking headers
		req.Header.Set("Authorization", r.Header.Get("Authorization"))
		req.Header.Set("X-Request-Id", r.Header.Get("X-Request-Id"))
		req.Header.Set("X-Tenant-Id", r.Header.Get("X-Tenant-Id"))
		req.Header.Set("X-User-Id", r.Header.Get("X-User-Id"))

		// Add Grompt identification
		req.Header.Set("X-Forwarded-By", "grompt-gateway")
		req.Header.Set("X-Original-Path", r.URL.Path)
	}

	// Set error handler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("GoBE proxy error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "gobe_proxy_error",
			"message": "Failed to reach GoBE backend",
			"details": err.Error(),
		})
	}

	// Execute proxy
	proxy.ServeHTTP(w, r)
}

// Helper functions

func (h *GromptV1Handlers) buildPromptEngineeringSystem(purpose string) string {
	baseSystem := `You are an expert prompt engineer. Your task is to transform user ideas into clear, effective prompts for AI models.

Apply proper prompt engineering techniques:
- Clear structure and organization
- Specific instructions and context
- Expected output format
- Examples where helpful
- Appropriate persona and tone`

	switch purpose {
	case "code":
		return baseSystem + `

Focus on coding tasks:
- Specify programming languages and frameworks
- Include technical requirements and constraints
- Define expected output format (functions, classes, etc.)
- Mention testing and documentation needs`

	case "creative":
		return baseSystem + `

Focus on creative writing:
- Establish tone, style, and genre
- Define target audience and length
- Include structure and formatting requirements
- Specify any creative constraints or themes`

	case "analysis":
		return baseSystem + `

Focus on analytical tasks:
- Define analysis methodology and frameworks
- Specify data sources and formats
- Include visualization or reporting requirements
- Mention accuracy and objectivity standards`

	default:
		return baseSystem + `

Create versatile prompts that can be adapted for various use cases.`
	}
}

func (h *GromptV1Handlers) buildPromptFromIdeas(ideas []string, context map[string]interface{}) string {
	var prompt strings.Builder

	prompt.WriteString("Transform these ideas into a well-structured prompt:\n\n")

	for i, idea := range ideas {
		prompt.WriteString(fmt.Sprintf("%d. %s\n", i+1, idea))
	}

	if len(context) > 0 {
		prompt.WriteString("\nAdditional context:\n")
		for key, value := range context {
			prompt.WriteString(fmt.Sprintf("- %s: %v\n", key, value))
		}
	}

	prompt.WriteString("\nCreate a comprehensive, well-engineered prompt that incorporates these ideas effectively.")

	return prompt.String()
}

func (h *GromptV1Handlers) checkGobeHealth() map[string]interface{} {
	if h.gobeBaseURL == "" {
		return map[string]interface{}{
			"status":  "not_configured",
			"message": "GoBE proxy URL not configured",
		}
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Make simple health check request
	req, err := http.NewRequestWithContext(ctx, "GET", h.gobeBaseURL+"/health", nil)
	if err != nil {
		return map[string]interface{}{
			"status": "unhealthy",
			"error":  fmt.Sprintf("failed to create request: %v", err),
		}
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{
			"status": "unhealthy",
			"error":  fmt.Sprintf("connection failed: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return map[string]interface{}{
			"status": "unhealthy",
			"error":  fmt.Sprintf("health check failed with status: %d", resp.StatusCode),
		}
	}

	return map[string]interface{}{
		"status":        "healthy",
		"response_time": "< 5s",
	}
}

// validateGenerateRequest performs comprehensive validation
func (h *GromptV1Handlers) validateGenerateRequest(req *GenerateRequest) error {
	if req.Provider == "" {
		return fmt.Errorf("provider is required")
	}

	if len(req.Ideas) == 0 {
		return fmt.Errorf("at least one idea is required")
	}

	// Validate ideas length
	for i, idea := range req.Ideas {
		if strings.TrimSpace(idea) == "" {
			return fmt.Errorf("idea %d cannot be empty", i+1)
		}
		if len(idea) > 1000 {
			return fmt.Errorf("idea %d exceeds maximum length of 1000 characters", i+1)
		}
	}

	// Validate purpose if provided
	if req.Purpose != "" {
		validPurposes := map[string]bool{
			"code":     true,
			"creative": true,
			"analysis": true,
			"general":  true,
		}
		if !validPurposes[req.Purpose] {
			return fmt.Errorf("invalid purpose '%s', must be one of: code, creative, analysis, general", req.Purpose)
		}
	}

	// Validate temperature range
	if req.Temperature < 0 || req.Temperature > 2.0 {
		return fmt.Errorf("temperature must be between 0 and 2.0")
	}

	// Validate max tokens
	if req.MaxTokens < 0 || req.MaxTokens > 32000 {
		return fmt.Errorf("max_tokens must be between 0 and 32000")
	}

	return nil
}

// sanitizeAndSetDefaults applies defaults and sanitizes the request
func (h *GromptV1Handlers) sanitizeAndSetDefaults(req *GenerateRequest) {
	// Set default temperature
	if req.Temperature == 0 {
		req.Temperature = 0.7
	}

	// Set default purpose
	if req.Purpose == "" {
		req.Purpose = "general"
	}

	// Set default model if not provided
	if req.Model == "" {
		cfg := h.registry.GetConfig()
		if providerConfig, exists := cfg.Providers[req.Provider]; exists {
			req.Model = providerConfig.DefaultModel
		}
	}

	// Sanitize ideas
	for i, idea := range req.Ideas {
		req.Ideas[i] = strings.TrimSpace(idea)
	}
}

// recordMetrics logs request metrics for observability
func (h *GromptV1Handlers) recordMetrics(endpoint, provider, model string, duration time.Duration, tokens int, cost float64, err error) {
	statusCode := "200"
	if err != nil {
		statusCode = "500"
	}

	log.Printf("[METRICS] endpoint=%s provider=%s model=%s duration=%v tokens=%d cost=%.6f status=%s",
		endpoint, provider, model, duration, tokens, cost, statusCode)
}

// withTimeout adds request timeout middleware
func (h *GromptV1Handlers) withTimeout(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timeout := 120 * time.Second // Default 2 minutes
		if r.URL.Path == "/v1/generate/stream" {
			timeout = 300 * time.Second // 5 minutes for streaming
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		// Channel to capture handler completion
		done := make(chan struct{})
		go func() {
			defer close(done)
			handler(w, r.WithContext(ctx))
		}()

		select {
		case <-done:
			// Handler completed normally
		case <-ctx.Done():
			// Timeout occurred
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
			log.Printf("Request timeout: %s %s", r.Method, r.URL.Path)
		}
	}
}

// concurrentRequests tracks active requests
var (
	activeRequests        int64
	maxConcurrentRequests int64 = 50 // Configurable limit
)

// withConcurrencyLimit limits concurrent requests
func (h *GromptV1Handlers) withConcurrencyLimit(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Atomic increment
		current := atomic.AddInt64(&activeRequests, 1)
		defer atomic.AddInt64(&activeRequests, -1)

		if current > maxConcurrentRequests {
			http.Error(w, "Too many concurrent requests", http.StatusTooManyRequests)
			log.Printf("Concurrent request limit exceeded: %d/%d", current, maxConcurrentRequests)
			return
		}

		handler(w, r)
	}
}

// mustMarshalJSON marshals JSON safely
func mustMarshalJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		return `{"error": "json_marshal_failed"}`
	}
	return string(data)
}
