// Package metrics - GraphQL client for GitHub API heavy aggregations
package metrics

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AuthProvider interface for getting authentication tokens
type AuthProvider interface {
	GetAuthToken() (string, error)
}

// GraphQLClient provides GraphQL queries for complex metrics aggregations
type GraphQLClient struct {
	httpClient   *http.Client
	baseURL      string
	authProvider AuthProvider
}

// NewGraphQLClient creates a new GraphQL client
func NewGraphQLClient(authProvider AuthProvider, baseURL string) *GraphQLClient {
	if baseURL == "" {
		baseURL = "https://api.github.com/graphql"
	}

	return &GraphQLClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:      baseURL,
		authProvider: authProvider,
	}
}

// GraphQL request and response types

// GraphQLRequest represents a GraphQL query request
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse represents a GraphQL query response
type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors,omitempty"`
}

// GraphQLError represents a GraphQL error
type GraphQLError struct {
	Message   string      `json:"message"`
	Locations []Location  `json:"locations,omitempty"`
	Path      []string    `json:"path,omitempty"`
	Extensions interface{} `json:"extensions,omitempty"`
}

// Location represents error location in GraphQL query
type Location struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Complex aggregation data types

// RepositoryMetricsData represents comprehensive repository data from GraphQL
type RepositoryMetricsData struct {
	Repository struct {
		Name              string    `json:"name"`
		Owner             Owner     `json:"owner"`
		CreatedAt         time.Time `json:"createdAt"`
		UpdatedAt         time.Time `json:"updatedAt"`
		PrimaryLanguage   Language  `json:"primaryLanguage"`
		Languages         Languages `json:"languages"`
		DefaultBranchRef  BranchRef `json:"defaultBranchRef"`
		PullRequests      PullRequests `json:"pullRequests"`
		Issues            Issues    `json:"issues"`
		Releases          Releases  `json:"releases"`
		Deployments       Deployments `json:"deployments"`
		Collaborators     Collaborators `json:"collaborators"`
		CommitComments    CommitComments `json:"commitComments"`
		DiskUsage         int       `json:"diskUsage"`
		ForkCount         int       `json:"forkCount"`
		StargazerCount    int       `json:"stargazerCount"`
		WatcherCount      int       `json:"watchers"`
	} `json:"repository"`
}

// Supporting types for GraphQL responses

type Owner struct {
	Login string `json:"login"`
	Type  string `json:"__typename"`
}

type Language struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Languages struct {
	TotalSize int `json:"totalSize"`
	Edges     []LanguageEdge `json:"edges"`
}

type LanguageEdge struct {
	Size int      `json:"size"`
	Node Language `json:"node"`
}

type BranchRef struct {
	Name   string `json:"name"`
	Target Target `json:"target"`
}

type Target struct {
	Oid     string    `json:"oid"`
	History History   `json:"history"`
}

type History struct {
	TotalCount int     `json:"totalCount"`
	Nodes      []GraphQLCommit `json:"nodes"`
}

type GraphQLCommit struct {
	Oid             string    `json:"oid"`
	Message         string    `json:"message"`
	CommittedDate   time.Time `json:"committedDate"`
	Author          GitActor  `json:"author"`
	Committer       GitActor  `json:"committer"`
	Additions       int       `json:"additions"`
	Deletions       int       `json:"deletions"`
	ChangedFiles    int       `json:"changedFiles"`
	AssociatedPullRequests PullRequestConnection `json:"associatedPullRequests"`
}

type GitActor struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
	User  User      `json:"user"`
}

type User struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type PullRequests struct {
	TotalCount int                   `json:"totalCount"`
	Nodes      []GraphQLPullRequest         `json:"nodes"`
}

