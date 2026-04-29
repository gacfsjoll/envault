// Package storage provides encrypted persistence for envault secrets.
//
// A Vault is bound to a file path and a passphrase. Secrets are stored as a
// JSON map that is encrypted with AES-GCM before being written to disk, using
// the crypto package for all cryptographic operations.
//
// Basic usage:
//
//	v := storage.New(".envault", "my-passphrase")
//
//	// Load existing secrets (returns empty vault if file is missing)
//	data, err := v.Load()
//
//	// Mutate secrets
//	data.Secrets["API_KEY"] = "abc123"
//
//	// Persist back to disk
//	err = v.Save(data)
//
// The vault file is written with 0600 permissions to limit access to the
// owning user.
package storage
