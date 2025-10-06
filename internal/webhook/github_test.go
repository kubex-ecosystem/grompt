package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
)

// MockHandler implements the Handler interface for testing
type MockHandler struct {
	lastEvent Event
	shouldErr bool
}

func (m *MockHandler) HandleEvent(ctx context.Context, event Event) error {
	m.lastEvent = event
	if m.shouldErr {
		return fmt.Errorf("mock error")
	}
	return nil
}

func TestGitHubHandler_HandleGitHubWebhook(t *testing.T) {
	secret := "test-webhook-secret"
	mockHandler := &MockHandler{}
	githubHandler := NewGitHubHandler(secret, mockHandler)

	// Create test payload
	payload := map[string]interface{}{
		"action": "opened",
		"pull_request": map[string]interface{}{
			"number": 123,
			"title":  "Test PR",
			"merged": false,
		},
		"repository": map[string]interface{}{
			"full_name": "test-owner/test-repo",
			"name":      "test-repo",
			"owner": map[string]interface{}{
				"login": "test-owner",
			},
		},
	}

	payloadBytes, _ := json.Marshal(payload)

	// Compute signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payloadBytes)
	signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	tests := []struct {
		name           string
		method         string
		event          string
		delivery       string
		signature      string
		payload        []byte
		expectedStatus int
		expectEvent    bool
	}{
		{
			name:           "valid pull request webhook",
			method:         "POST",
			event:          "pull_request",
			delivery:       "test-delivery-1",
			signature:      signature,
			payload:        payloadBytes,
			expectedStatus: 200,
			expectEvent:    true,
		},
		{
			name:           "invalid method",
			method:         "GET",
			event:          "pull_request",
			delivery:       "test-delivery-2",
			signature:      signature,
			payload:        payloadBytes,
			expectedStatus: 405,
			expectEvent:    false,
		},
		{
			name:           "invalid signature",
			method:         "POST",
			event:          "pull_request",
			delivery:       "test-delivery-3",
			signature:      "sha256=invalid",
			payload:        payloadBytes,
			expectedStatus: 401,
			expectEvent:    false,
		},
		{
			name:           "duplicate delivery",
			method:         "POST",
			event:          "pull_request",
			delivery:       "test-delivery-1", // Same as first test
			signature:      signature,
			payload:        payloadBytes,
			expectedStatus: 200, // Should return 200 for duplicates
			expectEvent:    false,
		},
		{
			name:           "invalid JSON",
			method:         "POST",
			event:          "pull_request",
			delivery:       "test-delivery-4",
			signature:      computeSignature([]byte("invalid json"), secret),
			payload:        []byte("invalid json"),
			expectedStatus: 400,
			expectEvent:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/webhook/github", bytes.NewReader(tt.payload))
			req.Header.Set("X-GitHub-Event", tt.event)
			req.Header.Set("X-GitHub-Delivery", tt.delivery)
			req.Header.Set("X-Hub-Signature-256", tt.signature)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			githubHandler.HandleGitHubWebhook(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectEvent && mockHandler.lastEvent.ID == "" {
				t.Error("Expected event to be processed but none was")
			}
			if !tt.expectEvent && tt.expectedStatus == 200 && mockHandler.lastEvent.ID != "" {
				// Reset for next test if this was a duplicate (which doesn't process)
				mockHandler.lastEvent = Event{}
			}
		})
	}
}

func TestGitHubHandler_MapGitHubEventType(t *testing.T) {
	handler := &GitHubHandler{}

	tests := []struct {
		name         string
		githubEvent  string
		payload      map[string]interface{}
		expectedType string
	}{
		{
			name:         "push event",
			githubEvent:  "push",
			payload:      map[string]interface{}{},
			expectedType: "push",
		},
		{
			name:        "pull request opened",
			githubEvent: "pull_request",
			payload: map[string]interface{}{
				"action": "opened",
				"pull_request": map[string]interface{}{
					"merged": false,
				},
			},
			expectedType: "pull_request_opened",
		},
		{
			name:        "pull request merged",
			githubEvent: "pull_request",
			payload: map[string]interface{}{
				"pull_request": map[string]interface{}{
					"merged": true,
				},
			},
			expectedType: "pull_request_merged",
		},
		{
			name:         "deployment",
			githubEvent:  "deployment",
			payload:      map[string]interface{}{},
			expectedType: "deployment",
		},
		{
			name:         "workflow run",
			githubEvent:  "workflow_run",
			payload:      map[string]interface{}{},
			expectedType: "workflow_run",
		},
		{
			name:        "issue opened",
			githubEvent: "issues",
			payload: map[string]interface{}{
				"action": "opened",
			},
			expectedType: "issue_opened",
		},
		{
			name:         "unknown event",
			githubEvent:  "unknown",
			payload:      map[string]interface{}{},
			expectedType: "github_unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.mapGitHubEventType(tt.githubEvent, tt.payload)
			if result != tt.expectedType {
				t.Errorf("Expected event type %s, got %s", tt.expectedType, result)
			}
		})
	}
}

