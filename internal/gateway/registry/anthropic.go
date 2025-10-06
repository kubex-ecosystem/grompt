package registry

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	providers "github.com/kubex-ecosystem/grompt/internal/types"
)

// anthropicProvider implements the Provider interface for Anthropic Claude
type anthropicProvider struct {
	name         string
	apiKey       string
	defaultModel string
	baseURL      string
	client       *http.Client
	mu           sync.Mutex
}

// NewAnthropicProvider creates a new Anthropic provider using REST API
func NewAnthropicProvider(name, baseURL, key, model string) (*anthropicProvider, error) {
	if key == "" {
		return nil, errors.New("API key is required for Anthropic provider")
	}
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}
	if model == "" {
		model = "claude-3-5-sonnet-20241022" // Latest Claude 3.5 Sonnet
	}

	return &anthropicProvider{
		name:         name,
		apiKey:       key,
		defaultModel: model,
		baseURL:      baseURL,
		client: &http.Client{
			Timeout: time.Minute * 5,
		},
	}, nil
}

func (p *anthropicProvider) Name() string {
	return p.name
}

func (p *anthropicProvider) Available() bool {
	if p.apiKey == "" {
		return false
	}
	return true
}

func (p *anthropicProvider) Execute(ctx context.Context, prompt string) (string, error) {
	// Implement command execution logic here if applicable
	return "", nil
}

func (p *anthropicProvider) GetCapabilities(ctx context.Context) *providers.Capabilities {
	if p == nil {
		return nil
	}

	models := []string{
		"claude-3-5-sonnet-20241022",
		"claude-3-haiku-20241022",
		"claude-3-opus-20241022",
	}

	return &providers.Capabilities{
		MaxTokens:         4096,
		SupportsBatch:     false,
		SupportsStreaming: true,
		Models:            models,
		Pricing: &providers.Pricing{
			InputCostPer1K:  0.003, // Default to Claude 3.5 Sonnet rates
			OutputCostPer1K: 0.015, // Default to Claude 3.5 Sonnet rates
			Currency:        "USD",
		},
		CanChat:       true,
		CanStream:     true,
		CanNotify:     false,
		CanExecute:    false,
		SupportsTools: false,
	}
}

// anthropicMessage represents a message in Anthropic's format
type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// anthropicRequest represents the request to Anthropic API
type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []anthropicMessage `json:"messages"`
	Stream    bool               `json:"stream"`
	System    string             `json:"system,omitempty"`
	Temp      float32            `json:"temperature,omitempty"`
}

