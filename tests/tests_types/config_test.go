// Pacote de teste externo para types.
package types_test

import (
	"os"
	"testing"

	typesx "github.com/kubex-ecosystem/grompt/internal/types"
)

func TestConfig_GetPort_DefaultAndCustom(t *testing.T) {
	cfg := &typesx.Config{}
	if got := cfg.GetPort(); got != "8080" {
		t.Fatalf("GetPort() default = %q, want 8080", got)
	}
	cfg.Port = "9090"
	if got := cfg.GetPort(); got != "9090" {
		t.Fatalf("GetPort() custom = %q, want 9090", got)
	}
}

func TestConfig_SetAndGetAPIKey(t *testing.T) {
	cfg := &typesx.Config{}
	if err := cfg.SetAPIKey("openai", "abc"); err != nil {
		t.Fatalf("SetAPIKey() error = %v", err)
	}
	if got := cfg.GetAPIKey("openai"); got != "abc" {
		t.Fatalf("GetAPIKey() = %q, want 'abc'", got)
	}
}

func TestConfig_GetAPIEndpoint_Ollama(t *testing.T) {
	cfg := &typesx.Config{}
	cfg.OllamaEndpoint = "http://localhost:11434"
	if got := cfg.GetAPIEndpoint("ollama"); got != "http://localhost:11434" {
		t.Fatalf("GetAPIEndpoint(ollama) = %q, want http://localhost:11434", got)
	}
}

func TestConfig_GetAPIKey_EnvFallback(t *testing.T) {
	// Garante restauração do ambiente
	orig := os.Getenv("GEMINI_API_KEY")
	t.Cleanup(func() { _ = os.Setenv("GEMINI_API_KEY", orig) })

	_ = os.Setenv("GEMINI_API_KEY", "env-val")

	cfg := &typesx.Config{GeminiAPIKey: ""}
	if got := cfg.GetAPIKey("gemini"); got != "env-val" {
		t.Fatalf("GetAPIKey(gemini) = %q, want 'env-val' (fallback env)", got)
	}
}

func TestConfig_GetBaseGenerationPrompt_NotEmpty(t *testing.T) {
	cfg := &typesx.Config{}
	got := cfg.GetBaseGenerationPrompt([]string{"a", "b"}, "purpose", "Code", "pt", 1000)
	if got == "" {
		t.Fatalf("GetBaseGenerationPrompt() returned empty string")
	}
	if wantSub := "PROPÓSITO DO PROMPT"; got != "" && !contains(got, wantSub) {
		t.Fatalf("prompt should contain %q", wantSub)
	}
}

// Helper simples para evitar importar strings direto
func contains(s, sub string) bool {
	return len(s) >= len(sub) && (func() bool { return (len(sub) == 0) || (len(s) >= len(sub) && (indexOf(s, sub) >= 0)) })()
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
