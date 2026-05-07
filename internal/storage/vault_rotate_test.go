package storage

import (
	"path/filepath"
	"testing"
)

func TestVaultAll(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.vault")
	v, err := New(path, "passphrase")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	v.Set("FOO", "bar")
	v.Set("BAZ", "qux")

	all := v.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
	if all["FOO"] != "bar" {
		t.Errorf("FOO: want bar, got %q", all["FOO"])
	}
	if all["BAZ"] != "qux" {
		t.Errorf("BAZ: want qux, got %q", all["BAZ"])
	}
}

func TestVaultAllReturnsCopy(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.vault")
	v, _ := New(path, "passphrase")
	v.Set("X", "original")

	all := v.All()
	all["X"] = "mutated"

	// Original vault should be unaffected
	if v.Get("X") != "original" {
		t.Errorf("All() should return a copy, not a reference")
	}
}

func TestVaultGetMissingKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.vault")
	v, _ := New(path, "passphrase")

	if got := v.Get("MISSING"); got != "" {
		t.Errorf("expected empty string for missing key, got %q", got)
	}
}

func TestVaultAllEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.vault")
	v, _ := New(path, "passphrase")

	all := v.All()
	if len(all) != 0 {
		t.Errorf("expected empty map, got %d entries", len(all))
	}
}
