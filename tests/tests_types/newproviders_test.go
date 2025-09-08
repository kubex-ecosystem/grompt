package types_test

import (
    "testing"

    m "github.com/kubex-ecosystem/grompt/mocks"
    typesx "github.com/kubex-ecosystem/grompt/internal/types"
)

func TestNewProviders_FiltersByConfiguredAndAvailable(t *testing.T) {
    cfg := &m.ConfigMock{
        Keys: map[string]string{
            "openai":  "k1",
            "claude":  "k2",
            "gemini":  "k3",
            "deepseek": "k4",
            "ollama":  "http://localhost:11434",
            // chatgpt intentionally missing to be excluded
        },
        APIByName: map[string]typesx.IAPIConfig{
            // mark only some as available
            "openai":  &m.APIConfigMock{Available: true},
            "claude":  &m.APIConfigMock{Available: false},
            "gemini":  &m.APIConfigMock{Available: true},
            "deepseek": &m.APIConfigMock{Available: true},
            "ollama":  &m.APIConfigMock{Available: true},
        },
    }

    got := typesx.NewProviders(cfg)
    if got == nil { t.Fatalf("NewProviders() returned nil") }

    names := make(map[string]bool)
    for _, p := range got {
        names[p.Name()] = true
    }

    // Expected: include openai, gemini, deepseek, ollama; exclude claude (not available) and chatgpt (no key)
    if !names["openai"] || !names["gemini"] || !names["deepseek"] || !names["ollama"] {
        t.Fatalf("expected providers missing, got: %#v", names)
    }
    if names["claude"] || names["chatgpt"] {
        t.Fatalf("unexpected providers included, got: %#v", names)
    }
}
