package cli

import (
	"context"
	"fmt"
	"html/template"
	"strings"

	"github.com/kubex-ecosystem/grompt/internal/engine"
	"github.com/kubex-ecosystem/grompt/internal/interfaces"
)

type Gateway struct {
	engine *engine.Engine
}

func NewGateway(engine *engine.Engine) *Gateway {
	return &Gateway{
		engine: engine,
	}
}

func (g *Gateway) ProcessRequest(ctx context.Context, request *interfaces.ChatRequest) (*interfaces.Result, error) {
	if g.engine == nil {
		return nil, fmt.Errorf("engine is not initialized")
	}

	template, err := template.New("prompt").Parse(request.PromptTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse prompt template: %w", err)
	}

	var promptBuilder strings.Builder
	err = template.Execute(&promptBuilder, request.Messages)
	if err != nil {
		return nil, fmt.Errorf("failed to execute prompt template: %w", err)
	}

	// Process the prompt using the engine
	result, err := g.engine.ProcessPrompt(ctx, request.PromptTemplate, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to process prompt: %w", err)
	}

	// Create and return the response
	response := &interfaces.Result{
		ID:        result.ID,
		Prompt:    result.Prompt,
		Response:  result.Response,
		Provider:  result.Provider,
		Variables: result.Variables,
		Timestamp: result.Timestamp,
	}

	return response, nil
}
