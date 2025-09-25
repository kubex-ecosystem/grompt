// Package scorecard implements the central repository intelligence engine.
package scorecard

import (
	"context"
	"fmt"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/metrics"
	"github.com/kubex-ecosystem/grompt/internal/types"
)

// Engine orchestrates repository analysis and scorecard generation
type Engine struct {
	doraCalculator *metrics.DORACalculator
	chiCalculator  *metrics.CHICalculator
	aiCalculator   *metrics.AIMetricsCalculator
}

// NewEngine creates a new scorecard engine
func NewEngine(
	dora *metrics.DORACalculator,
	chi *metrics.CHICalculator,
	ai *metrics.AIMetricsCalculator,
) *Engine {
	return &Engine{
		doraCalculator: dora,
		chiCalculator:  chi,
		aiCalculator:   ai,
	}
}

// GenerateScorecard creates a comprehensive repository scorecard
func (e *Engine) GenerateScorecard(ctx context.Context, repo types.Repository, user string, periodDays int) (*types.Scorecard, error) {
	// Calculate DORA metrics
	doraMetrics, err := e.doraCalculator.Calculate(ctx, repo, periodDays)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate DORA metrics: %w", err)
	}

	// Calculate Code Health Index
	chiMetrics, err := e.chiCalculator.Calculate(ctx, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate CHI metrics: %w", err)
	}

	// Calculate AI Impact metrics
	aiMetrics, err := e.aiCalculator.Calculate(ctx, repo, user, periodDays)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate AI metrics: %w", err)
	}

	// Calculate additional metrics
	busFactor := e.calculateBusFactor(ctx, repo)
	firstReviewP50 := e.calculateFirstReviewP50(ctx, repo, periodDays)

	// Calculate confidence scores
	confidence := e.calculateConfidence(doraMetrics, chiMetrics, aiMetrics)

	scorecard := &types.Scorecard{
		SchemaVersion:       "scorecard@1.0.0",
		Repository:          repo,
		DORA:                *doraMetrics,
		CHI:                 *chiMetrics,
		AI:                  *aiMetrics,
		BusFactor:           busFactor,
		FirstReviewP50Hours: firstReviewP50,
		Confidence:          confidence,
		GeneratedAt:         time.Now(),
	}

	return scorecard, nil
}

// GenerateExecutiveReport generates P1 Executive "North & South" report
func (e *Engine) GenerateExecutiveReport(ctx context.Context, scorecard *types.Scorecard, hotspots []string) (*types.ExecutiveReport, error) {
	summary := e.generateExecutiveSummary(scorecard)
	topFocus := e.generateTopFocus(scorecard)
	quickWins := e.generateQuickWins(scorecard, hotspots)
	risks := e.generateRisks(scorecard)
	callToAction := e.generateCallToAction(scorecard)

	return &types.ExecutiveReport{
		Summary:      summary,
		TopFocus:     topFocus,
		QuickWins:    quickWins,
		Risks:        risks,
		CallToAction: callToAction,
	}, nil
}

// GenerateCodeHealthReport generates P2 Code Health Deep Dive report
func (e *Engine) GenerateCodeHealthReport(ctx context.Context, scorecard *types.Scorecard, hotspots []string) (*types.CodeHealthReport, error) {
	drivers := e.identifyCHIDrivers(scorecard)
	refactorPlan := e.generateRefactorPlan(scorecard, drivers)
	guardrails := e.generateGuardrails(scorecard)
	milestones := e.generateMilestones(refactorPlan)

	return &types.CodeHealthReport{
		CHINow:       scorecard.CHI.Score,
		Drivers:      drivers,
		RefactorPlan: refactorPlan,
		Guardrails:   guardrails,
		Milestones:   milestones,
	}, nil
}

// GenerateDORAReport generates P3 DORA & Ops report
func (e *Engine) GenerateDORAReport(ctx context.Context, scorecard *types.Scorecard) (*types.DORAReport, error) {
	bottlenecks := e.identifyBottlenecks(scorecard)
	playbook := e.generatePlaybook(scorecard, bottlenecks)
	experiments := e.generateExperiments(scorecard)

	return &types.DORAReport{
		LeadTimeP95Hours:        scorecard.DORA.LeadTimeP95Hours,
		DeploymentFrequencyWeek: scorecard.DORA.DeploymentFrequencyWeek,
		ChangeFailRatePercent:   scorecard.DORA.ChangeFailRatePercent,
		MTTRHours:               scorecard.DORA.MTTRHours,
		Bottlenecks:             bottlenecks,
		Playbook:                playbook,
		Experiments:             experiments,
	}, nil
}

