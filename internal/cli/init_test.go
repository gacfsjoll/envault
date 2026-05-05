package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envault/envault/internal/cli"
	"github.com/envault/envault/internal/project"
)

func chdir(t *testing.T, dir string) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
}

func TestRunInitCreatesConfig(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	err := cli.RunInit(cli.InitOptions{VaultName: "myproject", EnvFile: ".env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg, err := project.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if cfg.VaultName != "myproject" {
		t.Errorf("vault name = %q, want %q", cfg.VaultName, "myproject")
	}
	if cfg.EnvFile != ".env" {
		t.Errorf("env file = %q, want %q", cfg.EnvFile, ".env")
	}
}

func TestRunInitDefaultsToDirectoryName(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	err := cli.RunInit(cli.InitOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg, err := project.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if cfg.VaultName != filepath.Base(tmp) {
		t.Errorf("vault name = %q, want %q", cfg.VaultName, filepath.Base(tmp))
	}
}

func TestRunInitRefusesOverwrite(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	if err := cli.RunInit(cli.InitOptions{VaultName: "first"}); err != nil {
		t.Fatalf("first init: %v", err)
	}

	err := cli.RunInit(cli.InitOptions{VaultName: "second"})
	if err == nil {
		t.Fatal("expected error on second init without --force")
	}
}

func TestRunInitForceOverwrites(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	if err := cli.RunInit(cli.InitOptions{VaultName: "first"}); err != nil {
		t.Fatalf("first init: %v", err)
	}

	if err := cli.RunInit(cli.InitOptions{VaultName: "second", Force: true}); err != nil {
		t.Fatalf("force init: %v", err)
	}

	cfg, err := project.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.VaultName != "second" {
		t.Errorf("vault name = %q, want %q", cfg.VaultName, "second")
	}
}
