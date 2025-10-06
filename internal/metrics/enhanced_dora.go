// Package metrics - Enhanced DORA calculator with timezone, caching, and GraphQL support
package metrics

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/types"
)

// EnhancedDORACalculator calculates DORA metrics with advanced features
type EnhancedDORACalculator struct {
	githubClient GitHubClient
	graphqlClient *GraphQLClient
	cache        *CacheMiddleware
	timeUtils    *TimeUtils
	config       DORAConfig
}

// DORAConfig configures the enhanced DORA calculator
type DORAConfig struct {
	DefaultTimezone       string        `json:"default_timezone"`
	BusinessHoursStart    int           `json:"business_hours_start"`    // 9 AM
	BusinessHoursEnd      int           `json:"business_hours_end"`      // 5 PM
	ExcludeWeekends       bool          `json:"exclude_weekends"`
	IncidentThresholdHours float64      `json:"incident_threshold_hours"` // Time before considering it an incident
	EnableGraphQL         bool          `json:"enable_graphql"`
	CacheEnabled          bool          `json:"cache_enabled"`
	DefaultCacheTTL       time.Duration `json:"default_cache_ttl"`
	MaxDataPoints         int           `json:"max_data_points"`
}

// NewEnhancedDORACalculator creates a new enhanced DORA calculator
func NewEnhancedDORACalculator(
	githubClient GitHubClient,
	graphqlClient *GraphQLClient,
	cache *CacheMiddleware,
	config DORAConfig,
) *EnhancedDORACalculator {
	// Set defaults
	if config.DefaultTimezone == "" {
		config.DefaultTimezone = "UTC"
	}
	if config.BusinessHoursStart == 0 {
		config.BusinessHoursStart = 9
	}
	if config.BusinessHoursEnd == 0 {
		config.BusinessHoursEnd = 17
	}
	if config.IncidentThresholdHours == 0 {
		config.IncidentThresholdHours = 4.0
	}
	if config.DefaultCacheTTL == 0 {
		config.DefaultCacheTTL = 15 * time.Minute
	}
	if config.MaxDataPoints == 0 {
		config.MaxDataPoints = 1000
	}

	return &EnhancedDORACalculator{
		githubClient:  githubClient,
		graphqlClient: graphqlClient,
		cache:         cache,
		timeUtils:     NewTimeUtils(config.DefaultTimezone),
		config:        config,
	}
}

// Calculate computes enhanced DORA metrics
func (edc *EnhancedDORACalculator) Calculate(ctx context.Context, request MetricsRequest) (*EnhancedDORAMetrics, error) {
	// Use cache if enabled
	if edc.config.CacheEnabled && edc.cache != nil {
		result, cacheInfo, err := edc.cache.CacheOrCompute(ctx, "dora", request, func() (interface{}, error) {
			return edc.calculateInternal(ctx, request)
		})
		if err != nil {
			return nil, err
		}

		metrics := result.(*EnhancedDORAMetrics)
		metrics.CacheInfo = cacheInfo
		return metrics, nil
	}

	return edc.calculateInternal(ctx, request)
}

