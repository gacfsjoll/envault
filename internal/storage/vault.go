// Package storage handles reading and writing encrypted vault files.
package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/envault/envault/internal/crypto"
)

const defaultVaultFile = ".envault"

// Vault represents an encrypted store of key-value secret pairs.
type Vault struct {
	path string
	pass string
}

// VaultData is the serialized structure stored on disk.
type VaultData struct {
	Secrets map[string]string `json:"secrets"`
}

// New returns a Vault bound to the given file path and passphrase.
func New(path, passphrase string) *Vault {
	if path == "" {
		path = defaultVaultFile
	}
	return &Vault{path: path, pass: passphrase}
}

// Load decrypts and deserializes the vault from disk.
// Returns an empty VaultData if the file does not exist yet.
func (v *Vault) Load() (*VaultData, error) {
	ciphertext, err := os.ReadFile(v.path)
	if errors.Is(err, os.ErrNotExist) {
		return &VaultData{Secrets: make(map[string]string)}, nil
	}
	if err != nil {
		return nil, err
	}

	key, err := crypto.DeriveKey(v.pass, nil)
	if err != nil {
		return nil, err
	}

	plaintext, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		return nil, err
	}

	var data VaultData
	if err := json.Unmarshal(plaintext, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// Save serializes and encrypts the vault data to disk.
func (v *Vault) Save(data *VaultData) error {
	plaintext, err := json.Marshal(data)
	if err != nil {
		return err
	}

	key, err := crypto.DeriveKey(v.pass, nil)
	if err != nil {
		return err
	}

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(v.path), 0o700); err != nil {
		return err
	}
	return os.WriteFile(v.path, ciphertext, 0o600)
}
