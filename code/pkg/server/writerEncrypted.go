package slogserver

import (
	"os"
	"strings"

	slogcrypto "github.com/FrancoLoyola/noroff-fdp/code/pkg/crypto"
	slogfile "github.com/FrancoLoyola/noroff-fdp/code/pkg/file"

	slogconfig "github.com/FrancoLoyola/noroff-fdp/code/pkg/config"
	"go.uber.org/zap"
)

// writerEncrypted :
//
// Writes the ServerData structs into a file in the encrypted format. Relies on the encoders.
// Reads the last line of the file and calculates the hash for it, then uses that one as a reference
// for the first log. Then each new generated line is hashed for the next struct, thus preserving the chain
func writerEncrypted(file string, data []ServerData, version int, conf slogconfig.ServerConfig, logger *zap.Logger) error {

	// Get last hash of the file before anything, the hashes are from the line WITHOUT the "\n"
	lastLine := slogfile.GetLastLine(file, logger)
	lastLine = strings.TrimSpace(lastLine)
	_, prevHash, err := slogcrypto.GenerateSHA256Hash(lastLine, logger)
	if err != nil {
		logger.Error("Failed to generate hash for the last line", zap.Error(err))
	}

	// Start
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		logger.Error("Failed to open the file: ", zap.Error(err))
		return err
	}

	pass := conf.Password
	compress := conf.CompressedBlob
	keyPath := conf.StoragePubKey
	pubKey, err := slogcrypto.ReadRsaPub(keyPath, logger)
	if err != nil {
		logger.Error("Failed to read public key", zap.Error(err))
		return err
	}

	// Store data based on version releases
	for _, item := range data {
		str := "Invalid encoder version, not writing anything :)"
		// Attach reference to the previous log
		item.PreviousHash = prevHash
		switch version {
		case 1:
			str, err = storageEncoderV1(item, pass, pubKey, compress, logger)
			if err != nil {
				logger.Error("Failed to encode struct", zap.Error(err))
				continue
			}
		}
		f.WriteString(str + "\n")
		// Hash for the previous line WITHOUT "\n", see replayProcessor
		_, prevHash, _ = slogcrypto.GenerateSHA256Hash(str, logger)
	}
	logger.Info("Encrypted data written", zap.Int("messages", len(data)))
	return nil
}
