package checks

import (
	"strconv"

	"go.uber.org/zap"
)

// IsPortValid :
//
// Check that the port is a number, within range
// and warns if it is a low number port
func IsPortValid(port string, logger *zap.Logger) bool {
	i, err := strconv.Atoi(port)
	if err != nil {
		logger.Error("Port is not a number", zap.Error(err))
		return false
	}
	if i < 1 || i > 65535 {
		logger.Error("Port is not valid", zap.Error(err))
		return false
	}
	if i < 1024 {
		logger.Warn("Port is under 1023")
	}
	return true
}
