package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rafa-mori/grompt/internal/services/agents"
	ii "github.com/rafa-mori/grompt/internal/types"
)

type Handlers struct {
	config      ii.IConfig
	claudeAPI   ii.IAPIConfig
	openaiAPI   ii.IAPIConfig
	deepseekAPI ii.IAPIConfig
	chatGPTAPI  ii.IAPIConfig
	ollamaAPI   ii.IAPIConfig
	agentStore  *agents.Store
}

// Unified request structure
type UnifiedRequest struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
	Model     string `json:"model"`
	Provider  string `json:"provider"`
}

type UnifiedResponse struct {
	Response string     `json:"response"`
	Provider string     `json:"provider"`
	Model    string     `json:"model"`
	Usage    *UsageInfo `json:"usage,omitempty"`
}

type UsageInfo struct {
	PromptTokens     int     `json:"prompt_tokens,omitempty"`
	CompletionTokens int     `json:"completion_tokens,omitempty"`
	TotalTokens      int     `json:"total_tokens,omitempty"`
	EstimatedCost    float64 `json:"estimated_cost,omitempty"`
}

var llmKeyMap map[string]string

func NewHandlers(cfg ii.IConfig) *Handlers {
	hndr := &Handlers{}
	if cfg == nil {
		return &Handlers{}
	} else {
		hndr.config = cfg
	}

	llmKeyMap = hndr.getLLMAPIKeyMap(cfg)
	hndr.claudeAPI = ii.NewClaudeAPI(llmKeyMap["claude"])
	hndr.openaiAPI = ii.NewOpenAIAPI(llmKeyMap["openai"])
	hndr.chatGPTAPI = ii.NewChatGPTAPI(llmKeyMap["chatgpt"])
	hndr.deepseekAPI = ii.NewDeepSeekAPI(llmKeyMap["deepseek"])
	hndr.ollamaAPI = ii.NewOllamaAPI(llmKeyMap["ollama"])
	hndr.agentStore = agents.NewStore("agents.json")

	return hndr
}

func (h *Handlers) getLLMAPIKeyMap(cfg ii.IConfig) map[string]string {
	if llmKeyMap == nil {
		llmKeyMap = make(map[string]string)
	}
	if h.config == nil && cfg != nil {
		h.config = cfg
	}
	for _, provider := range []string{"claude", "openai", "chatgpt", "deepseek", "ollama"} {
		providerAPIKey := cfg.GetAPIKey(provider)
		llmKeyMap[provider] = providerAPIKey
	}
	return llmKeyMap
}

func (h *Handlers) isAPIEnabled(provider string) bool {

	switch provider {
	case "claude":
		return llmKeyMap["claude"] != ""
	case "openai":
		return llmKeyMap["openai"] != ""
	case "chatgpt":
		return llmKeyMap["chatgpt"] != ""
	case "deepseek":
		return llmKeyMap["deepseek"] != ""
	case "ollama":
		return llmKeyMap["ollama"] != ""
	default:
		return false
	}
}

func (h *Handlers) getAPIConfigKey(provider string) string {
	if h.config == nil {
		return ""
	}
	switch provider {
	case "claude":
		return h.config.GetAPIKey("claude")
	case "openai":
		return h.config.GetAPIKey("openai")
	case "chatgpt":
		return h.config.GetAPIKey("chatgpt")
	case "deepseek":
		return h.config.GetAPIKey("deepseek")
	case "ollama":
		return h.config.GetAPIKey("ollama")
	default:
		return ""
	}
}

