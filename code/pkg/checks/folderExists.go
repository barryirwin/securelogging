package checks

import (
	"os"

	"go.uber.org/zap"
)

// FolderExists :
//
// Checks if a folder exists
func FolderExists(path string, logger *zap.Logger) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		logger.Warn("Folder not found:", zap.Error(err))
		return false
	}
	return info.IsDir()
}
