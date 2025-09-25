// Package github provides GitHub authentication with App JWT and PAT fallback.
package github

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthProvider handles GitHub authentication
type AuthProvider struct {
	config     *Config
	httpClient *http.Client

	// App authentication
	privateKey *rsa.PrivateKey
	appJWT     string
	jwtExpiry  time.Time

	// Installation tokens cache
	installationTokens map[int64]*InstallationToken
	tokenMutex         sync.RWMutex
}

// InstallationToken represents a GitHub App installation token
type InstallationToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// NewAuthProvider creates a new GitHub authentication provider
func NewAuthProvider(config *Config, httpClient *http.Client) (*AuthProvider, error) {
	auth := &AuthProvider{
		config:             config,
		httpClient:         httpClient,
		installationTokens: make(map[int64]*InstallationToken),
	}

	// Parse private key if App authentication is configured
	if config.IsAppAuthConfigured() {
		privateKey, err := parsePrivateKey(config.AppPrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse GitHub App private key: %w", err)
		}
		auth.privateKey = privateKey
	}

	return auth, nil
}

// GetAuthToken returns an authentication token for API requests
func (a *AuthProvider) GetAuthToken(installationID int64) (string, error) {
	// Use PAT if App auth is not configured
	if !a.config.IsAppAuthConfigured() {
		if !a.config.IsPATAuthConfigured() {
			return "", fmt.Errorf("no authentication configured")
		}
		return a.config.PersonalAccessToken, nil
	}

	// Use installation token for App auth
	return a.getInstallationToken(installationID)
}

// GetAppJWT returns a JWT for GitHub App authentication
func (a *AuthProvider) GetAppJWT() (string, error) {
	if !a.config.IsAppAuthConfigured() {
		return "", fmt.Errorf("GitHub App authentication not configured")
	}

	// Check if we have a valid cached JWT
	if a.appJWT != "" && time.Now().Before(a.jwtExpiry.Add(-1*time.Minute)) {
		return a.appJWT, nil
	}

	// Generate new JWT
	now := time.Now()
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(10 * time.Minute).Unix(), // JWTs are valid for 10 minutes
		"iss": a.config.AppID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	jwtString, err := token.SignedString(a.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	// Cache the JWT
	a.appJWT = jwtString
	a.jwtExpiry = now.Add(10 * time.Minute)

	return jwtString, nil
}

// getInstallationToken gets a valid installation token
func (a *AuthProvider) getInstallationToken(installationID int64) (string, error) {
	if installationID == 0 {
		installationID = a.config.InstallationID
	}

	if installationID == 0 {
		return "", fmt.Errorf("no installation ID provided")
	}

	a.tokenMutex.RLock()
	if token, exists := a.installationTokens[installationID]; exists {
		// Check if token is still valid (with 5-minute buffer)
		if time.Now().Before(token.ExpiresAt.Add(-5 * time.Minute)) {
			a.tokenMutex.RUnlock()
			return token.Token, nil
		}
	}
	a.tokenMutex.RUnlock()

	// Need to refresh the token
	return a.refreshInstallationToken(installationID)
}

// refreshInstallationToken requests a new installation token
func (a *AuthProvider) refreshInstallationToken(installationID int64) (string, error) {
	a.tokenMutex.Lock()
	defer a.tokenMutex.Unlock()

	// Double-check if another goroutine already refreshed the token
	if token, exists := a.installationTokens[installationID]; exists {
		if time.Now().Before(token.ExpiresAt.Add(-5 * time.Minute)) {
			return token.Token, nil
		}
	}

	// Get App JWT for authentication
	appJWT, err := a.GetAppJWT()
	if err != nil {
		return "", fmt.Errorf("failed to get App JWT: %w", err)
	}

	// Request installation token
	url := fmt.Sprintf("%s/app/installations/%d/access_tokens", a.config.BaseURL, installationID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+appJWT)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("X-GitHub-Api-Version", a.config.APIVersion)
	req.Header.Set("User-Agent", a.config.UserAgent)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to request installation token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	var tokenResponse struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	// Cache the token
	a.installationTokens[installationID] = &InstallationToken{
		Token:     tokenResponse.Token,
		ExpiresAt: tokenResponse.ExpiresAt,
	}

	return tokenResponse.Token, nil
}

// GetInstallationID gets the installation ID for a repository
func (a *AuthProvider) GetInstallationID(owner, repo string) (int64, error) {
	if !a.config.IsAppAuthConfigured() {
		return 0, fmt.Errorf("GitHub App authentication not configured")
	}

	// Use configured installation ID if available
	if a.config.InstallationID > 0 {
		return a.config.InstallationID, nil
	}

	// Get App JWT for authentication
	appJWT, err := a.GetAppJWT()
	if err != nil {
		return 0, fmt.Errorf("failed to get App JWT: %w", err)
	}

	// Request installation for repository
	url := fmt.Sprintf("%s/repos/%s/%s/installation", a.config.BaseURL, owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+appJWT)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("X-GitHub-Api-Version", a.config.APIVersion)
	req.Header.Set("User-Agent", a.config.UserAgent)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to get installation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	var installation struct {
		ID int64 `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&installation); err != nil {
		return 0, fmt.Errorf("failed to decode installation response: %w", err)
	}

	return installation.ID, nil
}

// IsUsingAppAuth returns true if App authentication is being used
func (a *AuthProvider) IsUsingAppAuth() bool {
	return a.config.IsAppAuthConfigured()
}

// IsUsingPAT returns true if Personal Access Token authentication is being used
func (a *AuthProvider) IsUsingPAT() bool {
	return !a.config.IsAppAuthConfigured() && a.config.IsPATAuthConfigured()
}

// parsePrivateKey parses a PEM-encoded RSA private key
func parsePrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS#8 format
		parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}

		rsaKey, ok := parsedKey.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is not RSA")
		}
		return rsaKey, nil
	}

	return key, nil
}