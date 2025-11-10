package server

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	ii "github.com/kubex-ecosystem/grompt/internal/interfaces"
	it "github.com/kubex-ecosystem/grompt/internal/types"
)

type Handlers struct {
	config      ii.IConfig
	claudeAPI   ii.IAPIConfig
	openaiAPI   ii.IAPIConfig
	deepseekAPI ii.IAPIConfig
	chatGPTAPI  ii.IAPIConfig
	geminiAPI   ii.IAPIConfig
	ollamaAPI   ii.IAPIConfig
	// agentStore  *agents.Store
}

// UnifiedRequest represents a request structure that supports both direct prompts
// and prompt engineering from raw ideas. Either 'prompt' OR 'ideas' must be provided.
//
// Usage examples:
//   - Direct prompt: {"prompt": "Write a story about cats", "max_tokens": 500}
//   - Prompt engineering: {"ideas": ["cats", "adventure", "friendship"], "purpose": "Creative writing", "max_tokens": 500}
type UnifiedRequest struct {
	Lang        string   `json:"lang,omitempty"`         // Response language (default: "português")
	Purpose     string   `json:"purpose,omitempty"`      // Specific purpose description
	PurposeType string   `json:"purpose_type,omitempty"` // Type category (e.g., "Tutorial", "Creative writing")
	Ideas       []string `json:"ideas,omitempty"`        // Raw ideas for prompt engineering (alternative to prompt)
	Prompt      string   `json:"prompt,omitempty"`       // Direct prompt text (alternative to ideas)
	MaxTokens   int      `json:"max_tokens,omitempty"`   // Maximum response tokens
	Model       string   `json:"model,omitempty"`        // AI model to use
	Provider    string   `json:"provider,omitempty"`     // AI provider (for unified endpoint)
}

type UnifiedResponse struct {
	Response string     `json:"response"`
	Provider string     `json:"provider"`
	Model    string     `json:"model"`
	Mode     string     `json:"mode,omitempty"` // "byok", "server", or "demo"
	Usage    *UsageInfo `json:"usage,omitempty"`
}

type UsageInfo struct {
	PromptTokens     int     `json:"prompt_tokens,omitempty"`
	CompletionTokens int     `json:"completion_tokens,omitempty"`
	TotalTokens      int     `json:"total_tokens,omitempty"`
	EstimatedCost    float64 `json:"estimated_cost,omitempty"`
}

var llmKeyMap map[string]string

type providerDescriptor struct {
	Key              string
	DisplayName      string
	DefaultModel     string
	Models           []string
	SupportsBYOK     bool
	RequiresEndpoint bool
}

type providerSummary struct {
	Name         string   `json:"name"`
	DisplayName  string   `json:"display_name"`
	Available    bool     `json:"available"`
	Configured   bool     `json:"configured"`
	Models       []string `json:"models"`
	DefaultModel string   `json:"default_model,omitempty"`
	Endpoint     string   `json:"endpoint,omitempty"`
	Status       string   `json:"status"`
	Mode         string   `json:"mode"`
	SupportsBYOK bool     `json:"supports_byok"`
}

var providerCatalog = []providerDescriptor{
	{
		Key:          "openai",
		DisplayName:  "OpenAI",
		DefaultModel: "gpt-4o-mini",
		Models: []string{
			"gpt-4o",
			"gpt-4o-mini",
			"gpt-4-turbo",
			"gpt-3.5-turbo",
		},
		SupportsBYOK: true,
	},
	{
		Key:          "claude",
		DisplayName:  "Anthropic Claude",
		DefaultModel: "claude-3-5-sonnet-20241022",
		Models: []string{
			"claude-3-5-sonnet-20241022",
			"claude-3-sonnet-20240229",
			"claude-3-haiku-20240307",
		},
		SupportsBYOK: true,
	},
	{
		Key:          "gemini",
		DisplayName:  "Google Gemini",
		DefaultModel: "gemini-2.0-flash",
		Models: []string{
			"gemini-2.0-flash",
			"gemini-2.0-flash-exp",
			"gemini-2.5-pro",
		},
		SupportsBYOK: true,
	},
	{
		Key:          "deepseek",
		DisplayName:  "DeepSeek",
		DefaultModel: "deepseek-v3",
		Models: []string{
			"deepseek-v3",
			"deepseek-chat",
			"deepseek-reasoner",
		},
		SupportsBYOK: true,
	},
	{
		Key:          "chatgpt",
		DisplayName:  "ChatGPT (OpenAI)",
		DefaultModel: "gpt-4o-mini",
		Models: []string{
			"gpt-4o-mini",
			"gpt-4o",
			"gpt-4.1-mini",
		},
		SupportsBYOK: true,
	},
	{
		Key:              "ollama",
		DisplayName:      "Ollama (Local)",
		DefaultModel:     "llama3.2",
		Models:           []string{"llama3.2", "phi3", "mistral"},
		SupportsBYOK:     false,
		RequiresEndpoint: true,
	},
}

