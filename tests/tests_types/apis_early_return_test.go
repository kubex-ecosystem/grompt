package types_test

import (
	"testing"

	typesx "github.com/kubex-ecosystem/gemx/grompt/internal/types"
)

func TestAPIs_Complete_ReturnsErrorWithoutKey(t *testing.T) {
	// OpenAI
	if _, err := typesx.NewOpenAIAPI("").Complete("hi", 1, ""); err == nil {
		t.Fatalf("OpenAI Complete should error without API key")
	}
	// Claude
	if _, err := typesx.NewClaudeAPI("").Complete("hi", 1, ""); err == nil {
		t.Fatalf("Claude Complete should error without API key")
	}
	// DeepSeek
	if _, err := typesx.NewDeepSeekAPI("").Complete("hi", 1, ""); err == nil {
		t.Fatalf("DeepSeek Complete should error without API key")
	}
	// ChatGPT
	if _, err := typesx.NewChatGPTAPI("").Complete("hi", 1, ""); err == nil {
		t.Fatalf("ChatGPT Complete should error without API key")
	}
	// Gemini
	if _, err := typesx.NewGeminiAPI("").Complete("hi", 1, ""); err == nil {
		t.Fatalf("Gemini Complete should error without API key")
	}
}

func TestAPIs_ListModels_ReturnsErrorWithoutKey_WhenApplicable(t *testing.T) {
	// OpenAI requires key for ListModels
	if _, err := typesx.NewOpenAIAPI("").ListModels(); err == nil {
		t.Fatalf("OpenAI ListModels should error without API key")
	}
	// ChatGPT requires key for ListModels
	if _, err := typesx.NewChatGPTAPI("").ListModels(); err == nil {
		t.Fatalf("ChatGPT ListModels should error without API key")
	}
	// Gemini returns common models without remote call; should not error
	if models, err := typesx.NewGeminiAPI("").ListModels(); err != nil || len(models) == 0 {
		t.Fatalf("Gemini ListModels should not error without key, got (%v, %v)", models, err)
	}
}
