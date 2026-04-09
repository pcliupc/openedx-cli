package cli

import (
	"github.com/spf13/cobra"
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

	return rootCmd
}
