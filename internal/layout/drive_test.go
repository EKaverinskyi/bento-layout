package layout

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestCheckDrive(t *testing.T) {
	l := &Layout{Folders: []Folder{{Path: "a"}, {Path: "a/b"}, {Path: "c"}, {Path: "d"}, {Path: "d/e"}, {Path: "f"}}}
	root := t.TempDir()
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	must(os.MkdirAll(filepath.Join(root, "a", "b"), 0o755))
	must(os.WriteFile(filepath.Join(root, "c"), nil, 0o644))    // file where a folder belongs
	must(os.WriteFile(filepath.Join(root, "d"), nil, 0o644))    // so d/e is blocked, not missing
	must(os.MkdirAll(filepath.Join(root, "a", "extra"), 0o755)) // nested extra
	must(os.MkdirAll(filepath.Join(root, "zzz"), 0o755))        // top-level extra
	must(os.MkdirAll(filepath.Join(root, ".git"), 0o755))       // hidden, ignored

	rep, err := CheckDrive(l, root)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Present != 2 || !slices.Equal(rep.Missing, []string{"f"}) ||
		!slices.Equal(rep.Conflicts, []string{"c", "d"}) ||
		!slices.Equal(rep.Extra, []string{"a/extra", "zzz"}) || !rep.Failed() {
		t.Errorf("report = %+v", rep)
	}

	for _, bad := range []string{filepath.Join(root, "nope"), filepath.Join(root, "c")} {
		if _, err := CheckDrive(l, bad); err == nil {
			t.Errorf("want error for %s", bad)
		}
	}
}

func TestCheckDriveExtraDoesNotFail(t *testing.T) {
	l := &Layout{Folders: []Folder{{Path: "a"}}}
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "a"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "x"), 0o755); err != nil {
		t.Fatal(err)
	}
	rep, err := CheckDrive(l, root)
	if err != nil || rep.Failed() || len(rep.Extra) != 1 {
		t.Errorf("%v, %+v", err, rep)
	}
}
