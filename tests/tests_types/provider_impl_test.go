// Pacote de teste externo para types.
package types_test

import (
	"errors"
	"testing"

	typesx "github.com/kubex-ecosystem/gemx/grompt/internal/types"
	m "github.com/kubex-ecosystem/gemx/grompt/mocks"
)

func TestProviderImpl_NameAndVersion(t *testing.T) {
	p := &typesx.ProviderImpl{VName: "openai", VVersion: "v1"}
	if p.Name() != "openai" {
		t.Fatalf("Name() = %q, want 'openai'", p.Name())
	}
	if p.Version() != "v1" {
		t.Fatalf("Version() = %q, want 'v1'", p.Version())
	}
}

func TestProviderImpl_Execute_WithAndWithoutAPI(t *testing.T) {
	// Sem API
	p := &typesx.ProviderImpl{}
	if _, err := p.Execute("hello"); err == nil {
		t.Fatalf("Execute() expected error when VAPI == nil")
	}

	// Com API mock
	p = &typesx.ProviderImpl{VName: "gemini", VAPI: &m.APIConfigMock{Resp: "ok"}}
	if got, err := p.Execute("hi"); err != nil || got != "ok" {
		t.Fatalf("Execute() got (%q,%v), want (ok,nil)", got, err)
	}
}

func TestProviderImpl_IsAvailable(t *testing.T) {
	p := &typesx.ProviderImpl{}
	if p.IsAvailable() {
		t.Fatalf("IsAvailable() = true, want false when VAPI nil")
	}

	p = &typesx.ProviderImpl{VAPI: &m.APIConfigMock{Available: true}}
	if !p.IsAvailable() {
		t.Fatalf("IsAvailable() = false, want true")
	}
}

func TestProviderImpl_GetCapabilities_DirectAPI(t *testing.T) {
	mock := &m.APIConfigMock{Models: []string{"m1", "m2"}}
	p := &typesx.ProviderImpl{VName: "openai", VAPI: mock}
	caps := p.GetCapabilities()
	if caps == nil {
		t.Fatalf("GetCapabilities() = nil, want non-nil")
	}
	if len(caps.Models) != 2 {
		t.Fatalf("models len = %d, want 2", len(caps.Models))
	}
	if caps.MaxTokens <= 0 {
		t.Fatalf("MaxTokens should be > 0")
	}
}

func TestProviderImpl_GetCapabilities_LazyInitFromConfig(t *testing.T) {
	// Sem VAPI, mas com VConfig que retorna um mock
	api := &m.APIConfigMock{Models: []string{"x"}}
	cfg := &m.ConfigMock{APIByName: map[string]typesx.IAPIConfig{"claude": api}}
	p := &typesx.ProviderImpl{VName: "claude", VConfig: cfg}
	caps := p.GetCapabilities()
	if caps == nil {
		t.Fatalf("GetCapabilities() = nil, want non-nil via lazy init")
	}
	if len(caps.Models) != 1 || caps.Models[0] != "x" {
		t.Fatalf("unexpected models: %#v", caps.Models)
	}
}

func TestProviderImpl_GetCapabilities_UnknownProvider(t *testing.T) {
	cfg := &m.ConfigMock{}
	p := &typesx.ProviderImpl{VName: "unknown", VConfig: cfg}
	if caps := p.GetCapabilities(); caps != nil {
		t.Fatalf("GetCapabilities() for unknown should be nil")
	}
}

func TestProviderImpl_GetCapabilities_ListModelsError(t *testing.T) {
	api := &m.APIConfigMock{RespErr: errors.New("boom")}
	p := &typesx.ProviderImpl{VName: "openai", VAPI: api}
	if caps := p.GetCapabilities(); caps != nil {
		t.Fatalf("GetCapabilities() should return nil on ListModels error")
	}
}
