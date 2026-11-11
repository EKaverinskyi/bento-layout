package layout

import (
	"fmt"
	"os"
	"path/filepath"
)

// Result reports what Create did, or would do in a dry run.
type Result struct {
	Created []string // folders created (or to be created)
	Existed []string // folders that were already there
}

// Create makes every folder of l below root, creating root itself if needed.
// It never deletes, renames or changes anything that exists. It fails before
// creating anything if a layout path is occupied by a non-directory. With
// dryRun nothing is written.
func Create(l *Layout, root string, dryRun bool) (*Result, error) {
	if st, err := os.Stat(root); err == nil && !st.IsDir() {
		return nil, fmt.Errorf("%s exists and is not a directory", root)
	}
	res := &Result{}
	var todo []string
	for _, f := range l.Folders {
		p := filepath.Join(root, filepath.FromSlash(f.Path))
		st, err := os.Stat(p)
		switch {
		case err == nil && st.IsDir():
			res.Existed = append(res.Existed, f.Path)
		case err == nil:
			return nil, fmt.Errorf("%s exists and is not a directory", p)
		case os.IsNotExist(err):
			res.Created = append(res.Created, f.Path)
			todo = append(todo, p)
		default:
			return nil, err
		}
	}
	if dryRun {
		return res, nil
	}
	for _, p := range todo {
		if err := os.MkdirAll(p, 0o755); err != nil {
			return res, err
		}
	}
	return res, nil
}
