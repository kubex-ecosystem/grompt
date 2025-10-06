// Package middleware provides production-grade middleware for the gateway including rate limiting and circuit breakers.
package middleware

import (
	"fmt"
	"sync"
	"time"
)

// TokenBucket implements the token bucket algorithm for rate limiting
type TokenBucket struct {
	capacity   int        // Maximum number of tokens
	tokens     int        // Current number of tokens
	refillRate int        // Tokens added per second
	lastRefill time.Time  // Last time tokens were added
	mu         sync.Mutex // Thread safety
}

// NewTokenBucket creates a new token bucket rate limiter
func NewTokenBucket(capacity, refillRate int) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity, // Start with full bucket
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request can proceed (consumes 1 token if available)
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Refill tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	// Calculate tokens to add
	tokensToAdd := int(elapsed.Seconds()) * tb.refillRate
	if tokensToAdd > 0 {
		tb.tokens = min(tb.capacity, tb.tokens+tokensToAdd)
		tb.lastRefill = now
	}

	// Check if we have tokens available
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// Tokens returns current number of available tokens
func (tb *TokenBucket) Tokens() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.tokens
}

// RateLimiter manages rate limiting for multiple providers
type RateLimiter struct {
	buckets map[string]*TokenBucket
	mu      sync.RWMutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string]*TokenBucket),
	}
}

// SetLimit configures rate limit for a specific provider
func (rl *RateLimiter) SetLimit(provider string, capacity, refillRate int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.buckets[provider] = NewTokenBucket(capacity, refillRate)
	fmt.Printf("[RateLimit] Configured %s: %d tokens, %d/sec refill\n",
		provider, capacity, refillRate)
}

// Allow checks if a request to the given provider should be allowed
func (rl *RateLimiter) Allow(provider string) bool {
	rl.mu.RLock()
	bucket, exists := rl.buckets[provider]
	rl.mu.RUnlock()

	if !exists {
		// No rate limit configured, allow by default
		return true
	}

	allowed := bucket.Allow()
	if !allowed {
		fmt.Printf("[RateLimit] BLOCKED request to %s - rate limit exceeded\n", provider)
	}

	return allowed
}

// GetStatus returns the current rate limit status for a provider
func (rl *RateLimiter) GetStatus(provider string) (int, int, bool) {
	rl.mu.RLock()
	bucket, exists := rl.buckets[provider]
	rl.mu.RUnlock()

	if !exists {
		return 0, 0, false
	}

	tokens := bucket.Tokens()
	return tokens, bucket.capacity, true
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
