package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ClaudeAPI struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type ClaudeRequest struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type ClaudeAPIRequest struct {
	Model     string `json:"model"`
	MaxTokens int    `json:"max_tokens"`
	Messages  []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type ClaudeAPIResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func NewClaudeAPI(apiKey string) *ClaudeAPI {
	return &ClaudeAPI{
		apiKey:  apiKey,
		baseURL: "https://api.anthropic.com/v1/messages",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *ClaudeAPI) Complete(prompt string, maxTokens int) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("API key não configurada")
	}

	requestBody := ClaudeAPIRequest{
		Model:     "claude-3-sonnet-20240229",
		MaxTokens: maxTokens,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("erro ao serializar request: %v", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("erro ao criar request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro na requisição: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API retornou status %d: %s", resp.StatusCode, string(body))
	}

	var response ClaudeAPIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta: %v", err)
	}

	if len(response.Content) == 0 {
		return "", fmt.Errorf("resposta vazia da API")
	}

	return response.Content[0].Text, nil
}
