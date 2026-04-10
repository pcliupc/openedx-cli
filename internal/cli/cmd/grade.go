package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/openedx/cli/internal/normalize"
)

// NewGradeCmd creates the "grade" command group with all its subcommands.
func NewGradeCmd(execFn ExecuteFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grade",
		Short: "Manage course grades",
		Long:  "List grades and view gradebooks for Open edX courses.",
	}

	cmd.AddCommand(
		newGradeListCmd(execFn),
		newGradebookCmd(execFn),
	)

	return cmd
}

// --- grade list ---

func newGradeListCmd(execFn ExecuteFunc) *cobra.Command {
	var courseID, username string
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List grades for a course",
		Long:  "List student grades for a specific course with optional filters.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{
				"course_id": courseID,
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

			data, err := execFn(cmd.Context(), "grade.list", cmdArgs)
			if err != nil {
				return err
			}

			grades, err := normalize.GradeListFromJSON(data)
			if err != nil {
				return fmt.Errorf("failed to normalize grade list: %w", err)
			}

			return printOutput(cmd, grades)
		},
	}

	cmd.Flags().StringVar(&courseID, "course-id", "", "course identifier (required)")
	cmd.Flags().StringVar(&username, "username", "", "filter by username")
	cmd.Flags().IntVar(&page, "page", 0, "page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "number of results per page")
	_ = cmd.MarkFlagRequired("course-id")

	return cmd
}

// --- grade gradebook ---

func newGradebookCmd(execFn ExecuteFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gradebook",
		Short: "Gradebook operations",
		Long:  "View course gradebooks.",
	}

	cmd.AddCommand(newGradebookGetCmd(execFn))
	return cmd
}

func newGradebookGetCmd(execFn ExecuteFunc) *cobra.Command {
	var courseID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get course gradebook",
		Long:  "Retrieve the full gradebook for a specific course.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{
				"course_id": courseID,
			}

			data, err := execFn(cmd.Context(), "gradebook.get", cmdArgs)
			if err != nil {
				return err
			}

			gradebook, err := normalize.GradebookFromJSON(data)
			if err != nil {
				return fmt.Errorf("failed to normalize gradebook: %w", err)
			}

			return printOutput(cmd, gradebook)
		},
	}

	cmd.Flags().StringVar(&courseID, "course-id", "", "course identifier (required)")
	_ = cmd.MarkFlagRequired("course-id")

	return cmd
}
