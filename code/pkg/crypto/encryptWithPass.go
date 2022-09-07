package slogcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"

	"go.uber.org/zap"
)

// EncryptWithPass :
//
// Encrypts the plaintext with a hash of the password as input for AES. The hashing of the password is done by GenerateSHA256Hash
func EncryptWithPass(plaintext, pass string, logger *zap.Logger) (string, error) {
	// Based on: https://www.melvinvivas.com/how-to-encrypt-and-decrypt-data-using-aes/
	ciphertext := []byte{}
	// Generate a hash for the password, use that as key for the AES key
	hb, _, err := GenerateSHA256Hash(pass, logger)
	if err != nil {
		logger.Error("Failed to generate password hash, returning plaintext as is", zap.Error(err))
		return plaintext, err
	}

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(hb)
	if err != nil {
		logger.Error("Failed to create cipher block, returning plaintext as is", zap.Error(err))
		return plaintext, err
	}

	// Create a new GCM :
	// https://en.wikipedia.org/wiki/Galois/Counter_Mode
	// https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		logger.Error("Failed to create GCM, returning plaintext as is", zap.Error(err))
		return plaintext, err
	}

	// Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		logger.Error("Failed to create nonce, returning plaintext as is", zap.Error(err))
		return plaintext, err
	}

	// Encrypt the data using aesGCM.Seal
	// Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data.
	// The first nonce argument in Seal is the prefix.
	bytetext := []byte(plaintext)
	ciphertext = aesGCM.Seal(nonce, nonce, bytetext, nil)

	out := hex.EncodeToString(ciphertext)
	return out, nil
}
