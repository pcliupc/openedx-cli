package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/openedx/cli/internal/normalize"
)

// NewEnrollmentCmd creates the "enrollment" command group with all its subcommands.
func NewEnrollmentCmd(execFn ExecuteFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enrollment",
		Short: "Manage Open edX enrollments",
		Long:  "Add, remove, list, and inspect course enrollments in an Open edX deployment.",
	}

	cmd.AddCommand(
		newEnrollmentAddCmd(execFn),
		newEnrollmentListCmd(execFn),
		newEnrollmentRemoveCmd(execFn),
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

// --- enrollment list ---

func newEnrollmentListCmd(execFn ExecuteFunc) *cobra.Command {
	var courseID, username string
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List enrollments",
		Long:  "List course enrollments with optional filters.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{}
			if courseID != "" {
				cmdArgs["course_id"] = courseID
			}
			if username != "" {
				cmdArgs["username"] = username
			}
			if page > 0 {
				cmdArgs["page"] = fmt.Sprintf("%d", page)
			}
			if pageSize > 0 {
				cmdArgs["page_size"] = fmt.Sprintf("%d", pageSize)
			}

			data, err := execFn(cmd.Context(), "enrollment.list", cmdArgs)
			if err != nil {
				return err
			}

			enrollments, err := normalize.EnrollmentListFromJSON(data)
			if err != nil {
				return fmt.Errorf("failed to normalize enrollment list: %w", err)
			}

			return printOutput(cmd, enrollments)
		},
	}

	cmd.Flags().StringVar(&courseID, "course-id", "", "filter by course ID")
	cmd.Flags().StringVar(&username, "username", "", "filter by username")
	cmd.Flags().IntVar(&page, "page", 0, "page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "number of results per page")

	return cmd
}

// --- enrollment remove ---

func newEnrollmentRemoveCmd(execFn ExecuteFunc) *cobra.Command {
	var courseID, username string

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a user from a course",
		Long:  "Deactivate a user's enrollment in a course.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{
				"course_id": courseID,
				"username":  username,
			}

			data, err := execFn(cmd.Context(), "enrollment.remove", cmdArgs)
			if err != nil {
				return err
			}

			var raw json.RawMessage = data
			return printOutput(cmd, &raw)
		},
	}

	cmd.Flags().StringVar(&courseID, "course-id", "", "course identifier (required)")
	cmd.Flags().StringVar(&username, "username", "", "username to remove (required)")
	_ = cmd.MarkFlagRequired("course-id")
	_ = cmd.MarkFlagRequired("username")

	return cmd
}
