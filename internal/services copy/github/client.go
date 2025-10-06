// Package github provides a unified GitHub client with cache, retry, and rate limiting.
package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Client provides a unified GitHub API client
type Client struct {
	config       *Config
	httpClient   *http.Client
	auth         *AuthProvider
	rateLimit    *rate.Limiter
	circuitBreaker *CircuitBreaker

	// Cache
	cache      map[string]*CacheEntry
	cacheMutex sync.RWMutex
}

// CacheEntry represents a cached API response
type CacheEntry struct {
	Data      []byte
	ETag      string
	ExpiresAt time.Time
}

// CircuitBreaker implements a simple circuit breaker pattern
type CircuitBreaker struct {
	maxFailures int
	resetTime   time.Duration
	failures    int
	lastFailure time.Time
	state       string // "closed", "open", "half-open"
	mutex       sync.Mutex
}

// NewClient creates a new GitHub client
func NewClient(config *Config) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	auth, err := NewAuthProvider(config, httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth provider: %w", err)
	}

	var rateLimiter *rate.Limiter
	if config.EnableRateLimit {
		// GitHub API rate limit: 5000 requests per hour for authenticated requests
		// We'll be more conservative with burst to avoid hitting secondary limits
		rateLimiter = rate.NewLimiter(rate.Limit(5000.0/3600.0), config.RateLimitBurst)
	}

	circuitBreaker := &CircuitBreaker{
		maxFailures: 5,
		resetTime:   60 * time.Second,
		state:       "closed",
	}

	client := &Client{
		config:         config,
		httpClient:     httpClient,
		auth:           auth,
		rateLimit:      rateLimiter,
		circuitBreaker: circuitBreaker,
		cache:          make(map[string]*CacheEntry),
	}

	return client, nil
}

// Get performs a GET request with caching, retries, and rate limiting
func (c *Client) Get(ctx context.Context, path string, installationID int64) ([]byte, error) {
	return c.request(ctx, "GET", path, nil, installationID)
}

// Post performs a POST request with retries and rate limiting
func (c *Client) Post(ctx context.Context, path string, body []byte, installationID int64) ([]byte, error) {
	return c.request(ctx, "POST", path, body, installationID)
}

// Put performs a PUT request with retries and rate limiting
func (c *Client) Put(ctx context.Context, path string, body []byte, installationID int64) ([]byte, error) {
	return c.request(ctx, "PUT", path, body, installationID)
}

// Patch performs a PATCH request with retries and rate limiting
func (c *Client) Patch(ctx context.Context, path string, body []byte, installationID int64) ([]byte, error) {
	return c.request(ctx, "PATCH", path, body, installationID)
}

// Delete performs a DELETE request with retries and rate limiting
func (c *Client) Delete(ctx context.Context, path string, installationID int64) ([]byte, error) {
	return c.request(ctx, "DELETE", path, nil, installationID)
}

// request performs an HTTP request with all the bells and whistles
func (c *Client) request(ctx context.Context, method, path string, body []byte, installationID int64) ([]byte, error) {
	// Check circuit breaker
	if !c.circuitBreaker.canRequest() {
		return nil, fmt.Errorf("circuit breaker is open")
	}

	// Check rate limit
	if c.rateLimit != nil {
		if err := c.rateLimit.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limit error: %w", err)
		}
	}

	// Check cache for GET requests
	if method == "GET" {
		if cachedData, found := c.getFromCache(path); found {
			return cachedData, nil
		}
	}

	// Perform request with retries
	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := time.Duration(attempt) * c.config.GetRetryBackoff()
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		data, err := c.doRequest(ctx, method, path, body, installationID)
		if err == nil {
			c.circuitBreaker.recordSuccess()
			return data, nil
		}

		lastErr = err

		// Don't retry on certain errors
		if isNonRetryableError(err) {
			break
		}
	}

	c.circuitBreaker.recordFailure()
	return nil, fmt.Errorf("request failed after %d attempts: %w", c.config.MaxRetries+1, lastErr)
}