// calculateInternal performs the actual DORA calculation
func (edc *EnhancedDORACalculator) calculateInternal(ctx context.Context, request MetricsRequest) (*EnhancedDORAMetrics, error) {
	start := time.Now()

	repo := request.Repository
	timeRange := request.TimeRange

	// Validate time range
	if err := edc.timeUtils.ValidateTimezone(timeRange.Timezone); err != nil {
		return nil, fmt.Errorf("invalid timezone: %w", err)
	}

	// Get data using GraphQL if enabled and available, otherwise fallback to REST
	var pullRequests []PullRequest
	var deployments []Deployment
	var workflowRuns []WorkflowRun
	var err error

	if edc.config.EnableGraphQL && edc.graphqlClient != nil {
		pullRequests, deployments, workflowRuns, err = edc.getDataViaGraphQL(ctx, repo, timeRange)
	} else {
		pullRequests, deployments, workflowRuns, err = edc.getDataViaREST(ctx, repo, timeRange)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get data: %w", err)
	}

	// Calculate basic DORA metrics
	leadTimeP95 := edc.calculateEnhancedLeadTime(pullRequests, timeRange)
	deploymentFreq := edc.calculateEnhancedDeploymentFrequency(deployments, timeRange)
	changeFailureRate := edc.calculateEnhancedChangeFailureRate(workflowRuns, deployments)
	mttr := edc.calculateEnhancedMTTR(workflowRuns, deployments, timeRange)

	// Calculate additional metrics
	incidentCount, failedDeployments := edc.analyzeIncidents(workflowRuns, deployments)
	deploymentTrends := edc.calculateDeploymentTrends(deployments, timeRange)
	timeSeries := edc.generateTimeSeries(pullRequests, deployments, workflowRuns, timeRange)
	incidentBreakdown := edc.classifyIncidents(workflowRuns, deployments, timeRange)

	// Calculate confidence and data quality
	confidence := edc.calculateConfidence(pullRequests, deployments, workflowRuns)
	dataQuality := edc.assessDataQuality(pullRequests, deployments, workflowRuns, timeRange)

	// Create enhanced metrics
	enhanced := &EnhancedDORAMetrics{
		DORAMetrics: types.DORAMetrics{
			LeadTimeP95Hours:        leadTimeP95,
			DeploymentFrequencyWeek: deploymentFreq,
			ChangeFailRatePercent:   changeFailureRate,
			MTTRHours:               mttr,
			Period:                  int(timeRange.Duration().Hours() / 24),
			CalculatedAt:            time.Now(),
		},
		TimeRange:              timeRange,
		Granularity:            request.Granularity,
		Timezone:               timeRange.Timezone,
		IncidentCount:          incidentCount,
		FailedDeployments:      failedDeployments,
		TotalDeployments:       len(deployments),
		MeanLeadTimeHours:      edc.calculateMeanLeadTime(pullRequests),
		MedianLeadTimeHours:    edc.calculateMedianLeadTime(pullRequests),
		TimeSeries:             timeSeries,
		IncidentBreakdown:      incidentBreakdown,
		DeploymentTrends:       deploymentTrends,
		Confidence:             confidence,
		DataQuality:            dataQuality,
		CacheInfo: CacheInfo{
			CacheHit:      false,
			ComputeTimeMs: time.Since(start).Milliseconds(),
			DataSources:   edc.getDataSources(),
		},
	}

	return enhanced, nil
}

// Data retrieval methods

func (edc *EnhancedDORACalculator) getDataViaGraphQL(ctx context.Context, repo types.Repository, timeRange TimeRange) ([]PullRequest, []Deployment, []WorkflowRun, error) {
	// Use GraphQL to get comprehensive data
	data, err := edc.graphqlClient.GetRepositoryMetrics(ctx, repo.Owner, repo.Name, timeRange.Start)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("GraphQL query failed: %w", err)
	}

	// Convert GraphQL data to internal types
	pullRequests := edc.convertGraphQLPullRequests(data.Repository.PullRequests.Nodes)
	deployments := edc.convertGraphQLDeployments(data.Repository.Deployments.Nodes)
	workflowRuns := edc.extractWorkflowRunsFromCommits(data.Repository.DefaultBranchRef.Target.History.Nodes)

	return pullRequests, deployments, workflowRuns, nil
}

func (edc *EnhancedDORACalculator) getDataViaREST(ctx context.Context, repo types.Repository, timeRange TimeRange) ([]PullRequest, []Deployment, []WorkflowRun, error) {
	// Use REST API as fallback
	pullRequests, err := edc.githubClient.GetPullRequests(ctx, repo.Owner, repo.Name, timeRange.Start)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get pull requests: %w", err)
	}

	deployments, err := edc.githubClient.GetDeployments(ctx, repo.Owner, repo.Name, timeRange.Start)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get deployments: %w", err)
	}

	workflowRuns, err := edc.githubClient.GetWorkflowRuns(ctx, repo.Owner, repo.Name, timeRange.Start)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get workflow runs: %w", err)
	}

	return pullRequests, deployments, workflowRuns, nil
}

