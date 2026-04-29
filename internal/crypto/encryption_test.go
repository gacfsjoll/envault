package crypto

import (
	"bytes"
	"testing"
)

func TestDeriveKey(t *testing.T) {
	key := DeriveKey("my-secret-passphrase")
	if len(key) != 32 {
		t.Errorf("expected key length 32, got %d", len(key))
	}

	// Same passphrase should produce same key
	key2 := DeriveKey("my-secret-passphrase")
	if !bytes.Equal(key, key2) {
		t.Error("expected deterministic key derivation")
	}

	// Different passphrase should produce different key
	key3 := DeriveKey("different-passphrase")
	if bytes.Equal(key, key3) {
		t.Error("expected different keys for different passphrases")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := DeriveKey("test-passphrase")
	plaintext := []byte("DB_PASSWORD=supersecret\nAPI_KEY=abc123")

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Error("ciphertext should not equal plaintext")
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestEncryptProducesUniqueOutput(t *testing.T) {
	key := DeriveKey("test-passphrase")
	plaintext := []byte("SECRET=value")

	c1, _ := Encrypt(key, plaintext)
	c2, _ := Encrypt(key, plaintext)

	if bytes.Equal(c1, c2) {
		t.Error("expected unique ciphertexts due to random nonce")
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	key := DeriveKey("correct-passphrase")
	wrongKey := DeriveKey("wrong-passphrase")
	plaintext := []byte("SECRET=value")

	ciphertext, _ := Encrypt(key, plaintext)

	_, err := Decrypt(wrongKey, ciphertext)
	if err == nil {
		t.Error("expected error when decrypting with wrong key")
	}
}

func TestDecryptShortCiphertext(t *testing.T) {
	key := DeriveKey("test-passphrase")
	_, err := Decrypt(key, []byte("short"))
	if err == nil {
		t.Error("expected error for short ciphertext")
	}
}
