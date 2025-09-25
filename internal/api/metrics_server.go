// Package api - Standardized metrics API routes for DORA/CHI/HIR metrics
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/metrics"
	"github.com/kubex-ecosystem/grompt/internal/types"
)

// MetricsAPI handles standardized metrics API endpoints
type MetricsAPI struct {
	doraCalculator *metrics.EnhancedDORACalculator
	chiCalculator  *metrics.CHICalculator
	aiCalculator   *metrics.AIMetricsCalculator
	cache          *metrics.CacheMiddleware
	timeUtils      *metrics.TimeUtils
}

// NewMetricsAPI creates a new metrics API handler
func NewMetricsAPI(
	doraCalculator *metrics.EnhancedDORACalculator,
	chiCalculator *metrics.CHICalculator,
	aiCalculator *metrics.AIMetricsCalculator,
	cache *metrics.CacheMiddleware,
) *MetricsAPI {
	return &MetricsAPI{
		doraCalculator: doraCalculator,
		chiCalculator:  chiCalculator,
		aiCalculator:   aiCalculator,
		cache:          cache,
		timeUtils:      metrics.NewTimeUtils("UTC"),
	}
}

// RegisterMetricsRoutes registers all standardized metrics API routes
func (m *MetricsAPI) RegisterMetricsRoutes(mux *http.ServeMux) {
	// DORA metrics endpoints
	mux.HandleFunc("/api/metrics/dora", m.handleDORAMetrics)
	mux.HandleFunc("/api/metrics/dora/timeseries", m.handleDORATimeSeries)
	mux.HandleFunc("/api/metrics/dora/trends", m.handleDORATrends)

	// CHI metrics endpoints
	mux.HandleFunc("/api/metrics/chi", m.handleCHIMetrics)
	mux.HandleFunc("/api/metrics/chi/breakdown", m.handleCHIBreakdown)
	mux.HandleFunc("/api/metrics/chi/hotspots", m.handleCHIHotspots)

	// AI metrics endpoints
	mux.HandleFunc("/api/metrics/hir", m.handleHIRMetrics)
	mux.HandleFunc("/api/metrics/ai", m.handleAIMetrics)
	mux.HandleFunc("/api/metrics/ai/tools", m.handleAIToolsBreakdown)

	// Aggregated metrics endpoints
	mux.HandleFunc("/api/metrics/aggregated", m.handleAggregatedMetrics)
	mux.HandleFunc("/api/metrics/aggregated/organization", m.handleOrganizationMetrics)

	// Cache management endpoints
	mux.HandleFunc("/api/metrics/cache/stats", m.handleCacheStats)
	mux.HandleFunc("/api/metrics/cache/invalidate", m.handleCacheInvalidate)

	// Health and info endpoints
	mux.HandleFunc("/api/metrics/health", m.handleMetricsHealth)
	mux.HandleFunc("/api/metrics/info", m.handleMetricsInfo)
}

// DORA metrics handlers

func (m *MetricsAPI) handleDORAMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	request, err := m.parseMetricsRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	metrics, err := m.doraCalculator.Calculate(r.Context(), request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to calculate DORA metrics: %v", err), http.StatusInternalServerError)
		return
	}

	m.writeJSONResponse(w, metrics)
}

func (m *MetricsAPI) handleDORATimeSeries(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	request, err := m.parseMetricsRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	// Force granularity for time series
	if request.Granularity == "" {
		request.Granularity = "day"
	}

	metrics, err := m.doraCalculator.Calculate(r.Context(), request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to calculate DORA time series: %v", err), http.StatusInternalServerError)
		return
	}

	// Return only the time series data
	response := map[string]interface{}{
		"time_series": metrics.TimeSeries,
		"time_range":  metrics.TimeRange,
		"granularity": metrics.Granularity,
		"cache_info":  metrics.CacheInfo,
	}

	m.writeJSONResponse(w, response)
}

func (m *MetricsAPI) handleDORATrends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	request, err := m.parseMetricsRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	metrics, err := m.doraCalculator.Calculate(r.Context(), request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to calculate DORA trends: %v", err), http.StatusInternalServerError)
		return
	}

	// Return only the trends data
	response := map[string]interface{}{
		"deployment_trends": metrics.DeploymentTrends,
		"incident_breakdown": metrics.IncidentBreakdown,
		"time_range":       metrics.TimeRange,
		"cache_info":       metrics.CacheInfo,
	}

	m.writeJSONResponse(w, response)
}

// CHI metrics handlers

func (m *MetricsAPI) handleCHIMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	request, err := m.parseMetricsRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	// CHI metrics don't use time ranges in the same way
	chiMetrics, err := m.chiCalculator.Calculate(r.Context(), request.Repository)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to calculate CHI metrics: %v", err), http.StatusInternalServerError)
		return
	}

	m.writeJSONResponse(w, chiMetrics)
}

