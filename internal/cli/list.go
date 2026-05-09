package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/envault/envault/internal/project"
	"github.com/envault/envault/internal/storage"
)

// RunList prints all secret keys stored in the vault for the current project.
// If showValues is true, decrypted values are also printed.
// Secrets are printed in alphabetical order by key.
func RunList(passphrase string, showValues bool) error {
	cfg, err := project.Load(".envault.json")
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	vault, err := storage.New(cfg.VaultName, passphrase)
	if err != nil {
		return fmt.Errorf("failed to open vault: %w", err)
	}

	secrets := vault.List()
	if len(secrets) == 0 {
		fmt.Fprintln(os.Stdout, "No secrets stored in vault.")
		return nil
	}

	keys := sortedKeys(secrets)

	for _, k := range keys {
		if showValues {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, secrets[k])
		} else {
			fmt.Fprintf(os.Stdout, "%s\n", k)
		}
	}
	return nil
}

// sortedKeys returns the keys of the given map sorted alphabetically.
func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
