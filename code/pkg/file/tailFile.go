package slogfile

import (
	"github.com/hpcloud/tail"
	"go.uber.org/zap"
)

// TailFile :
//
// Follows a file forever. Writes the content to the out channel.
func TailFile(path string, out chan string, logger *zap.Logger) {

	t, err := tail.TailFile(path, tail.Config{Follow: true, ReOpen: true, Poll: true})
	if err != nil {
		logger.Error("Error trying to open: "+path, zap.Error(err))
	}
	defer t.Cleanup()

	for line := range t.Lines {
		out <- line.Text
	}
}
