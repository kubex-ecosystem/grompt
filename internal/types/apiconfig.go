package types

import (
	"context"
	"errors"
)

// ---------- API config implementation ----------

type apiConfig struct {
	provider string
	cfg      *ServerConfigImpl
}

func (a *apiConfig) IsAvailable() bool {
	if a == nil || a.cfg == nil {
		return false
	}
	return a.cfg.GetAPIKey(a.provider) != ""
}

func (a *apiConfig) IsDemoMode() bool { return false }

func (a *apiConfig) Version() string { return "gateway-v1" }

func (a *apiConfig) ListModels() ([]string, error) {
	model := a.cfg.DefaultModels[a.provider]
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
	if a == nil || a.cfg == nil || a.cfg.Engine == nil {
		return "", errors.New("prompt engine not initialized")
	}

	vars := map[string]interface{}{
		"provider": a.provider,
	}
	if model != "" {
		vars["model"] = model
	}

	result, err := a.cfg.Engine.InvokeProvider(context.Background(), a.provider, prompt, vars)
	if err != nil {
		return "", err
	}
	return result.Response, nil
}
