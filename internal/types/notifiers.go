package types

import "time"

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
