package project_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envault/internal/project"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "envault-project-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestSaveAndLoad(t *testing.T) {
	dir := tempDir(t)
	cfg := &project.Config{
		VaultName:      "my-app",
		DefaultProfile: "development",
		EnvFile:        ".env",
	}
	if err := project.Save(dir, cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := project.Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.VaultName != cfg.VaultName {
		t.Errorf("VaultName: got %q, want %q", loaded.VaultName, cfg.VaultName)
	}
	if loaded.DefaultProfile != cfg.DefaultProfile {
		t.Errorf("DefaultProfile: got %q, want %q", loaded.DefaultProfile, cfg.DefaultProfile)
	}
	if loaded.EnvFile != cfg.EnvFile {
		t.Errorf("EnvFile: got %q, want %q", loaded.EnvFile, cfg.EnvFile)
	}
}

func TestLoadMissingFile(t *testing.T) {
	dir := tempDir(t)
	_, err := project.Load(dir)
	if err != project.ErrNoConfig {
		t.Errorf("expected ErrNoConfig, got %v", err)
	}
}

func TestSaveDefaultEnvFile(t *testing.T) {
	dir := tempDir(t)
	cfg := &project.Config{VaultName: "demo"}
	if err := project.Save(dir, cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := project.Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.EnvFile != ".env" {
		t.Errorf("expected default EnvFile '.env', got %q", loaded.EnvFile)
	}
}

func TestSaveRequiresVaultName(t *testing.T) {
	dir := tempDir(t)
	err := project.Save(dir, &project.Config{})
	if err == nil {
		t.Error("expected error for empty vault_name, got nil")
	}
}

func TestSaveCreatesFile(t *testing.T) {
	dir := tempDir(t)
	cfg := &project.Config{VaultName: "test-vault", EnvFile: ".env.local"}
	if err := project.Save(dir, cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}
	path := filepath.Join(dir, ".envault.json")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("config file not created: %v", err)
	}
	if info.Mode().Perm() != 0644 {
		t.Errorf("expected mode 0644, got %v", info.Mode().Perm())
	}
}
