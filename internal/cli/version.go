package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/EKaverinskyi/bento-layout/internal/version"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the bento-layout version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "bento-layout %s\n", version.Version)
			return err
		},
	}
}
