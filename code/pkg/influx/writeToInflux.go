package slogdb

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"go.uber.org/zap"
)

// WriteToInflux :
//
// Writes to an Influx client a batch of data points, closes the connection to the DB
func WriteToInflux(c client.Client, b client.BatchPoints, logger *zap.Logger) error {
	// Write the batch
	if err := c.Write(b); err != nil {
		logger.Error("Error writing to Influx:", zap.Error(err))
		return err
	}
	// Close client resources
	if err := c.Close(); err != nil {
		logger.Error("Error closing Influx:", zap.Error(err))
		return err
	}
	return nil
}
