package cli

import (
	"fmt"
	"strings"

	"github.com/user/envault/internal/project"
	"github.com/user/envault/internal/storage"
)

// RunSet sets one or more key=value pairs directly in the vault without
// requiring a local .env file. This is useful for adding or updating
// individual secrets from the command line.
func RunSet(passphrase string, pairs []string) error {
	if len(pairs) == 0 {
		return fmt.Errorf("no key=value pairs provided")
	}

	cfg, err := project.Load(".envault.json")
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	vault, err := storage.New(cfg.VaultPath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to open vault: %w", err)
	}

	for _, pair := range pairs {
		key, value, err := parsePair(pair)
		if err != nil {
			return err
		}
		vault.Set(key, value)
	}

	if err := vault.Save(passphrase); err != nil {
		return fmt.Errorf("failed to save vault: %w", err)
	}

	if len(pairs) == 1 {
		key, _, _ := parsePair(pairs[0])
		fmt.Printf("Set %s in vault %q\n", key, cfg.VaultName)
	} else {
		fmt.Printf("Set %d keys in vault %q\n", len(pairs), cfg.VaultName)
	}

	return nil
}

// parsePair splits a "KEY=VALUE" string into its key and value components.
// The key must be non-empty. The value may be empty.
func parsePair(pair string) (string, string, error) {
	idx := strings.IndexByte(pair, '=')
	if idx < 0 {
		return "", "", fmt.Errorf("invalid format %q: expected KEY=VALUE", pair)
	}
	key := strings.TrimSpace(pair[:idx])
	if key == "" {
		return "", "", fmt.Errorf("invalid format %q: key must not be empty", pair)
	}
	value := pair[idx+1:]
	return key, value, nil
}
