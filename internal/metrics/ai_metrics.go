// Package metrics - AI Impact Metrics (HIR/AAC/TPH) calculation engine
package metrics

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/types"
)

// AIMetricsCalculator calculates Human Input Ratio, AI Assist Coverage, and Throughput per Human-hour
type AIMetricsCalculator struct {
	wakatimeClient WakaTimeClient
	gitClient      GitClient
	ideClient      IDEClient
}

// WakaTimeClient interface for time tracking data
type WakaTimeClient interface {
	GetCodingTime(ctx context.Context, user, repo string, since time.Time) (*CodingTime, error)
}

// GitClient interface for Git commit analysis
type GitClient interface {
	GetCommits(ctx context.Context, owner, repo string, since time.Time) ([]Commit, error)
}

// IDEClient interface for IDE telemetry data
type IDEClient interface {
	GetAIAssistData(ctx context.Context, user, repo string, since time.Time) (*AIAssistData, error)
}

// CodingTime represents time tracking data from WakaTime
type CodingTime struct {
	TotalHours  float64        `json:"total_hours"`
	CodingHours float64        `json:"coding_hours"`
	Period      int            `json:"period_days"`
	Languages   []LanguageTime `json:"languages"`
	Projects    []ProjectTime  `json:"projects"`
}

type LanguageTime struct {
	Name  string  `json:"name"`
	Hours float64 `json:"hours"`
}

type ProjectTime struct {
	Name  string  `json:"name"`
	Hours float64 `json:"hours"`
}

