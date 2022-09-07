package slogcrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"

	"go.uber.org/zap"
)

// DecryptWithRSAkey :
//
// Decrypt the ciphertext with the given private RSA key, returns the cleartext.
//
// If it fails returns the ciphertext as is.
//
// https://medium.com/rahasak/golang-rsa-cryptography-1f1897ada311
func DecryptWithRSAkey(ciphertext string, key rsa.PrivateKey, logger *zap.Logger) (string, error) {

	// Prepare data
	msg, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		logger.Error("Failed to base64 decode the ciphertext", zap.Error(err))
		return ciphertext, err
	}
	rnd := rand.Reader
	hash := sha512.New()

	// Decrypt with OAEP
	plainText, err := rsa.DecryptOAEP(hash, rnd, &key, msg, nil)
	if err != nil {
		logger.Error("Failed to decrypt the ciphertext", zap.Error(err))
		return ciphertext, err
	}

	cleartext := string(plainText)
	return cleartext, nil
}
