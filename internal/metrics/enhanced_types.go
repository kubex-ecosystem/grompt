// Package metrics - Enhanced types for mature DORA/CHI/HIR metrics with timezone and caching support
package metrics

import (
	"time"

	"github.com/kubex-ecosystem/grompt/internal/types"
)

// TimeRange represents a time period with timezone support
type TimeRange struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Timezone string    `json:"timezone"` // IANA timezone (e.g., "America/New_York")
}

// MetricsRequest represents a request for metrics calculation
type MetricsRequest struct {
	Repository  types.Repository `json:"repository"`
	TimeRange   TimeRange        `json:"time_range"`
	Granularity string           `json:"granularity"` // "hour", "day", "week", "month"
	UseCache    bool             `json:"use_cache"`
	CacheTTL    time.Duration    `json:"cache_ttl"`
}

// Enhanced metrics with timezone and aggregation support

// EnhancedDORAMetrics extends DORAMetrics with timezone and granularity support
type EnhancedDORAMetrics struct {
	types.DORAMetrics
	TimeRange              TimeRange                `json:"time_range"`
	Granularity            string                   `json:"granularity"`
	Timezone               string                   `json:"timezone"`
	IncidentCount          int                      `json:"incident_count"`
	FailedDeployments      int                      `json:"failed_deployments"`
	TotalDeployments       int                      `json:"total_deployments"`
	MeanLeadTimeHours      float64                  `json:"mean_lead_time_hours"`
	MedianLeadTimeHours    float64                  `json:"median_lead_time_hours"`
	TimeSeries             []DORATimeSeriesPoint    `json:"time_series,omitempty"`
	IncidentBreakdown      []IncidentClassification `json:"incident_breakdown,omitempty"`
	DeploymentTrends       []DeploymentTrend        `json:"deployment_trends,omitempty"`
	Confidence             float64                  `json:"confidence"`
	DataQuality            DataQuality              `json:"data_quality"`
	CacheInfo              CacheInfo                `json:"cache_info"`
}

// EnhancedCHIMetrics extends CHIMetrics with detailed analysis
type EnhancedCHIMetrics struct {
	types.CHIMetrics
	TimeRange           TimeRange             `json:"time_range"`
	Granularity         string                `json:"granularity"`
	Timezone            string                `json:"timezone"`
	FileMetrics         []FileMetric          `json:"file_metrics,omitempty"`
	LanguageBreakdown   []LanguageMetric      `json:"language_breakdown,omitempty"`
	ComplexityHotspots  []ComplexityHotspot   `json:"complexity_hotspots,omitempty"`
	TechnicalDebtItems  []TechnicalDebtItem   `json:"technical_debt_items,omitempty"`
	TestCoverageDetail  TestCoverageDetail    `json:"test_coverage_detail"`
	Trends              CHITrendAnalysis      `json:"trends"`
	Confidence          float64               `json:"confidence"`
	DataQuality         DataQuality           `json:"data_quality"`
	CacheInfo           CacheInfo             `json:"cache_info"`
}

// EnhancedAIMetrics extends AIMetrics with detailed AI assistance analysis
type EnhancedAIMetrics struct {
	types.AIMetrics
	TimeRange              TimeRange               `json:"time_range"`
	Granularity            string                  `json:"granularity"`
	Timezone               string                  `json:"timezone"`
	AIToolBreakdown        []AIToolUsage           `json:"ai_tool_breakdown,omitempty"`
	ProductivityMetrics    ProductivityMetrics     `json:"productivity_metrics"`
	CodeQualityImpact      CodeQualityImpact       `json:"code_quality_impact"`
	DeveloperEfficiency    DeveloperEfficiency     `json:"developer_efficiency"`
	AIAssistanceTimeline   []AIAssistancePoint     `json:"ai_assistance_timeline,omitempty"`
	HumanVsAIContributions HumanVsAIContributions `json:"human_vs_ai_contributions"`
	Confidence             float64                 `json:"confidence"`
	DataQuality            DataQuality             `json:"data_quality"`
	CacheInfo              CacheInfo               `json:"cache_info"`
}

// Time series and breakdown types

// DORATimeSeriesPoint represents a point in DORA metrics time series
type DORATimeSeriesPoint struct {
	Timestamp           time.Time `json:"timestamp"`
	LeadTimeHours       float64   `json:"lead_time_hours"`
	DeploymentCount     int       `json:"deployment_count"`
	FailureCount        int       `json:"failure_count"`
	RecoveryTimeHours   float64   `json:"recovery_time_hours"`
	ChangeFailureRate   float64   `json:"change_failure_rate"`
}

