package sync_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envault/internal/storage"
	"github.com/user/envault/internal/sync"
)

const testPassphrase = "test-passphrase-123"
const testNamespace = "myapp"

func newTestVault(t *testing.T) (*storage.Vault, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.enc")
	v, err := storage.New(path, testPassphrase)
	if err != nil {
		t.Fatalf("create vault: %v", err)
	}
	return v, dir
}

func writeEnvFile(t *testing.T, dir, content string) string {
	t.Helper()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write env file: %v", err)
	}
	return path
}

func TestPushToVault(t *testing.T) {
	v, dir := newTestVault(t)
	envPath := writeEnvFile(t, dir, "FOO=bar\nBAZ=qux\n")

	result, err := sync.PushToVault(v, testNamespace, envPath)
	if err != nil {
		t.Fatalf("PushToVault: %v", err)
	}

	if len(result.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(result.Added))
	}

	val, ok := v.GetSecret(testNamespace, "FOO")
	if !ok || val != "bar" {
		t.Errorf("expected FOO=bar, got %q ok=%v", val, ok)
	}
}

func TestPushToVaultUpdatesExistingKeys(t *testing.T) {
	v, dir := newTestVault(t)
	v.SetSecret(testNamespace, "FOO", "old")

	envPath := writeEnvFile(t, dir, "FOO=new\n")
	result, err := sync.PushToVault(v, testNamespace, envPath)
	if err != nil {
		t.Fatalf("PushToVault: %v", err)
	}

	if len(result.Updated) != 1 || result.Updated[0] != "FOO" {
		t.Errorf("expected FOO in updated, got %v", result.Updated)
	}

	val, _ := v.GetSecret(testNamespace, "FOO")
	if val != "new" {
		t.Errorf("expected FOO=new, got %q", val)
	}
}

func TestPushRemovesDeletedKeys(t *testing.T) {
	v, dir := newTestVault(t)
	v.SetSecret(testNamespace, "GONE", "value")

	envPath := writeEnvFile(t, dir, "FOO=bar\n")
	result, err := sync.PushToVault(v, testNamespace, envPath)
	if err != nil {
		t.Fatalf("PushToVault: %v", err)
	}

	if len(result.Removed) != 1 || result.Removed[0] != "GONE" {
		t.Errorf("expected GONE removed, got %v", result.Removed)
	}
}

func TestPullFromVault(t *testing.T) {
	v, dir := newTestVault(t)
	v.SetSecret(testNamespace, "KEY1", "val1")
	v.SetSecret(testNamespace, "KEY2", "val2")

	envPath := filepath.Join(dir, "out.env")
	result, err := sync.PullFromVault(v, testNamespace, envPath)
	if err != nil {
		t.Fatalf("PullFromVault: %v", err)
	}

	if len(result.Added) != 2 {
		t.Errorf("expected 2 entries written, got %d", len(result.Added))
	}

	info, err := os.Stat(envPath)
	if err != nil {
		t.Fatalf("stat output file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected permissions 0600, got %v", info.Mode().Perm())
	}
}
