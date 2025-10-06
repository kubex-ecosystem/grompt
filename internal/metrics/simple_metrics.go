// Package metrics provides simple metrics collection for the Grompt Gateway
package metrics

import (
	"log"
	"sync"
	"time"
)

// SimpleMetrics provides basic metrics collection without external dependencies
type SimpleMetrics struct {
	// Counters
	requestsTotal      map[string]int64
	tokensGenerated    map[string]int64
	generationErrors   map[string]int64
	proxyRequests      map[string]int64

	// Cost tracking
	estimatedCostUSD   map[string]float64

	// Performance tracking
	latencies          map[string][]time.Duration

	mutex sync.RWMutex
}

// NewSimpleMetrics creates a new metrics collector
func NewSimpleMetrics() *SimpleMetrics {
	return &SimpleMetrics{
		requestsTotal:    make(map[string]int64),
		tokensGenerated:  make(map[string]int64),
		generationErrors: make(map[string]int64),
		proxyRequests:    make(map[string]int64),
		estimatedCostUSD: make(map[string]float64),
		latencies:        make(map[string][]time.Duration),
	}
}

// RecordRequest logs request metrics
func (m *SimpleMetrics) RecordRequest(endpoint, method, statusCode string, duration time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	key := endpoint + ":" + method + ":" + statusCode
	m.requestsTotal[key]++

	latencyKey := endpoint + ":" + method
	m.latencies[latencyKey] = append(m.latencies[latencyKey], duration)

	log.Printf("[METRICS] Request: endpoint=%s method=%s status=%s duration=%v",
		endpoint, method, statusCode, duration)
}

// RecordTokensGenerated logs token generation
func (m *SimpleMetrics) RecordTokensGenerated(provider, model, endpoint string, tokens int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	key := provider + ":" + model + ":" + endpoint
	m.tokensGenerated[key] += int64(tokens)

	log.Printf("[METRICS] Tokens generated: provider=%s model=%s endpoint=%s tokens=%d",
		provider, model, endpoint, tokens)
}

// RecordGenerationLatency logs generation timing
func (m *SimpleMetrics) RecordGenerationLatency(provider, model, endpoint string, duration time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	key := provider + ":" + model + ":" + endpoint
	m.latencies[key] = append(m.latencies[key], duration)

	log.Printf("[METRICS] Generation latency: provider=%s model=%s endpoint=%s duration=%v",
		provider, model, endpoint, duration)
}

// RecordGenerationError logs generation errors
func (m *SimpleMetrics) RecordGenerationError(provider, model, errorType string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	key := provider + ":" + model + ":" + errorType
	m.generationErrors[key]++

	log.Printf("[METRICS] Generation error: provider=%s model=%s error_type=%s",
		provider, model, errorType)
}

// RecordEstimatedCost logs cost information
func (m *SimpleMetrics) RecordEstimatedCost(provider, model string, costUSD float64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	key := provider + ":" + model
	m.estimatedCostUSD[key] += costUSD

	log.Printf("[METRICS] Estimated cost: provider=%s model=%s cost_usd=%.6f",
		provider, model, costUSD)
}

// RecordProxyRequest logs proxy metrics
func (m *SimpleMetrics) RecordProxyRequest(targetPath, method, statusCode string, duration time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	key := targetPath + ":" + method + ":" + statusCode
	m.proxyRequests[key]++

	log.Printf("[METRICS] Proxy request: target=%s method=%s status=%s duration=%v",
		targetPath, method, statusCode, duration)
}

// GetSummary returns a summary of current metrics
func (m *SimpleMetrics) GetSummary() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Calculate totals
	totalRequests := int64(0)
	for _, count := range m.requestsTotal {
		totalRequests += count
	}

	totalTokens := int64(0)
	for _, count := range m.tokensGenerated {
		totalTokens += count
	}

	totalCost := float64(0)
	for _, cost := range m.estimatedCostUSD {
		totalCost += cost
	}

	totalErrors := int64(0)
	for _, count := range m.generationErrors {
		totalErrors += count
	}

	return map[string]interface{}{
		"total_requests":      totalRequests,
		"total_tokens":        totalTokens,
		"total_cost_usd":      totalCost,
		"total_errors":        totalErrors,
		"requests_by_endpoint": m.requestsTotal,
		"tokens_by_provider":   m.tokensGenerated,
		"cost_by_provider":     m.estimatedCostUSD,
		"errors_by_type":       m.generationErrors,
	}
}

// Global metrics instance
var (
	globalMetrics *SimpleMetrics
	metricsOnce   sync.Once
)

// GetGlobalMetrics returns the global metrics instance (singleton)
func GetGlobalMetrics() *SimpleMetrics {
	metricsOnce.Do(func() {
		globalMetrics = NewSimpleMetrics()
	})
	return globalMetrics
}

// Convenience functions for global metrics usage

// RecordRequestGlobal records a request using global metrics
func RecordRequestGlobal(endpoint, method, statusCode string, duration time.Duration) {
	GetGlobalMetrics().RecordRequest(endpoint, method, statusCode, duration)
}

// RecordTokensGeneratedGlobal records tokens generated using global metrics
func RecordTokensGeneratedGlobal(provider, model, endpoint string, tokens int) {
	GetGlobalMetrics().RecordTokensGenerated(provider, model, endpoint, tokens)
}

// RecordGenerationLatencyGlobal records generation latency using global metrics
func RecordGenerationLatencyGlobal(provider, model, endpoint string, duration time.Duration) {
	GetGlobalMetrics().RecordGenerationLatency(provider, model, endpoint, duration)
}

// RecordEstimatedCostGlobal records estimated cost using global metrics
func RecordEstimatedCostGlobal(provider, model string, costUSD float64) {
	GetGlobalMetrics().RecordEstimatedCost(provider, model, costUSD)
}

// RecordProxyRequestGlobal records proxy request using global metrics
func RecordProxyRequestGlobal(targetPath, method, statusCode string, duration time.Duration) {
	GetGlobalMetrics().RecordProxyRequest(targetPath, method, statusCode, duration)
}