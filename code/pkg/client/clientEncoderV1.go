package slogclient

import (
	"crypto/rsa"
	"os"
	"time"

	slogcrypto "github.com/FrancoLoyola/noroff-fdp/code/pkg/crypto"

	"go.uber.org/zap"
)

// ClientEncoderV1 :
//
// Generates the blob and RSA encrypted password, in this output.
//
// "AES-256 blob,RSA encrypted AES-256 pass"
//
// The password is randomly generated and stored as the second field after the comma.
// The blob is constructed as follows:
// timestamp,hostname,logfile,checksum(of the logline),logline
//
// If it fails to generate the string, return "logline,logfile", so the server at least can get the line, but flag it
func ClientEncoderV1(logline, logfile string, pubKey *rsa.PublicKey, logger *zap.Logger) (string, error) {
	errStr := logfile + "," + logline

	// Get timestamp
	ts := time.Now()
	timestamp := ts.Format("2006-01-02T15:04:05.000-07")
	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		logger.Error("Failed to get the hostname", zap.Error(err))
		return errStr, err
	}
	// Generate checksum of the logline
	_, checksum, err := slogcrypto.GenerateSHA256Hash(logline, logger)
	if err != nil {
		logger.Error("Failed to generate the SHA256 checksum for the log file", zap.Error(err))
		return errStr, err
	}
	// Generate random 32 byte pass for the AES-256 and encrypt
	aesPass := slogcrypto.GenerateRandomPass(32)
	plaintext := timestamp + "," + hostname + "," + logfile + "," + checksum + "," + logline
	cipherAES, err := slogcrypto.EncryptWithPass(plaintext, string(aesPass), logger)
	if err != nil {
		logger.Error("Failed to AES-256 encrypt", zap.Error(err))
		return errStr, err
	}

	// Encrypt the AES-256 pass with RSA public key
	cipherRSA, err := slogcrypto.EncryptWithRSAkey(aesPass, *pubKey, logger)
	if err != nil {
		logger.Error("Failed to RSA encrypt the AES password", zap.Error(err))
		return errStr, err
	}

	out := cipherAES + "," + cipherRSA
	return out, nil
}
