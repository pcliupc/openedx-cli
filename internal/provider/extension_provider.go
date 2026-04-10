package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/openedx/cli/internal/config"
)

// ExtensionProvider executes commands against a custom extension endpoint
// defined in the user's configuration. Unlike PublicProvider, the extension
// mapping provides a full URL so no path template resolution is needed.
type ExtensionProvider struct {
	Client *http.Client
}

// NewExtensionProvider creates an ExtensionProvider with the given HTTP client.
// If client is nil, http.DefaultClient is used.
func NewExtensionProvider(client *http.Client) *ExtensionProvider {
	if client == nil {
		client = http.DefaultClient
	}
	return &ExtensionProvider{Client: client}
}

// Execute sends an HTTP request to the extension endpoint defined by ext.
// It uses ext.Method and ext.URL directly, sends args as JSON for POST
// requests, and attaches the Bearer token.
func (p *ExtensionProvider) Execute(ctx context.Context, token string, ext config.ExtensionMapping, args map[string]string) (*ProviderResult, error) {
	reqURL := ext.URL
	var bodyReader io.Reader
	if ext.Method == "POST" || ext.Method == "post" {
		jsonBody, err := json.Marshal(args)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	} else if len(args) > 0 {
		params := url.Values{}
		for key, val := range args {
			params.Set(key, val)
		}
		if strings.Contains(reqURL, "?") {
			reqURL += "&" + params.Encode()
		} else {
			reqURL += "?" + params.Encode()
		}
	}

	req, err := http.NewRequestWithContext(ctx, ext.Method, reqURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if bodyReader != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &ProviderError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	return &ProviderResult{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}