// Enhanced calculation methods

func (edc *EnhancedDORACalculator) calculateEnhancedLeadTime(pullRequests []PullRequest, timeRange TimeRange) float64 {
	var leadTimes []float64

	for _, pr := range pullRequests {
		if pr.MergedAt == nil {
			continue
		}

		// Calculate lead time considering business hours if configured
		var leadTime float64
		if edc.config.ExcludeWeekends {
			workingHours, err := edc.timeUtils.CalculateWorkingHours(
				pr.CreatedAt, *pr.MergedAt, timeRange.Timezone,
				edc.config.BusinessHoursStart, edc.config.BusinessHoursEnd, true,
			)
			if err == nil {
				leadTime = workingHours
			} else {
				leadTime = pr.MergedAt.Sub(pr.CreatedAt).Hours()
			}
		} else {
			leadTime = pr.MergedAt.Sub(pr.CreatedAt).Hours()
		}

		leadTimes = append(leadTimes, leadTime)
	}

	if len(leadTimes) == 0 {
		return 0
	}

	// Calculate P95
	sort.Float64s(leadTimes)
	index := int(math.Ceil(0.95 * float64(len(leadTimes))))
	if index >= len(leadTimes) {
		index = len(leadTimes) - 1
	}

	return leadTimes[index]
}

func (edc *EnhancedDORACalculator) calculateEnhancedDeploymentFrequency(deployments []Deployment, timeRange TimeRange) float64 {
	if len(deployments) == 0 {
		return 0
	}

	// Count successful deployments
	successfulDeployments := 0
	for _, deployment := range deployments {
		if deployment.State == "success" || deployment.State == "active" {
			successfulDeployments++
		}
	}

	// Calculate per week frequency
	weeks := timeRange.Duration().Hours() / (7 * 24)
	if weeks == 0 {
		weeks = 1
	}

	return float64(successfulDeployments) / weeks
}

func (edc *EnhancedDORACalculator) calculateEnhancedChangeFailureRate(workflowRuns []WorkflowRun, deployments []Deployment) float64 {
	totalChanges := len(workflowRuns)
	if totalChanges == 0 {
		totalChanges = len(deployments)
	}
	if totalChanges == 0 {
		return 0
	}

	failures := 0
	for _, run := range workflowRuns {
		if run.Conclusion == "failure" {
			failures++
		}
	}

	// Also count failed deployments
	for _, deployment := range deployments {
		if deployment.State == "failure" || deployment.State == "error" {
			failures++
		}
	}

	return (float64(failures) / float64(totalChanges)) * 100
}

func (edc *EnhancedDORACalculator) calculateEnhancedMTTR(workflowRuns []WorkflowRun, deployments []Deployment, timeRange TimeRange) float64 {
	var recoveryTimes []float64

	// Analyze workflow runs for recovery patterns
	for i, run := range workflowRuns {
		if run.Conclusion != "failure" {
			continue
		}

		// Look for the next successful run to calculate recovery time
		for j := i + 1; j < len(workflowRuns); j++ {
			nextRun := workflowRuns[j]
			if nextRun.Conclusion == "success" {
				recoveryTime := nextRun.UpdatedAt.Sub(run.UpdatedAt).Hours()
				if recoveryTime > 0 && recoveryTime < edc.config.IncidentThresholdHours*24 { // Max 24x threshold
					recoveryTimes = append(recoveryTimes, recoveryTime)
				}
				break
			}
		}
	}

	// Analyze deployments for recovery patterns
	for i, deployment := range deployments {
		if deployment.State != "failure" && deployment.State != "error" {
			continue
		}

		// Look for the next successful deployment
		for j := i + 1; j < len(deployments); j++ {
			nextDeployment := deployments[j]
			if nextDeployment.State == "success" {
				recoveryTime := nextDeployment.UpdatedAt.Sub(deployment.UpdatedAt).Hours()
				if recoveryTime > 0 && recoveryTime < edc.config.IncidentThresholdHours*24 {
					recoveryTimes = append(recoveryTimes, recoveryTime)
				}
				break
			}
		}
	}

	if len(recoveryTimes) == 0 {
		return 0
	}

	// Calculate mean recovery time
	total := 0.0
	for _, time := range recoveryTimes {
		total += time
	}

	return total / float64(len(recoveryTimes))
}

