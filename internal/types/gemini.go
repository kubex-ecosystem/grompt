// Package types provides Google Gemini API implementation.
package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	gl "github.com/rafa-mori/grompt/internal/module/logger"
)

type GeminiAPI struct{ *APIConfig }

type GeminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
	GenerationConfig struct {
		MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
		Temperature     float64 `json:"temperature,omitempty"`
	} `json:"generationConfig,omitempty"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
}

type GeminiErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func NewGeminiAPI(apiKey string) IAPIConfig {
	return &GeminiAPI{
		APIConfig: &APIConfig{
			apiKey:  apiKey,
			baseURL: "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent",
			httpClient: &http.Client{
				Timeout: 60 * time.Second,
			},
		},
	}
}

// Complete sends a completion request to the Gemini API
func (g *GeminiAPI) Complete(prompt string, maxTokens int, model string) (string, error) {
	if g.apiKey == "" {
		gl.Log("debug", "Gemini API key not configured")
		return "", fmt.Errorf("gemini API key not configured")
	}

	// Define default model if not specified
	if model == "" {
		model = "gemini-1.5-flash"
	}

	// Update baseURL with model
	baseURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent", model)

	requestBody := GeminiRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: prompt},
				},
			},
		},
	}

	// Set generation config if maxTokens specified
	if maxTokens > 0 {
		requestBody.GenerationConfig.MaxOutputTokens = maxTokens
		requestBody.GenerationConfig.Temperature = 0.7
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		gl.Log("error", "Failed to serialize Gemini request: %v", err)
		return "", fmt.Errorf("error serializing request: %v", err)
	}

	// Create request with API key as query parameter
	requestURL := fmt.Sprintf("%s?key=%s", baseURL, g.apiKey)
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		gl.Log("error", "Failed to create Gemini request: %v", err)
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		gl.Log("error", "Gemini API request error: %v", err)
		return "", fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		gl.Log("error", "Failed to read Gemini response: %v", err)
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Try to parse Gemini error response
		var errorResp GeminiErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			gl.Log("error", "Gemini API error: %s (code: %d)", errorResp.Error.Message, errorResp.Error.Code)
			return "", fmt.Errorf("gemini API error: %s (code: %d)", errorResp.Error.Message, errorResp.Error.Code)
		}
		gl.Log("error", "API returned status %d: %s", resp.StatusCode, string(body))
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var response GeminiResponse
	if err := json.Unmarshal(body, &response); err != nil {
		gl.Log("error", "Failed to parse Gemini response: %v", err)
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		gl.Log("error", "No response generated from Gemini API")
		return "", fmt.Errorf("no response generated from Gemini API")
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}

// IsAvailable checks if the Gemini API is available
func (g *GeminiAPI) IsAvailable() bool {
	return true

	if g.apiKey == "" {
		gl.Log("notice", "Gemini API key not configured, assuming not available")
		return false
	}

	// Make a simple request to check API availability
	testReq := GeminiRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: "test"},
				},
			},
		},
		GenerationConfig: struct {
			MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
			Temperature     float64 `json:"temperature,omitempty"`
		}{
			MaxOutputTokens: 1,
		},
	}

	jsonData, err := json.Marshal(testReq)
	if err != nil {
		gl.Log("error", "Failed to serialize Gemini test request: %v", err)
		return false
	}

	requestURL := fmt.Sprintf("%s?key=%s", g.baseURL, g.apiKey)
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		gl.Log("error", "Failed to create Gemini test request: %v", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		gl.Log("error", "Gemini test request error: %v", err)
		return false
	}
	defer resp.Body.Close()

	gl.Log("debug", "Gemini API test request status: %d", resp.StatusCode)
	return resp.StatusCode == http.StatusOK
}

// GetCommonModels returns a list of common Gemini models
func (g *GeminiAPI) GetCommonModels() []string {
	return []string{
		//"gemini-2.5-flash",
		"gemini-1.5-flash",
		// "gemini-1.5-pro",
		// "gemini-1.0-pro",
	}
}

// ListModels returns available Gemini models
func (g *GeminiAPI) ListModels() ([]string, error) {
	// For now, return common models
	// In the future, could make API call to list available models
	return g.GetCommonModels(), nil
}

// GetVersion returns the API version
func (g *GeminiAPI) GetVersion() string {
	return "v1beta"
}

// IsDemoMode returns false as this is not demo mode
func (g *GeminiAPI) IsDemoMode() bool {
	return false
}
