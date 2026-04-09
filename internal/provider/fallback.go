package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/openedx/cli/internal/auth"
	"github.com/openedx/cli/internal/config"
	"github.com/openedx/cli/internal/registry"
)

// FallbackProvider wraps a public provider and an extension provider,
// implementing a fallback policy: if the public API returns an
// "unavailable" status code (404, 405, 501) and an extension mapping
// exists for the command, the request is retried against the extension.
type FallbackProvider struct {
	Public    *PublicProvider
	Extension *ExtensionProvider
}

// NewFallbackProvider creates a FallbackProvider with the given HTTP client.
// If client is nil, http.DefaultClient is used for both inner providers.
func NewFallbackProvider(client *http.Client) *FallbackProvider {
	if client == nil {
		client = http.DefaultClient
	}
	return &FallbackProvider{
		Public:    nil, // set per-command via cmd metadata
		Extension: NewExtensionProvider(client),
	}
}

// Execute dispatches a command using the fallback policy:
//  1. Acquire a token from the token provider.
//  2. Execute against the public API.
//  3. If the public API returns 2xx, return the result immediately.
//  4. If the public API returns an "unavailable" error (404, 405, 501) AND
//     an extension mapping is provided, retry against the extension endpoint.
//  5. For all other errors (400, 401, 403, 500, etc.) or when no extension
//     mapping exists, return the public error directly.
func (f *FallbackProvider) Execute(ctx context.Context, profile config.Profile, tokenProvider auth.TokenProvider, cmdKey string, cmd registry.CommandMeta, extMapping *config.ExtensionMapping, args map[string]string) (*ProviderResult, error) {
	// Step 1: Acquire token.
	token, err := tokenProvider.Token(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire token: %w", err)
	}

	// Step 2: Execute against public provider.
	public := &PublicProvider{Client: f.Extension.Client, Cmd: cmd}
	result, err := public.Execute(ctx, profile.BaseURL, token.AccessToken, args)
	if err == nil {
		// Step 3: Public succeeded.
		return result, nil
	}

	// Step 4: Check if the error is "unavailable" for fallback.
	provErr, ok := err.(*ProviderError)
	if !ok || !provErr.IsUnavailable() {
		// Not an unavailable error (400, 401, 403, 500, etc.) - return as-is.
		return nil, err
	}

	// Step 5: If no extension mapping, return the public error.
	if extMapping == nil {
		return nil, err
	}

	// Step 6: Retry via extension provider.
	extResult, extErr := f.Extension.Execute(ctx, token.AccessToken, *extMapping, args)
	if extErr != nil {
		return nil, extErr
	}
	return extResult, nil
}
