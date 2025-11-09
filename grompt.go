// Package grompt provides a simple and flexible way to create interactive command-line prompts in Go.
package grompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kubex-ecosystem/grompt/internal/interfaces"
	itypes "github.com/kubex-ecosystem/grompt/internal/types"
	"github.com/kubex-ecosystem/grompt/types"
	logz "github.com/kubex-ecosystem/logz"
)

// Result represents the outcome of a processed prompt.
type Result = interfaces.Result

// Capabilities describes provider abilities (legacy compatibility).
type Capabilities = interfaces.Capabilities

// Pricing captures basic usage pricing metadata.
type Pricing = interfaces.Pricing

// PromptEngine exposes legacy prompt-processing capabilities.
type PromptEngine = types.PromptEngine

// Provider is the legacy provider interface exposed to consumers.
type Provider = interfaces.Provider

// APIConfig mirrors the legacy engine API configuration contract.
type IAPIConfig = interfaces.IAPIConfig
type APIConfig = itypes.APIConfig

// Config mirrors the legacy engine configuration contract.
type Config = types.Config

// Prompt represents a command-line prompt with a message and an optional default value.
type Prompt struct {
	Message      string
	DefaultValue string
}

// ---------- Public constructors ----------

// DefaultConfig rebuilds a legacy-compatible configuration.
func DefaultConfig(configFilePath string) Config {
	return types.DefaultConfig(configFilePath)
}

// NewConfig constructs a configuration using explicit parameters.
func NewConfig(
	bindAddr string,
	port string,
	openAIKey string,
	deepSeekKey string,
	ollamaEndpoint string,
	claudeKey string,
	geminiKey string,
	chatGPTKey string,
	logger logz.Logger,
) Config {
	return types.NewConfig(
		bindAddr,
		port,
		openAIKey,
		deepSeekKey,
		ollamaEndpoint,
		claudeKey,
		geminiKey,
		chatGPTKey,
		logger,
	)
}

// NewPromptEngine returns a legacy-compatible engine backed by the new gateway stack.
func NewPromptEngine(cfg Config) PromptEngine {
	return types.NewPromptEngine(cfg)
}

// NewPrompt creates a new Prompt instance with the given message and default value.
func NewPrompt(message, defaultValue string) *Prompt {
	return &Prompt{
		Message:      message,
		DefaultValue: defaultValue,
	}
}

// Show displays the prompt to the user and captures their input.
// If the user provides no input, the default value is returned.
func (p *Prompt) Show() string {
	reader := bufio.NewReader(os.Stdin)
	if p.DefaultValue != "" {
		fmt.Printf("%s [%s]: ", p.Message, p.DefaultValue)
	} else {
		fmt.Printf("%s: ", p.Message)
	}
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return p.DefaultValue
	}
	return input
}
