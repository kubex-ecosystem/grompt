package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestService_CreatePR(t *testing.T) {
	tests := []struct {
		name           string
		request        CreatePRRequest
		expectedStatus int
		expectError    bool
		errorContains  string
	}{
		{
			name: "valid PR creation",
			request: CreatePRRequest{
				Title: "Test PR",
				Head:  "feature-branch",
				Base:  "main",
				Body:  "Test description",
			},
			expectedStatus: 201,
			expectError:    false,
		},
		{
			name: "missing title",
			request: CreatePRRequest{
				Head: "feature-branch",
				Base: "main",
			},
			expectError:   true,
			errorContains: "title is required",
		},
		{
			name: "missing head branch",
			request: CreatePRRequest{
				Title: "Test PR",
				Base:  "main",
			},
			expectError:   true,
			errorContains: "head branch is required",
		},
		{
			name: "missing base branch",
			request: CreatePRRequest{
				Title: "Test PR",
				Head:  "feature-branch",
			},
			expectError:   true,
			errorContains: "base branch is required",
		},
		{
			name: "API error response",
			request: CreatePRRequest{
				Title: "Test PR",
				Head:  "feature-branch",
				Base:  "main",
			},
			expectedStatus: 422,
			expectError:    true,
			errorContains:  "failed to create PR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.expectedStatus == 422 {
					w.WriteHeader(422)
					json.NewEncoder(w).Encode(map[string]string{"message": "Validation Failed"})
					return
				}

				if r.URL.Path == "/repos/test/repo/pulls" && r.Method == "POST" {
					w.WriteHeader(201)
					response := GitHubPullRequest{
						Number:       123,
						Title:        tt.request.Title,
						State:        "open",
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
						Additions:    10,
						Deletions:    5,
						ChangedFiles: 2,
					}
					json.NewEncoder(w).Encode(response)
					return
				}

				w.WriteHeader(404)
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
				EnableRateLimit:     true,
				RateLimitBurst:      100,
			}

			service, err := NewService(config)
			if err != nil {
				t.Fatalf("Failed to create service: %v", err)
			}

			ctx := context.Background()
			pr, err := service.CreatePR(ctx, "test", "repo", tt.request)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error but got: %v", err)
				return
			}

			if pr == nil {
				t.Error("Expected PR but got nil")
				return
			}

			if pr.Number != 123 {
				t.Errorf("Expected PR number 123, got %d", pr.Number)
			}

			if pr.Title != tt.request.Title {
				t.Errorf("Expected title '%s', got '%s'", tt.request.Title, pr.Title)
			}
		})
	}
}

func TestService_UpdatePR(t *testing.T) {
	tests := []struct {
		name           string
		prNumber       int
		request        UpdatePRRequest
		expectedStatus int
		expectError    bool
		errorContains  string
	}{
		{
			name:     "update title",
			prNumber: 123,
			request: UpdatePRRequest{
				Title: stringPtr("Updated Title"),
			},
			expectedStatus: 200,
			expectError:    false,
		},
		{
			name:     "update state to closed",
			prNumber: 123,
			request: UpdatePRRequest{
				State: stringPtr("closed"),
			},
			expectedStatus: 200,
			expectError:    false,
		},
		{
			name:     "invalid state",
			prNumber: 123,
			request: UpdatePRRequest{
				State: stringPtr("invalid"),
			},
			expectError:   true,
			errorContains: "invalid state",
		},
		{
			name:     "no fields to update",
			prNumber: 123,
			request:  UpdatePRRequest{},
			expectError:   true,
			errorContains: "no fields to update",
		},
		{
			name:     "PR not found",
			prNumber: 999,
			request: UpdatePRRequest{
				Title: stringPtr("Updated Title"),
			},
			expectedStatus: 404,
			expectError:    true,
			errorContains:  "failed to update PR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.expectedStatus == 404 {
					w.WriteHeader(404)
					json.NewEncoder(w).Encode(map[string]string{"message": "Not Found"})
					return
				}

				if strings.Contains(r.URL.Path, "/pulls/123") && r.Method == "PATCH" {
					w.WriteHeader(200)
					response := GitHubPullRequest{
						Number:    123,
						Title:     "Updated Title",
						State:     "open",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}
					if tt.request.State != nil && *tt.request.State == "closed" {
						response.State = "closed"
					}
					json.NewEncoder(w).Encode(response)
					return
				}

				w.WriteHeader(404)
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
				EnableRateLimit:     true,
				RateLimitBurst:      100,
			}

			service, err := NewService(config)
			if err != nil {
				t.Fatalf("Failed to create service: %v", err)
			}

			ctx := context.Background()
			pr, err := service.UpdatePR(ctx, "test", "repo", tt.prNumber, tt.request)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error but got: %v", err)
				return
			}

			if pr == nil {
				t.Error("Expected PR but got nil")
				return
			}

			if pr.Number != 123 {
				t.Errorf("Expected PR number 123, got %d", pr.Number)
			}
		})
	}
}

