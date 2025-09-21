// Package repositories implements real data integrations with GitHub, Jira, WakaTime, and Git.
package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/metrics"
	"github.com/kubex-ecosystem/grompt/internal/services/github"
)

// GitHubClient implements real GitHub API integration using the new GitHub service
type GitHubClient struct {
	service *github.Service
}

// NewGitHubClient creates a new GitHub API client using the enhanced service
func NewGitHubClient(token string) *GitHubClient {
	// Create a config for backward compatibility with PAT token
	config := &github.Config{
		PersonalAccessToken: token,
		BaseURL:             "https://api.github.com",
		APIVersion:          "2022-11-28",
		UserAgent:           "GemX-Grompt/1.0.0",
		Timeout:             30 * time.Second,
		MaxRetries:          3,
		RetryBackoffMs:      1000,
		CacheTTLMinutes:     15,
		EnableRateLimit:     true,
		RateLimitBurst:      100,
	}

	service, err := github.NewService(config)
	if err != nil {
		// Fallback to a basic implementation in case of configuration errors
		// This maintains backward compatibility
		return &GitHubClient{service: nil}
	}

	return &GitHubClient{
		service: service,
	}
}

// NewGitHubClientFromService creates a GitHub client from an existing service
func NewGitHubClientFromService(service *github.Service) *GitHubClient {
	return &GitHubClient{service: service}
}

// GetPullRequests fetches pull requests from GitHub API
func (g *GitHubClient) GetPullRequests(ctx context.Context, owner, repo string, since time.Time) ([]metrics.PullRequest, error) {
	if g.service == nil {
		return nil, fmt.Errorf("GitHub service not initialized")
	}

	return g.service.GetPullRequests(ctx, owner, repo, since)
}

// GetDeployments fetches deployments from GitHub API
func (g *GitHubClient) GetDeployments(ctx context.Context, owner, repo string, since time.Time) ([]metrics.Deployment, error) {
	if g.service == nil {
		return nil, fmt.Errorf("GitHub service not initialized")
	}

	return g.service.GetDeployments(ctx, owner, repo, since)
}

// GetWorkflowRuns fetches workflow runs from GitHub Actions API
func (g *GitHubClient) GetWorkflowRuns(ctx context.Context, owner, repo string, since time.Time) ([]metrics.WorkflowRun, error) {
	if g.service == nil {
		return nil, fmt.Errorf("GitHub service not initialized")
	}

	return g.service.GetWorkflowRuns(ctx, owner, repo, since)
}

// The GitHub API response types are now defined in the github service package.
// This maintains backward compatibility while using the enhanced service.

// JiraClient implements real Jira API integration
type JiraClient struct {
	baseURL    string
	username   string
	apiToken   string
	httpClient *http.Client
}

