// Package diagnostics provides health-check ("doctor") functions for the
// Open edX CLI. Each check probes a specific dependency (base URL reachability,
// token acquisition, command mapping) and returns a structured result.
package diagnostics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/openedx/cli/internal/auth"
	"github.com/openedx/cli/internal/config"
	"github.com/openedx/cli/internal/registry"
)

// CheckResult represents the outcome of a single diagnostic check.
type CheckResult struct {
	Name    string `json:"name"`
	Status  string `json:"status"`  // "ok" or "error"
	Message string `json:"message"`
}

// DoctorResult holds the aggregate results of all diagnostic checks.
type DoctorResult struct {
	Checks []CheckResult `json:"checks"`
}

const checkTimeout = 10 * time.Second

// CheckBaseURL performs an HTTP GET against the profile's base URL to verify
// that the Open edX deployment is reachable.
func CheckBaseURL(ctx context.Context, baseURL string) CheckResult {
	if baseURL == "" {
		return CheckResult{
			Name:    "base_url",
			Status:  "error",
			Message: "base_url is not configured",
		}
	}

	ctx, cancel := context.WithTimeout(ctx, checkTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return CheckResult{
			Name:    "base_url",
			Status:  "error",
			Message: fmt.Sprintf("invalid base_url: %v", err),
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return CheckResult{
			Name:    "base_url",
			Status:  "error",
			Message: fmt.Sprintf("base_url unreachable: %v", err),
		}
	}
	defer resp.Body.Close()

	return CheckResult{
		Name:    "base_url",
		Status:  "ok",
		Message: fmt.Sprintf("base_url reachable (status %d)", resp.StatusCode),
	}
}

// CheckTokenAcquisition attempts to acquire an OAuth token using the provided
// token provider and profile.
func CheckTokenAcquisition(ctx context.Context, profile config.Profile, tokenProvider auth.TokenProvider) CheckResult {
	if tokenProvider == nil {
		return CheckResult{
			Name:    "token_acquisition",
			Status:  "error",
			Message: "token provider is not configured",
		}
	}

	token, err := tokenProvider.Token(ctx, profile)
	if err != nil {
		return CheckResult{
			Name:    "token_acquisition",
			Status:  "error",
			Message: fmt.Sprintf("token acquisition failed: %v", err),
		}
	}

	return CheckResult{
		Name:    "token_acquisition",
		Status:  "ok",
		Message: fmt.Sprintf("token acquired (type=%s, expires_in=%d)", token.TokenType, token.ExpiresIn),
	}
}

// CheckCommand verifies that a CLI command key has a valid mapping in either
// the built-in public registry or the configured extension mappings.
func CheckCommand(ctx context.Context, cmdKey string, profile config.Profile, extensions map[string]config.ExtensionMapping) CheckResult {
	reg := registry.LatestRegistry()

	if _, ok := reg[cmdKey]; ok {
		return CheckResult{
			Name:    "command_mapping",
			Status:  "ok",
			Message: fmt.Sprintf("command %q found in public registry", cmdKey),
		}
	}

	if ext, ok := extensions[cmdKey]; ok {
		if ext.URL == "" {
			return CheckResult{
				Name:    "command_mapping",
				Status:  "error",
				Message: fmt.Sprintf("extension mapping for %q exists but has empty URL", cmdKey),
			}
		}
		return CheckResult{
			Name:    "command_mapping",
			Status:  "ok",
			Message: fmt.Sprintf("command %q found in extension mappings (method=%s, url=%s)", cmdKey, ext.Method, ext.URL),
		}
	}

	return CheckResult{
		Name:    "command_mapping",
		Status:  "error",
		Message: fmt.Sprintf("command %q not found in public registry or extension mappings", cmdKey),
	}
}

// RunAllChecks runs the standard set of diagnostic checks: base URL
// reachability and token acquisition. It returns a DoctorResult containing
// the outcome of each individual check.
func RunAllChecks(ctx context.Context, profile config.Profile, tokenProvider auth.TokenProvider, extensions map[string]config.ExtensionMapping) *DoctorResult {
	result := &DoctorResult{}

	result.Checks = append(result.Checks, CheckBaseURL(ctx, profile.BaseURL))
	result.Checks = append(result.Checks, CheckTokenAcquisition(ctx, profile, tokenProvider))

	return result
}
