// Package sync provides functionality for syncing .env files with the encrypted vault.
package sync

import (
	"fmt"
	"os"

	"github.com/user/envault/internal/env"
	"github.com/user/envault/internal/storage"
)

// Result holds the outcome of a sync operation.
type Result struct {
	Added   []string
	Updated []string
	Removed []string
}

// PushToVault reads the .env file at envPath and stores all secrets into the vault
// under the given namespace. Returns a Result describing what changed.
func PushToVault(v *storage.Vault, namespace, envPath string) (Result, error) {
	f, err := os.Open(envPath)
	if err != nil {
		return Result{}, fmt.Errorf("open env file: %w", err)
	}
	defer f.Close()

	parsed, err := env.Parse(f)
	if err != nil {
		return Result{}, fmt.Errorf("parse env file: %w", err)
	}

	existing := v.ListSecrets(namespace)
	existingSet := make(map[string]struct{}, len(existing))
	for _, k := range existing {
		existingSet[k] = struct{}{}
	}

	var result Result
	for _, entry := range parsed {
		if _, found := existingSet[entry.Key]; found {
			result.Updated = append(result.Updated, entry.Key)
		} else {
			result.Added = append(result.Added, entry.Key)
		}
		v.SetSecret(namespace, entry.Key, entry.Value)
	}

	// Detect keys removed from the env file.
	incomingSet := make(map[string]struct{}, len(parsed))
	for _, entry := range parsed {
		incomingSet[entry.Key] = struct{}{}
	}
	for _, k := range existing {
		if _, found := incomingSet[k]; !found {
			result.Removed = append(result.Removed, k)
			v.DeleteSecret(namespace, k)
		}
	}

	return result, nil
}

// PullFromVault writes all secrets stored in the vault under namespace to the
// .env file at envPath, creating or overwriting it.
func PullFromVault(v *storage.Vault, namespace, envPath string) (Result, error) {
	keys := v.ListSecrets(namespace)
	entries := make([]env.Entry, 0, len(keys))
	for _, k := range keys {
		val, ok := v.GetSecret(namespace, k)
		if !ok {
			continue
		}
		entries = append(entries, env.Entry{Key: k, Value: val})
	}

	f, err := os.OpenFile(envPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return Result{}, fmt.Errorf("open env file for writing: %w", err)
	}
	defer f.Close()

	if err := env.Write(f, entries); err != nil {
		return Result{}, fmt.Errorf("write env file: %w", err)
	}

	var result Result
	for _, e := range entries {
		result.Added = append(result.Added, e.Key)
	}
	return result, nil
}
