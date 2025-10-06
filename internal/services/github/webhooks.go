// Package github provides webhook signature validation and management.
package github

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// WebhookValidator handles GitHub webhook signature validation and deduplication
type WebhookValidator struct {
	secret string

	// Deduplication tracking
	deliveryMutex sync.RWMutex
	deliveries    map[string]time.Time
	maxAge        time.Duration
}

// NewWebhookValidator creates a new webhook validator with the given secret
func NewWebhookValidator(secret string) *WebhookValidator {
	return &WebhookValidator{
		secret:     secret,
		deliveries: make(map[string]time.Time),
		maxAge:     24 * time.Hour, // Keep delivery IDs for 24 hours
	}
}

// ValidateSignature validates the GitHub webhook signature
func (wv *WebhookValidator) ValidateSignature(payload []byte, signature string) error {
	if wv.secret == "" {
		return fmt.Errorf("webhook secret not configured")
	}

	if signature == "" {
		return fmt.Errorf("missing X-Hub-Signature-256 header")
	}

	// GitHub sends signature as "sha256=<hash>"
	if !strings.HasPrefix(signature, "sha256=") {
		return fmt.Errorf("invalid signature format, expected sha256= prefix")
	}

	expectedSignature := wv.computeSignature(payload)
	receivedSignature := signature[7:] // Remove "sha256=" prefix

	// Use constant-time comparison to prevent timing attacks
	if !hmac.Equal([]byte(expectedSignature), []byte(receivedSignature)) {
		return fmt.Errorf("signature validation failed")
	}

	return nil
}

