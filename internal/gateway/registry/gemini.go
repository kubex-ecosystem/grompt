package registry

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	genai "cloud.google.com/go/vertexai/genai"
	providers "github.com/kubex-ecosystem/grompt/internal/types"
	"google.golang.org/api/option"
)

// geminiProvider implements the Provider interface for Google Gemini
type geminiProvider struct {
	name         string
	apiKey       string
	defaultModel string
	baseURL      string
	client       *genai.Client
	mu           sync.Mutex
}

// NewGeminiProvider creates a new Gemini provider using the SDK
func NewGeminiProvider(name, baseURL, key, model string) (*geminiProvider, error) {
	if key == "" {
		return nil, errors.New("API key is required for Gemini provider")
	}
	if model == "" {
		model = "gemini-1.5-flash"
	}

	// Create a client for the entire provider instance
	ctx := context.Background()
	client, err := genai.NewClient(ctx, "587138832075", "southamerica-east1", option.WithAPIKey(key))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &geminiProvider{
		name:         name,
		apiKey:       key,
		defaultModel: model,
		baseURL:      baseURL,
		client:       client,
	}, nil
}

// Name returns the provider name
func (g *geminiProvider) Name() string {
	return g.name
}

// Available checks if the provider is available
func (g *geminiProvider) Available() bool {
	if g.apiKey == "" {
		return false
	}
	return true
}

// Chat performs a chat completion request using Gemini's streaming API with the SDK
func (g *geminiProvider) Chat(ctx context.Context, req providers.ChatRequest) (<-chan providers.ChatChunk, error) {
	model := req.Model
	if model == "" {
		model = g.defaultModel
	}

	// Create a new model instance for each request to set specific parameters
	geminiModel := g.client.GenerativeModel(model)

	// Set generation configurations
	geminiModel.SetTemperature(float32(req.Temp))
	geminiModel.SetMaxOutputTokens(int32(8192))
	geminiModel.SetTopP(0.95)

	var safetySettings []*genai.SafetySetting
	var schema *genai.Schema
	var candidateCount = func() *int32 { var i int32 = 1; return &i }()

	// Set safety settings
	// safetySettings = append(safetySettings, &genai.SafetySetting{
	// 	Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockMediumAndAbove,
	// })
	// safetySettings = append(safetySettings, &genai.SafetySetting{
	// 	Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockMediumAndAbove,
	// })
	// safetySettings = append(safetySettings, &genai.SafetySetting{
	// 	Category: genai.HarmCategorySexual, Threshold: genai.HarmBlockMediumAndAbove,
	// })
	// safetySettings = append(safetySettings, &genai.SafetySetting{
	// 	Category: genai.HarmCategoryToxicity, Threshold: genai.HarmBlockLowAndAbove,
	// })
	geminiModel.SafetySettings = safetySettings

	// Convert messages to Gemini SDK format - CREATE PARTS FOR STREAMING!
	var parts []genai.Part

	// Handle special analysis requests (your genius feature!)
	if analysisType, ok := req.Meta["analysisType"]; ok {
		if projectContext, hasContext := req.Meta["projectContext"]; hasContext {
			prompt := g.getAnalysisPrompt(projectContext.(string), analysisType.(string), req.Meta)
			parts = append(parts, genai.Text(prompt))

			// Configure for analysis
			geminiModel.SetTemperature(0.3)

			// Add response schema if structured output is requested
			if req.Meta["useStructuredOutput"] == true {

				schema = &genai.Schema{
					Type:       genai.TypeObject,
					Properties: make(map[string]*genai.Schema),
					Required:   []string{"projectName", "summary", "strengths", "weaknesses", "recommendations"},
				}

				// Define properties based on analysis type
				geminiModel.CandidateCount = candidateCount
				geminiModel.ResponseMIMEType = "application/json"
				geminiModel.ResponseSchema = schema
			}
		}
	} else {
		// Normal chat - convert messages to parts
		for _, msg := range req.Messages {
			if msg.Content != "" {
				parts = append(parts, genai.Text(msg.Content))
			}
		}
	}

	// Validation: ensure we have content to send
	if len(parts) == 0 {
		return nil, fmt.Errorf("no valid content to send to Gemini")
	}

	ch := make(chan providers.ChatChunk, 8)

	go func() {
		defer close(ch)
		startTime := time.Now()

		// Call the SDK's streaming method with PARTS not CONTENTS!
		iter := geminiModel.GenerateContentStream(ctx, parts...)

		totalTokens := 0
		var fullContent strings.Builder

		// Iterate through streaming response
		for {
			resp, err := iter.Next()
			if err != nil {
				// Check if iteration is complete - this is NORMAL end of stream
				if strings.Contains(err.Error(), "done") ||
					strings.Contains(err.Error(), "EOF") ||
					strings.Contains(err.Error(), "no more items") {
					break // Normal completion, not an error
				}
				ch <- providers.ChatChunk{Done: true, Error: fmt.Sprintf("streaming error: %v", err)}
				return
			}

			if resp == nil {
				continue
			}

			// Extract content from response
			if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
				for _, part := range resp.Candidates[0].Content.Parts {
					if text, ok := part.(genai.Text); ok {
						chunk := string(text)
						ch <- providers.ChatChunk{Content: chunk}
						fullContent.WriteString(chunk)
					}
				}
			}

			// Extract usage metadata when available
			if resp.UsageMetadata != nil {
				totalTokens = int(resp.UsageMetadata.PromptTokenCount + resp.UsageMetadata.CandidatesTokenCount)
			}
		}

		// If no usage metadata was provided, estimate tokens from the final text
		if totalTokens == 0 {
			totalTokens = g.estimateTokens(fullContent.String())
		}

		// Send final chunk with metrics
		latencyMs := time.Since(startTime).Milliseconds()
		ch <- providers.ChatChunk{
			Done: true,
			Usage: &providers.Usage{
				Tokens:   totalTokens,
				Ms:       latencyMs,
				CostUSD:  g.estimateCost(model, totalTokens),
				Provider: g.name,
				Model:    model,
			},
		}
	}()

	return ch, nil
}

