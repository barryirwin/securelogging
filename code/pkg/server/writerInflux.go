package slogserver

import (
	slogconfig "github.com/FrancoLoyola/noroff-fdp/code/pkg/config"
	slogdb "github.com/FrancoLoyola/noroff-fdp/code/pkg/influx"
	"go.uber.org/zap"
)

// writerInflux :
//
// Generates all the Influx points using the ServerData.ToInflux().
// Then sends the data to the DB in the conf struct.
func writerInflux(data []ServerData, conf slogconfig.ServerConfig, logger *zap.Logger) error {
	// Create client and batch to write
	dbClient, err := slogdb.CreateInfluxClient(conf.InfluxIP, logger)
	if err != nil {
		logger.Error("Error creating the Influx DB client", zap.Error(err))
		return err
	}
	batch, err := slogdb.CreateInfluxBatch(conf.InfluxDBname, "ns", logger)
	if err != nil {
		logger.Error("Error creating Influx batch", zap.Error(err))
		return err
	}

	if !conf.StoreToInflux {
		logger.Info("Influx writing is not enabled, skipping Influx write")
		return nil
	}
	// Get all points from data
	replay := conf.ReplayMode
	for _, item := range data {
		p, err := item.ToInfluxPoint(replay)
		if err != nil {
			logger.Warn("Failed to create Influx point", zap.Error(err))
		}
		batch.AddPoint(p)
	}
	logger.Info("Will attempt to write", zap.Int("datapoints", len(batch.Points())))
	// Write
	err = slogdb.WriteToInflux(dbClient, batch, logger)
	if err != nil {
		logger.Error("Error writing the Influx batch", zap.Error(err))
		return err
	}

	logger.Info("Successfully written batch to Influx")
	return nil
}
