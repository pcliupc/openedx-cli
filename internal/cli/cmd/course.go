// Package cmd implements Cobra subcommands for the Open edX CLI.
// Each command group (course, user, enrollment, etc.) is defined in its own
// file and wires flags to a provider-backed executor function.
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/openedx/cli/internal/model"
	"github.com/openedx/cli/internal/normalize"
)

// ExecuteFunc is the signature for executing a command against a provider.
// In production this loads config, acquires a token, looks up the command in
// the registry, calls the fallback provider, and returns raw response bytes.
// In tests this can be replaced with a mock that returns fixture data.
type ExecuteFunc func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error)

// NewCourseCmd creates the "course" command group with all its subcommands.
func NewCourseCmd(execFn ExecuteFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "course",
		Short: "Manage Open edX courses",
		Long:  "List, create, import, export, and inspect courses in an Open edX deployment.",
	}

	cmd.AddCommand(
		newCourseListCmd(execFn),
		newCourseGetCmd(execFn),
		newCourseCreateCmd(execFn),
		newCourseImportCmd(execFn),
		newCourseExportCmd(execFn),
		newCourseOutlineCmd(execFn),
	)

	return cmd
}

// --- course list ---

func newCourseListCmd(execFn ExecuteFunc) *cobra.Command {
	var page, pageSize int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List courses",
		Long:  "List courses available in the configured Open edX deployment.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{}
			if all {
				cmdArgs["all"] = "true"
			} else {
				if page > 0 {
					cmdArgs["page"] = fmt.Sprintf("%d", page)
				}
				if pageSize > 0 {
					cmdArgs["page_size"] = fmt.Sprintf("%d", pageSize)
				}
			}

			data, err := execFn(cmd.Context(), "course.list", cmdArgs)
			if err != nil {
				return err
			}

			courses, err := normalize.CourseListFromJSON(data)
			if err != nil {
				return fmt.Errorf("failed to normalize course list: %w", err)
			}

			return printOutput(cmd, courses)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "number of results per page")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all results")

	return cmd
}

// --- course get ---

func newCourseGetCmd(execFn ExecuteFunc) *cobra.Command {
	var courseID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get course details",
		Long:  "Retrieve details for a specific course by its course ID.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID == "" {
				return fmt.Errorf("required flag(s) \"course-id\" not set")
			}

			cmdArgs := map[string]string{
				"course_id": courseID,
			}

			data, err := execFn(cmd.Context(), "course.get", cmdArgs)
			if err != nil {
				return err
			}

			course, err := normalize.CourseFromJSON(data)
			if err != nil {
				return fmt.Errorf("failed to normalize course: %w", err)
			}

			return printOutput(cmd, course)
		},
	}

	cmd.Flags().StringVar(&courseID, "course-id", "", "course identifier (required)")
	_ = cmd.MarkFlagRequired("course-id")

	return cmd
}

// --- course create ---

func newCourseCreateCmd(execFn ExecuteFunc) *cobra.Command {
	var org, number, run, title string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new course",
		Long:  "Create a new course in the configured Open edX deployment.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{
				"org":    org,
				"number": number,
				"run":    run,
				"title":  title,
			}

			data, err := execFn(cmd.Context(), "course.create", cmdArgs)
			if err != nil {
				return err
			}

			course, err := normalize.CourseFromJSON(data)
			if err != nil {
				return fmt.Errorf("failed to normalize course: %w", err)
			}

			return printOutput(cmd, course)
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "organization (required)")
	cmd.Flags().StringVar(&number, "number", "", "course number (required)")
	cmd.Flags().StringVar(&run, "run", "", "course run (required)")
	cmd.Flags().StringVar(&title, "title", "", "course title (required)")
	_ = cmd.MarkFlagRequired("org")
	_ = cmd.MarkFlagRequired("number")
	_ = cmd.MarkFlagRequired("run")
	_ = cmd.MarkFlagRequired("title")

	return cmd
}

// --- course import ---

func newCourseImportCmd(execFn ExecuteFunc) *cobra.Command {
	var courseID, file string

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import a course",
		Long:  "Import a course from a tar.gz archive file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID == "" {
				return fmt.Errorf("required flag(s) \"course-id\" not set")
			}
			if file == "" {
				return fmt.Errorf("required flag(s) \"file\" not set")
			}

			cmdArgs := map[string]string{
				"course_id": courseID,
				"file":      file,
			}

			data, err := execFn(cmd.Context(), "course.import", cmdArgs)
			if err != nil {
				return err
			}

			job, err := normalizeJob(data)
			if err != nil {
				return fmt.Errorf("failed to normalize job: %w", err)
			}

			return printOutput(cmd, job)
		},
	}

	cmd.Flags().StringVar(&courseID, "course-id", "", "course identifier (required)")
	cmd.Flags().StringVar(&file, "file", "", "path to course archive file (required)")
	_ = cmd.MarkFlagRequired("course-id")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

// --- course export ---

func newCourseExportCmd(execFn ExecuteFunc) *cobra.Command {
	var courseID, output string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export a course",
		Long:  "Export a course to a tar.gz archive.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID == "" {
				return fmt.Errorf("required flag(s) \"course-id\" not set")
			}

			cmdArgs := map[string]string{
				"course_id": courseID,
			}
			if output != "" {
				cmdArgs["output"] = output
			}

			data, err := execFn(cmd.Context(), "course.export", cmdArgs)
			if err != nil {
				return err
			}

			job, err := normalizeJob(data)
			if err != nil {
				return fmt.Errorf("failed to normalize job: %w", err)
			}

			return printOutput(cmd, job)
		},
	}

	cmd.Flags().StringVar(&courseID, "course-id", "", "course identifier (required)")
	cmd.Flags().StringVar(&output, "output", "", "output directory for the exported archive")
	_ = cmd.MarkFlagRequired("course-id")

	return cmd
}

// --- course outline ---

func newCourseOutlineCmd(execFn ExecuteFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "outline",
		Short: "Course outline operations",
		Long:  "Inspect and manage course outline structures.",
	}

	cmd.AddCommand(newCourseOutlineGetCmd(execFn))
	return cmd
}

func newCourseOutlineGetCmd(execFn ExecuteFunc) *cobra.Command {
	var courseID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get course outline",
		Long:  "Retrieve the outline structure for a specific course.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if courseID == "" {
				return fmt.Errorf("required flag(s) \"course-id\" not set")
			}

			cmdArgs := map[string]string{
				"course_id": courseID,
			}

			data, err := execFn(cmd.Context(), "course.outline.get", cmdArgs)
			if err != nil {
				return err
			}

			outline, err := normalize.OutlineFromJSON(data)
			if err != nil {
				return fmt.Errorf("failed to normalize outline: %w", err)
			}

			return printOutput(cmd, outline)
		},
	}

	cmd.Flags().StringVar(&courseID, "course-id", "", "course identifier (required)")
	_ = cmd.MarkFlagRequired("course-id")

	return cmd
}

// --- helpers ---

// normalizeJob parses raw JSON bytes into a model.Job.
func normalizeJob(data []byte) (*model.Job, error) {
	var job model.Job
	if err := json.Unmarshal(data, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// printOutput dispatches to JSON output using the format flag from the command.
func printOutput(cmd *cobra.Command, v interface{}) error {
	format, _ := cmd.Root().PersistentFlags().GetString("format")
	return printJSON(cmd.OutOrStdout(), format, v)
}

// printJSON writes the value as formatted JSON to the writer. It accepts a
// format parameter for forward compatibility but currently always outputs JSON.
func printJSON(w io.Writer, format string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\n"))
	return err
}