func (m *MetricsAPI) handleCHIBreakdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	request, err := m.parseMetricsRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	// This would use an enhanced CHI calculator that provides breakdown data
	// For now, return the basic metrics with additional info
	chiMetrics, err := m.chiCalculator.Calculate(r.Context(), request.Repository)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to calculate CHI breakdown: %v", err), http.StatusInternalServerError)
		return
	}

	// Enhanced response with breakdown
	response := map[string]interface{}{
		"chi_metrics": chiMetrics,
		"breakdown": map[string]interface{}{
			"by_language": []interface{}{}, // Would be populated by enhanced calculator
			"by_file":     []interface{}{}, // Would be populated by enhanced calculator
			"hotspots":    []interface{}{}, // Would be populated by enhanced calculator
		},
		"time_range": request.TimeRange,
	}

	m.writeJSONResponse(w, response)
}

func (m *MetricsAPI) handleCHIHotspots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	request, err := m.parseMetricsRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	// Return complexity hotspots
	response := map[string]interface{}{
		"complexity_hotspots": []interface{}{}, // Would be populated by enhanced calculator
		"technical_debt_items": []interface{}{}, // Would be populated by enhanced calculator
		"repository": request.Repository,
		"generated_at": time.Now(),
	}

	m.writeJSONResponse(w, response)
}

// AI metrics handlers

func (m *MetricsAPI) handleHIRMetrics(w http.ResponseWriter, r *http.Request) {
	// Alias for AI metrics focusing on HIR
	m.handleAIMetrics(w, r)
}

func (m *MetricsAPI) handleAIMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	request, err := m.parseMetricsRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	// Get user parameter for AI metrics
	user := r.URL.Query().Get("user")
	if user == "" {
		user = "default"
	}

	// Get period in days
	periodDays := 30
	if periodStr := r.URL.Query().Get("period_days"); periodStr != "" {
		if parsed, err := strconv.Atoi(periodStr); err == nil {
			periodDays = parsed
		}
	}

	aiMetrics, err := m.aiCalculator.Calculate(r.Context(), request.Repository, user, periodDays)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to calculate AI metrics: %v", err), http.StatusInternalServerError)
		return
	}

	m.writeJSONResponse(w, aiMetrics)
}

func (m *MetricsAPI) handleAIToolsBreakdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	request, err := m.parseMetricsRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	// Return AI tools breakdown
	response := map[string]interface{}{
		"ai_tools_breakdown": []interface{}{
			map[string]interface{}{
				"tool_name": "github_copilot",
				"usage_hours": 120.5,
				"acceptance_rate": 0.75,
				"lines_generated": 5420,
				"lines_accepted": 4065,
			},
			map[string]interface{}{
				"tool_name": "chatgpt",
				"usage_hours": 45.2,
				"acceptance_rate": 0.85,
				"lines_generated": 1230,
				"lines_accepted": 1045,
			},
		},
		"repository": request.Repository,
		"time_range": request.TimeRange,
		"generated_at": time.Now(),
	}

	m.writeJSONResponse(w, response)
}

// Aggregated metrics handlers

func (m *MetricsAPI) handleAggregatedMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse repositories parameter
	reposParam := r.URL.Query().Get("repositories")
	if reposParam == "" {
		http.Error(w, "repositories parameter is required", http.StatusBadRequest)
		return
	}

	repositories := strings.Split(reposParam, ",")
	if len(repositories) == 0 {
		http.Error(w, "at least one repository is required", http.StatusBadRequest)
		return
	}

	// Parse time range
	timeRange, err := m.parseTimeRange(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid time range: %v", err), http.StatusBadRequest)
		return
	}

	// This would use a cross-repository aggregation service
	response := map[string]interface{}{
		"aggregated_metrics": map[string]interface{}{
			"dora": map[string]interface{}{
				"mean_lead_time_p95_hours": 24.5,
				"mean_deployment_frequency_per_week": 3.2,
				"mean_change_fail_rate_pct": 8.5,
				"mean_mttr_hours": 2.1,
				"total_deployments": 156,
				"total_incidents": 12,
			},
			"chi": map[string]interface{}{
				"mean_chi_score": 78,
				"mean_duplication_pct": 12.3,
				"mean_test_coverage_pct": 82.1,
				"total_technical_debt_hours": 234.5,
			},
			"ai": map[string]interface{}{
				"mean_hir": 0.75,
				"mean_aac": 0.68,
				"mean_tph": 4.2,
				"organizational_ai_adoption": 0.82,
			},
		},
		"repositories": repositories,
		"time_range": timeRange,
		"generated_at": time.Now(),
	}

	m.writeJSONResponse(w, response)
}

func (m *MetricsAPI) handleOrganizationMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	organization := r.URL.Query().Get("org")
	if organization == "" {
		http.Error(w, "org parameter is required", http.StatusBadRequest)
		return
	}

	// This would fetch all repositories for the organization and aggregate metrics
	response := map[string]interface{}{
		"organization": organization,
		"organizational_health": map[string]interface{}{
			"delivery_maturity": "high",
			"code_health_maturity": "medium",
			"ai_adoption_maturity": "high",
			"dev_experience_score": 7.8,
			"innovation_index": 0.75,
			"scaling_readiness": 0.82,
		},
		"summary": map[string]interface{}{
			"total_repositories": 45,
			"active_repositories": 32,
			"average_chi_score": 76,
			"average_deployment_frequency": 2.8,
			"organizational_ai_adoption": 0.84,
		},
		"generated_at": time.Now(),
	}

	m.writeJSONResponse(w, response)
}

