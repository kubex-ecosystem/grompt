package types_test

import (
    "testing"

    typesx "github.com/kubex-ecosystem/grompt/internal/types"
)

func TestGetAPIConfig_KnownProviders_NotNil(t *testing.T) {
    cfg := &typesx.Config{}
    cases := []string{"openai", "deepseek", "ollama", "claude", "gemini", "chatgpt"}
    for _, name := range cases {
        if api := cfg.GetAPIConfig(name); api == nil {
            t.Fatalf("GetAPIConfig(%q) = nil, want non-nil", name)
        }
    }
}

func TestGetAPIConfig_UnknownProvider_Nil(t *testing.T) {
    cfg := &typesx.Config{}
    if api := cfg.GetAPIConfig("unknown"); api != nil {
        t.Fatalf("GetAPIConfig(unknown) = %#v, want nil", api)
    }
}
