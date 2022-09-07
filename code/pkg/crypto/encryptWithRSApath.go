package slogcrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"

	"go.uber.org/zap"
)

// EncryptWithRSApath :
//
// Encrypt the cleartext with the given public RSA key path, returns the ciphertext of it.
// If it fails, returns the cleartext as is
//
// https://medium.com/rahasak/golang-rsa-cryptography-1f1897ada311
func EncryptWithRSApath(cleartext, keyPath string, logger *zap.Logger) (string, error) {
	// Get key
	key, err := ReadRsaPub(keyPath, logger)
	if err != nil {
		logger.Error("Failed to read or get the public key", zap.Error(err))
		return cleartext, err
	}

	// Prepare data
	msg := []byte(cleartext)
	rnd := rand.Reader
	hash := sha512.New()

	// Encrypt
	ciperText, err := rsa.EncryptOAEP(hash, rnd, &key, msg, nil)
	if err != nil {
		logger.Error("Failed to encrypt cleartext", zap.Error(err))
		return "", err
	}

	ciphertext := base64.StdEncoding.EncodeToString(ciperText)
	return ciphertext, nil
}
