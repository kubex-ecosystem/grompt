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

// groqProvider implements the Provider interface for Groq's ultra-fast inference
type groqProvider struct {
	name         string
	apiKey       string
	defaultModel string
	baseURL      string
	client       *http.Client
	mu           sync.Mutex
}

// NewGroqProvider creates a new Groq provider for lightning-fast inference
func NewGroqProvider(name, baseURL, key, model string) (*groqProvider, error) {
	if key == "" {
		return nil, errors.New("API key is required for Groq provider")
	}
	if baseURL == "" {
		baseURL = "https://api.groq.com"
	}
	if model == "" {
		model = "llama-3.1-70b-versatile" // Default to Llama 3.1 70B
	}

	return &groqProvider{
		name:         name,
		apiKey:       key,
		defaultModel: model,
		baseURL:      baseURL,
		client: &http.Client{
			Timeout: time.Minute * 2, // Groq is so fast we can use shorter timeout
		},
	}, nil
}

func (p *groqProvider) Name() string {
	return p.name
}

func (p *groqProvider) Available() error {
	if p.apiKey == "" {
		return errors.New("groq API key not configured")
	}
	return nil
}

func (p *groqProvider) Notify(ctx context.Context, event providers.NotificationEvent) error {
	// Implement notification logic here
	return nil
}

// groqRequest represents the request to Groq API (OpenAI-compatible)
type groqRequest struct {
	Model       string        `json:"model"`
	Messages    []groqMessage `json:"messages"`
	Stream      bool          `json:"stream"`
	Temperature *float32      `json:"temperature,omitempty"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
	TopP        *float32      `json:"top_p,omitempty"`
}

// groqMessage represents a message in Groq's format (OpenAI-compatible)
type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// groqResponse represents the response from Groq API
type groqResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// groqStreamChunk represents a streaming chunk from Groq
type groqStreamChunk struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage,omitempty"`
}

func (p *groqProvider) Chat(ctx context.Context, req providers.ChatRequest) (<-chan providers.ChatChunk, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Validate request
	if len(req.Messages) == 0 {
		return nil, errors.New("at least one message is required")
	}

	// Convert messages to Groq format (same as OpenAI)
	messages := make([]groqMessage, 0, len(req.Messages))
	for _, msg := range req.Messages {
		messages = append(messages, groqMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Prepare request
	model := req.Model
	if model == "" {
		model = p.defaultModel
	}

	groqReq := groqRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
	}

	// Add optional parameters
	if req.Temp > 0 {
		groqReq.Temperature = &req.Temp
	}

	// Groq supports high token limits
	maxTokens := 8192
	groqReq.MaxTokens = &maxTokens

	// Create request body
	reqBody, err := json.Marshal(groqReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/openai/v1/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers (OpenAI-compatible)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
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
				Error:   fmt.Sprintf("Groq API error %d: %s", resp.StatusCode, string(body)),
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

			// Check for end of stream
			if data == "[DONE]" {
				break
			}

			// Parse chunk
			var chunk groqStreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue // Skip invalid JSON
			}

			// Process choices
			if len(chunk.Choices) > 0 {
				choice := chunk.Choices[0]

				// Send content chunk
				if choice.Delta.Content != "" {
					responseChunk := providers.ChatChunk{
						Content: choice.Delta.Content,
						Done:    false,
					}

					select {
					case responseChan <- responseChunk:
					case <-ctx.Done():
						return
					}
				}

				// Handle completion
				if choice.FinishReason != nil && *choice.FinishReason != "" {
					// This is the final chunk, extract usage if available
					if chunk.Usage != nil {
						inputTokens = chunk.Usage.PromptTokens
						outputTokens = chunk.Usage.CompletionTokens
						totalTokens = chunk.Usage.TotalTokens
					}
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
				CostUSD:    calculateGroqCost(model, inputTokens, outputTokens),
				Provider:   p.name,
				Model:      model,
			},
		}

		select {
		case responseChan <- finalChunk:
		case <-ctx.Done():
		}

		// Log completion with speed info
		tokensPerSecond := float64(totalTokens) / (float64(latencyMs) / 1000.0)
		fmt.Printf("[Groq] ⚡ LIGHTNING FAST! Model: %s, Tokens: %d, Duration: %v, Speed: %.1f tok/s\n",
			model, totalTokens, time.Since(startTime), tokensPerSecond)
	}()

	return responseChan, nil
}

func (p *groqProvider) Close() error {
	// HTTP client doesn't require explicit cleanup
	return nil
}

// calculateGroqCost calculates the cost for Groq API usage
// Groq has very competitive pricing, especially for open-source models
func calculateGroqCost(model string, inputTokens, outputTokens int) float64 {
	var inputRate, outputRate float64

	switch {
	case strings.Contains(model, "llama-3.1-70b"):
		inputRate = 0.59 / 1000000  // $0.59 per 1M input tokens
		outputRate = 0.79 / 1000000 // $0.79 per 1M output tokens
	case strings.Contains(model, "llama-3.1-8b"):
		inputRate = 0.05 / 1000000  // $0.05 per 1M input tokens
		outputRate = 0.08 / 1000000 // $0.08 per 1M output tokens
	case strings.Contains(model, "mixtral-8x7b"):
		inputRate = 0.24 / 1000000  // $0.24 per 1M input tokens
		outputRate = 0.24 / 1000000 // $0.24 per 1M output tokens
	case strings.Contains(model, "gemma"):
		inputRate = 0.10 / 1000000  // $0.10 per 1M input tokens
		outputRate = 0.10 / 1000000 // $0.10 per 1M output tokens
	default:
		// Default to Llama 3.1 70B pricing
		inputRate = 0.59 / 1000000
		outputRate = 0.79 / 1000000
	}

	return float64(inputTokens)*inputRate + float64(outputTokens)*outputRate
}
