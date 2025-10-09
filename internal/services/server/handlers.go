package server

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/services/agents"
	"github.com/kubex-ecosystem/grompt/internal/services/squad"
	ii "github.com/kubex-ecosystem/grompt/internal/types"
)

type Handlers struct {
	config      ii.IConfig
	claudeAPI   ii.IAPIConfig
	openaiAPI   ii.IAPIConfig
	deepseekAPI ii.IAPIConfig
	chatGPTAPI  ii.IAPIConfig
	geminiAPI   ii.IAPIConfig
	ollamaAPI   ii.IAPIConfig
	agentStore  *agents.Store
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
` + "```bash" + `
export OPENAI_API_KEY=sk-...
export CLAUDE_API_KEY=sk-ant-...
export GEMINI_API_KEY=AIza...
./grompt start
` + "```" + `

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
	hndr.claudeAPI = ii.NewClaudeAPI(llmKeyMap["claude"])
	hndr.openaiAPI = ii.NewOpenAIAPI(llmKeyMap["openai"])
	hndr.chatGPTAPI = ii.NewChatGPTAPI(llmKeyMap["chatgpt"])
	hndr.deepseekAPI = ii.NewDeepSeekAPI(llmKeyMap["deepseek"])
	hndr.ollamaAPI = ii.NewOllamaAPI(llmKeyMap["ollama"])
	hndr.geminiAPI = ii.NewGeminiAPI(llmKeyMap["gemini"])
	hndr.agentStore = agents.NewStore("agents.json")

	return hndr
}

