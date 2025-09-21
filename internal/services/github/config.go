// Package github provides GitHub API integration with App authentication and PAT fallback.
package github

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

// Config holds GitHub integration configuration
type Config struct {
	// GitHub App configuration
	AppID          int64  `json:"app_id"`
	AppPrivateKey  string `json:"app_private_key"`
	InstallationID int64  `json:"installation_id"`
	WebhookSecret  string `json:"webhook_secret"`

	// PAT fallback configuration
	PersonalAccessToken string `json:"personal_access_token"`

	// API configuration
	BaseURL         string        `json:"base_url"`
	APIVersion      string        `json:"api_version"`
	UserAgent       string        `json:"user_agent"`
	Timeout         time.Duration `json:"timeout"`
	MaxRetries      int           `json:"max_retries"`
	RetryBackoffMs  int           `json:"retry_backoff_ms"`
	CacheTTLMinutes int           `json:"cache_ttl_minutes"`

	// Rate limiting
	EnableRateLimit bool `json:"enable_rate_limit"`
	RateLimitBurst  int  `json:"rate_limit_burst"`
}

// LoadConfig loads GitHub configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		// Default values
		BaseURL:         "https://api.github.com",
		APIVersion:      "2022-11-28",
		UserAgent:       "GemX-Grompt/1.0.0",
		Timeout:         30 * time.Second,
		MaxRetries:      3,
		RetryBackoffMs:  1000,
		CacheTTLMinutes: 15,
		EnableRateLimit: true,
		RateLimitBurst:  100,
	}

	// GitHub App configuration
	if appIDStr := os.Getenv("GITHUB_APP_ID"); appIDStr != "" {
		appID, err := strconv.ParseInt(appIDStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid GITHUB_APP_ID: %w", err)
		}
		config.AppID = appID
	}

	config.AppPrivateKey = os.Getenv("GITHUB_APP_PRIVATE_KEY")

	if installationIDStr := os.Getenv("GITHUB_INSTALLATION_ID"); installationIDStr != "" {
		installationID, err := strconv.ParseInt(installationIDStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid GITHUB_INSTALLATION_ID: %w", err)
		}
		config.InstallationID = installationID
	}

	config.WebhookSecret = os.Getenv("GITHUB_WEBHOOK_SECRET")

	// PAT fallback
	config.PersonalAccessToken = os.Getenv("GITHUB_TOKEN")

	// API configuration overrides
	if baseURL := os.Getenv("GITHUB_BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}

	if apiVersion := os.Getenv("GITHUB_API_VERSION"); apiVersion != "" {
		config.APIVersion = apiVersion
	}

	if userAgent := os.Getenv("GITHUB_USER_AGENT"); userAgent != "" {
		config.UserAgent = userAgent
	}

	// Timeout configuration
	if timeoutStr := os.Getenv("GITHUB_TIMEOUT_SECONDS"); timeoutStr != "" {
		timeoutSec, err := strconv.Atoi(timeoutStr)
		if err == nil && timeoutSec > 0 {
			config.Timeout = time.Duration(timeoutSec) * time.Second
		}
	}

	// Retry configuration
	if retriesStr := os.Getenv("GITHUB_MAX_RETRIES"); retriesStr != "" {
		retries, err := strconv.Atoi(retriesStr)
		if err == nil && retries >= 0 {
			config.MaxRetries = retries
		}
	}

	if backoffStr := os.Getenv("GITHUB_RETRY_BACKOFF_MS"); backoffStr != "" {
		backoff, err := strconv.Atoi(backoffStr)
		if err == nil && backoff > 0 {
			config.RetryBackoffMs = backoff
		}
	}

	// Cache configuration
	if ttlStr := os.Getenv("GITHUB_CACHE_TTL_MINUTES"); ttlStr != "" {
		ttl, err := strconv.Atoi(ttlStr)
		if err == nil && ttl > 0 {
			config.CacheTTLMinutes = ttl
		}
	}

	// Rate limiting configuration
	if rateLimitStr := os.Getenv("GITHUB_ENABLE_RATE_LIMIT"); rateLimitStr != "" {
		config.EnableRateLimit = rateLimitStr == "true"
	}

	if burstStr := os.Getenv("GITHUB_RATE_LIMIT_BURST"); burstStr != "" {
		burst, err := strconv.Atoi(burstStr)
		if err == nil && burst > 0 {
			config.RateLimitBurst = burst
		}
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate base URL
	if _, err := url.Parse(c.BaseURL); err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}

	// Check authentication configuration
	hasAppAuth := c.AppID > 0 && c.AppPrivateKey != ""
	hasPATAuth := c.PersonalAccessToken != ""

	if !hasAppAuth && !hasPATAuth {
		return fmt.Errorf("no authentication configured: need either GitHub App (GITHUB_APP_ID + GITHUB_APP_PRIVATE_KEY) or Personal Access Token (GITHUB_TOKEN)")
	}

	// Validate GitHub App configuration if provided
	if hasAppAuth {
		if c.AppID <= 0 {
			return fmt.Errorf("invalid GitHub App ID: must be positive")
		}
		if len(c.AppPrivateKey) < 100 {
			return fmt.Errorf("invalid GitHub App private key: too short")
		}
		// Installation ID is optional for App auth (can be determined dynamically)
	}

	// Validate timeout
	if c.Timeout <= 0 {
		return fmt.Errorf("invalid timeout: must be positive")
	}

	// Validate retry configuration
	if c.MaxRetries < 0 {
		return fmt.Errorf("invalid max retries: must be non-negative")
	}
	if c.RetryBackoffMs <= 0 {
		return fmt.Errorf("invalid retry backoff: must be positive")
	}

	// Validate cache TTL
	if c.CacheTTLMinutes <= 0 {
		return fmt.Errorf("invalid cache TTL: must be positive")
	}

	// Validate rate limit configuration
	if c.EnableRateLimit && c.RateLimitBurst <= 0 {
		return fmt.Errorf("invalid rate limit burst: must be positive when rate limiting is enabled")
	}

	return nil
}

// IsAppAuthConfigured returns true if GitHub App authentication is configured
func (c *Config) IsAppAuthConfigured() bool {
	return c.AppID > 0 && c.AppPrivateKey != ""
}

// IsPATAuthConfigured returns true if Personal Access Token authentication is configured
func (c *Config) IsPATAuthConfigured() bool {
	return c.PersonalAccessToken != ""
}

// IsEnterpriseServer returns true if this is a GitHub Enterprise Server instance
func (c *Config) IsEnterpriseServer() bool {
	return c.BaseURL != "https://api.github.com"
}

// GetWebhookSecret returns the webhook secret for signature validation
func (c *Config) GetWebhookSecret() string {
	return c.WebhookSecret
}

// GetCacheTTL returns the cache TTL as a duration
func (c *Config) GetCacheTTL() time.Duration {
	return time.Duration(c.CacheTTLMinutes) * time.Minute
}

// GetRetryBackoff returns the retry backoff as a duration
func (c *Config) GetRetryBackoff() time.Duration {
	return time.Duration(c.RetryBackoffMs) * time.Millisecond
}