// Analysis methods

func (edc *EnhancedDORACalculator) analyzeIncidents(workflowRuns []WorkflowRun, deployments []Deployment) (int, int) {
	incidents := 0
	failedDeployments := 0

	// Count workflow failures that exceed incident threshold
	for _, run := range workflowRuns {
		if run.Conclusion == "failure" {
			incidents++
		}
	}

	// Count failed deployments
	for _, deployment := range deployments {
		if deployment.State == "failure" || deployment.State == "error" {
			failedDeployments++
			if deployment.UpdatedAt.Sub(deployment.CreatedAt).Hours() > edc.config.IncidentThresholdHours {
				incidents++
			}
		}
	}

	return incidents, failedDeployments
}

func (edc *EnhancedDORACalculator) calculateDeploymentTrends(deployments []Deployment, timeRange TimeRange) []DeploymentTrend {
	if len(deployments) == 0 {
		return nil
	}

	// Split time range into periods based on granularity
	periods, err := edc.timeUtils.GetPeriodBoundaries(timeRange.End, "week", timeRange.Timezone, 4)
	if err != nil {
		return nil
	}

	var trends []DeploymentTrend
	for _, period := range periods {
		deploymentCount := 0
		successfulDeployments := 0
		totalLeadTime := 0.0

		for _, deployment := range deployments {
			if period.Contains(deployment.CreatedAt) {
				deploymentCount++
				if deployment.State == "success" {
					successfulDeployments++
				}
			}
		}

		successRate := 0.0
		if deploymentCount > 0 {
			successRate = float64(successfulDeployments) / float64(deploymentCount) * 100
		}

		// Determine trend direction (simplified)
		direction := "stable"
		if successRate > 80 {
			direction = "improving"
		} else if successRate < 50 {
			direction = "declining"
		}

		trends = append(trends, DeploymentTrend{
			Period:          "week",
			DeploymentCount: deploymentCount,
			SuccessRate:     successRate,
			AverageLeadTime: totalLeadTime / float64(deploymentCount),
			TrendDirection:  direction,
		})
	}

	return trends
}

func (edc *EnhancedDORACalculator) generateTimeSeries(pullRequests []PullRequest, deployments []Deployment, workflowRuns []WorkflowRun, timeRange TimeRange) []DORATimeSeriesPoint {
	// Generate time series based on granularity
	periods, err := edc.timeUtils.GetPeriodBoundaries(timeRange.End, "day", timeRange.Timezone, int(timeRange.Duration().Hours()/24))
	if err != nil {
		return nil
	}

	var timeSeries []DORATimeSeriesPoint
	for _, period := range periods {
		point := DORATimeSeriesPoint{
			Timestamp: period.Start,
		}

		// Calculate metrics for this time period
		for _, pr := range pullRequests {
			if pr.MergedAt != nil && period.Contains(*pr.MergedAt) {
				leadTime := pr.MergedAt.Sub(pr.CreatedAt).Hours()
				if point.LeadTimeHours == 0 || leadTime > point.LeadTimeHours {
					point.LeadTimeHours = leadTime
				}
			}
		}

		for _, deployment := range deployments {
			if period.Contains(deployment.CreatedAt) {
				point.DeploymentCount++
				if deployment.State == "failure" || deployment.State == "error" {
					point.FailureCount++
				}
			}
		}

		for _, run := range workflowRuns {
			if period.Contains(run.UpdatedAt) && run.Conclusion == "failure" {
				point.RecoveryTimeHours = run.UpdatedAt.Sub(run.CreatedAt).Hours()
			}
		}

		if point.DeploymentCount > 0 {
			point.ChangeFailureRate = float64(point.FailureCount) / float64(point.DeploymentCount) * 100
		}

		timeSeries = append(timeSeries, point)
	}

	return timeSeries
}

