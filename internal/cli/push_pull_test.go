package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envault/envault/internal/project"
)

const testPassphrase = "test-secret-passphrase"

func setupPushPullDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	old := chdir(t, dir)
	t.Cleanup(func() { _ = os.Chdir(old) })
	return dir
}

func writeConfig(t *testing.T, vaultName, envFile string) {
	t.Helper()
	cfg := &project.Config{VaultName: vaultName, EnvFile: envFile}
	if err := project.Save(".envault.json", cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}
}

func TestRunPushCreatesVault(t *testing.T) {
	setupPushPullDir(t)
	writeConfig(t, "myproject", ".env")

	if err := os.WriteFile(".env", []byte("KEY=value\nFOO=bar\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := RunPush(testPassphrase); err != nil {
		t.Fatalf("RunPush: %v", err)
	}

	if _, err := os.Stat("myproject.vault"); err != nil {
		t.Fatalf("expected vault file to exist: %v", err)
	}
}

func TestRunPullRestoresEnvFile(t *testing.T) {
	setupPushPullDir(t)
	writeConfig(t, "myproject", ".env")

	original := "KEY=value\nFOO=bar\n"
	if err := os.WriteFile(".env", []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := RunPush(testPassphrase); err != nil {
		t.Fatalf("RunPush: %v", err)
	}

	if err := os.Remove(".env"); err != nil {
		t.Fatal(err)
	}

	if err := RunPull(testPassphrase); err != nil {
		t.Fatalf("RunPull: %v", err)
	}

	data, err := os.ReadFile(".env")
	if err != nil {
		t.Fatalf("read restored .env: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("restored .env file is empty")
	}
}

func TestRunPushMissingEnvFile(t *testing.T) {
	setupPushPullDir(t)
	writeConfig(t, "myproject", ".env")

	if err := RunPush(testPassphrase); err == nil {
		t.Fatal("expected error when .env is missing, got nil")
	}
}

func TestRunPullMissingVault(t *testing.T) {
	setupPushPullDir(t)
	writeConfig(t, "myproject", ".env")

	if err := RunPull(testPassphrase); err == nil {
		t.Fatal("expected error when vault is missing, got nil")
	}
}

func TestRunPullWrongPassphrase(t *testing.T) {
	dir := setupPushPullDir(t)
	_ = dir
	writeConfig(t, "myproject", ".env")

	if err := os.WriteFile(".env", []byte("SECRET=abc\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := RunPush(testPassphrase); err != nil {
		t.Fatalf("RunPush: %v", err)
	}

	if err := RunPull("wrong-passphrase"); err == nil {
		t.Fatal("expected error with wrong passphrase, got nil")
	}
}

// chdir is defined in init_test.go; replicated here for package-level use.
var _ = filepath.Join // ensure filepath imported
