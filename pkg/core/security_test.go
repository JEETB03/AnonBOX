package core

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	originalText := []byte("Hello, AnonBOX!")

	encrypted, err := Encrypt(originalText, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if !bytes.Equal(originalText, decrypted) {
		t.Errorf("Decrypted text does not match original. Got %s, want %s", decrypted, originalText)
	}
}

func TestDecryptInvalidKey(t *testing.T) {
	key1, _ := GenerateKey()
	key2, _ := GenerateKey()

	originalText := []byte("Secret Data")
	encrypted, _ := Encrypt(originalText, key1)

	_, err := Decrypt(encrypted, key2)
	if err == nil {
		t.Error("Decryption should fail with wrong key, but it succeeded")
	}
}
