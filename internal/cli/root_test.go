package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// layoutsDir holds the minimal layout used by all tests.
var layoutsDir = filepath.Join("..", "..", "testdata", "layouts")

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

func TestInitByName(t *testing.T) {
	target := t.TempDir()
	code, out, errOut := run(t, "init", "mini", target, "--layouts-dir", layoutsDir)
	if code != 0 {
		t.Fatalf("code = %d, stderr = %q", code, errOut)
	}
	if !strings.Contains(out, "created 4, 0 already existed") {
		t.Errorf("stdout = %q", out)
	}
	if st, err := os.Stat(filepath.Join(target, "alpha", "inner")); err != nil || !st.IsDir() {
		t.Errorf("alpha/inner missing: %v", err)
	}
}

func TestInitDryRun(t *testing.T) {
	target := filepath.Join(t.TempDir(), "x")
	dir := filepath.Join(layoutsDir, "mini")
	code, out, _ := run(t, "init", dir, target, "-n")
	if code != 0 || !strings.Contains(out, "would create "+filepath.Join(target, "alpha")+"\n") {
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
	dir, err := filepath.Abs(filepath.Join(layoutsDir, "mini"))
	if err != nil {
		t.Fatal(err)
	}
	cwd := t.TempDir()
	t.Chdir(cwd)
	code, out, _ := run(t, "init", dir, ".", "-n")
	if code != 0 || !strings.Contains(out, "would create alpha\n") ||
		!strings.Contains(out, "would create alpha"+string(filepath.Separator)+"inner\n") {
		t.Fatalf("target .: code = %d, stdout = %q", code, out)
	}
	code, out, _ = run(t, "init", dir, "sub/bento", "-n")
	want := "would create " + filepath.Join("sub", "bento", "alpha") + "\n"
	if code != 0 || !strings.Contains(out, want) {
		t.Fatalf("target sub/bento: code = %d, stdout = %q, want %q", code, out, want)
	}
	code, out, _ = run(t, "init", dir, filepath.Join(cwd, "abs"), "-n")
	if code != 0 || !strings.Contains(out, "would create abs"+string(filepath.Separator)+"alpha\n") {
		t.Fatalf("absolute target: code = %d, stdout = %q", code, out)
	}
}

func TestInitRequiresTarget(t *testing.T) {
	code, _, errOut := run(t, "init", "mini", "--layouts-dir", layoutsDir)
	if code != 1 || !strings.Contains(errOut, "accepts 2 arg(s), received 1") {
		t.Errorf("code = %d, stderr = %q", code, errOut)
	}
}

func TestInitEmptyTarget(t *testing.T) {
	dir, err := filepath.Abs(filepath.Join(layoutsDir, "mini"))
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(t.TempDir())
	code, _, errOut := run(t, "init", dir, "")
	if code == 0 || !strings.Contains(errOut, "target is empty") {
		t.Errorf("code = %d, stderr = %q", code, errOut)
	}
	if entries, _ := os.ReadDir("."); len(entries) != 0 {
		t.Errorf("created %d entries despite error", len(entries))
	}
}

func TestInitRepairsDamagedLayout(t *testing.T) {
	target := t.TempDir()
	if code, _, errOut := run(t, "init", "mini", target, "--layouts-dir", layoutsDir); code != 0 {
		t.Fatalf("first init: code = %d, stderr = %q", code, errOut)
	}
	// Damage the layout: remove one folder and put a file in another.
	if err := os.Remove(filepath.Join(target, "gamma")); err != nil {
		t.Fatal(err)
	}
	keep := filepath.Join(target, "alpha", "inner", "keep.txt")
	if err := os.WriteFile(keep, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}

	code, out, errOut := run(t, "init", "mini", target, "--layouts-dir", layoutsDir)
	if code != 0 {
		t.Fatalf("second init: code = %d, stderr = %q", code, errOut)
	}
	if !strings.Contains(out, "created "+filepath.Join(target, "gamma")+"\n") {
		t.Errorf("stdout = %q, want gamma reported as created", out)
	}
	for _, other := range []string{"alpha", "inner", "beta"} {
		if strings.Contains(out, "created "+filepath.Join(target, other)+"\n") {
			t.Errorf("stdout = %q, %s must not be recreated", out, other)
		}
	}
	if !strings.Contains(out, "created 1, 3 already existed") {
		t.Errorf("stdout = %q, want summary \"created 1, 3 already existed\"", out)
	}
	if st, err := os.Stat(filepath.Join(target, "gamma")); err != nil || !st.IsDir() {
		t.Errorf("gamma not restored: %v", err)
	}
	if data, err := os.ReadFile(keep); err != nil || string(data) != "data" {
		t.Errorf("existing file touched: %q, %v", data, err)
	}
}

func TestCheckPasses(t *testing.T) {
	layouts := filepath.Join("..", "..", "layouts")
	drive := t.TempDir()
	if code, _, errOut := run(t, "init", "personal", drive, "--layouts-dir", layouts); code != 0 {
		t.Fatalf("init: %q", errOut)
	}
	code, out, errOut := run(t, "check", "--layout", "personal", "--path", drive, "--layouts-dir", layouts)
	if code != 0 {
		t.Fatalf("code = %d, stdout = %q, stderr = %q", code, out, errOut)
	}
	if !strings.Contains(out, "-- personal: 17 rules, 0 failing") {
		t.Errorf("stdout = %q", out)
	}
}

func TestCheckFailingLayoutExitsNonZero(t *testing.T) {
	// mini has top-level folders no scenario uses, so U08 fails.
	code, out, _ := run(t, "check", "--layout", "mini", "--path", t.TempDir(), "--layouts-dir", layoutsDir, "--rules", filepath.Join("..", "..", "layouts", "rules.toml"))
	if code != 1 {
		t.Fatalf("code = %d, want 1", code)
	}
	if !strings.Contains(out, "FAIL  U08") || !strings.Contains(out, "no scenario uses top-level folder: beta") {
		t.Errorf("stdout = %q", out)
	}
}

func TestCheckMissingRules(t *testing.T) {
	code, _, errOut := run(t, "check", "--layout", "mini", "--path", t.TempDir(), "--layouts-dir", layoutsDir)
	if code != 1 || !strings.Contains(errOut, "rules.toml") {
		t.Errorf("code = %d, stderr = %q", code, errOut)
	}
}

func TestCheckRequiresBothFlags(t *testing.T) {
	layouts := filepath.Join("..", "..", "layouts")
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"no flags", []string{"check", "--layouts-dir", layouts}, "required flag(s)"},
		{"no path", []string{"check", "--layout", "personal", "--layouts-dir", layouts}, `"path"`},
		{"no layout", []string{"check", "--path", t.TempDir(), "--layouts-dir", layouts}, `"layout"`},
		{"positional", []string{"check", "personal", "--layout", "personal", "--path", t.TempDir()}, "unknown command"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, errOut := run(t, tt.args...)
			if code != 1 || !strings.Contains(errOut, tt.want) {
				t.Errorf("code = %d, stderr = %q, want containing %q", code, errOut, tt.want)
			}
		})
	}
}

