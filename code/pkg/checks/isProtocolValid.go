package checks

import (
	"go.uber.org/zap"
)

// IsProtocolValid :
//
// Check that the protocol is valid or not
func IsProtocolValid(protocol string, logger *zap.Logger) bool {
	validProtocols := []string{"udp", "UDP", "tcp", "TCP"}
	for _, item := range validProtocols {
		if protocol == item {
			return true
		}
	}
	logger.Error("Protocol is not valid", zap.String("Got", protocol))
	return false
}
