// Package metrics implements calculation engines for DORA, CHI, and AI impact metrics.
package metrics

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/types"
)

// DORACalculator calculates DevOps Research and Assessment metrics
type DORACalculator struct {
	githubClient GitHubClient
	jiraClient   JiraClient
}

// GitHubClient interface for repository data access
type GitHubClient interface {
	GetPullRequests(ctx context.Context, owner, repo string, since time.Time) ([]PullRequest, error)
	GetDeployments(ctx context.Context, owner, repo string, since time.Time) ([]Deployment, error)
	GetWorkflowRuns(ctx context.Context, owner, repo string, since time.Time) ([]WorkflowRun, error)
}

// JiraClient interface for issue tracking data
type JiraClient interface {
	GetIssues(ctx context.Context, project string, since time.Time) ([]Issue, error)
}

// PullRequest represents a GitHub pull request
type PullRequest struct {
	Number        int        `json:"number"`
	Title         string     `json:"title"`
	State         string     `json:"state"` // open, closed, merged
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	MergedAt      *time.Time `json:"merged_at"`
	ClosedAt      *time.Time `json:"closed_at"`
	Commits       int        `json:"commits"`
	Additions     int        `json:"additions"`
	Deletions     int        `json:"deletions"`
	ChangedFiles  int        `json:"changed_files"`
	FirstReviewAt *time.Time `json:"first_review_at"`
}

// Deployment represents a deployment event
type Deployment struct {
	ID          int       `json:"id"`
	Environment string    `json:"environment"`
	State       string    `json:"state"` // success, failure, error, pending
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	SHA         string    `json:"sha"`
}

// WorkflowRun represents a CI/CD pipeline run
type WorkflowRun struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`     // completed, in_progress, queued
	Conclusion string    `json:"conclusion"` // success, failure, cancelled, skipped
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	SHA        string    `json:"sha"`
}

// Issue represents a Jira issue
type Issue struct {
	Key        string     `json:"key"`
	Type       string     `json:"type"`
	Status     string     `json:"status"`
	Priority   string     `json:"priority"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ResolvedAt *time.Time `json:"resolved_at"`
}

// NewDORACalculator creates a new DORA metrics calculator
func NewDORACalculator(github GitHubClient, jira JiraClient) *DORACalculator {
	return &DORACalculator{
		githubClient: github,
		jiraClient:   jira,
	}
}

// Calculate computes DORA metrics for a repository
func (d *DORACalculator) Calculate(ctx context.Context, repo types.Repository, periodDays int) (*types.DORAMetrics, error) {
	since := time.Now().AddDate(0, 0, -periodDays)

	// Get data from GitHub
	prs, err := d.githubClient.GetPullRequests(ctx, repo.Owner, repo.Name, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull requests: %w", err)
	}

	deployments, err := d.githubClient.GetDeployments(ctx, repo.Owner, repo.Name, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployments: %w", err)
	}

	workflows, err := d.githubClient.GetWorkflowRuns(ctx, repo.Owner, repo.Name, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow runs: %w", err)
	}

	// Calculate each DORA metric
	leadTime := d.calculateLeadTimeP95(prs)
	deployFreq := d.calculateDeploymentFrequency(deployments, periodDays)
	changeFailRate := d.calculateChangeFailureRate(workflows)
	mttr := d.calculateMTTR(workflows, deployments)

	return &types.DORAMetrics{
		LeadTimeP95Hours:        leadTime,
		DeploymentFrequencyWeek: deployFreq,
		ChangeFailRatePercent:   changeFailRate,
		MTTRHours:               mttr,
		Period:                  periodDays,
		CalculatedAt:            time.Now(),
	}, nil
}

// calculateLeadTimeP95 calculates the 95th percentile lead time in hours
func (d *DORACalculator) calculateLeadTimeP95(prs []PullRequest) float64 {
	var leadTimes []float64

	for _, pr := range prs {
		if pr.MergedAt != nil {
			leadTime := pr.MergedAt.Sub(pr.CreatedAt).Hours()
			leadTimes = append(leadTimes, leadTime)
		}
	}

	if len(leadTimes) == 0 {
		return 0
	}

	sort.Float64s(leadTimes)
	p95Index := int(math.Ceil(0.95*float64(len(leadTimes)))) - 1
	if p95Index >= len(leadTimes) {
		p95Index = len(leadTimes) - 1
	}

	return leadTimes[p95Index]
}

// calculateDeploymentFrequency calculates deployments per week
func (d *DORACalculator) calculateDeploymentFrequency(deployments []Deployment, periodDays int) float64 {
	successfulDeploys := 0
	for _, deploy := range deployments {
		if deploy.State == "success" {
			successfulDeploys++
		}
	}

	weeks := float64(periodDays) / 7.0
	if weeks == 0 {
		return 0
	}

	return float64(successfulDeploys) / weeks
}

// calculateChangeFailureRate calculates the percentage of failed deployments
func (d *DORACalculator) calculateChangeFailureRate(workflows []WorkflowRun) float64 {
	totalRuns := 0
	failedRuns := 0

	for _, run := range workflows {
		if run.Status == "completed" {
			totalRuns++
			if run.Conclusion == "failure" {
				failedRuns++
			}
		}
	}

	if totalRuns == 0 {
		return 0
	}

	return (float64(failedRuns) / float64(totalRuns)) * 100.0
}

// calculateMTTR calculates Mean Time To Recovery in hours
func (d *DORACalculator) calculateMTTR(workflows []WorkflowRun, deployments []Deployment) float64 {
	var recoveryTimes []float64

	// Find failure -> success patterns in workflows
	var lastFailureTime *time.Time

	// Sort by creation time
	sort.Slice(workflows, func(i, j int) bool {
		return workflows[i].CreatedAt.Before(workflows[j].CreatedAt)
	})

	for _, run := range workflows {
		if run.Status == "completed" {
			if run.Conclusion == "failure" {
				lastFailureTime = &run.UpdatedAt
			} else if run.Conclusion == "success" && lastFailureTime != nil {
				recoveryTime := run.UpdatedAt.Sub(*lastFailureTime).Hours()
				recoveryTimes = append(recoveryTimes, recoveryTime)
				lastFailureTime = nil // Reset
			}
		}
	}

	if len(recoveryTimes) == 0 {
		return 0
	}

	// Calculate mean
	sum := 0.0
	for _, rt := range recoveryTimes {
		sum += rt
	}

	return sum / float64(len(recoveryTimes))
}
