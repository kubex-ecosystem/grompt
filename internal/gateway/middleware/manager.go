package middleware

import (
	"context"
	"fmt"
	"time"
)

// ProductionConfig holds all production middleware configuration
type ProductionConfig struct {
	RateLimit struct {
		Enabled bool `yaml:"enabled"`
		Default struct {
			Capacity   int `yaml:"capacity"`    // requests per bucket
			RefillRate int `yaml:"refill_rate"` // tokens per second
		} `yaml:"default"`
		PerProvider map[string]struct {
			Capacity   int `yaml:"capacity"`
			RefillRate int `yaml:"refill_rate"`
		} `yaml:"per_provider"`
	} `yaml:"rate_limit"`

	CircuitBreaker struct {
		Enabled bool `yaml:"enabled"`
		Default struct {
			MaxFailures      int `yaml:"max_failures"`
			ResetTimeoutSec  int `yaml:"reset_timeout_sec"`
			SuccessThreshold int `yaml:"success_threshold"`
		} `yaml:"default"`
		PerProvider map[string]struct {
			MaxFailures      int `yaml:"max_failures"`
			ResetTimeoutSec  int `yaml:"reset_timeout_sec"`
			SuccessThreshold int `yaml:"success_threshold"`
		} `yaml:"per_provider"`
	} `yaml:"circuit_breaker"`

	HealthCheck struct {
		Enabled     bool `yaml:"enabled"`
		IntervalSec int  `yaml:"interval_sec"`
		TimeoutSec  int  `yaml:"timeout_sec"`
	} `yaml:"health_check"`

	Retry struct {
		Enabled     bool    `yaml:"enabled"`
		MaxRetries  int     `yaml:"max_retries"`
		BaseDelayMs int     `yaml:"base_delay_ms"`
		MaxDelayMs  int     `yaml:"max_delay_ms"`
		Multiplier  float64 `yaml:"multiplier"`
	} `yaml:"retry"`
}

// DefaultProductionConfig returns a sensible default configuration
func DefaultProductionConfig() ProductionConfig {
	config := ProductionConfig{}

	// Rate limiting defaults
	config.RateLimit.Enabled = true
	config.RateLimit.Default.Capacity = 100
	config.RateLimit.Default.RefillRate = 10

	// Circuit breaker defaults
	config.CircuitBreaker.Enabled = true
	config.CircuitBreaker.Default.MaxFailures = 5
	config.CircuitBreaker.Default.ResetTimeoutSec = 60
	config.CircuitBreaker.Default.SuccessThreshold = 3

	// Health check defaults
	config.HealthCheck.Enabled = true
	config.HealthCheck.IntervalSec = 30
	config.HealthCheck.TimeoutSec = 10

	// Retry defaults
	config.Retry.Enabled = true
	config.Retry.MaxRetries = 3
	config.Retry.BaseDelayMs = 100
	config.Retry.MaxDelayMs = 5000
	config.Retry.Multiplier = 2.0

	return config
}

// ProductionMiddleware wraps all production middleware functionality
type ProductionMiddleware struct {
	config         ProductionConfig
	rateLimiter    *RateLimiter
	circuitBreaker *CircuitBreakerManager
	healthMonitor  *HealthMonitor
	retryConfig    RetryConfig
}

// NewProductionMiddleware creates a new production middleware manager
func NewProductionMiddleware(config ProductionConfig) *ProductionMiddleware {
	pm := &ProductionMiddleware{
		config: config,
	}

	// Initialize rate limiter
	if config.RateLimit.Enabled {
		pm.rateLimiter = NewRateLimiter()
	}

	// Initialize circuit breaker
	if config.CircuitBreaker.Enabled {
		pm.circuitBreaker = NewCircuitBreakerManager()
	}

	// Initialize health monitor
	if config.HealthCheck.Enabled {
		interval := time.Duration(config.HealthCheck.IntervalSec) * time.Second
		pm.healthMonitor = NewHealthMonitor(interval)
	}

	// Setup retry config
	if config.Retry.Enabled {
		pm.retryConfig = RetryConfig{
			MaxRetries: config.Retry.MaxRetries,
			BaseDelay:  time.Duration(config.Retry.BaseDelayMs) * time.Millisecond,
			MaxDelay:   time.Duration(config.Retry.MaxDelayMs) * time.Millisecond,
			Multiplier: config.Retry.Multiplier,
		}
	}

	fmt.Println("[ProductionMiddleware] Initialized with enterprise features:")
	if config.RateLimit.Enabled {
		fmt.Printf("  ✅ Rate Limiting: %d capacity, %d/sec refill\n",
			config.RateLimit.Default.Capacity, config.RateLimit.Default.RefillRate)
	}
	if config.CircuitBreaker.Enabled {
		fmt.Printf("  ✅ Circuit Breaker: %d max failures, %ds reset timeout\n",
			config.CircuitBreaker.Default.MaxFailures, config.CircuitBreaker.Default.ResetTimeoutSec)
	}
	if config.HealthCheck.Enabled {
		fmt.Printf("  ✅ Health Checks: every %ds\n", config.HealthCheck.IntervalSec)
	}
	if config.Retry.Enabled {
		fmt.Printf("  ✅ Retry Logic: %d max retries with exponential backoff\n", config.Retry.MaxRetries)
	}

	return pm
}

