package cli

import (
	"fmt"
	"os"

	"github.com/envault/envault/internal/project"
	"github.com/envault/envault/internal/storage"
)

// RunRotate re-encrypts the vault with a new passphrase.
// It reads the vault with the old passphrase, then writes it back
// with the new passphrase, preserving all stored secrets.
func RunRotate(dir, oldPassphrase, newPassphrase string) error {
	cfg, err := project.Load(dir)
	if err != nil {
		return fmt.Errorf("load project config: %w", err)
	}

	vaultPath := cfg.VaultFile
	if vaultPath == "" {
		vaultPath = cfg.VaultName + ".vault"
	}

	// Open vault with old passphrase
	oldVault, err := storage.New(vaultPath, oldPassphrase)
	if err != nil {
		return fmt.Errorf("open vault with old passphrase: %w", err)
	}

	secrets := oldVault.All()

	// Remove existing vault file so we can create it fresh
	if err := os.Remove(vaultPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove old vault: %w", err)
	}

	// Create new vault with new passphrase
	newVault, err := storage.New(vaultPath, newPassphrase)
	if err != nil {
		return fmt.Errorf("create vault with new passphrase: %w", err)
	}

	for k, v := range secrets {
		newVault.Set(k, v)
	}

	if err := newVault.Save(); err != nil {
		return fmt.Errorf("save rotated vault: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Passphrase rotated successfully. %d secret(s) re-encrypted.\n", len(secrets))
	return nil
}
