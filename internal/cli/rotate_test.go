package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envault/envault/internal/project"
	"github.com/envault/envault/internal/storage"
)

func setupRotateDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	cfg := &project.Config{
		VaultName: "test",
		VaultFile: filepath.Join(dir, "test.vault"),
		EnvFile:   ".env",
	}
	if err := project.Save(dir, cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}
	return dir
}

func TestRunRotateChangesPassphrase(t *testing.T) {
	dir := setupRotateDir(t)
	cfg, _ := project.Load(dir)

	vault, err := storage.New(cfg.VaultFile, "old-pass")
	if err != nil {
		t.Fatalf("create vault: %v", err)
	}
	vault.Set("KEY", "value")
	if err := vault.Save(); err != nil {
		t.Fatalf("save vault: %v", err)
	}

	if err := RunRotate(dir, "old-pass", "new-pass"); err != nil {
		t.Fatalf("RunRotate: %v", err)
	}

	// Old passphrase should no longer work
	_, err = storage.New(cfg.VaultFile, "old-pass")
	if err == nil {
		t.Error("expected error opening vault with old passphrase, got nil")
	}

	// New passphrase should work and preserve secrets
	newVault, err := storage.New(cfg.VaultFile, "new-pass")
	if err != nil {
		t.Fatalf("open with new passphrase: %v", err)
	}
	if got := newVault.Get("KEY"); got != "value" {
		t.Errorf("expected KEY=value, got %q", got)
	}
}

func TestRunRotateWrongOldPassphrase(t *testing.T) {
	dir := setupRotateDir(t)
	cfg, _ := project.Load(dir)

	vault, _ := storage.New(cfg.VaultFile, "correct-pass")
	vault.Set("A", "1")
	vault.Save()

	if err := RunRotate(dir, "wrong-pass", "new-pass"); err == nil {
		t.Error("expected error with wrong old passphrase")
	}
}

func TestRunRotateMissingConfig(t *testing.T) {
	dir := t.TempDir()
	if err := RunRotate(dir, "old", "new"); err == nil {
		t.Error("expected error with missing config")
	}
}

func TestRunRotatePreservesAllSecrets(t *testing.T) {
	dir := setupRotateDir(t)
	cfg, _ := project.Load(dir)

	vault, _ := storage.New(cfg.VaultFile, "pass")
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	for k, v := range secrets {
		vault.Set(k, v)
	}
	vault.Save()

	if err := RunRotate(dir, "pass", "newpass"); err != nil {
		t.Fatalf("RunRotate: %v", err)
	}

	newVault, _ := storage.New(cfg.VaultFile, "newpass")
	for k, want := range secrets {
		if got := newVault.Get(k); got != want {
			t.Errorf("key %s: want %q, got %q", k, want, got)
		}
	}
	_ = os.Getenv("_") // suppress unused import
}
