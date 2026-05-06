package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envault/envault/internal/project"
	"github.com/envault/envault/internal/storage"
)

func setupListDir(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "envault-list-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	return dir, func() {
		os.Chdir(origDir)
		os.RemoveAll(dir)
	}
}

// captureStdout redirects os.Stdout to a pipe, calls fn, then restores
// os.Stdout and returns everything written to it as a string.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf strings.Builder
	buf.ReadFrom(r)
	return buf.String()
}

func TestRunListNoSecrets(t *testing.T) {
	dir, cleanup := setupListDir(t)
	defer cleanup()

	writeConfig(t, dir, "myapp", ".env")

	const pass = "listpass"
	_, err := storage.New(filepath.Join(dir, "myapp.vault"), pass)
	if err != nil {
		t.Fatalf("failed to create vault: %v", err)
	}

	// Should not error even when vault is empty.
	if err := RunList(pass, false); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunListShowsKeys(t *testing.T) {
	dir, cleanup := setupListDir(t)
	defer cleanup()

	writeConfig(t, dir, "myapp", ".env")
	writeEnvFile(t, dir, "KEY1=value1\nKEY2=value2\n")

	const pass = "listpass"
	if err := RunPush(pass); err != nil {
		t.Fatalf("push failed: %v", err)
	}

	out := captureStdout(t, func() {
		if err := RunList(pass, false); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "KEY1") || !strings.Contains(out, "KEY2") {
		t.Errorf("expected KEY1 and KEY2 in output, got: %q", out)
	}
	if strings.Contains(out, "value1") {
		t.Errorf("values should not appear when showValues=false, got: %q", out)
	}
}

func TestRunListShowsValues(t *testing.T) {
	dir, cleanup := setupListDir(t)
	defer cleanup()

	writeConfig(t, dir, "myapp", ".env")
	writeEnvFile(t, dir, "SECRET=hunter2\n")

	const pass = "listpass"
	if err := RunPush(pass); err != nil {
		t.Fatalf("push failed: %v", err)
	}

	out := captureStdout(t, func() {
		if err := RunList(pass, true); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "SECRET=hunter2") {
		t.Errorf("expected SECRET=hunter2 in output, got: %q", out)
	}
}

func TestRunListMissingConfig(t *testing.T) {
	_, cleanup := setupListDir(t)
	defer cleanup()
	// No config written — should return an error.
	if err := RunList("pass", false); err == nil {
		t.Error("expected error when config is missing")
	}
}

func init() {
	// Ensure project package is used (avoids import cycle lint warnings).
	_ = project.Load
}