func TestService_MergePR(t *testing.T) {
	tests := []struct {
		name           string
		prNumber       int
		request        MergePRRequest
		expectedStatus int
		expectError    bool
		errorContains  string
	}{
		{
			name:     "successful merge",
			prNumber: 123,
			request: MergePRRequest{
				CommitTitle:   "Merge feature",
				CommitMessage: "Adding new feature",
				MergeMethod:   "merge",
			},
			expectedStatus: 200,
			expectError:    false,
		},
		{
			name:     "squash merge",
			prNumber: 123,
			request: MergePRRequest{
				MergeMethod: "squash",
			},
			expectedStatus: 200,
			expectError:    false,
		},
		{
			name:     "invalid merge method",
			prNumber: 123,
			request: MergePRRequest{
				MergeMethod: "invalid",
			},
			expectError:   true,
			errorContains: "invalid merge method",
		},
		{
			name:     "merge conflict",
			prNumber: 456,
			request: MergePRRequest{
				MergeMethod: "merge",
			},
			expectedStatus: 409,
			expectError:    true,
			errorContains:  "failed to merge PR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.expectedStatus == 409 {
					w.WriteHeader(409)
					json.NewEncoder(w).Encode(map[string]string{"message": "Merge conflict"})
					return
				}

				if strings.Contains(r.URL.Path, "/merge") && r.Method == "PUT" {
					w.WriteHeader(200)
					response := GitHubMergeResult{
						SHA:    "abc123",
						Merged: true,
						Message: "Pull Request successfully merged",
					}
					json.NewEncoder(w).Encode(response)
					return
				}

				w.WriteHeader(404)
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
				EnableRateLimit:     true,
				RateLimitBurst:      100,
			}

			service, err := NewService(config)
			if err != nil {
				t.Fatalf("Failed to create service: %v", err)
			}

			ctx := context.Background()
			result, err := service.MergePR(ctx, "test", "repo", tt.prNumber, tt.request)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error but got: %v", err)
				return
			}

			if result == nil {
				t.Error("Expected merge result but got nil")
				return
			}

			if !result.Merged {
				t.Error("Expected merge to be successful")
			}

			if result.SHA != "abc123" {
				t.Errorf("Expected SHA 'abc123', got '%s'", result.SHA)
			}
		})
	}
}

func TestService_ValidatePR(t *testing.T) {
	tests := []struct {
		name          string
		prNumber      int
		policy        PRValidationPolicy
		prData        GitHubPullRequest
		expectError   bool
		errorContains string
	}{
		{
			name:     "valid PR - all checks pass",
			prNumber: 123,
			policy: PRValidationPolicy{
				MaxPRSize:       100,
				MaxFilesChanged: 10,
			},
			prData: GitHubPullRequest{
				Number:       123,
				Additions:    30,
				Deletions:    20,
				ChangedFiles: 5,
			},
			expectError: false,
		},
		{
			name:     "PR too large",
			prNumber: 123,
			policy: PRValidationPolicy{
				MaxPRSize: 50,
			},
			prData: GitHubPullRequest{
				Number:    123,
				Additions: 60,
				Deletions: 20,
			},
			expectError:   true,
			errorContains: "PR size (80 lines) exceeds maximum allowed (50 lines)",
		},
		{
			name:     "too many files changed",
			prNumber: 123,
			policy: PRValidationPolicy{
				MaxFilesChanged: 5,
			},
			prData: GitHubPullRequest{
				Number:       123,
				ChangedFiles: 10,
			},
			expectError:   true,
			errorContains: "PR changes 10 files, exceeds maximum allowed (5 files)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.Contains(r.URL.Path, fmt.Sprintf("/pulls/%d", tt.prNumber)) && r.Method == "GET" {
					w.WriteHeader(200)
					json.NewEncoder(w).Encode(tt.prData)
					return
				}
				w.WriteHeader(404)
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
				EnableRateLimit:     true,
				RateLimitBurst:      100,
			}

			service, err := NewService(config)
			if err != nil {
				t.Fatalf("Failed to create service: %v", err)
			}

			ctx := context.Background()
			err = service.ValidatePR(ctx, "test", "repo", tt.prNumber, tt.policy)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestService_GetPR(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/test/repo/pulls/123" && r.Method == "GET" {
			w.WriteHeader(200)
			response := GitHubPullRequest{
				Number:    123,
				Title:     "Test PR",
				State:     "open",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(404)
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
		EnableRateLimit:     true,
		RateLimitBurst:      100,
	}

	service, err := NewService(config)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	ctx := context.Background()
	pr, err := service.GetPR(ctx, "test", "repo", 123)

	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
		return
	}

	if pr == nil {
		t.Error("Expected PR but got nil")
		return
	}

	if pr.Number != 123 {
		t.Errorf("Expected PR number 123, got %d", pr.Number)
	}

	if pr.Title != "Test PR" {
		t.Errorf("Expected title 'Test PR', got '%s'", pr.Title)
	}
}

// Helper function for string pointers
func stringPtr(s string) *string {
	return &s
}