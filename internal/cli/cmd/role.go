package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

// NewRoleCmd creates the "role" command group with all its subcommands.
func NewRoleCmd(execFn ExecuteFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role",
		Short: "Manage Open edX roles",
		Long:  "Assign, revoke, and inspect user roles in an Open edX deployment.",
	}

	cmd.AddCommand(
		newRoleAssignCmd(execFn),
	)

	return cmd
}

// --- role assign ---

func newRoleAssignCmd(execFn ExecuteFunc) *cobra.Command {
	var courseID, username, role string

	cmd := &cobra.Command{
		Use:   "assign",
		Short: "Assign a role to a user",
		Long:  "Assign a role to a user on a specific course.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{
				"course_id": courseID,
				"username":  username,
				"role":      role,
			}

			data, err := execFn(cmd.Context(), "role.assign", cmdArgs)
			if err != nil {
				return err
			}

			// No dedicated normalizer for role yet; output raw JSON.
			var raw json.RawMessage = data
			return printOutput(cmd, &raw)
		},
	}

	cmd.Flags().StringVar(&courseID, "course-id", "", "course identifier (required)")
	cmd.Flags().StringVar(&username, "username", "", "username to assign the role to (required)")
	cmd.Flags().StringVar(&role, "role", "", "role to assign (e.g. instructor, staff) (required)")
	_ = cmd.MarkFlagRequired("course-id")
	_ = cmd.MarkFlagRequired("username")
	_ = cmd.MarkFlagRequired("role")

	return cmd
}
