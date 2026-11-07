package cli

import (
	"bytes"
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