// Commit represents a Git commit with AI assistance indicators
type Commit struct {
	SHA       string    `json:"sha"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	Date      time.Time `json:"date"`
	Files     []string  `json:"files"`
	Additions int       `json:"additions"`
	Deletions int       `json:"deletions"`
	// AI assistance indicators
	CoAuthoredBy []string `json:"co_authored_by"`
	AIAssisted   bool     `json:"ai_assisted"`
	AIProvider   string   `json:"ai_provider"`
}

// AIAssistData represents AI assistance data from IDE
type AIAssistData struct {
	TotalSuggestions    int     `json:"total_suggestions"`
	AcceptedSuggestions int     `json:"accepted_suggestions"`
	AcceptanceRate      float64 `json:"acceptance_rate"`
	TimeWithAI          float64 `json:"time_with_ai_hours"`
	LinesGenerated      int     `json:"lines_generated"`
	Provider            string  `json:"provider"` // copilot, codewhisperer, codeium, etc.
}

// NewAIMetricsCalculator creates a new AI metrics calculator
func NewAIMetricsCalculator(wakatime WakaTimeClient, git GitClient, ide IDEClient) *AIMetricsCalculator {
	return &AIMetricsCalculator{
		wakatimeClient: wakatime,
		gitClient:      git,
		ideClient:      ide,
	}
}

// Calculate computes AI impact metrics for a repository
func (a *AIMetricsCalculator) Calculate(ctx context.Context, repo types.Repository, user string, periodDays int) (*types.AIMetrics, error) {
	since := time.Now().AddDate(0, 0, -periodDays)

	// Get time tracking data
	codingTime, err := a.wakatimeClient.GetCodingTime(ctx, user, repo.Name, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get coding time: %w", err)
	}

	// Get commit data
	commits, err := a.gitClient.GetCommits(ctx, repo.Owner, repo.Name, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	// Get AI assistance data
	aiData, err := a.ideClient.GetAIAssistData(ctx, user, repo.Name, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI assist data: %w", err)
	}

	// Calculate metrics
	hir := a.calculateHIR(codingTime, aiData)
	aac := a.calculateAAC(commits)
	tph := a.calculateTPH(commits, codingTime)

	return &types.AIMetrics{
		HIR:          hir,
		AAC:          aac,
		TPH:          tph,
		HumanHours:   codingTime.CodingHours - aiData.TimeWithAI,
		AIHours:      aiData.TimeWithAI,
		Period:       periodDays,
		CalculatedAt: time.Now(),
	}, nil
}

// calculateHIR calculates Human Input Ratio (0.0-1.0)
// HIR = human_edit_time / (human_edit_time + ai_assist_time)
func (a *AIMetricsCalculator) calculateHIR(codingTime *CodingTime, aiData *AIAssistData) float64 {
	humanHours := codingTime.CodingHours - aiData.TimeWithAI
	if humanHours < 0 {
		humanHours = 0
	}

	totalHours := humanHours + aiData.TimeWithAI
	if totalHours == 0 {
		return 1.0 // No AI assistance means 100% human input
	}

	hir := humanHours / totalHours
	if hir < 0 {
		hir = 0
	}
	if hir > 1 {
		hir = 1
	}

	return hir
}

// calculateAAC calculates AI Assist Coverage (0.0-1.0)
// AAC = % of commits/PRs with AI assistance
func (a *AIMetricsCalculator) calculateAAC(commits []Commit) float64 {
	if len(commits) == 0 {
		return 0
	}

	aiAssistedCommits := 0

	for _, commit := range commits {
		if a.isAIAssisted(commit) {
			aiAssistedCommits++
		}
	}

	return float64(aiAssistedCommits) / float64(len(commits))
}

// calculateTPH calculates Throughput per Human-hour
// TPH = approved_artifacts / human_hours
func (a *AIMetricsCalculator) calculateTPH(commits []Commit, codingTime *CodingTime) float64 {
	humanHours := codingTime.CodingHours
	if humanHours == 0 {
		return 0
	}

	// Calculate useful artifacts (commits, LOC, etc.)
	totalLOC := 0
	for _, commit := range commits {
		totalLOC += commit.Additions
	}

	// Use commits as artifacts for now
	// In production, could use PRs, issues resolved, features delivered, etc.
	artifacts := float64(len(commits))

	return artifacts / humanHours
}

// isAIAssisted determines if a commit was AI-assisted
func (a *AIMetricsCalculator) isAIAssisted(commit Commit) bool {
	// Check explicit AI assisted flag
	if commit.AIAssisted {
		return true
	}

	// Check co-authored-by patterns
	for _, coAuthor := range commit.CoAuthoredBy {
		if a.isAICoAuthor(coAuthor) {
			return true
		}
	}

	// Check commit message patterns
	aiPatterns := []string{
		"copilot",
		"codewhisperer",
		"codeium",
		"ai-assisted",
		"ai-generated",
		"with ai",
		"assisted by",
	}

	message := strings.ToLower(commit.Message)
	for _, pattern := range aiPatterns {
		if strings.Contains(message, pattern) {
			return true
		}
	}

	return false
}

// isAICoAuthor checks if a co-author is an AI system
func (a *AIMetricsCalculator) isAICoAuthor(coAuthor string) bool {
	aiCoAuthors := []string{
		"github-copilot",
		"copilot",
		"codewhisperer",
		"codeium",
		"tabnine",
		"ai-assistant",
	}

	coAuthorLower := strings.ToLower(coAuthor)
	for _, ai := range aiCoAuthors {
		if strings.Contains(coAuthorLower, ai) {
			return true
		}
	}

	return false
}

// AnalyzeAIImpact provides insights on AI usage patterns
func (a *AIMetricsCalculator) AnalyzeAIImpact(hir, aac, tph float64) AIImpactAnalysis {
	analysis := AIImpactAnalysis{}

	// Analyze HIR
	if hir > 0.8 {
		analysis.HIRInsight = "High human input - good manual craftsmanship"
		analysis.HIRRecommendation = "Consider AI assistance for repetitive tasks"
	} else if hir > 0.5 {
		analysis.HIRInsight = "Balanced human-AI collaboration"
		analysis.HIRRecommendation = "Good leverage of AI assistance"
	} else {
		analysis.HIRInsight = "High AI assistance usage"
		analysis.HIRRecommendation = "Monitor code quality and avoid over-prompting"
	}

	// Analyze AAC
	if aac > 0.7 {
		analysis.AACInsight = "High AI adoption across commits"
		analysis.AACRecommendation = "Ensure AI assists with quality, not just quantity"
	} else if aac > 0.3 {
		analysis.AACInsight = "Moderate AI usage"
		analysis.AACRecommendation = "Consider expanding AI assistance to repetitive tasks"
	} else {
		analysis.AACInsight = "Low AI adoption"
		analysis.AACRecommendation = "Explore AI tools to increase productivity"
	}

	// Analyze TPH
	if tph > 2.0 {
		analysis.TPHInsight = "High throughput per hour"
		analysis.TPHRecommendation = "Maintain current productivity patterns"
	} else if tph > 1.0 {
		analysis.TPHInsight = "Moderate throughput"
		analysis.TPHRecommendation = "Look for automation opportunities"
	} else {
		analysis.TPHInsight = "Low throughput"
		analysis.TPHRecommendation = "Focus on removing blockers and improving flow"
	}

	// Overall assessment
	if hir < 0.5 && aac > 0.5 && tph > 1.5 {
		analysis.OverallAssessment = "Good AI leverage with high productivity"
	} else if hir > 0.8 && aac < 0.3 {
		analysis.OverallAssessment = "Manual craftsmanship - consider AI assistance"
	} else if hir < 0.3 && tph < 1.0 {
		analysis.OverallAssessment = "Possible over-prompting - focus on quality"
	} else {
		analysis.OverallAssessment = "Balanced development approach"
	}

	return analysis
}

// AIImpactAnalysis provides insights on AI usage patterns
type AIImpactAnalysis struct {
	HIRInsight        string `json:"hir_insight"`
	HIRRecommendation string `json:"hir_recommendation"`
	AACInsight        string `json:"aac_insight"`
	AACRecommendation string `json:"aac_recommendation"`
	TPHInsight        string `json:"tph_insight"`
	TPHRecommendation string `json:"tph_recommendation"`
	OverallAssessment string `json:"overall_assessment"`
}
