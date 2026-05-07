package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envault/envault/internal/project"
	"github.com/envault/envault/internal/storage"
)

func setupExecDir(t *testing.T) (string, func()) {
	t.Helper()
	dir := t.TempDir()
	old := chdir(t, dir)
	return dir, func() { os.Chdir(old) }
}

func TestRunExecNoCommand(t *testing.T) {
	_, cleanup := setupExecDir(t)
	defer cleanup()

	err := RunExec("passphrase", []string{})
	if err == nil {
		t.Fatal("expected error for empty args")
	}
}

func TestRunExecMissingConfig(t *testing.T) {
	_, cleanup := setupExecDir(t)
	defer cleanup()

	err := RunExec("passphrase", []string{"env"})
	if err == nil {
		t.Fatal("expected error when config is missing")
	}
}

func TestRunExecInjectsSecrets(t *testing.T) {
	dir, cleanup := setupExecDir(t)
	defer cleanup()

	cfg := &project.Config{VaultName: filepath.Join(dir, "test.vault"), EnvFile: ".env"}
	if err := project.Save(".", cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}

	vault, err := storage.New(cfg.VaultName, "secret")
	if err != nil {
		t.Fatalf("new vault: %v", err)
	}
	vault.Set("EXEC_TEST_VAR", "hello_from_vault")
	if err := vault.Save(); err != nil {
		t.Fatalf("save vault: %v", err)
	}

	// Run a command that will succeed only if the env var is present.
	err = RunExec("secret", []string{"sh", "-c", "test \"$EXEC_TEST_VAR\" = \"hello_from_vault\""})
	if err != nil {
		t.Errorf("expected command to succeed with injected secret, got: %v", err)
	}
}

func TestRunExecWrongPassphrase(t *testing.T) {
	dir, cleanup := setupExecDir(t)
	defer cleanup()

	cfg := &project.Config{VaultName: filepath.Join(dir, "test.vault"), EnvFile: ".env"}
	if err := project.Save(".", cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}

	vault, err := storage.New(cfg.VaultName, "correct")
	if err != nil {
		t.Fatalf("new vault: %v", err)
	}
	if err := vault.Save(); err != nil {
		t.Fatalf("save vault: %v", err)
	}

	err = RunExec("wrong", []string{"echo", "hi"})
	if err == nil {
		t.Fatal("expected error with wrong passphrase")
	}
}