// IncidentClassification represents incident analysis
type IncidentClassification struct {
	Type                string        `json:"type"` // "deployment_failure", "hotfix", "rollback", "outage"
	Severity            string        `json:"severity"` // "critical", "high", "medium", "low"
	Count               int           `json:"count"`
	MeanResolutionTime  time.Duration `json:"mean_resolution_time"`
	TotalDowntimeHours  float64       `json:"total_downtime_hours"`
	AffectedDeployments []string      `json:"affected_deployments,omitempty"`
}

// DeploymentTrend represents deployment frequency trends
type DeploymentTrend struct {
	Period           string  `json:"period"` // "week", "month", "quarter"
	DeploymentCount  int     `json:"deployment_count"`
	SuccessRate      float64 `json:"success_rate"`
	AverageLeadTime  float64 `json:"average_lead_time_hours"`
	TrendDirection   string  `json:"trend_direction"` // "improving", "stable", "declining"
}

// CHI breakdown types

// FileMetric represents metrics for a single file
type FileMetric struct {
	Path                 string  `json:"path"`
	Language             string  `json:"language"`
	LinesOfCode          int     `json:"lines_of_code"`
	CyclomaticComplexity int     `json:"cyclomatic_complexity"`
	TestCoverage         float64 `json:"test_coverage"`
	DuplicationScore     float64 `json:"duplication_score"`
	MaintainabilityIndex float64 `json:"maintainability_index"`
	TechnicalDebtHours   float64 `json:"technical_debt_hours"`
	LastModified         time.Time `json:"last_modified"`
}

// LanguageMetric represents metrics per programming language
type LanguageMetric struct {
	Language             string  `json:"language"`
	FileCount            int     `json:"file_count"`
	TotalLinesOfCode     int     `json:"total_lines_of_code"`
	AverageComplexity    float64 `json:"average_complexity"`
	TestCoverage         float64 `json:"test_coverage"`
	TechnicalDebtHours   float64 `json:"technical_debt_hours"`
}

// ComplexityHotspot represents high-complexity code areas
type ComplexityHotspot struct {
	File                 string  `json:"file"`
	Function             string  `json:"function"`
	CyclomaticComplexity int     `json:"cyclomatic_complexity"`
	LinesOfCode          int     `json:"lines_of_code"`
	EstimatedRefactorHours float64 `json:"estimated_refactor_hours"`
	Priority             string  `json:"priority"` // "critical", "high", "medium", "low"
}

// TechnicalDebtItem represents a specific technical debt item
type TechnicalDebtItem struct {
	Type                 string  `json:"type"` // "complexity", "duplication", "test_coverage", "maintainability"
	Description          string  `json:"description"`
	Location             string  `json:"location"`
	EstimatedEffortHours float64 `json:"estimated_effort_hours"`
	ImpactLevel          string  `json:"impact_level"` // "critical", "high", "medium", "low"
	RecommendedAction    string  `json:"recommended_action"`
}

// TestCoverageDetail represents detailed test coverage analysis
type TestCoverageDetail struct {
	LinesCovered         int     `json:"lines_covered"`
	LinesTotal           int     `json:"lines_total"`
	BranchesCovered      int     `json:"branches_covered"`
	BranchesTotal        int     `json:"branches_total"`
	FunctionsCovered     int     `json:"functions_covered"`
	FunctionsTotal       int     `json:"functions_total"`
	UncoveredFiles       []string `json:"uncovered_files,omitempty"`
	TestFileCount        int     `json:"test_file_count"`
	TestToCodeRatio      float64 `json:"test_to_code_ratio"`
}

// CHITrendAnalysis represents CHI trends over time
type CHITrendAnalysis struct {
	ScoreTrend           string  `json:"score_trend"` // "improving", "stable", "declining"
	ComplexityTrend      string  `json:"complexity_trend"`
	TestCoverageTrend    string  `json:"test_coverage_trend"`
	TechnicalDebtTrend   string  `json:"technical_debt_trend"`
	MonthlyScoreChange   float64 `json:"monthly_score_change"`
	RecommendedActions   []string `json:"recommended_actions,omitempty"`
}

// AI metrics breakdown types

// AIToolUsage represents usage of specific AI tools
type AIToolUsage struct {
	ToolName            string  `json:"tool_name"` // "copilot", "chatgpt", "codeium", etc.
	UsageHours          float64 `json:"usage_hours"`
	AcceptanceRate      float64 `json:"acceptance_rate"`
	LinesGenerated      int     `json:"lines_generated"`
	LinesAccepted       int     `json:"lines_accepted"`
	CodeQualityScore    float64 `json:"code_quality_score"`
	ProductivityBoost   float64 `json:"productivity_boost"` // Percentage increase
}

