package layout

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// fixture is a small tree that every rule type can be run against.
func fixture() *Layout {
	return &Layout{
		Name: "fx",
		Folders: []Folder{
			{Path: "dev", Description: "Code."},
			{Path: "dev/work", Description: "Work code."},
			{Path: "dev/personal", Description: "Own code."},
			{Path: "programs", Description: "Apps."},
			{Path: "programs/bin", Description: "Executables.", OnPath: true},
			{Path: "Bad_Name", Description: ""},
			{Path: "progs", Description: "Typo twin.", ManifestExclude: true},
			{Path: "a/b/c/d"},
		},
		Scenarios: []Scenario{
			{Item: "x", Home: "dev/work"},
			{Item: "x", Home: "programs/bin"},
			{Item: "y", Home: "nowhere"},
		},
	}
}

func TestRuleTypes(t *testing.T) {
	tests := []struct {
		name   string
		typ    string
		params map[string]any
		want   []string // substrings expected in problems, in order; nil means pass
	}{
		{"max_top_level pass", "max_top_level", map[string]any{"value": int64(10)}, nil},
		{"max_top_level fail", "max_top_level", map[string]any{"value": int64(2)}, []string{"5 top-level folders: [Bad_Name, a, dev, programs, progs]"}},
		{"top_level_name", "top_level_name", map[string]any{"pattern": "[a-z]+"}, []string{"bad top-level name: Bad_Name"}},
		{"name", "name", map[string]any{"pattern": "[a-z]+"}, []string{"bad name: Bad_Name"}},
		{"max_depth", "max_depth", map[string]any{"value": int64(3)}, []string{"too deep: a/b/c/d"}},
		{"has_description", "has_description", nil, []string{"no description: Bad_Name", "no description: a/b/c/d"}},
		{"forbidden_names pass", "forbidden_names", map[string]any{"names": []any{"tmp"}}, nil},
		{"forbidden_names fail", "forbidden_names", map[string]any{"names": []any{"dev"}}, []string{"forbidden top-level name: dev"}},
		{"on_path_exact pass", "on_path_exact", map[string]any{"paths": []any{"programs/bin"}}, nil},
		{"on_path_exact fail", "on_path_exact", map[string]any{"paths": []any{"dev"}}, []string{"on_path is [programs/bin], expected [dev]"}},
		{"manifest_exclude_exact pass", "manifest_exclude_exact", map[string]any{"paths": []any{"progs"}}, nil},
		{"manifest_exclude_exact fail", "manifest_exclude_exact", map[string]any{"paths": []any{}}, []string{"manifest_exclude is [progs], expected []"}},
		{"exact_children pass", "exact_children", map[string]any{"parent": "dev", "children": []any{"work", "personal"}}, nil},
		{"exact_children fail", "exact_children", map[string]any{"parent": "dev", "children": []any{"work"}}, []string{"dev has [personal, work], expected [work]"}},
		{"all_exist", "all_exist", map[string]any{"paths": []any{"dev", "nope"}}, []string{"missing: nope"}},
		{"no_children pass", "no_children", map[string]any{"path": "programs/bin"}, nil},
		{"no_children fail", "no_children", map[string]any{"path": "dev"}, []string{"dev has children [personal, work]"}},
		{"scenarios_resolve", "scenarios_resolve", nil, []string{"two homes for: x", "home does not exist: nowhere (y)"}},
		{"sibling_prefix", "sibling_prefix", map[string]any{"max_shared": int64(3)}, []string{"/: programs and progs share the first 4 characters"}},
		{"no_dead_top_level", "no_dead_top_level", nil, []string{"no scenario uses top-level folder: Bad_Name"}},
	}
	tr := newTree(fixture())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := evalRule(tr, Rule{ID: "X", Type: tt.typ, Text: "t", Params: tt.params})
			if len(tt.want) == 0 {
				if len(got) != 0 {
					t.Fatalf("want pass, got %v", got)
				}
				return
			}
			// Every expected substring must appear, in order.
			i := 0
			for _, g := range got {
				if i < len(tt.want) && strings.Contains(g, tt.want[i]) {
					i++
				}
			}
			if i != len(tt.want) {
				t.Errorf("problems = %q, want containing %q in order", got, tt.want)
			}
		})
	}
}