func (edc *EnhancedDORACalculator) classifyIncidents(workflowRuns []WorkflowRun, deployments []Deployment, timeRange TimeRange) []IncidentClassification {
	incidentMap := make(map[string]*IncidentClassification)

	// Classify workflow failures
	for _, run := range workflowRuns {
		if run.Conclusion == "failure" {
			key := "workflow_failure"
			if incident, exists := incidentMap[key]; exists {
				incident.Count++
			} else {
				incidentMap[key] = &IncidentClassification{
					Type:     "workflow_failure",
					Severity: edc.determineSeverity(run.UpdatedAt.Sub(run.CreatedAt)),
					Count:    1,
				}
			}
		}
	}

	// Classify deployment failures
	for _, deployment := range deployments {
		if deployment.State == "failure" || deployment.State == "error" {
			key := "deployment_failure"
			if incident, exists := incidentMap[key]; exists {
				incident.Count++
			} else {
				incidentMap[key] = &IncidentClassification{
					Type:     "deployment_failure",
					Severity: edc.determineSeverity(deployment.UpdatedAt.Sub(deployment.CreatedAt)),
					Count:    1,
				}
			}
		}
	}

	// Convert map to slice
	var incidents []IncidentClassification
	for _, incident := range incidentMap {
		incidents = append(incidents, *incident)
	}

	return incidents
}

// Helper methods

