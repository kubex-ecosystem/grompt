// Package integration provides GoBE backend integration for Repository Intelligence
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/types"
)

// GoBeClient handles communication with GoBE backend services
type GoBeClient struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

// NewGoBeClient creates a new GoBE client
func NewGoBeClient(baseURL, apiKey string) *GoBeClient {
	return &GoBeClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// RepositoryIntelligenceRequest represents a request for repository analysis
type RepositoryIntelligenceRequest struct {
	RepoURL        string                 `json:"repo_url"`
	AnalysisType   string                 `json:"analysis_type"`
	Configuration  map[string]interface{} `json:"configuration,omitempty"`
	ScheduledBy    string                 `json:"scheduled_by,omitempty"`
	NotifyChannels []string               `json:"notify_channels,omitempty"`
}

// NotificationRequest represents a notification to be sent
type NotificationRequest struct {
	Type        string                 `json:"type"` // "discord", "email", "webhook"
	Recipients  []string               `json:"recipients"`
	Subject     string                 `json:"subject"`
	Message     string                 `json:"message"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Priority    string                 `json:"priority,omitempty"` // "low", "normal", "high", "urgent"
	AttachFiles []string               `json:"attach_files,omitempty"`
}

// AgentRegistration represents grompt registration in AI squad
type AgentRegistration struct {
	Name         string            `json:"name"`
	Type         string            `json:"type"` // "grompt", "improver", "orchestrator"
	Capabilities []string          `json:"capabilities"`
	Endpoints    map[string]string `json:"endpoints"`
	Config       AgentConfig       `json:"config"`
}

// AgentConfig represents agent configuration
type AgentConfig struct {
	AutoSchedule  bool                   `json:"auto_schedule"`
	ScheduleCron  string                 `json:"schedule_cron,omitempty"`
	RetryPolicy   RetryPolicy            `json:"retry_policy"`
	Notifications NotificationConfig     `json:"notifications"`
	Integrations  map[string]interface{} `json:"integrations"`
}

// RetryPolicy defines retry behavior for failed operations
type RetryPolicy struct {
	MaxRetries    int           `json:"max_retries"`
	BackoffPolicy string        `json:"backoff_policy"` // "linear", "exponential"
	InitialDelay  time.Duration `json:"initial_delay"`
	MaxDelay      time.Duration `json:"max_delay"`
}

// NotificationConfig defines notification preferences
type NotificationConfig struct {
	OnSuccess       []string `json:"on_success,omitempty"`
	OnFailure       []string `json:"on_failure"`
	OnScheduled     []string `json:"on_scheduled,omitempty"`
	DiscordWebhook  string   `json:"discord_webhook,omitempty"`
	EmailRecipients []string `json:"email_recipients,omitempty"`
}

// SquadStatus represents current AI squad status
type SquadStatus struct {
	ActiveAgents int           `json:"active_agents"`
	RunningJobs  int           `json:"running_jobs"`
	QueuedJobs   int           `json:"queued_jobs"`
	Agents       []AgentStatus `json:"agents"`
	SystemHealth string        `json:"system_health"`
	LastUpdate   time.Time     `json:"last_update"`
}

// AgentStatus represents individual agent status
type AgentStatus struct {
	Name          string             `json:"name"`
	Status        string             `json:"status"` // "active", "idle", "error", "maintenance"
	LastActivity  time.Time          `json:"last_activity"`
	ProcessedJobs int                `json:"processed_jobs"`
	ErrorCount    int                `json:"error_count"`
	CurrentJob    *types.AnalysisJob `json:"current_job,omitempty"`
}

// ScheduleAnalysis schedules a repository analysis in GoBE
func (c *GoBeClient) ScheduleAnalysis(ctx context.Context, req RepositoryIntelligenceRequest) (*types.AnalysisJob, error) {
	url := fmt.Sprintf("%s/api/v1/mcp/grompt/schedule", c.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var job types.AnalysisJob
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &job, nil
}

// SendNotification sends notification through GoBE notification system
func (c *GoBeClient) SendNotification(ctx context.Context, notification NotificationRequest) error {
	url := fmt.Sprintf("%s/api/v1/notifications/send", c.baseURL)

	jsonData, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("notification failed with status: %d", resp.StatusCode)
	}

	return nil
}

// RegisterAgent registers grompt as an AI agent in GoBE squad system
func (c *GoBeClient) RegisterAgent(ctx context.Context, agent AgentRegistration) error {
	url := fmt.Sprintf("%s/api/v1/mcp/squad/register", c.baseURL)

	jsonData, err := json.Marshal(agent)
	if err != nil {
		return fmt.Errorf("failed to marshal agent registration: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to register agent: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("agent registration failed with status: %d", resp.StatusCode)
	}

	return nil
}

// GetSquadStatus retrieves current AI squad status from GoBE
func (c *GoBeClient) GetSquadStatus(ctx context.Context) (*SquadStatus, error) {
	url := fmt.Sprintf("%s/api/v1/mcp/squad/status", c.baseURL)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get squad status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get squad status: %d", resp.StatusCode)
	}

	var status SquadStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode squad status: %w", err)
	}

	return &status, nil
}