// GenerateCommunityReport generates P4 Community & Bus Factor report
func (e *Engine) GenerateCommunityReport(ctx context.Context, scorecard *types.Scorecard) (*types.CommunityReport, error) {
	roadmap := e.generateCommunityRoadmap(scorecard)
	visibility := e.generateVisibilityItems(scorecard)

	return &types.CommunityReport{
		BusFactor:         scorecard.BusFactor,
		OnboardingP50Days: int(scorecard.FirstReviewP50Hours / 24), // Convert hours to days
		Roadmap:           roadmap,
		Visibility:        visibility,
	}, nil
}

// calculateBusFactor estimates the bus factor (key person dependency)
func (e *Engine) calculateBusFactor(ctx context.Context, repo types.Repository) int {
	// Simplified bus factor calculation
	// In production, analyze commit authorship distribution
	return 2 // Default conservative estimate
}

// calculateFirstReviewP50 calculates median first review time
func (e *Engine) calculateFirstReviewP50(ctx context.Context, repo types.Repository, periodDays int) float64 {
	// Simplified calculation
	// In production, analyze PR review times from GitHub API
	return 8.0 // 8 hours default
}

// calculateConfidence calculates confidence scores for metrics
func (e *Engine) calculateConfidence(dora *types.DORAMetrics, chi *types.CHIMetrics, ai *types.AIMetrics) types.Confidence {
	// Confidence based on data completeness and recency
	doraConfidence := 0.8
	chiConfidence := 0.9
	aiConfidence := 0.7

	if dora.Period < 30 {
		doraConfidence -= 0.2
	}
	if ai.HumanHours < 10 {
		aiConfidence -= 0.3
	}

	groupConfidence := (doraConfidence + chiConfidence + aiConfidence) / 3.0

	return types.Confidence{
		DORA:  doraConfidence,
		CHI:   chiConfidence,
		AI:    aiConfidence,
		Group: groupConfidence,
	}
}

// generateExecutiveSummary creates executive summary
func (e *Engine) generateExecutiveSummary(scorecard *types.Scorecard) types.ExecutiveSummary {
	grade := "C"
	if scorecard.CHI.Score >= 80 && scorecard.DORA.LeadTimeP95Hours < 48 {
		grade = "A"
	} else if scorecard.CHI.Score >= 60 && scorecard.DORA.LeadTimeP95Hours < 72 {
		grade = "B"
	} else if scorecard.CHI.Score < 40 || scorecard.DORA.LeadTimeP95Hours > 120 {
		grade = "D"
	}

	return types.ExecutiveSummary{
		Grade:            grade,
		CHI:              scorecard.CHI.Score,
		LeadTimeP95Hours: scorecard.DORA.LeadTimeP95Hours,
		DeploysPerWeek:   scorecard.DORA.DeploymentFrequencyWeek,
	}
}

// generateTopFocus identifies top 3 focus areas
func (e *Engine) generateTopFocus(scorecard *types.Scorecard) []types.FocusArea {
	var focus []types.FocusArea

	// Focus 1: Code Health if CHI is low
	if scorecard.CHI.Score < 60 {
		focus = append(focus, types.FocusArea{
			Title:      "Improve Code Health",
			Why:        "CHI score below threshold indicates technical debt",
			KPI:        "CHI Score",
			Target:     "≥ 70",
			Confidence: 0.9,
		})
	}

	// Focus 2: Lead Time if too high
	if scorecard.DORA.LeadTimeP95Hours > 72 {
		focus = append(focus, types.FocusArea{
			Title:      "Reduce Lead Time",
			Why:        "P95 lead time above 3 days impacts delivery speed",
			KPI:        "Lead Time P95",
			Target:     "≤ 48 hours",
			Confidence: 0.8,
		})
	}

	// Focus 3: Bus Factor if too low
	if scorecard.BusFactor <= 1 {
		focus = append(focus, types.FocusArea{
			Title:      "Increase Bus Factor",
			Why:        "Single point of failure in team knowledge",
			KPI:        "Bus Factor",
			Target:     "≥ 2",
			Confidence: 0.7,
		})
	}

	// Ensure we have exactly 3 focus areas
	if len(focus) < 3 {
		focus = append(focus, types.FocusArea{
			Title:      "Maintain Quality",
			Why:        "Continue current good practices",
			KPI:        "Overall Score",
			Target:     "Maintain current level",
			Confidence: 0.8,
		})
	}

	if len(focus) > 3 {
		focus = focus[:3]
	}

	return focus
}

