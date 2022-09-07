package slogserver

import (
	"bufio"
	"os"

	slogconfig "github.com/FrancoLoyola/noroff-fdp/code/pkg/config"
	slogcrypto "github.com/FrancoLoyola/noroff-fdp/code/pkg/crypto"
	"go.uber.org/zap"
)

// ReplayProcessor :
//
// Reads a .slog file and re-process it
func ReplayProcessor(fileIn, fileOut string, conf slogconfig.ServerConfig, logger *zap.Logger) error {
	// Get key
	privKey, err := slogcrypto.ReadRsaPriv(conf.StoragePrivKey, logger)
	if err != nil {
		logger.Error("Failed to read the private key, can't replay anything", zap.Error(err))
		return err
	}

	// From the generateSHA256Hash_test.go
	fileStartHash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	// Placeholders
	prevHash := ""
	nSuspicious := 0  // Replay suspicious
	inSuspicious := 0 // Incomming suspicious
	nLines := 0
	batchSize := 1000
	data := []ServerData{}

	// start/end timestamps?

	// Iterate through file
	file, err := os.Open(fileIn)
	if err != nil {
		logger.Error("Failed to open the file: "+fileIn, zap.Error(err))
		return err
	}
	defer file.Close()

	logger.Info("Bear in mind that if the segment read is was not the start of the file, the order of the first log in the file can't be verified. So will always be flagged as suspicious.")
	logger.Info("To solve this, just start slightly before the segment of interest and disregard the first log failure. Or replay the whole original file.")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		nLines++

		currServerData := storageDecoderV1(line, conf.Password, privKey, logger)
		if currServerData.PreviousHash == fileStartHash {
			logger.Info("Looks like this is the start of the file")
			prevHash = fileStartHash
		}

		if currServerData.PreviousHash != prevHash {
			logger.Warn("Missmatch when comparing this line previous hash: ", zap.Int("Line Number:", nLines), zap.String("Struct Hash", currServerData.PreviousHash), zap.String("Calculated hash", prevHash))
			currServerData.SuspiciousRe = true
			currServerData.DecryptedData = "Replay Processor - " + currServerData.DecryptedData
		}

		if currServerData.SuspiciousRe {
			nSuspicious++
		}
		if currServerData.SuspiciousIn {
			inSuspicious++
		}
		data = append(data, currServerData)

		// Time to write and reset
		if nLines%batchSize == 0 {
			err = writerTxt(fileOut, data, logger)
			if err != nil {
				logger.Error("Failed to write the replay file, stopping the replay", zap.Error(err))
				return err
			}
			err = writerInflux(data, conf, logger)
			if err != nil {
				logger.Error("Failed to write replay file to Influx", zap.Error(err))
				return err
			}
			data = []ServerData{}
		}

		// The current line is the verification hash for the next iteration, WITHOUT "\n", see the writerEncrypted
		_, prevHash, _ = slogcrypto.GenerateSHA256Hash(line, logger)
	}

	// Write once again whatever was left from the batch
	err = writerTxt(fileOut, data, logger)
	if err != nil {
		logger.Error("Failed to write the replay file, stopping the replay", zap.Error(err))
		return err
	}
	err = writerInflux(data, conf, logger)
	if err != nil {
		logger.Error("Failed to write replay file to Influx", zap.Error(err))
		return err
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Error while reading the file", zap.Error(err))
		return err
	}

	logger.Info("Replay Summary:", zap.Int("Total lines processed", nLines), zap.Int("Replay Suspicious Lines", nSuspicious), zap.Int("Suspicious packets (When data arrived)", inSuspicious))

	return nil
}
