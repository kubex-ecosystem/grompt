package types

import (
	"time"

	iit "github.com/rafa-mori/grompt/internal/core/provider"
)

// Config holds engine configuration
type Config struct {
	Port           string
	ClaudeAPIKey   string
	OpenAIAPIKey   string
	DeepSeekAPIKey string
	OllamaEndpoint string
	HistoryPath    string
	GeminiAPIKey   string
	ChatGPTAPIKey  string
	TemplatesPath  string
	Debug          bool
}

func (c *Config) GetAPIConfig(name string) iit.IAPIConfig {
	switch name {
	// case "claude":
	// 	return &APIConfig{APIKey: c.ClaudeAPIKey}
	// case "openai":
	// 	return &APIConfig{APIKey: c.OpenAIAPIKey}
	// case "deepseek":
	// 	return &APIConfig{APIKey: c.DeepSeekAPIKey}
	// case "ollama":
	// 	return &APIConfig{Endpoint: c.OllamaEndpoint}
	// case "gemini":
	// 	return &APIConfig{APIKey: c.GeminiAPIKey}
	// case "chatgpt":
	// 	return &APIConfig{APIKey: c.ChatGPTAPIKey}
	default:
		return nil
	}
}

func (c *Config) GetAPIEndpoint(name string) string {
	switch name {
	case "claude":
		return c.ClaudeAPIKey
	case "openai":
		return c.OpenAIAPIKey
	case "deepseek":
		return c.DeepSeekAPIKey
	case "ollama":
		return c.OllamaEndpoint
	case "gemini":
		return c.GeminiAPIKey
	case "chatgpt":
		return c.ChatGPTAPIKey
	default:
		return ""
	}
}

func (c *Config) GetAPIKey(name string) string {
	switch name {
	case "claude":
		return c.ClaudeAPIKey
	case "openai":
		return c.OpenAIAPIKey
	case "deepseek":
		return c.DeepSeekAPIKey
	case "ollama":
		return c.OllamaEndpoint
	case "gemini":
		return c.GeminiAPIKey
	case "chatgpt":
		return c.ChatGPTAPIKey
	default:
		return ""
	}
}

func (c *Config) GetBaseGenerationPrompt([]string, string, string, string, int) string {
	return "Base generation prompt"
}

func (c *Config) GetPort() string {
	return c.Port
}

func (c *Config) SetAPIKey(name, key string) error {
	switch name {
	case "claude":
		c.ClaudeAPIKey = key
	case "openai":
		c.OpenAIAPIKey = key
	case "deepseek":
		c.DeepSeekAPIKey = key
	case "ollama":
		c.OllamaEndpoint = key
	case "gemini":
		c.GeminiAPIKey = key
	case "chatgpt":
		c.ChatGPTAPIKey = key
	}

	return nil
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
