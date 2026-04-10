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
		newUserListCmd(execFn),
		newUserGetCmd(execFn),
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

// --- user list ---

func newUserListCmd(execFn ExecuteFunc) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users",
		Long:  "List user accounts in the configured Open edX deployment.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{}
			if page > 0 {
				cmdArgs["page"] = fmt.Sprintf("%d", page)
			}
			if pageSize > 0 {
				cmdArgs["page_size"] = fmt.Sprintf("%d", pageSize)
			}

			data, err := execFn(cmd.Context(), "user.list", cmdArgs)
			if err != nil {
				return err
			}

			users, err := normalize.UserListFromJSON(data)
			if err != nil {
				return fmt.Errorf("failed to normalize user list: %w", err)
			}

			return printOutput(cmd, users)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "number of results per page")

	return cmd
}

// --- user get ---

func newUserGetCmd(execFn ExecuteFunc) *cobra.Command {
	var username string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get user details",
		Long:  "Retrieve details for a specific user by username.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{
				"username": username,
			}

			data, err := execFn(cmd.Context(), "user.get", cmdArgs)
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

	cmd.Flags().StringVar(&username, "username", "", "username to look up (required)")
	_ = cmd.MarkFlagRequired("username")

	return cmd
}
