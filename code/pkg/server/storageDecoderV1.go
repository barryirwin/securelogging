package slogserver

import (
	"crypto/rsa"
	"fmt"
	"strconv"
	"strings"
	"time"

	slogcrypto "github.com/FrancoLoyola/noroff-fdp/code/pkg/crypto"
	"github.com/FrancoLoyola/noroff-fdp/code/pkg/decompress"

	slogstring "github.com/FrancoLoyola/noroff-fdp/code/pkg/string"
	"go.uber.org/zap"
)

// storageDecoderV1 :
//
// First version of the decoder to transform the slogserver.ServerData base64 encrypted string/line to a struct.
//
// Decoding assumes that all lines are good/valid/pristine until some error occurs, then the field SuspiciousRe is set to true to flag the error/corruption, but decoding is still being attempted. If decoding fails and wouldn't make sense to continue, it will return a struct flagged and with the "data" content being: "Decoder v1 - Something is fishy or corrupted..." instead of error.
//
// NOTE: The string contains a newline in the end. This will be trimmed, as the checksum is for the blob without it.
//
// Data formatted as following:
//
// timestamp,version,compress,AES-256 blob length,checksum (AES-256),AES-256 password wrapped on RSA,AES-256 blob
//
// AES-256 blob is the ciphertext of of the base64 ServerData json string representation:
// struct -> json encoded -> base64 string -> AES-256 encrypted -> Compression (Optional)
func storageDecoderV1(line, pass string, privKey rsa.PrivateKey, logger *zap.Logger) ServerData {
	// Default error struct to return if something is really fishy/wrong
	errStruct := NewServerData(NewWireData("Decoder v1 - Something is fishy or corrupted...", "127.0.0.1", "-", false))
	errStruct.SuspiciousRe = true
	errStruct.DecryptedData = "Decoder v1 - Something is fishy or corrupted..."

	// Cleanup and separation
	line = slogstring.TrimWhitespace(line)
	fields := strings.Split(line, ",")

	// Bail out if we have more or less fields than expected, can cause more problems
	if len(fields) != 7 {
		logger.Warn("Unexpected number of fields, want 7, got: " + fmt.Sprint(len(fields)))
		return errStruct
	}

	// Assume that all is ok until an error happens
	suspiciousRe := false

	// Line format
	ts := fields[0]
	encVer := fields[1]
	comp := fields[2]
	lineBlobLen := fields[3]
	lineChecksum := fields[4]
	aesRSApass := fields[5]
	encryptedBlob := fields[6]

	// First verify that the password provided by the user is the same as the stored in the file
	aesRSApass, err := slogcrypto.DecryptWithRSAkey(aesRSApass, privKey, logger)
	if err != nil {
		logger.Warn("Failed to RSA decrypt", zap.Error(err))
		return errStruct
	}
	_, rsaPassHash, err := slogcrypto.GenerateSHA256Hash(aesRSApass, logger)
	if err != nil {
		logger.Warn("Failed to generate hash for the stored password", zap.Error(err))
		return errStruct
	}
	_, passHash, err := slogcrypto.GenerateSHA256Hash(pass, logger)
	if err != nil {
		logger.Warn("Failed to generate hash for the provided password", zap.Error(err))
		return errStruct
	}
	if passHash != rsaPassHash {
		logger.Warn("Hash for the provided password and stored one do not match")
		return errStruct
	}
	// Extract the timestamp of the line to later compare against the struct timestamp
	layout := "2006-01-02T15:04:05.0000-07"
	lineTs, err := time.Parse(layout, ts)
	if err != nil {
		logger.Warn("Failed to parse the timestamp in the start of the line", zap.Error(err))
		suspiciousRe = true
	}

	// Version
	encodingVersion, err := strconv.Atoi(encVer)
	if err != nil {
		logger.Warn("Failed to parse the enconding version number", zap.Error(err))
		suspiciousRe = true
		encodingVersion = 1
	}
	if encodingVersion != 1 {
		logger.Warn("Unexpected encoding version: " + fmt.Sprint(encodingVersion))
		suspiciousRe = true
	}

	// Onwards is not in the same order, as the checksum of the blob is generated before any compression/encryption (see encoderV1)
	// But the length is after compression, hence the need to first check length, uncompress if required, decrypt (RSA then password) and finally check the checksum

	// Length
	blobLen, err := strconv.Atoi(lineBlobLen)
	if err != nil {
		logger.Warn("Encrypted/Compressed blob length couldn't be parsed", zap.Error(err))
		suspiciousRe = true
	}
	if blobLen != len(encryptedBlob) {
		logger.Warn("Length stored in the line: " + lineBlobLen + " is different from the actual blob: " + fmt.Sprint(len(encryptedBlob)))
		suspiciousRe = true
	}

	// Compression
	compressed := false
	compressed, err = strconv.ParseBool(comp)
	if err != nil {
		logger.Warn("Compression fields couldn't be parsed", zap.Error(err))
		suspiciousRe = true
	}
	if compressed {
		encryptedBlob, err = decompress.Gzip(encryptedBlob, logger)
		if err != nil {
			logger.Warn("Failed to decompress the encrypted blob, will try to parse it anyway...", zap.Error(err))
			suspiciousRe = true
		}
	}

	// AES-256 Decryption
	clearBlob, err := slogcrypto.DecryptWithPass(encryptedBlob, pass, logger)
	if err != nil {
		logger.Warn("Failed to Password decrypt", zap.Error(err))
		return errStruct
	}

	// If it made it here, should be ok to attempt to decode the string to a struct and compare that stuff is correct
	out, err := ServerDataFromBase64(clearBlob)
	if err != nil {
		logger.Warn("Failed to parse the base64 string to ServerData", zap.Error(err))
		return errStruct
	}
	// Timestamps
	if lineTs.Format(layout) != out.WireData.ReceivedAt.Format(layout) {
		logger.Warn("Timestamps between the stored data and the struct do not match")
		suspiciousRe = true
	}
	// Checksum
	_, chk, err := out.ToBase64()
	if err != nil {
		logger.Warn("Failed to re-generate the checksum", zap.Error(err))
		suspiciousRe = true
	}
	structChecksum := strconv.FormatUint(uint64(chk), 10)
	if lineChecksum != structChecksum {
		logger.Warn("Checksums from the stored data and the struct do not match")
		suspiciousRe = true
	}

	out.SuspiciousRe = suspiciousRe
	if suspiciousRe {
		out.DecryptedData = errStruct.DecryptedData + out.DecryptedData
	}
	return out
}
