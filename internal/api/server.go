// Package api implements Repository Intelligence HTTP APIs.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/scorecard"
	"github.com/kubex-ecosystem/grompt/internal/types"
)

// GromptAPI handles Repository Intelligence API endpoints
type GromptAPI struct {
	engine *scorecard.Engine
}

// NewGromptAPI creates a new grompt API handler
func NewGromptAPI(engine *scorecard.Engine) *GromptAPI {
	return &GromptAPI{
		engine: engine,
	}
}

// RegisterRoutes registers all grompt API routes
func (a *GromptAPI) RegisterRoutes(mux *http.ServeMux) {
	// Core scorecard endpoints
	mux.HandleFunc("/api/v1/scorecard", a.handleScorecard)
	mux.HandleFunc("/api/v1/scorecard/advice", a.handleScorecardAdvice)
	mux.HandleFunc("/api/v1/metrics/ai", a.handleAIMetrics)

	// Asset endpoints
	mux.HandleFunc("/api/v1/scorecard/assets/", a.handleAssets)

	// Health endpoint
	mux.HandleFunc("/api/v1/health", a.handleHealth)
}

// handleScorecard handles GET /api/v1/scorecard
func (a *GromptAPI) handleScorecard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	repo := r.URL.Query().Get("repo")
	if repo == "" {
		http.Error(w, "Missing 'repo' parameter", http.StatusBadRequest)
		return
	}

	periodStr := r.URL.Query().Get("period")
	period := 60 // Default to 60 days
	if periodStr != "" {
		if p, err := strconv.Atoi(periodStr); err == nil && p > 0 {
			period = p
		}
	}

	user := r.URL.Query().Get("user")
	if user == "" {
		user = "current-user" // Default or extract from auth
	}

	// Parse repository
	repository, err := parseRepository(repo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid repository format: %v", err), http.StatusBadRequest)
		return
	}

	// Generate scorecard
	scorecard, err := a.engine.GenerateScorecard(r.Context(), *repository, user, period)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate scorecard: %v", err), http.StatusInternalServerError)
		return
	}

	// Set headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Schema-Version", "scorecard@1.0.0")
	w.Header().Set("X-Server-Version", "grompt-v1.0.0")
	w.Header().Set("Cache-Control", "max-age=300") // 5 minutes cache

	// Return scorecard
	json.NewEncoder(w).Encode(scorecard)
}

// handleScorecardAdvice handles POST /api/v1/scorecard/advice
func (a *GromptAPI) handleScorecardAdvice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req ScorecardAdviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if req.Scorecard == nil {
		http.Error(w, "Missing 'scorecard' in request body", http.StatusBadRequest)
		return
	}

	// Generate appropriate report based on mode
	var response interface{}
	var err error

	switch req.Mode {
	case "exec", "executive":
		response, err = a.engine.GenerateExecutiveReport(r.Context(), req.Scorecard, req.Hotspots)
	case "code", "health":
		response, err = a.engine.GenerateCodeHealthReport(r.Context(), req.Scorecard, req.Hotspots)
	case "ops", "dora":
		response, err = a.engine.GenerateDORAReport(r.Context(), req.Scorecard)
	case "community", "bus":
		response, err = a.engine.GenerateCommunityReport(r.Context(), req.Scorecard)
	default:
		http.Error(w, "Invalid mode. Use: exec, code, ops, or community", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate report: %v", err), http.StatusInternalServerError)
		return
	}

	// Set headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Schema-Version", "advice@1.0.0")
	w.Header().Set("X-Server-Version", "grompt-v1.0.0")

	// Return report
	json.NewEncoder(w).Encode(response)
}

// handleAIMetrics handles GET /api/v1/metrics/ai
func (a *GromptAPI) handleAIMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	repo := r.URL.Query().Get("repo")
	if repo == "" {
		http.Error(w, "Missing 'repo' parameter", http.StatusBadRequest)
		return
	}

	periodStr := r.URL.Query().Get("period")
	period := 60 // Default to 60 days
	if periodStr != "" {
		if p, err := strconv.Atoi(periodStr); err == nil && p > 0 {
			period = p
		}
	}

	user := r.URL.Query().Get("user")
	if user == "" {
		user = "current-user"
	}

	// Parse repository
	repository, err := parseRepository(repo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid repository format: %v", err), http.StatusBadRequest)
		return
	}

	// For AI metrics, we need a scorecard first
	scorecard, err := a.engine.GenerateScorecard(r.Context(), *repository, user, period)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate scorecard: %v", err), http.StatusInternalServerError)
		return
	}

	// Create AI metrics response
	aiResponse := AIMetricsResponse{
		SchemaVersion: "ai_metrics@1.0.0",
		Owner:         repository.Owner,
		Repo:          repository.Name,
		PeriodDays:    period,
		Contributors: []ContributorMetrics{
			{
				User: user,
				HIR:  scorecard.AI.HIR,
				AAC:  calculateAAC(scorecard), // Helper function
				TPH:  scorecard.AI.TPH,
				Hours: HoursBreakdown{
					Human: scorecard.AI.HumanHours,
					AI:    scorecard.AI.AIHours,
				},
				Commits: 0, // TODO: Get from Git analysis
			},
		},
		Aggregates: AggregateMetrics{
			HIRP50: scorecard.AI.HIR,
			HIRP90: scorecard.AI.HIR, // TODO: Calculate from multiple contributors
			AAC:    calculateAAC(scorecard),
			TPHP50: scorecard.AI.TPH,
		},
		Provenance: ProvenanceInfo{
			Sources: []string{"git", "wakatime", "ide_telemetry"},
		},
		Confidence: ConfidenceMetrics{
			HIR: scorecard.Confidence.AI,
			AAC: scorecard.Confidence.AI,
			TPH: scorecard.Confidence.AI,
		},
	}

	// Set headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Schema-Version", "ai_metrics@1.0.0")
	w.Header().Set("X-Server-Version", "grompt-v1.0.0")
	w.Header().Set("Cache-Control", "max-age=300")

	// Return AI metrics
	json.NewEncoder(w).Encode(aiResponse)
}

