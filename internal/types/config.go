// Package types defines the configuration and versioning for the Grompt application.
package types

import (
	"context"
	"time"

	vs "github.com/rafa-mori/grompt/version"
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

type Config struct {
	Port           string
	OpenAIAPIKey   string
	DeepSeekAPIKey string
	ClaudeAPIKey   string
	OllamaEndpoint string
	ChatGPTAPIKey  string
	Debug          bool
}

type APIConfig struct {
	OpenAIAvailable   bool   `json:"openai_available"`
	DeepSeekAvailable bool   `json:"deepseek_available"`
	ClaudeAvailable   bool   `json:"claude_available"`
	OllamaAvailable   bool   `json:"ollama_available"`
	ChatGPTAvailable  bool   `json:"chatgpt_available"`
	DemoMode          bool   `json:"demo_mode"`
	Version           string `json:"version"`
}

func (c *Config) GetAPIConfig() *APIConfig {
	return &APIConfig{
		OpenAIAvailable:   c.OpenAIAPIKey != "",
		DeepSeekAvailable: c.DeepSeekAPIKey != "",
		ClaudeAvailable:   c.ClaudeAPIKey != "",
		OllamaAvailable:   c.checkOllamaConnection(),
		ChatGPTAvailable:  c.ChatGPTAPIKey != "",
		DemoMode:          false,
		Version:           AppVersion,
	}
}

func (c *Config) checkOllamaConnection() bool {
	// Implementar verificação de conexão com Ollama
	// Por simplicidade, retorna false por enquanto
	return false
}
