package github

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Test with minimal configuration
	os.Setenv("GITHUB_TOKEN", "test-token")
	defer os.Unsetenv("GITHUB_TOKEN")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if config.PersonalAccessToken != "test-token" {
		t.Errorf("Expected PersonalAccessToken to be 'test-token', got %s", config.PersonalAccessToken)
	}

	if config.BaseURL != "https://api.github.com" {
		t.Errorf("Expected default BaseURL to be 'https://api.github.com', got %s", config.BaseURL)
	}

	if config.Timeout != 30*time.Second {
		t.Errorf("Expected default Timeout to be 30s, got %v", config.Timeout)
	}
}

func TestLoadConfigWithAppAuth(t *testing.T) {
	// Test with GitHub App configuration
	os.Setenv("GITHUB_APP_ID", "12345")
	os.Setenv("GITHUB_APP_PRIVATE_KEY", "test-private-key")
	os.Setenv("GITHUB_INSTALLATION_ID", "67890")
	defer func() {
		os.Unsetenv("GITHUB_APP_ID")
		os.Unsetenv("GITHUB_APP_PRIVATE_KEY")
		os.Unsetenv("GITHUB_INSTALLATION_ID")
	}()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if config.AppID != 12345 {
		t.Errorf("Expected AppID to be 12345, got %d", config.AppID)
	}

	if config.AppPrivateKey != "test-private-key" {
		t.Errorf("Expected AppPrivateKey to be 'test-private-key', got %s", config.AppPrivateKey)
	}

	if config.InstallationID != 67890 {
		t.Errorf("Expected InstallationID to be 67890, got %d", config.InstallationID)
	}
}

func TestLoadConfigWithCustomSettings(t *testing.T) {
	// Test with custom configuration
	os.Setenv("GITHUB_TOKEN", "test-token")
	os.Setenv("GITHUB_BASE_URL", "https://github.enterprise.com/api/v3")
	os.Setenv("GITHUB_TIMEOUT_SECONDS", "60")
	os.Setenv("GITHUB_MAX_RETRIES", "5")
	os.Setenv("GITHUB_CACHE_TTL_MINUTES", "30")
	defer func() {
		os.Unsetenv("GITHUB_TOKEN")
		os.Unsetenv("GITHUB_BASE_URL")
		os.Unsetenv("GITHUB_TIMEOUT_SECONDS")
		os.Unsetenv("GITHUB_MAX_RETRIES")
		os.Unsetenv("GITHUB_CACHE_TTL_MINUTES")
	}()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if config.BaseURL != "https://github.enterprise.com/api/v3" {
		t.Errorf("Expected BaseURL to be custom URL, got %s", config.BaseURL)
	}

	if config.Timeout != 60*time.Second {
		t.Errorf("Expected Timeout to be 60s, got %v", config.Timeout)
	}

	if config.MaxRetries != 5 {
		t.Errorf("Expected MaxRetries to be 5, got %d", config.MaxRetries)
	}

	if config.CacheTTLMinutes != 30 {
		t.Errorf("Expected CacheTTLMinutes to be 30, got %d", config.CacheTTLMinutes)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "valid PAT config",
			config: &Config{
				PersonalAccessToken: "test-token",
				BaseURL:             "https://api.github.com",
				Timeout:             30 * time.Second,
				MaxRetries:          3,
				RetryBackoffMs:      1000,
				CacheTTLMinutes:     15,
				RateLimitBurst:      100,
				EnableRateLimit:     true,
			},
			expectError: false,
		},
		{
			name: "valid App config",
			config: &Config{
				AppID:           12345,
				AppPrivateKey:   generateTestPrivateKey(),
				BaseURL:         "https://api.github.com",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
				RetryBackoffMs:  1000,
				CacheTTLMinutes: 15,
				RateLimitBurst:  100,
				EnableRateLimit: true,
			},
			expectError: false,
		},
		{
			name: "no authentication",
			config: &Config{
				BaseURL:         "https://api.github.com",
				Timeout:         30 * time.Second,
				MaxRetries:      3,
				RetryBackoffMs:  1000,
				CacheTTLMinutes: 15,
				RateLimitBurst:  100,
				EnableRateLimit: true,
			},
			expectError: true,
		},
		{
			name: "invalid base URL",
			config: &Config{
				PersonalAccessToken: "test-token",
				BaseURL:             "://invalid-url",
				Timeout:             30 * time.Second,
				MaxRetries:          3,
				RetryBackoffMs:      1000,
				CacheTTLMinutes:     15,
				RateLimitBurst:      100,
				EnableRateLimit:     true,
			},
			expectError: true,
		},
		{
			name: "invalid timeout",
			config: &Config{
				PersonalAccessToken: "test-token",
				BaseURL:             "https://api.github.com",
				Timeout:             0,
				MaxRetries:          3,
				RetryBackoffMs:      1000,
				CacheTTLMinutes:     15,
				RateLimitBurst:      100,
				EnableRateLimit:     true,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError && err == nil {
				t.Errorf("Expected validation error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no validation error but got: %v", err)
			}
		})
	}
}

func TestConfigHelperMethods(t *testing.T) {
	config := &Config{
		AppID:               12345,
		AppPrivateKey:       "test-key",
		PersonalAccessToken: "test-token",
		BaseURL:             "https://github.enterprise.com/api/v3",
		WebhookSecret:       "test-secret",
		CacheTTLMinutes:     30,
		RetryBackoffMs:      2000,
	}

	if !config.IsAppAuthConfigured() {
		t.Error("Expected IsAppAuthConfigured() to return true")
	}

	if !config.IsPATAuthConfigured() {
		t.Error("Expected IsPATAuthConfigured() to return true")
	}

	if !config.IsEnterpriseServer() {
		t.Error("Expected IsEnterpriseServer() to return true")
	}

	if config.GetWebhookSecret() != "test-secret" {
		t.Errorf("Expected webhook secret to be 'test-secret', got %s", config.GetWebhookSecret())
	}

	if config.GetCacheTTL() != 30*time.Minute {
		t.Errorf("Expected cache TTL to be 30 minutes, got %v", config.GetCacheTTL())
	}

	if config.GetRetryBackoff() != 2*time.Second {
		t.Errorf("Expected retry backoff to be 2 seconds, got %v", config.GetRetryBackoff())
	}
}

// generateTestPrivateKey generates a test RSA private key for testing
func generateTestPrivateKey() string {
	return `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4qLxBDHoKI0DsVAo+7Y2Ec/Y9XhTEQY5hqD2fR8yyQhRDDa3
ZCOUxM2Y+QQ9VKnj9XU1Z8LBJ1zYvKvQfHvPQbYJ3YXhPYRjQ9KhLJ8XYbPFqvY
+YLcNYQ9KhLJ8XYbPFqvY+YLcNYQ9KhLJ8XYbPFqvY+YLcNYQ9KhLJ8XYbPFqvY
+YLcNYQ9KhLJ8XYbPFqvY+YLcNYQ9KhLJ8XYbPFqvY+YLcNYQ9KhLJ8XYbPFqvY
+YLcNYQ9KhLJ8XYbPFqvY+YLcNYQ9KhLJ8XYbPFqvY+YLcNYQ9KhLJ8XYbPFqvY
+YLcNYQ9KhLJ8XYbPFqvY+YLcNYQ9KhLJ8XYbPFqvY+YLcNYQ9KhLJ8XYbPFqvY
wIDAQABAoIBADQgI9YLfI1K2OwYjV1RV1Q5Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y
1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1
Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y
1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y
1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y
1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y
ECgYEA+xKOFI1/Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1Y1
-----END RSA PRIVATE KEY-----`
}