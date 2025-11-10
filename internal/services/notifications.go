package services

import (
	"context"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/interfaces"
	providers "github.com/kubex-ecosystem/grompt/internal/types"
)

// NotificationService handles sending notifications via different providers
type NotificationService struct {
	provider *interfaces.Provider
	timeout  time.Duration
}

// NewNotificationService creates a new NotificationService with the given provider and timeout_seconds
func NewNotificationService(config *providers.Config) *NotificationService {
	if config.Defaults.NotificationTimeoutSeconds <= 0 {
		config.Defaults.NotificationTimeoutSeconds = 60 // Default to 60 seconds if invalid
	}
	var notificationProvider interfaces.Provider
	notProvider := config.Defaults.NotificationProvider
	if notProvider == nil {
		notificationProvider = nil
	} else {
		provider, ok := notProvider.(interfaces.Provider)
		if !ok {
			notificationProvider = nil
		} else {
			notificationProvider = provider
		}
	}
	return &NotificationService{
		provider: &notificationProvider,
		timeout:  time.Duration(config.Defaults.NotificationTimeoutSeconds) * time.Second,
	}
}

// SendNotification sends a notification message using the configured provider
func (n *NotificationService) SendNotification(ctx context.Context, event interfaces.NotificationEvent) error {
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
