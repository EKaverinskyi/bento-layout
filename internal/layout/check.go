package layout

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RequiredFiles must all exist in a layout directory before it is checked.
var RequiredFiles = []string{metaFile, foldersFile, rulesFile, scenariosFile}

// Outcome is the result of one rule.
type Outcome struct {
	Rule     Rule
	Problems []string // empty means the rule passed
}

// Report is the result of checking one layout directory.
type Report struct {
	Name        string
	Dir         string
	Description string
	// Incomplete lists why the layout could not be checked at all (missing
	// files, unreadable TOML, empty metadata). When set, nothing else is.
	Incomplete []string
	// Structure lists problems with the folder tree itself.
	Structure []string
	Outcomes  []Outcome
}

// Failed counts the failing checks: the structure check (if it found
// anything), each failing rule, or 1 for an incomplete layout.
func (r *Report) Failed() int {
	n := 0
	if len(r.Incomplete) > 0 {
		return 1
	}
	if len(r.Structure) > 0 {
		n++
	}
	for _, o := range r.Outcomes {
		if len(o.Problems) > 0 {
			n++
		}
	}
	return n
}

// Check verifies the layout in dir against the universal rules followed by the
// layout's own rules. It never returns an error: everything wrong with the
// layout is reported in the Report.
func Check(dir string, universal []Rule) *Report {
	abs, err := filepath.Abs(dir)
	if err != nil {
		abs = dir
	}
	rep := &Report{Name: filepath.Base(abs), Dir: dir}

	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		rep.Incomplete = []string{"not a directory: " + dir}
		return rep
	}
	for _, name := range RequiredFiles {
		if st, err := os.Stat(filepath.Join(dir, name)); err != nil || !st.Mode().IsRegular() {
			rep.Incomplete = append(rep.Incomplete, "missing required file: "+name)
		}
	}
	if len(rep.Incomplete) > 0 {
		return rep
	}

	l, err := Parse(dir)
	if err != nil {
		rep.Incomplete = []string{err.Error()}
		return rep
	}
	if strings.TrimSpace(l.Name) == "" {
		rep.Incomplete = append(rep.Incomplete, "layout.toml: [layout] needs a non-empty name")
	} else if l.Name != rep.Name {
		rep.Incomplete = append(rep.Incomplete, fmt.Sprintf("layout.toml: name %q does not match directory %q", l.Name, rep.Name))
	}
	if strings.TrimSpace(l.Description) == "" {
		rep.Incomplete = append(rep.Incomplete, "layout.toml: [layout] needs a non-empty description")
	}
	if len(l.Folders) == 0 {
		rep.Incomplete = append(rep.Incomplete, "folders.toml: defines no [[folder]]")
	}
	if len(l.Scenarios) == 0 {
		rep.Incomplete = append(rep.Incomplete, "scenarios.toml: defines no [[scenario]]")
	}
	if len(rep.Incomplete) > 0 {
		return rep
	}

	rep.Description = l.Description
	rep.Structure = l.structureProblems()
	t := newTree(l)
	for _, r := range append(append([]Rule{}, universal...), l.Rules...) {
		rep.Outcomes = append(rep.Outcomes, Outcome{Rule: r, Problems: evalRule(t, r)})
	}
	return rep
}

// evalRule runs one rule. A malformed rule counts as a failure.
func evalRule(t *tree, r Rule) []string {
	if r.ID == "" || r.Type == "" || r.Text == "" {
		return []string{"rule needs id, type and text"}
	}
	fn, ok := ruleTypes[r.Type]
	if !ok {
		return []string{"unknown rule type: " + r.Type}
	}
	problems, err := fn(t, &params{m: r.Params})
	if err != nil {
		return []string{"invalid rule: " + err.Error()}
	}
	return problems
}
