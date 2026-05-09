package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envault/internal/project"
	"github.com/user/envault/internal/storage"
)

func setupSetDir(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "envault-set-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	old := chdir(t, dir)
	return dir, func() {
		os.Chdir(old)
		os.RemoveAll(dir)
	}
}

func TestRunSetSingleKey(t *testing.T) {
	dir, cleanup := setupSetDir(t)
	defer cleanup()

	cfg := &project.Config{VaultName: "myapp", VaultPath: filepath.Join(dir, "myapp.vault"), EnvFile: ".env"}
	if err := cfg.Save(".envault.json"); err != nil {
		t.Fatalf("save config: %v", err)
	}

	if err := RunSet("secret", []string{"API_KEY=abc123"}); err != nil {
		t.Fatalf("RunSet: %v", err)
	}

	vault, err := storage.New(cfg.VaultPath, "secret")
	if err != nil {
		t.Fatalf("open vault: %v", err)
	}
	if got, ok := vault.Get("API_KEY"); !ok || got != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q (ok=%v)", got, ok)
	}
}

func TestRunSetMultipleKeys(t *testing.T) {
	dir, cleanup := setupSetDir(t)
	defer cleanup()

	cfg := &project.Config{VaultName: "myapp", VaultPath: filepath.Join(dir, "myapp.vault"), EnvFile: ".env"}
	if err := cfg.Save(".envault.json"); err != nil {
		t.Fatalf("save config: %v", err)
	}

	pairs := []string{"FOO=bar", "BAZ=qux", "EMPTY="}
	if err := RunSet("secret", pairs); err != nil {
		t.Fatalf("RunSet: %v", err)
	}

	vault, err := storage.New(cfg.VaultPath, "secret")
	if err != nil {
		t.Fatalf("open vault: %v", err)
	}
	for _, pair := range pairs {
		key, val, _ := parsePair(pair)
		if got, ok := vault.Get(key); !ok || got != val {
			t.Errorf("expected %s=%q, got %q (ok=%v)", key, val, got, ok)
		}
	}
}

func TestRunSetNoPairs(t *testing.T) {
	_, cleanup := setupSetDir(t)
	defer cleanup()

	if err := RunSet("secret", []string{}); err == nil {
		t.Error("expected error for empty pairs, got nil")
	}
}

func TestRunSetInvalidFormat(t *testing.T) {
	dir, cleanup := setupSetDir(t)
	defer cleanup()

	cfg := &project.Config{VaultName: "myapp", VaultPath: filepath.Join(dir, "myapp.vault"), EnvFile: ".env"}
	if err := cfg.Save(".envault.json"); err != nil {
		t.Fatalf("save config: %v", err)
	}

	if err := RunSet("secret", []string{"NOKEYVALUE"}); err == nil {
		t.Error("expected error for missing '=', got nil")
	}
}

func TestRunSetMissingConfig(t *testing.T) {
	_, cleanup := setupSetDir(t)
	defer cleanup()

	if err := RunSet("secret", []string{"KEY=val"}); err == nil {
		t.Error("expected error for missing config, got nil")
	}
}

func TestParsePair(t *testing.T) {
	cases := []struct {
		input   string
		wantKey string
		wantVal string
		wantErr bool
	}{
		{"KEY=VALUE", "KEY", "VALUE", false},
		{"KEY=", "KEY", "", false},
		{"KEY=foo=bar", "KEY", "foo=bar", false},
		{"NOEQUALS", "", "", true},
		{"=NOKEY", "", "", true},
	}
	for _, c := range cases {
		key, val, err := parsePair(c.input)
		if c.wantErr {
			if err == nil {
				t.Errorf("parsePair(%q): expected error", c.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("parsePair(%q): unexpected error: %v", c.input, err)
			continue
		}
		if key != c.wantKey || val != c.wantVal {
			t.Errorf("parsePair(%q): got (%q, %q), want (%q, %q)", c.input, key, val, c.wantKey, c.wantVal)
		}
	}
}
