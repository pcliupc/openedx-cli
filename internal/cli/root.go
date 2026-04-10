package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/openedx/cli/internal/auth"
	"github.com/openedx/cli/internal/cli/cmd"
	"github.com/openedx/cli/internal/config"
	"github.com/openedx/cli/internal/diagnostics"
	"github.com/openedx/cli/internal/provider"
	"github.com/openedx/cli/internal/registry"
)

// NewRootCmd creates and returns the root command for the openedx CLI.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "openedx",
		Short: "Open edX CLI for CI pipelines and automation",
		Long:  "A Go-based CLI for Open edX designed for CI pipelines and coding agents.",
	}

	rootCmd.PersistentFlags().StringP("profile", "p", "", "config profile to use")
	rootCmd.PersistentFlags().StringP("format", "f", "json", "output format (json, table)")
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file path")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	// loadCfg reads config based on the --config flag value.
	loadCfg := func() (*config.Config, error) {
		configPath, _ := rootCmd.PersistentFlags().GetString("config")
		return config.Load(configPath)
	}

	// resolveProfile loads config and returns the named profile.
	resolveProfile := func() (*config.Config, *config.Profile, error) {
		profileName, _ := rootCmd.PersistentFlags().GetString("profile")
		cfg, err := loadCfg()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load config: %w", err)
		}
		p, ok := cfg.Profiles[profileName]
		if !ok {
			return nil, nil, fmt.Errorf("profile %q not found in config", profileName)
		}
		return cfg, &p, nil
	}

	// execFn is the production ExecuteFunc: loads config, acquires a token,
	// looks up the command in the registry, and dispatches through the
	// fallback provider.
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		cfg, profile, err := resolveProfile()
		if err != nil {
			return nil, err
		}

		reg := registry.LatestRegistry()
		cmdMeta, ok := reg[cmdKey]
		if !ok {
			return nil, fmt.Errorf("unknown command: %s", cmdKey)
		}

		var extMapping *config.ExtensionMapping
		if ext, found := cfg.Extensions[cmdKey]; found {
			extMapping = &ext
		}

		tokenProvider := auth.NewCachingTokenProvider(auth.NewHTTPTokenClient(nil), nil)
		fb := provider.NewFallbackProvider(nil)

		result, err := fb.Execute(ctx, *profile, tokenProvider, cmdKey, cmdMeta, extMapping, args)
		if err != nil {
			return nil, err
		}

		return result.Body, nil
	}

	// extProvider returns extension mappings from the loaded config.
	extProvider := func() map[string]config.ExtensionMapping {
		cfg, err := loadCfg()
		if err != nil {
			return nil
		}
		return cfg.Extensions
	}

	// doctorFn is the production DoctorCheckFunc.
	doctorFn := func(ctx context.Context, args []string) (*diagnostics.DoctorResult, error) {
		cfg, profile, err := resolveProfile()
		if err != nil {
			return nil, err
		}

		tokenProvider := auth.NewCachingTokenProvider(auth.NewHTTPTokenClient(nil), nil)

		// "doctor verify <command>" checks a specific command mapping.
		if len(args) >= 2 && args[0] == "verify" {
			check := diagnostics.CheckCommand(ctx, args[1], *profile, cfg.Extensions)
			return &diagnostics.DoctorResult{Checks: []diagnostics.CheckResult{check}}, nil
		}

		return diagnostics.RunAllChecks(ctx, *profile, tokenProvider, cfg.Extensions), nil
	}

	// Register command groups.
	rootCmd.AddCommand(
		cmd.NewCourseCmd(execFn),
		cmd.NewUserCmd(execFn),
		cmd.NewEnrollmentCmd(execFn),
		cmd.NewRoleCmd(execFn),
		cmd.NewGradeCmd(execFn),
		cmd.NewCertificateCmd(execFn),
		cmd.NewSchemaCmd(extProvider),
		cmd.NewDoctorCmd(doctorFn),
	)

	return rootCmd
}
