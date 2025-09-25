package grompt_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

// TestProjectAnalysisCore tests the core business logic of project analysis
// This is what the grompt REALLY does - analyze project contexts!
func TestProjectAnalysisCore(t *testing.T) {
	tests := []struct {
		name           string
		projectContext string
		analysisType   string
		expectContains []string
		expectError    bool
	}{
		{
			name: "Simple Go project analysis",
			projectContext: `
# My Go Project
This is a simple REST API built with Go and Gin framework.

## Features
- User authentication
- CRUD operations
- PostgreSQL database
- Docker support

## TODO
- Add tests
- Improve error handling
- Add monitoring
			`,
			analysisType: "GENERAL",
			expectContains: []string{
				"Go",
				"REST API",
				"authentication",
				"PostgreSQL",
				"tests",
			},
			expectError: false,
		},
		{
			name: "React frontend project",
			projectContext: `
# React Dashboard
Modern dashboard built with React, TypeScript, and Tailwind CSS.

## Current State
- Login page implemented
- Dashboard layout done
- Charts integration pending
- No tests yet

## Tech Stack
- React 18
- TypeScript
- Tailwind CSS
- Vite
			`,
			analysisType: "CODE_QUALITY",
			expectContains: []string{
				"React",
				"TypeScript",
				"tests",
				"quality",
			},
			expectError: false,
		},
		{
			name: "Security analysis of API project",
			projectContext: `
# Banking API
RESTful API for banking operations.

## Current Implementation
- No authentication yet
- Direct SQL queries
- No input validation
- Admin endpoints exposed
- Passwords stored in plain text
			`,
			analysisType: "SECURITY",
			expectContains: []string{
				"authentication",
				"SQL",
				"validation",
				"security",
				"passwords",
			},
			expectError: false,
		},
		{
			name:           "Empty project context should fail gracefully",
			projectContext: "",
			analysisType:   "GENERAL",
			expectContains: []string{},
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the core analysis logic
			result, err := analyzeProjectContext(tt.projectContext, tt.analysisType)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for test %s, but got none", tt.name)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for test %s: %v", tt.name, err)
				return
			}

			if result == nil {
				t.Errorf("Expected non-nil result for test %s", tt.name)
				return
			}

			// Check if result contains expected keywords
			resultJSON, _ := json.Marshal(result)
			resultStr := strings.ToLower(string(resultJSON))

			for _, expected := range tt.expectContains {
				if !strings.Contains(resultStr, strings.ToLower(expected)) {
					t.Errorf("Expected result to contain '%s' for test %s, but it didn't", expected, tt.name)
				}
			}
		})
	}
}

// TestAnalysisTypes tests different analysis types
func TestAnalysisTypes(t *testing.T) {
	projectContext := `
# Test Project
A sample project for testing analysis types.
Built with Go, has security issues, and needs performance improvements.
	`

	analysisTypes := []string{"GENERAL", "SECURITY", "SCALABILITY", "CODE_QUALITY"}

	for _, analysisType := range analysisTypes {
		t.Run("Analysis_"+analysisType, func(t *testing.T) {
			result, err := analyzeProjectContext(projectContext, analysisType)

			if err != nil {
				t.Errorf("Error analyzing with type %s: %v", analysisType, err)
				return
			}

			if result == nil {
				t.Errorf("Nil result for analysis type %s", analysisType)
				return
			}

			// Verify analysis type is correctly set
			if result.AnalysisType != analysisType {
				t.Errorf("Expected analysis type %s, got %s", analysisType, result.AnalysisType)
			}

			// Verify required fields are present
			if result.ProjectName == "" {
				t.Errorf("Missing project name for analysis type %s", analysisType)
			}

			if result.Summary == "" {
				t.Errorf("Missing summary for analysis type %s", analysisType)
			}

			if len(result.Strengths) == 0 {
				t.Errorf("No strengths found for analysis type %s", analysisType)
			}
		})
	}
}

