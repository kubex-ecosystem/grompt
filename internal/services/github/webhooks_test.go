package github

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWebhookValidator_ValidateSignature(t *testing.T) {
	secret := "test-secret"
	validator := NewWebhookValidator(secret)

	payload := []byte(`{"test": "payload"}`)

	// Compute expected signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	tests := []struct {
		name        string
		payload     []byte
		signature   string
		expectError bool
	}{
		{
			name:        "valid signature",
			payload:     payload,
			signature:   expectedSignature,
			expectError: false,
		},
		{
			name:        "invalid signature",
			payload:     payload,
			signature:   "sha256=invalid",
			expectError: true,
		},
		{
			name:        "missing sha256 prefix",
			payload:     payload,
			signature:   hex.EncodeToString(mac.Sum(nil)),
			expectError: true,
		},
		{
			name:        "empty signature",
			payload:     payload,
			signature:   "",
			expectError: true,
		},
		{
			name:        "wrong payload",
			payload:     []byte(`{"different": "payload"}`),
			signature:   expectedSignature,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateSignature(tt.payload, tt.signature)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestWebhookValidator_DeduplicationTracking(t *testing.T) {
	validator := NewWebhookValidator("test-secret")

	deliveryID := "test-delivery-123"

	// Initially should not be duplicate
	if validator.IsDuplicate(deliveryID) {
		t.Error("Expected delivery ID to not be duplicate initially")
	}

	// Mark as processed
	validator.MarkDelivery(deliveryID)

	// Now should be duplicate
	if !validator.IsDuplicate(deliveryID) {
		t.Error("Expected delivery ID to be duplicate after marking")
	}

	// Test empty delivery ID handling
	if validator.IsDuplicate("") {
		t.Error("Expected empty delivery ID to not be duplicate")
	}

	validator.MarkDelivery("")
	if validator.IsDuplicate("") {
		t.Error("Expected empty delivery ID to remain not duplicate")
	}
}

func TestWebhookValidator_ValidateAndProcess(t *testing.T) {
	secret := "test-secret"
	validator := NewWebhookValidator(secret)

	payload := []byte(`{"test": "payload"}`)

	// Compute valid signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	validSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	tests := []struct {
		name        string
		signature   string
		deliveryID  string
		expectError bool
		errorType   string
	}{
		{
			name:        "valid first delivery",
			signature:   validSignature,
			deliveryID:  "delivery-1",
			expectError: false,
		},
		{
			name:        "duplicate delivery",
			signature:   validSignature,
			deliveryID:  "delivery-1", // Same as above
			expectError: true,
			errorType:   "duplicate",
		},
		{
			name:        "invalid signature",
			signature:   "sha256=invalid",
			deliveryID:  "delivery-2",
			expectError: true,
			errorType:   "signature",
		},
		{
			name:        "valid new delivery",
			signature:   validSignature,
			deliveryID:  "delivery-3",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock request
			req := httptest.NewRequest("POST", "/webhook", nil)
			req.Header.Set("X-Hub-Signature-256", tt.signature)
			req.Header.Set("X-GitHub-Delivery", tt.deliveryID)

			err := validator.ValidateAndProcess(req, payload)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if tt.expectError && tt.errorType == "duplicate" {
				if err.Error() != "duplicate delivery ID: "+tt.deliveryID {
					t.Errorf("Expected duplicate error, got: %v", err)
				}
			}
		})
	}
}

func TestWebhookValidator_CleanupOldDeliveries(t *testing.T) {
	validator := NewWebhookValidator("test-secret")
	validator.maxAge = 100 * time.Millisecond // Very short for testing

	// Add old delivery
	validator.MarkDelivery("old-delivery")

	// Wait for it to expire
	time.Sleep(150 * time.Millisecond)

	// Add recent delivery (this will trigger cleanup)
	validator.MarkDelivery("recent-delivery")

	// Check results
	validator.deliveryMutex.RLock()
	deliveryCount := len(validator.deliveries)
	_, oldExists := validator.deliveries["old-delivery"]
	_, recentExists := validator.deliveries["recent-delivery"]
	validator.deliveryMutex.RUnlock()

	if oldExists {
		t.Error("Expected old delivery to be cleaned up")
	}
	if !recentExists {
		t.Error("Expected recent delivery to still exist")
	}
	if deliveryCount != 1 {
		t.Errorf("Expected 1 delivery, got %d", deliveryCount)
	}
}

func TestNewWebhookValidator(t *testing.T) {
	secret := "test-secret"
	validator := NewWebhookValidator(secret)

	if validator.secret != secret {
		t.Errorf("Expected secret %s, got %s", secret, validator.secret)
	}

	if validator.deliveries == nil {
		t.Error("Expected deliveries map to be initialized")
	}

	if validator.maxAge != 24*time.Hour {
		t.Errorf("Expected maxAge to be 24 hours, got %v", validator.maxAge)
	}
}

func TestWebhookValidator_NoSecret(t *testing.T) {
	validator := NewWebhookValidator("")
	payload := []byte(`{"test": "payload"}`)

	err := validator.ValidateSignature(payload, "sha256=anything")
	if err == nil {
		t.Error("Expected error when secret is not configured")
	}

	if err.Error() != "webhook secret not configured" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}