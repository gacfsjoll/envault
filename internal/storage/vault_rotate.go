package storage

// All returns a copy of all key-value pairs stored in the vault.
// This is used during passphrase rotation to export secrets before
// re-encrypting them with a new passphrase.
func (v *Vault) All() map[string]string {
	copy := make(map[string]string, len(v.data))
	for k, val := range v.data {
		copy[k] = val
	}
	return copy
}

// Get returns the value for the given key, or an empty string if not found.
func (v *Vault) Get(key string) string {
	return v.data[key]
}