// ProductivityMetrics represents productivity analysis
type ProductivityMetrics struct {
	CommitsPerHour          float64 `json:"commits_per_hour"`
	LinesPerHour            float64 `json:"lines_per_hour"`
	FeaturesPerSprint       float64 `json:"features_per_sprint"`
	BugsPerFeature          float64 `json:"bugs_per_feature"`
	TimeToFirstReview       float64 `json:"time_to_first_review_hours"`
	CodeReviewCycles        float64 `json:"code_review_cycles"`
	HumanOnlyProductivity   float64 `json:"human_only_productivity"`
	AIAssistedProductivity  float64 `json:"ai_assisted_productivity"`
	ProductivityImprovement float64 `json:"productivity_improvement"` // Percentage
}

// CodeQualityImpact represents AI impact on code quality
type CodeQualityImpact struct {
	BugDensityReduction     float64 `json:"bug_density_reduction"`
	TestCoverageImprovement float64 `json:"test_coverage_improvement"`
	CodeComplexityChange    float64 `json:"code_complexity_change"`
	RefactoringFrequency    float64 `json:"refactoring_frequency"`
	CodeReviewPassRate      float64 `json:"code_review_pass_rate"`
	SecurityVulnerabilities int     `json:"security_vulnerabilities"`
}

// DeveloperEfficiency represents developer efficiency metrics
type DeveloperEfficiency struct {
	FocusTimeHours          float64 `json:"focus_time_hours"`
	InterruptionFrequency   float64 `json:"interruption_frequency"`
	ContextSwitchingPenalty float64 `json:"context_switching_penalty"`
	FlowStateAchievement    float64 `json:"flow_state_achievement"`
	LearningCurveReduction  float64 `json:"learning_curve_reduction"`
	OnboardingTimeReduction float64 `json:"onboarding_time_reduction"`
}

// AIAssistancePoint represents AI assistance over time
type AIAssistancePoint struct {
	Timestamp       time.Time `json:"timestamp"`
	HIR             float64   `json:"hir"`
	AAC             float64   `json:"aac"`
	TPH             float64   `json:"tph"`
	ActiveAITools   []string  `json:"active_ai_tools"`
	ProductivityIndex float64 `json:"productivity_index"`
}

// HumanVsAIContributions represents the breakdown of human vs AI contributions
type HumanVsAIContributions struct {
	HumanCommits        int     `json:"human_commits"`
	AIAssistedCommits   int     `json:"ai_assisted_commits"`
	HumanLinesAdded     int     `json:"human_lines_added"`
	AILinesAdded        int     `json:"ai_lines_added"`
	HumanTestsWritten   int     `json:"human_tests_written"`
	AITestsWritten      int     `json:"ai_tests_written"`
	HumanBugsFixed      int     `json:"human_bugs_fixed"`
	AIBugsFixed         int     `json:"ai_bugs_fixed"`
	CollaborationScore  float64 `json:"collaboration_score"` // How well human and AI work together
}

// Common support types

// DataQuality represents data quality metrics
type DataQuality struct {
	Completeness    float64 `json:"completeness"` // 0.0-1.0
	Accuracy        float64 `json:"accuracy"`     // 0.0-1.0
	Timeliness      float64 `json:"timeliness"`   // 0.0-1.0
	Consistency     float64 `json:"consistency"`  // 0.0-1.0
	DataPoints      int     `json:"data_points"`
	MissingData     int     `json:"missing_data"`
	QualityWarnings []string `json:"quality_warnings,omitempty"`
}

// CacheInfo represents cache metadata
type CacheInfo struct {
	CacheHit       bool          `json:"cache_hit"`
	CacheKey       string        `json:"cache_key,omitempty"`
	CachedAt       *time.Time    `json:"cached_at,omitempty"`
	TTL            time.Duration `json:"ttl,omitempty"`
	ExpiresAt      *time.Time    `json:"expires_at,omitempty"`
	ComputeTimeMs  int64         `json:"compute_time_ms"`
	DataSources    []string      `json:"data_sources,omitempty"`
}

// AggregatedMetrics represents cross-repository aggregated metrics
type AggregatedMetrics struct {
	Repositories         []types.Repository    `json:"repositories"`
	TimeRange            TimeRange             `json:"time_range"`
	AggregatedDORA       AggregatedDORAMetrics `json:"aggregated_dora"`
	AggregatedCHI        AggregatedCHIMetrics  `json:"aggregated_chi"`
	AggregatedAI         AggregatedAIMetrics   `json:"aggregated_ai"`
	CrossRepoInsights    CrossRepoInsights     `json:"cross_repo_insights"`
	OrganizationalHealth OrganizationalHealth  `json:"organizational_health"`
	CacheInfo            CacheInfo             `json:"cache_info"`
}