// computeSignature computes the HMAC-SHA256 signature for the payload
func (wv *WebhookValidator) computeSignature(payload []byte) string {
	mac := hmac.New(sha256.New, []byte(wv.secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

// IsDuplicate checks if a delivery ID has been processed before
func (wv *WebhookValidator) IsDuplicate(deliveryID string) bool {
	if deliveryID == "" {
		return false // If no delivery ID, assume not duplicate
	}

	wv.deliveryMutex.RLock()
	_, exists := wv.deliveries[deliveryID]
	wv.deliveryMutex.RUnlock()

	return exists
}

// MarkDelivery marks a delivery ID as processed
func (wv *WebhookValidator) MarkDelivery(deliveryID string) {
	if deliveryID == "" {
		return
	}

	wv.deliveryMutex.Lock()
	wv.deliveries[deliveryID] = time.Now()
	wv.deliveryMutex.Unlock()

	// Clean up old entries
	wv.cleanupOldDeliveries()
}

// cleanupOldDeliveries removes delivery IDs older than maxAge
func (wv *WebhookValidator) cleanupOldDeliveries() {
	now := time.Now()
	cutoff := now.Add(-wv.maxAge)

	wv.deliveryMutex.Lock()
	defer wv.deliveryMutex.Unlock()

	for deliveryID, timestamp := range wv.deliveries {
		if timestamp.Before(cutoff) {
			delete(wv.deliveries, deliveryID)
		}
	}
}

// ValidateAndProcess validates webhook signature and checks for duplicates
func (wv *WebhookValidator) ValidateAndProcess(r *http.Request, payload []byte) error {
	// Get headers
	signature := r.Header.Get("X-Hub-Signature-256")
	deliveryID := r.Header.Get("X-GitHub-Delivery")

	// Validate signature first
	if err := wv.ValidateSignature(payload, signature); err != nil {
		return fmt.Errorf("signature validation failed: %w", err)
	}

	// Check for duplicate delivery
	if wv.IsDuplicate(deliveryID) {
		return fmt.Errorf("duplicate delivery ID: %s", deliveryID)
	}

	// Mark as processed
	wv.MarkDelivery(deliveryID)

	return nil
}

// WebhookManager handles webhook registration and management
type WebhookManager struct {
	service *Service
	config  *Config
}

// NewWebhookManager creates a new webhook manager
func NewWebhookManager(service *Service, config *Config) *WebhookManager {
	return &WebhookManager{
		service: service,
		config:  config,
	}
}

// WebhookConfig represents webhook configuration
type WebhookConfig struct {
	URL         string   `json:"url"`
	ContentType string   `json:"content_type"`
	Secret      string   `json:"secret,omitempty"`
	InsecureSSL bool     `json:"insecure_ssl"`
	Events      []string `json:"events"`
	Active      bool     `json:"active"`
}

// CreateWebhook creates a webhook for a repository
func (wm *WebhookManager) CreateWebhook(owner, repo string, webhookConfig WebhookConfig) (*GitHubWebhook, error) {
	installationID := wm.service.installationID
	if installationID == 0 && wm.service.client.auth.IsUsingAppAuth() {
		var err error
		installationID, err = wm.service.client.auth.GetInstallationID(owner, repo)
		if err != nil {
			return nil, fmt.Errorf("failed to get installation ID: %w", err)
		}
	}

	// Prepare webhook payload
	payload := map[string]interface{}{
		"name":   "web",
		"active": webhookConfig.Active,
		"events": webhookConfig.Events,
		"config": map[string]interface{}{
			"url":          webhookConfig.URL,
			"content_type": webhookConfig.ContentType,
			"insecure_ssl": "0",
		},
	}

	if webhookConfig.Secret != "" {
		payload["config"].(map[string]interface{})["secret"] = webhookConfig.Secret
	}

	if webhookConfig.InsecureSSL {
		payload["config"].(map[string]interface{})["insecure_ssl"] = "1"
	}

	// Marshal payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	// Create webhook
	path := fmt.Sprintf("/repos/%s/%s/hooks", owner, repo)
	data, err := wm.service.client.Post(nil, path, payloadBytes, installationID)
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}

	var webhook GitHubWebhook
	if err := json.Unmarshal(data, &webhook); err != nil {
		return nil, fmt.Errorf("failed to parse webhook response: %w", err)
	}

	return &webhook, nil
}

// ListWebhooks lists webhooks for a repository
func (wm *WebhookManager) ListWebhooks(owner, repo string) ([]GitHubWebhook, error) {
	installationID := wm.service.installationID
	if installationID == 0 && wm.service.client.auth.IsUsingAppAuth() {
		var err error
		installationID, err = wm.service.client.auth.GetInstallationID(owner, repo)
		if err != nil {
			return nil, fmt.Errorf("failed to get installation ID: %w", err)
		}
	}

	path := fmt.Sprintf("/repos/%s/%s/hooks", owner, repo)
	data, err := wm.service.client.Get(nil, path, installationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list webhooks: %w", err)
	}

	var webhooks []GitHubWebhook
	if err := json.Unmarshal(data, &webhooks); err != nil {
		return nil, fmt.Errorf("failed to parse webhooks response: %w", err)
	}

	return webhooks, nil
}

// DeleteWebhook deletes a webhook from a repository
func (wm *WebhookManager) DeleteWebhook(owner, repo string, webhookID int) error {
	installationID := wm.service.installationID
	if installationID == 0 && wm.service.client.auth.IsUsingAppAuth() {
		var err error
		installationID, err = wm.service.client.auth.GetInstallationID(owner, repo)
		if err != nil {
			return fmt.Errorf("failed to get installation ID: %w", err)
		}
	}

	path := fmt.Sprintf("/repos/%s/%s/hooks/%d", owner, repo, webhookID)
	_, err := wm.service.client.Delete(nil, path, installationID)
	if err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}

	return nil
}

// UpdateWebhook updates a webhook configuration
func (wm *WebhookManager) UpdateWebhook(owner, repo string, webhookID int, webhookConfig WebhookConfig) (*GitHubWebhook, error) {
	installationID := wm.service.installationID
	if installationID == 0 && wm.service.client.auth.IsUsingAppAuth() {
		var err error
		installationID, err = wm.service.client.auth.GetInstallationID(owner, repo)
		if err != nil {
			return nil, fmt.Errorf("failed to get installation ID: %w", err)
		}
	}

	// Prepare webhook payload
	payload := map[string]interface{}{
		"active": webhookConfig.Active,
		"events": webhookConfig.Events,
		"config": map[string]interface{}{
			"url":          webhookConfig.URL,
			"content_type": webhookConfig.ContentType,
			"insecure_ssl": "0",
		},
	}

	if webhookConfig.Secret != "" {
		payload["config"].(map[string]interface{})["secret"] = webhookConfig.Secret
	}

	if webhookConfig.InsecureSSL {
		payload["config"].(map[string]interface{})["insecure_ssl"] = "1"
	}

	// Marshal payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	// Update webhook
	path := fmt.Sprintf("/repos/%s/%s/hooks/%d", owner, repo, webhookID)
	data, err := wm.service.client.Patch(nil, path, payloadBytes, installationID)
	if err != nil {
		return nil, fmt.Errorf("failed to update webhook: %w", err)
	}

	var webhook GitHubWebhook
	if err := json.Unmarshal(data, &webhook); err != nil {
		return nil, fmt.Errorf("failed to parse webhook response: %w", err)
	}

	return &webhook, nil
}

// GitHubWebhook represents a GitHub webhook
type GitHubWebhook struct {
	ID        int                    `json:"id"`
	Name      string                 `json:"name"`
	Active    bool                   `json:"active"`
	Events    []string               `json:"events"`
	Config    map[string]interface{} `json:"config"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	URL       string                 `json:"url"`
	TestURL   string                 `json:"test_url"`
	PingURL   string                 `json:"ping_url"`
}

// AutoSetupWebhook automatically sets up a webhook for the grompt
func (wm *WebhookManager) AutoSetupWebhook(owner, repo, callbackURL string) (*GitHubWebhook, error) {
	// Define the events we want to listen to
	events := []string{
		"push",
		"pull_request",
		"pull_request_review",
		"deployment",
		"deployment_status",
		"workflow_run",
		"check_suite",
		"issues",
		"issue_comment",
		"release",
	}

	webhookConfig := WebhookConfig{
		URL:         callbackURL,
		ContentType: "json",
		Secret:      wm.config.WebhookSecret,
		InsecureSSL: false,
		Events:      events,
		Active:      true,
	}

	return wm.CreateWebhook(owner, repo, webhookConfig)
}