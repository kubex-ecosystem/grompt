package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	config := &Config{
		PersonalAccessToken: "test-token",
		BaseURL:             "https://api.github.com",
		APIVersion:          "2022-11-28",
		UserAgent:           "test-agent",
		Timeout:             30 * time.Second,
		MaxRetries:          3,
		RetryBackoffMs:      1000,
		CacheTTLMinutes:     15,
		EnableRateLimit:     true,
		RateLimitBurst:      100,
	}

	service, err := NewService(config)
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}

	if service == nil {
		t.Fatal("Expected service, got nil")
	}

	if service.client == nil {
		t.Error("Expected client to be initialized")
	}
}

func TestNewServiceFromEnv(t *testing.T) {
	// This test requires environment setup, so it's more of an integration test
	// For now, we'll test the basic functionality

	// Set up minimal environment
	t.Setenv("GITHUB_TOKEN", "test-token")

	service, err := NewServiceFromEnv()
	if err != nil {
		t.Fatalf("NewServiceFromEnv() failed: %v", err)
	}

	if service == nil {
		t.Fatal("Expected service, got nil")
	}
}

func TestServiceGetRepositoryWithMockServer(t *testing.T) {
	// Create a mock GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/test-owner/test-repo" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"id": 12345,
				"name": "test-repo",
				"full_name": "test-owner/test-repo",
				"owner": {
					"id": 1,
					"login": "test-owner",
					"type": "User"
				},
				"private": false,
				"html_url": "https://github.com/test-owner/test-repo",
				"description": "Test repository",
				"language": "Go",
				"default_branch": "main",
				"created_at": "2023-01-01T00:00:00Z",
				"updated_at": "2023-01-02T00:00:00Z",
				"pushed_at": "2023-01-02T00:00:00Z"
			}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	config := &Config{
		PersonalAccessToken: "test-token",
		BaseURL:             server.URL,
		APIVersion:          "2022-11-28",
		UserAgent:           "test-agent",
		Timeout:             30 * time.Second,
		MaxRetries:          3,
		RetryBackoffMs:      1000,
		CacheTTLMinutes:     15,
		EnableRateLimit:     false, // Disable rate limiting for tests
		RateLimitBurst:      100,
	}

	service, err := NewService(config)
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}

	ctx := context.Background()
	repo, err := service.GetRepository(ctx, "test-owner", "test-repo")
	if err != nil {
		t.Fatalf("GetRepository() failed: %v", err)
	}

	if repo == nil {
		t.Fatal("Expected repository, got nil")
	}

	if repo.Name != "test-repo" {
		t.Errorf("Expected repository name 'test-repo', got %s", repo.Name)
	}

	if repo.FullName != "test-owner/test-repo" {
		t.Errorf("Expected full name 'test-owner/test-repo', got %s", repo.FullName)
	}

	if repo.Owner.Login != "test-owner" {
		t.Errorf("Expected owner 'test-owner', got %s", repo.Owner.Login)
	}
}

func TestServiceGetPullRequestsWithMockServer(t *testing.T) {
	// Create a mock GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/test-owner/test-repo/pulls" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{
					"number": 1,
					"title": "Test PR",
					"state": "open",
					"created_at": "2023-01-01T00:00:00Z",
					"updated_at": "2023-01-02T00:00:00Z",
					"merged_at": null,
					"closed_at": null,
					"commits": 3,
					"additions": 100,
					"deletions": 50,
					"changed_files": 5,
					"review_comments": 0
				}
			]`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	config := &Config{
		PersonalAccessToken: "test-token",
		BaseURL:             server.URL,
		APIVersion:          "2022-11-28",
		UserAgent:           "test-agent",
		Timeout:             30 * time.Second,
		MaxRetries:          3,
		RetryBackoffMs:      1000,
		CacheTTLMinutes:     15,
		EnableRateLimit:     false, // Disable rate limiting for tests
		RateLimitBurst:      100,
	}

	service, err := NewService(config)
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}

	ctx := context.Background()
	since := time.Now().AddDate(0, 0, -30) // 30 days ago

	prs, err := service.GetPullRequests(ctx, "test-owner", "test-repo", since)
	if err != nil {
		t.Fatalf("GetPullRequests() failed: %v", err)
	}

	if len(prs) != 1 {
		t.Errorf("Expected 1 pull request, got %d", len(prs))
	}

	if len(prs) > 0 {
		pr := prs[0]
		if pr.Number != 1 {
			t.Errorf("Expected PR number 1, got %d", pr.Number)
		}
		if pr.Title != "Test PR" {
			t.Errorf("Expected PR title 'Test PR', got %s", pr.Title)
		}
		if pr.State != "open" {
			t.Errorf("Expected PR state 'open', got %s", pr.State)
		}
	}
}

func TestServiceHelperMethods(t *testing.T) {
	config := &Config{
		PersonalAccessToken: "test-token",
		InstallationID:      12345,
		BaseURL:             "https://api.github.com",
		APIVersion:          "2022-11-28",
		UserAgent:           "test-agent",
		Timeout:             30 * time.Second,
		MaxRetries:          3,
		RetryBackoffMs:      1000,
		CacheTTLMinutes:     15,
		EnableRateLimit:     true,
		RateLimitBurst:      100,
	}

	service, err := NewService(config)
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}

	// Test GetInstallationID
	if service.GetInstallationID() != 12345 {
		t.Errorf("Expected installation ID 12345, got %d", service.GetInstallationID())
	}

	// Test SetInstallationID
	service.SetInstallationID(67890)
	if service.GetInstallationID() != 67890 {
		t.Errorf("Expected installation ID 67890, got %d", service.GetInstallationID())
	}

	// Test GetClient
	client := service.GetClient()
	if client == nil {
		t.Error("Expected client, got nil")
	}
}