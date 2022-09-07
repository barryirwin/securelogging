package slogcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"

	"go.uber.org/zap"
)

// DecryptWithPass :
//
// Decrypts the ciphertext with a hash of the password as input for AES. The hashing of the password is done by GenerateSHA256Hash
func DecryptWithPass(ciphertext, pass string, logger *zap.Logger) (string, error) {
	// Based on: https://www.melvinvivas.com/how-to-encrypt-and-decrypt-data-using-aes/
	plaintext := []byte{}
	// Generate a hash for the password, use that as key for the AES key
	hb, _, err := GenerateSHA256Hash(pass, logger)
	if err != nil {
		logger.Error("Failed to generate password hash, returning ciphertext as is", zap.Error(err))
		return ciphertext, err
	}
	enc, err := hex.DecodeString(ciphertext)
	if err != nil {
		logger.Error("Failed to decode ciphertext to hex, returning ciphertext as is", zap.Error(err))
		return ciphertext, err
	}

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(hb)
	if err != nil {
		logger.Error("Failed to create cipher block, returning ciphertext as is", zap.Error(err))
		return ciphertext, err
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		logger.Error("Failed to create GCM, returning ciphertext as is", zap.Error(err))
		return ciphertext, err
	}

	// Get the nonce size
	nonceSize := aesGCM.NonceSize()

	// Extract the nonce from the encrypted data
	nonce, ctext := enc[:nonceSize], enc[nonceSize:]

	// Decrypt the data
	plaintext, err = aesGCM.Open(nil, nonce, ctext, nil)
	if err != nil {
		logger.Error("Failed to decrypt, returning ciphertext as is", zap.Error(err))
		return ciphertext, err
	}

	out := string(plaintext)
	return out, nil
}
