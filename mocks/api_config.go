package mocks

import (
	t "github.com/kubex-ecosystem/grompt/internal/types"
)

// APIConfigMock implements t.IAPIConfig for tests.
type APIConfigMock struct {
	Available    bool
	DemoMode     bool
	Version      string
	Models       []string
	CommonModels []string
	Resp         string
	RespErr      error
}

func (m *APIConfigMock) IsAvailable() bool             { return m.Available }
func (m *APIConfigMock) IsDemoMode() bool              { return m.DemoMode }
func (m *APIConfigMock) GetVersion() string            { return m.Version }
func (m *APIConfigMock) ListModels() ([]string, error) { return m.Models, m.RespErr }
func (m *APIConfigMock) GetCommonModels() []string     { return m.CommonModels }
func (m *APIConfigMock) Complete(prompt string, max int, model string) (string, error) {
	return m.Resp, m.RespErr
}

// ConfigMock implements t.IConfig for tests.
type ConfigMock struct {
	Port       string
	Keys       map[string]string
	Endpoints  map[string]string
	APIByName  map[string]t.IAPIConfig
	BasePrompt string
}

func (c *ConfigMock) GetAPIConfig(provider string) t.IAPIConfig {
	if c.APIByName == nil {
		return nil
	}
	return c.APIByName[provider]
}

func (c *ConfigMock) GetPort() string {
	if c.Port == "" {
		return "8080"
	}
	return c.Port
}

func (c *ConfigMock) GetAPIKey(provider string) string {
	if c.Keys == nil {
		return ""
	}
	return c.Keys[provider]
}

func (c *ConfigMock) SetAPIKey(provider string, key string) error {
	if c.Keys == nil {
		c.Keys = map[string]string{}
	}
	c.Keys[provider] = key
	return nil
}

func (c *ConfigMock) GetAPIEndpoint(provider string) string {
	if c.Endpoints == nil {
		return ""
	}
	return c.Endpoints[provider]
}

func (c *ConfigMock) GetBaseGenerationPrompt(ideas []string, purpose, purposeType, lang string, maxLength int) string {
	return c.BasePrompt
}
