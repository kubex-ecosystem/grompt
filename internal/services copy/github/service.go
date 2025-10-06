// Package github provides a service wrapper that implements existing interfaces.
package github

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/metrics"
)

// Service provides GitHub integration using the new client architecture
type Service struct {
	client         *Client
	installationID int64
}

// NewService creates a new GitHub service
func NewService(config *Config) (*Service, error) {
	client, err := NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}

	return &Service{
		client:         client,
		installationID: config.InstallationID,
	}, nil
}

// NewServiceFromEnv creates a new GitHub service using environment configuration
func NewServiceFromEnv() (*Service, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return NewService(config)
}

// GetPullRequests implements the metrics.GitHubClient interface
func (s *Service) GetPullRequests(ctx context.Context, owner, repo string, since time.Time) ([]metrics.PullRequest, error) {
	// Get installation ID for this repository if not set
	installationID := s.installationID
	if installationID == 0 && s.client.auth.IsUsingAppAuth() {
		var err error
		installationID, err = s.client.auth.GetInstallationID(owner, repo)
		if err != nil {
			return nil, fmt.Errorf("failed to get installation ID: %w", err)
		}
	}

	path := fmt.Sprintf("/repos/%s/%s/pulls?state=all&since=%s&sort=updated&direction=desc&per_page=100",
		owner, repo, since.Format(time.RFC3339))

	var allPRs []metrics.PullRequest
	page := 1

	for {
		pagePath := fmt.Sprintf("%s&page=%d", path, page)
		data, err := s.client.Get(ctx, pagePath, installationID)
		if err != nil {
			return nil, fmt.Errorf("failed to get pull requests: %w", err)
		}

		var githubPRs []GitHubPullRequest
		if err := json.Unmarshal(data, &githubPRs); err != nil {
			return nil, fmt.Errorf("failed to parse pull requests: %w", err)
		}

		if len(githubPRs) == 0 {
			break
		}

		for _, gpr := range githubPRs {
			pr := metrics.PullRequest{
				Number:       gpr.Number,
				Title:        gpr.Title,
				State:        gpr.State,
				CreatedAt:    gpr.CreatedAt,
				UpdatedAt:    gpr.UpdatedAt,
				MergedAt:     gpr.MergedAt,
				ClosedAt:     gpr.ClosedAt,
				Commits:      gpr.Commits,
				Additions:    gpr.Additions,
				Deletions:    gpr.Deletions,
				ChangedFiles: gpr.ChangedFiles,
			}

			// Get first review time if there are review comments
			if gpr.ReviewComments > 0 {
				firstReview, err := s.getFirstReviewTime(ctx, owner, repo, gpr.Number, installationID)
				if err == nil && firstReview != nil {
					pr.FirstReviewAt = firstReview
				}
			}

			allPRs = append(allPRs, pr)
		}

		page++
		if len(githubPRs) < 100 {
			break
		}
	}

	return allPRs, nil
}

// GetDeployments implements the metrics.GitHubClient interface
func (s *Service) GetDeployments(ctx context.Context, owner, repo string, since time.Time) ([]metrics.Deployment, error) {
	installationID := s.installationID
	if installationID == 0 && s.client.auth.IsUsingAppAuth() {
		var err error
		installationID, err = s.client.auth.GetInstallationID(owner, repo)
		if err != nil {
			return nil, fmt.Errorf("failed to get installation ID: %w", err)
		}
	}

	path := fmt.Sprintf("/repos/%s/%s/deployments", owner, repo)
	data, err := s.client.Get(ctx, path, installationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployments: %w", err)
	}

	var githubDeployments []GitHubDeployment
	if err := json.Unmarshal(data, &githubDeployments); err != nil {
		return nil, fmt.Errorf("failed to parse deployments: %w", err)
	}

	var deployments []metrics.Deployment
	for _, gd := range githubDeployments {
		if gd.CreatedAt.Before(since) {
			continue
		}

		deployment := metrics.Deployment{
			ID:          gd.ID,
			Environment: gd.Environment,
			State:       gd.State,
			CreatedAt:   gd.CreatedAt,
			UpdatedAt:   gd.UpdatedAt,
			SHA:         gd.SHA,
		}
		deployments = append(deployments, deployment)
	}

	return deployments, nil
}