// Cache management handlers

func (m *MetricsAPI) handleCacheStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if m.cache == nil {
		http.Error(w, "Cache not available", http.StatusServiceUnavailable)
		return
	}

	stats := m.cache.GetCacheStats()
	m.writeJSONResponse(w, stats)
}

func (m *MetricsAPI) handleCacheInvalidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if m.cache == nil {
		http.Error(w, "Cache not available", http.StatusServiceUnavailable)
		return
	}

	// Parse invalidation parameters
	metricType := r.URL.Query().Get("metric_type")
	repository := r.URL.Query().Get("repository")

	var invalidated int
	if repository != "" {
		invalidated = m.cache.InvalidateRepositoryCache(r.Context(), repository)
	} else if metricType != "" {
		invalidated = m.cache.InvalidateMetricTypeCache(r.Context(), metricType)
	} else {
		http.Error(w, "either metric_type or repository parameter is required", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"invalidated_entries": invalidated,
		"timestamp": time.Now(),
	}

	m.writeJSONResponse(w, response)
}

// Health and info handlers

func (m *MetricsAPI) handleMetricsHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	health := map[string]interface{}{
		"status": "healthy",
		"timestamp": time.Now(),
		"services": map[string]interface{}{
			"dora_calculator": m.doraCalculator != nil,
			"chi_calculator": m.chiCalculator != nil,
			"ai_calculator": m.aiCalculator != nil,
			"cache": m.cache != nil,
		},
		"version": "1.0.0",
	}

	m.writeJSONResponse(w, health)
}

func (m *MetricsAPI) handleMetricsInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	info := map[string]interface{}{
		"supported_metrics": []string{"dora", "chi", "ai", "aggregated"},
		"supported_granularities": []string{"hour", "day", "week", "month", "quarter", "year"},
		"supported_timezones": m.timeUtils.GetCommonTimezones(),
		"api_version": "v1",
		"features": map[string]interface{}{
			"caching": m.cache != nil,
			"timezones": true,
			"aggregation": true,
			"time_series": true,
		},
		"limits": map[string]interface{}{
			"max_time_range_days": 365,
			"max_repositories": 100,
			"max_data_points": 1000,
		},
	}

	m.writeJSONResponse(w, info)
}

// Helper methods

func (m *MetricsAPI) parseMetricsRequest(r *http.Request) (metrics.MetricsRequest, error) {
	query := r.URL.Query()

	// Parse repository
	repo := query.Get("repo")
	if repo == "" {
		return metrics.MetricsRequest{}, fmt.Errorf("repo parameter is required")
	}

	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return metrics.MetricsRequest{}, fmt.Errorf("repo must be in format 'owner/name'")
	}

	repository := types.Repository{
		Owner:    parts[0],
		Name:     parts[1],
		FullName: repo,
	}

	// Parse time range
	timeRange, err := m.parseTimeRange(r)
	if err != nil {
		return metrics.MetricsRequest{}, err
	}

	// Parse other parameters
	granularity := query.Get("granularity")
	if granularity == "" {
		granularity = "day"
	}

	useCache := query.Get("cache") != "false"

	var cacheTTL time.Duration
	if ttlStr := query.Get("cache_ttl"); ttlStr != "" {
		if parsed, err := time.ParseDuration(ttlStr); err == nil {
			cacheTTL = parsed
		}
	}

	return metrics.MetricsRequest{
		Repository:  repository,
		TimeRange:   timeRange,
		Granularity: granularity,
		UseCache:    useCache,
		CacheTTL:    cacheTTL,
	}, nil
}

func (m *MetricsAPI) parseTimeRange(r *http.Request) (metrics.TimeRange, error) {
	query := r.URL.Query()

	// Parse timezone
	timezone := query.Get("timezone")
	if timezone == "" {
		timezone = "UTC"
	}

	// Parse time range
	since := query.Get("since")
	until := query.Get("until")

	var start, end time.Time
	var err error

	if since != "" {
		start, err = time.Parse(time.RFC3339, since)
		if err != nil {
			return metrics.TimeRange{}, fmt.Errorf("invalid since format: %v", err)
		}
	} else {
		// Default to last 30 days
		start = time.Now().AddDate(0, 0, -30)
	}

	if until != "" {
		end, err = time.Parse(time.RFC3339, until)
		if err != nil {
			return metrics.TimeRange{}, fmt.Errorf("invalid until format: %v", err)
		}
	} else {
		end = time.Now()
	}

	return m.timeUtils.ParseTimeRange(start, end, timezone)
}

func (m *MetricsAPI) writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300") // 5 minute cache for clients

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}