package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

// NewEnrollmentCmd creates the "enrollment" command group with all its subcommands.
func NewEnrollmentCmd(execFn ExecuteFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enrollment",
		Short: "Manage Open edX enrollments",
		Long:  "Add, remove, and inspect course enrollments in an Open edX deployment.",
	}

	cmd.AddCommand(
		newEnrollmentAddCmd(execFn),
	)

	return cmd
}

// --- enrollment add ---

func newEnrollmentAddCmd(execFn ExecuteFunc) *cobra.Command {
	var courseID, username, mode string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Enroll a user in a course",
		Long:  "Enroll a user in a course with the specified enrollment mode.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{
				"course_id": courseID,
				"username":  username,
				"mode":      mode,
			}

			data, err := execFn(cmd.Context(), "enrollment.add", cmdArgs)
			if err != nil {
				return err
			}

			// No dedicated normalizer for enrollment yet; output raw JSON.
			var raw json.RawMessage = data
			return printOutput(cmd, &raw)
		},
	}

	cmd.Flags().StringVar(&courseID, "course-id", "", "course identifier (required)")
	cmd.Flags().StringVar(&username, "username", "", "username to enroll (required)")
	cmd.Flags().StringVar(&mode, "mode", "audit", "enrollment mode (e.g. audit, verified)")
	_ = cmd.MarkFlagRequired("course-id")
	_ = cmd.MarkFlagRequired("username")

	return cmd
}
