package slogconfig

import (
	"strconv"

	"go.uber.org/zap"
)

// isBuffSizeOk :
//
// Check the sanity of the buffer provided in the config:
// Non-zero, non-negative, too small, too big
func isBuffSizeOk(size string, logger *zap.Logger) bool {
	i, err := strconv.Atoi(size)
	if err != nil {
		logger.Error("Buffer size is not a number", zap.Error(err))
		return false
	}
	if i < 1 || i > 65535 {
		logger.Error("Buffer size is either negative or too big", zap.Error(err))
		return false
	}
	if i > 1024 {
		logger.Warn("Buffer size is over 1024, check memory usage to avoid surprises")
	}
	return true
}
