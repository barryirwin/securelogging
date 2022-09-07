package slogserver

import (
	slogconfig "github.com/FrancoLoyola/noroff-fdp/code/pkg/config"
	slogcrypto "github.com/FrancoLoyola/noroff-fdp/code/pkg/crypto"
	"go.uber.org/zap"
)

// WireDataProcessor :
//
// Reads encrypted data coming on the wire and coordinates the decryption.
// Loads the comms private key to be able to decrypt data arriving
//
// Creates a go routine per received packet to decrypt all of them in parallell.
// ServerData channel is used by the go routines created here, only passed over.
func WireDataProcessor(conf slogconfig.ServerConfig, chanWireData chan WireData, chanServerData chan ServerData, logger *zap.Logger) {
	n := 1
	// Load the key
	privKey, err := slogcrypto.ReadRsaPriv(conf.CommsPrivKey, logger)
	if err != nil {
		logger.Fatal("Failed to load comms private key, can't parse any arriving data", zap.Error(err))
		return
	}
	for {
		select {
		case d := <-chanWireData:
			n++
			if n%1000 == 0 {
				logger.Debug("Got 1000 more packets", zap.Int("total", n))
			}
			go decryptWorker(conf, privKey, d, chanServerData, logger)
		}
	}
}
