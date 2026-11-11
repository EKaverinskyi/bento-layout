// Package layout loads layout definitions (a directory of TOML files) and
// creates their folder structure on disk.
package layout

// Layout is a parsed layout definition.
type Layout struct {
	Name        string
	Description string
	Folders     []Folder
	Rules       []Rule
	Scenarios   []Scenario
}

// Folder is one folder of the layout, addressed by a slash-separated path
// relative to the layout root.
type Folder struct {
	Path            string `toml:"path"`
	Description     string `toml:"description"`
	OnPath          bool   `toml:"on_path"`
	ManifestExclude bool   `toml:"manifest_exclude"`
}

// Rule is a layout-specific rule. Its parameters depend on Type and are kept
// raw; checking rules is not done at import time.
type Rule struct {
	ID     string         `toml:"id"`
	Type   string         `toml:"type"`
	Text   string         `toml:"text"`
	Params map[string]any `toml:"-"`
}

// Scenario is an item together with the one folder it belongs in.
type Scenario struct {
	Item string `toml:"item"`
	Home string `toml:"home"`
}