// generateDemoResponse creates a fallback demo response when no API keys are available
func (h *Handlers) generateDemoResponse(prompt string, purpose string) string {
	return fmt.Sprintf(`# 🎭 Demo Mode Response

**Original Request:** %s

## ⚠️ Notice
You are currently running Grompt in **DEMO MODE**. This is a simulated response demonstrating the application's capabilities.

## To Use Real AI Features

Configure an API key using one of these methods:

### Option 1: Server Configuration (Recommended)
Set environment variables before starting Grompt:
`+"```bash"+`
export OPENAI_API_KEY=sk-...
export CLAUDE_API_KEY=sk-ant-...
export GEMINI_API_KEY=AIza...
./grompt start
`+"```"+`

### Option 2: BYOK (Bring Your Own Key)
Use the 🔑 button in the interface to provide your API key per request.

## Supported Providers
- OpenAI (GPT-4, GPT-4o, GPT-3.5-turbo)
- Anthropic Claude (Claude 3.5 Sonnet, Haiku)
- Google Gemini (Gemini 2.0 Flash, Pro)
- DeepSeek (DeepSeek Chat, Coder)
- Ollama (Local models)

---

*This response was generated in demo mode. Configure an API key to access real AI-powered features.*
`, prompt)
}

func NewHandlers(cfg ii.IConfig) *Handlers {
	hndr := &Handlers{}
	if cfg == nil {
		return &Handlers{}
	} else {
		hndr.config = cfg
	}

	llmKeyMap = hndr.getLLMAPIKeyMap(cfg)
	hndr.claudeAPI = it.NewClaudeAPI(llmKeyMap["claude"])
	hndr.openaiAPI = it.NewOpenAIAPI(llmKeyMap["openai"])
	hndr.chatGPTAPI = it.NewChatGPTAPI(llmKeyMap["chatgpt"])
	hndr.deepseekAPI = it.NewDeepSeekAPI(llmKeyMap["deepseek"])
	hndr.ollamaAPI = it.NewOllamaAPI(llmKeyMap["ollama"])
	hndr.geminiAPI = it.NewGeminiAPI(llmKeyMap["gemini"])
	// hndr.agentStore = agents.NewStore("agents.json")

	return hndr
}