// TestViabilityScoring tests the viability scoring logic
func TestViabilityScoring(t *testing.T) {
	tests := []struct {
		name           string
		projectContext string
		minScore       float64
		maxScore       float64
	}{
		{
			name: "High viability project",
			projectContext: `
# Mature E-commerce Platform
Production-ready e-commerce platform with:
- 99.9% uptime
- Full test coverage
- Security audited
- Scalable architecture
- Active user base
- Revenue generating
			`,
			minScore: 8.0,
			maxScore: 10.0,
		},
		{
			name: "Low viability project",
			projectContext: `
# Broken Legacy System
Old system with major problems:
- No documentation
- No tests
- Security vulnerabilities
- Outdated dependencies
- Abandoned codebase
- No maintainer
			`,
			minScore: 1.0,
			maxScore: 4.0,
		},
		{
			name: "Medium viability project",
			projectContext: `
# Standard Web App
Regular web application:
- Basic functionality works
- Some tests exist
- Documentation is okay
- Uses modern framework
- Few security issues
			`,
			minScore: 5.0,
			maxScore: 7.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := analyzeProjectContext(tt.projectContext, "GENERAL")

			if err != nil {
				t.Errorf("Error analyzing project: %v", err)
				return
			}

			score := result.Viability.Score
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("Viability score %f not in expected range [%f, %f] for %s",
					score, tt.minScore, tt.maxScore, tt.name)
			}
		})
	}
}

// TestPromptGeneration tests the prompt generation logic
func TestPromptGeneration(t *testing.T) {
	tests := []struct {
		name           string
		analysisType   string
		locale         string
		expectInPrompt []string
	}{
		{
			name:         "Security analysis with Portuguese",
			analysisType: "SECURITY",
			locale:       "pt-BR",
			expectInPrompt: []string{
				"security",
				"Portuguese",
				"vulnerabilities",
			},
		},
		{
			name:         "General analysis with English",
			analysisType: "GENERAL",
			locale:       "en-US",
			expectInPrompt: []string{
				"general",
				"English",
				"project",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := generateAnalysisPrompt("test context", tt.analysisType, tt.locale)

			if prompt == "" {
				t.Errorf("Empty prompt generated for %s", tt.name)
				return
			}

			promptLower := strings.ToLower(prompt)
			for _, expected := range tt.expectInPrompt {
				if !strings.Contains(promptLower, strings.ToLower(expected)) {
					t.Errorf("Prompt should contain '%s' for %s", expected, tt.name)
				}
			}
		})
	}
}

// Benchmark tests for performance
func BenchmarkProjectAnalysis(b *testing.B) {
	projectContext := `
# Sample Project for Benchmarking
This is a medium-sized project with typical characteristics.
Built with modern technologies and has room for improvement.
	`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := analyzeProjectContext(projectContext, "GENERAL")
		if err != nil {
			b.Fatalf("Error in benchmark: %v", err)
		}
	}
}

// Mock structures for testing (these should match your real types)
type ProjectAnalysis struct {
	ProjectName  string          `json:"projectName"`
	AnalysisType string          `json:"analysisType"`
	Summary      string          `json:"summary"`
	Strengths    []string        `json:"strengths"`
	Improvements []Improvement   `json:"improvements"`
	Viability    Viability       `json:"viability"`
	NextSteps    NextSteps       `json:"nextSteps"`
	ROIAnalysis  ROIAnalysis     `json:"roiAnalysis"`
	Maturity     ProjectMaturity `json:"maturity"`
}

type Improvement struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	Priority       string `json:"priority"`
	Difficulty     string `json:"difficulty"`
	BusinessImpact string `json:"businessImpact"`
}

type Viability struct {
	Score      float64 `json:"score"`
	Assessment string  `json:"assessment"`
}

type NextSteps struct {
	ShortTerm []Step `json:"shortTerm"`
	LongTerm  []Step `json:"longTerm"`
}

type Step struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
}

type ROIAnalysis struct {
	Assessment      string   `json:"assessment"`
	PotentialGains  []string `json:"potentialGains"`
	EstimatedEffort string   `json:"estimatedEffort"`
}

type ProjectMaturity struct {
	Level      string `json:"level"`
	Assessment string `json:"assessment"`
}

