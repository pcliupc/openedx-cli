// Package auth provides OAuth token acquisition and caching for the Open edX CLI.
// It implements the client credentials flow used for machine-to-machine
// authentication in CI pipelines and automation scenarios.
package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/openedx/cli/internal/config"
)

// Token represents an OAuth access token returned by the token endpoint.
type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// TokenProvider is the interface for acquiring OAuth tokens. Implementations
// may add caching, retries, or other middleware around the core token request.
type TokenProvider interface {
	Token(ctx context.Context, profile config.Profile) (*Token, error)
}

// HTTPTokenClient acquires tokens via the OAuth client credentials grant using
// an *http.Client. By default it uses http.DefaultClient.
type HTTPTokenClient struct {
	HTTPClient *http.Client
}

// NewHTTPTokenClient creates a new HTTPTokenClient with the given HTTP client.
// If client is nil, http.DefaultClient is used.
func NewHTTPTokenClient(client *http.Client) *HTTPTokenClient {
	if client == nil {
		client = http.DefaultClient
	}
	return &HTTPTokenClient{HTTPClient: client}
}

// Token sends a client credentials token request to the profile's token endpoint
// and returns the parsed token response.
func (c *HTTPTokenClient) Token(ctx context.Context, profile config.Profile) (*Token, error) {
	clientID := os.Getenv(profile.ClientIDEnv)
	clientSecret := os.Getenv(profile.ClientSecretEnv)

	if strings.TrimSpace(clientID) == "" {
		return nil, fmt.Errorf("missing env var: %s is not set", profile.ClientIDEnv)
	}
	if strings.TrimSpace(clientSecret) == "" {
		return nil, fmt.Errorf("missing env var: %s is not set", profile.ClientSecretEnv)
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, profile.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token endpoint unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token request failed (status %d): %s", resp.StatusCode, string(body))
	}

	var token Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("invalid token response: %w", err)
	}

	if strings.TrimSpace(token.AccessToken) == "" {
		return nil, fmt.Errorf("invalid token response: access_token is missing")
	}

	return &token, nil
}
