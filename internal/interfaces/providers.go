package interfaces

import "context"

// Provider interface defines the contract for AI providers
type Provider interface {
	Name() string
	Chat(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error)
	KeyEnv() string
	Type() string
	IsAvailable() bool
	Version() string
	Execute(ctx context.Context, template string, vars map[string]any) (*Result, error)
	Notify(ctx context.Context, event NotificationEvent) error
	GetCapabilities(ctx context.Context) *Capabilities
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Headers  map[string]string `json:"-"`
	Provider string            `json:"provider"`
	Model    string            `json:"model"`
	PromptTemplate string          `json:"prompt_template"`
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
	// ToolCall *ToolCall `json:"toolCall,omitempty"`
}

// Capabilities describes what a provider can do
type Capabilities struct {
	MaxTokens         int      `json:"max_tokens"`
	SupportsBatch     bool     `json:"supports_batch"`
	SupportsStreaming bool     `json:"supports_streaming"`
	Models            map[string]any `json:"models"`
	Pricing           *Pricing `json:"pricing,omitempty"`
}

// Pricing information for the provider
type Pricing struct {
	InputCostPer1K  float64 `json:"input_cost_per_1k"`
	OutputCostPer1K float64 `json:"output_cost_per_1k"`
	Currency        string  `json:"currency"`
}
