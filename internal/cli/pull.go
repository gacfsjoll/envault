package cli

import (
	"fmt"
	"os"

	"github.com/envault/envault/internal/project"
	"github.com/envault/envault/internal/storage"
	"github.com/envault/envault/internal/sync"
)

// RunPull decrypts secrets from the vault and writes them to the local .env file.
// If the env file already exists it is overwritten.
// passphrase must match the one used when the vault was created.
func RunPull(passphrase string) error {
	cfg, err := project.Load(".envault.json")
	if err != nil {
		return fmt.Errorf("load project config: %w", err)
	}

	vaultPath := cfg.VaultName + ".vault"
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		return fmt.Errorf("vault file %q not found — run push first", vaultPath)
	}

	vault, err := storage.New(vaultPath, passphrase)
	if err != nil {
		return fmt.Errorf("open vault %s: %w", vaultPath, err)
	}

	envFile := cfg.EnvFile
	if envFile == "" {
		envFile = ".env"
	}

	if err := sync.PullFromVault(vault, passphrase, envFile); err != nil {
		return fmt.Errorf("pull from vault: %w", err)
	}

	fmt.Printf("✔  Pulled %s → %s\n", vaultPath, envFile)
	return nil
}
