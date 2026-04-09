// Package config defines the configuration types and loading logic for the
// Open edX CLI. Configuration is loaded from YAML files and secrets are
// referenced by environment variable names (resolved at runtime).
package config

// Config is the top-level configuration structure for the CLI.
type Config struct {
	Version    int                         `mapstructure:"version"`
	Profiles   map[string]Profile          `mapstructure:"profiles"`
	Extensions map[string]ExtensionMapping `mapstructure:"extensions"`
}

// Profile holds connection details for a specific Open edX deployment.
// ClientIDEnv and ClientSecretEnv store the *names* of environment variables
// that contain the actual credentials, not the credential values themselves.
type Profile struct {
	BaseURL        string `mapstructure:"base_url"`
	TokenURL       string `mapstructure:"token_url"`
	ClientIDEnv    string `mapstructure:"client_id_env"`
	ClientSecretEnv string `mapstructure:"client_secret_env"`
	DefaultFormat  string `mapstructure:"default_format"`
}

// ExtensionMapping defines a custom API endpoint that replaces a built-in
// command when the official public API does not provide the required
// functionality.
type ExtensionMapping struct {
	Method string `mapstructure:"method"`
	URL    string `mapstructure:"url"`
}
