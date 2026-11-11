package layout

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func write(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mini(t *testing.T, folders string) string {
	t.Helper()
	dir := t.TempDir()
	write(t, dir, "layout.toml", "[layout]\nname = \"mini\"\n")
	write(t, dir, "folders.toml", folders)
	return dir
}

func TestLoadRepoPersonal(t *testing.T) {
	l, err := Load(filepath.Join("..", "..", "layouts", "personal"))
	if err != nil {
		t.Fatal(err)
	}
	if l.Name != "personal" || len(l.Folders) != 35 || len(l.Rules) == 0 || len(l.Scenarios) == 0 {
		t.Errorf("unexpected layout: %q, %d folders, %d rules, %d scenarios",
			l.Name, len(l.Folders), len(l.Rules), len(l.Scenarios))
	}
	var onPath int
	for _, f := range l.Folders {
		if f.OnPath {
			onPath++
		}
	}
	if onPath != 2 {
		t.Errorf("on_path folders = %d, want 2", onPath)
	}
}

func TestLoadErrors(t *testing.T) {
	tests := []struct {
		name, folders, want string
	}{
		{"missing parent", "[[folder]]\npath = \"a/b\"\n", "parent"},
		{"duplicate", "[[folder]]\npath = \"a\"\n[[folder]]\npath = \"a\"\n", "twice"},
		{"traversal", "[[folder]]\npath = \"../x\"\n", "clean"},
		{"absolute", "[[folder]]\npath = \"/x\"\n", "relative"},
		{"unknown key", "[[folder]]\npath = \"a\"\ndesc = \"typo\"\n", "desc"},
		{"empty", "", "no folders"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Load(mini(t, tt.folders))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("err = %v, want containing %q", err, tt.want)
			}
		})
	}
}

func TestLoadMissingDir(t *testing.T) {
	if _, err := Load(filepath.Join(t.TempDir(), "nope")); err == nil {
		t.Fatal("want error")
	}
}

func TestCreate(t *testing.T) {
	l, err := Load(mini(t, "[[folder]]\npath = \"a\"\n[[folder]]\npath = \"a/b\"\n[[folder]]\npath = \"c\"\n"))
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(t.TempDir(), "bento")

	res, err := Create(l, root, true)
	if err != nil || len(res.Created) != 3 {
		t.Fatalf("dry run: %v, %+v", err, res)
	}
	if _, err := os.Stat(root); err == nil {
		t.Fatal("dry run created root")
	}

	if _, err := Create(l, root, false); err != nil {
		t.Fatal(err)
	}
	if st, err := os.Stat(filepath.Join(root, "a", "b")); err != nil || !st.IsDir() {
		t.Fatalf("a/b not created: %v", err)
	}
	res, err = Create(l, root, false)
	if err != nil || len(res.Created) != 0 || len(res.Existed) != 3 {
		t.Fatalf("second run: %v, %+v", err, res)
	}
}

func TestCreateConflictCreatesNothing(t *testing.T) {
	l, err := Load(mini(t, "[[folder]]\npath = \"a\"\n[[folder]]\npath = \"b\"\n"))
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	write(t, root, "b", "i am a file")
	if _, err := Create(l, root, false); err == nil {
		t.Fatal("want error")
	}
	if _, err := os.Stat(filepath.Join(root, "a")); err == nil {
		t.Error("a was created despite conflict")
	}
}
