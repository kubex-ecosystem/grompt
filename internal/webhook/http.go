// Package webhook implements HTTP handlers for webhook endpoints.
package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// HTTPHandler handles HTTP webhook requests
type HTTPHandler struct {
	handler *Handler
}

// NewHTTPHandler creates a new HTTP webhook handler
func NewHTTPHandler(handler *Handler) *HTTPHandler {
	return &HTTPHandler{
		handler: handler,
	}
}

// HandleWebhook processes incoming webhook HTTP requests
func (h *HTTPHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse webhook event
	var rawEvent map[string]interface{}
	if err := json.Unmarshal(body, &rawEvent); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Convert to internal event format
	event, err := h.parseWebhookEvent(r, rawEvent)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse event: %v", err), http.StatusBadRequest)
		return
	}

	// Process event asynchronously
	if err := h.handler.HandleEvent(r.Context(), event); err != nil {
		http.Error(w, fmt.Sprintf("Failed to process event: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]interface{}{
		"status":    "accepted",
		"event_id":  event.ID,
		"message":   "Event queued for processing",
		"timestamp": time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

// parseWebhookEvent converts raw webhook data to internal Event format
func (h *HTTPHandler) parseWebhookEvent(r *http.Request, rawEvent map[string]interface{}) (Event, error) {
	// Generate unique event ID
	eventID := uuid.New().String()

	// Detect webhook source and type
	source, eventType := h.detectWebhookSource(r, rawEvent)
	repository := h.extractRepository(rawEvent)

	// Determine analysis configuration based on event type
	analysisTypes, priority, expectedLatency := h.determineAnalysisConfig(eventType, rawEvent)

	event := Event{
		ID:         eventID,
		Type:       eventType,
		Source:     source,
		Repository: repository,
		Timestamp:  time.Now().UTC(),
		Payload:    rawEvent,
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

// detectWebhookSource identifies the source and type of webhook
func (h *HTTPHandler) detectWebhookSource(r *http.Request, payload map[string]interface{}) (string, string) {
	// GitHub webhook detection
	if githubEvent := r.Header.Get("X-GitHub-Event"); githubEvent != "" {
		return "github", h.mapGitHubEvent(githubEvent, payload)
	}

	// GitLab webhook detection
	if gitlabEvent := r.Header.Get("X-Gitlab-Event"); gitlabEvent != "" {
		return "gitlab", h.mapGitLabEvent(gitlabEvent, payload)
	}

	// Jenkins webhook detection
	if userAgent := r.Header.Get("User-Agent"); userAgent == "Jenkins" {
		return "jenkins", "deployment"
	}

	// Custom internal webhooks
	if r.Header.Get("X-Grompt-Source") == "internal" {
		if eventType, ok := payload["type"].(string); ok {
			return "internal", eventType
		}
	}

	// Default fallback
	return "unknown", "generic"
}

// mapGitHubEvent maps GitHub webhook events to internal event types
func (h *HTTPHandler) mapGitHubEvent(githubEvent string, payload map[string]interface{}) string {
	switch githubEvent {
	case "push":
		return "push"
	case "pull_request":
		// Check if PR was merged
		if pr, ok := payload["pull_request"].(map[string]interface{}); ok {
			if merged, ok := pr["merged"].(bool); ok && merged {
				return "pull_request_merged"
			}
		}
		return "pull_request"
	case "deployment":
		return "deployment"
	case "deployment_status":
		return "deployment_status"
	case "workflow_run":
		return "workflow_run"
	case "release":
		return "release"
	case "issues":
		return "issue"
	default:
		return "github_" + githubEvent
	}
}

// mapGitLabEvent maps GitLab webhook events to internal event types
func (h *HTTPHandler) mapGitLabEvent(gitlabEvent string, payload map[string]interface{}) string {
	switch gitlabEvent {
	case "Push Hook":
		return "push"
	case "Merge Request Hook":
		return "pull_request"
	case "Pipeline Hook":
		return "pipeline"
	case "Deployment Hook":
		return "deployment"
	case "Release Hook":
		return "release"
	default:
		return "gitlab_" + gitlabEvent
	}
}

// extractRepository extracts repository information from webhook payload
func (h *HTTPHandler) extractRepository(payload map[string]interface{}) string {
	// GitHub format
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

	// GitLab format
	if project, ok := payload["project"].(map[string]interface{}); ok {
		if pathWithNamespace, ok := project["path_with_namespace"].(string); ok {
			return pathWithNamespace
		}
	}

	// Jenkins or custom format
	if repo, ok := payload["repository"].(string); ok {
		return repo
	}

	return "unknown/repository"
}

// determineAnalysisConfig determines what analysis should be performed
func (h *HTTPHandler) determineAnalysisConfig(eventType string, payload map[string]interface{}) ([]string, string, string) {
	var analysisTypes []string
	var priority string
	var expectedLatency string

	switch eventType {
	case "push":
		analysisTypes = []string{"chi", "incremental_dora"}
		priority = "normal"
		expectedLatency = "minutes"

	case "pull_request_merged":
		analysisTypes = []string{"dora", "chi", "ai"}
		priority = "high"
		expectedLatency = "minutes"

	case "deployment":
		analysisTypes = []string{"dora", "executive"}
		priority = "high"
		expectedLatency = "instant"

	case "deployment_status":
		// Check if deployment failed
		if status := h.extractDeploymentStatus(payload); status == "failure" {
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
		if conclusion := h.extractWorkflowConclusion(payload); conclusion == "failure" {
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

	case "execution_result":
		// Meta-analysis event
		analysisTypes = []string{"execution_analysis", "meta_insights"}
		priority = "normal"
		expectedLatency = "minutes"

	default:
		// Generic analysis for unknown events
		analysisTypes = []string{"chi"}
		priority = "low"
		expectedLatency = "hours"
	}

	return analysisTypes, priority, expectedLatency
}

// extractDeploymentStatus extracts deployment status from webhook payload
func (h *HTTPHandler) extractDeploymentStatus(payload map[string]interface{}) string {
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
func (h *HTTPHandler) extractWorkflowConclusion(payload map[string]interface{}) string {
	if workflowRun, ok := payload["workflow_run"].(map[string]interface{}); ok {
		if conclusion, ok := workflowRun["conclusion"].(string); ok {
			return conclusion
		}
	}
	return "unknown"
}

// HealthCheck provides a simple health check endpoint
func (h *HTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"service":   "webhook-handler",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
