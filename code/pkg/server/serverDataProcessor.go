package slogserver

import (
	"os"
	"sort"

	"github.com/FrancoLoyola/noroff-fdp/code/pkg/checks"
	slogconfig "github.com/FrancoLoyola/noroff-fdp/code/pkg/config"
	"go.uber.org/zap"
)

// ServerDataProcessor :
//
// Reads the already decrypted server data, sorts them when a certain limit is reached
// Coordinates the storage and output to csv and Influx
func ServerDataProcessor(chanServerData chan ServerData, conf slogconfig.ServerConfig, logger *zap.Logger) error {
	slice := []ServerData{}
	// These maybe could be parameter based?
	filename := conf.StorageFolder + "slog-logging"
	version := 1

	// Verify path, Create it if required
	if !checks.FolderExists(conf.StorageFolder, logger) {
		logger.Warn("StorageFolder doesn't exist, creating it...")
		err := os.MkdirAll(conf.StorageFolder, 0755)
		if err != nil {
			logger.Error("Could not create the required folder(s) for StorageFolder", zap.Error(err))
			return err
		}
	}
	// The processing
	for {
		select {
		case d := <-chanServerData:
			slice = append(slice, d)
			// Process the batch when the size is reached, maybe sort the batch before storing it.
			if len(slice) == conf.FileBatchSize {
				logger.Info("Filled a new batch", zap.Int("size", len(slice)))
				// Sort by time of arrival
				// https://stackoverflow.com/a/57095730
				sort.Slice(slice, func(i, j int) bool { return slice[i].WireData.ReceivedAt.Before(slice[j].WireData.ReceivedAt) })
				// Write
				go writerTxt(filename+".txt", slice, logger)
				go writerEncrypted(filename+".slog", slice, version, conf, logger)
				go writerInflux(slice, conf, logger)
				// Reset the slice for the new chunk
				slice = []ServerData{}
			}
		}
	}
}
