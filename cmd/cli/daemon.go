// Package cli provides the daemon command for background service operations
package cli

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/daemon"
	"github.com/spf13/cobra"
)

var (
	gobeURL             string
	gobeAPIKey          string
	autoScheduleEnabled bool
	scheduleCron        string
	notifyChannels      []string
	healthCheckInterval time.Duration
)

// NewDaemonCommand creates the daemon command
func NewDaemonCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "Start grompt as background daemon with GoBE integration",
		Long: `Start the grompt as a background daemon service that integrates with GoBE backend.

The daemon provides:
• Automatic repository analysis scheduling
• Integration with KubeX AI Squad system
• Discord/WhatsApp/Email notifications
• Health monitoring and reporting
• Meta-recursivity coordination with lookatni/grompt

Examples:
  grompt daemon --gobe-url=http://localhost:3000 --gobe-api-key=abc123
  grompt daemon --auto-schedule --schedule-cron="0 2 * * *"
  grompt daemon --notify-channels=discord,email`,
		RunE: runDaemon,
	}

	// GoBE Integration flags
	cmd.Flags().StringVar(&gobeURL, "gobe-url",
		getEnvOrDefault("GOBE_URL", "http://localhost:3000"),
		"GoBE backend URL")
	cmd.Flags().StringVar(&gobeAPIKey, "gobe-api-key",
		os.Getenv("GOBE_API_KEY"),
		"GoBE API key for authentication")

	// Scheduling flags
	cmd.Flags().BoolVar(&autoScheduleEnabled, "auto-schedule", false,
		"Enable automatic repository analysis scheduling")
	cmd.Flags().StringVar(&scheduleCron, "schedule-cron", "0 2 * * *",
		"Cron expression for automatic scheduling (default: daily at 2 AM)")

	// Notification flags
	cmd.Flags().StringSliceVar(&notifyChannels, "notify-channels",
		[]string{"discord"},
		"Notification channels (discord,email,webhook)")

	// Health monitoring flags
	cmd.Flags().DurationVar(&healthCheckInterval, "health-interval",
		5*time.Minute,
		"Health check interval")

	return cmd
}

func runDaemon(cmd *cobra.Command, args []string) error {
	// Validate required flags
	if gobeAPIKey == "" {
		return fmt.Errorf("--gobe-api-key is required (or set GOBE_API_KEY env var)")
	}

	// Create daemon configuration
	config := daemon.DaemonConfig{
		GoBeURL:              gobeURL,
		GoBeAPIKey:           gobeAPIKey,
		AutoScheduleEnabled:  autoScheduleEnabled,
		ScheduleCron:         scheduleCron,
		NotificationChannels: notifyChannels,
		HealthCheckInterval:  healthCheckInterval,
	}

	// Create and start daemon
	d := daemon.NewGromptDaemon(config)

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start daemon
	if err := d.Start(); err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}

	// Print startup information
	printDaemonInfo(config)

	// Wait for shutdown signal
	<-sigChan
	log.Println("🔥 Received shutdown signal...")

	// Graceful shutdown
	d.Stop()
	log.Println("✅ Grompt daemon stopped gracefully")

	return nil
}

func printDaemonInfo(config daemon.DaemonConfig) {
	fmt.Println()
	fmt.Println("🚀 ============================================================")
	fmt.Println("🤖   GROMPT DAEMON - Repository Intelligence Platform")
	fmt.Println("🚀 ============================================================")
	fmt.Println()
	fmt.Printf("🏗️  GoBE Integration: %s\n", config.GoBeURL)
	fmt.Printf("📅 Auto Schedule: %v", config.AutoScheduleEnabled)
	if config.AutoScheduleEnabled {
		fmt.Printf(" (%s)", config.ScheduleCron)
	}
	fmt.Println()
	fmt.Printf("🔔 Notifications: %v\n", config.NotificationChannels)
	fmt.Printf("🏥 Health Checks: every %v\n", config.HealthCheckInterval)
	fmt.Println()
	fmt.Println("📊 CAPABILITIES:")
	fmt.Println("   • Repository Intelligence Analysis")
	fmt.Println("   • DORA Metrics Collection")
	fmt.Println("   • Code Health Index (CHI)")
	fmt.Println("   • AI Impact Analysis")
	fmt.Println("   • Automated Scheduling")
	fmt.Println("   • Multi-channel Notifications")
	fmt.Println("   • KubeX AI Squad Integration")
	fmt.Println("   • Meta-recursivity Coordination")
	fmt.Println()
	fmt.Println("🎯 INTEGRATION POINTS:")
	fmt.Println("   • GoBE Backend APIs")
	fmt.Println("   • Discord Webhooks")
	fmt.Println("   • Email Notifications")
	fmt.Println("   • GitHub Events")
	fmt.Println("   • Jira Workflows (planned)")
	fmt.Println("   • WakaTime Analytics (planned)")
	fmt.Println()
	fmt.Println("🔄 META-RECURSIVITY:")
	fmt.Println("   • Coordinates with lookatni (analysis)")
	fmt.Println("   • Orchestrates grompt (improvement)")
	fmt.Println("   • Manages continuous optimization")
	fmt.Println()
	fmt.Println("✅ Daemon running... Press Ctrl+C to stop")
	fmt.Println()
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
