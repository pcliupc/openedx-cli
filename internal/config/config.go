package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Load reads and validates the CLI configuration. If configPath is non-empty it
// is used directly; otherwise the function searches the default paths
// (./openedx.yaml and ~/.openedx/config.yaml). Returns an error if the file
// cannot be found, parsed, or fails validation.
func Load(configPath string) (*Config, error) {
	// Use a key delimiter that will never appear in our config keys so that
	// Viper does not treat dots in extension names (e.g. "course.create")
	// as nested map separators.
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))

	v.SetConfigType("yaml")

	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// Default search path: local directory
		v.SetConfigName("openedx")
		v.AddConfigPath(".")

		// Default search path: user home directory
		home, err := os.UserHomeDir()
		if err == nil {
			v.AddConfigPath(filepath.Join(home, ".openedx"))
		}
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: searched ./openedx.yaml and ~/.openedx/config.yaml")
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validate checks that every profile has the required base_url and token_url
// fields.
func validate(cfg *Config) error {
	for name, profile := range cfg.Profiles {
		if strings.TrimSpace(profile.BaseURL) == "" {
			return fmt.Errorf("profile %q: base_url is required", name)
		}
		if strings.TrimSpace(profile.TokenURL) == "" {
			return fmt.Errorf("profile %q: token_url is required", name)
		}
	}
	return nil
}
