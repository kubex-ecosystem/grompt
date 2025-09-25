// Package types defines interfaces and types for AI providers
package types

import "context"

type ServerCORS struct {
	AllowOrigins []string `yaml:"allow_origins"`
}

type ServerConfig struct {
	Addr  string     `yaml:"addr"`
	CORS  ServerCORS `yaml:"cors"`
	Debug bool       `yaml:"debug"`
}

type DefaultsConfig struct {
	TenantID                   string    `yaml:"tenant_id"`
	UserID                     string    `yaml:"user_id"`
	Byok                       string    `yaml:"byok"`
	NotificationProvider       *Provider `yaml:"notification_provider"`
	NotificationTimeoutSeconds int       `yaml:"notification_timeout_seconds"`
}

type ToolCall struct {
	Name string      `json:"name"`
	Args interface{} `json:"args"` // geralmente map[string]any
}

// ProviderConfig holds configuration for a specific provider
type ProviderConfig struct {
	BaseURL      string `yaml:"base_url"`
	KeyEnv       string `yaml:"key_env"`
	DefaultModel string `yaml:"default_model"`
	Type         string `yaml:"type"` // "openai", "anthropic", "groq", "openrouter", "ollama"
}

// Config holds the complete provider configuration
type Config struct {
	Server    *ServerConfig             `yaml:"server"`
	Defaults  *DefaultsConfig           `yaml:"defaults"`
	Providers map[string]ProviderConfig `yaml:"providers"`
}

// Provider interface defines the contract for AI providers
type Provider interface {
	Name() string
	Chat(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error)
	Available() error
	Notify(ctx context.Context, event NotificationEvent) error
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Headers  map[string]string `json:"-"`
	Provider string            `json:"provider"`
	Model    string            `json:"model"`
	Messages []Message         `json:"messages"`
	Temp     float32           `json:"temperature"`
	Stream   bool              `json:"stream"`
	Meta     map[string]any    `json:"meta"`
}

// Message represents a single chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Usage represents token usage and cost information
type Usage struct {
	Completion int     `json:"completion_tokens"`
	Prompt     int     `json:"prompt_tokens"`
	Tokens     int     `json:"tokens"`
	Ms         int64   `json:"latency_ms"`
	CostUSD    float64 `json:"cost_usd"`
	Provider   string  `json:"provider"`
	Model      string  `json:"model"`
}

// ChatChunk represents a streaming response chunk
type ChatChunk struct {
	Content  string    `json:"content,omitempty"`
	Done     bool      `json:"done"`
	Usage    *Usage    `json:"usage,omitempty"`
	Error    string    `json:"error,omitempty"`
	ToolCall *ToolCall `json:"toolCall,omitempty"`
}
