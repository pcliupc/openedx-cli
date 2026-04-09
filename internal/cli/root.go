package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/openedx/cli/internal/cli/cmd"
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

	// Register command groups.
	rootCmd.AddCommand(
		cmd.NewCourseCmd(DefaultExecuteFunc),
		cmd.NewUserCmd(DefaultExecuteFunc),
		cmd.NewEnrollmentCmd(DefaultExecuteFunc),
		cmd.NewRoleCmd(DefaultExecuteFunc),
	)

	return rootCmd
}

// DefaultExecuteFunc is the production implementation of ExecuteFunc.
// It loads configuration, acquires a token, looks up the command in the
// registry, and dispatches through the fallback provider.
// TODO: Wire this to the full config/auth/provider stack in a future task.
var DefaultExecuteFunc cmd.ExecuteFunc = func(_ context.Context, cmdKey string, args map[string]string) ([]byte, error) {
	return nil, fmt.Errorf("command execution not yet configured: %s", cmdKey)
}