func TestGitHubHandler_ExtractRepository(t *testing.T) {
	handler := &GitHubHandler{}

	tests := []struct {
		name     string
		payload  map[string]interface{}
		expected string
	}{
		{
			name: "valid repository with full_name",
			payload: map[string]interface{}{
				"repository": map[string]interface{}{
					"full_name": "owner/repo",
				},
			},
			expected: "owner/repo",
		},
		{
			name: "valid repository with owner and name",
			payload: map[string]interface{}{
				"repository": map[string]interface{}{
					"name": "repo",
					"owner": map[string]interface{}{
						"login": "owner",
					},
				},
			},
			expected: "owner/repo",
		},
		{
			name:     "missing repository",
			payload:  map[string]interface{}{},
			expected: "unknown/repository",
		},
		{
			name: "invalid repository structure",
			payload: map[string]interface{}{
				"repository": "invalid",
			},
			expected: "unknown/repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.extractRepository(tt.payload)
			if result != tt.expected {
				t.Errorf("Expected repository %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGitHubHandler_DetermineAnalysisConfig(t *testing.T) {
	handler := &GitHubHandler{}

	tests := []struct {
		name              string
		eventType         string
		payload           map[string]interface{}
		expectedAnalysis  []string
		expectedPriority  string
		expectedLatency   string
	}{
		{
			name:              "push event",
			eventType:         "push",
			payload:           map[string]interface{}{},
			expectedAnalysis:  []string{"chi", "incremental_dora", "ai"},
			expectedPriority:  "normal",
			expectedLatency:   "minutes",
		},
		{
			name:              "pull request merged",
			eventType:         "pull_request_merged",
			payload:           map[string]interface{}{},
			expectedAnalysis:  []string{"dora", "chi", "ai"},
			expectedPriority:  "high",
			expectedLatency:   "minutes",
		},
		{
			name:      "deployment failure",
			eventType: "deployment_status",
			payload: map[string]interface{}{
				"deployment_status": map[string]interface{}{
					"state": "failure",
				},
			},
			expectedAnalysis: []string{"dora", "incident_analysis"},
			expectedPriority: "critical",
			expectedLatency:  "instant",
		},
		{
			name:      "workflow failure",
			eventType: "workflow_run",
			payload: map[string]interface{}{
				"workflow_run": map[string]interface{}{
					"conclusion": "failure",
				},
			},
			expectedAnalysis: []string{"dora", "failure_analysis"},
			expectedPriority: "high",
			expectedLatency:  "instant",
		},
		{
			name:              "release",
			eventType:         "release",
			payload:           map[string]interface{}{},
			expectedAnalysis:  []string{"dora", "chi", "ai", "executive"},
			expectedPriority:  "high",
			expectedLatency:   "minutes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis, priority, latency := handler.determineAnalysisConfig(tt.eventType, tt.payload)

			if len(analysis) != len(tt.expectedAnalysis) {
				t.Errorf("Expected %d analysis types, got %d", len(tt.expectedAnalysis), len(analysis))
			}

			for i, expected := range tt.expectedAnalysis {
				if i >= len(analysis) || analysis[i] != expected {
					t.Errorf("Expected analysis[%d] = %s, got %s", i, expected, analysis[i])
				}
			}

			if priority != tt.expectedPriority {
				t.Errorf("Expected priority %s, got %s", tt.expectedPriority, priority)
			}

			if latency != tt.expectedLatency {
				t.Errorf("Expected latency %s, got %s", tt.expectedLatency, latency)
			}
		})
	}
}

func TestGitHubHandler_GetWebhookInfo(t *testing.T) {
	handler := &GitHubHandler{}

	req := httptest.NewRequest("POST", "/webhook", nil)
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-GitHub-Delivery", "test-delivery")
	req.Header.Set("X-Hub-Signature-256", "sha256=test")
	req.Header.Set("User-Agent", "GitHub-Hookshot/test")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Hook-ID", "123")

	info := handler.GetWebhookInfo(req)

	expected := map[string]string{
		"event":         "push",
		"delivery":      "test-delivery",
		"signature":     "sha256=test",
		"user_agent":    "GitHub-Hookshot/test",
		"content_type":  "application/json",
		"github_hook_id": "123",
	}

	for key, expectedValue := range expected {
		if info[key] != expectedValue {
			t.Errorf("Expected %s = %s, got %s", key, expectedValue, info[key])
		}
	}
}

func TestNewGitHubHandler(t *testing.T) {
	secret := "test-secret"
	mockHandler := &MockHandler{}

	githubHandler := NewGitHubHandler(secret, mockHandler)

	if githubHandler.validator == nil {
		t.Error("Expected validator to be initialized")
	}

	if githubHandler.handler == nil {
		t.Error("Expected handler to be set correctly")
	}
}

// Helper function to compute signature for tests
func computeSignature(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}