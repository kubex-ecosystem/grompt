package middleware

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

// HealthStatus represents the health status of a provider
type HealthStatus int

const (
	// HealthUnknown - initial state or unknown health
	HealthUnknown HealthStatus = iota
	// HealthHealthy - provider is responding normally
	HealthHealthy
	// HealthUnhealthy - provider is experiencing issues
	HealthUnhealthy
	// HealthDegraded - provider is working but with reduced performance
	HealthDegraded
)

// String returns string representation of health status
func (hs HealthStatus) String() string {
	switch hs {
	case HealthUnknown:
		return "UNKNOWN"
	case HealthHealthy:
		return "HEALTHY"
	case HealthUnhealthy:
		return "UNHEALTHY"
	case HealthDegraded:
		return "DEGRADED"
	default:
		return "UNKNOWN"
	}
}

// HealthCheck represents a health check result
type HealthCheck struct {
	Provider     string        `json:"provider"`
	Status       HealthStatus  `json:"status"`
	LastCheck    time.Time     `json:"last_check"`
	ResponseTime time.Duration `json:"response_time"`
	ErrorMsg     string        `json:"error_msg,omitempty"`
	Uptime       float64       `json:"uptime_percentage"`
}

// HealthMonitor monitors the health of providers
type HealthMonitor struct {
	checks   map[string]*HealthCheck
	history  map[string][]bool // Recent success/failure history for uptime calculation
	mu       sync.RWMutex
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewHealthMonitor creates a new health monitor
func NewHealthMonitor(checkInterval time.Duration) *HealthMonitor {
	ctx, cancel := context.WithCancel(context.Background())

	hm := &HealthMonitor{
		checks:   make(map[string]*HealthCheck),
		history:  make(map[string][]bool),
		interval: checkInterval,
		ctx:      ctx,
		cancel:   cancel,
	}

	// Start background health checking
	go hm.runHealthChecks()

	return hm
}

// RegisterProvider adds a provider to health monitoring
func (hm *HealthMonitor) RegisterProvider(provider string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	hm.checks[provider] = &HealthCheck{
		Provider:  provider,
		Status:    HealthUnknown,
		LastCheck: time.Now(),
		Uptime:    100.0,
	}
	hm.history[provider] = make([]bool, 0, 100) // Keep last 100 checks

	fmt.Printf("[HealthMonitor] Registered provider: %s\n", provider)
}

// RecordCheck records the result of a health check
func (hm *HealthMonitor) RecordCheck(provider string, success bool, responseTime time.Duration, errorMsg string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	check, exists := hm.checks[provider]
	if !exists {
		return
	}

	// Update health check
	check.LastCheck = time.Now()
	check.ResponseTime = responseTime
	check.ErrorMsg = errorMsg

	// Update status based on response time and success
	if success {
		if responseTime < 1*time.Second {
			check.Status = HealthHealthy
		} else if responseTime < 5*time.Second {
			check.Status = HealthDegraded
		} else {
			check.Status = HealthUnhealthy
		}
	} else {
		check.Status = HealthUnhealthy
	}

	// Update history for uptime calculation
	history := hm.history[provider]
	history = append(history, success)

	// Keep only last 100 checks
	if len(history) > 100 {
		history = history[1:]
	}
	hm.history[provider] = history

	// Calculate uptime percentage
	if len(history) > 0 {
		successCount := 0
		for _, s := range history {
			if s {
				successCount++
			}
		}
		check.Uptime = float64(successCount) / float64(len(history)) * 100
	}
}

// GetHealth returns the current health status of a provider
func (hm *HealthMonitor) GetHealth(provider string) (*HealthCheck, bool) {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	check, exists := hm.checks[provider]
	if !exists {
		return nil, false
	}

	// Return a copy to avoid race conditions
	return &HealthCheck{
		Provider:     check.Provider,
		Status:       check.Status,
		LastCheck:    check.LastCheck,
		ResponseTime: check.ResponseTime,
		ErrorMsg:     check.ErrorMsg,
		Uptime:       check.Uptime,
	}, true
}

// GetAllHealth returns health status for all providers
func (hm *HealthMonitor) GetAllHealth() map[string]*HealthCheck {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	result := make(map[string]*HealthCheck)
	for provider, check := range hm.checks {
		result[provider] = &HealthCheck{
			Provider:     check.Provider,
			Status:       check.Status,
			LastCheck:    check.LastCheck,
			ResponseTime: check.ResponseTime,
			ErrorMsg:     check.ErrorMsg,
			Uptime:       check.Uptime,
		}
	}
	return result
}

// runHealthChecks runs periodic health checks
func (hm *HealthMonitor) runHealthChecks() {
	ticker := time.NewTicker(hm.interval)
	defer ticker.Stop()

	for {
		select {
		case <-hm.ctx.Done():
			return
		case <-ticker.C:
			// Note: In a real implementation, you'd perform actual health checks here
			// For now, we'll let the providers report their own health via RecordCheck
		}
	}
}

// Stop stops the health monitor
func (hm *HealthMonitor) Stop() {
	hm.cancel()
}

// RetryConfig holds configuration for retry logic
type RetryConfig struct {
	MaxRetries int           // Maximum number of retry attempts
	BaseDelay  time.Duration // Base delay between retries
	MaxDelay   time.Duration // Maximum delay between retries
	Multiplier float64       // Exponential backoff multiplier
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries: 3,
		BaseDelay:  100 * time.Millisecond,
		MaxDelay:   5 * time.Second,
		Multiplier: 2.0,
	}
}

// RetryWithBackoff executes a function with exponential backoff retry logic
func RetryWithBackoff(ctx context.Context, config RetryConfig, operation func() error) error {
	var lastErr error

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// Try the operation
		err := operation()
		if err == nil {
			return nil // Success!
		}

		lastErr = err

		// Don't wait after the last attempt
		if attempt == config.MaxRetries {
			break
		}

		// Calculate delay with exponential backoff
		delay := time.Duration(float64(config.BaseDelay) * math.Pow(config.Multiplier, float64(attempt)))
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}

		fmt.Printf("[Retry] Attempt %d failed: %v. Retrying in %v...\n", attempt+1, err, delay)

		// Wait for the delay or context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return fmt.Errorf("operation failed after %d attempts: %w", config.MaxRetries+1, lastErr)
}
