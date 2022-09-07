package slogfile

import (
	"io"
	"os"

	"go.uber.org/zap"
)

// GetLastLine :
//
// Seeks the end of the file and reads backwards byte by byte until a newline char is found.
// If the file doesn't exist, returns an empty string
// https://stackoverflow.com/a/51328256
func GetLastLine(filepath string, logger *zap.Logger) string {

	fileHandle, err := os.Open(filepath)
	if err != nil {
		logger.Warn("Couldn't open the file, probably doesn't exist yet", zap.Error(err))
		return ""
	}
	defer fileHandle.Close()

	line := ""
	var cursor int64 = 0
	stat, _ := fileHandle.Stat()
	filesize := stat.Size()

	if filesize == 0 {
		logger.Info("Looks like this is an empty file")
		return line
	}

	// Reading backwards
	for {
		cursor--
		fileHandle.Seek(cursor, io.SeekEnd)

		char := make([]byte, 1)
		fileHandle.Read(char)

		// Stop if we find a new line
		if cursor != -1 && (char[0] == 10 || char[0] == 13) {
			break
		}

		line = string(char) + line

		// Stop if we are at the begining
		if cursor == -filesize {
			break
		}
	}
	return line
}
