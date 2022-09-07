package checks

import (
	"os"

	"go.uber.org/zap"
)

// FileExists :
//
// Checks if a file exists and is not a directory before we try using it to prevent further errors.
func FileExists(filename string, logger *zap.Logger) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		logger.Error("File not found:", zap.Error(err))
		return false
	}
	return !info.IsDir()
}
