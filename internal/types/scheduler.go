package types

import "time"

// AnalysisRequest represents a repository analysis task
type AnalysisRequest struct {
	ID          string                 `json:"id"`
	ProjectPath string                 `json:"project_path"`
	Type        string                 `json:"type"` // "full", "dora", "chi", "ai_metrics"
	Scheduled   bool                   `json:"scheduled"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
}

// OrchestrationTask represents a task for external tool coordination
type OrchestrationTask struct {
	ID        string                 `json:"id"`
	Tool      string                 `json:"tool"` // "lookatni", "grompt", "squad_agent"
	Action    string                 `json:"action"`
	ProjectID string                 `json:"project_id"`
	Payload   map[string]interface{} `json:"payload"`
	CreatedAt time.Time              `json:"created_at"`
}