// RegisterProvider registers a provider with all middleware components
func (pm *ProductionMiddleware) RegisterProvider(provider string) {
	// Set up rate limiting
	if pm.rateLimiter != nil {
		capacity := pm.config.RateLimit.Default.Capacity
		refillRate := pm.config.RateLimit.Default.RefillRate

		// Check for provider-specific configuration
		if providerConfig, exists := pm.config.RateLimit.PerProvider[provider]; exists {
			capacity = providerConfig.Capacity
			refillRate = providerConfig.RefillRate
		}

		pm.rateLimiter.SetLimit(provider, capacity, refillRate)
	}

	// Set up circuit breaker
	if pm.circuitBreaker != nil {
		maxFailures := pm.config.CircuitBreaker.Default.MaxFailures
		resetTimeout := time.Duration(pm.config.CircuitBreaker.Default.ResetTimeoutSec) * time.Second
		successThreshold := pm.config.CircuitBreaker.Default.SuccessThreshold

		// Check for provider-specific configuration
		if providerConfig, exists := pm.config.CircuitBreaker.PerProvider[provider]; exists {
			maxFailures = providerConfig.MaxFailures
			resetTimeout = time.Duration(providerConfig.ResetTimeoutSec) * time.Second
			successThreshold = providerConfig.SuccessThreshold
		}

		pm.circuitBreaker.SetCircuitBreaker(provider, CircuitBreakerConfig{
			MaxFailures:      maxFailures,
			ResetTimeout:     resetTimeout,
			SuccessThreshold: successThreshold,
		})
	}

	// Register with health monitor
	if pm.healthMonitor != nil {
		pm.healthMonitor.RegisterProvider(provider)
	}
}

// WrapProvider wraps a provider call with all production middleware
func (pm *ProductionMiddleware) WrapProvider(provider string, operation func() error) error {
	startTime := time.Now()

	// 1. Check rate limit
	if pm.rateLimiter != nil {
		if !pm.rateLimiter.Allow(provider) {
			return fmt.Errorf("rate limit exceeded for provider %s", provider)
		}
	}

	// 2. Check circuit breaker
	if pm.circuitBreaker != nil {
		if err := pm.circuitBreaker.Allow(provider); err != nil {
			return fmt.Errorf("circuit breaker blocked request to %s: %w", provider, err)
		}
	}

	// 3. Execute with retry logic
	var err error
	if pm.config.Retry.Enabled {
		ctx := context.Background()
		err = RetryWithBackoff(ctx, pm.retryConfig, operation)
	} else {
		err = operation()
	}

	// 4. Record results
	responseTime := time.Since(startTime)
	success := err == nil

	// Record circuit breaker result
	if pm.circuitBreaker != nil {
		if success {
			pm.circuitBreaker.RecordSuccess(provider)
		} else {
			pm.circuitBreaker.RecordFailure(provider)
		}
	}

	// Record health check result
	if pm.healthMonitor != nil {
		errorMsg := ""
		if err != nil {
			errorMsg = err.Error()
		}
		pm.healthMonitor.RecordCheck(provider, success, responseTime, errorMsg)
	}

	return err
}

// GetStatus returns comprehensive status for all middleware components
func (pm *ProductionMiddleware) GetStatus() map[string]interface{} {
	status := make(map[string]interface{})

	// Rate limit status
	if pm.rateLimiter != nil {
		rateLimitStatus := make(map[string]interface{})
		// Note: You'd need to implement a way to get all provider names
		// For now, we'll just indicate that rate limiting is enabled
		rateLimitStatus["enabled"] = true
		status["rate_limit"] = rateLimitStatus
	}

	// Circuit breaker status
	if pm.circuitBreaker != nil {
		circuitBreakerStatus := make(map[string]interface{})
		circuitBreakerStatus["enabled"] = true
		status["circuit_breaker"] = circuitBreakerStatus
	}

	// Health check status
	if pm.healthMonitor != nil {
		healthStatus := pm.healthMonitor.GetAllHealth()
		status["health_checks"] = healthStatus
	}

	return status
}

// GetHealthMonitor returns the health monitor instance
func (pm *ProductionMiddleware) GetHealthMonitor() *HealthMonitor {
	return pm.healthMonitor
}

// Stop gracefully stops all middleware components
func (pm *ProductionMiddleware) Stop() {
	if pm.healthMonitor != nil {
		pm.healthMonitor.Stop()
	}
	fmt.Println("[ProductionMiddleware] Stopped all components")
}
