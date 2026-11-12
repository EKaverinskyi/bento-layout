package layout

import (
	"fmt"
	"path"
	"regexp"
	"slices"
	"sort"
	"strings"
)

// tree is a read-only index over a layout's folders that rules query.
type tree struct {
	l      *Layout
	byPath map[string]Folder
	paths  []string // unique folder paths, in file order
	top    []string // sorted unique top-level names
}

func newTree(l *Layout) *tree {
	t := &tree{l: l, byPath: make(map[string]Folder, len(l.Folders))}
	topSet := map[string]bool{}
	for _, f := range l.Folders {
		if _, dup := t.byPath[f.Path]; !dup {
			t.paths = append(t.paths, f.Path)
		}
		t.byPath[f.Path] = f
		topSet[strings.SplitN(f.Path, "/", 2)[0]] = true
	}
	for n := range topSet {
		t.top = append(t.top, n)
	}
	sort.Strings(t.top)
	return t
}

// children returns the sorted names of the direct subfolders of parent.
func (t *tree) children(parent string) []string {
	var out []string
	for _, p := range t.paths {
		if rest, ok := strings.CutPrefix(p, parent+"/"); ok && !strings.Contains(rest, "/") {
			out = append(out, rest)
		}
	}
	sort.Strings(out)
	return out
}

// flagged returns the sorted paths of the folders for which pick is true.
func (t *tree) flagged(pick func(Folder) bool) []string {
	var out []string
	for _, p := range t.paths {
		if pick(t.byPath[p]) {
			out = append(out, p)
		}
	}
	sort.Strings(out)
	return out
}

// params gives typed access to a rule's free-form parameters. The first
// missing or mistyped parameter is recorded in err and zero values are
// returned from then on, so a rule can read all its parameters and check err
// once.
type params struct {
	m   map[string]any
	err error
}

func (p *params) fail(key, want string) {
	if p.err == nil {
		p.err = fmt.Errorf("parameter %q must be %s", key, want)
	}
}

func (p *params) str(key string) string {
	v, ok := p.m[key].(string)
	if !ok {
		p.fail(key, "a string")
	}
	return v
}

func (p *params) int(key string) int {
	switch v := p.m[key].(type) {
	case int64:
		return int(v)
	case int:
		return v
	}
	p.fail(key, "an integer")
	return 0
}

func (p *params) strs(key string) []string {
	raw, ok := p.m[key].([]any)
	if !ok {
		p.fail(key, "a list of strings")
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		s, ok := v.(string)
		if !ok {
			p.fail(key, "a list of strings")
			return nil
		}
		out = append(out, s)
	}
	return out
}

// pattern compiles a regular expression that must match a whole string.
func (p *params) pattern(key string) *regexp.Regexp {
	s := p.str(key)
	re, err := regexp.Compile(`^(?:` + s + `)$`)
	if err != nil {
		if p.err == nil {
			p.err = fmt.Errorf("parameter %q: %w", key, err)
		}
		return nil
	}
	return re
}

// checkFunc evaluates one rule against a tree and returns the problems found
// (none means the rule passes). A non-nil error means the rule itself is
// malformed.
type checkFunc func(t *tree, p *params) ([]string, error)

// ruleTypes lists every rule type that can appear in a rules file.
var ruleTypes = map[string]checkFunc{
	"max_top_level":          checkMaxTopLevel,
	"top_level_name":         checkTopLevelName,
	"name":                   checkName,
	"max_depth":              checkMaxDepth,
	"has_description":        checkHasDescription,
	"forbidden_names":        checkForbiddenNames,
	"on_path_exact":          checkOnPathExact,
	"exact_children":         checkExactChildren,
	"all_exist":              checkAllExist,
	"manifest_exclude_exact": checkManifestExcludeExact,
	"no_children":            checkNoChildren,
	"scenarios_resolve":      checkScenariosResolve,
	"sibling_prefix":         checkSiblingPrefix,
	"no_dead_top_level":      checkNoDeadTopLevel,
}

// list formats a string slice as [a, b].
func list(ss []string) string {
	return "[" + strings.Join(ss, ", ") + "]"
}

func checkMaxTopLevel(t *tree, p *params) ([]string, error) {
	limit := p.int("value")
	if p.err != nil {
		return nil, p.err
	}
	if len(t.top) > limit {
		return []string{fmt.Sprintf("%d top-level folders: %s", len(t.top), list(t.top))}, nil
	}
	return nil, nil
}

func checkTopLevelName(t *tree, p *params) ([]string, error) {
	re := p.pattern("pattern")
	if p.err != nil {
		return nil, p.err
	}
	var out []string
	for _, n := range t.top {
		if !re.MatchString(n) {
			out = append(out, "bad top-level name: "+n)
		}
	}
	return out, nil
}

func checkName(t *tree, p *params) ([]string, error) {
	re := p.pattern("pattern")
	if p.err != nil {
		return nil, p.err
	}
	var out []string
	for _, fp := range t.paths {
		for seg := range strings.SplitSeq(fp, "/") {
			if !re.MatchString(seg) {
				out = append(out, "bad name: "+fp)
			}
		}
	}
	return out, nil
}

