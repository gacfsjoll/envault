package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envault/internal/storage"
)

func setupDeleteDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	writeConfig(t, dir, "delete-test", ".env")
	return dir
}

func populateVaultForDelete(t *testing.T, dir string, passphrase string, secrets map[string]string) {
	t.Helper()
	cfgPath := filepath.Join(dir, ".envault.json")
	_ = cfgPath
	vaultPath := filepath.Join(dir, "delete-test.vault")
	v, err := storage.New(vaultPath, passphrase)
	if err != nil {
		t.Fatalf("failed to create vault: %v", err)
	}
	for k, val := range secrets {
		v.Set(k, val)
	}
	if err := v.Save(); err != nil {
		t.Fatalf("failed to save vault: %v", err)
	}
}

func TestRunDeleteRemovesKey(t *testing.T) {
	dir := setupDeleteDir(t)
	populateVaultForDelete(t, dir, "pass", map[string]string{"FOO": "bar", "BAZ": "qux"})

	if err := RunDelete(dir, []string{"FOO"}, "pass"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vaultPath := filepath.Join(dir, "delete-test.vault")
	v, err := storage.New(vaultPath, "pass")
	if err != nil {
		t.Fatalf("failed to reload vault: %v", err)
	}
	if v.Has("FOO") {
		t.Error("expected FOO to be deleted")
	}
	if !v.Has("BAZ") {
		t.Error("expected BAZ to still exist")
	}
}

func TestRunDeleteMultipleKeys(t *testing.T) {
	dir := setupDeleteDir(t)
	populateVaultForDelete(t, dir, "pass", map[string]string{"A": "1", "B": "2", "C": "3"})

	if err := RunDelete(dir, []string{"A", "B"}, "pass"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vaultPath := filepath.Join(dir, "delete-test.vault")
	v, _ := storage.New(vaultPath, "pass")
	if v.Has("A") || v.Has("B") {
		t.Error("expected A and B to be deleted")
	}
	if !v.Has("C") {
		t.Error("expected C to remain")
	}
}

func TestRunDeleteNoKeys(t *testing.T) {
	dir := setupDeleteDir(t)
	err := RunDelete(dir, []string{}, "pass")
	if err == nil {
		t.Fatal("expected error when no keys provided")
	}
}

func TestRunDeleteMissingConfig(t *testing.T) {
	dir := t.TempDir()
	err := RunDelete(dir, []string{"FOO"}, "pass")
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestRunDeleteWrongPassphrase(t *testing.T) {
	dir := setupDeleteDir(t)
	populateVaultForDelete(t, dir, "correct", map[string]string{"X": "y"})

	err := RunDelete(dir, []string{"X"}, "wrong")
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

func TestRunDeleteNonExistentKeyWarns(t *testing.T) {
	dir := setupDeleteDir(t)
	populateVaultForDelete(t, dir, "pass", map[string]string{"REAL": "value"})

	// Deleting a missing key alongside a real one should succeed
	if err := RunDelete(dir, []string{"MISSING", "REAL"}, "pass"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vaultPath := filepath.Join(dir, "delete-test.vault")
	v, _ := storage.New(vaultPath, "pass")
	if v.Has("REAL") {
		t.Error("expected REAL to be deleted")
	}
}

func TestRunDeleteAllMissingKeys(t *testing.T) {
	dir := setupDeleteDir(t)
	populateVaultForDelete(t, dir, "pass", map[string]string{"FOO": "bar"})

	err := RunDelete(dir, []string{"NONEXISTENT"}, "pass")
	if err == nil {
		t.Fatal("expected error when all keys are missing")
	}
	_ = os.Stderr
}
