package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/EKaverinskyi/bento-layout/internal/layout"
)

func newInitCmd() *cobra.Command {
	var layoutsDir string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "init <layout> [target]",
		Short: "Create the folder structure of a layout",
		Long: `Create the folder structure of a layout below target (default: the current directory).

<layout> is either a path to a layout directory (it contains layout.toml and
folders.toml), or a layout name that is looked up in --layouts-dir.

Existing folders are left untouched; nothing is ever deleted or renamed.`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 2 && args[1] == "" {
				return errors.New("target is empty; omit it to use the current directory")
			}
			dir, err := resolveLayoutDir(args[0], layoutsDir)
			if err != nil {
				return err
			}
			l, err := layout.Load(dir)
			if err != nil {
				return err
			}
			target := "."
			if len(args) == 2 {
				target = args[1]
			}
			res, err := layout.Create(l, target, dryRun)
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			verb := "created"
			if dryRun {
				verb = "would create"
			}
			for _, p := range res.Created {
				if _, err := fmt.Fprintf(out, "%s %s\n", verb, relToCwd(filepath.Join(target, filepath.FromSlash(p)))); err != nil {
					return err
				}
			}
			_, err = fmt.Fprintf(out, "layout %q in %s: %s %d, %d already existed\n",
				l.Name, target, verb, len(res.Created), len(res.Existed))
			return err
		},
	}
	cmd.Flags().StringVar(&layoutsDir, "layouts-dir", "layouts", "directory searched for layouts given by name")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "show what would be created without changing anything")
	return cmd
}

// resolveLayoutDir turns the <layout> argument into a directory. An argument
// containing a path separator, or naming an existing directory in the working
// directory, is used as a path; anything else is a name inside layoutsDir.
func resolveLayoutDir(arg, layoutsDir string) (string, error) {
	if strings.ContainsRune(arg, filepath.Separator) || arg == "." || arg == ".." {
		return arg, nil
	}
	byName := filepath.Join(layoutsDir, arg)
	if st, err := os.Stat(byName); err == nil && st.IsDir() {
		return byName, nil
	}
	if st, err := os.Stat(arg); err == nil && st.IsDir() {
		return arg, nil
	}
	return "", fmt.Errorf("layout %q not found (looked in %s)", arg, layoutsDir)
}

// relToCwd returns p relative to the current directory. If p is outside the
// current directory (the result would start with ".."), or on error, p is
// returned unchanged.
func relToCwd(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	cwd, err := os.Getwd()
	if err != nil {
		return p
	}
	rel, err := filepath.Rel(cwd, abs)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return p
	}
	return rel
}