func TestCheckPathMustExist(t *testing.T) {
	code, _, errOut := run(t, "check", "--layout", "personal", "--path", filepath.Join(t.TempDir(), "nope"),
		"--layouts-dir", filepath.Join("..", "..", "layouts"))
	if code != 1 || errOut == "" {
		t.Errorf("code = %d, stderr = %q", code, errOut)
	}
}

func TestCheckPathAgainstCreatedLayout(t *testing.T) {
	layouts := filepath.Join("..", "..", "layouts")
	drive := t.TempDir()
	if code, _, errOut := run(t, "init", "personal", drive, "--layouts-dir", layouts); code != 0 {
		t.Fatalf("init: %q", errOut)
	}
	code, out, errOut := run(t, "check", "--layout", "personal", "--path", drive, "--layouts-dir", layouts)
	if code != 0 || !strings.Contains(out, "matches the layout") {
		t.Fatalf("code = %d, stdout = %q, stderr = %q", code, out, errOut)
	}

	if err := os.Remove(filepath.Join(drive, "config", "secrets")); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(drive, "stray"), 0o755); err != nil {
		t.Fatal(err)
	}
	code, out, _ = run(t, "check", "--layout", "personal", "--path", drive, "--layouts-dir", layouts)
	if code != 1 || !strings.Contains(out, "MISSING   config/secrets") || !strings.Contains(out, "extra     stray") {
		t.Errorf("code = %d, stdout = %q", code, out)
	}
}
