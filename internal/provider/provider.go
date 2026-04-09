// Package provider implements the execution layer that dispatches CLI commands
// to backend API endpoints. It defines the Provider interface, concrete HTTP
// implementations for public and extension backends, and a fallback policy
// that transparently retries unavailable commands against extension endpoints.
package provider

import "context"

// ProviderResult holds the response from a provider execution.
type ProviderResult struct {
	StatusCode int
	Body       []byte
}

// Provider executes a CLI command against a backend API.
type Provider interface {
	Execute(ctx context.Context, baseURL string, token string, args map[string]string) (*ProviderResult, error)
}

// ProviderError captures an HTTP-level error from a provider execution.
// It implements the error interface and provides an IsUnavailable method
// to check whether the error qualifies for fallback to an extension provider.
type ProviderError struct {
	StatusCode int
	Message    string
}

// Error returns a human-readable description of the provider error.
func (e *ProviderError) Error() string {
	return e.Message
}

// IsUnavailable returns true for HTTP status codes that indicate the public
// API does not support the requested operation, making it eligible for
// fallback to an extension provider. These codes are:
//   - 404 Not Found
//   - 405 Method Not Allowed
//   - 501 Not Implemented
func (e *ProviderError) IsUnavailable() bool {
	return e.StatusCode == 404 || e.StatusCode == 405 || e.StatusCode == 501
}
