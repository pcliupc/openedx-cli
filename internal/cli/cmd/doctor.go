package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/openedx/cli/internal/diagnostics"
)

// DoctorCheckFunc is the signature for running diagnostic checks. The args
// parameter contains positional arguments after the "doctor" subcommand:
//   - empty: run all checks
//   - "verify <command>": verify a specific command mapping
type DoctorCheckFunc func(ctx context.Context, args []string) (*diagnostics.DoctorResult, error)

// NewDoctorCmd creates the "doctor" command for running health diagnostics.
func NewDoctorCmd(checkFn DoctorCheckFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor [verify <command>]",
		Short: "Run health diagnostics",
		Long: `Run diagnostic checks to verify that the CLI can reach the configured
Open edX deployment and acquire authentication tokens.

Use "doctor verify <command>" to check whether a specific command key has
a valid mapping in the public registry or extension configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if checkFn == nil {
				return fmt.Errorf("doctor check function not configured")
			}

			result, err := checkFn(cmd.Context(), args)
			if err != nil {
				return err
			}

			return printOutput(cmd, result)
		},
	}

	return cmd
}
