package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/openedx/cli/internal/normalize"
)

// NewUserCmd creates the "user" command group with all its subcommands.
func NewUserCmd(execFn ExecuteFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage Open edX users",
		Long:  "Create, list, and inspect users in an Open edX deployment.",
	}

	cmd.AddCommand(
		newUserCreateCmd(execFn),
	)

	return cmd
}

// --- user create ---

func newUserCreateCmd(execFn ExecuteFunc) *cobra.Command {
	var username, email, name string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		Long:  "Create a new user in the configured Open edX deployment.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{
				"username": username,
				"email":    email,
			}
			if name != "" {
				cmdArgs["name"] = name
			}

			data, err := execFn(cmd.Context(), "user.create", cmdArgs)
			if err != nil {
				return err
			}

			user, err := normalize.UserFromJSON(data)
			if err != nil {
				return fmt.Errorf("failed to normalize user: %w", err)
			}

			return printOutput(cmd, user)
		},
	}

	cmd.Flags().StringVar(&username, "username", "", "username for the new user (required)")
	cmd.Flags().StringVar(&email, "email", "", "email address for the new user (required)")
	cmd.Flags().StringVar(&name, "name", "", "full name for the new user")
	_ = cmd.MarkFlagRequired("username")
	_ = cmd.MarkFlagRequired("email")

	return cmd
}
