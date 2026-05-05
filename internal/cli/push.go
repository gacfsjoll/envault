package cli

import (
	"fmt"
	"os"

	"github.com/envault/envault/internal/project"
	"github.com/envault/envault/internal/storage"
	"github.com/envault/envault/internal/sync"
)

// RunPush reads the local .env file and pushes its secrets into the encrypted vault.
// passphrase is used to open (or create) the vault file on disk.
func RunPush(passphrase string) error {
	cfg, err := project.Load(".envault.json")
	if err != nil {
		return fmt.Errorf("load project config: %w", err)
	}

	vaultPath := cfg.VaultName + ".vault"
	vault, err := storage.New(vaultPath, passphrase)
	if err != nil {
		return fmt.Errorf("open vault %s: %w", vaultPath, err)
	}

	envFile := cfg.EnvFile
	if envFile == "" {
		envFile = ".env"
	}

	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		return fmt.Errorf("env file %q not found", envFile)
	}

	if err := sync.PushToVault(envFile, vault, passphrase); err != nil {
		return fmt.Errorf("push to vault: %w", err)
	}

	fmt.Printf("✔  Pushed %s → %s\n", envFile, vaultPath)
	return nil
}
