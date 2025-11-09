package cli

import (
	"context"
	"fmt"
	"html/template"
	"strings"

	"github.com/kubex-ecosystem/grompt/internal/engine"
	"github.com/kubex-ecosystem/grompt/internal/interfaces"
	"github.com/spf13/cobra"
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

func GatewayCmd() *cobra.Command {
	var gatewayCmd = &cobra.Command{
		Use:   "gateway",
		Short: "Gateway commands",
		Long:  "Commands for managing the gateway",
	}

	gatewayCmd.AddCommand(startGatewayServerCmd())

	return gatewayCmd
}

func startGatewayServerCmd() *cobra.Command {
	var port, bindAddr string

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the gateway server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Starting gateway server on %s:%s...\n", bindAddr, port)
			// Implement the logic to start the gateway server here
		},
	}

	startCmd.Flags().StringVarP(&bindAddr, "bind", "b", "localhost", "Address to bind the gateway server to")
	startCmd.Flags().StringVarP(&port, "port", "p", "9090", "Port to run the gateway server on")

	return startCmd
}
