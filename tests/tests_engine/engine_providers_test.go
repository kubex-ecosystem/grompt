package engine_test

import (
	"testing"

	eng "github.com/kubex-ecosystem/gemx/grompt/internal/engine"
	m "github.com/kubex-ecosystem/gemx/grompt/mocks"
)

func TestEngine_GetProviders_FromConfig(t *testing.T) {
	cfg := &m.ConfigMock{Keys: map[string]string{
		"openai":   "k1",
		"claude":   "k2",
		"gemini":   "k3",
		"deepseek": "k4",
		// ollama handled separately; engine adds one regardless of key
	}}
	e := eng.NewEngine(cfg)
	if e == nil {
		t.Fatalf("NewEngine returned nil")
	}

	provs := e.GetProviders()
	if len(provs) == 0 {
		t.Fatalf("expected at least one provider")
	}

	names := map[string]bool{}
	for _, p := range provs {
		names[p.Name()] = true
	}

	// Expected: all configured + ollama
	for _, want := range []string{"openai", "claude", "gemini", "deepseek", "ollama"} {
		if !names[want] {
			t.Fatalf("provider %q not found in engine", want)
		}
	}
}

func TestEngine_History_SaveAndGet(t *testing.T) {
	cfg := &m.ConfigMock{}
	e := eng.NewEngine(cfg)
	if e == nil {
		t.Fatalf("NewEngine returned nil")
	}

	if err := e.SaveToHistory("hello", "world"); err != nil {
		t.Fatalf("SaveToHistory error: %v", err)
	}
	h := e.GetHistory()
	if len(h) != 1 {
		t.Fatalf("history len = %d, want 1", len(h))
	}
	if h[0].Prompt != "hello" || h[0].Response != "world" {
		t.Fatalf("unexpected history entry: %#v", h[0])
	}
}
