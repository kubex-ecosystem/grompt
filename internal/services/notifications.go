package services

import (
	"context"
	"time"

	providers "github.com/kubex-ecosystem/grompt/internal/types"
)

// NotificationService handles sending notifications via different providers
type NotificationService struct {
	provider *providers.Provider
	timeout  time.Duration
}

// NewNotificationService creates a new NotificationService with the given provider and timeout_seconds
func NewNotificationService(config *providers.Config) *NotificationService {
	if config.Defaults.NotificationTimeoutSeconds <= 0 {
		config.Defaults.NotificationTimeoutSeconds = 60 // Default to 60 seconds if invalid
	}
	return &NotificationService{
		provider: config.Defaults.NotificationProvider,
		timeout:  time.Duration(config.Defaults.NotificationTimeoutSeconds) * time.Second,
	}
}

// SendNotification sends a notification message using the configured provider
func (n *NotificationService) SendNotification(ctx context.Context, event providers.NotificationEvent) error {
	ctx, cancel := context.WithTimeout(ctx, n.timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if n.provider == nil {
			return nil // No provider configured, nothing to do
		}
		// Lógica temporária para evitar problemas de 'non implemented'
		p := *n.provider
		return p.Notify(ctx, event)
	}
}

func (n *NotificationService) Name() string {
	return "NotificationService"
}