func checkMaxDepth(t *tree, p *params) ([]string, error) {
	limit := p.int("value")
	if p.err != nil {
		return nil, p.err
	}
	var out []string
	for _, fp := range t.paths {
		if strings.Count(fp, "/")+1 > limit {
			out = append(out, "too deep: "+fp)
		}
	}
	return out, nil
}

func checkHasDescription(t *tree, _ *params) ([]string, error) {
	var out []string
	for _, fp := range t.paths {
		if strings.TrimSpace(t.byPath[fp].Description) == "" {
			out = append(out, "no description: "+fp)
		}
	}
	return out, nil
}

func checkForbiddenNames(t *tree, p *params) ([]string, error) {
	names := p.strs("names")
	if p.err != nil {
		return nil, p.err
	}
	var out []string
	for _, n := range t.top {
		if slices.Contains(names, n) {
			out = append(out, "forbidden top-level name: "+n)
		}
	}
	return out, nil
}

// compareSets reports when the actual set differs from the expected one.
func compareSets(what string, actual, expected []string) []string {
	want := slices.Sorted(slices.Values(expected))
	if slices.Equal(actual, want) {
		return nil
	}
	return []string{fmt.Sprintf("%s %s, expected %s", what, list(actual), list(want))}
}

func checkOnPathExact(t *tree, p *params) ([]string, error) {
	paths := p.strs("paths")
	if p.err != nil {
		return nil, p.err
	}
	actual := t.flagged(func(f Folder) bool { return f.OnPath })
	return compareSets("on_path is", actual, paths), nil
}

func checkManifestExcludeExact(t *tree, p *params) ([]string, error) {
	paths := p.strs("paths")
	if p.err != nil {
		return nil, p.err
	}
	actual := t.flagged(func(f Folder) bool { return f.ManifestExclude })
	return compareSets("manifest_exclude is", actual, paths), nil
}

func checkExactChildren(t *tree, p *params) ([]string, error) {
	parent, children := p.str("parent"), p.strs("children")
	if p.err != nil {
		return nil, p.err
	}
	return compareSets(parent+" has", t.children(parent), children), nil
}

func checkAllExist(t *tree, p *params) ([]string, error) {
	paths := p.strs("paths")
	if p.err != nil {
		return nil, p.err
	}
	var out []string
	for _, fp := range paths {
		if _, ok := t.byPath[fp]; !ok {
			out = append(out, "missing: "+fp)
		}
	}
	return out, nil
}

func checkNoChildren(t *tree, p *params) ([]string, error) {
	fp := p.str("path")
	if p.err != nil {
		return nil, p.err
	}
	if c := t.children(fp); len(c) > 0 {
		return []string{fmt.Sprintf("%s has children %s", fp, list(c))}, nil
	}
	return nil, nil
}

func checkScenariosResolve(t *tree, _ *params) ([]string, error) {
	var out []string
	if len(t.l.Scenarios) == 0 {
		out = append(out, "layout has no scenarios")
	}
	seen := map[string]string{}
	for _, s := range t.l.Scenarios {
		if _, ok := t.byPath[s.Home]; !ok {
			out = append(out, fmt.Sprintf("home does not exist: %s (%s)", s.Home, s.Item))
		}
		if home, ok := seen[s.Item]; ok && home != s.Home {
			out = append(out, "two homes for: "+s.Item)
		}
		seen[s.Item] = s.Home
	}
	return out, nil
}

func checkSiblingPrefix(t *tree, p *params) ([]string, error) {
	maxShared := p.int("max_shared")
	if p.err != nil {
		return nil, p.err
	}
	groups := map[string][]string{}
	for _, fp := range t.paths {
		parent := path.Dir(fp)
		if parent == "." {
			parent = ""
		}
		groups[parent] = append(groups[parent], path.Base(fp))
	}
	parents := make([]string, 0, len(groups))
	for parent := range groups {
		parents = append(parents, parent)
	}
	sort.Strings(parents)
	var out []string
	for _, parent := range parents {
		names := groups[parent]
		sort.Strings(names)
		for i, a := range names {
			for _, b := range names[i+1:] {
				if k := commonPrefixLen(a, b); k > maxShared {
					shown := parent
					if shown == "" {
						shown = "/"
					}
					out = append(out, fmt.Sprintf("%s: %s and %s share the first %d characters", shown, a, b, k))
				}
			}
		}
	}
	return out, nil
}

func commonPrefixLen(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	k := 0
	for k < len(ra) && k < len(rb) && ra[k] == rb[k] {
		k++
	}
	return k
}

func checkNoDeadTopLevel(t *tree, _ *params) ([]string, error) {
	used := map[string]bool{}
	for _, s := range t.l.Scenarios {
		used[strings.SplitN(s.Home, "/", 2)[0]] = true
	}
	var out []string
	for _, n := range t.top {
		if !used[n] {
			out = append(out, "no scenario uses top-level folder: "+n)
		}
	}
	return out, nil
}