// anthropicResponse represents the response from Anthropic API
type anthropicResponse struct {
	Type    string `json:"type"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model string `json:"model"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// anthropicStreamEvent represents a streaming event from Anthropic
type anthropicStreamEvent struct {
	Type  string `json:"type"`
	Delta struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"delta"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	Message struct {
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	} `json:"message"`
}

func (p *anthropicProvider) Chat(ctx context.Context, req providers.ChatRequest) (<-chan providers.ChatChunk, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Validate request
	if len(req.Messages) == 0 {
		return nil, errors.New("at least one message is required")
	}

	// Convert messages to Anthropic format
	messages := make([]anthropicMessage, 0, len(req.Messages))
	var systemMessage string

	for _, msg := range req.Messages {
		switch msg.Role {
		case "user", "assistant":
			messages = append(messages, anthropicMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		case "system":
			// Anthropic handles system messages separately
			systemMessage = msg.Content
		}
	}

	// Prepare request
	model := req.Model
	if model == "" {
		model = p.defaultModel
	}

	anthropicReq := anthropicRequest{
		Model:     model,
		MaxTokens: 4096,
		Messages:  messages,
		Stream:    true,
	}

	if systemMessage != "" {
		anthropicReq.System = systemMessage
	}

	if req.Temp > 0 {
		anthropicReq.Temp = req.Temp
	}

	// Create request body
	reqBody, err := json.Marshal(anthropicReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/v1/messages", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	httpReq.Header.Set("Accept", "text/event-stream")

	// Create response channel
	responseChan := make(chan providers.ChatChunk, 100)

	// Start streaming request in goroutine
	go func() {
		defer close(responseChan)

		startTime := time.Now()
		var totalTokens int
		var inputTokens int
		var outputTokens int

		// Make request
		resp, err := p.client.Do(httpReq)
		if err != nil {
			responseChan <- providers.ChatChunk{
				Content: "",
				Done:    true,
				Error:   fmt.Sprintf("HTTP request failed: %v", err),
			}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			responseChan <- providers.ChatChunk{
				Content: "",
				Done:    true,
				Error:   fmt.Sprintf("API error %d: %s", resp.StatusCode, string(body)),
			}
			return
		}

		// Handle streaming response
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()

			// Skip empty lines and non-data lines
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			// Remove "data: " prefix
			data := strings.TrimPrefix(line, "data: ")

			// Skip heartbeat
			if data == "[DONE]" {
				break
			}

			// Parse event
			var event anthropicStreamEvent
			if err := json.Unmarshal([]byte(data), &event); err != nil {
				continue // Skip invalid JSON
			}

			// Handle different event types
			switch event.Type {
			case "content_block_delta":
				if event.Delta.Type == "text_delta" {
					chunk := providers.ChatChunk{
						Content: event.Delta.Text,
						Done:    false,
					}

					select {
					case responseChan <- chunk:
					case <-ctx.Done():
						return
					}
				}

			case "message_start":
				if event.Message.Usage.InputTokens > 0 {
					inputTokens = event.Message.Usage.InputTokens
				}

			case "message_delta":
				if event.Usage.OutputTokens > 0 {
					outputTokens = event.Usage.OutputTokens
				}
			}
		}

		if err := scanner.Err(); err != nil {
			responseChan <- providers.ChatChunk{
				Content: "",
				Done:    true,
				Error:   fmt.Sprintf("Stream reading error: %v", err),
			}
			return
		}

		// Calculate final metrics
		totalTokens = inputTokens + outputTokens
		latencyMs := time.Since(startTime).Milliseconds()

		// Send final chunk with usage
		finalChunk := providers.ChatChunk{
			Content: "",
			Done:    true,
			Usage: &providers.Usage{
				Completion: outputTokens,
				Prompt:     inputTokens,
				Tokens:     totalTokens,
				Ms:         latencyMs,
				CostUSD:    calculateAnthropicCost(model, inputTokens, outputTokens),
				Provider:   p.name,
				Model:      model,
			},
		}

		select {
		case responseChan <- finalChunk:
		case <-ctx.Done():
		}

		// Log completion
		fmt.Printf("[Anthropic] Request completed - Model: %s, Tokens: %d, Duration: %v\n",
			model, totalTokens, time.Since(startTime))
	}()

	return responseChan, nil
}

func (p *anthropicProvider) Notify(ctx context.Context, event providers.NotificationEvent) error {
	// Implement notification logic here
	return nil
}

func (p *anthropicProvider) Close() error {
	// HTTP client doesn't require explicit cleanup
	return nil
}

// calculateAnthropicCost calculates the cost for Anthropic API usage
// Based on Claude pricing as of 2024
func calculateAnthropicCost(model string, inputTokens, outputTokens int) float64 {
	var inputRate, outputRate float64

	switch {
	case strings.Contains(model, "claude-3-5-sonnet"):
		inputRate = 3.0 / 1000000   // $3 per 1M input tokens
		outputRate = 15.0 / 1000000 // $15 per 1M output tokens
	case strings.Contains(model, "claude-3-haiku"):
		inputRate = 0.25 / 1000000  // $0.25 per 1M input tokens
		outputRate = 1.25 / 1000000 // $1.25 per 1M output tokens
	case strings.Contains(model, "claude-3-opus"):
		inputRate = 15.0 / 1000000  // $15 per 1M input tokens
		outputRate = 75.0 / 1000000 // $75 per 1M output tokens
	default:
		// Default to Sonnet pricing
		inputRate = 3.0 / 1000000
		outputRate = 15.0 / 1000000
	}

	return float64(inputTokens)*inputRate + float64(outputTokens)*outputRate
}
