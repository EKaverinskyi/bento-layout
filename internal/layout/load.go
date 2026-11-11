package layout

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// File names inside a layout directory.
const (
	metaFile      = "layout.toml"
	foldersFile   = "folders.toml"
	rulesFile     = "rules.toml"
	scenariosFile = "scenarios.toml"
)

// Load reads and validates the layout stored in dir. layout.toml and
// folders.toml are required; rules.toml and scenarios.toml are optional.
// Unknown keys in layout.toml, folders.toml and scenarios.toml are errors, so
// typos are not silently ignored.
func Load(dir string) (*Layout, error) {
	var meta struct {
		Layout struct {
			Name        string `toml:"name"`
			Description string `toml:"description"`
		} `toml:"layout"`
	}
	if err := decodeFile(dir, metaFile, true, &meta); err != nil {
		return nil, err
	}
	var folders struct {
		Folder []Folder `toml:"folder"`
	}
	if err := decodeFile(dir, foldersFile, true, &folders); err != nil {
		return nil, err
	}
	var scenarios struct {
		Scenario []Scenario `toml:"scenario"`
	}
	if err := decodeFile(dir, scenariosFile, false, &scenarios); err != nil {
		return nil, err
	}
	var rules struct {
		Rule []map[string]any `toml:"rule"`
	}
	if err := decodeFile(dir, rulesFile, false, &rules); err != nil {
		return nil, err
	}

	l := &Layout{
		Name:        meta.Layout.Name,
		Description: meta.Layout.Description,
		Folders:     folders.Folder,
		Scenarios:   scenarios.Scenario,
	}
	for _, raw := range rules.Rule {
		r := Rule{Params: map[string]any{}}
		for k, v := range raw {
			switch k {
			case "id":
				r.ID, _ = v.(string)
			case "type":
				r.Type, _ = v.(string)
			case "text":
				r.Text, _ = v.(string)
			default:
				r.Params[k] = v
			}
		}
		l.Rules = append(l.Rules, r)
	}
	if err := l.Validate(); err != nil {
		return nil, fmt.Errorf("layout %s: %w", dir, err)
	}
	return l, nil
}

// decodeFile decodes dir/name into v. A missing optional file is not an error.
func decodeFile(dir, name string, required bool, v any) error {
	p := filepath.Join(dir, name)
	data, err := os.ReadFile(p)
	if err != nil {
		if !required && errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read %s: %w", p, err)
	}
	dec := toml.NewDecoder(bytes.NewReader(data))
	// Rules carry free-form parameters, so only the other files are strict.
	if name != rulesFile {
		dec.DisallowUnknownFields()
	}
	if err := dec.Decode(v); err != nil {
		var sme *toml.StrictMissingError
		if errors.As(err, &sme) {
			return fmt.Errorf("parse %s: %s", p, sme.String())
		}
		return fmt.Errorf("parse %s: %w", p, err)
	}
	return nil
}

// Validate checks that the layout can be safely created on disk: it has a
// name, every folder path is a clean relative path without ".." or duplicates,
// every parent folder is listed, and every scenario home exists.
func (l *Layout) Validate() error {
	if l.Name == "" {
		return errors.New("layout name is empty")
	}
	if len(l.Folders) == 0 {
		return errors.New("layout has no folders")
	}
	seen := make(map[string]bool, len(l.Folders))
	for _, f := range l.Folders {
		if err := validatePath(f.Path); err != nil {
			return fmt.Errorf("folder %q: %w", f.Path, err)
		}
		if seen[f.Path] {
			return fmt.Errorf("folder %q is listed twice", f.Path)
		}
		seen[f.Path] = true
	}
	for _, f := range l.Folders {
		if parent := path.Dir(f.Path); parent != "." && !seen[parent] {
			return fmt.Errorf("folder %q: parent %q is not listed", f.Path, parent)
		}
	}
	for _, s := range l.Scenarios {
		if !seen[s.Home] {
			return fmt.Errorf("scenario %q: home %q is not a folder of the layout", s.Item, s.Home)
		}
	}
	return nil
}

func validatePath(p string) error {
	switch {
	case p == "":
		return errors.New("path is empty")
	case strings.Contains(p, `\`):
		return errors.New("path must use / as separator")
	case path.IsAbs(p):
		return errors.New("path must be relative")
	case path.Clean(p) != p:
		return errors.New("path must be clean (no '.', '..', empty or trailing segments)")
	}
	for _, seg := range strings.Split(p, "/") {
		if seg == ".." || seg == "." {
			return errors.New("path must be clean (no '.', '..', empty or trailing segments)")
		}
	}
	return nil
}
