package slogcrypto

import (
	"crypto/sha256"
	"encoding/hex"

	"go.uber.org/zap"
)

// GenerateSHA256Hash :
//
// Generates a 32 bytes SHA256 hash for any given string.
// Returns the byte and string representation of it
func GenerateSHA256Hash(str string, logger *zap.Logger) ([]byte, string, error) {
	h := sha256.New()
	_, err := h.Write([]byte(str))
	if err != nil {
		logger.Error("", zap.Error(err))
		return nil, "", err
	}
	sha256Hash := h.Sum(nil)
	sha256String := hex.EncodeToString(sha256Hash)

	return sha256Hash, sha256String, nil
}
