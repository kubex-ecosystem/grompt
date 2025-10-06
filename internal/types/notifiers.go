package types

import (
	"context"
	"time"
)

// NotificationEvent represents a notification to be sent
type NotificationEvent struct {
	Type      string                 `json:"type"` // "discord", "whatsapp", "email"
	Recipient string                 `json:"recipient"`
	Subject   string                 `json:"subject"`
	Content   string                 `json:"content"`
	Priority  string                 `json:"priority"` // "low", "medium", "high", "critical"
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
}

type Notifier interface {
	Notify(ctx context.Context, event NotificationEvent) error
}

// NotifierConfig holds configuration for different notifiers
type NotifierConfig struct {
	DiscordWebhookURL string `json:"discord_webhook_url"`
	WhatsAppAPIKey    string `json:"whatsapp_api_key"`
	EmailSMTPServer   string `json:"email_smtp_server"`
	EmailUsername     string `json:"email_username"`
	EmailPassword     string `json:"email_password"`
	FromAddress       string `json:"from_address"`
}
