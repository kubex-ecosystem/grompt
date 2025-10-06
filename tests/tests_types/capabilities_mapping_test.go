package types_test

import (
	"context"
	"testing"

	typesx "github.com/kubex-ecosystem/grompt/internal/types"
	m "github.com/kubex-ecosystem/grompt/mocks"
)

func TestCapabilities_MaxTokens_ByProvider(t *testing.T) {
	// Reuse the same mock API returning one model to satisfy ListModels
	api := &m.APIConfigMock{Models: []string{"dummy"}}
	cases := []struct {
		name     string
		provider string
		wantMin  int
	}{
		{"openai", "openai", 4000},
		{"claude", "claude", 8000},
		{"gemini", "gemini", 8000},
		{"deepseek", "deepseek", 4000},
		{"ollama", "ollama", 2000},
	}
	for _, c := range cases {
		p := &typesx.ProviderImpl{VName: c.provider, VAPI: api}
		caps := p.GetCapabilities(context.Background())
		if caps == nil {
			t.Fatalf("%s: caps is nil", c.name)
		}
		if caps.MaxTokens < c.wantMin {
			t.Fatalf("%s: MaxTokens=%d < wantMin=%d", c.name, caps.MaxTokens, c.wantMin)
		}
	}
}
