// This file contains all the functions to create the keys, since they are only used here.
//
// The whole process is based on: https://medium.com/rahasak/golang-rsa-cryptography-1f1897ada311
// Adding logger and some parameter configurations and easier names, last one is preference though...

package slogcrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"go.uber.org/zap"
)

// GenerateRSAkeys :
//
// Generates a key pair with the desired names and size
func GenerateRSAkeys(privFileName, pubFileName string, keySize int, logger *zap.Logger) error {
	// Generate private key, pass it to generate the public one, finally store
	privKey, err := generatePrivKey(keySize, logger)
	if err != nil {
		logger.Error("Failed to generate the private key", zap.Error(err))
		return err
	}

	err = savePrivKey(privFileName, privKey, logger)
	if err != nil {
		logger.Error("Failed to store the private key", zap.Error(err))
		return err
	}

	err = savePubKey(pubFileName, privKey, logger)
	if err != nil {
		logger.Error("Failed to store the public key", zap.Error(err))
		return err
	}
	return nil
}

func generatePrivKey(keySize int, logger *zap.Logger) (*rsa.PrivateKey, error) {
	// Generate
	privKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		logger.Error("Failed to generate private key", zap.Error(err))
		return nil, err
	}
	// Validate
	err = privKey.Validate()
	if err != nil {
		logger.Error("Failed to validate private key", zap.Error(err))
		return nil, err
	}

	return privKey, nil
}

func savePrivKey(fileName string, privKey *rsa.PrivateKey, logger *zap.Logger) error {
	// Private key stream
	privateKeyBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	}

	// Create file
	f, err := os.Create(fileName)
	if err != nil {
		logger.Error("Failed to create the file", zap.Error(err))
		return err
	}

	// Store
	err = pem.Encode(f, privateKeyBlock)
	if err != nil {
		logger.Error("Failed to store encoded key to file", zap.Error(err))
		return err
	}

	return nil
}

func savePubKey(fileName string, keyPair *rsa.PrivateKey, logger *zap.Logger) error {
	// Public key stream
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&keyPair.PublicKey)
	if err != nil {
		logger.Error("Failed to generate the public key stream", zap.Error(err))
		return err
	}

	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}

	// Create file
	f, err := os.Create(fileName)
	if err != nil {
		logger.Error("Failed to create the file", zap.Error(err))
		return err
	}

	err = pem.Encode(f, publicKeyBlock)
	if err != nil {
		logger.Error("Failed to store encoded key to file", zap.Error(err))
		return err
	}

	return nil
}