// generateQuickWins identifies quick wins
func (e *Engine) generateQuickWins(scorecard *types.Scorecard, hotspots []string) []types.QuickWin {
	var wins []types.QuickWin

	if scorecard.CHI.DuplicationPercent > 20 {
		wins = append(wins, types.QuickWin{
			Action:       "Add duplication detection to CI",
			Effort:       "S",
			ExpectedGain: "Prevent new duplications",
		})
	}

	if scorecard.DORA.LeadTimeP95Hours > 72 {
		wins = append(wins, types.QuickWin{
			Action:       "Implement PR size limits (≤300 LOC)",
			Effort:       "M",
			ExpectedGain: "Reduce lead time by 20%",
		})
	}

	if scorecard.CHI.TestCoverage < 50 {
		wins = append(wins, types.QuickWin{
			Action:       "Add coverage reporting to CI",
			Effort:       "S",
			ExpectedGain: "Visibility into test gaps",
		})
	}

	return wins
}

// generateRisks identifies top risks
func (e *Engine) generateRisks(scorecard *types.Scorecard) []types.Risk {
	var risks []types.Risk

	if scorecard.BusFactor <= 1 {
		risks = append(risks, types.Risk{
			Risk:       "Single point of failure in team knowledge",
			Mitigation: "Document processes and cross-train team members",
		})
	}

	if scorecard.DORA.ChangeFailRatePercent > 15 {
		risks = append(risks, types.Risk{
			Risk:       "High change failure rate indicates instability",
			Mitigation: "Implement feature flags and better testing",
		})
	}

	if scorecard.CHI.TechnicalDebt > 100 {
		risks = append(risks, types.Risk{
			Risk:       "High technical debt slowing development",
			Mitigation: "Schedule dedicated refactoring sprints",
		})
	}

	return risks
}

// generateCallToAction creates call to action
func (e *Engine) generateCallToAction(scorecard *types.Scorecard) string {
	if scorecard.CHI.Score < 50 {
		return "Immediate action required: Code health is critically low. Focus on reducing duplication and complexity."
	}
	if scorecard.DORA.LeadTimeP95Hours > 120 {
		return "Urgent: Lead time is too high. Implement smaller PRs and faster reviews."
	}
	return "Continue monitoring metrics and maintain current quality standards."
}

// identifyCHIDrivers identifies what's driving CHI score
func (e *Engine) identifyCHIDrivers(scorecard *types.Scorecard) []types.CHIDriver {
	var drivers []types.CHIDriver

	if scorecard.CHI.DuplicationPercent > 15 {
		drivers = append(drivers, types.CHIDriver{
			Metric: "duplication_pct",
			Value:  scorecard.CHI.DuplicationPercent,
			Impact: "high",
		})
	}

	if scorecard.CHI.CyclomaticComplexity > 10 {
		drivers = append(drivers, types.CHIDriver{
			Metric: "cyclomatic_avg",
			Value:  scorecard.CHI.CyclomaticComplexity,
			Impact: "medium",
		})
	}

	if scorecard.CHI.MaintainabilityIndex < 60 {
		drivers = append(drivers, types.CHIDriver{
			Metric: "mi",
			Value:  scorecard.CHI.MaintainabilityIndex,
			Impact: "medium",
		})
	}

	return drivers
}

// generateRefactorPlan creates incremental refactoring plan
func (e *Engine) generateRefactorPlan(scorecard *types.Scorecard, drivers []types.CHIDriver) []types.RefactorStep {
	var plan []types.RefactorStep
	step := 1

	for _, driver := range drivers {
		switch driver.Metric {
		case "duplication_pct":
			plan = append(plan, types.RefactorStep{
				Step:    step,
				Theme:   "duplication",
				Actions: []string{"Identify duplicated code blocks", "Extract common functions", "Create shared libraries"},
				KPI:     "Duplication Percentage",
				Target:  "≤ 15%",
			})
			step++
		case "cyclomatic_avg":
			plan = append(plan, types.RefactorStep{
				Step:    step,
				Theme:   "complexity",
				Actions: []string{"Break down complex functions", "Extract helper methods", "Simplify conditional logic"},
				KPI:     "Average Complexity",
				Target:  "≤ 8",
			})
			step++
		}
	}

	return plan
}

