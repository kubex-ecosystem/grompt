// Package types defines the configuration and versioning for the Grompt application.
package types

import (
	"context"
	"fmt"
	"net/http"
	"time"

	vs "github.com/rafa-mori/grompt/internal/module/version"
)

var (
	CurrentVersion string    = vs.GetVersion()
	LatestVersion  string    = vs.GetLatestVersionFromGit()
	LastCheckTime  time.Time = time.Now()
)

func init() {
	// Initialize the CurrentVersion and LatestVersion
	// This will run in a goroutine to avoid blocking the main execution
	// and will check for the latest version from Git if needed.
	// If the CurrentVersion is not set, it will use the version from the version package
	// and will update the LatestVersion if the last check was more than 24 hours ago
	// or if it is the first run.
	ctx := context.Background()
	cancel, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	// Use a goroutine to avoid blocking the main execution
	go func(cancel context.Context) {
		select {
		case <-cancel.Done():
			return // Exit if the context is cancelled
		default:
			if CurrentVersion == "" {
				CurrentVersion = vs.GetVersion()
			}
			if LastCheckTime.IsZero() || LastCheckTime.Before(time.Now().Add(-24*time.Hour)) {
				// Check for the latest version from Git
				LatestVersion = vs.GetLatestVersionFromGit()
				if LatestVersion == "" {
					LatestVersion = CurrentVersion // Fallback to current version if check fails
				}
				LastCheckTime = time.Now()
			}
		}
	}(cancel)

	// // Ensure that the CurrentVersion is always set
	// if CurrentVersion == "" {
	// 	CurrentVersion = vs.GetVersion()
	// }
	// if LatestVersion == "" {
	// 	LatestVersion = vs.GetLatestVersionFromGit()
	// }
}

const (
	AppName     = "Grompt"
	AppVersion  = "1.0.0"
	DefaultPort = "8080"
)

type IAPIConfig interface {
	IsAvailable() bool
	IsDemoMode() bool
	GetVersion() string
	ListModels() ([]string, error)
	GetCommonModels() []string
	Complete(prompt string, maxTokens int, model string) (string, error)
}

type APIConfig struct {
	apiKey     string
	baseURL    string
	version    string
	httpClient *http.Client
	demoMode   bool
}

type IConfig interface {
	GetAPIConfig(provider string) IAPIConfig
	GetPort() string
	GetAPIKey(provider string) string
	SetAPIKey(provider string, key string) error
	GetAPIEndpoint(provider string) string
}

type Config struct {
	Port           string `json:"port" gorm:"default:8080"`
	OpenAIAPIKey   string `json:"openai_api_key" gorm:"default:''"`
	DeepSeekAPIKey string `json:"deepseek_api_key" gorm:"default:''"`
	OllamaEndpoint string `json:"ollama_endpoint" gorm:"default:'http://localhost:11434'"`
	ClaudeAPIKey   string `json:"claude_api_key" gorm:"default:''"`
	GeminiAPIKey   string `json:"gemini_api_key" gorm:"default:''"`
	ChatGPTAPIKey  string `json:"chatgpt_api_key" gorm:"default:''"`
	Debug          bool   `json:"debug" gorm:"default:false"`
}

func NewConfig(port, openAIKey, deepSeekKey, ollamaEndpoint, claudeKey, geminiKey string) *Config {
	return &Config{
		Port:           port,
		OpenAIAPIKey:   openAIKey,
		DeepSeekAPIKey: deepSeekKey,
		OllamaEndpoint: ollamaEndpoint,
		ClaudeAPIKey:   claudeKey,
		GeminiAPIKey:   geminiKey,
	}
}

func (c *Config) GetAPIConfig(provider string) IAPIConfig {
	switch provider {
	case "openai":
		return NewOpenAIAPI(c.OpenAIAPIKey)
	case "deepseek":
		return NewDeepSeekAPI(c.DeepSeekAPIKey)
	case "ollama":
		return NewOllamaAPI(c.OllamaEndpoint)
	case "claude":
		return NewClaudeAPI(c.ClaudeAPIKey)
	case "gemini":
		return NewGeminiAPI(c.GeminiAPIKey)
	default:
		return nil
	}
}

func (c *Config) GetPort() string {
	if c.Port == "" {
		return DefaultPort
	}
	return c.Port
}

func (c *Config) GetAPIKey(provider string) string {
	switch provider {
	case "openai":
		return c.OpenAIAPIKey
	case "deepseek":
		return c.DeepSeekAPIKey
	case "ollama":
		return c.OllamaEndpoint
	case "claude":
		return c.ClaudeAPIKey
	case "gemini":
		return c.GeminiAPIKey
	default:
		return ""
	}
}

func (c *Config) SetAPIKey(provider string, key string) error {
	switch provider {
	case "openai":
		c.OpenAIAPIKey = key
	case "deepseek":
		c.DeepSeekAPIKey = key
	case "ollama":
		c.OllamaEndpoint = key
	case "claude":
		c.ClaudeAPIKey = key
	case "gemini":
		c.GeminiAPIKey = key
	default:
		return fmt.Errorf("unknown provider: %s", provider)
	}
	return nil
}

func (c *Config) GetAPIEndpoint(provider string) string {
	if provider == "ollama" {
		return c.OllamaEndpoint
	}
	return ""
}

func (c *Config) checkOllamaConnection() bool {
	// Implementar verificação de conexão com Ollama
	// Por simplicidade, retorna false por enquanto
	return false
}
