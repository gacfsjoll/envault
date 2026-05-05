package storage

// List returns a copy of all key-value pairs currently stored in the vault.
// The returned map is safe to modify without affecting the vault's internal state.
func (v *Vault) List() map[string]string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	copy := make(map[string]string, len(v.data.Secrets))
	for k, val := range v.data.Secrets {
		copy[k] = val
	}
	return copy
}