// generateGuardrails creates quality guardrails
func (e *Engine) generateGuardrails(scorecard *types.Scorecard) []string {
	var guardrails []string

	guardrails = append(guardrails, "Add lint checks to CI pipeline")

	if scorecard.CHI.DuplicationPercent > 15 {
		guardrails = append(guardrails, "Threshold: duplication_pct ≤ 20%")
	}

	if scorecard.CHI.TestCoverage < 60 {
		guardrails = append(guardrails, "Minimum test coverage per package: 60%")
	}

	guardrails = append(guardrails, "Maximum PR size: 300 LOC")

	return guardrails
}

// generateMilestones creates refactoring milestones
func (e *Engine) generateMilestones(plan []types.RefactorStep) []types.Milestone {
	return []types.Milestone{
		{InDays: 14, Goal: "Complete first refactoring step"},
		{InDays: 30, Goal: "Achieve 10-point CHI improvement"},
		{InDays: 60, Goal: "Implement all quality guardrails"},
	}
}

// identifyBottlenecks identifies DORA bottlenecks
func (e *Engine) identifyBottlenecks(scorecard *types.Scorecard) []types.Bottleneck {
	var bottlenecks []types.Bottleneck

	if scorecard.FirstReviewP50Hours > 12 {
		bottlenecks = append(bottlenecks, types.Bottleneck{
			Area:     "review",
			Evidence: "First review P50 > 12 hours",
		})
	}

	if scorecard.DORA.LeadTimeP95Hours > 72 {
		bottlenecks = append(bottlenecks, types.Bottleneck{
			Area:     "batch_size",
			Evidence: "Large PRs causing long lead times",
		})
	}

	return bottlenecks
}

// generatePlaybook creates operational playbook
func (e *Engine) generatePlaybook(scorecard *types.Scorecard, bottlenecks []types.Bottleneck) []types.PlaybookItem {
	var playbook []types.PlaybookItem

	playbook = append(playbook, types.PlaybookItem{
		Name:           "Smaller PRs",
		Policy:         "Max 300 LOC or 5 files per PR",
		ExpectedEffect: "Lead time P95 -20%",
	})

	playbook = append(playbook, types.PlaybookItem{
		Name:           "Review SLA",
		Policy:         "P50 ≤ 8h, P90 ≤ 24h",
		ExpectedEffect: "Lead time -15%",
	})

	return playbook
}

// generateExperiments creates A/B experiments
func (e *Engine) generateExperiments(scorecard *types.Scorecard) []types.Experiment {
	return []types.Experiment{
		{
			AB:           "Pair programming sessions vs individual work",
			Metric:       "lead_time_p95",
			DurationDays: 14,
		},
		{
			AB:           "Automated tests vs manual reviews",
			Metric:       "CFR",
			DurationDays: 21,
		},
	}
}

// generateCommunityRoadmap creates community growth roadmap
func (e *Engine) generateCommunityRoadmap(scorecard *types.Scorecard) []types.RoadmapItem {
	var roadmap []types.RoadmapItem

	roadmap = append(roadmap, types.RoadmapItem{
		Item:          "CONTRIBUTING guide + PR templates",
		Why:           "Reduce onboarding friction",
		SuccessMetric: "First review P50 ≤ 8h",
	})

	if scorecard.BusFactor <= 1 {
		roadmap = append(roadmap, types.RoadmapItem{
			Item:          "Create good-first-issues (5 issues)",
			Why:           "Attract new contributors",
			SuccessMetric: "Bus factor ≥ 2",
		})
	}

	return roadmap
}

// generateVisibilityItems creates visibility improvement items
func (e *Engine) generateVisibilityItems(scorecard *types.Scorecard) []types.VisibilityItem {
	return []types.VisibilityItem{
		{
			Asset:  "README with clear value proposition",
			KPI:    "Stars growth trend",
			Effort: "S",
		},
		{
			Asset:  "Documentation site",
			KPI:    "Contributor onboarding time",
			Effort: "M",
		},
	}
}
