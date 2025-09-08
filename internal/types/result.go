package engine

import "time"

// Config holds engine configuration
type Config struct {
	Port           string
	ClaudeAPIKey   string
	OpenAIAPIKey   string
	DeepSeekAPIKey string
	OllamaEndpoint string
	HistoryPath    string
	TemplatesPath  string
}

// Result represents the result of prompt processing
type Result struct {
	ID        string                 `json:"id"`
	Prompt    string                 `json:"prompt"`
	Response  string                 `json:"response"`
	Provider  string                 `json:"provider"`
	Variables map[string]interface{} `json:"variables,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}
