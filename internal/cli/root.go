// Package cli defines the bento-layout command tree.
package cli

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/EKaverinskyi/bento-layout/internal/version"
)

// NewRootCmd builds the root command. Output streams are injected for testing.
func NewRootCmd(stdout, stderr io.Writer) *cobra.Command {
	root := &cobra.Command{
		Use:           "bento-layout",
		Short:         "Create and check a standard folder layout for a portable drive",
		Version:       version.Version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetVersionTemplate("bento-layout {{.Version}}\n")
	root.AddCommand(newVersionCmd())
	return root
}

// Execute runs the CLI with args and returns the process exit code.
func Execute(args []string, stdout, stderr io.Writer) int {
	root := NewRootCmd(stdout, stderr)
	root.SetArgs(args)
	if err := root.Execute(); err != nil {
		_, _ = io.WriteString(stderr, "Error: "+err.Error()+"\n")
		return 1
	}
	return 0
}
