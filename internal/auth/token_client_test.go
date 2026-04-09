package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/openedx/cli/internal/config"
	"github.com/stretchr/testify/assert"
)

// testProfile returns a config.Profile configured to talk to the given tokenURL.
func testProfile(tokenURL string) config.Profile {
	return config.Profile{
		BaseURL:        "https://example.com",
		TokenURL:       tokenURL,
		ClientIDEnv:    "TEST_CLIENT_ID",
		ClientSecretEnv: "TEST_CLIENT_SECRET",
	}
}

// setupEnv sets the test environment variables and returns a cleanup function.
func setupEnv(t *testing.T) {
	t.Helper()
	t.Setenv("TEST_CLIENT_ID", "test-id")
	t.Setenv("TEST_CLIENT_SECRET", "test-secret")
}

func TestClientCredentialsTokenRequest(t *testing.T) {
	setupEnv(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		if err := r.ParseForm(); err != nil {
			t.Fatalf("failed to parse form: %v", err)
		}
		assert.Equal(t, "client_credentials", r.FormValue("grant_type"))
		assert.Equal(t, "test-id", r.FormValue("client_id"))
		assert.Equal(t, "test-secret", r.FormValue("client_secret"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "abc123",
			"expires_in":   3600,
			"token_type":   "Bearer",
		})
	}))
	defer server.Close()

	client := NewHTTPTokenClient(server.Client())
	profile := testProfile(server.URL)

	token, err := client.Token(context.Background(), profile)
	assert.NoError(t, err)
	assert.Equal(t, "abc123", token.AccessToken)
	assert.Equal(t, 3600, token.ExpiresIn)
	assert.Equal(t, "Bearer", token.TokenType)
}

func TestTokenCacheReturnsUnexpiredToken(t *testing.T) {
	setupEnv(t)

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": fmt.Sprintf("token-%d", callCount),
			"expires_in":   3600,
			"token_type":   "Bearer",
		})
	}))
	defer server.Close()

	now := time.Now()
	inner := NewHTTPTokenClient(server.Client())
	cache := NewCachingTokenProvider(inner, func() time.Time { return now })

	profile := testProfile(server.URL)

	// First call should fetch from the server.
	token1, err := cache.Token(context.Background(), profile)
	assert.NoError(t, err)
	assert.Equal(t, "token-1", token1.AccessToken)
	assert.Equal(t, 1, callCount)

	// Second call should return the cached token without hitting the server.
	token2, err := cache.Token(context.Background(), profile)
	assert.NoError(t, err)
	assert.Equal(t, "token-1", token2.AccessToken)
	assert.Equal(t, 1, callCount) // still 1, no additional server call
}

func TestTokenCacheRefreshesExpiredToken(t *testing.T) {
	setupEnv(t)

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": fmt.Sprintf("token-%d", callCount),
			"expires_in":   3600,
			"token_type":   "Bearer",
		})
	}))
	defer server.Close()

	now := time.Now()
	inner := NewHTTPTokenClient(server.Client())
	cache := NewCachingTokenProvider(inner, func() time.Time { return now })

	profile := testProfile(server.URL)

	// First call fetches a token.
	token1, err := cache.Token(context.Background(), profile)
	assert.NoError(t, err)
	assert.Equal(t, "token-1", token1.AccessToken)
	assert.Equal(t, 1, callCount)

	// Advance time past the token expiry (including the 60s margin).
	// Token had expires_in=3600, so we advance 3600 seconds to make it expired.
	now = now.Add(3600 * time.Second)

	// This call should detect expiry and fetch a new token.
	token2, err := cache.Token(context.Background(), profile)
	assert.NoError(t, err)
	assert.Equal(t, "token-2", token2.AccessToken)
	assert.Equal(t, 2, callCount)
}

func TestMissingEnvVarReturnsError(t *testing.T) {
	// Ensure the env vars are NOT set for this test.
	os.Unsetenv("TEST_CLIENT_ID")
	os.Unsetenv("TEST_CLIENT_SECRET")

	client := NewHTTPTokenClient(nil)
	profile := config.Profile{
		TokenURL:       "https://example.com/oauth2/token",
		ClientIDEnv:    "TEST_CLIENT_ID",
		ClientSecretEnv: "TEST_CLIENT_SECRET",
	}

	_, err := client.Token(context.Background(), profile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing env var")
	assert.Contains(t, err.Error(), "TEST_CLIENT_ID")
}

func TestTokenEndpointUnavailableReturnsError(t *testing.T) {
	setupEnv(t)

	// Create a server that is immediately closed to simulate an unavailable endpoint.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()

	client := NewHTTPTokenClient(server.Client())
	profile := testProfile(server.URL)

	_, err := client.Token(context.Background(), profile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token endpoint unavailable")
}

func TestInvalidTokenResponseReturnsError(t *testing.T) {
	setupEnv(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`not valid json`))
	}))
	defer server.Close()

	client := NewHTTPTokenClient(server.Client())
	profile := testProfile(server.URL)

	_, err := client.Token(context.Background(), profile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token response")
}

func TestMissingAccessTokenFieldReturnsError(t *testing.T) {
	setupEnv(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"expires_in": 3600,
			"token_type": "Bearer",
		})
	}))
	defer server.Close()

	client := NewHTTPTokenClient(server.Client())
	profile := testProfile(server.URL)

	_, err := client.Token(context.Background(), profile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access_token is missing")
}

func TestNonOKStatusCodeReturnsError(t *testing.T) {
	setupEnv(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"invalid_client"}`))
	}))
	defer server.Close()

	client := NewHTTPTokenClient(server.Client())
	profile := testProfile(server.URL)

	_, err := client.Token(context.Background(), profile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token request failed")
	assert.Contains(t, err.Error(), "401")
}