// NewJiraClient creates a new Jira API client
func NewJiraClient(baseURL, username, apiToken string) *JiraClient {
	return &JiraClient{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		username:   username,
		apiToken:   apiToken,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// GetIssues fetches issues from Jira API
func (j *JiraClient) GetIssues(ctx context.Context, project string, since time.Time) ([]metrics.Issue, error) {
	jql := fmt.Sprintf("project = %s AND created >= %s ORDER BY created DESC",
		project, since.Format("2006-01-02"))

	url := fmt.Sprintf("%s/rest/api/3/search?jql=%s&fields=key,issuetype,status,priority,created,updated,resolutiondate",
		j.baseURL, jql)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(j.username, j.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := j.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jira API error: %d", resp.StatusCode)
	}

	var response JiraSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	var issues []metrics.Issue
	for _, ji := range response.Issues {
		issue := metrics.Issue{
			Key:       ji.Key,
			Type:      ji.Fields.IssueType.Name,
			Status:    ji.Fields.Status.Name,
			Priority:  ji.Fields.Priority.Name,
			CreatedAt: ji.Fields.Created,
			UpdatedAt: ji.Fields.Updated,
		}

		if ji.Fields.ResolutionDate != nil {
			issue.ResolvedAt = ji.Fields.ResolutionDate
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

// JiraSearchResponse API response types
type JiraSearchResponse struct {
	Issues []JiraIssue `json:"issues"`
}

type JiraIssue struct {
	Key    string     `json:"key"`
	Fields JiraFields `json:"fields"`
}

type JiraFields struct {
	IssueType      JiraIssueType `json:"issuetype"`
	Status         JiraStatus    `json:"status"`
	Priority       JiraPriority  `json:"priority"`
	Created        time.Time     `json:"created"`
	Updated        time.Time     `json:"updated"`
	ResolutionDate *time.Time    `json:"resolutiondate"`
}

type JiraIssueType struct {
	Name string `json:"name"`
}

type JiraStatus struct {
	Name string `json:"name"`
}

type JiraPriority struct {
	Name string `json:"name"`
}

// WakaTimeClient implements real WakaTime API integration
type WakaTimeClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewWakaTimeClient creates a new WakaTime API client
func NewWakaTimeClient(apiKey string) *WakaTimeClient {
	return &WakaTimeClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    "https://wakatime.com/api/v1",
	}
}

// GetCodingTime fetches coding time from WakaTime API
func (w *WakaTimeClient) GetCodingTime(ctx context.Context, user, repo string, since time.Time) (*metrics.CodingTime, error) {
	// WakaTime API for summaries
	start := since.Format("2006-01-02")
	end := time.Now().Format("2006-01-02")

	url := fmt.Sprintf("%s/users/%s/summaries?start=%s&end=%s&project=%s",
		w.baseURL, user, start, end, repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+w.apiKey)

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("WakaTime API error: %d", resp.StatusCode)
	}

	var response WakaTimeSummariesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	totalSeconds := 0.0
	codingSeconds := 0.0

	var languages []metrics.LanguageTime
	var projects []metrics.ProjectTime

	for _, day := range response.Data {
		totalSeconds += day.GrandTotal.TotalSeconds

		for _, lang := range day.Languages {
			codingSeconds += lang.TotalSeconds

			// Aggregate languages
			found := false
			for i, existing := range languages {
				if existing.Name == lang.Name {
					languages[i].Hours += lang.TotalSeconds / 3600.0
					found = true
					break
				}
			}
			if !found {
				languages = append(languages, metrics.LanguageTime{
					Name:  lang.Name,
					Hours: lang.TotalSeconds / 3600.0,
				})
			}
		}

		for _, proj := range day.Projects {
			// Aggregate projects
			found := false
			for i, existing := range projects {
				if existing.Name == proj.Name {
					projects[i].Hours += proj.TotalSeconds / 3600.0
					found = true
					break
				}
			}
			if !found {
				projects = append(projects, metrics.ProjectTime{
					Name:  proj.Name,
					Hours: proj.TotalSeconds / 3600.0,
				})
			}
		}
	}

	periodDays := int(time.Since(since).Hours() / 24)
	if periodDays == 0 {
		periodDays = 1
	}

	return &metrics.CodingTime{
		TotalHours:  totalSeconds / 3600.0,
		CodingHours: codingSeconds / 3600.0,
		Period:      periodDays,
		Languages:   languages,
		Projects:    projects,
	}, nil
}

// WakaTimeSummariesResponse API response types
type WakaTimeSummariesResponse struct {
	Data []WakaTimeDaySummary `json:"data"`
}

type WakaTimeDaySummary struct {
	GrandTotal WakaTimeGrandTotal `json:"grand_total"`
	Languages  []WakaTimeLanguage `json:"languages"`
	Projects   []WakaTimeProject  `json:"projects"`
}

type WakaTimeGrandTotal struct {
	TotalSeconds float64 `json:"total_seconds"`
}

type WakaTimeLanguage struct {
	Name         string  `json:"name"`
	TotalSeconds float64 `json:"total_seconds"`
}

type WakaTimeProject struct {
	Name         string  `json:"name"`
	TotalSeconds float64 `json:"total_seconds"`
}
