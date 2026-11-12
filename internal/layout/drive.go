package layout

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DriveReport compares a directory tree on disk with a layout.
type DriveReport struct {
	Root      string
	Present   int      // layout folders that exist as directories
	Missing   []string // layout folders not on disk
	Conflicts []string // layout paths occupied by something that is not a directory
	Extra     []string // directories on disk that the layout does not list (informational)
}

// Failed reports whether the drive differs from the layout in a way that
// counts as a failure: missing folders or conflicts. Extra folders do not.
func (r *DriveReport) Failed() bool {
	return len(r.Missing) > 0 || len(r.Conflicts) > 0
}

// CheckDrive compares the directory root with l. It only reads; nothing on
// disk is changed. Hidden entries (starting with ".") are ignored.
func CheckDrive(l *Layout, root string) (*DriveReport, error) {
	st, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !st.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", root)
	}
	rep := &DriveReport{Root: root}
	known := make(map[string]bool, len(l.Folders))
	for _, f := range l.Folders {
		known[f.Path] = true
	}
	blocked := func(p string) bool {
		for _, c := range rep.Conflicts {
			if strings.HasPrefix(p, c+"/") {
				return true
			}
		}
		return false
	}
	for _, f := range l.Folders {
		fi, err := os.Lstat(filepath.Join(root, filepath.FromSlash(f.Path)))
		switch {
		case blocked(f.Path):
			// Reported through the conflicting ancestor.
		case err == nil && fi.IsDir():
			rep.Present++
		case err == nil:
			rep.Conflicts = append(rep.Conflicts, f.Path)
		case os.IsNotExist(err):
			rep.Missing = append(rep.Missing, f.Path)
		default:
			return nil, err
		}
	}
	// Directories under the root and under every present layout folder that
	// the layout does not list.
	parents := []string{""}
	for _, f := range l.Folders {
		if fi, err := os.Lstat(filepath.Join(root, filepath.FromSlash(f.Path))); err == nil && fi.IsDir() {
			parents = append(parents, f.Path)
		}
	}
	for _, parent := range parents {
		entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(parent)))
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
				continue
			}
			p := e.Name()
			if parent != "" {
				p = parent + "/" + p
			}
			if !known[p] {
				rep.Extra = append(rep.Extra, p)
			}
		}
	}
	sort.Strings(rep.Extra)
	return rep, nil
}