func (edc *EnhancedDORACalculator) calculateMeanLeadTime(pullRequests []PullRequest) float64 {
	if len(pullRequests) == 0 {
		return 0
	}

	total := 0.0
	count := 0
	for _, pr := range pullRequests {
		if pr.MergedAt != nil {
			total += pr.MergedAt.Sub(pr.CreatedAt).Hours()
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return total / float64(count)
}

func (edc *EnhancedDORACalculator) calculateMedianLeadTime(pullRequests []PullRequest) float64 {
	var leadTimes []float64
	for _, pr := range pullRequests {
		if pr.MergedAt != nil {
			leadTimes = append(leadTimes, pr.MergedAt.Sub(pr.CreatedAt).Hours())
		}
	}

	if len(leadTimes) == 0 {
		return 0
	}

	sort.Float64s(leadTimes)
	mid := len(leadTimes) / 2
	if len(leadTimes)%2 == 0 {
		return (leadTimes[mid-1] + leadTimes[mid]) / 2
	}
	return leadTimes[mid]
}

func (edc *EnhancedDORACalculator) calculateConfidence(pullRequests []PullRequest, deployments []Deployment, workflowRuns []WorkflowRun) float64 {
	// Calculate confidence based on data completeness and quality
	dataPoints := len(pullRequests) + len(deployments) + len(workflowRuns)
	if dataPoints == 0 {
		return 0.0
	}

	// Higher confidence with more data points, up to a maximum
	maxPoints := 100
	confidence := float64(dataPoints) / float64(maxPoints)
	if confidence > 1.0 {
		confidence = 1.0
	}

	// Reduce confidence if data seems incomplete
	if len(pullRequests) == 0 {
		confidence *= 0.7
	}
	if len(deployments) == 0 {
		confidence *= 0.8
	}

	return confidence
}

func (edc *EnhancedDORACalculator) assessDataQuality(pullRequests []PullRequest, deployments []Deployment, workflowRuns []WorkflowRun, timeRange TimeRange) DataQuality {
	totalDataPoints := len(pullRequests) + len(deployments) + len(workflowRuns)
	missingData := 0
	var warnings []string

	// Check for missing timestamps
	for _, pr := range pullRequests {
		if pr.CreatedAt.IsZero() || (pr.State == "merged" && pr.MergedAt == nil) {
			missingData++
		}
	}

	for _, deployment := range deployments {
		if deployment.CreatedAt.IsZero() {
			missingData++
		}
	}

	// Check data coverage
	expectedDays := int(timeRange.Duration().Hours() / 24)
	if expectedDays > 30 && len(pullRequests) == 0 {
		warnings = append(warnings, "No pull request data found for extended period")
	}

	completeness := 1.0 - (float64(missingData) / float64(totalDataPoints))
	if completeness < 0 {
		completeness = 0
	}

	return DataQuality{
		Completeness:    completeness,
		Accuracy:        0.9, // Assume high accuracy from GitHub API
		Timeliness:      1.0, // Real-time data
		Consistency:     0.95,
		DataPoints:      totalDataPoints,
		MissingData:     missingData,
		QualityWarnings: warnings,
	}
}

func (edc *EnhancedDORACalculator) determineSeverity(duration time.Duration) string {
	hours := duration.Hours()
	if hours > 24 {
		return "critical"
	} else if hours > 8 {
		return "high"
	} else if hours > 2 {
		return "medium"
	}
	return "low"
}

func (edc *EnhancedDORACalculator) getDataSources() []string {
	sources := []string{"github_rest_api"}
	if edc.config.EnableGraphQL {
		sources = append(sources, "github_graphql_api")
	}
	return sources
}

// Conversion methods for GraphQL data

func (edc *EnhancedDORACalculator) convertGraphQLPullRequests(graphqlPRs []GraphQLPullRequest) []PullRequest {
	var pullRequests []PullRequest
	for _, gpr := range graphqlPRs {
		pr := PullRequest{
			Number:        gpr.Number,
			Title:         gpr.Title,
			State:         gpr.State,
			CreatedAt:     gpr.CreatedAt,
			UpdatedAt:     gpr.UpdatedAt,
			MergedAt:      gpr.MergedAt,
			ClosedAt:      gpr.ClosedAt,
			Commits:       gpr.Commits.TotalCount,
			Additions:     gpr.Additions,
			Deletions:     gpr.Deletions,
			ChangedFiles:  gpr.ChangedFiles,
		}

		// Get first review time if available
		if len(gpr.Reviews.Nodes) > 0 {
			firstReview := gpr.Reviews.Nodes[0].SubmittedAt
			pr.FirstReviewAt = &firstReview
		}

		pullRequests = append(pullRequests, pr)
	}
	return pullRequests
}

func (edc *EnhancedDORACalculator) convertGraphQLDeployments(graphqlDeployments []GraphQLDeployment) []Deployment {
	var deployments []Deployment
	for _, gd := range graphqlDeployments {
		deployment := Deployment{
			ID:          0, // GraphQL doesn't provide numeric ID
			Environment: gd.Environment,
			State:       gd.State,
			CreatedAt:   gd.CreatedAt,
			UpdatedAt:   gd.UpdatedAt,
			SHA:         gd.Ref.Target.Oid,
		}
		deployments = append(deployments, deployment)
	}
	return deployments
}

func (edc *EnhancedDORACalculator) extractWorkflowRunsFromCommits(commits []GraphQLCommit) []WorkflowRun {
	// This is a simplified extraction - in practice, you'd need additional GraphQL queries
	// to get actual workflow run data
	var workflowRuns []WorkflowRun
	for _, commit := range commits {
		// Create a mock workflow run based on commit data
		// In reality, you'd query the actual workflow runs
		run := WorkflowRun{
			ID:         0, // Would come from actual workflow API
			Name:       "ci",
			Status:     "completed",
			Conclusion: "success", // Assume success unless we have evidence otherwise
			CreatedAt:  commit.CommittedDate,
			UpdatedAt:  commit.CommittedDate,
			SHA:        commit.Oid,
		}
		workflowRuns = append(workflowRuns, run)
	}
	return workflowRuns
}