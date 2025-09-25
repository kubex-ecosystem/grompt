package middleware

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// CircuitState represents the state of a circuit breaker
type CircuitState int

const (
	// CircuitClosed - normal operation, requests are allowed
	CircuitClosed CircuitState = iota
	// CircuitOpen - circuit is open, requests are blocked
	CircuitOpen
	// CircuitHalfOpen - testing if service has recovered
	CircuitHalfOpen
)

// String returns the string representation of circuit state
func (cs CircuitState) String() string {
	switch cs {
	case CircuitClosed:
		return "CLOSED"
	case CircuitOpen:
		return "OPEN"
	case CircuitHalfOpen:
		return "HALF-OPEN"
	default:
		return "UNKNOWN"
	}
}

// CircuitBreakerConfig holds configuration for a circuit breaker
type CircuitBreakerConfig struct {
	MaxFailures      int           // Number of failures before opening
	ResetTimeout     time.Duration // Time before trying half-open
	SuccessThreshold int           // Successes needed to close from half-open
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config       CircuitBreakerConfig
	state        CircuitState
	failures     int
	successes    int
	lastFailTime time.Time
	mu           sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  CircuitClosed,
	}
}

// Allow checks if a request should be allowed through the circuit
func (cb *CircuitBreaker) Allow() error {
	cb.mu.RLock()
	state := cb.state
	lastFailTime := cb.lastFailTime
	cb.mu.RUnlock()

	switch state {
	case CircuitClosed:
		return nil // Always allow when closed

	case CircuitOpen:
		// Check if enough time has passed to try half-open
		if time.Since(lastFailTime) >= cb.config.ResetTimeout {
			cb.mu.Lock()
			// Double-check state hasn't changed
			if cb.state == CircuitOpen {
				cb.state = CircuitHalfOpen
				cb.successes = 0
				fmt.Printf("[CircuitBreaker] Moving to HALF-OPEN state\n")
			}
			cb.mu.Unlock()
			return nil
		}
		return errors.New("circuit breaker is OPEN")

	case CircuitHalfOpen:
		return nil // Allow limited requests in half-open state

	default:
		return errors.New("unknown circuit breaker state")
	}
}

// RecordSuccess records a successful operation
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case CircuitClosed:
		// Reset failure count on success
		cb.failures = 0

	case CircuitHalfOpen:
		cb.successes++
		if cb.successes >= cb.config.SuccessThreshold {
			cb.state = CircuitClosed
			cb.failures = 0
			cb.successes = 0
			fmt.Printf("[CircuitBreaker] Moving to CLOSED state after %d successes\n",
				cb.config.SuccessThreshold)
		}
	}
}

// RecordFailure records a failed operation
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailTime = time.Now()

	switch cb.state {
	case CircuitClosed:
		if cb.failures >= cb.config.MaxFailures {
			cb.state = CircuitOpen
			fmt.Printf("[CircuitBreaker] Moving to OPEN state after %d failures\n", cb.failures)
		}

	case CircuitHalfOpen:
		// Any failure in half-open state immediately opens the circuit
		cb.state = CircuitOpen
		fmt.Printf("[CircuitBreaker] Moving to OPEN state from HALF-OPEN after failure\n")
	}
}

// GetState returns the current state and metrics
func (cb *CircuitBreaker) GetState() (CircuitState, int, int) {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state, cb.failures, cb.successes
}

// CircuitBreakerManager manages circuit breakers for multiple providers
type CircuitBreakerManager struct {
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
}

// NewCircuitBreakerManager creates a new circuit breaker manager
func NewCircuitBreakerManager() *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// SetCircuitBreaker configures a circuit breaker for a provider
func (cbm *CircuitBreakerManager) SetCircuitBreaker(provider string, config CircuitBreakerConfig) {
	cbm.mu.Lock()
	defer cbm.mu.Unlock()

	cbm.breakers[provider] = NewCircuitBreaker(config)
	fmt.Printf("[CircuitBreaker] Configured %s: %d max failures, %v reset timeout\n",
		provider, config.MaxFailures, config.ResetTimeout)
}

// Allow checks if a request to the provider should be allowed
func (cbm *CircuitBreakerManager) Allow(provider string) error {
	cbm.mu.RLock()
	breaker, exists := cbm.breakers[provider]
	cbm.mu.RUnlock()

	if !exists {
		// No circuit breaker configured, allow by default
		return nil
	}

	return breaker.Allow()
}

// RecordSuccess records a successful operation for a provider
func (cbm *CircuitBreakerManager) RecordSuccess(provider string) {
	cbm.mu.RLock()
	breaker, exists := cbm.breakers[provider]
	cbm.mu.RUnlock()

	if exists {
		breaker.RecordSuccess()
	}
}

// RecordFailure records a failed operation for a provider
func (cbm *CircuitBreakerManager) RecordFailure(provider string) {
	cbm.mu.RLock()
	breaker, exists := cbm.breakers[provider]
	cbm.mu.RUnlock()

	if exists {
		breaker.RecordFailure()
	}
}

// GetStatus returns the circuit breaker status for a provider
func (cbm *CircuitBreakerManager) GetStatus(provider string) (CircuitState, int, int, bool) {
	cbm.mu.RLock()
	breaker, exists := cbm.breakers[provider]
	cbm.mu.RUnlock()

	if !exists {
		return CircuitClosed, 0, 0, false
	}

	state, failures, successes := breaker.GetState()
	return state, failures, successes, true
}