// GetWorkflowRuns implements the metrics.GitHubClient interface
func (s *Service) GetWorkflowRuns(ctx context.Context, owner, repo string, since time.Time) ([]metrics.WorkflowRun, error) {
	installationID := s.installationID
	if installationID == 0 && s.client.auth.IsUsingAppAuth() {
		var err error
		installationID, err = s.client.auth.GetInstallationID(owner, repo)
		if err != nil {
			return nil, fmt.Errorf("failed to get installation ID: %w", err)
		}
	}

	path := fmt.Sprintf("/repos/%s/%s/actions/runs?created=>>=%s&per_page=100",
		owner, repo, since.Format("2006-01-02"))

	data, err := s.client.Get(ctx, path, installationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow runs: %w", err)
	}

	var response GitHubWorkflowRunsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse workflow runs: %w", err)
	}

	var runs []metrics.WorkflowRun
	for _, gwr := range response.WorkflowRuns {
		run := metrics.WorkflowRun{
			ID:         gwr.ID,
			Name:       gwr.Name,
			Status:     gwr.Status,
			Conclusion: gwr.Conclusion,
			CreatedAt:  gwr.CreatedAt,
			UpdatedAt:  gwr.UpdatedAt,
			SHA:        gwr.HeadSHA,
		}
		runs = append(runs, run)
	}

	return runs, nil
}

// getFirstReviewTime gets the first review time for a PR
func (s *Service) getFirstReviewTime(ctx context.Context, owner, repo string, prNumber int, installationID int64) (*time.Time, error) {
	path := fmt.Sprintf("/repos/%s/%s/pulls/%d/reviews", owner, repo, prNumber)
	data, err := s.client.Get(ctx, path, installationID)
	if err != nil {
		return nil, err
	}

	var reviews []GitHubReview
	if err := json.Unmarshal(data, &reviews); err != nil {
		return nil, err
	}

	if len(reviews) > 0 {
		return &reviews[0].SubmittedAt, nil
	}

	return nil, nil
}

// GetRepository gets repository information
func (s *Service) GetRepository(ctx context.Context, owner, repo string) (*GitHubRepository, error) {
	installationID := s.installationID
	if installationID == 0 && s.client.auth.IsUsingAppAuth() {
		var err error
		installationID, err = s.client.auth.GetInstallationID(owner, repo)
		if err != nil {
			return nil, fmt.Errorf("failed to get installation ID: %w", err)
		}
	}

	path := fmt.Sprintf("/repos/%s/%s", owner, repo)
	data, err := s.client.Get(ctx, path, installationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	var repoInfo GitHubRepository
	if err := json.Unmarshal(data, &repoInfo); err != nil {
		return nil, fmt.Errorf("failed to parse repository: %w", err)
	}

	return &repoInfo, nil
}

// ListRepositories lists repositories accessible to the authenticated user/app
func (s *Service) ListRepositories(ctx context.Context) ([]GitHubRepository, error) {
	var path string
	installationID := s.installationID

	if s.client.auth.IsUsingAppAuth() {
		// For GitHub Apps, list installation repositories
		if installationID == 0 {
			return nil, fmt.Errorf("installation ID required for App authentication")
		}
		path = "/installation/repositories?per_page=100"
	} else {
		// For PAT, list user repositories
		path = "/user/repos?per_page=100&sort=updated"
	}

	var allRepos []GitHubRepository
	page := 1

	for {
		pagePath := fmt.Sprintf("%s&page=%d", path, page)
		data, err := s.client.Get(ctx, pagePath, installationID)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories: %w", err)
		}

		var repos []GitHubRepository
		if s.client.auth.IsUsingAppAuth() {
			// GitHub App response format
			var response struct {
				Repositories []GitHubRepository `json:"repositories"`
			}
			if err := json.Unmarshal(data, &response); err != nil {
				return nil, fmt.Errorf("failed to parse repositories: %w", err)
			}
			repos = response.Repositories
		} else {
			// PAT response format
			if err := json.Unmarshal(data, &repos); err != nil {
				return nil, fmt.Errorf("failed to parse repositories: %w", err)
			}
		}

		if len(repos) == 0 {
			break
		}

		allRepos = append(allRepos, repos...)
		page++

		if len(repos) < 100 {
			break
		}
	}

	return allRepos, nil
}