// doRequest performs a single HTTP request
func (c *Client) doRequest(ctx context.Context, method, path string, body []byte, installationID int64) ([]byte, error) {
	url := strings.TrimSuffix(c.config.BaseURL, "/") + "/" + strings.TrimPrefix(path, "/")

	var bodyReader io.Reader
	if body != nil {
		bodyReader = strings.NewReader(string(body))
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("X-GitHub-Api-Version", c.config.APIVersion)
	req.Header.Set("User-Agent", c.config.UserAgent)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set ETag header for cache validation
	if method == "GET" {
		if etag := c.getCachedETag(path); etag != "" {
			req.Header.Set("If-None-Match", etag)
		}
	}

	// Set authentication
	token, err := c.auth.GetAuthToken(installationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth token: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Perform request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %w", err)
	}
	defer resp.Body.Close()

	// Handle 304 Not Modified
	if resp.StatusCode == http.StatusNotModified {
		if cachedData, found := c.getFromCache(path); found {
			return cachedData, nil
		}
	}

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for API errors
	if resp.StatusCode >= 400 {
		return nil, c.handleAPIError(resp.StatusCode, respBody)
	}

	// Cache successful GET responses
	if method == "GET" && resp.StatusCode == http.StatusOK {
		c.cacheResponse(path, respBody, resp.Header.Get("ETag"))
	}

	return respBody, nil
}

// getFromCache retrieves data from cache if valid
func (c *Client) getFromCache(path string) ([]byte, bool) {
	c.cacheMutex.RLock()
	defer c.cacheMutex.RUnlock()

	entry, exists := c.cache[path]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry.Data, true
}

// getCachedETag retrieves the ETag for a cached entry
func (c *Client) getCachedETag(path string) string {
	c.cacheMutex.RLock()
	defer c.cacheMutex.RUnlock()

	entry, exists := c.cache[path]
	if !exists {
		return ""
	}

	return entry.ETag
}

// cacheResponse stores a response in cache
func (c *Client) cacheResponse(path string, data []byte, etag string) {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	c.cache[path] = &CacheEntry{
		Data:      data,
		ETag:      etag,
		ExpiresAt: time.Now().Add(c.config.GetCacheTTL()),
	}
}

// handleAPIError handles GitHub API errors
func (c *Client) handleAPIError(statusCode int, body []byte) error {
	var apiError struct {
		Message          string `json:"message"`
		DocumentationURL string `json:"documentation_url"`
	}

	if err := json.Unmarshal(body, &apiError); err != nil {
		return fmt.Errorf("GitHub API error %d: %s", statusCode, string(body))
	}

	return fmt.Errorf("GitHub API error %d: %s", statusCode, apiError.Message)
}

// isNonRetryableError checks if an error should not be retried
func isNonRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	// Don't retry on authentication errors
	if strings.Contains(errStr, "401") || strings.Contains(errStr, "403") {
		return true
	}

	// Don't retry on not found errors
	if strings.Contains(errStr, "404") {
		return true
	}

	// Don't retry on client errors (4xx)
	if strings.Contains(errStr, "GitHub API error 4") {
		return true
	}

	return false
}

// canRequest checks if the circuit breaker allows requests
func (cb *CircuitBreaker) canRequest() bool {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()

	switch cb.state {
	case "closed":
		return true
	case "open":
		if now.After(cb.lastFailure.Add(cb.resetTime)) {
			cb.state = "half-open"
			return true
		}
		return false
	case "half-open":
		return true
	default:
		return true
	}
}

// recordSuccess records a successful request
func (cb *CircuitBreaker) recordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.failures = 0
	cb.state = "closed"
}

// recordFailure records a failed request
func (cb *CircuitBreaker) recordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()

	if cb.failures >= cb.maxFailures {
		cb.state = "open"
	}
}

// GetRateLimitInfo returns current rate limit information
func (c *Client) GetRateLimitInfo(ctx context.Context, installationID int64) (*RateLimitInfo, error) {
	data, err := c.Get(ctx, "/rate_limit", installationID)
	if err != nil {
		return nil, err
	}

	var response struct {
		Rate RateLimitInfo `json:"rate"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse rate limit response: %w", err)
	}

	return &response.Rate, nil
}

// RateLimitInfo represents GitHub API rate limit information
type RateLimitInfo struct {
	Limit     int       `json:"limit"`
	Used      int       `json:"used"`
	Remaining int       `json:"remaining"`
	ResetTime time.Time `json:"reset"`
}

// ClearCache clears the HTTP cache
func (c *Client) ClearCache() {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	c.cache = make(map[string]*CacheEntry)
}