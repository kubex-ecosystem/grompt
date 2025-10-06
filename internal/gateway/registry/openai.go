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
	"time"

	providers "github.com/kubex-ecosystem/grompt/internal/types"
)

// openaiProvider implements the Provider interface for OpenAI-compatible APIs
type openaiProvider struct {
	name         string
	baseURL      string
	apiKey       string
	defaultModel string
	client       *http.Client
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(name, baseURL, key, model string) (*openaiProvider, error) {
	if key == "" {
		return nil, errors.New("API key is required for OpenAI provider")
	}
	if baseURL == "" {
		baseURL = "https://api.openai.com"
	}

	return &openaiProvider{
		name:         name,
		baseURL:      baseURL,
		apiKey:       key,
		defaultModel: model,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// Name returns the provider name
func (o *openaiProvider) Name() string {
	return o.name
}

// Available checks if the provider is available
func (o *openaiProvider) Available() bool {
	return o.apiKey != ""
}

// Chat performs a chat completion request
func (o *openaiProvider) Chat(ctx context.Context, req providers.ChatRequest) (<-chan providers.ChatChunk, error) {
	model := req.Model
	if model == "" {
		model = o.defaultModel
	}

	body := map[string]interface{}{
		"model":       model,
		"messages":    toOpenAIMessages(req.Messages),
		"temperature": req.Temp,
		"stream":      true,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := o.baseURL + "/v1/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+o.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	ch := make(chan providers.ChatChunk, 8)

	go func() {
		defer close(ch)
		startTime := time.Now()

		resp, err := o.client.Do(httpReq)
		if err != nil {
			ch <- providers.ChatChunk{Done: true, Error: fmt.Sprintf("request failed: %v", err)}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			ch <- providers.ChatChunk{Done: true, Error: fmt.Sprintf("API error %d: %s", resp.StatusCode, string(body))}
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		totalTokens := 0

		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			var chunk openaiStreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue // Skip malformed chunks
			}

			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				ch <- providers.ChatChunk{Content: chunk.Choices[0].Delta.Content}
			}

			// Track token usage from usage field if present
			if chunk.Usage != nil {
				totalTokens = chunk.Usage.TotalTokens
			}
		}

		// Send final chunk with usage info
		latencyMs := time.Since(startTime).Milliseconds()
		ch <- providers.ChatChunk{
			Done: true,
			Usage: &providers.Usage{
				Tokens:   totalTokens,
				Ms:       latencyMs,
				CostUSD:  estimateCost(model, totalTokens), // Simple cost estimation
				Provider: o.name,
				Model:    model,
			},
		}
	}()

	return ch, nil
}

// Notify sends a notification event (not implemented)
func (o *openaiProvider) Notify(ctx context.Context, event providers.NotificationEvent) error {
	// Implement notification logic here
	return nil
}

func (o *openaiProvider) Execute(ctx context.Context, command string) (string, error) {
	// Implement command execution logic here if applicable
	return "", nil
}

func (o *openaiProvider) GetCapabilities(ctx context.Context) *providers.Capabilities {
	return &providers.Capabilities{
		SupportsBatch:     true,
		SupportsStreaming: true,
		Models: []string{
			o.defaultModel,
			"gpt-4",
			"gpt-4-turbo",
			"gpt-3.5-turbo",
			"gpt-3.5",
			"gpt-3.5-16k",
		},
		MaxTokens: 4096,
		Pricing: &providers.Pricing{
			InputCostPer1K:  0.0015,
			OutputCostPer1K: 0.002,
			Currency:        "USD",
		},
		CanChat:       true,
		CanStream:     true,
		CanNotify:     false,
		CanExecute:    false,
		SupportsTools: false,
	}
}

// toOpenAIMessages converts generic messages to OpenAI format
func toOpenAIMessages(messages []providers.Message) []map[string]string {
	result := make([]map[string]string, len(messages))
	for i, msg := range messages {
		result[i] = map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}
	return result
}

// openaiStreamChunk represents a streaming response chunk from OpenAI
type openaiStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
	Usage *struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage,omitempty"`
}

// estimateCost provides a rough cost estimation (simplified)
func estimateCost(model string, tokens int) float64 {
	// Simplified cost estimation - in production you'd want more accurate pricing
	costPerToken := 0.000002 // Default ~$2/1M tokens

	switch {
	case strings.Contains(model, "gpt-4"):
		costPerToken = 0.00003 // $30/1M tokens
	case strings.Contains(model, "gpt-3.5"):
		costPerToken = 0.000002 // $2/1M tokens
	}

	return float64(tokens) * costPerToken
}