// GetClient returns the underlying GitHub client for advanced operations
func (s *Service) GetClient() *Client {
	return s.client
}

// GetInstallationID returns the installation ID (if using App auth)
func (s *Service) GetInstallationID() int64 {
	return s.installationID
}

// SetInstallationID sets the installation ID for App auth
func (s *Service) SetInstallationID(id int64) {
	s.installationID = id
}

// PR Operations for Checkpoint 3

// CreatePRRequest represents a request to create a pull request
type CreatePRRequest struct {
	Title               string   `json:"title"`
	Head                string   `json:"head"`
	Base                string   `json:"base"`
	Body                string   `json:"body,omitempty"`
	MaintainerCanModify bool     `json:"maintainer_can_modify,omitempty"`
	Draft               bool     `json:"draft,omitempty"`
	Assignees           []string `json:"assignees,omitempty"`
	Reviewers           []string `json:"reviewers,omitempty"`
	TeamReviewers       []string `json:"team_reviewers,omitempty"`
	Labels              []string `json:"labels,omitempty"`
}

// UpdatePRRequest represents a request to update a pull request
type UpdatePRRequest struct {
	Title     *string  `json:"title,omitempty"`
	Body      *string  `json:"body,omitempty"`
	State     *string  `json:"state,omitempty"`
	Base      *string  `json:"base,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
	Labels    []string `json:"labels,omitempty"`
}

// MergePRRequest represents a request to merge a pull request
type MergePRRequest struct {
	CommitTitle   string `json:"commit_title,omitempty"`
	CommitMessage string `json:"commit_message,omitempty"`
	SHA           string `json:"sha,omitempty"`
	MergeMethod   string `json:"merge_method,omitempty"` // merge, squash, rebase
}

// PRValidationPolicy represents validation rules for pull requests
type PRValidationPolicy struct {
	RequiredStatusChecks []string `json:"required_status_checks"`
	RequiredReviewers    int      `json:"required_reviewers"`
	MaxPRSize            int      `json:"max_pr_size"` // Max lines changed
	MaxFilesChanged      int      `json:"max_files_changed"`
	RequiredLabels       []string `json:"required_labels"`
	BlockedLabels        []string `json:"blocked_labels"`
}

// CreatePR creates a new pull request
func (s *Service) CreatePR(ctx context.Context, owner, repo string, request CreatePRRequest) (*GitHubPullRequest, error) {
	installationID := s.installationID
	if installationID == 0 && s.client.auth.IsUsingAppAuth() {
		var err error
		installationID, err = s.client.auth.GetInstallationID(owner, repo)
		if err != nil {
			return nil, fmt.Errorf("failed to get installation ID: %w", err)
		}
	}

	// Validate required fields
	if request.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if request.Head == "" {
		return nil, fmt.Errorf("head branch is required")
	}
	if request.Base == "" {
		return nil, fmt.Errorf("base branch is required")
	}

	// Create PR payload
	payload := map[string]interface{}{
		"title":                 request.Title,
		"head":                  request.Head,
		"base":                  request.Base,
		"maintainer_can_modify": request.MaintainerCanModify,
		"draft":                 request.Draft,
	}

	if request.Body != "" {
		payload["body"] = request.Body
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal PR payload: %w", err)
	}

	// Create the PR
	path := fmt.Sprintf("/repos/%s/%s/pulls", owner, repo)
	data, err := s.client.Post(ctx, path, payloadBytes, installationID)
	if err != nil {
		return nil, fmt.Errorf("failed to create PR: %w", err)
	}

	var pr GitHubPullRequest
	if err := json.Unmarshal(data, &pr); err != nil {
		return nil, fmt.Errorf("failed to parse PR response: %w", err)
	}

	// Handle additional operations if needed
	if len(request.Assignees) > 0 || len(request.Reviewers) > 0 || len(request.TeamReviewers) > 0 {
		if err := s.addPRAssigneesAndReviewers(ctx, owner, repo, pr.Number, request, installationID); err != nil {
			// Log warning but don't fail the PR creation
			fmt.Printf("Warning: failed to add assignees/reviewers: %v\n", err)
		}
	}

	if len(request.Labels) > 0 {
		if err := s.addPRLabels(ctx, owner, repo, pr.Number, request.Labels, installationID); err != nil {
			// Log warning but don't fail the PR creation
			fmt.Printf("Warning: failed to add labels: %v\n", err)
		}
	}

	return &pr, nil
}

// UpdatePR updates an existing pull request
func (s *Service) UpdatePR(ctx context.Context, owner, repo string, prNumber int, request UpdatePRRequest) (*GitHubPullRequest, error) {
	installationID := s.installationID
	if installationID == 0 && s.client.auth.IsUsingAppAuth() {
		var err error
		installationID, err = s.client.auth.GetInstallationID(owner, repo)
		if err != nil {
			return nil, fmt.Errorf("failed to get installation ID: %w", err)
		}
	}

	// Build update payload
	payload := make(map[string]interface{})
	if request.Title != nil {
		payload["title"] = *request.Title
	}
	if request.Body != nil {
		payload["body"] = *request.Body
	}
	if request.State != nil {
		if *request.State != "open" && *request.State != "closed" {
			return nil, fmt.Errorf("invalid state: must be 'open' or 'closed'")
		}
		payload["state"] = *request.State
	}
	if request.Base != nil {
		payload["base"] = *request.Base
	}

	if len(payload) == 0 && len(request.Assignees) == 0 && len(request.Labels) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	var pr *GitHubPullRequest
	var err error

	// Update basic PR fields if any
	if len(payload) > 0 {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal update payload: %w", err)
		}

		path := fmt.Sprintf("/repos/%s/%s/pulls/%d", owner, repo, prNumber)
		data, err := s.client.Patch(ctx, path, payloadBytes, installationID)
		if err != nil {
			return nil, fmt.Errorf("failed to update PR: %w", err)
		}

		if err := json.Unmarshal(data, &pr); err != nil {
			return nil, fmt.Errorf("failed to parse PR response: %w", err)
		}
	} else {
		// Get current PR if only updating assignees/labels
		pr, err = s.GetPR(ctx, owner, repo, prNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to get PR: %w", err)
		}
	}

	// Update assignees if specified
	if len(request.Assignees) > 0 {
		if err := s.updatePRAssignees(ctx, owner, repo, prNumber, request.Assignees, installationID); err != nil {
			return nil, fmt.Errorf("failed to update assignees: %w", err)
		}
	}

	// Update labels if specified
	if len(request.Labels) > 0 {
		if err := s.updatePRLabels(ctx, owner, repo, prNumber, request.Labels, installationID); err != nil {
			return nil, fmt.Errorf("failed to update labels: %w", err)
		}
	}

	return pr, nil
}

// MergePR merges a pull request
func (s *Service) MergePR(ctx context.Context, owner, repo string, prNumber int, request MergePRRequest) (*GitHubMergeResult, error) {
	installationID := s.installationID
	if installationID == 0 && s.client.auth.IsUsingAppAuth() {
		var err error
		installationID, err = s.client.auth.GetInstallationID(owner, repo)
		if err != nil {
			return nil, fmt.Errorf("failed to get installation ID: %w", err)
		}
	}

	// Validate merge method
	if request.MergeMethod == "" {
		request.MergeMethod = "merge"
	}
	if request.MergeMethod != "merge" && request.MergeMethod != "squash" && request.MergeMethod != "rebase" {
		return nil, fmt.Errorf("invalid merge method: must be 'merge', 'squash', or 'rebase'")
	}

	// Build merge payload
	payload := map[string]interface{}{
		"merge_method": request.MergeMethod,
	}

	if request.CommitTitle != "" {
		payload["commit_title"] = request.CommitTitle
	}
	if request.CommitMessage != "" {
		payload["commit_message"] = request.CommitMessage
	}
	if request.SHA != "" {
		payload["sha"] = request.SHA
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal merge payload: %w", err)
	}

	// Perform the merge
	path := fmt.Sprintf("/repos/%s/%s/pulls/%d/merge", owner, repo, prNumber)
	data, err := s.client.Put(ctx, path, payloadBytes, installationID)
	if err != nil {
		return nil, fmt.Errorf("failed to merge PR: %w", err)
	}

	var result GitHubMergeResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse merge response: %w", err)
	}

	return &result, nil
}

// GetPR gets a single pull request
func (s *Service) GetPR(ctx context.Context, owner, repo string, prNumber int) (*GitHubPullRequest, error) {
	installationID := s.installationID
	if installationID == 0 && s.client.auth.IsUsingAppAuth() {
		var err error
		installationID, err = s.client.auth.GetInstallationID(owner, repo)
		if err != nil {
			return nil, fmt.Errorf("failed to get installation ID: %w", err)
		}
	}

	path := fmt.Sprintf("/repos/%s/%s/pulls/%d", owner, repo, prNumber)
	data, err := s.client.Get(ctx, path, installationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR: %w", err)
	}

	var pr GitHubPullRequest
	if err := json.Unmarshal(data, &pr); err != nil {
		return nil, fmt.Errorf("failed to parse PR response: %w", err)
	}

	return &pr, nil
}

// Helper methods for PR operations

// addPRAssigneesAndReviewers adds assignees and reviewers to a PR
func (s *Service) addPRAssigneesAndReviewers(ctx context.Context, owner, repo string, prNumber int, request CreatePRRequest, installationID int64) error {
	// Add assignees
	if len(request.Assignees) > 0 {
		if err := s.updatePRAssignees(ctx, owner, repo, prNumber, request.Assignees, installationID); err != nil {
			return fmt.Errorf("failed to add assignees: %w", err)
		}
	}

	// Add reviewers
	if len(request.Reviewers) > 0 || len(request.TeamReviewers) > 0 {
		reviewPayload := map[string]interface{}{}
		if len(request.Reviewers) > 0 {
			reviewPayload["reviewers"] = request.Reviewers
		}
		if len(request.TeamReviewers) > 0 {
			reviewPayload["team_reviewers"] = request.TeamReviewers
		}

		payloadBytes, err := json.Marshal(reviewPayload)
		if err != nil {
			return fmt.Errorf("failed to marshal review payload: %w", err)
		}

		path := fmt.Sprintf("/repos/%s/%s/pulls/%d/requested_reviewers", owner, repo, prNumber)
		_, err = s.client.Post(ctx, path, payloadBytes, installationID)
		if err != nil {
			return fmt.Errorf("failed to add reviewers: %w", err)
		}
	}

	return nil
}

// addPRLabels adds labels to a PR
func (s *Service) addPRLabels(ctx context.Context, owner, repo string, prNumber int, labels []string, installationID int64) error {
	return s.updatePRLabels(ctx, owner, repo, prNumber, labels, installationID)
}

// updatePRAssignees updates the assignees of a PR
func (s *Service) updatePRAssignees(ctx context.Context, owner, repo string, prNumber int, assignees []string, installationID int64) error {
	payload := map[string]interface{}{
		"assignees": assignees,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal assignees payload: %w", err)
	}

	path := fmt.Sprintf("/repos/%s/%s/issues/%d/assignees", owner, repo, prNumber)
	_, err = s.client.Post(ctx, path, payloadBytes, installationID)
	if err != nil {
		return fmt.Errorf("failed to update assignees: %w", err)
	}

	return nil
}

// updatePRLabels updates the labels of a PR
func (s *Service) updatePRLabels(ctx context.Context, owner, repo string, prNumber int, labels []string, installationID int64) error {
	payload := map[string]interface{}{
		"labels": labels,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal labels payload: %w", err)
	}

	path := fmt.Sprintf("/repos/%s/%s/issues/%d/labels", owner, repo, prNumber)
	_, err = s.client.Put(ctx, path, payloadBytes, installationID)
	if err != nil {
		return fmt.Errorf("failed to update labels: %w", err)
	}

	return nil
}

// ValidatePR validates a PR against the given policy
func (s *Service) ValidatePR(ctx context.Context, owner, repo string, prNumber int, policy PRValidationPolicy) error {
	// Get the PR details
	pr, err := s.GetPR(ctx, owner, repo, prNumber)
	if err != nil {
		return fmt.Errorf("failed to get PR for validation: %w", err)
	}

	// Validate PR size
	if policy.MaxPRSize > 0 {
		totalChanges := pr.Additions + pr.Deletions
		if totalChanges > policy.MaxPRSize {
			return fmt.Errorf("PR size (%d lines) exceeds maximum allowed (%d lines)", totalChanges, policy.MaxPRSize)
		}
	}

	// Validate files changed
	if policy.MaxFilesChanged > 0 && pr.ChangedFiles > policy.MaxFilesChanged {
		return fmt.Errorf("PR changes %d files, exceeds maximum allowed (%d files)", pr.ChangedFiles, policy.MaxFilesChanged)
	}

	// Validate status checks
	if len(policy.RequiredStatusChecks) > 0 {
		if err := s.validatePRStatusChecks(ctx, owner, repo, prNumber, policy.RequiredStatusChecks); err != nil {
			return fmt.Errorf("status check validation failed: %w", err)
		}
	}

	// Validate reviewers
	if policy.RequiredReviewers > 0 {
		if err := s.validatePRReviewers(ctx, owner, repo, prNumber, policy.RequiredReviewers); err != nil {
			return fmt.Errorf("reviewer validation failed: %w", err)
		}
	}

	// Validate labels
	if len(policy.RequiredLabels) > 0 || len(policy.BlockedLabels) > 0 {
		if err := s.validatePRLabels(ctx, owner, repo, prNumber, policy.RequiredLabels, policy.BlockedLabels); err != nil {
			return fmt.Errorf("label validation failed: %w", err)
		}
	}

	return nil
}

// validatePRStatusChecks validates that required status checks are passing
func (s *Service) validatePRStatusChecks(ctx context.Context, owner, repo string, prNumber int, requiredChecks []string) error {
	installationID := s.installationID
	if installationID == 0 && s.client.auth.IsUsingAppAuth() {
		var err error
		installationID, err = s.client.auth.GetInstallationID(owner, repo)
		if err != nil {
			return fmt.Errorf("failed to get installation ID: %w", err)
		}
	}

	// Get PR details to get the head SHA
	pr, err := s.GetPR(ctx, owner, repo, prNumber)
	if err != nil {
		return fmt.Errorf("failed to get PR: %w", err)
	}

	// Get commit status
	path := fmt.Sprintf("/repos/%s/%s/commits/%s/status", owner, repo, pr.Head.SHA)
	data, err := s.client.Get(ctx, path, installationID)
	if err != nil {
		return fmt.Errorf("failed to get commit status: %w", err)
	}

	var status GitHubCommitStatus
	if err := json.Unmarshal(data, &status); err != nil {
		return fmt.Errorf("failed to parse commit status: %w", err)
	}

	// Check each required status check
	statusMap := make(map[string]string)
	for _, check := range status.Statuses {
		statusMap[check.Context] = check.State
	}

	for _, required := range requiredChecks {
		state, exists := statusMap[required]
		if !exists {
			return fmt.Errorf("required status check '%s' not found", required)
		}
		if state != "success" {
			return fmt.Errorf("required status check '%s' is not passing (state: %s)", required, state)
		}
	}

	return nil
}

// validatePRReviewers validates that PR has sufficient reviewers
func (s *Service) validatePRReviewers(ctx context.Context, owner, repo string, prNumber int, requiredCount int) error {
	installationID := s.installationID
	if installationID == 0 && s.client.auth.IsUsingAppAuth() {
		var err error
		installationID, err = s.client.auth.GetInstallationID(owner, repo)
		if err != nil {
			return fmt.Errorf("failed to get installation ID: %w", err)
		}
	}

	// Get PR reviews
	path := fmt.Sprintf("/repos/%s/%s/pulls/%d/reviews", owner, repo, prNumber)
	data, err := s.client.Get(ctx, path, installationID)
	if err != nil {
		return fmt.Errorf("failed to get PR reviews: %w", err)
	}

	var reviews []GitHubReview
	if err := json.Unmarshal(data, &reviews); err != nil {
		return fmt.Errorf("failed to parse reviews: %w", err)
	}

	// Count approved reviews
	approvedReviews := 0
	for _, review := range reviews {
		if review.State == "APPROVED" {
			approvedReviews++
		}
	}

	if approvedReviews < requiredCount {
		return fmt.Errorf("PR has %d approved reviews, requires %d", approvedReviews, requiredCount)
	}

	return nil
}

// validatePRLabels validates PR labels against policy
func (s *Service) validatePRLabels(ctx context.Context, owner, repo string, prNumber int, requiredLabels, blockedLabels []string) error {
	installationID := s.installationID
	if installationID == 0 && s.client.auth.IsUsingAppAuth() {
		var err error
		installationID, err = s.client.auth.GetInstallationID(owner, repo)
		if err != nil {
			return fmt.Errorf("failed to get installation ID: %w", err)
		}
	}

	// Get PR labels
	path := fmt.Sprintf("/repos/%s/%s/issues/%d/labels", owner, repo, prNumber)
	data, err := s.client.Get(ctx, path, installationID)
	if err != nil {
		return fmt.Errorf("failed to get PR labels: %w", err)
	}

	var labels []GitHubLabel
	if err := json.Unmarshal(data, &labels); err != nil {
		return fmt.Errorf("failed to parse labels: %w", err)
	}

	// Build label set
	labelSet := make(map[string]bool)
	for _, label := range labels {
		labelSet[label.Name] = true
	}

	// Check required labels
	for _, required := range requiredLabels {
		if !labelSet[required] {
			return fmt.Errorf("required label '%s' is missing", required)
		}
	}

	// Check blocked labels
	for _, blocked := range blockedLabels {
		if labelSet[blocked] {
			return fmt.Errorf("blocked label '%s' is present", blocked)
		}
	}

	return nil
}

// GitHub API response types (reuse from existing integrations.go)

// GitHubPullRequest represents a GitHub pull request
type GitHubPullRequest struct {
	Number         int        `json:"number"`
	Title          string     `json:"title"`
	State          string     `json:"state"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	MergedAt       *time.Time `json:"merged_at"`
	ClosedAt       *time.Time `json:"closed_at"`
	Commits        int        `json:"commits"`
	Additions      int        `json:"additions"`
	Deletions      int        `json:"deletions"`
	ChangedFiles   int        `json:"changed_files"`
	ReviewComments int        `json:"review_comments"`
	Head           struct {
		SHA string `json:"sha"`
	} `json:"head"`
}

