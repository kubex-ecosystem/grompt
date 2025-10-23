package grompt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/metrics"
	"github.com/kubex-ecosystem/grompt/internal/scorecard"
	"github.com/kubex-ecosystem/grompt/internal/types"
)

// MockGitHubClient for testing
type MockGitHubClient struct{}

func (m *MockGitHubClient) GetPullRequests(ctx context.Context, owner, repo string, since time.Time) ([]metrics.PullRequest, error) {
	// Return realistic test data
	now := time.Now()
	return []metrics.PullRequest{
		{
			Number:        123,
			Title:         "Add user authentication",
			State:         "merged",
			CreatedAt:     now.AddDate(0, 0, -7),
			MergedAt:      &[]time.Time{now.AddDate(0, 0, -5)}[0],
			Commits:       5,
			Additions:     150,
			Deletions:     30,
			ChangedFiles:  8,
			FirstReviewAt: &[]time.Time{now.AddDate(0, 0, -6)}[0],
		},
		{
			Number:        124,
			Title:         "Fix payment validation bug",
			State:         "merged",
			CreatedAt:     now.AddDate(0, 0, -3),
			MergedAt:      &[]time.Time{now.AddDate(0, 0, -1)}[0],
			Commits:       2,
			Additions:     25,
			Deletions:     10,
			ChangedFiles:  3,
			FirstReviewAt: &[]time.Time{now.AddDate(0, 0, -2)}[0],
		},
	}, nil
}

func (m *MockGitHubClient) GetDeployments(ctx context.Context, owner, repo string, since time.Time) ([]metrics.Deployment, error) {
	now := time.Now()
	return []metrics.Deployment{
		{
			ID:          1,
			Environment: "production",
			State:       "success",
			CreatedAt:   now.AddDate(0, 0, -5),
			UpdatedAt:   now.AddDate(0, 0, -5),
			SHA:         "abc123",
		},
		{
			ID:          2,
			Environment: "production",
			State:       "success",
			CreatedAt:   now.AddDate(0, 0, -1),
			UpdatedAt:   now.AddDate(0, 0, -1),
			SHA:         "def456",
		},
	}, nil
}

func (m *MockGitHubClient) GetWorkflowRuns(ctx context.Context, owner, repo string, since time.Time) ([]metrics.WorkflowRun, error) {
	now := time.Now()
	return []metrics.WorkflowRun{
		{
			ID:         1,
			Name:       "CI",
			Status:     "completed",
			Conclusion: "success",
			CreatedAt:  now.AddDate(0, 0, -5),
			UpdatedAt:  now.AddDate(0, 0, -5),
			SHA:        "abc123",
		},
		{
			ID:         2,
			Name:       "CI",
			Status:     "completed",
			Conclusion: "failure",
			CreatedAt:  now.AddDate(0, 0, -3),
			UpdatedAt:  now.AddDate(0, 0, -3),
			SHA:        "xyz789",
		},
		{
			ID:         3,
			Name:       "CI",
			Status:     "completed",
			Conclusion: "success",
			CreatedAt:  now.AddDate(0, 0, -1),
			UpdatedAt:  now.AddDate(0, 0, -1),
			SHA:        "def456",
		},
	}, nil
}

// MockJiraClient for testing
type MockJiraClient struct{}

func (m *MockJiraClient) GetIssues(ctx context.Context, project string, since time.Time) ([]metrics.Issue, error) {
	now := time.Now()
	return []metrics.Issue{
		{
			Key:        "PROJ-123",
			Type:       "Bug",
			Status:     "Done",
			Priority:   "High",
			CreatedAt:  now.AddDate(0, 0, -10),
			UpdatedAt:  now.AddDate(0, 0, -5),
			ResolvedAt: &[]time.Time{now.AddDate(0, 0, -5)}[0],
		},
		{
			Key:       "PROJ-124",
			Type:      "Feature",
			Status:    "In Progress",
			Priority:  "Medium",
			CreatedAt: now.AddDate(0, 0, -7),
			UpdatedAt: now.AddDate(0, 0, -1),
		},
	}, nil
}

// MockWakaTimeClient for AI metrics testing
type MockWakaTimeClient struct{}

