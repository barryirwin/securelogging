package slogcrypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"go.uber.org/zap"
)

// ReadRsaPriv :
//
// Reads the key file, and returns the key to be used (copy, not pointer)
//
// https://medium.com/rahasak/golang-rsa-cryptography-1f1897ada311
func ReadRsaPriv(keyPath string, logger *zap.Logger) (rsa.PrivateKey, error) {
	keyData, err := ioutil.ReadFile(keyPath)
	if err != nil {
		logger.Error("Failed to read the private key file", zap.Error(err))
		return rsa.PrivateKey{}, err
	}

	keyBlock, _ := pem.Decode(keyData)
	if keyBlock == nil {
		logger.Error("Failed to decode the key, invalid key")
		return rsa.PrivateKey{}, fmt.Errorf("Failed to decode the key, invalid key content")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		logger.Error("Failed to parse the private key", zap.Error(err))
		return rsa.PrivateKey{}, err
	}

	return *privateKey, nil
}