// GitHubDeployment represents a GitHub deployment
type GitHubDeployment struct {
	ID          int       `json:"id"`
	Environment string    `json:"environment"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	SHA         string    `json:"sha"`
}

// GitHubWorkflowRun represents a GitHub Actions workflow run
type GitHubWorkflowRun struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	Conclusion string    `json:"conclusion"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	HeadSHA    string    `json:"head_sha"`
}

// GitHubWorkflowRunsResponse represents the response for workflow runs
type GitHubWorkflowRunsResponse struct {
	TotalCount   int                 `json:"total_count"`
	WorkflowRuns []GitHubWorkflowRun `json:"workflow_runs"`
}

// GitHubReview represents a GitHub pull request review
type GitHubReview struct {
	ID          int       `json:"id"`
	State       string    `json:"state"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// GitHubRepository represents a GitHub repository
type GitHubRepository struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	FullName      string    `json:"full_name"`
	Owner         GitHubUser `json:"owner"`
	Private       bool      `json:"private"`
	HTMLURL       string    `json:"html_url"`
	Description   string    `json:"description"`
	Language      string    `json:"language"`
	DefaultBranch string    `json:"default_branch"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	PushedAt      time.Time `json:"pushed_at"`
}

// GitHubUser represents a GitHub user
type GitHubUser struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Type  string `json:"type"`
}

// GitHubMergeResult represents the result of a PR merge operation
type GitHubMergeResult struct {
	SHA     string `json:"sha"`
	Merged  bool   `json:"merged"`
	Message string `json:"message"`
}

// GitHubCommitStatus represents commit status information
type GitHubCommitStatus struct {
	State    string                 `json:"state"`
	Statuses []GitHubStatusCheckItem `json:"statuses"`
}

// GitHubStatusCheckItem represents a single status check
type GitHubStatusCheckItem struct {
	State       string `json:"state"`
	Context     string `json:"context"`
	Description string `json:"description"`
}

// GitHubLabel represents a GitHub label
type GitHubLabel struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}