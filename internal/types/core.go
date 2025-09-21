// Package types defines core domain types for the repository intelligence platform.
package types

import "time"

// Repository represents a source code repository
type Repository struct {
	Owner         string    `json:"owner"`
	Name          string    `json:"name"`
	FullName      string    `json:"full_name"`
	CloneURL      string    `json:"clone_url"`
	DefaultBranch string    `json:"default_branch"`
	Language      string    `json:"language"`
	IsPrivate     bool      `json:"is_private"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// DORAMetrics - DevOps Research and Assessment metrics
type DORAMetrics struct {
	LeadTimeP95Hours        float64   `json:"lead_time_p95_hours"`
	DeploymentFrequencyWeek float64   `json:"deployment_frequency_per_week"`
	ChangeFailRatePercent   float64   `json:"change_fail_rate_pct"`
	MTTRHours               float64   `json:"mttr_hours"`
	Period                  int       `json:"period_days"`
	CalculatedAt            time.Time `json:"calculated_at"`
}

// CHIMetrics Index metrics
type CHIMetrics struct {
	Score                int       `json:"chi_score"` // 0-100
	DuplicationPercent   float64   `json:"duplication_pct"`
	CyclomaticComplexity float64   `json:"cyclomatic_avg"`
	TestCoverage         float64   `json:"test_coverage_pct"`
	MaintainabilityIndex float64   `json:"maintainability_index"`
	TechnicalDebt        float64   `json:"technical_debt_hours"`
	Period               int       `json:"period_days"`
	CalculatedAt         time.Time `json:"calculated_at"`
}

// AIMetrics Metrics - Human vs AI development analysis
type AIMetrics struct {
	HIR          float64   `json:"hir"` // Human Input Ratio (0.0-1.0)
	AAC          float64   `json:"aac"` // AI Assist Coverage (0.0-1.0)
	TPH          float64   `json:"tph"` // Throughput per Human-hour
	HumanHours   float64   `json:"human_hours"`
	AIHours      float64   `json:"ai_hours"`
	Period       int       `json:"period_days"`
	CalculatedAt time.Time `json:"calculated_at"`
}

// Scorecard combines all metrics for a repository
type Scorecard struct {
	SchemaVersion       string      `json:"schema_version"`
	Repository          Repository  `json:"repository"`
	DORA                DORAMetrics `json:"dora"`
	CHI                 CHIMetrics  `json:"chi"`
	AI                  AIMetrics   `json:"ai"`
	BusFactor           int         `json:"bus_factor"`
	FirstReviewP50Hours float64     `json:"first_review_p50_hours"`
	Confidence          Confidence  `json:"confidence"`
	GeneratedAt         time.Time   `json:"generated_at"`
}

// Confidence levels for metrics accuracy
type Confidence struct {
	DORA  float64 `json:"dora"`  // 0.0-1.0
	CHI   float64 `json:"chi"`   // 0.0-1.0
	AI    float64 `json:"ai"`    // 0.0-1.0
	Group float64 `json:"group"` // Overall confidence
}

// ExecutiveReport report types (P1-P4 prompts)
type ExecutiveReport struct {
	Summary      ExecutiveSummary `json:"summary"`
	TopFocus     []FocusArea      `json:"top_focus"`
	QuickWins    []QuickWin       `json:"quick_wins"`
	Risks        []Risk           `json:"risks"`
	CallToAction string           `json:"call_to_action"`
}

type ExecutiveSummary struct {
	Grade            string  `json:"grade"` // A, B, C, D, F
	CHI              int     `json:"chi"`
	LeadTimeP95Hours float64 `json:"lead_time_p95_hours"`
	DeploysPerWeek   float64 `json:"deploys_per_week"`
}

type FocusArea struct {
	Title      string  `json:"title"`
	Why        string  `json:"why"`
	KPI        string  `json:"kpi"`
	Target     string  `json:"target"`
	Confidence float64 `json:"confidence"`
}

type QuickWin struct {
	Action       string `json:"action"`
	Effort       string `json:"effort"` // S, M, L
	ExpectedGain string `json:"expected_gain"`
}

type Risk struct {
	Risk       string `json:"risk"`
	Mitigation string `json:"mitigation"`
}

// CodeHealthReport Code Health Deep Dive report
type CodeHealthReport struct {
	CHINow       int            `json:"chi_now"`
	Drivers      []CHIDriver    `json:"drivers"`
	RefactorPlan []RefactorStep `json:"refactor_plan"`
	Guardrails   []string       `json:"guardrails"`
	Milestones   []Milestone    `json:"milestones"`
}

type CHIDriver struct {
	Metric string  `json:"metric"` // mi|duplication_pct|cyclomatic_avg
	Value  float64 `json:"value"`
	Impact string  `json:"impact"` // high|med|low
}

type RefactorStep struct {
	Step    int      `json:"step"`
	Theme   string   `json:"theme"` // duplication|complexity|tests
	Actions []string `json:"actions"`
	KPI     string   `json:"kpi"`
	Target  string   `json:"target"`
}

type Milestone struct {
	InDays int    `json:"in_days"`
	Goal   string `json:"goal"`
}

// DORAReport DORA & Ops report
type DORAReport struct {
	LeadTimeP95Hours        float64        `json:"lead_time_p95_hours"`
	DeploymentFrequencyWeek float64        `json:"deployment_frequency_per_week"`
	ChangeFailRatePercent   float64        `json:"change_fail_rate_pct"`
	MTTRHours               float64        `json:"mttr_hours"`
	Bottlenecks             []Bottleneck   `json:"bottlenecks"`
	Playbook                []PlaybookItem `json:"playbook"`
	Experiments             []Experiment   `json:"experiments"`
}

type Bottleneck struct {
	Area     string `json:"area"` // review|pipeline|batch_size|release
	Evidence string `json:"evidence"`
}

type PlaybookItem struct {
	Name           string `json:"name"`
	Policy         string `json:"policy"`
	ExpectedEffect string `json:"expected_effect"`
}

type Experiment struct {
	AB           string `json:"A/B"`
	Metric       string `json:"metric"` // lead_time_p95|CFR|MTTR
	DurationDays int    `json:"duration_days"`
}

// CommunityReport Community & Bus Factor report
type CommunityReport struct {
	BusFactor         int              `json:"bus_factor"`
	OnboardingP50Days int              `json:"onboarding_p50_days"`
	Roadmap           []RoadmapItem    `json:"roadmap"`
	Visibility        []VisibilityItem `json:"visibility"`
}

type RoadmapItem struct {
	Item          string `json:"item"`
	Why           string `json:"why"`
	SuccessMetric string `json:"success_metric"`
}

type VisibilityItem struct {
	Asset  string `json:"asset"`
	KPI    string `json:"kpi"`
	Effort string `json:"effort"`
}

// AnalysisJob represents a scheduled or running repository analysis job
type AnalysisJob struct {
	ID           string                 `json:"id"`
	RepoURL      string                 `json:"repo_url"`
	AnalysisType string                 `json:"analysis_type"`
	Status       string                 `json:"status"` // "scheduled", "running", "completed", "failed"
	Progress     float64                `json:"progress"`
	CreatedAt    time.Time              `json:"created_at"`
	StartedAt    *time.Time             `json:"started_at,omitempty"`
	CompletedAt  *time.Time             `json:"completed_at,omitempty"`
	Error        string                 `json:"error,omitempty"`
	Results      map[string]interface{} `json:"results,omitempty"`
	ScheduledBy  string                 `json:"scheduled_by,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}
