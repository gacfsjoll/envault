package storage_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envault/envault/internal/storage"
)

func tempVaultPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, ".envault")
}

func TestLoadEmptyVault(t *testing.T) {
	v := storage.New(tempVaultPath(t), "passphrase")
	data, err := v.Load()
	if err != nil {
		t.Fatalf("expected no error loading missing vault, got: %v", err)
	}
	if len(data.Secrets) != 0 {
		t.Errorf("expected empty secrets, got %d entries", len(data.Secrets))
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := tempVaultPath(t)
	v := storage.New(path, "correct-passphrase")

	original := &storage.VaultData{
		Secrets: map[string]string{
			"DB_HOST": "localhost",
			"API_KEY": "supersecret",
		},
	}

	if err := v.Save(original); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := v.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	for k, want := range original.Secrets {
		got, ok := loaded.Secrets[k]
		if !ok {
			t.Errorf("missing key %q after load", k)
			continue
		}
		if got != want {
			t.Errorf("key %q: got %q, want %q", k, got, want)
		}
	}
}

func TestLoadWithWrongPassphrase(t *testing.T) {
	path := tempVaultPath(t)

	v := storage.New(path, "correct")
	if err := v.Save(&storage.VaultData{Secrets: map[string]string{"X": "y"}}); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	wrong := storage.New(path, "wrong")
	_, err := wrong.Load()
	if err == nil {
		t.Error("expected error when loading with wrong passphrase, got nil")
	}
}

func TestSaveCreatesFileWithRestrictedPermissions(t *testing.T) {
	path := tempVaultPath(t)
	v := storage.New(path, "passphrase")

	if err := v.Save(&storage.VaultData{Secrets: map[string]string{}}); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("expected file permissions 0600, got %o", perm)
	}
}
