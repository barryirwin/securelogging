package slogclient

import (
	"crypto/rsa"
	"strings"

	slogcrypto "github.com/FrancoLoyola/noroff-fdp/code/pkg/crypto"
	"go.uber.org/zap"
)

// ClientDecoderV1 :
//
// Reads the incoming string, should be formatted as:
//
// "AES-256 blob,RSA encrypted AES-256 pass hash"
//
// Attempts to decode the payload, if ok, returns the cleartext,false,nil.
// If something is fishy, returns the payload as received,true,nil.
// If something fails, returns the payload as received,true,error.
func ClientDecoderV1(payload string, privKey *rsa.PrivateKey, logger *zap.Logger) (string, bool, error) {

	// Data should be only 2 fields separated by a comma
	tmp := strings.Split(payload, ",")
	if len(tmp) != 2 {
		logger.Debug("Data from the wire is not in the expected format, different than 2 fields:", zap.Int("Got", len(tmp)))
		return payload, true, nil
	}

	// Use the private key to get the AES-256 password
	pass, err := slogcrypto.DecryptWithRSAkey(tmp[1], *privKey, logger)
	if err != nil {
		logger.Warn("Failed to retrieve the AES-256 password", zap.Error(err))
		return payload, true, err
	}

	// Hard-coded password for the test script,
	//pass := "56b4fbe06471a2a378bb654e94925c39ed04bf70846b7632441500f8c3c4fa5e"

	// Decrypt the AES-256 part of the payload
	cleartext, err := slogcrypto.DecryptWithPass(tmp[0], pass, logger)
	if err != nil {
		logger.Warn("Error decrypting the AES-256 part", zap.Error(err))
		return payload, true, err
	}

	// Verify that there are 5 fields as the encoder
	elements := strings.Split(cleartext, ",")

	if len(elements) != 5 {
		logger.Debug("Decrypted data succesfully, but the number of fields is not the expected", zap.Int("Got", len(elements)))
		return payload, true, nil
	}

	ts := elements[0]
	host := elements[1]
	file := elements[2]
	checksum := elements[3]
	decodedLog := elements[4]

	if len(ts) == 0 {
		logger.Debug("Decrypted data succesfully, but the timestamp is empty")
		return payload, true, nil
	}
	if len(host) == 0 {
		logger.Debug("Decrypted data succesfully, but the host is empty")
		return payload, true, nil
	}
	if len(file) == 0 {
		logger.Debug("Decrypted data succesfully, but the file name is empty")
		return payload, true, nil
	}
	_, chksum, _ := slogcrypto.GenerateSHA256Hash(decodedLog, logger)
	if checksum != chksum {
		logger.Debug("Decrypted data succesfully, but the checksums do not match")
		return payload, true, nil
	}

	return cleartext, false, nil
}