func TestMalformedRules(t *testing.T) {
	tr := newTree(fixture())
	tests := []struct {
		name string
		rule Rule
		want string
	}{
		{"unknown type", Rule{ID: "X", Type: "bogus", Text: "t"}, "unknown rule type: bogus"},
		{"missing id", Rule{Type: "all_exist", Text: "t"}, "needs id"},
		{"missing param", Rule{ID: "X", Type: "all_exist", Text: "t"}, `"paths"`},
		{"wrong param type", Rule{ID: "X", Type: "max_depth", Text: "t", Params: map[string]any{"value": "3"}}, `"value"`},
		{"bad regexp", Rule{ID: "X", Type: "name", Text: "t", Params: map[string]any{"pattern": "("}}, `"pattern"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := evalRule(tr, tt.rule)
			if len(got) != 1 || !strings.Contains(got[0], tt.want) {
				t.Errorf("problems = %q, want one containing %q", got, tt.want)
			}
		})
	}
}

func TestPatternMatchesWholeString(t *testing.T) {
	tr := newTree(&Layout{Folders: []Folder{{Path: "ab1"}}})
	got := evalRule(tr, Rule{ID: "X", Type: "top_level_name", Text: "t", Params: map[string]any{"pattern": "[a-z]+"}})
	if len(got) != 1 {
		t.Errorf("partial match must not pass, got %v", got)
	}
}

func TestLoadRules(t *testing.T) {
	rules, err := LoadRules(filepath.Join("..", "..", "layouts", "rules.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) == 0 {
		t.Fatal("no rules loaded")
	}
	for _, r := range rules {
		if _, ok := ruleTypes[r.Type]; !ok {
			t.Errorf("rule %s has unimplemented type %q", r.ID, r.Type)
		}
	}
	if _, err := LoadRules(filepath.Join(t.TempDir(), "missing.toml")); err == nil {
		t.Error("want error for missing file")
	}
}

func TestCheckPersonalPasses(t *testing.T) {
	universal, err := LoadRules(filepath.Join("..", "..", "layouts", "rules.toml"))
	if err != nil {
		t.Fatal(err)
	}
	rep := Check(filepath.Join("..", "..", "layouts", "personal"), universal)
	if rep.Failed() != 0 {
		t.Errorf("personal layout fails: %+v", rep)
	}
	if len(rep.Outcomes) < len(universal) {
		t.Errorf("outcomes = %d, want at least %d", len(rep.Outcomes), len(universal))
	}
}

// layoutDir writes a complete layout directory named name.
func layoutDir(t *testing.T, name, folders, scenarios string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), name)
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, dir, "layout.toml", "[layout]\nname = \""+name+"\"\ndescription = \"d\"\n")
	write(t, dir, "folders.toml", folders)
	write(t, dir, "rules.toml", "")
	write(t, dir, "scenarios.toml", scenarios)
	return dir
}

func TestCheckIncomplete(t *testing.T) {
	ok := "[[folder]]\npath = \"a\"\n"
	sc := "[[scenario]]\nitem = \"i\"\nhome = \"a\"\n"

	t.Run("missing file", func(t *testing.T) {
		dir := layoutDir(t, "x", ok, sc)
		_ = os.Remove(filepath.Join(dir, "scenarios.toml"))
		rep := Check(dir, nil)
		if !slices.Contains(rep.Incomplete, "missing required file: scenarios.toml") || rep.Failed() != 1 {
			t.Errorf("report = %+v", rep)
		}
	})
	t.Run("name mismatch", func(t *testing.T) {
		dir := layoutDir(t, "x", ok, sc)
		write(t, dir, "layout.toml", "[layout]\nname = \"other\"\ndescription = \"d\"\n")
		rep := Check(dir, nil)
		if len(rep.Incomplete) != 1 || !strings.Contains(rep.Incomplete[0], "does not match directory") {
			t.Errorf("incomplete = %q", rep.Incomplete)
		}
	})
	t.Run("no scenarios", func(t *testing.T) {
		rep := Check(layoutDir(t, "x", ok, ""), nil)
		if len(rep.Incomplete) != 1 || !strings.Contains(rep.Incomplete[0], "no [[scenario]]") {
			t.Errorf("incomplete = %q", rep.Incomplete)
		}
	})
	t.Run("invalid toml", func(t *testing.T) {
		rep := Check(layoutDir(t, "x", "[[folder", sc), nil)
		if len(rep.Incomplete) != 1 || !strings.Contains(rep.Incomplete[0], "folders.toml") {
			t.Errorf("incomplete = %q", rep.Incomplete)
		}
	})
	t.Run("not a directory", func(t *testing.T) {
		rep := Check(filepath.Join(t.TempDir(), "nope"), nil)
		if len(rep.Incomplete) != 1 {
			t.Errorf("incomplete = %q", rep.Incomplete)
		}
	})
}

func TestCheckReportsStructure(t *testing.T) {
	dir := layoutDir(t, "x",
		"[[folder]]\npath = \"a/b\"\n[[folder]]\npath = \"c\"\n[[folder]]\npath = \"c\"\n",
		"[[scenario]]\nitem = \"i\"\nhome = \"c\"\n")
	rep := Check(dir, nil)
	if len(rep.Incomplete) != 0 || len(rep.Structure) != 2 || rep.Failed() != 1 {
		t.Errorf("report = %+v", rep)
	}
}
