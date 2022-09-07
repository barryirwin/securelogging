package slogserver

import (
	"crypto/rsa"

	slogclient "github.com/FrancoLoyola/noroff-fdp/code/pkg/client"
	slogconfig "github.com/FrancoLoyola/noroff-fdp/code/pkg/config"
	"go.uber.org/zap"
)

// decryptWorker :
//
// Reads the WireData, decrypts the password used by the client using the key, then the data with that password
// and sends the result to the ServerData channel
//
// The data arriving on the wire should be the AES-256 blob,AES pass RSA encrypted
func decryptWorker(conf slogconfig.ServerConfig, commsPrivKey rsa.PrivateKey, d WireData, chanServerData chan ServerData, logger *zap.Logger) {

	out := NewServerData(d)
	// No need to process anything if the data is from the syslog ports
	if out.WireData.SyslogData {
		out.DecryptedData = out.WireData.Data
		chanServerData <- out
		return
	}

	cleartext, fishy, err := slogclient.ClientDecoderV1(out.WireData.Data, &commsPrivKey, logger)
	// Always attach the returned string to the struct. TO-DO: check if this is required or not
	out.DecryptedData = cleartext

	// Was all ok?
	if err != nil {
		logger.Warn("Decoding of the data from the wire failed", zap.Error(err))
		out.SuspiciousIn = true
		chanServerData <- out
		return
	}
	if fishy {
		logger.Debug("Data from the wire is fishy", zap.String("Source", out.WireData.SrcIP), zap.String("Data", out.WireData.Data))
		out.SuspiciousIn = true
		chanServerData <- out
		return
	}

	chanServerData <- out
}
