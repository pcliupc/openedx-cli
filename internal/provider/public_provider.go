package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/openedx/cli/internal/registry"
)

// PublicProvider executes commands against the built-in public Open edX API.
type PublicProvider struct {
	Client *http.Client
	Cmd    registry.CommandMeta
}

// NewPublicProvider creates a PublicProvider for the given command metadata.
// If client is nil, http.DefaultClient is used.
func NewPublicProvider(client *http.Client, cmd registry.CommandMeta) *PublicProvider {
	if client == nil {
		client = http.DefaultClient
	}
	return &PublicProvider{Client: client, Cmd: cmd}
}

// Execute sends an HTTP request to the public API. It builds the URL from
// baseURL + resolved path template, sets the HTTP method from the command
// metadata, attaches the Bearer token, and for POST requests sends args
// as a JSON body.
func (p *PublicProvider) Execute(ctx context.Context, baseURL string, token string, args map[string]string) (*ProviderResult, error) {
	path := resolvePath(p.Cmd.Path, args)
	url := strings.TrimRight(baseURL, "/") + path

	var bodyReader io.Reader
	if p.Cmd.Method == "POST" || p.Cmd.Method == "post" {
		jsonBody, err := json.Marshal(args)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, p.Cmd.Method, url, bodyReader)
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

// resolvePath replaces {param} placeholders in the path template with values
// from the args map.
func resolvePath(pathTmpl string, args map[string]string) string {
	result := pathTmpl
	for key, val := range args {
		result = strings.ReplaceAll(result, "{"+key+"}", val)
	}
	return result
}
