package cli

import (
	"fmt"
	"os"

	"github.com/user/envault/internal/project"
	"github.com/user/envault/internal/storage"
)

// RunDelete removes one or more keys from the vault.
// Usage: envault delete <key> [key2 ...] --passphrase <pass>
func RunDelete(dir string, keys []string, passphrase string) error {
	if len(keys) == 0 {
		return fmt.Errorf("at least one key must be specified")
	}

	cfg, err := project.Load(dir)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	vaultPath := cfg.VaultPath(dir)
	vault, err := storage.New(vaultPath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to open vault: %w", err)
	}

	deleted := 0
	for _, key := range keys {
		if !vault.Has(key) {
			fmt.Fprintf(os.Stderr, "warning: key %q not found in vault, skipping\n", key)
			continue
		}
		vault.Delete(key)
		deleted++
	}

	if deleted == 0 {
		return fmt.Errorf("no keys were deleted")
	}

	if err := vault.Save(); err != nil {
		return fmt.Errorf("failed to save vault: %w", err)
	}

	if deleted == 1 {
		fmt.Printf("Deleted 1 key from vault.\n")
	} else {
		fmt.Printf("Deleted %d keys from vault.\n", deleted)
	}
	return nil
}