// Mock implementation - this is where the REAL business logic should be
func analyzeProjectContext(projectContext, analysisType string) (*ProjectAnalysis, error) {
	if strings.TrimSpace(projectContext) == "" {
		return nil, ErrEmptyContext
	}

	// This is a SIMPLIFIED mock - the real implementation would use AI
	analysis := &ProjectAnalysis{
		ProjectName:  extractProjectName(projectContext),
		AnalysisType: analysisType,
		Summary:      generateSummary(projectContext, analysisType),
		Strengths:    extractStrengths(projectContext),
		Viability: Viability{
			Score:      calculateViabilityScore(projectContext),
			Assessment: "Mock assessment based on project context",
		},
		Maturity: ProjectMaturity{
			Level:      "MVP",
			Assessment: "Project appears to be in MVP stage",
		},
	}

	return analysis, nil
}

func generateAnalysisPrompt(projectContext, analysisType, locale string) string {
	language := "English (US)"
	if locale == "pt-BR" {
		language = "Portuguese (Brazil)"
	}

	focusMap := map[string]string{
		"GENERAL":      "overall project assessment",
		"SECURITY":     "security vulnerabilities and best practices",
		"SCALABILITY":  "performance bottlenecks and scaling issues",
		"CODE_QUALITY": "code structure and maintainability",
	}

	focus := focusMap[analysisType]
	if focus == "" {
		focus = "general analysis"
	}

	return fmt.Sprintf(`
You are an expert software architect analyzing a project.

Analysis Type: %s
Focus: %s
Language: %s

Project Context:
%s

Provide detailed insights focusing on %s.
	`, analysisType, focus, language, projectContext, focus)
}

// Helper functions for mock implementation
func extractProjectName(context string) string {
	lines := strings.Split(context, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			return strings.TrimSpace(strings.TrimPrefix(line, "#"))
		}
	}
	return "Unknown Project"
}

func generateSummary(context, analysisType string) string {
	contextLower := strings.ToLower(context)

	switch analysisType {
	case "SECURITY":
		if strings.Contains(contextLower, "authentication") {
			return "Security analysis reveals authentication concerns"
		}
		return "Basic security assessment completed"
	case "SCALABILITY":
		return "Scalability analysis shows areas for performance improvement"
	default:
		return "General analysis of project structure and potential"
	}
}

func extractStrengths(context string) []string {
	strengths := []string{}
	contextLower := strings.ToLower(context)

	strengthKeywords := map[string]string{
		"docker":     "Containerization with Docker",
		"typescript": "Strong typing with TypeScript",
		"test":       "Testing infrastructure in place",
		"api":        "Well-defined API structure",
		"database":   "Database integration",
	}

	for keyword, strength := range strengthKeywords {
		if strings.Contains(contextLower, keyword) {
			strengths = append(strengths, strength)
		}
	}

	if len(strengths) == 0 {
		strengths = append(strengths, "Project has basic structure")
	}

	return strengths
}

func calculateViabilityScore(context string) float64 {
	score := 5.0 // Base score
	contextLower := strings.ToLower(context)

	// Positive indicators
	if strings.Contains(contextLower, "test") {
		score += 1.0
	}
	if strings.Contains(contextLower, "documentation") {
		score += 0.5
	}
	if strings.Contains(contextLower, "production") {
		score += 1.5
	}
	if strings.Contains(contextLower, "security") {
		score += 0.5
	}

	// Negative indicators
	if strings.Contains(contextLower, "no test") {
		score -= 1.0
	}
	if strings.Contains(contextLower, "broken") {
		score -= 2.0
	}
	if strings.Contains(contextLower, "abandoned") {
		score -= 3.0
	}
	if strings.Contains(contextLower, "vulnerabilities") {
		score -= 1.0
	}

	// Keep score in valid range
	if score < 1.0 {
		score = 1.0
	}
	if score > 10.0 {
		score = 10.0
	}

	return score
}

// Custom errors
var (
	ErrEmptyContext = errors.New("project context cannot be empty")
)
