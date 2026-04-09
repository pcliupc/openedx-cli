package diagnostics

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openedx/cli/internal/auth"
	"github.com/openedx/cli/internal/config"
)

// --- mock token provider ---

type mockTokenProvider struct {
	token *auth.Token
	err   error
}

func (m *mockTokenProvider) Token(_ context.Context, _ config.Profile) (*auth.Token, error) {
	return m.token, m.err
}

// --- CheckBaseURL tests ---

func TestCheckBaseURLReachable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result := CheckBaseURL(context.Background(), server.URL)

	if result.Status != "ok" {
		t.Errorf("expected status ok, got %q: %s", result.Status, result.Message)
	}
	if result.Name != "base_url" {
		t.Errorf("expected name base_url, got %q", result.Name)
	}
}

func TestCheckBaseURLUnreachable(t *testing.T) {
	result := CheckBaseURL(context.Background(), "http://127.0.0.1:0")

	if result.Status != "error" {
		t.Errorf("expected status error, got %q", result.Status)
	}
}

func TestCheckBaseURLEmpty(t *testing.T) {
	result := CheckBaseURL(context.Background(), "")

	if result.Status != "error" {
		t.Errorf("expected status error, got %q", result.Status)
	}
	if result.Message != "base_url is not configured" {
		t.Errorf("unexpected message: %s", result.Message)
	}
}

// --- CheckTokenAcquisition tests ---

func TestCheckTokenAcquisitionSuccess(t *testing.T) {
	provider := &mockTokenProvider{
		token: &auth.Token{
			AccessToken: "test-token",
			ExpiresIn:   3600,
			TokenType:   "Bearer",
		},
	}

	profile := config.Profile{BaseURL: "https://example.com"}
	result := CheckTokenAcquisition(context.Background(), profile, provider)

	if result.Status != "ok" {
		t.Errorf("expected status ok, got %q: %s", result.Status, result.Message)
	}
	if result.Name != "token_acquisition" {
		t.Errorf("expected name token_acquisition, got %q", result.Name)
	}
}

func TestCheckTokenAcquisitionFailure(t *testing.T) {
	provider := &mockTokenProvider{
		err: context.DeadlineExceeded,
	}

	profile := config.Profile{BaseURL: "https://example.com"}
	result := CheckTokenAcquisition(context.Background(), profile, provider)

	if result.Status != "error" {
		t.Errorf("expected status error, got %q", result.Status)
	}
}

func TestCheckTokenAcquisitionNilProvider(t *testing.T) {
	profile := config.Profile{BaseURL: "https://example.com"}
	result := CheckTokenAcquisition(context.Background(), profile, nil)

	if result.Status != "error" {
		t.Errorf("expected status error, got %q", result.Status)
	}
	if result.Message != "token provider is not configured" {
		t.Errorf("unexpected message: %s", result.Message)
	}
}

// --- CheckCommand tests ---

func TestDoctorVerifyCommandChecksExtensionAvailability(t *testing.T) {
	ctx := context.Background()
	profile := config.Profile{BaseURL: "https://example.com"}

	// Test command in public registry.
	result := CheckCommand(ctx, "course.list", profile, nil)
	if result.Status != "ok" {
		t.Errorf("expected status ok for public registry command, got %q: %s", result.Status, result.Message)
	}

	// Test command in extension mappings.
	extensions := map[string]config.ExtensionMapping{
		"course.custom": {
			Method: "GET",
			URL:    "https://custom.example.com/api/course",
		},
	}
	result = CheckCommand(ctx, "course.custom", profile, extensions)
	if result.Status != "ok" {
		t.Errorf("expected status ok for extension command, got %q: %s", result.Status, result.Message)
	}

	// Test command not found.
	result = CheckCommand(ctx, "nonexistent.command", profile, extensions)
	if result.Status != "error" {
		t.Errorf("expected status error for unknown command, got %q", result.Status)
	}

	// Test extension mapping with empty URL.
	extensions["broken.command"] = config.ExtensionMapping{
		Method: "POST",
		URL:    "",
	}
	result = CheckCommand(ctx, "broken.command", profile, extensions)
	if result.Status != "error" {
		t.Errorf("expected status error for extension with empty URL, got %q", result.Status)
	}
}

// --- RunAllChecks tests ---

func TestDoctorRunAllChecks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	provider := &mockTokenProvider{
		token: &auth.Token{
			AccessToken: "test-token",
			ExpiresIn:   3600,
			TokenType:   "Bearer",
		},
	}

	profile := config.Profile{
		BaseURL: server.URL,
	}

	extensions := map[string]config.ExtensionMapping{}

	result := RunAllChecks(context.Background(), profile, provider, extensions)

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Checks) != 2 {
		t.Fatalf("expected 2 checks, got %d", len(result.Checks))
	}

	// Check base URL result.
	if result.Checks[0].Name != "base_url" {
		t.Errorf("expected first check name 'base_url', got %q", result.Checks[0].Name)
	}
	if result.Checks[0].Status != "ok" {
		t.Errorf("expected base_url status ok, got %q: %s", result.Checks[0].Status, result.Checks[0].Message)
	}

	// Check token acquisition result.
	if result.Checks[1].Name != "token_acquisition" {
		t.Errorf("expected second check name 'token_acquisition', got %q", result.Checks[1].Name)
	}
	if result.Checks[1].Status != "ok" {
		t.Errorf("expected token_acquisition status ok, got %q: %s", result.Checks[1].Status, result.Checks[1].Message)
	}
}

func TestDoctorRunAllChecksFailures(t *testing.T) {
	profile := config.Profile{
		BaseURL: "http://127.0.0.1:0",
	}

	provider := &mockTokenProvider{
		err: context.DeadlineExceeded,
	}

	result := RunAllChecks(context.Background(), profile, provider, nil)

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Checks) != 2 {
		t.Fatalf("expected 2 checks, got %d", len(result.Checks))
	}

	if result.Checks[0].Status != "error" {
		t.Errorf("expected base_url status error, got %q", result.Checks[0].Status)
	}
	if result.Checks[1].Status != "error" {
		t.Errorf("expected token_acquisition status error, got %q", result.Checks[1].Status)
	}
}
