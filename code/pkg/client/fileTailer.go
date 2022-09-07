package slogclient

import (
	slogconfig "github.com/FrancoLoyola/noroff-fdp/code/pkg/config"
	slogcrypto "github.com/FrancoLoyola/noroff-fdp/code/pkg/crypto"
	slogfile "github.com/FrancoLoyola/noroff-fdp/code/pkg/file"

	"go.uber.org/zap"
)

// FileTailer :
//
// Tails the file, encodes the data and sends it to the slog server
func FileTailer(filePath string, conf slogconfig.ClientConfig, logger *zap.Logger) {
	// General preparations
	pubKey, err := slogcrypto.ReadRsaPub(conf.CommsPubKey, logger)
	if err != nil {
		logger.Error("Error reading the public key", zap.Error(err))
		return
	}
	serverAddr := conf.ServerIP + ":" + conf.Port

	// Read file
	chanTail := make(chan string, conf.BufferLines)
	logger.Debug("Starting the tailing of file " + filePath + " using protocol " + conf.Protocol)
	go slogfile.TailFile(filePath, chanTail, logger)

	for {
		line := <-chanTail
		// Encode lines
		ciphertext, err := ClientEncoderV1(line, filePath, &pubKey, logger)
		if err != nil {
			logger.Warn("Failed to encode line", zap.Error(err))
		}
		// Send data
		if conf.Protocol == "TCP" {
			err = TCPWriter(ciphertext, serverAddr, logger)
			if err != nil {
				logger.Error("Failed to send data to the server via TCP", zap.Error(err))
			}
		}
		if conf.Protocol == "UDP" {
			err = UDPWriter(ciphertext, serverAddr, logger)
			if err != nil {
				logger.Error("Failed to send data to the server via UDP", zap.Error(err))
			}
		}
	}
}
