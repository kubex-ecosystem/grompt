package types

import (
	"context"
	"errors"
)

// ---------- API config implementation ----------

type apiConfig struct {
	Provider string 											 `json:"provider,omitempty" yaml:"provider,omitempty" mapstructure:"provider,omitempty"`
	APIServerConfig      *ServerConfigImpl `json:"-" yaml:"-" mapstructure:"-"`
}

func (a *apiConfig) IsAvailable() bool {
	if a == nil || a.APIServerConfig == nil {
		return false
	}
	return a.APIServerConfig.GetAPIKey(a.Provider) != ""
}

func (a *apiConfig) IsDemoMode() bool { return false }

func (a *apiConfig) Version() string { return "gateway-v1" }

func (a *apiConfig) ListModels() ([]string, error) {
	model := a.APIServerConfig.DefaultModels[a.Provider]
	if model == "" {
		return []string{}, nil
	}
	return []string{model}, nil
}

func (a *apiConfig) GetCommonModels() []string {
	models, _ := a.ListModels()
	return models
}

func (a *apiConfig) Complete(prompt string, maxTokens int, model string) (string, error) {
	if a == nil || a.APIServerConfig == nil || a.APIServerConfig.Engine == nil {
		return "", errors.New("prompt engine not initialized")
	}

	vars := map[string]interface{}{
		"provider": a.Provider,
	}
	if model != "" {
		vars["model"] = model
	}

	result, err := a.APIServerConfig.Engine.InvokeProvider(context.Background(), a.Provider, prompt, vars)
	if err != nil {
		return "", err
	}
	return result.Response, nil
}