func (h *Handlers) HandleRoot(buildFS fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Bloqueia chamadas para /api/*
		p := strings.TrimPrefix(r.URL.Path, "/")
		if strings.HasPrefix(p, "api/") {
			http.NotFound(w, r)
			return
		}

		// 🧼 normaliza caminho e bloqueia traversal
		p = path.Clean(p)
		if p == "" || p == "/" || p == "." {
			p = "index.html"
		}

		if strings.Contains(p, "..") || !fs.ValidPath(p) {
			http.Error(w, "bad path", http.StatusBadRequest)
			return
		}

		if strings.Contains(p, "v1") && r.Method == http.MethodGet {
			http.NotFound(w, r)
			return
		} else if strings.Contains(p, "v1") && r.Method != http.MethodGet {
			// Define headers para evitar caching durante desenvolvimento
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
			w.Header().Set("Pragma", "no-cache")
		} else {
			// Define headers para caching de arquivos estáticos
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			expires := time.Now().Add(365 * 24 * time.Hour).Format(http.TimeFormat)
			w.Header().Set("Expires", expires)

			if strings.HasSuffix(p, "/") {
				p += "index.html"
			}
		}

		// Hide file from address bar
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", path.Base(p)))

		// tenta arquivo exato
		if f, err := buildFS.Open(p); err == nil {
			defer f.Close()
			if fi, _ := f.Stat(); fi != nil {
				if fi.IsDir() {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					fmt.Printf("⚙️  ServeFileFS para %s\n", p+"/index.html")
					http.ServeFileFS(w, r, buildFS, p+"/index.html")
					return
				} else {
					if strings.HasSuffix(p, ".txt") {
						p = strings.TrimSuffix(p, ".txt") + ".html"
					}

					// Define headers estáticos
					SetStaticHeaders(w, p)
					fmt.Printf("⚙️  ServeFileFS para %s\n", p)
					http.ServeFileFS(w, r, buildFS, p)
					return
				}
			}
		}

		// tenta index.html em subdiretório
		if f, err := buildFS.Open(path.Join(p, "index.html")); err == nil {
			defer f.Close()
			if fi, _ := f.Stat(); fi != nil {
				if !fi.IsDir() {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					fmt.Printf("⚙️  ServeFileFS para %s\n", p+"/index.html")
					http.ServeFileFS(w, r, buildFS, p)
					return
				}
			}
		}

		fmt.Printf("⚙️  ServeFileFS para %s\n", "index.html")
		http.ServeFileFS(w, r, buildFS, "index.html")
	}
}

func (h *Handlers) getLLMAPIKeyMap(cfg ii.IConfig) map[string]string {
	if llmKeyMap == nil {
		llmKeyMap = make(map[string]string)
	}
	if h.config == nil && cfg != nil {
		h.config = cfg
	}
	for _, provider := range []string{"claude", "openai", "chatgpt", "deepseek", "ollama", "gemini"} {
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
	case "gemini":
		return llmKeyMap["gemini"] != ""
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
	case "deepseek":
		return h.config.GetAPIKey("deepseek")
	case "ollama":
		return h.config.GetAPIKey("ollama")
	case "gemini":
		return h.config.GetAPIKey("gemini")
	case "chatgpt":
		return h.config.GetAPIKey("chatgpt")
	default:
		return ""
	}
}

func (h *Handlers) HandleConfig(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	if c.Request.Method == "POST" {
		// Handle POST request to update config
		var updateReq map[string]interface{}
		if err := json.NewDecoder(c.Request.Body).Decode(&updateReq); err != nil {
			http.Error(c.Writer, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Update configuration
		for key, value := range updateReq {
			if strValue, ok := value.(string); ok {
				switch key {
				case "gemini_api_key":
					h.config.SetAPIKey("gemini", strValue)
				case "openai_api_key":
					h.config.SetAPIKey("openai", strValue)
				case "claude_api_key":
					h.config.SetAPIKey("claude", strValue)
				case "deepseek_api_key":
					h.config.SetAPIKey("deepseek", strValue)
				case "chatgpt_api_key":
					h.config.SetAPIKey("chatgpt", strValue)
				case "ollama_endpoint":
					h.config.SetAPIKey("ollama", strValue)
				}
			}
		}

		// Refresh the API map
		h.getLLMAPIKeyMap(h.config)
	}

	providerDetails, providerOrder, readyProviders := h.describeProviders()
	demoMode := len(readyProviders) == 0
	defaultProvider := "openai"
	if len(readyProviders) > 0 {
		defaultProvider = readyProviders[0]
	} else if len(providerOrder) > 0 {
		defaultProvider = providerOrder[0]
	}

	serverInfo := map[string]string{
		"name":    "Grompt Server",
		"version": it.AppVersion,
		"port":    h.config.GetPort(),
		"status":  "ready",
	}
	if demoMode {
		serverInfo["status"] = "demo"
	}

	// Prepare response with both raw config and availability flags
	config := h.config
	response := map[string]interface{}{
		"server":              serverInfo,
		"environment":         map[string]bool{"demo_mode": demoMode},
		"providers":           providerDetails,
		"available_providers": providerOrder,
		"default_provider":    defaultProvider,
		"port":                h.config.GetPort(),
		"openai_api_key":      config.GetAPIKey("openai"),
		"deepseek_api_key":    config.GetAPIKey("deepseek"),
		"ollama_endpoint":     config.GetAPIEndpoint("ollama"),
		"claude_api_key":      config.GetAPIKey("claude"),
		"gemini_api_key":      config.GetAPIKey("gemini"),
		"chatgpt_api_key":     config.GetAPIKey("chatgpt"),
		"debug":               false,
		// Availability flags for frontend (legacy)
		"openai_available":   providerDetails["openai"].Available,
		"deepseek_available": providerDetails["deepseek"].Available,
		"ollama_available":   providerDetails["ollama"].Available,
		"claude_available":   providerDetails["claude"].Available,
		"gemini_available":   providerDetails["gemini"].Available,
		"chatgpt_available":  providerDetails["chatgpt"].Available,
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusOK, response)
}

func (h *Handlers) describeProviders() (map[string]providerSummary, []string, []string) {
	providers := make(map[string]providerSummary, len(providerCatalog))
	order := make([]string, 0, len(providerCatalog))
	ready := make([]string, 0, len(providerCatalog))

	for _, descriptor := range providerCatalog {
		order = append(order, descriptor.Key)

		var configured bool
		var endpoint string
		if descriptor.RequiresEndpoint {
			endpoint = h.config.GetAPIEndpoint(descriptor.Key)
			configured = strings.TrimSpace(endpoint) != ""
		} else {
			configured = strings.TrimSpace(h.config.GetAPIKey(descriptor.Key)) != ""
		}

		status := "needs_api_key"
		mode := "byok"
		available := configured

		if configured {
			status = "ready"
			mode = "server"
			ready = append(ready, descriptor.Key)
		} else if descriptor.RequiresEndpoint {
			status = "offline"
			mode = "offline"
		}

		providers[descriptor.Key] = providerSummary{
			Name:         descriptor.Key,
			DisplayName:  descriptor.DisplayName,
			Available:    available,
			Configured:   configured,
			Models:       descriptor.Models,
			DefaultModel: descriptor.DefaultModel,
			Endpoint:     endpoint,
			Status:       status,
			Mode:         mode,
			SupportsBYOK: descriptor.SupportsBYOK,
		}
	}

	if summary, ok := providers["claude"]; ok {
		alias := summary
		alias.Name = "anthropic"
		providers["anthropic"] = alias
	}

	return providers, order, ready
}

func (h *Handlers) HandleClaude(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	var req UnifiedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Validate input: must have either prompt or ideas
	if req.Prompt == "" && len(req.Ideas) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either 'prompt' or 'ideas' must be provided"})
		return
	}

	prompt := ""
	if req.Prompt != "" {
		prompt = req.Prompt
	} else {
		prompt = h.config.GetBaseGenerationPrompt(req.Ideas, req.Purpose, req.PurposeType, req.Lang, req.MaxTokens)
	}
	if prompt == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate prompt"})
		return
	}

	if key := h.config.GetAPIKey("claude"); key == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Claude API Key not configured"})
		return
	}

	response, err := h.claudeAPI.Complete(prompt, req.MaxTokens, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error in Claude API: %v", err)})
		return
	}

	result := UnifiedResponse{
		Response: response,
		Provider: "claude",
		Model:    "claude-3-5-sonnet-20241022",
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusOK, result)
}

func (h *Handlers) HandleOpenAI(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	var req UnifiedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	prompt := ""
	if req.Prompt != "" {
		prompt = req.Prompt
	} else {
		prompt = h.config.GetBaseGenerationPrompt(req.Ideas, req.Purpose, req.PurposeType, req.Lang, req.MaxTokens)
	}
	if prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prompt is required"})
		return
	}

	if key := h.config.GetAPIKey("openai"); key == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "OpenAI API Key not configured"})
		return
	}

	// Use default model if not specified
	model := req.Model
	if model == "" {
		model = "gpt-4o-mini"
	}

	response, err := h.openaiAPI.Complete(prompt, req.MaxTokens, model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error in OpenAI API: %v", err)})
		return
	}

	result := UnifiedResponse{
		Response: response,
		Provider: "openai",
		Model:    model,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusOK, result)
}

func (h *Handlers) HandleDeepSeek(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	var req UnifiedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	prompt := ""
	if req.Prompt != "" {
		prompt = req.Prompt
	} else {
		prompt = h.config.GetBaseGenerationPrompt(req.Ideas, req.Purpose, req.PurposeType, req.Lang, req.MaxTokens)
	}
	if prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prompt is required"})
		return
	}

	if key := h.config.GetAPIKey("deepseek"); key == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "DeepSeek API Key not configured"})
		return
	}

	// Use default model if not specified
	model := req.Model
	if model == "" {
		model = "deepseek-chat"
	}

	response, err := h.deepseekAPI.Complete(prompt, req.MaxTokens, model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error in DeepSeek API: %v", err)})
		return
	}

	result := UnifiedResponse{
		Response: response,
		Provider: "deepseek",
		Model:    model,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusOK, result)
}

func (h *Handlers) HandleGemini(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	var req UnifiedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	prompt := ""
	if req.Prompt != "" {
		prompt = req.Prompt
	} else {
		prompt = h.config.GetBaseGenerationPrompt(req.Ideas, req.Purpose, req.PurposeType, req.Lang, req.MaxTokens)
	}
	if prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prompt is required"})
		return
	}

	if key := h.config.GetAPIKey("gemini"); key == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Gemini API Key not configured"})
		return
	}

	// Use default model if not specified
	model := req.Model
	if model == "" {
		model = "gemini-2.0-flash"
	}

	response, err := h.geminiAPI.Complete(prompt, req.MaxTokens, model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error in Gemini API: %v", err)})
		return
	}

	result := UnifiedResponse{
		Response: response,
		Provider: "gemini",
		Model:    model,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusOK, result)
}

func (h *Handlers) HandleChatGPT(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	var req UnifiedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	prompt := ""
	if req.Prompt != "" {
		prompt = req.Prompt
	} else {
		prompt = h.config.GetBaseGenerationPrompt(req.Ideas, req.Purpose, req.PurposeType, req.Lang, req.MaxTokens)
	}
	if prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prompt is required"})
		return
	}

	if key := h.config.GetAPIKey("chatgpt"); key == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ChatGPT API Key not configured"})
		return
	}

	// Use default model if not specified
	model := req.Model
	if model == "" {
		model = "gpt-4o-mini"
	}

	response, err := h.chatGPTAPI.Complete(prompt, req.MaxTokens, model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error in ChatGPT API: %v", err)})
		return
	}

	result := UnifiedResponse{
		Response: response,
		Provider: "chatgpt",
		Model:    model,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusOK, result)
}

func (h *Handlers) HandleOllama(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	var req UnifiedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	prompt := ""
	if req.Prompt != "" {
		prompt = req.Prompt
	} else {
		prompt = h.config.GetBaseGenerationPrompt(req.Ideas, req.Purpose, req.PurposeType, req.Lang, req.MaxTokens)
	}
	if prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prompt is required"})
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

	response, err := h.ollamaAPI.Complete(model, maxTokens, prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error in Ollama API: %v", err)})
		return
	}

	result := UnifiedResponse{
		Response: response,
		Provider: "ollama",
		Model:    model,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusOK, result)
}

// HandleUnified processes requests for multiple providers in a unified manner
func (h *Handlers) HandleUnified(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	var req UnifiedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	prompt := ""
	if req.Prompt != "" {
		prompt = req.Prompt
	} else {
		prompt = h.config.GetBaseGenerationPrompt(req.Ideas, req.Purpose, req.PurposeType, req.Lang, req.MaxTokens)
	}
	if prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prompt is required"})
		return
	}

	// Validate provider
	if req.Provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider not specified"})
		return
	}

	// BYOK Support: Check for external API key in headers
	// Supports both generic X-API-Key and provider-specific X-{PROVIDER}-Key headers
	externalKey := c.Request.Header.Get("X-API-Key")
	if externalKey == "" {
		externalKey = c.Request.Header.Get("X-" + strings.ToUpper(req.Provider) + "-Key")
	}

	var response string
	var err error
	var model = req.Model
	var maxTokens = req.MaxTokens
	var mode = "server" // default to server mode

	// Determine API key source and mode
	var finalAPIKey string
	if externalKey != "" {
		finalAPIKey = externalKey
		mode = "byok"
	} else if h.config != nil && h.config.GetAPIKey(req.Provider) != "" {
		finalAPIKey = h.config.GetAPIKey(req.Provider)
		mode = "server"
	}

	switch req.Provider {
	case "claude":
		if finalAPIKey == "" {
			// Demo mode fallback
			mode = "demo"
			if model == "" {
				model = "claude-3-5-sonnet-20241022"
			}
			response = h.generateDemoResponse(prompt, req.Purpose)
		} else {
			// Real API call
			api := h.claudeAPI
			if mode == "byok" {
				api = it.NewClaudeAPI(finalAPIKey)
			}
			if model == "" {
				model = "claude-3-5-sonnet-20241022"
			}
			response, err = api.Complete(prompt, maxTokens, model)
		}

	case "openai":
		if finalAPIKey == "" {
			// Demo mode fallback
			mode = "demo"
			if model == "" {
				model = "gpt-4o-mini"
			}
			response = h.generateDemoResponse(prompt, req.Purpose)
		} else {
			// Real API call
			api := h.openaiAPI
			if mode == "byok" {
				api = it.NewOpenAIAPI(finalAPIKey)
			}
			if model == "" {
				model = "gpt-4o-mini"
			}
			response, err = api.Complete(prompt, req.MaxTokens, model)
		}

	case "deepseek":
		if finalAPIKey == "" {
			// Demo mode fallback
			mode = "demo"
			if model == "" {
				model = "deepseek-chat"
			}
			response = h.generateDemoResponse(prompt, req.Purpose)
		} else {
			// Real API call
			api := h.deepseekAPI
			if mode == "byok" {
				api = it.NewDeepSeekAPI(finalAPIKey)
			}
			if model == "" {
				model = "deepseek-chat"
			}
			response, err = api.Complete(prompt, req.MaxTokens, model)
		}

	case "gemini":
		if finalAPIKey == "" {
			// Demo mode fallback
			mode = "demo"
			if model == "" {
				model = "gemini-2.0-flash-exp"
			}
			response = h.generateDemoResponse(prompt, req.Purpose)
		} else {
			// Real API call
			api := h.geminiAPI
			if mode == "byok" {
				api = it.NewGeminiAPI(finalAPIKey)
			}
			if model == "" {
				model = "gemini-2.0-flash-exp"
			}
			response, err = api.Complete(prompt, req.MaxTokens, model)
		}

	case "chatgpt":
		if finalAPIKey == "" {
			// Demo mode fallback
			mode = "demo"
			if model == "" {
				model = "gpt-4o-mini"
			}
			response = h.generateDemoResponse(prompt, req.Purpose)
		} else {
			// Real API call
			api := h.chatGPTAPI
			if mode == "byok" {
				api = it.NewChatGPTAPI(finalAPIKey)
			}
			if model == "" {
				model = "gpt-4o-mini"
			}
			response, err = api.Complete(prompt, req.MaxTokens, model)
		}

	case "ollama":
		// Ollama doesn't require API key (local instance)
		if model == "" {
			model = "llama3.2"
		}
		response, err = h.ollamaAPI.Complete(model, maxTokens, prompt)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported provider: " + req.Provider})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error in %s API: %v", req.Provider, err)})
		return
	}

	result := UnifiedResponse{
		Response: response,
		Provider: req.Provider,
		Model:    model,
		Mode:     mode, // Include mode in response
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusOK, result)
}

// HandleAsk provides a simpler alias to ask a direct question to a provider.
// Request JSON:
//
//	{
//	  "question": "...",            // required
//	  "provider": "openai|claude|...", // optional (auto-pick if omitted)
//	  "model": "gpt-4o-mini",       // optional
//	  "max_tokens": 1000             // optional
//	}
//
// Response JSON matches UnifiedResponse: { response, provider, model }
func (h *Handlers) HandleAsk(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	var req struct {
		Question  string `json:"question"`
		Provider  string `json:"provider,omitempty"`
		Model     string `json:"model,omitempty"`
		MaxTokens int    `json:"max_tokens,omitempty"`
		Lang      string `json:"lang,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	if strings.TrimSpace(req.Question) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'question' is required"})
		return
	}

	prompt := req.Question
	provider := strings.TrimSpace(req.Provider)
	model := strings.TrimSpace(req.Model)
	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 1000
	}

	var response string
	var err error

	// If provider not specified, pick the first available in a sensible order
	if provider == "" {
		for _, p := range []string{"openai", "claude", "deepseek", "ollama", "gemini", "chatgpt"} {
			if h.isAPIEnabled(p) {
				provider = p
				break
			}
		}
		if provider == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "No AI provider configured"})
			return
		}
	}

	switch provider {
	case "openai":
		if h.config.GetAPIKey("openai") == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "OpenAI API Key not configured"})
			return
		}
		if model == "" {
			model = "gpt-4o-mini"
		}
		response, err = h.openaiAPI.Complete(prompt, maxTokens, model)
	case "claude":
		if h.config.GetAPIKey("claude") == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Claude API Key not configured"})
			return
		}
		if model == "" {
			model = "claude-3-5-sonnet-20241022"
		}
		response, err = h.claudeAPI.Complete(prompt, maxTokens, model)
	case "deepseek":
		if h.config.GetAPIKey("deepseek") == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "DeepSeek API Key not configured"})
			return
		}
		if model == "" {
			model = "deepseek-chat"
		}
		response, err = h.deepseekAPI.Complete(prompt, maxTokens, model)
	case "ollama":
		if model == "" {
			model = "llama3.2"
		}
		response, err = h.ollamaAPI.Complete(model, maxTokens, prompt)
	case "gemini":
		if h.config.GetAPIKey("gemini") == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Gemini API Key not configured"})
			return
		}
		if model == "" {
			model = "gemini-1.5-flash"
		}
		response, err = h.geminiAPI.Complete(prompt, maxTokens, model)
	case "chatgpt":
		if h.config.GetAPIKey("chatgpt") == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ChatGPT API Key not configured"})
			return
		}
		if model == "" {
			model = "gpt-4o-mini"
		}
		response, err = h.chatGPTAPI.Complete(prompt, maxTokens, model)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported provider: " + provider})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error in %s API: %v", provider, err)})
		return
	}

	res := UnifiedResponse{Response: response, Provider: provider, Model: model}
	c.JSON(http.StatusOK, res)
}

// HandleSquad is a stable alias that generates a squad of agents from a description,
// delegating to the existing agents generation logic.
// Request JSON:
// { "description": "..." } or { "requirements": "..." }
// Response JSON: { agents: [...], markdown: "..." }
func (h *Handlers) HandleSquad(c *gin.Context) {
	h.setCORSHeaders(c)
	if c.Request.Method == "OPTIONS" {
		return
	}
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	var req struct {
		Description  string `json:"description"`
		Requirements string `json:"requirements"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	content := strings.TrimSpace(req.Requirements)
	if content == "" {
		content = strings.TrimSpace(req.Description)
	}
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'description' or 'requirements' is required"})
		return
	}

	// Reuse HandleAgentsGenerate inner logic via local function
	// Build LLM function
	// llmFunc := func(prompt string) (string, error) {
	// 	if h.config.GetAPIKey("claude") != "" {
	// 		return h.claudeAPI.Complete(prompt, 4000, "claude-2")
	// 	}
	// 	if h.config.GetAPIKey("openai") != "" {
	// 		return h.openaiAPI.Complete(prompt, 4000, "gpt-4")
	// 	}
	// 	if h.config.GetAPIKey("deepseek") != "" {
	// 		return h.deepseekAPI.Complete(prompt, 4000, "deepseek-chat")
	// 	}
	// 	return "", fmt.Errorf("no LLM API available")
	// }

	// Use the same squad package as /api/v1/agents/generate
	// sqAgents, err := squad.ParseRequirementsWithLLM(content, llmFunc)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// markdown := squad.GenerateMarkdown(sqAgents)

	// resp := struct {
	// 	Agents   []squad.Agent `json:"agents"`
	// 	Markdown string        `json:"markdown"`
	// }{Agents: sqAgents, Markdown: markdown}

	c.Header("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(resp)
	c.JSON(http.StatusOK, gin.H{"agents": []string{}, "markdown": "Feature under development."})
}

func (h *Handlers) HandleHealth(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	healthStatus := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"version":   it.AppVersion,
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

	c.JSON(http.StatusOK, healthStatus)
}

// HandleModels retrieves available models for a specific provider or all providers
func (h *Handlers) HandleModels(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	provider := c.Query("provider")

	var models map[string]any
	var err error

	switch provider {
	case "openai":
		if h.isAPIEnabled("openai") {
			models, err = h.openaiAPI.ListModels()
			if err != nil {
				// Fallback to common models
				models = map[string]any{
					"common": h.openaiAPI.GetCommonModels(),
				}
			}
		} else {
			models = map[string]any{
				"common": h.openaiAPI.GetCommonModels(),
			}
		}

	case "deepseek":
		models = map[string]any{
			"common": h.deepseekAPI.GetCommonModels(),
		}

	case "claude":
		models = map[string]any{
			"claude-3-5-sonnet-20241022": struct{}{},
			"claude-3-5-haiku-20241022":  struct{}{},
			"claude-3-opus-20240229":     struct{}{},
			"claude-3-sonnet-20240229":   struct{}{},
			"claude-3-haiku-20240307":    struct{}{},
		}

	case "ollama":
		models = map[string]any{
			"llama3.2":    struct{}{},
			"llama3.1":    struct{}{},
			"codellama":   struct{}{},
			"mistral":     struct{}{},
			"neural-chat": struct{}{},
			"vicuna":      struct{}{},
			"wizardcoder": struct{}{},
			"llama2":      struct{}{},
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

		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, allModels)
		return
	}

	result := map[string]interface{}{
		"provider": provider,
		"models":   models,
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, result)
}

// HandleTest checks the availability of the specified provider
func (h *Handlers) HandleTest(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	provider := c.Query("provider")

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider not specified or invalid"})
		return
	}

	result := map[string]interface{}{
		"provider":  provider,
		"available": available,
		"message":   message,
		"timestamp": time.Now().Unix(),
	}

	c.JSON(http.StatusOK, result)
}

// setCORSHeaders sets the CORS headers for the response.
func (h *Handlers) setCORSHeaders(c *gin.Context) {
	// CORS headers
	// These headers allow cross-origin requests from any domain
	// Adjust as necessary for your security requirements.
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	// BYOK Support: Allow custom API key headers
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key, X-OPENAI-Key, X-CLAUDE-Key, X-GEMINI-Key, X-DEEPSEEK-Key, X-CHATGPT-Key")
	c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self'; connect-src 'self' https://api.openai.com https://api.deepseek.com https://api.ollama.com; frame-ancestors 'none'")
	c.Writer.Header().Set("Referrer-Policy", "no-referrer")
	c.Writer.Header().Set("X-Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self'; connect-src 'self' https://api.openai.com https://api.deepseek.com https://api.ollama.com; frame-ancestors 'none'")
	c.Writer.Header().Set("X-Content-Security-Policy-Report-Only", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self'; connect-src 'self' https://api.openai.com https://api.deepseek.com https://api.ollama.com; frame-ancestors 'none'")
	c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
	c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	c.Writer.Header().Set("X-Frame-Options", "DENY")
}

// HandleVersion returns the current application version
func (h *Handlers) HandleVersion(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	appVersion := it.CurrentVersion
	if appVersion == "" {
		appVersion = it.AppVersion
	}

	versionInfo := map[string]string{
		"version": appVersion,
		"name":    it.AppName,
	}

	c.JSON(http.StatusOK, versionInfo)
}

// HandleDocs returns the documentation URL
func (h *Handlers) HandleDocs(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	docs := map[string]string{
		"docs": "https://github.com/kubex-ecosystem/grompt/blob/main/README.md",
	}

	c.JSON(http.StatusOK, docs)
}

// HandleSupport returns the support URL
func (h *Handlers) HandleSupport(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	supportInfo := map[string]string{
		"support": "https://github.com/kubex-ecosystem/grompt/blob/main/README.md",
	}

	c.JSON(http.StatusOK, supportInfo)
}

// HandleAbout returns information about the application
func (h *Handlers) HandleAbout(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	aboutInfo := map[string]string{
		"name":        "Grompt",
		"description": "A tool for building prompts with AI assistance using real engineering practices.",
		"version":     it.AppVersion,
		"author":      "Rafa Mori",
		"license":     "MIT",
	}

	c.JSON(http.StatusOK, aboutInfo)
}

// HandleStatus returns the current status of the application
func (h *Handlers) HandleStatus(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	status := map[string]string{
		"status":    "running",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   it.AppVersion,
	}

	c.JSON(http.StatusOK, status)
}

// HandleHelp returns the help URL or information
func (h *Handlers) HandleHelp(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	helpInfo := map[string]string{
		"help": "<HELP_URL>",
	}

	c.JSON(http.StatusOK, helpInfo)
}

// HandleFeedback returns the feedback URL
func (h *Handlers) HandleFeedback(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	feedbackInfo := map[string]string{
		"feedback": "<FEEDBACK_URL>",
	}

	c.JSON(http.StatusOK, feedbackInfo)
}

// HandleContact returns the contact URL
func (h *Handlers) HandleContact(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	contactInfo := map[string]string{
		"contact": "<CONTACT_URL>",
	}

	c.JSON(http.StatusOK, contactInfo)
}

// HandlePrivacy returns the privacy policy URL
func (h *Handlers) HandlePrivacy(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	privInfo := map[string]string{
		"privacy": "<PRIVACY_POLICY_URL>",
	}

	c.JSON(http.StatusOK, privInfo)
}

// HandleTerms returns the terms of service URL
func (h *Handlers) HandleTerms(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	termsInfo := map[string]string{
		"terms": "<TERMS_OF_SERVICE_URL>",
	}

	c.JSON(http.StatusOK, termsInfo)
}

// HandleRateLimit returns rate limit information
func (h *Handlers) HandleRateLimit(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	rateLimitInfo := map[string]string{
		"rate_limit": "<RATE_LIMIT_URL>",
	}

	c.JSON(http.StatusOK, rateLimitInfo)
}

// HandleError returns a generic error message
func (h *Handlers) HandleError(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	errorInfo := map[string]string{
		"error": "An unexpected error occurred. Please try again later.",
	}

	c.JSON(http.StatusInternalServerError, errorInfo)
}

// HandleNotFound returns a 404 not found error
func (h *Handlers) HandleNotFound(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	notFoundInfo := map[string]string{
		"error": "Resource not found. Check the URL and try again.",
	}

	c.JSON(http.StatusNotFound, notFoundInfo)
}

// HandleMethodNotAllowed returns a 405 method not allowed error_count
func (h *Handlers) HandleMethodNotAllowed(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	methodNotAllowedInfo := map[string]string{
		"error": "Method not allowed. Check the API documentation.",
	}

	c.JSON(http.StatusMethodNotAllowed, methodNotAllowedInfo)
}

// HandleInternalServerError returns a 500 internal server error
func (h *Handlers) HandleInternalServerError(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	errorInfo := map[string]string{
		"error": "An unexpected error occurred. Please try again later.",
	}

	c.JSON(http.StatusInternalServerError, errorInfo)
}

// HandleBadRequest returns a 400 bad request error
func (h *Handlers) HandleBadRequest(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	badRequestInfo := map[string]string{
		"error": "Invalid request. Check the parameters and try again.",
	}

	c.JSON(http.StatusBadRequest, badRequestInfo)
}

// HandleUnauthorized returns a 401 unauthorized error
func (h *Handlers) HandleUnauthorized(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	unauthorizedInfo := map[string]string{
		"error": "Unauthorized access. Check your credentials.",
	}

	c.JSON(http.StatusUnauthorized, unauthorizedInfo)
}

// HandleForbidden returns a 403 forbidden error
func (h *Handlers) HandleForbidden(c *gin.Context) {
	h.setCORSHeaders(c)

	if c.Request.Method == "OPTIONS" {
		return
	}

	forbiddenInfo := map[string]string{
		"error": "Access forbidden. You don't have permission to access this resource.",
	}

	c.JSON(http.StatusForbidden, forbiddenInfo)
}
