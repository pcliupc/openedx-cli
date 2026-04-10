package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/openedx/cli/internal/normalize"
)

// NewCertificateCmd creates the "certificate" command group with all its subcommands.
func NewCertificateCmd(execFn ExecuteFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificate",
		Short: "Manage course certificates",
		Long:  "List and inspect course completion certificates.",
	}

	cmd.AddCommand(
		newCertificateListCmd(execFn),
	)

	return cmd
}

// --- certificate list ---

func newCertificateListCmd(execFn ExecuteFunc) *cobra.Command {
	var username string
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List certificates",
		Long:  "List certificates for a specific user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{
				"username": username,
			}
			if page > 0 {
				cmdArgs["page"] = fmt.Sprintf("%d", page)
			}
			if pageSize > 0 {
				cmdArgs["page_size"] = fmt.Sprintf("%d", pageSize)
			}

			data, err := execFn(cmd.Context(), "certificate.list", cmdArgs)
			if err != nil {
				return err
			}

			certs, err := normalize.CertificateListFromJSON(data)
			if err != nil {
				return fmt.Errorf("failed to normalize certificate list: %w", err)
			}

			return printOutput(cmd, certs)
		},
	}

	cmd.Flags().StringVar(&username, "username", "", "username to look up certificates for (required)")
	cmd.Flags().IntVar(&page, "page", 0, "page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "number of results per page")
	_ = cmd.MarkFlagRequired("username")

	return cmd
}
