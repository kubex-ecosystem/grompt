// Package webhook implements the webhook receiver for meta-recursive analysis triggers.
package webhook

import (
	"context"
	"fmt"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/types"
)

// Event represents a webhook event that triggers meta-analysis
type Event struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`       // "push", "pull_request", "deployment", "metrics_update"
	Source     string                 `json:"source"`     // "github", "gitlab", "jenkins", "internal"
	Repository string                 `json:"repository"` // "owner/repo"
	Timestamp  time.Time              `json:"timestamp"`
	Payload    map[string]interface{} `json:"payload"`
	Metadata   EventMetadata          `json:"metadata"`
}

// EventMetadata contains meta-information about the event
type EventMetadata struct {
	TriggerLevel    int      `json:"trigger_level"`    // 1=direct, 2=meta, 3=meta-meta
	AnalysisTypes   []string `json:"analysis_types"`   // ["dora", "chi", "ai", "executive"]
	Priority        string   `json:"priority"`         // "low", "normal", "high", "critical"
	RecursionDepth  int      `json:"recursion_depth"`  // How deep in the meta-loop
	ParentEventID   string   `json:"parent_event_id"`  // For tracking causality
	ExpectedLatency string   `json:"expected_latency"` // "instant", "minutes", "hours"
}

// Handler processes webhook events and triggers meta-analysis
type Handler struct {
	eventQueue  EventQueue
	grompt    GromptActor
	recommender RecommenderActor
	executor    ExecutorActor
}

// EventQueue interface for background job processing
type EventQueue interface {
	Enqueue(ctx context.Context, event Event, priority int) error
	Process(ctx context.Context, handler func(Event) error) error
}

// GromptActor interface for the analysis component
type GromptActor interface {
	TriggerAnalysis(ctx context.Context, event Event) (*AnalysisResult, error)
}

// RecommenderActor interface for recommendation generation
type RecommenderActor interface {
	GenerateRecommendations(ctx context.Context, analysis AnalysisResult) (*RecommendationSet, error)
}

// ExecutorActor interface for executing recommendations
type ExecutorActor interface {
	ExecuteRecommendations(ctx context.Context, recommendations RecommendationSet) (*ExecutionResult, error)
}

// AnalysisResult represents the output of analysis
type AnalysisResult struct {
	EventID     string            `json:"event_id"`
	Repository  string            `json:"repository"`
	Scorecard   *types.Scorecard  `json:"scorecard"`
	Insights    []AnalysisInsight `json:"insights"`
	Metadata    AnalysisMetadata  `json:"metadata"`
	GeneratedAt time.Time         `json:"generated_at"`
}

// AnalysisInsight represents a specific insight from analysis
type AnalysisInsight struct {
	Type        string  `json:"type"`     // "trend", "anomaly", "recommendation"
	Severity    string  `json:"severity"` // "info", "warning", "critical"
	Category    string  `json:"category"` // "performance", "quality", "security"
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Confidence  float64 `json:"confidence"` // 0.0-1.0
	Impact      string  `json:"impact"`     // "low", "medium", "high"
	Effort      string  `json:"effort"`     // "S", "M", "L", "XL"
}

// AnalysisMetadata contains meta-information about the analysis
type AnalysisMetadata struct {
	ProcessingTimeMs int      `json:"processing_time_ms"`
	DataSources      []string `json:"data_sources"`
	Confidence       float64  `json:"confidence"`
	Completeness     float64  `json:"completeness"`
}

// RecommendationSet represents a set of actionable recommendations
type RecommendationSet struct {
	EventID         string                 `json:"event_id"`
	Repository      string                 `json:"repository"`
	Recommendations []Recommendation       `json:"recommendations"`
	Metadata        RecommendationMetadata `json:"metadata"`
	GeneratedAt     time.Time              `json:"generated_at"`
}

// Recommendation represents a single actionable recommendation
type Recommendation struct {
	ID           string               `json:"id"`
	Type         string               `json:"type"` // "code_fix", "process_change", "tool_adoption"
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	Action       RecommendationAction `json:"action"`
	Priority     string               `json:"priority"`     // "low", "normal", "high", "critical"
	Effort       string               `json:"effort"`       // "S", "M", "L", "XL"
	Impact       string               `json:"impact"`       // "low", "medium", "high"
	Confidence   float64              `json:"confidence"`   // 0.0-1.0
	Dependencies []string             `json:"dependencies"` // IDs of other recommendations
}

// RecommendationAction defines how to execute a recommendation
type RecommendationAction struct {
	Type       string                 `json:"type"`       // "pull_request", "config_change", "notification"
	Target     string                 `json:"target"`     // File, service, or endpoint to modify
	Parameters map[string]interface{} `json:"parameters"` // Action-specific parameters
	Validation string                 `json:"validation"` // How to verify success
}

// RecommendationMetadata contains meta-information about recommendations
type RecommendationMetadata struct {
	TotalRecommendations int     `json:"total_recommendations"`
	HighPriorityCount    int     `json:"high_priority_count"`
	EstimatedEffortHours float64 `json:"estimated_effort_hours"`
	ConfidenceScore      float64 `json:"confidence_score"`
}

// ExecutionResult represents the result of executing recommendations
type ExecutionResult struct {
	EventID         string            `json:"event_id"`
	Repository      string            `json:"repository"`
	ExecutedActions []ExecutedAction  `json:"executed_actions"`
	FailedActions   []FailedAction    `json:"failed_actions"`
	Metadata        ExecutionMetadata `json:"metadata"`
	CompletedAt     time.Time         `json:"completed_at"`
}

// ExecutedAction represents a successfully executed action
type ExecutedAction struct {
	RecommendationID string                 `json:"recommendation_id"`
	ActionType       string                 `json:"action_type"`
	Result           map[string]interface{} `json:"result"`
	ExecutedAt       time.Time              `json:"executed_at"`
}

// FailedAction represents a failed action execution
type FailedAction struct {
	RecommendationID string    `json:"recommendation_id"`
	ActionType       string    `json:"action_type"`
	Error            string    `json:"error"`
	RetryCount       int       `json:"retry_count"`
	FailedAt         time.Time `json:"failed_at"`
}

// ExecutionMetadata contains meta-information about execution
type ExecutionMetadata struct {
	SuccessRate       float64 `json:"success_rate"` // 0.0-1.0
	TotalActions      int     `json:"total_actions"`
	SuccessfulActions int     `json:"successful_actions"`
	FailedActions     int     `json:"failed_actions"`
	ExecutionTimeMs   int     `json:"execution_time_ms"`
}

// NewHandler creates a new webhook handler with meta-recursive capabilities
func NewHandler(queue EventQueue, grompt GromptActor, recommender RecommenderActor, executor ExecutorActor) *Handler {
	return &Handler{
		eventQueue:  queue,
		grompt:    grompt,
		recommender: recommender,
		executor:    executor,
	}
}

// HandleEvent processes an incoming webhook event
func (h *Handler) HandleEvent(ctx context.Context, event Event) error {
	// Determine priority based on event metadata
	priority := h.calculatePriority(event)

	// Enqueue event for asynchronous processing
	if err := h.eventQueue.Enqueue(ctx, event, priority); err != nil {
		return fmt.Errorf("failed to enqueue event: %w", err)
	}

	return nil
}

// ProcessEvent executes the meta-recursive analysis loop
func (h *Handler) ProcessEvent(ctx context.Context, event Event) error {
	// Step 1: Trigger Analysis
	analysisResult, err := h.grompt.TriggerAnalysis(ctx, event)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	// Step 2: Generate Recommendations
	recommendations, err := h.recommender.GenerateRecommendations(ctx, *analysisResult)
	if err != nil {
		return fmt.Errorf("recommendation generation failed: %w", err)
	}

	// Step 3: Execute Recommendations (if auto-execution is enabled)
	if h.shouldAutoExecute(event, *recommendations) {
		executionResult, err := h.executor.ExecuteRecommendations(ctx, *recommendations)
		if err != nil {
			return fmt.Errorf("recommendation execution failed: %w", err)
		}

		// Step 4: Meta-recursive trigger - analyze the execution results
		if h.shouldTriggerMetaAnalysis(executionResult) {
			metaEvent := h.createMetaEvent(event, *executionResult)
			return h.HandleEvent(ctx, metaEvent)
		}
	}

	return nil
}

// calculatePriority determines event priority based on metadata and content
func (h *Handler) calculatePriority(event Event) int {
	basePriority := 50 // Normal priority

	// Adjust based on event type
	switch event.Type {
	case "deployment":
		basePriority += 30 // High priority for deployments
	case "pull_request":
		basePriority += 10 // Medium-high for PRs
	case "push":
		basePriority += 5 // Slightly higher for commits
	}

	// Adjust based on metadata
	switch event.Metadata.Priority {
	case "critical":
		basePriority += 40
	case "high":
		basePriority += 20
	case "low":
		basePriority -= 20
	}

	// Increase priority for deeper recursion (they become more important)
	basePriority += event.Metadata.RecursionDepth * 10

	return basePriority
}

// shouldAutoExecute determines if recommendations should be executed automatically
func (h *Handler) shouldAutoExecute(event Event, recommendations RecommendationSet) bool {
	// Only auto-execute low-risk, high-confidence recommendations
	for _, rec := range recommendations.Recommendations {
		if rec.Priority == "critical" || rec.Impact == "high" {
			return false // Require human approval for high-impact changes
		}
		if rec.Confidence < 0.8 {
			return false // Require human approval for low-confidence recommendations
		}
	}

	// Auto-execute only for internal events and simple changes
	return event.Source == "internal" && recommendations.Metadata.TotalRecommendations <= 3
}

// shouldTriggerMetaAnalysis determines if execution results should trigger meta-analysis
func (h *Handler) shouldTriggerMetaAnalysis(result *ExecutionResult) bool {
	// Trigger meta-analysis if:
	// 1. Success rate is below threshold
	// 2. There are failed actions that need analysis
	// 3. Execution revealed new insights

	if result.Metadata.SuccessRate < 0.8 {
		return true
	}

	if len(result.FailedActions) > 0 {
		return true
	}

	// Always trigger meta-analysis for learning
	return true
}

// createMetaEvent creates a new event for meta-analysis
func (h *Handler) createMetaEvent(originalEvent Event, executionResult ExecutionResult) Event {
	return Event{
		ID:         fmt.Sprintf("meta-%s-%d", originalEvent.ID, time.Now().Unix()),
		Type:       "execution_result",
		Source:     "internal",
		Repository: originalEvent.Repository,
		Timestamp:  time.Now(),
		Payload: map[string]interface{}{
			"original_event":   originalEvent,
			"execution_result": executionResult,
		},
		Metadata: EventMetadata{
			TriggerLevel:    originalEvent.Metadata.TriggerLevel + 1,
			AnalysisTypes:   []string{"execution_analysis", "performance_analysis"},
			Priority:        "normal",
			RecursionDepth:  originalEvent.Metadata.RecursionDepth + 1,
			ParentEventID:   originalEvent.ID,
			ExpectedLatency: "minutes",
		},
	}
}