func (g *geminiProvider) Notify(ctx context.Context, event providers.NotificationEvent) error {
	// Implement notification logic here
	return nil
}

func (g *geminiProvider) Execute(ctx context.Context, prompt string) (string, error) {
	// Implement command execution logic here if applicable
	return "", nil
}

func (g *geminiProvider) GetCapabilities(ctx context.Context) *providers.Capabilities {
	return &providers.Capabilities{
		SupportsBatch:     true,
		SupportsStreaming: true,
		Models: []string{
			"gemini-1",
			"gemini-1.5",
			"gemini-2",
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

// toGeminiContents converts generic messages to Gemini SDK format
func (g *geminiProvider) toGeminiContents(messages []providers.Message) []genai.Part {
	contents := make([]genai.Part, 0, len(messages))

	for _, msg := range messages {
		if msg.Content == "" {
			continue
		}

		//role := "user"
		// if msg.Role == "assistant" || msg.Role == "model" {
		// 	role = "model"
		// }

		contents = append(contents, genai.Text(msg.Content))
	}
	return contents
}

// getAnalysisPrompt generates analysis prompts (your original logic, cleaned up)
func (g *geminiProvider) getAnalysisPrompt(projectContext, analysisType string, meta map[string]interface{}) string {
	locale := "en-US"
	if l, ok := meta["locale"]; ok {
		if localeStr, ok := l.(string); ok {
			locale = localeStr
		}
	}
	language := "English (US)"
	if locale == "pt-BR" {
		language = "Portuguese (Brazil)"
	}
	return fmt.Sprintf(`You are a world-class senior software architect and project management consultant with 20 years of experience.

**Task:** Analyze the following software project based on the provided context.
**Analysis Type:** %s
**Response Language:** %s

**Project Context:**
%s

**Instructions:**
- Provide detailed, actionable insights
- Focus on practical recommendations
- Structure your response clearly
- Be specific and concrete in your suggestions

Analyze thoroughly and provide valuable insights.`, analysisType, language, projectContext)
}

// getResponseSchema returns the expected JSON schema for structured responses
func (g *geminiProvider) getResponseSchema(analysisType string) map[string]interface{} {
	baseSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"projectName": map[string]string{"type": "string"},
			"summary":     map[string]string{"type": "string"},
			"strengths": map[string]interface{}{
				"type":  "array",
				"items": map[string]string{"type": "string"},
			},
			"weaknesses": map[string]interface{}{
				"type":  "array",
				"items": map[string]string{"type": "string"},
			},
			"recommendations": map[string]interface{}{
				"type":  "array",
				"items": map[string]string{"type": "string"},
			},
		},
		"required": []string{"projectName", "summary", "strengths", "weaknesses", "recommendations"},
	}
	switch analysisType {
	case "security":
		props := baseSchema["properties"].(map[string]interface{})
		props["securityRisks"] = map[string]interface{}{
			"type":  "array",
			"items": map[string]string{"type": "string"},
		}
	case "scalability":
		props := baseSchema["properties"].(map[string]interface{})
		props["bottlenecks"] = map[string]interface{}{
			"type":  "array",
			"items": map[string]string{"type": "string"},
		}
	}
	return baseSchema
}

// estimateTokens provides a rough token estimation
func (g *geminiProvider) estimateTokens(text string) int {
	// Rough estimation: ~4 characters per token
	return len(text) / 4
}

// estimateCost provides cost estimation for Gemini models
func (g *geminiProvider) estimateCost(model string, tokens int) float64 {
	var costPerToken float64
	switch {
	case strings.Contains(model, "flash"):
		costPerToken = 0.000000125 // $0.125/1M tokens for Gemini Flash
	case strings.Contains(model, "pro"):
		costPerToken = 0.000001 // $1/1M tokens for Gemini Pro
	default:
		costPerToken = 0.000000125 // Default to Flash pricing
	}
	return float64(tokens) * costPerToken
}

// Close gracefully closes the Gemini client
func (g *geminiProvider) Close() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.client != nil {
		return g.client.Close()
	}
	return nil
}
