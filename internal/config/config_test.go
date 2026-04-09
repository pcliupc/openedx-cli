package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// exampleConfigPath returns the path to the example config fixture.
func exampleConfigPath() string {
	return filepath.Join("..", "..", "testdata", "config", "example.yaml")
}

func TestLoadConfigProfiles(t *testing.T) {
	cfg, err := Load(exampleConfigPath())
	require.NoError(t, err, "Load should succeed with the example config")

	require.NotNil(t, cfg.Profiles, "Profiles map should not be nil")
	assert.Contains(t, cfg.Profiles, "admin", "config should contain 'admin' profile")
	assert.Contains(t, cfg.Profiles, "ops", "config should contain 'ops' profile")

	admin := cfg.Profiles["admin"]
	assert.Equal(t, "https://openedx.example.com", admin.BaseURL)
	assert.Equal(t, "https://openedx.example.com/oauth2/access_token", admin.TokenURL)
	assert.Equal(t, "json", admin.DefaultFormat)

	ops := cfg.Profiles["ops"]
	assert.Equal(t, "https://openedx.example.com", ops.BaseURL)
	assert.Equal(t, "https://openedx.example.com/oauth2/access_token", ops.TokenURL)
}

func TestLoadConfigExtensions(t *testing.T) {
	cfg, err := Load(exampleConfigPath())
	require.NoError(t, err, "Load should succeed with the example config")

	require.NotNil(t, cfg.Extensions, "Extensions map should not be nil")
	assert.Contains(t, cfg.Extensions, "course.create")
	assert.Contains(t, cfg.Extensions, "course.import")
	assert.Contains(t, cfg.Extensions, "course.export")

	create := cfg.Extensions["course.create"]
	assert.Equal(t, "POST", create.Method)
	assert.Equal(t, "https://openedx.example.com/api/cli-ext/course/create", create.URL)

	imp := cfg.Extensions["course.import"]
	assert.Equal(t, "POST", imp.Method)
	assert.Equal(t, "https://openedx.example.com/api/cli-ext/course/import", imp.URL)

	exp := cfg.Extensions["course.export"]
	assert.Equal(t, "POST", exp.Method)
	assert.Equal(t, "https://openedx.example.com/api/cli-ext/course/export", exp.URL)
}

func TestSecretsComeFromEnvNames(t *testing.T) {
	cfg, err := Load(exampleConfigPath())
	require.NoError(t, err, "Load should succeed with the example config")

	admin := cfg.Profiles["admin"]
	// Secrets must be environment variable names, not actual secret values.
	// They should look like ENV_VAR style names (uppercase, underscores).
	assert.Equal(t, "OPENEDX_ADMIN_CLIENT_ID", admin.ClientIDEnv,
		"client_id_env should be an env var name, not a literal secret")
	assert.Equal(t, "OPENEDX_ADMIN_CLIENT_SECRET", admin.ClientSecretEnv,
		"client_secret_env should be an env var name, not a literal secret")

	ops := cfg.Profiles["ops"]
	assert.Equal(t, "OPENEDX_OPS_CLIENT_ID", ops.ClientIDEnv)
	assert.Equal(t, "OPENEDX_OPS_CLIENT_SECRET", ops.ClientSecretEnv)

	// Demonstrate that resolution happens via os.Getenv at runtime.
	// Set the env var and verify we can look it up using the stored name.
	os.Setenv("OPENEDX_ADMIN_CLIENT_ID", "test-client-id-value")
	defer os.Unsetenv("OPENEDX_ADMIN_CLIENT_ID")

	assert.Equal(t, "test-client-id-value", os.Getenv(admin.ClientIDEnv),
		"the stored env var name should resolve via os.Getenv")
}

func TestLoadConfigValidation(t *testing.T) {
	t.Run("missing base_url", func(t *testing.T) {
		dir := t.TempDir()
		configFile := filepath.Join(dir, "openedx.yaml")
		content := []byte(`version: 1
profiles:
  broken:
    token_url: https://example.com/oauth2/token
    client_id_env: CLIENT_ID
    client_secret_env: CLIENT_SECRET
`)
		err := os.WriteFile(configFile, content, 0644)
		require.NoError(t, err)

		_, err = Load(configFile)
		assert.Error(t, err, "should reject profile with missing base_url")
		assert.Contains(t, err.Error(), "base_url")
		assert.Contains(t, err.Error(), "broken")
	})

	t.Run("missing token_url", func(t *testing.T) {
		dir := t.TempDir()
		configFile := filepath.Join(dir, "openedx.yaml")
		content := []byte(`version: 1
profiles:
  broken:
    base_url: https://example.com
    client_id_env: CLIENT_ID
    client_secret_env: CLIENT_SECRET
`)
		err := os.WriteFile(configFile, content, 0644)
		require.NoError(t, err)

		_, err = Load(configFile)
		assert.Error(t, err, "should reject profile with missing token_url")
		assert.Contains(t, err.Error(), "token_url")
		assert.Contains(t, err.Error(), "broken")
	})

	t.Run("config file not found", func(t *testing.T) {
		_, err := Load("/nonexistent/path/config.yaml")
		assert.Error(t, err, "should return error for missing config file")
		assert.Contains(t, err.Error(), "config")
	})

	t.Run("valid config with empty profiles", func(t *testing.T) {
		dir := t.TempDir()
		configFile := filepath.Join(dir, "openedx.yaml")
		content := []byte(`version: 1
profiles: {}
extensions: {}
`)
		err := os.WriteFile(configFile, content, 0644)
		require.NoError(t, err)

		cfg, err := Load(configFile)
		require.NoError(t, err, "empty profiles/extensions should be valid")
		assert.Empty(t, cfg.Profiles)
		assert.Empty(t, cfg.Extensions)
	})
}
