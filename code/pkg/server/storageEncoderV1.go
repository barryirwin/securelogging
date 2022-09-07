package slogserver

import (
	"crypto/rsa"
	"strconv"

	c "github.com/FrancoLoyola/noroff-fdp/code/pkg/compress"
	slogcrypto "github.com/FrancoLoyola/noroff-fdp/code/pkg/crypto"

	"go.uber.org/zap"
)

// storageEncoderV1 :
//
// First version of the encoder to transform the slogserver.ServerData into a base64 encrypted string.
// NOTE: This string contains a newline in the end
//
// Data formatted as following:
// timestamp,version,compress,AES-256 blob length,checksum (AES-256),AES-256 password wrapped on RSA,AES-256 blob
//
// AES-256 blob is the ciphertext of the base64 ServerData json string representation:
// struct -> json encoded -> base64 string -> AES-256 encrypted -> Compression (Optional)
func storageEncoderV1(sd ServerData, pass string, pubKey rsa.PublicKey, compress bool, logger *zap.Logger) (string, error) {

	// Header (1st part)
	ts := sd.WireData.ReceivedAt.Format("2006-01-02T15:04:05.0000-07")
	ver := "1"
	comp := strconv.FormatBool(compress)
	errStr := ts + ",Encoder v1 error\n"

	// Generate Blob
	blob, checksum, err := sd.ToBase64()
	if err != nil {
		logger.Error("Base64 encoding failed", zap.Error(err))
		return errStr, err
	}

	// Encryption
	// First Pass
	blob, err = slogcrypto.EncryptWithPass(blob, pass, logger)
	if err != nil {
		logger.Error("Password encryption failed", zap.Error(err))
		return errStr, err
	}
	// Then RSA
	blobPass, err := slogcrypto.EncryptWithRSAkey(pass, pubKey, logger)
	if err != nil {
		logger.Error("Key encryption failed", zap.Error(err))
		return errStr, err
	}

	// Compress the blob if enabled
	if compress {
		blob, err = c.Gzip(blob, logger)
		if err != nil {
			logger.Error("Compression failed", zap.Error(err))
			return errStr, err
		}
	}

	// Header (2nd part)
	blobLen := strconv.Itoa(len(blob))
	chk := strconv.FormatUint(uint64(checksum), 10)

	// Build string accoding to the spec
	out := ts + "," + ver + "," + comp + "," + blobLen + "," + chk + "," + blobPass + "," + blob
	return out, nil
}