// handleAssets handles GET /api/v1/scorecard/assets/:repo/spark-chi.svg
func (a *GromptAPI) handleAssets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract asset type from path
	// TODO: Implement SVG chart generation for CHI sparklines

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "max-age=3600") // 1 hour cache
	w.Header().Set("ETag", fmt.Sprintf("\"%d\"", time.Now().Unix()))

	// Return simple SVG for now
	svg := `<svg width="120" height="25" xmlns="http://www.w3.org/2000/svg">
		<polyline fill="none" stroke="#00ff00" stroke-width="2"
		points="0,20 20,15 40,10 60,12 80,8 100,5 120,3"/>
	</svg>`

	w.Write([]byte(svg))
}

// handleHealth handles GET /api/v1/health
func (a *GromptAPI) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	health := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services: map[string]string{
			"scorecard_engine": "ok",
			"dora_calculator":  "ok",
			"chi_calculator":   "ok",
			"ai_metrics":       "ok",
		},
		Version: "grompt-v1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// Helper functions

// parseRepository parses repository string (owner/repo format)
func parseRepository(repo string) (*types.Repository, error) {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("repository must be in 'owner/repo' format")
	}

	return &types.Repository{
		Owner:         parts[0],
		Name:          parts[1],
		FullName:      repo,
		CloneURL:      fmt.Sprintf("https://github.com/%s.git", repo),
		DefaultBranch: "main", // Default
		Language:      "unknown",
		IsPrivate:     false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

// calculateAAC calculates AI Assist Coverage from scorecard
func calculateAAC(scorecard *types.Scorecard) float64 {
	// This would be calculated from commit analysis
	// For now, return a derived value
	return scorecard.AI.AAC
}

// Request/Response types

// ScorecardAdviceRequest represents the request for scorecard advice
type ScorecardAdviceRequest struct {
	Scorecard *types.Scorecard `json:"scorecard"`
	Hotspots  []string         `json:"hotspots,omitempty"`
	Mode      string           `json:"mode"` // exec|code|ops|community
}

// AIMetricsResponse represents the AI metrics API response
type AIMetricsResponse struct {
	SchemaVersion string               `json:"schema_version"`
	Owner         string               `json:"owner"`
	Repo          string               `json:"repo"`
	PeriodDays    int                  `json:"period_days"`
	Contributors  []ContributorMetrics `json:"contributors"`
	Aggregates    AggregateMetrics     `json:"aggregates"`
	Provenance    ProvenanceInfo       `json:"provenance"`
	Confidence    ConfidenceMetrics    `json:"confidence"`
}

// ContributorMetrics represents metrics for a single contributor
type ContributorMetrics struct {
	User    string         `json:"user"`
	HIR     float64        `json:"hir"`
	AAC     float64        `json:"aac"`
	TPH     float64        `json:"tph"`
	Hours   HoursBreakdown `json:"hours"`
	Commits int            `json:"commits"`
}

// HoursBreakdown represents time breakdown
type HoursBreakdown struct {
	Human float64 `json:"human"`
	AI    float64 `json:"ai"`
}

// AggregateMetrics represents aggregated team metrics
type AggregateMetrics struct {
	HIRP50 float64 `json:"hir_p50"`
	HIRP90 float64 `json:"hir_p90"`
	AAC    float64 `json:"aac"`
	TPHP50 float64 `json:"tph_p50"`
}

// ProvenanceInfo represents data source information
type ProvenanceInfo struct {
	Sources []string `json:"sources"`
}

// ConfidenceMetrics represents confidence in metrics
type ConfidenceMetrics struct {
	HIR float64 `json:"hir"`
	AAC float64 `json:"aac"`
	TPH float64 `json:"tph"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
}

/*

Cara... Olha os arquivos que inseri nesse contexto aqui!! Por favor!

Pausa só um pouco.. rsrs

*/
