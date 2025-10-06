// Package webhook provides GitHub-specific webhook handling with signature validation.
package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/services/github"
)

// EventHandler interface for processing webhook events
type EventHandler interface {
	HandleEvent(ctx context.Context, event Event) error
}

// GitHubHandler handles GitHub webhooks with signature validation and deduplication
type GitHubHandler struct {
	validator *github.WebhookValidator
	handler   EventHandler // The event handler interface
}

// NewGitHubHandler creates a new GitHub webhook handler
func NewGitHubHandler(webhookSecret string, handler EventHandler) *GitHubHandler {
	return &GitHubHandler{
		validator: github.NewWebhookValidator(webhookSecret),
		handler:   handler,
	}
}

// HandleGitHubWebhook processes incoming GitHub webhooks with security validation
func (gh *GitHubHandler) HandleGitHubWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate signature and check for duplicates
	if err := gh.validator.ValidateAndProcess(r, body); err != nil {
		if strings.Contains(err.Error(), "duplicate delivery ID") {
			// For duplicates, return 200 to avoid GitHub retry
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  "ignored",
				"reason":  "duplicate",
				"message": "Webhook delivery already processed",
			})
			return
		}

		// For signature validation failures, return 401
		http.Error(w, fmt.Sprintf("Webhook validation failed: %v", err), http.StatusUnauthorized)
		return
	}

	// Parse the webhook payload
	var rawEvent map[string]interface{}
	if err := json.Unmarshal(body, &rawEvent); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Convert to internal event format
	event, err := gh.convertGitHubEvent(r, rawEvent)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse event: %v", err), http.StatusBadRequest)
		return
	}

	// Process the event using the existing handler
	if err := gh.handler.HandleEvent(r.Context(), event); err != nil {
		http.Error(w, fmt.Sprintf("Failed to process event: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response with delivery information
	deliveryID := r.Header.Get("X-GitHub-Delivery")
	response := map[string]interface{}{
		"status":     "accepted",
		"event_id":   event.ID,
		"delivery":   deliveryID,
		"event_type": event.Type,
		"message":    "GitHub webhook processed successfully",
		"timestamp":  time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// convertGitHubEvent converts GitHub webhook payload to internal Event format
func (gh *GitHubHandler) convertGitHubEvent(r *http.Request, payload map[string]interface{}) (Event, error) {
	githubEvent := r.Header.Get("X-GitHub-Event")
	deliveryID := r.Header.Get("X-GitHub-Delivery")

	// Use delivery ID as event ID if available, otherwise generate one
	eventID := deliveryID
	if eventID == "" {
		eventID = fmt.Sprintf("github-%d", time.Now().UnixNano())
	}

	// Extract repository information
	repository := gh.extractRepository(payload)

	// Map GitHub event to internal event type and determine analysis config
	eventType := gh.mapGitHubEventType(githubEvent, payload)
	analysisTypes, priority, expectedLatency := gh.determineAnalysisConfig(eventType, payload)

	event := Event{
		ID:         eventID,
		Type:       eventType,
		Source:     "github",
		Repository: repository,
		Timestamp:  time.Now().UTC(),
		Payload:    payload,
		Metadata: EventMetadata{
			TriggerLevel:    1, // First level trigger
			AnalysisTypes:   analysisTypes,
			Priority:        priority,
			RecursionDepth:  0,  // Starting depth
			ParentEventID:   "", // No parent for external events
			ExpectedLatency: expectedLatency,
		},
	}

	return event, nil
}

// mapGitHubEventType maps GitHub webhook events to internal event types
func (gh *GitHubHandler) mapGitHubEventType(githubEvent string, payload map[string]interface{}) string {
	switch githubEvent {
	case "push":
		return "push"
	case "pull_request":
		// Check if PR was merged
		if pr, ok := payload["pull_request"].(map[string]interface{}); ok {
			if merged, ok := pr["merged"].(bool); ok && merged {
				return "pull_request_merged"
			}
			// Check PR action
			if action, ok := payload["action"].(string); ok {
				return fmt.Sprintf("pull_request_%s", action)
			}
		}
		return "pull_request"
	case "pull_request_review":
		return "pull_request_review"
	case "deployment":
		return "deployment"
	case "deployment_status":
		return "deployment_status"
	case "workflow_run":
		return "workflow_run"
	case "check_suite":
		return "check_suite"
	case "release":
		return "release"
	case "issues":
		if action, ok := payload["action"].(string); ok {
			return fmt.Sprintf("issue_%s", action)
		}
		return "issue"
	case "issue_comment":
		return "issue_comment"
	default:
		return fmt.Sprintf("github_%s", githubEvent)
	}
}

// extractRepository extracts repository information from GitHub webhook payload
func (gh *GitHubHandler) extractRepository(payload map[string]interface{}) string {
	if repo, ok := payload["repository"].(map[string]interface{}); ok {
		if fullName, ok := repo["full_name"].(string); ok {
			return fullName
		}
		// Fallback to owner/name construction
		if owner, ok := repo["owner"].(map[string]interface{}); ok {
			if ownerName, ok := owner["login"].(string); ok {
				if repoName, ok := repo["name"].(string); ok {
					return fmt.Sprintf("%s/%s", ownerName, repoName)
				}
			}
		}
	}

	return "unknown/repository"
}

// determineAnalysisConfig determines what analysis should be performed for GitHub events
func (gh *GitHubHandler) determineAnalysisConfig(eventType string, payload map[string]interface{}) ([]string, string, string) {
	var analysisTypes []string
	var priority string
	var expectedLatency string

	switch eventType {
	case "push":
		analysisTypes = []string{"chi", "incremental_dora", "ai"}
		priority = "normal"
		expectedLatency = "minutes"

	case "pull_request_opened", "pull_request_synchronize":
		analysisTypes = []string{"chi", "ai"}
		priority = "normal"
		expectedLatency = "minutes"

	case "pull_request_merged":
		analysisTypes = []string{"dora", "chi", "ai"}
		priority = "high"
		expectedLatency = "minutes"

	case "pull_request_review":
		analysisTypes = []string{"dora"}
		priority = "normal"
		expectedLatency = "instant"

	case "deployment":
		analysisTypes = []string{"dora", "executive"}
		priority = "high"
		expectedLatency = "instant"

	case "deployment_status":
		// Check if deployment failed
		if status := gh.extractDeploymentStatus(payload); status == "failure" {
			analysisTypes = []string{"dora", "incident_analysis"}
			priority = "critical"
			expectedLatency = "instant"
		} else {
			analysisTypes = []string{"dora"}
			priority = "normal"
			expectedLatency = "minutes"
		}

	case "workflow_run":
		// Check if workflow failed
		if conclusion := gh.extractWorkflowConclusion(payload); conclusion == "failure" {
			analysisTypes = []string{"dora", "failure_analysis"}
			priority = "high"
			expectedLatency = "instant"
		} else if conclusion == "success" {
			analysisTypes = []string{"incremental_dora"}
			priority = "low"
			expectedLatency = "minutes"
		} else {
			analysisTypes = []string{"incremental_dora"}
			priority = "low"
			expectedLatency = "minutes"
		}

	case "check_suite":
		if conclusion := gh.extractCheckSuiteConclusion(payload); conclusion == "failure" {
			analysisTypes = []string{"dora", "failure_analysis"}
			priority = "high"
			expectedLatency = "instant"
		} else {
			analysisTypes = []string{"incremental_dora"}
			priority = "low"
			expectedLatency = "minutes"
		}

	case "release":
		analysisTypes = []string{"dora", "chi", "ai", "executive"}
		priority = "high"
		expectedLatency = "minutes"

	case "issue_opened", "issue_closed":
		analysisTypes = []string{"community"}
		priority = "low"
		expectedLatency = "hours"

	case "issue_comment":
		analysisTypes = []string{"community"}
		priority = "low"
		expectedLatency = "hours"

	default:
		// Generic analysis for unknown events
		analysisTypes = []string{"chi"}
		priority = "low"
		expectedLatency = "hours"
	}

	return analysisTypes, priority, expectedLatency
}

// extractDeploymentStatus extracts deployment status from webhook payload
func (gh *GitHubHandler) extractDeploymentStatus(payload map[string]interface{}) string {
	if deployment, ok := payload["deployment"].(map[string]interface{}); ok {
		if state, ok := deployment["state"].(string); ok {
			return state
		}
	}
	if deploymentStatus, ok := payload["deployment_status"].(map[string]interface{}); ok {
		if state, ok := deploymentStatus["state"].(string); ok {
			return state
		}
	}
	return "unknown"
}

// extractWorkflowConclusion extracts workflow conclusion from webhook payload
func (gh *GitHubHandler) extractWorkflowConclusion(payload map[string]interface{}) string {
	if workflowRun, ok := payload["workflow_run"].(map[string]interface{}); ok {
		if conclusion, ok := workflowRun["conclusion"].(string); ok {
			return conclusion
		}
	}
	return "unknown"
}

// extractCheckSuiteConclusion extracts check suite conclusion from webhook payload
func (gh *GitHubHandler) extractCheckSuiteConclusion(payload map[string]interface{}) string {
	if checkSuite, ok := payload["check_suite"].(map[string]interface{}); ok {
		if conclusion, ok := checkSuite["conclusion"].(string); ok {
			return conclusion
		}
	}
	return "unknown"
}

// GetWebhookInfo returns information about webhook delivery for debugging
func (gh *GitHubHandler) GetWebhookInfo(r *http.Request) map[string]string {
	return map[string]string{
		"event":         r.Header.Get("X-GitHub-Event"),
		"delivery":      r.Header.Get("X-GitHub-Delivery"),
		"signature":     r.Header.Get("X-Hub-Signature-256"),
		"user_agent":    r.Header.Get("User-Agent"),
		"content_type":  r.Header.Get("Content-Type"),
		"github_hook_id": r.Header.Get("X-GitHub-Hook-ID"),
	}
}