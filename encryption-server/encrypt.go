// Handles the encryption side of things, most of this has been picked
// out of https://leanpub.com/gocrypto/read - good resource!
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	mathrand "math/rand"
)

const (
	nonceSize = aes.BlockSize
	macSize   = 32 // Output size of HMAC-SHA-256
	cKeySize  = 32 // Cipher key size - AES-256
	mKeySize  = 32 // HMAC key size - HMAC-SHA-256

	keySize = cKeySize + mKeySize
)

var (
	ErrorInvalidKeyLength   = errors.New("Key length is invalid")
	ErrorHMACNotMultiple    = errors.New("HMAC is not multiple of block size")
	ErrorBlockSizeIncorrect = errors.New("block size incorrect")
	ErrorHMACNotEqual       = errors.New("HMAC is not equal")
	ErrorCouldNotUnpad      = errors.New("Failed to unpad")
)

// Generates an encryption key (slice of bytes) at the specified
// length. Should be called ideally with 32 or 64 depending on usage
func randomBytes(length int) []byte {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		// Generate a byte between 0 and z
		bytes[i] = byte('0' + mathrand.Intn('z'-'0'))
	}
	return bytes
}

// Encrypt will take some data and an encryption key, it will return the
// encrypted blocks or an error
func encrypt(message, key []byte) ([]byte, error) {

	if len(key) != keySize {
		return nil, ErrorInvalidKeyLength
	}
	iv := randomBytes(nonceSize)

	// Pad the message to the block size
	padBy := len(message) % aes.BlockSize
	if padBy > 0 {
		padding := bytes.Repeat([]byte{0}, aes.BlockSize-padBy)
		message = append(message, padding...)
	}
	ct := make([]byte, len(message))

	// NewCipher only returns an error with an invalid key size,
	// but the key size was checked at the beginning of the function.
	c, _ := aes.NewCipher(key[:cKeySize])
	ctr := cipher.NewCBCEncrypter(c, iv)
	ctr.CryptBlocks(ct, message)

	h := hmac.New(sha256.New, key[cKeySize:])
	ct = append(iv, ct...)
	h.Write(ct)
	ct = h.Sum(ct)

	return ct, nil
}

// Decrypts some data with a given key and returns the decrypted
// message (or error)
func decrypt(message, key []byte) ([]byte, error) {

	// Various validations on the sizes
	if len(key) != keySize {
		return nil, ErrorInvalidKeyLength
	}

	if (len(message) % aes.BlockSize) != 0 {
		return nil, ErrorHMACNotMultiple
	}

	if len(message) < (4 * aes.BlockSize) {
		return nil, ErrorBlockSizeIncorrect
	}

	macStart := len(message) - macSize
	tag := message[macStart:]
	out := make([]byte, macStart-nonceSize)
	message = message[:macStart]

	h := hmac.New(sha256.New, key[cKeySize:])
	h.Write(message)
	mac := h.Sum(nil)
	if !hmac.Equal(mac, tag) {
		return nil, ErrorHMACNotEqual
	}

	// NewCipher only returns an error with an invalid key size,
	// but the key size was checked at the beginning of the function.
	c, _ := aes.NewCipher(key[:cKeySize])
	ctr := cipher.NewCBCDecrypter(c, message[:nonceSize])
	ctr.CryptBlocks(out, message[nonceSize:])

	// Remove padding
	for i := len(out) - 1; i > 0; i-- {
		if out[i] != byte(0) {
			out = out[:i+1]
			break
		}
	}

	return out, nil
}