type GraphQLPullRequest struct {
	Number        int       `json:"number"`
	Title         string    `json:"title"`
	State         string    `json:"state"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	MergedAt      *time.Time `json:"mergedAt"`
	ClosedAt      *time.Time `json:"closedAt"`
	Author        User      `json:"author"`
	Mergeable     string    `json:"mergeable"`
	Additions     int       `json:"additions"`
	Deletions     int       `json:"deletions"`
	ChangedFiles  int       `json:"changedFiles"`
	Reviews       Reviews   `json:"reviews"`
	Comments      Comments  `json:"comments"`
	Commits       PullRequestCommits `json:"commits"`
	Labels        Labels    `json:"labels"`
	Assignees     Assignees `json:"assignees"`
}

type PullRequestConnection struct {
	TotalCount int           `json:"totalCount"`
	Nodes      []GraphQLPullRequest `json:"nodes"`
}

type Reviews struct {
	TotalCount int      `json:"totalCount"`
	Nodes      []Review `json:"nodes"`
}

type Review struct {
	State       string    `json:"state"`
	SubmittedAt time.Time `json:"submittedAt"`
	Author      User      `json:"author"`
}

type Comments struct {
	TotalCount int       `json:"totalCount"`
	Nodes      []Comment `json:"nodes"`
}

type Comment struct {
	CreatedAt time.Time `json:"createdAt"`
	Author    User      `json:"author"`
	Body      string    `json:"body"`
}

type PullRequestCommits struct {
	TotalCount int      `json:"totalCount"`
	Nodes      []GraphQLCommit `json:"nodes"`
}

type Labels struct {
	TotalCount int     `json:"totalCount"`
	Nodes      []Label `json:"nodes"`
}

type Label struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

type Assignees struct {
	TotalCount int    `json:"totalCount"`
	Nodes      []User `json:"nodes"`
}

type Issues struct {
	TotalCount int     `json:"totalCount"`
	Nodes      []GraphQLIssue `json:"nodes"`
}

type GraphQLIssue struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ClosedAt  *time.Time `json:"closedAt"`
	Author    User      `json:"author"`
	Labels    Labels    `json:"labels"`
	Assignees Assignees `json:"assignees"`
	Comments  Comments  `json:"comments"`
}

type Releases struct {
	TotalCount int       `json:"totalCount"`
	Nodes      []Release `json:"nodes"`
}

type Release struct {
	Name         string    `json:"name"`
	TagName      string    `json:"tagName"`
	CreatedAt    time.Time `json:"createdAt"`
	PublishedAt  time.Time `json:"publishedAt"`
	Author       User      `json:"author"`
	IsPrerelease bool      `json:"isPrerelease"`
	IsDraft      bool      `json:"isDraft"`
}

type Deployments struct {
	TotalCount int          `json:"totalCount"`
	Nodes      []GraphQLDeployment `json:"nodes"`
}

type GraphQLDeployment struct {
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Environment   string    `json:"environment"`
	State         string    `json:"state"`
	Description   string    `json:"description"`
	Creator       User      `json:"creator"`
	Ref           Reference `json:"ref"`
	Statuses      DeploymentStatuses `json:"statuses"`
}

type Reference struct {
	Name   string `json:"name"`
	Target Target `json:"target"`
}

type DeploymentStatuses struct {
	TotalCount int                `json:"totalCount"`
	Nodes      []DeploymentStatus `json:"nodes"`
}

type DeploymentStatus struct {
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Description string    `json:"description"`
	Creator     User      `json:"creator"`
}

type Collaborators struct {
	TotalCount int    `json:"totalCount"`
	Nodes      []User `json:"nodes"`
}

type CommitComments struct {
	TotalCount int             `json:"totalCount"`
	Nodes      []CommitComment `json:"nodes"`
}

type CommitComment struct {
	CreatedAt time.Time `json:"createdAt"`
	Author    User      `json:"author"`
	Body      string    `json:"body"`
	Path      string    `json:"path"`
	Position  int       `json:"position"`
}

// GraphQL query methods

// GetRepositoryMetrics fetches comprehensive repository metrics using GraphQL
func (gc *GraphQLClient) GetRepositoryMetrics(ctx context.Context, owner, repo string, since time.Time) (*RepositoryMetricsData, error) {
	query := `
	query GetRepositoryMetrics($owner: String!, $name: String!, $since: DateTime!) {
		repository(owner: $owner, name: $name) {
			name
			owner {
				login
				__typename
			}
			createdAt
			updatedAt
			primaryLanguage {
				name
				color
			}
			languages(first: 10, orderBy: {field: SIZE, direction: DESC}) {
				totalSize
				edges {
					size
					node {
						name
						color
					}
				}
			}
			defaultBranchRef {
				name
				target {
					... on Commit {
						oid
						history(first: 100, since: $since) {
							totalCount
							nodes {
								oid
								message
								committedDate
								author {
									name
									email
									date
									user {
										login
										name
										email
									}
								}
								committer {
									name
									email
									date
									user {
										login
									}
								}
								additions
								deletions
								changedFiles
								associatedPullRequests(first: 5) {
									totalCount
									nodes {
										number
										title
										state
										mergedAt
									}
								}
							}
						}
					}
				}
			}
			pullRequests(first: 100, states: [MERGED, CLOSED], orderBy: {field: UPDATED_AT, direction: DESC}) {
				totalCount
				nodes {
					number
					title
					state
					createdAt
					updatedAt
					mergedAt
					closedAt
					author {
						login
						name
						email
					}
					mergeable
					additions
					deletions
					changedFiles
					reviews(first: 10) {
						totalCount
						nodes {
							state
							submittedAt
							author {
								login
							}
						}
					}
					comments {
						totalCount
					}
					commits {
						totalCount
					}
					labels(first: 10) {
						totalCount
						nodes {
							name
							color
							description
						}
					}
					assignees(first: 5) {
						totalCount
						nodes {
							login
							name
						}
					}
				}
			}
			issues(first: 50, states: [OPEN, CLOSED], orderBy: {field: UPDATED_AT, direction: DESC}) {
				totalCount
				nodes {
					number
					title
					state
					createdAt
					updatedAt
					closedAt
					author {
						login
						name
					}
					labels(first: 10) {
						nodes {
							name
							color
						}
					}
					assignees(first: 5) {
						nodes {
							login
						}
					}
					comments {
						totalCount
					}
				}
			}
			releases(first: 20, orderBy: {field: CREATED_AT, direction: DESC}) {
				totalCount
				nodes {
					name
					tagName
					createdAt
					publishedAt
					author {
						login
						name
					}
					isPrerelease
					isDraft
				}
			}
			deployments(first: 50, orderBy: {field: CREATED_AT, direction: DESC}) {
				totalCount
				nodes {
					createdAt
					updatedAt
					environment
					state
					description
					creator {
						login
					}
					ref {
						name
						target {
							... on Commit {
								oid
							}
						}
					}
					statuses(first: 5) {
						totalCount
						nodes {
							state
							createdAt
							updatedAt
							description
							creator {
								login
							}
						}
					}
				}
			}
			collaborators(first: 50) {
				totalCount
				nodes {
					login
					name
					email
				}
			}
			commitComments(first: 50) {
				totalCount
				nodes {
					createdAt
					author {
						login
					}
					body
					path
					position
				}
			}
			diskUsage
			forkCount
			stargazerCount
			watchers {
				totalCount
			}
		}
	}`

	variables := map[string]interface{}{
		"owner": owner,
		"name":  repo,
		"since": since.Format(time.RFC3339),
	}

	var response RepositoryMetricsData
	err := gc.executeQuery(ctx, query, variables, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to execute repository metrics query: %w", err)
	}

	return &response, nil
}

// GetMultipleRepositoryMetrics fetches metrics for multiple repositories in a single query
func (gc *GraphQLClient) GetMultipleRepositoryMetrics(ctx context.Context, repositories []string, since time.Time) (map[string]*RepositoryMetricsData, error) {
	// This would implement a more complex query to fetch multiple repositories
	// For now, we'll iterate over repositories (can be optimized later)

	results := make(map[string]*RepositoryMetricsData)

	for _, repoFullName := range repositories {
		parts := splitRepositoryName(repoFullName)
		if len(parts) != 2 {
			continue
		}

		data, err := gc.GetRepositoryMetrics(ctx, parts[0], parts[1], since)
		if err != nil {
			// Log error but continue with other repositories
			fmt.Printf("Warning: failed to get metrics for %s: %v\n", repoFullName, err)
			continue
		}

		results[repoFullName] = data
	}

	return results, nil
}

// GetOrganizationRepositories fetches all repositories for an organization
func (gc *GraphQLClient) GetOrganizationRepositories(ctx context.Context, org string, limit int) ([]string, error) {
	query := `
	query GetOrgRepositories($org: String!, $first: Int!) {
		organization(login: $org) {
			repositories(first: $first, orderBy: {field: UPDATED_AT, direction: DESC}) {
				totalCount
				nodes {
					name
					owner {
						login
					}
					updatedAt
					primaryLanguage {
						name
					}
				}
			}
		}
	}`

	variables := map[string]interface{}{
		"org":   org,
		"first": limit,
	}

	var response struct {
		Organization struct {
			Repositories struct {
				TotalCount int `json:"totalCount"`
				Nodes      []struct {
					Name  string `json:"name"`
					Owner Owner  `json:"owner"`
				} `json:"nodes"`
			} `json:"repositories"`
		} `json:"organization"`
	}

	err := gc.executeQuery(ctx, query, variables, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization repositories: %w", err)
	}

	var repositories []string
	for _, repo := range response.Organization.Repositories.Nodes {
		repositories = append(repositories, fmt.Sprintf("%s/%s", repo.Owner.Login, repo.Name))
	}

	return repositories, nil
}

// Private methods

// executeQuery executes a GraphQL query
func (gc *GraphQLClient) executeQuery(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error {
	request := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal GraphQL request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", gc.baseURL, bytes.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add authentication
	token, err := gc.authProvider.GetAuthToken()
	if err != nil {
		return fmt.Errorf("failed to get auth token: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := gc.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to execute GraphQL request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GraphQL request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	var graphqlResp GraphQLResponse
	if err := json.Unmarshal(responseBody, &graphqlResp); err != nil {
		return fmt.Errorf("failed to unmarshal GraphQL response: %w", err)
	}

	if len(graphqlResp.Errors) > 0 {
		return fmt.Errorf("GraphQL errors: %+v", graphqlResp.Errors)
	}

	if err := json.Unmarshal(graphqlResp.Data, result); err != nil {
		return fmt.Errorf("failed to unmarshal GraphQL data: %w", err)
	}

	return nil
}

// Helper functions

// splitRepositoryName splits "owner/repo" into ["owner", "repo"]
func splitRepositoryName(fullName string) []string {
	parts := make([]string, 0, 2)
	if fullName == "" {
		return parts
	}

	for _, part := range []string{fullName[:len(fullName)/2], fullName[len(fullName)/2:]} {
		if idx := len(part) - 1; idx >= 0 && part[idx] == '/' {
			parts = append(parts, part[:idx], part[idx+1:])
			break
		}
	}

	// Fallback to simple split
	if len(parts) == 0 {
		for i, char := range fullName {
			if char == '/' {
				return []string{fullName[:i], fullName[i+1:]}
			}
		}
	}

	return parts
}