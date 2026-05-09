package storage

// Set stores a key-value pair in the vault, overwriting any existing value
// for that key. The change is only persisted when Save is called.
func (v *Vault) Set(key, value string) {
	v.data[key] = value
}

// Get retrieves the value associated with key from the vault.
// It returns the value and true if the key exists, or an empty string
// and false if the key is not present.
func (v *Vault) Get(key string) (string, bool) {
	val, ok := v.data[key]
	return val, ok
}

// Delete removes a key from the vault. It is a no-op if the key does not
// exist. The change is only persisted when Save is called.
func (v *Vault) Delete(key string) {
	delete(v.data, key)
}