func (h *Handlers) HandleConfig(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	config := h.config
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

func (h *Handlers) HandleClaude(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UnifiedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if key := h.config.GetAPIKey("claude"); key == "" {
		http.Error(w, "Claude API Key not configured", http.StatusServiceUnavailable)
		return
	}

	response, err := h.claudeAPI.Complete(req.Prompt, req.MaxTokens, "")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in Claude API: %v", err), http.StatusInternalServerError)
		return
	}

	result := UnifiedResponse{
		Response: response,
		Provider: "claude",
		Model:    "claude-3-5-sonnet-20241022",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handlers) HandleOpenAI(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UnifiedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if key := h.config.GetAPIKey("openai"); key == "" {
		http.Error(w, "OpenAI API Key not configured", http.StatusServiceUnavailable)
		return
	}

	// Use default model if not specified
	model := req.Model
	if model == "" {
		model = "gpt-4o-mini"
	}

	response, err := h.openaiAPI.Complete(req.Prompt, req.MaxTokens, model)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in OpenAI API: %v", err), http.StatusInternalServerError)
		return
	}

	result := UnifiedResponse{
		Response: response,
		Provider: "openai",
		Model:    model,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handlers) HandleDeepSeek(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UnifiedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if key := h.config.GetAPIKey("deepseek"); key == "" {
		http.Error(w, "DeepSeek API Key not configured", http.StatusServiceUnavailable)
		return
	}

	// Use default model if not specified
	model := req.Model
	if model == "" {
		model = "deepseek-chat"
	}

	response, err := h.deepseekAPI.Complete(req.Prompt, req.MaxTokens, model)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in DeepSeek API: %v", err), http.StatusInternalServerError)
		return
	}

	result := UnifiedResponse{
		Response: response,
		Provider: "deepseek",
		Model:    model,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handlers) HandleOllama(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UnifiedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Use default model if not specified
	model := req.Model
	if model == "" {
		model = "llama3.2"
	}

	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 2048
	}

	response, err := h.ollamaAPI.Complete(model, maxTokens, req.Prompt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in Ollama API: %v", err), http.StatusInternalServerError)
		return
	}

	result := UnifiedResponse{
		Response: response,
		Provider: "ollama",
		Model:    model,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// HandleUnified processes requests for multiple providers in a unified manner
func (h *Handlers) HandleUnified(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UnifiedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate provider
	if req.Provider == "" {
		http.Error(w, "Provider not specified", http.StatusBadRequest)
		return
	}

	var response string
	var err error
	var model string = req.Model
	var maxTokens int = req.MaxTokens

	switch req.Provider {
	case "claude":
		if key := h.config.GetAPIKey("claude"); key == "" {
			http.Error(w, "Claude API Key not configured", http.StatusServiceUnavailable)
			return
		}
		response, err = h.claudeAPI.Complete(req.Prompt, maxTokens, model)
		if model == "" {
			model = "claude-3-5-sonnet-20241022"
		}

	case "openai":
		if key := h.config.GetAPIKey("openai"); key == "" {
			http.Error(w, "OpenAI API Key not configured", http.StatusServiceUnavailable)
			return
		}
		if model == "" {
			model = "gpt-4o-mini"
		}
		response, err = h.openaiAPI.Complete(req.Prompt, req.MaxTokens, model)

	case "deepseek":
		if key := h.config.GetAPIKey("deepseek"); key == "" {
			http.Error(w, "DeepSeek API Key not configured", http.StatusServiceUnavailable)
			return
		}
		if model == "" {
			model = "deepseek-chat"
		}
		response, err = h.deepseekAPI.Complete(req.Prompt, req.MaxTokens, model)

	case "ollama":
		if model == "" {
			model = "llama3.2"
		}
		response, err = h.ollamaAPI.Complete(model, maxTokens, req.Prompt)

	default:
		http.Error(w, "Unsupported provider: "+req.Provider, http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Error in %s API: %v", req.Provider, err), http.StatusInternalServerError)
		return
	}

	result := UnifiedResponse{
		Response: response,
		Provider: req.Provider,
		Model:    model,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handlers) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	healthStatus := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"version":   ii.AppVersion,
		"apis": map[string]interface{}{
			"claude": map[string]interface{}{
				"configured": h.config.GetAPIKey("claude") != "",
				"available":  h.isAPIEnabled("claude"),
			},
			"openai": map[string]interface{}{
				"configured": h.config.GetAPIKey("openai") != "",
				"available":  h.isAPIEnabled("openai"),
			},
			"deepseek": map[string]interface{}{
				"configured": h.config.GetAPIKey("deepseek") != "",
				"available":  h.isAPIEnabled("deepseek"),
			},
			"ollama": map[string]interface{}{
				"configured": h.config.GetAPIKey("ollama") != "",
				"available":  h.isAPIEnabled("ollama"),
				// "endpoint":   h.config.OllamaEndpoint,
			},
		},
		"features": map[string]bool{
			"unified_api":     true,
			"model_selection": true,
			"cost_estimation": true,
		},
	}

	json.NewEncoder(w).Encode(healthStatus)
}

// HandleModels retrieves available models for a specific provider or all providers
func (h *Handlers) HandleModels(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	provider := r.URL.Query().Get("provider")

	var models []string
	var err error

	switch provider {
	case "openai":
		if h.isAPIEnabled("openai") {
			models, err = h.openaiAPI.ListModels()
			if err != nil {
				// Fallback to common models
				models = h.openaiAPI.GetCommonModels()
			}
		} else {
			models = h.openaiAPI.GetCommonModels()
		}

	case "deepseek":
		models = h.deepseekAPI.GetCommonModels()

	case "claude":
		models = []string{
			"claude-3-5-sonnet-20241022",
			"claude-3-5-haiku-20241022",
			"claude-3-opus-20240229",
			"claude-3-sonnet-20240229",
			"claude-3-haiku-20240307",
		}

	case "ollama":
		models = []string{
			"llama3.2",
			"llama3.1",
			"codellama",
			"mistral",
			"neural-chat",
			"vicuna",
			"wizardcoder",
			"llama2",
		}

	default:
		// Return all models
		allModels := map[string][]string{
			"openai":   h.openaiAPI.GetCommonModels(),
			"deepseek": h.deepseekAPI.GetCommonModels(),
			"claude": {
				"claude-3-5-sonnet-20241022",
				"claude-3-5-haiku-20241022",
				"claude-3-opus-20240229",
				"claude-3-sonnet-20240229",
				"claude-3-haiku-20240307",
			},
			"ollama": {
				"llama3.2", "llama3.1", "codellama", "mistral", "neural-chat",
				"vicuna", "wizardcoder", "llama2",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(allModels)
		return
	}

	result := map[string]interface{}{
		"provider": provider,
		"models":   models,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// HandleTest checks the availability of the specified provider
func (h *Handlers) HandleTest(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	provider := r.URL.Query().Get("provider")

	var available bool
	var message string

	switch provider {
	case "claude":
		available = h.config.GetAPIKey("claude") != ""
		if available {
			message = "Claude API configured"
		} else {
			message = "Claude API Key not configured"
		}

	case "openai":
		available = h.config.GetAPIKey("openai") != "" && h.openaiAPI.IsAvailable()
		if h.config.GetAPIKey("openai") == "" {
			message = "OpenAI API Key not configured"
		} else if !h.openaiAPI.IsAvailable() {
			message = "OpenAI API is not responding"
		} else {
			message = "OpenAI API is working"
		}

	case "deepseek":
		available = h.config.GetAPIKey("deepseek") != "" && h.deepseekAPI.IsAvailable()
		if h.config.GetAPIKey("deepseek") == "" {
			message = "DeepSeek API Key not configured"
		} else if !h.deepseekAPI.IsAvailable() {
			message = "DeepSeek API is not responding"
		} else {
			message = "DeepSeek API is working"
		}

	case "ollama":
		ollama := h.config.GetAPIConfig("ollama")
		available = h.ollamaAPI.IsAvailable()
		if available && ollama != nil {
			message = "Ollama is working at " + h.config.GetAPIEndpoint("ollama")
		} else {
			message = "Ollama is not responding at " + h.config.GetAPIEndpoint("ollama")
		}

	default:
		http.Error(w, "Provider not specified or invalid", http.StatusBadRequest)
		return
	}

	result := map[string]interface{}{
		"provider":  provider,
		"available": available,
		"message":   message,
		"timestamp": time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// setCORSHeaders sets the CORS headers for the response.
func (h *Handlers) setCORSHeaders(w http.ResponseWriter) {
	// CORS headers
	// These headers allow cross-origin requests from any domain
	// Adjust as necessary for your security requirements.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self'; connect-src 'self' https://api.openai.com https://api.deepseek.com https://api.ollama.com; frame-ancestors 'none'")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("X-Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self'; connect-src 'self' https://api.openai.com https://api.deepseek.com https://api.ollama.com; frame-ancestors 'none'")
	w.Header().Set("X-Content-Security-Policy-Report-Only", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self'; connect-src 'self' https://api.openai.com https://api.deepseek.com https://api.ollama.com; frame-ancestors 'none'")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
}

// HandleVersion returns the current application version
func (h *Handlers) HandleVersion(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	appVersion := ii.CurrentVersion
	if appVersion == "" {
		appVersion = ii.AppVersion
	}

	versionInfo := map[string]string{
		"version": appVersion,
		"name":    ii.AppName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versionInfo)
}

// HandleDocs returns the documentation URL
func (h *Handlers) HandleDocs(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	docs := map[string]string{
		"docs": "https://github.com/rafa-mori/grompt/blob/main/README.md",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}

// HandleSupport returns the support URL
func (h *Handlers) HandleSupport(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	supportInfo := map[string]string{
		"support": "https://github.com/rafa-mori/grompt/blob/main/README.md",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(supportInfo)
}

// HandleAbout returns information about the application
func (h *Handlers) HandleAbout(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	aboutInfo := map[string]string{
		"name":        "Grompt",
		"description": "A tool for building prompts with AI assistance using real engineering practices.",
		"version":     ii.AppVersion,
		"author":      "Rafa Mori",
		"license":     "MIT",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(aboutInfo)
}

// HandleStatus returns the current status of the application
func (h *Handlers) HandleStatus(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	status := map[string]string{
		"status":    "running",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   ii.AppVersion,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// HandleHelp returns the help URL or information
func (h *Handlers) HandleHelp(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	helpInfo := map[string]string{
		"help": "<HELP_URL>",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(helpInfo)
}

// HandleFeedback returns the feedback URL
func (h *Handlers) HandleFeedback(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	feedbackInfo := map[string]string{
		"feedback": "<FEEDBACK_URL>",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedbackInfo)
}

// HandleContact returns the contact URL
func (h *Handlers) HandleContact(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	contactInfo := map[string]string{
		"contact": "<CONTACT_URL>",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contactInfo)
}

// HandlePrivacy returns the privacy policy URL
func (h *Handlers) HandlePrivacy(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	privInfo := map[string]string{
		"privacy": "<PRIVACY_POLICY_URL>",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(privInfo)
}

// HandleTerms returns the terms of service URL
func (h *Handlers) HandleTerms(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	termsInfo := map[string]string{
		"terms": "<TERMS_OF_SERVICE_URL>",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(termsInfo)
}

// HandleRateLimit returns rate limit information
func (h *Handlers) HandleRateLimit(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	rateLimitInfo := map[string]string{
		"rate_limit": "<RATE_LIMIT_URL>",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rateLimitInfo)
}

// HandleError returns a generic error message
func (h *Handlers) HandleError(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	errorInfo := map[string]string{
		"error": "An unexpected error occurred. Please try again later.",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(errorInfo)
}

// HandleNotFound returns a 404 not found error
func (h *Handlers) HandleNotFound(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	notFoundInfo := map[string]string{
		"error": "Resource not found. Check the URL and try again.",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(notFoundInfo)
}

// HandleMethodNotAllowed returns a 405 method not allowed error_count
func (h *Handlers) HandleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	methodNotAllowedInfo := map[string]string{
		"error": "Method not allowed. Check the API documentation.",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(methodNotAllowedInfo)
}

// HandleInternalServerError returns a 500 internal server error
func (h *Handlers) HandleInternalServerError(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	errorInfo := map[string]string{
		"error": "An unexpected error occurred. Please try again later.",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(errorInfo)
}

// HandleBadRequest returns a 400 bad request error_count
func (h *Handlers) HandleBadRequest(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	badRequestInfo := map[string]string{
		"error": "Invalid request. Check the parameters and try again.",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(badRequestInfo)
}

// HadleUnauthorized returns a 401 unauthorized error_count
func (h *Handlers) HandleUnauthorized(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	unauthorizedInfo := map[string]string{
		"error": "Unauthorized access. Check your credentials.",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(unauthorizedInfo)
}

// HandleForbidden returns a 403 forbidden error_count
func (h *Handlers) HandleForbidden(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	forbiddenInfo := map[string]string{
		"error": "Access forbidden. You don't have permission to access this resource.",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(forbiddenInfo)
}