func (h *Handlers) HandleRoot(buildFS fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		// Hide file from address bar
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", path.Base(p)))

		// tenta arquivo exato
		if f, err := buildFS.Open(p); err == nil {
			defer f.Close()
			if fi, _ := f.Stat(); fi != nil {
				if fi.IsDir() {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
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
					w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
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

func (h *Handlers) HandleConfig(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST" {
		// Handle POST request to update config
		var updateReq map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
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

	// Prepare response with both raw config and availability flags
	config := h.config
	response := map[string]interface{}{
		"port":             "8080",
		"openai_api_key":   config.GetAPIKey("openai"),
		"deepseek_api_key": config.GetAPIKey("deepseek"),
		"ollama_endpoint":  config.GetAPIKey("ollama"),
		"claude_api_key":   config.GetAPIKey("claude"),
		"gemini_api_key":   config.GetAPIKey("gemini"),
		"chatgpt_api_key":  config.GetAPIKey("chatgpt"),
		"debug":            false,
		// Availability flags for frontend
		"openai_available":   h.isAPIEnabled("openai"),
		"deepseek_available": h.isAPIEnabled("deepseek"),
		"ollama_available":   h.isAPIEnabled("ollama"),
		"claude_available":   h.isAPIEnabled("claude"),
		"gemini_available":   h.isAPIEnabled("gemini"),
		"chatgpt_available":  h.isAPIEnabled("chatgpt"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

	// Validate input: must have either prompt or ideas
	if req.Prompt == "" && len(req.Ideas) == 0 {
		http.Error(w, "Either 'prompt' or 'ideas' must be provided", http.StatusBadRequest)
		return
	}

	prompt := ""
	if req.Prompt != "" {
		prompt = req.Prompt
	} else {
		prompt = h.config.GetBaseGenerationPrompt(req.Ideas, req.Purpose, req.PurposeType, req.Lang, req.MaxTokens)
	}
	if prompt == "" {
		http.Error(w, "Failed to generate prompt", http.StatusInternalServerError)
		return
	}

	if key := h.config.GetAPIKey("claude"); key == "" {
		http.Error(w, "Claude API Key not configured", http.StatusServiceUnavailable)
		return
	}

	response, err := h.claudeAPI.Complete(prompt, req.MaxTokens, "")
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

	prompt := ""
	if req.Prompt != "" {
		prompt = req.Prompt
	} else {
		prompt = h.config.GetBaseGenerationPrompt(req.Ideas, req.Purpose, req.PurposeType, req.Lang, req.MaxTokens)
	}
	if prompt == "" {
		http.Error(w, "Prompt is required", http.StatusBadRequest)
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

	response, err := h.openaiAPI.Complete(prompt, req.MaxTokens, model)
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

	prompt := ""
	if req.Prompt != "" {
		prompt = req.Prompt
	} else {
		prompt = h.config.GetBaseGenerationPrompt(req.Ideas, req.Purpose, req.PurposeType, req.Lang, req.MaxTokens)
	}
	if prompt == "" {
		http.Error(w, "Prompt is required", http.StatusBadRequest)
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

	response, err := h.deepseekAPI.Complete(prompt, req.MaxTokens, model)
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

func (h *Handlers) HandleGemini(w http.ResponseWriter, r *http.Request) {
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

	prompt := ""
	if req.Prompt != "" {
		prompt = req.Prompt
	} else {
		prompt = h.config.GetBaseGenerationPrompt(req.Ideas, req.Purpose, req.PurposeType, req.Lang, req.MaxTokens)
	}
	if prompt == "" {
		http.Error(w, "Prompt is required", http.StatusBadRequest)
		return
	}

	if key := h.config.GetAPIKey("gemini"); key == "" {
		http.Error(w, "Gemini API Key not configured", http.StatusServiceUnavailable)
		return
	}

	// Use default model if not specified
	model := req.Model
	if model == "" {
		model = "gemini-2.0-flash"
	}

	response, err := h.geminiAPI.Complete(prompt, req.MaxTokens, model)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in Gemini API: %v", err), http.StatusInternalServerError)
		return
	}

	result := UnifiedResponse{
		Response: response,
		Provider: "gemini",
		Model:    model,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handlers) HandleChatGPT(w http.ResponseWriter, r *http.Request) {
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

	prompt := ""
	if req.Prompt != "" {
		prompt = req.Prompt
	} else {
		prompt = h.config.GetBaseGenerationPrompt(req.Ideas, req.Purpose, req.PurposeType, req.Lang, req.MaxTokens)
	}
	if prompt == "" {
		http.Error(w, "Prompt is required", http.StatusBadRequest)
		return
	}

	if key := h.config.GetAPIKey("chatgpt"); key == "" {
		http.Error(w, "ChatGPT API Key not configured", http.StatusServiceUnavailable)
		return
	}

	// Use default model if not specified
	model := req.Model
	if model == "" {
		model = "gpt-4o-mini"
	}

	response, err := h.chatGPTAPI.Complete(prompt, req.MaxTokens, model)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in ChatGPT API: %v", err), http.StatusInternalServerError)
		return
	}

	result := UnifiedResponse{
		Response: response,
		Provider: "chatgpt",
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

	prompt := ""
	if req.Prompt != "" {
		prompt = req.Prompt
	} else {
		prompt = h.config.GetBaseGenerationPrompt(req.Ideas, req.Purpose, req.PurposeType, req.Lang, req.MaxTokens)
	}
	if prompt == "" {
		http.Error(w, "Prompt is required", http.StatusBadRequest)
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

	prompt := ""
	if req.Prompt != "" {
		prompt = req.Prompt
	} else {
		prompt = h.config.GetBaseGenerationPrompt(req.Ideas, req.Purpose, req.PurposeType, req.Lang, req.MaxTokens)
	}
	if prompt == "" {
		http.Error(w, "Prompt is required", http.StatusBadRequest)
		return
	}

	// Validate provider
	if req.Provider == "" {
		http.Error(w, "Provider not specified", http.StatusBadRequest)
		return
	}

	// BYOK Support: Check for external API key in headers
	// Supports both generic X-API-Key and provider-specific X-{PROVIDER}-Key headers
	externalKey := r.Header.Get("X-API-Key")
	if externalKey == "" {
		externalKey = r.Header.Get("X-" + strings.ToUpper(req.Provider) + "-Key")
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
				api = ii.NewClaudeAPI(finalAPIKey)
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
				api = ii.NewOpenAIAPI(finalAPIKey)
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
				api = ii.NewDeepSeekAPI(finalAPIKey)
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
				api = ii.NewGeminiAPI(finalAPIKey)
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
				api = ii.NewChatGPTAPI(finalAPIKey)
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
		Mode:     mode, // Include mode in response
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
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
func (h *Handlers) HandleAsk(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Question  string `json:"question"`
		Provider  string `json:"provider,omitempty"`
		Model     string `json:"model,omitempty"`
		MaxTokens int    `json:"max_tokens,omitempty"`
		Lang      string `json:"lang,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Question) == "" {
		http.Error(w, "'question' is required", http.StatusBadRequest)
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
			http.Error(w, "No AI provider configured", http.StatusServiceUnavailable)
			return
		}
	}

	switch provider {
	case "openai":
		if h.config.GetAPIKey("openai") == "" {
			http.Error(w, "OpenAI API Key not configured", http.StatusServiceUnavailable)
			return
		}
		if model == "" {
			model = "gpt-4o-mini"
		}
		response, err = h.openaiAPI.Complete(prompt, maxTokens, model)
	case "claude":
		if h.config.GetAPIKey("claude") == "" {
			http.Error(w, "Claude API Key not configured", http.StatusServiceUnavailable)
			return
		}
		if model == "" {
			model = "claude-3-5-sonnet-20241022"
		}
		response, err = h.claudeAPI.Complete(prompt, maxTokens, model)
	case "deepseek":
		if h.config.GetAPIKey("deepseek") == "" {
			http.Error(w, "DeepSeek API Key not configured", http.StatusServiceUnavailable)
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
			http.Error(w, "Gemini API Key not configured", http.StatusServiceUnavailable)
			return
		}
		if model == "" {
			model = "gemini-1.5-flash"
		}
		response, err = h.geminiAPI.Complete(prompt, maxTokens, model)
	case "chatgpt":
		if h.config.GetAPIKey("chatgpt") == "" {
			http.Error(w, "ChatGPT API Key not configured", http.StatusServiceUnavailable)
			return
		}
		if model == "" {
			model = "gpt-4o-mini"
		}
		response, err = h.chatGPTAPI.Complete(prompt, maxTokens, model)
	default:
		http.Error(w, "Unsupported provider: "+provider, http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Error in %s API: %v", provider, err), http.StatusInternalServerError)
		return
	}

	res := UnifiedResponse{Response: response, Provider: provider, Model: model}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// HandleSquad is a stable alias that generates a squad of agents from a description,
// delegating to the existing agents generation logic.
// Request JSON:
// { "description": "..." } or { "requirements": "..." }
// Response JSON: { agents: [...], markdown: "..." }
func (h *Handlers) HandleSquad(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Description  string `json:"description"`
		Requirements string `json:"requirements"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	content := strings.TrimSpace(req.Requirements)
	if content == "" {
		content = strings.TrimSpace(req.Description)
	}
	if content == "" {
		http.Error(w, "'description' or 'requirements' is required", http.StatusBadRequest)
		return
	}

	// Reuse HandleAgentsGenerate inner logic via local function
	// Build LLM function
	llmFunc := func(prompt string) (string, error) {
		if h.config.GetAPIKey("claude") != "" {
			return h.claudeAPI.Complete(prompt, 4000, "claude-2")
		}
		if h.config.GetAPIKey("openai") != "" {
			return h.openaiAPI.Complete(prompt, 4000, "gpt-4")
		}
		if h.config.GetAPIKey("deepseek") != "" {
			return h.deepseekAPI.Complete(prompt, 4000, "deepseek-chat")
		}
		return "", fmt.Errorf("no LLM API available")
	}

	// Use the same squad package as /api/agents/generate
	sqAgents, err := squad.ParseRequirementsWithLLM(content, llmFunc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	markdown := squad.GenerateMarkdown(sqAgents)

	resp := struct {
		Agents   []squad.Agent `json:"agents"`
		Markdown string        `json:"markdown"`
	}{Agents: sqAgents, Markdown: markdown}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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
	// BYOK Support: Allow custom API key headers
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key, X-OPENAI-Key, X-CLAUDE-Key, X-GEMINI-Key, X-DEEPSEEK-Key, X-CHATGPT-Key")
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
		"docs": "https://github.com/kubex-ecosystem/grompt/blob/main/README.md",
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
		"support": "https://github.com/kubex-ecosystem/grompt/blob/main/README.md",
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

// HandleUnauthorized returns a 401 unauthorized error
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
