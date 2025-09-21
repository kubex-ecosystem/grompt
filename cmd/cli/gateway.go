package cli

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kubex-ecosystem/grompt/internal/gateway"
	"github.com/spf13/cobra"

	gl "github.com/kubex-ecosystem/grompt/internal/module/logger"
)

// GatewayCmds returns the gateway command with subcommands
func GatewayCmds() *cobra.Command {
	// Main gateway command
	var (
		bindingAddress string
		port           string
		configPath     string
		debug          bool
		enableCORS     bool
	)

	rootCmd := &cobra.Command{
		Use:   "gateway",
		Short: "🚀 Grompt Gateway - AI Provider Gateway with Repository Intelligence",
		Long: `Grompt Gateway provides a unified API for AI providers with enterprise features.

Features:
  • Multi-provider AI gateway (OpenAI, Anthropic, Gemini, Groq, etc.)
  • Repository Intelligence APIs (DORA metrics, Code Health, AI Impact)
  • Enterprise production features (rate limiting, circuit breaker, health checks)
  • Real-time streaming with Server-Sent Events (SSE)
  • BYOK (Bring Your Own Key) support
  • Tenant and user isolation`,
		Example: `  # Start gateway with default settings
  grompt gateway serve

  # Start with custom config and address
  grompt gateway serve --addr :8080 --config ./config/config.example.yml

  # Start with debug mode and CORS enabled
  grompt gateway serve --debug --cors`,
	}

	// Serve subcommand
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the gateway server (with GUI)",
		Long:  "Start the Grompt Gateway server with enterprise features (with GUI)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load environment variables from .env file if exists
			loadEnv(".env")

			// Validate port
			if _, err := net.LookupPort("tcp", port); err != nil {
				return fmt.Errorf("invalid port %q: %w", port, err)
			}

			return startGateway(&gateway.ServerConfig{
				Addr:            net.JoinHostPort(bindingAddress, port),
				ProvidersConfig: configPath,
				Debug:           debug,
				EnableCORS:      enableCORS,
			})
		},
	}

	// Add flags to serve command
	serveCmd.Flags().StringVarP(&bindingAddress, "binding", "b", getEnv("ADDR", "0.0.0.0"), "Server address")
	serveCmd.Flags().StringVarP(&port, "port", "p", getEnv("PORT", "8080"), "Server port")
	serveCmd.Flags().BoolVar(&enableCORS, "cors", true, "Enable CORS headers")
	serveCmd.Flags().StringVarP(&configPath, "config", "c", getEnv("PROVIDERS_CFG", "config/config.example.yml"), "Providers config file")
	serveCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")

	// Add status subcommand
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Check gateway status",
		Long:  "Check the health and status of the running gateway",
		RunE:  statusCommand,
	}

	// Add advise subcommand (legacy support)
	adviseCmd := &cobra.Command{
		Use:   "advise",
		Short: "Generate repository advice using AI",
		Long:  "Generate repository advice using AI providers with scorecard data",
		RunE:  adviseCommand,
	}

	// Add subcommands
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(adviseCmd)

	return rootCmd
}

// startGateway starts the gateway server with given configuration
func startGateway(config *gateway.ServerConfig) error {
	server, err := gateway.NewServer(config)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	return server.Start()
}

// statusCommand checks the gateway status
func statusCommand(cmd *cobra.Command, args []string) error {
	port := getEnv("PORT", "8080")
	bindingAddress := getEnv("ADDR", "0.0.0.0")
	targetAddress := ""
	if bindingAddress == "0.0.0.0" {
		targetAddress = net.JoinHostPort("localhost", port)
	} else {
		targetAddress = net.JoinHostPort(bindingAddress, port)
	}

	resp, err := http.Get("http://" + targetAddress + "/healthz")
	if err != nil {
		return fmt.Errorf("gateway not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Gateway is healthy")

		// Also check Repository Intelligence endpoints
		riResp, err := http.Get("http://" + targetAddress + "/api/v1/health")
		if err == nil && riResp.StatusCode == http.StatusOK {
			fmt.Println("✅ Repository Intelligence APIs are available")
			riResp.Body.Close()
		} else {
			fmt.Println("⚠️  Repository Intelligence APIs not fully initialized")
		}

		return nil
	}

	return fmt.Errorf("gateway unhealthy (status: %d)", resp.StatusCode)
}

// adviseCommand provides legacy advise functionality
func adviseCommand(cmd *cobra.Command, args []string) error {
	fmt.Println("🤖 Repository Advice using AI")
	fmt.Println("This command provides repository advice using scorecard data and AI providers.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  grompt gateway advise --mode exec --provider openai --model gpt-4o-mini --scorecard ./scorecard.json")
	fmt.Println("")
	fmt.Println("Available modes: exec, code, ops, community")
	fmt.Println("Available providers: openai, anthropic, gemini, groq")

	// TODO: Implement full advise functionality using cmdAdvise
	return fmt.Errorf("advise command not fully implemented yet")
}

// getEnv returns environment variable value or default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loadEnv(configPath string) {
	// Initialize environment variables, set them inside environment
	if err := godotenv.Load(configPath); err != nil {
		gl.Log("warning", fmt.Sprintf("No .env file found at %s, proceeding with existing environment variables", configPath))
	}
}
