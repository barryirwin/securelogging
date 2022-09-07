package slogcrypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"go.uber.org/zap"
)

// ReadRsaPub :
//
// Reads the key file, and returns the key to be used (copy, not pointer)
//
// https://medium.com/rahasak/golang-rsa-cryptography-1f1897ada311
func ReadRsaPub(keyPath string, logger *zap.Logger) (rsa.PublicKey, error) {
	keyData, err := ioutil.ReadFile(keyPath)
	if err != nil {
		logger.Error("Failed to read the public key file", zap.Error(err))
		return rsa.PublicKey{}, err
	}

	keyBlock, _ := pem.Decode(keyData)
	if keyBlock == nil {
		logger.Error("Failed to decode the key, invalid key")
		return rsa.PublicKey{}, fmt.Errorf("Failed to decode the key, invalid key content")
	}

	publicKey, err := x509.ParsePKIXPublicKey(keyBlock.Bytes)
	if err != nil {
		logger.Error("Failed to parse the public key file", zap.Error(err))
		return rsa.PublicKey{}, err
	}
	switch publicKey := publicKey.(type) {
	case *rsa.PublicKey:
		return *publicKey, nil
	default:
		return rsa.PublicKey{}, fmt.Errorf("Key is not public?")
	}
}
