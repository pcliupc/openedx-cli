package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/openedx/cli/internal/config"
	"github.com/openedx/cli/internal/diagnostics"
)

// NewSchemaCmd creates the "schema" command for inspecting how CLI commands
// map to backend API endpoints and extension overrides.
func NewSchemaCmd(extensions map[string]config.ExtensionMapping) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema [command|all]",
		Short: "Show command schema and visibility",
		Long: `Display how a CLI command maps to its backend API endpoint.

Pass a command key (e.g. "course.create") to see a single command's schema,
or "all" to list schemas for every registered command.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			arg := args[0]

			if arg == "all" {
				schemas, err := diagnostics.GetAllCommandSchemas(extensions)
				if err != nil {
					return fmt.Errorf("failed to get command schemas: %w", err)
				}
				return printOutput(cmd, schemas)
			}

			schema, err := diagnostics.GetCommandSchema(arg, extensions)
			if err != nil {
				return err
			}
			return printOutput(cmd, schema)
		},
	}

	return cmd
}
