package types

// NotifierConfig holds configuration for different notifiers
type NotifierConfig struct {
	DiscordWebhookURL string `json:"discord_webhook_url"`
	WhatsAppAPIKey    string `json:"whatsapp_api_key"`
	EmailSMTPServer   string `json:"email_smtp_server"`
	EmailUsername     string `json:"email_username"`
	EmailPassword     string `json:"email_password"`
	FromAddress       string `json:"from_address"`
}
