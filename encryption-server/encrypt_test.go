package main

import (
	"reflect"
	"testing"
)

// Simple length checks
func TestGenerateEncryptionKey(t *testing.T) {
	{
		key := randomBytes(32)
		if len(key) != 32 {
			t.Errorf("Length of key does not match expected")
		}
	}
	{
		key := randomBytes(256)
		if len(key) != 256 {
			t.Errorf("Length of key does not match expected")
		}
	}
}

// Runs some tests which will utilise randomness, trying 100
// encryptions and decryptions and making sure they match up
func TestEncrypt(t *testing.T) {
	for i := 0; i < 100; i++ {
		data := []byte("The quick brown fox jumps over the lazy dog")
		key := randomBytes(64)

		encrypted, err := encrypt(data, key)
		if err != nil {
			t.Fatalf("Did not expect an error, got %s", err)
		}

		decrypted, err := decrypt(encrypted, key)
		if err != nil {
			t.Fatalf("Did not expect an error, got %s", err)
		}

		// The message we started with should match what we decrypt at the end
		if !reflect.DeepEqual(data, decrypted) {
			// Error as a string and bytes
			t.Errorf("Decrypted '%s' does not match expected input message '%s'", string(decrypted), string(data))
			t.Errorf("Decrypted '%x' does not match expected input message '%x'", decrypted, data)
		}
	}
}