func (m *MockWakaTimeClient) GetCodingTime(ctx context.Context, user, repo string, since time.Time) (*metrics.CodingTime, error) {
	return &metrics.CodingTime{
		TotalHours:  40.0,
		CodingHours: 32.0,
		Period:      30,
	}, nil
}

// MockGitClient for AI metrics testing
type MockGitClient struct{}

func (m *MockGitClient) GetCommits(ctx context.Context, owner, repo string, since time.Time) ([]metrics.Commit, error) {
	now := time.Now()
	return []metrics.Commit{
		{
			SHA:        "abc123",
			Message:    "Add user authentication",
			Author:     "test-user",
			Date:       now.AddDate(0, 0, -5),
			Additions:  150,
			Deletions:  30,
			AIAssisted: true,
			AIProvider: "copilot",
		},
		{
			SHA:        "def456",
			Message:    "Fix payment validation",
			Author:     "test-user",
			Date:       now.AddDate(0, 0, -1),
			Additions:  25,
			Deletions:  10,
			AIAssisted: false,
		},
	}, nil
}

// MockIDEClient for AI metrics testing
type MockIDEClient struct{}

func (m *MockIDEClient) GetAIAssistData(ctx context.Context, user, repo string, since time.Time) (*metrics.AIAssistData, error) {
	return &metrics.AIAssistData{
		TotalSuggestions:    100,
		AcceptedSuggestions: 75,
		AcceptanceRate:      0.75,
		TimeWithAI:          20.0,
		LinesGenerated:      300,
		Provider:            "copilot",
	}, nil
}

func main() {
	ctx := context.Background()

	// Create test repository
	repo := types.Repository{
		Owner:         "kubex-ecosystem",
		Name:          "grompt",
		FullName:      "kubex-ecosystem/grompt",
		DefaultBranch: "main",
		Language:      "Go",
		CreatedAt:     time.Now().AddDate(0, 0, -365),
		UpdatedAt:     time.Now(),
	}

	fmt.Println("🔍 Testing DORA Metrics Calculation...")

	// Test DORA calculator
	githubClient := &MockGitHubClient{}
	jiraClient := &MockJiraClient{}
	doraCalc := metrics.NewDORACalculator(githubClient, jiraClient)

	doraMetrics, err := doraCalc.Calculate(ctx, repo, 30)
	if err != nil {
		log.Fatalf("DORA calculation failed: %v", err)
	}

	doraJSON, _ := json.MarshalIndent(doraMetrics, "", "  ")
	fmt.Printf("✅ DORA Metrics:\n%s\n\n", doraJSON)

	fmt.Println("🔍 Testing CHI Metrics Calculation...")

	// Test CHI calculator with current repository
	chiCalc := metrics.NewCHICalculator("/srv/apps/LIFE/KUBEX/grompt")

	chiMetrics, err := chiCalc.Calculate(ctx, repo)
	if err != nil {
		log.Fatalf("CHI calculation failed: %v", err)
	}

	chiJSON, _ := json.MarshalIndent(chiMetrics, "", "  ")
	fmt.Printf("✅ CHI Metrics:\n%s\n\n", chiJSON)

	fmt.Println("🔍 Testing Scorecard Engine...")

	// Test AI metrics calculator with mock clients
	wakatimeClient := &MockWakaTimeClient{}
	gitClient := &MockGitClient{}
	ideClient := &MockIDEClient{}
	aiCalc := metrics.NewAIMetricsCalculator(wakatimeClient, gitClient, ideClient)

	// Create scorecard engine
	engine := scorecard.NewEngine(doraCalc, chiCalc, aiCalc)

	scorecard, err := engine.GenerateScorecard(ctx, repo, "test-user", 30)
	if err != nil {
		log.Fatalf("Scorecard generation failed: %v", err)
	}

	scorecardJSON, _ := json.MarshalIndent(scorecard, "", "  ")
	fmt.Printf("✅ Complete Scorecard:\n%s\n\n", scorecardJSON)

	fmt.Println("🎉 All metrics calculations completed successfully!")
	fmt.Println("✅ Day 1 metrics validation: PASSED")
}
