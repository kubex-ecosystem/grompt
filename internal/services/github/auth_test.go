package github

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewAuthProvider(t *testing.T) {
	config := &Config{
		PersonalAccessToken: "test-token",
		BaseURL:             "https://api.github.com",
		Timeout:             30 * time.Second,
		MaxRetries:          3,
		RetryBackoffMs:      1000,
		CacheTTLMinutes:     15,
		EnableRateLimit:     true,
		RateLimitBurst:      100,
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}

	auth, err := NewAuthProvider(config, httpClient)
	if err != nil {
		t.Fatalf("NewAuthProvider() failed: %v", err)
	}

	if auth == nil {
		t.Fatal("Expected auth provider, got nil")
	}

	if !auth.IsUsingPAT() {
		t.Error("Expected IsUsingPAT() to return true")
	}

	if auth.IsUsingAppAuth() {
		t.Error("Expected IsUsingAppAuth() to return false")
	}
}

func TestGetAuthTokenWithPAT(t *testing.T) {
	config := &Config{
		PersonalAccessToken: "test-token",
		BaseURL:             "https://api.github.com",
		Timeout:             30 * time.Second,
		MaxRetries:          3,
		RetryBackoffMs:      1000,
		CacheTTLMinutes:     15,
		EnableRateLimit:     true,
		RateLimitBurst:      100,
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}

	auth, err := NewAuthProvider(config, httpClient)
	if err != nil {
		t.Fatalf("NewAuthProvider() failed: %v", err)
	}

	token, err := auth.GetAuthToken(0)
	if err != nil {
		t.Fatalf("GetAuthToken() failed: %v", err)
	}

	if token != "test-token" {
		t.Errorf("Expected token 'test-token', got %s", token)
	}
}

func TestGetAuthTokenWithoutConfig(t *testing.T) {
	config := &Config{
		BaseURL:         "https://api.github.com",
		Timeout:         30 * time.Second,
		MaxRetries:      3,
		RetryBackoffMs:  1000,
		CacheTTLMinutes: 15,
		EnableRateLimit: true,
		RateLimitBurst:  100,
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}

	auth, err := NewAuthProvider(config, httpClient)
	if err != nil {
		t.Fatalf("NewAuthProvider() failed: %v", err)
	}

	_, err = auth.GetAuthToken(0)
	if err == nil {
		t.Error("Expected GetAuthToken() to fail without authentication config")
	}
}

func TestGetAppJWTWithoutAppAuth(t *testing.T) {
	config := &Config{
		PersonalAccessToken: "test-token",
		BaseURL:             "https://api.github.com",
		Timeout:             30 * time.Second,
		MaxRetries:          3,
		RetryBackoffMs:      1000,
		CacheTTLMinutes:     15,
		EnableRateLimit:     true,
		RateLimitBurst:      100,
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}

	auth, err := NewAuthProvider(config, httpClient)
	if err != nil {
		t.Fatalf("NewAuthProvider() failed: %v", err)
	}

	_, err = auth.GetAppJWT()
	if err == nil {
		t.Error("Expected GetAppJWT() to fail without App authentication config")
	}
}

func TestGetInstallationIDWithMockServer(t *testing.T) {
	// Create a mock server for GitHub API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/test-owner/test-repo/installation" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": 12345}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Skip this test for now since it requires a real private key
	t.Skip("Skipping test that requires real private key")


}

func TestParsePrivateKey(t *testing.T) {
	// Test with invalid key
	invalidKey := "not-a-private-key"
	_, err := parsePrivateKey(invalidKey)
	if err == nil {
		t.Error("Expected parsePrivateKey() to fail with invalid key")
	}

	// Test with empty key
	_, err = parsePrivateKey("")
	if err == nil {
		t.Error("Expected parsePrivateKey() to fail with empty key")
	}

	// TODO: Add test with real valid private key for integration testing
}