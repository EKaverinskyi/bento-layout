package cli

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/EKaverinskyi/bento-layout/internal/layout"
)

// errChecksFailed makes the command exit non-zero after the report is printed.
var errChecksFailed = errors.New("layout check failed")

func newCheckCmd() *cobra.Command {
	var layoutsDir, rulesFile, layoutArg, drivePath string
	cmd := &cobra.Command{
		Use:   "check --layout <layout> --path <dir>",
		Short: "Check a layout definition and a real directory against it",
		Long: `Check that a layout is sound, then compare a real directory with it.

--layout is a path to a layout directory, or the name of one in --layouts-dir.
Its definition is checked against the universal rules (rules.toml in
--layouts-dir, or --rules) and its own rules and scenarios.

--path is the directory to compare with the layout, for example ~/bento. Missing
folders, and files where a folder belongs, are failures. Folders the layout does
not list are shown but do not fail the check. Nothing on disk is changed. The
directory is not compared if the layout itself fails.

The exit code is 1 if the layout is incomplete, a rule fails, or the directory
does not match.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if rulesFile == "" {
				rulesFile = filepath.Join(layoutsDir, "rules.toml")
			}
			universal, err := layout.LoadRules(rulesFile)
			if err != nil {
				return err
			}
			dir, err := resolveLayoutDir(layoutArg, layoutsDir)
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			rep := layout.Check(dir, universal)
			if err := printReport(out, rep); err != nil {
				return err
			}
			if rep.Failed() > 0 {
				if _, err := fmt.Fprintf(out, "-- %s: not checked (layout fails)\n", drivePath); err != nil {
					return err
				}
				return errChecksFailed
			}
			l, err := layout.Load(dir)
			if err != nil {
				return err
			}
			drep, err := layout.CheckDrive(l, drivePath)
			if err != nil {
				return err
			}
			if err := printDriveReport(out, l.Name, drep); err != nil {
				return err
			}
			if drep.Failed() {
				return errChecksFailed
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&layoutArg, "layout", "", "layout to check: a name in --layouts-dir or a path to a layout directory (required)")
	cmd.Flags().StringVar(&drivePath, "path", "", "directory to compare with the layout, e.g. ~/bento (required)")
	cmd.Flags().StringVar(&layoutsDir, "layouts-dir", "layouts", "directory searched for layouts given by name, and for the universal rules")
	cmd.Flags().StringVar(&rulesFile, "rules", "", "universal rules file (default: rules.toml in --layouts-dir)")
	_ = cmd.MarkFlagRequired("layout")
	_ = cmd.MarkFlagRequired("path")
	return cmd
}

func printReport(w io.Writer, rep *layout.Report) error {
	var err error
	p := func(format string, a ...any) {
		if err == nil {
			_, err = fmt.Fprintf(w, format+"\n", a...)
		}
	}
	if len(rep.Incomplete) > 0 {
		p("== %s: incomplete layout (%s)", rep.Name, rep.Dir)
		p("FAIL  files")
		for _, s := range rep.Incomplete {
			p("      - %s", s)
		}
		p("-- %s: not checked\n", rep.Name)
		return err
	}
	p("== %s: %s", rep.Name, rep.Description)
	if len(rep.Structure) > 0 {
		p("FAIL  structure")
		for _, s := range rep.Structure {
			p("      - %s", s)
		}
	}
	for _, o := range rep.Outcomes {
		if len(o.Problems) == 0 {
			p("ok    %s  %s", o.Rule.ID, o.Rule.Text)
			continue
		}
		p("FAIL  %s  %s", o.Rule.ID, o.Rule.Text)
		for _, s := range o.Problems {
			p("      - %s", s)
		}
	}
	p("-- %s: %d rules, %d failing\n", rep.Name, len(rep.Outcomes), rep.Failed())
	return err
}

func printDriveReport(w io.Writer, name string, r *layout.DriveReport) error {
	var err error
	p := func(format string, a ...any) {
		if err == nil {
			_, err = fmt.Fprintf(w, format+"\n", a...)
		}
	}
	p("== %s in %s", name, r.Root)
	for _, s := range r.Missing {
		p("MISSING   %s", s)
	}
	for _, s := range r.Conflicts {
		p("CONFLICT  %s  (exists and is not a directory)", s)
	}
	for _, s := range r.Extra {
		p("extra     %s  (not in the layout)", s)
	}
	status := "matches"
	if r.Failed() {
		status = "does not match"
	}
	p("-- %s %s the layout: %d present, %d missing, %d conflicts, %d extra\n",
		r.Root, status, r.Present, len(r.Missing), len(r.Conflicts), len(r.Extra))
	return err
}