// AggregatedDORAMetrics represents DORA metrics across repositories
type AggregatedDORAMetrics struct {
	MeanLeadTimeP95Hours        float64 `json:"mean_lead_time_p95_hours"`
	MeanDeploymentFrequencyWeek float64 `json:"mean_deployment_frequency_per_week"`
	MeanChangeFailRatePercent   float64 `json:"mean_change_fail_rate_pct"`
	MeanMTTRHours               float64 `json:"mean_mttr_hours"`
	TotalDeployments            int     `json:"total_deployments"`
	TotalIncidents              int     `json:"total_incidents"`
	BestPerformingRepo          string  `json:"best_performing_repo"`
	WorstPerformingRepo         string  `json:"worst_performing_repo"`
	Percentiles                 DORAPercentiles `json:"percentiles"`
}

// AggregatedCHIMetrics represents CHI metrics across repositories
type AggregatedCHIMetrics struct {
	MeanCHIScore                int     `json:"mean_chi_score"`
	MeanDuplicationPercent      float64 `json:"mean_duplication_pct"`
	MeanCyclomaticComplexity    float64 `json:"mean_cyclomatic_avg"`
	MeanTestCoverage            float64 `json:"mean_test_coverage_pct"`
	MeanMaintainabilityIndex    float64 `json:"mean_maintainability_index"`
	TotalTechnicalDebtHours     float64 `json:"total_technical_debt_hours"`
	HealthiestRepo              string  `json:"healthiest_repo"`
	MostTechnicalDebtRepo       string  `json:"most_technical_debt_repo"`
	LanguageHealthBreakdown     []LanguageHealthMetric `json:"language_health_breakdown"`
}

// AggregatedAIMetrics represents AI metrics across repositories
type AggregatedAIMetrics struct {
	MeanHIR                     float64 `json:"mean_hir"`
	MeanAAC                     float64 `json:"mean_aac"`
	MeanTPH                     float64 `json:"mean_tph"`
	TotalHumanHours             float64 `json:"total_human_hours"`
	TotalAIHours                float64 `json:"total_ai_hours"`
	MostAIAssistedRepo          string  `json:"most_ai_assisted_repo"`
	LeastAIAssistedRepo         string  `json:"least_ai_assisted_repo"`
	OrganizationalAIAdoption    float64 `json:"organizational_ai_adoption"`
	AverageProductivityBoost    float64 `json:"average_productivity_boost"`
}

// Supporting aggregated types

// DORAPercentiles represents percentile analysis of DORA metrics
type DORAPercentiles struct {
	LeadTimeP50        float64 `json:"lead_time_p50"`
	LeadTimeP75        float64 `json:"lead_time_p75"`
	LeadTimeP90        float64 `json:"lead_time_p90"`
	LeadTimeP95        float64 `json:"lead_time_p95"`
	DeployFreqP50      float64 `json:"deploy_freq_p50"`
	DeployFreqP75      float64 `json:"deploy_freq_p75"`
	DeployFreqP90      float64 `json:"deploy_freq_p90"`
	ChangeFailRateP50  float64 `json:"change_fail_rate_p50"`
	ChangeFailRateP75  float64 `json:"change_fail_rate_p75"`
	MTTRP50            float64 `json:"mttr_p50"`
	MTTRP75            float64 `json:"mttr_p75"`
}

// LanguageHealthMetric represents health metrics per language across repos
type LanguageHealthMetric struct {
	Language                string  `json:"language"`
	RepositoryCount         int     `json:"repository_count"`
	AverageCHIScore         int     `json:"average_chi_score"`
	AverageComplexity       float64 `json:"average_complexity"`
	AverageTestCoverage     float64 `json:"average_test_coverage"`
	TotalTechnicalDebtHours float64 `json:"total_technical_debt_hours"`
	HealthRanking           int     `json:"health_ranking"`
}

// CrossRepoInsights represents insights across repositories
type CrossRepoInsights struct {
	CommonPatterns          []string `json:"common_patterns"`
	SharedTechnicalDebt     []string `json:"shared_technical_debt"`
	BestPracticesSharing    []string `json:"best_practices_sharing"`
	KnowledgeTransferOpps   []string `json:"knowledge_transfer_opportunities"`
	StandardizationOpps     []string `json:"standardization_opportunities"`
	CollaborationHotspots   []string `json:"collaboration_hotspots"`
}

// OrganizationalHealth represents organization-wide health metrics
type OrganizationalHealth struct {
	DeliveryMaturity        string  `json:"delivery_maturity"` // "elite", "high", "medium", "low"
	CodeHealthMaturity      string  `json:"code_health_maturity"`
	AIAdoptionMaturity      string  `json:"ai_adoption_maturity"`
	DevExperienceScore      float64 `json:"dev_experience_score"`
	InnovationIndex         float64 `json:"innovation_index"`
	ScalingReadiness        float64 `json:"scaling_readiness"`
	TalentRetentionRisk     string  `json:"talent_retention_risk"`
	CompetitiveAdvantage    string  `json:"competitive_advantage"`
	StrategicRecommendations []string `json:"strategic_recommendations"`
}