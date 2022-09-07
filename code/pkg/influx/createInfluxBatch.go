package slogdb

import (
	"go.uber.org/zap"

	client "github.com/influxdata/influxdb1-client/v2"
)

// CreateInfluxBatch :
//
// Returns a batch for adding data points to it for later storage into Influx
func CreateInfluxBatch(dbName, precision string, logger *zap.Logger) (client.BatchPoints, error) {
	newBatch, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  dbName,
		Precision: precision,
	})
	if err != nil {
		logger.Error("Error creating the Influx batch:", zap.Error(err))
		return nil, err
	}
	return newBatch, nil
}
