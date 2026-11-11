package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func run(t *testing.T, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	var out, errb bytes.Buffer
	code = Execute(args, &out, &errb)
	return code, out.String(), errb.String()
}

func TestVersionCommand(t *testing.T) {
	code, out, _ := run(t, "version")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.HasPrefix(out, "bento-layout ") {
		t.Errorf("stdout = %q, want prefix %q", out, "bento-layout ")
	}
}

func TestVersionFlag(t *testing.T) {
	code, out, _ := run(t, "--version")
	if code != 0 || !strings.HasPrefix(out, "bento-layout ") {
		t.Errorf("code = %d, stdout = %q", code, out)
	}
}

func TestUnknownCommand(t *testing.T) {
	code, _, errOut := run(t, "nope")
	if code == 0 {
		t.Fatal("exit code = 0, want non-zero")
	}
	if !strings.Contains(errOut, "unknown command") {
		t.Errorf("stderr = %q, want unknown command error", errOut)
	}
}

func TestInitPersonal(t *testing.T) {
	target := t.TempDir()
	layouts := filepath.Join("..", "..", "layouts")
	code, out, errOut := run(t, "init", "personal", target, "--layouts-dir", layouts)
	if code != 0 {
		t.Fatalf("code = %d, stderr = %q", code, errOut)
	}
	if !strings.Contains(out, "35 created") {
		t.Errorf("stdout = %q", out)
	}
	if st, err := os.Stat(filepath.Join(target, "media", "pictures", "memes")); err != nil || !st.IsDir() {
		t.Errorf("media/pictures/memes missing: %v", err)
	}
}

func TestInitDryRun(t *testing.T) {
	target := filepath.Join(t.TempDir(), "x")
	dir := filepath.Join("..", "..", "layouts", "personal")
	code, out, _ := run(t, "init", dir, target, "-n")
	if code != 0 || !strings.Contains(out, "would create "+filepath.Join(target, "dev")+"\n") {
		t.Fatalf("code = %d, stdout = %q", code, out)
	}
	if _, err := os.Stat(target); err == nil {
		t.Error("dry run created target")
	}
}

func TestInitUnknownLayout(t *testing.T) {
	code, _, errOut := run(t, "init", "nope", t.TempDir(), "--layouts-dir", t.TempDir())
	if code == 0 || !strings.Contains(errOut, "not found") {
		t.Errorf("code = %d, stderr = %q", code, errOut)
	}
}

func TestInitDryRunPathsRelativeToCwd(t *testing.T) {
	dir, err := filepath.Abs(filepath.Join("..", "..", "layouts", "personal"))
	if err != nil {
		t.Fatal(err)
	}
	cwd := t.TempDir()
	t.Chdir(cwd)
	code, out, _ := run(t, "init", dir, "-n")
	if code != 0 || !strings.Contains(out, "would create dev\n") ||
		!strings.Contains(out, "would create dev"+string(filepath.Separator)+"work\n") {
		t.Fatalf("target .: code = %d, stdout = %q", code, out)
	}
	code, out, _ = run(t, "init", dir, "sub/bento", "-n")
	want := "would create " + filepath.Join("sub", "bento", "dev") + "\n"
	if code != 0 || !strings.Contains(out, want) {
		t.Fatalf("target sub/bento: code = %d, stdout = %q, want %q", code, out, want)
	}
	code, out, _ = run(t, "init", dir, filepath.Join(cwd, "abs"), "-n")
	if code != 0 || !strings.Contains(out, "would create abs"+string(filepath.Separator)+"dev\n") {
		t.Fatalf("absolute target: code = %d, stdout = %q", code, out)
	}
}
