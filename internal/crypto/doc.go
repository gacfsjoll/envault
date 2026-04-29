// Package crypto provides AES-256-GCM encryption and decryption utilities
// used by envault to protect .env file contents at rest.
//
// Key derivation uses SHA-256 to convert a user-supplied passphrase into a
// 32-byte symmetric key suitable for AES-256.
//
// Encrypted blobs are formatted as:
//
//	[ nonce (12 bytes) | ciphertext + auth tag ]
//
// The GCM authentication tag ensures both confidentiality and integrity,
// so any tampering or use of the wrong key will result in a decryption error.
//
// Usage:
//
//	key := crypto.DeriveKey(passphrase)
//	encrypted, err := crypto.Encrypt(key, []byte(envContents))
//	decrypted, err := crypto.Decrypt(key, encrypted)
package crypto
