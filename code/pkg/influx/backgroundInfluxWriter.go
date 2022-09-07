package slogdb

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"go.uber.org/zap"
)

// BackgroundInfluxWriter :
//
// Writes to the Influx client the batch of data points arriving through the channel
// until a stop signal is received
func BackgroundInfluxWriter(b chan client.BatchPoints, stop chan bool, c client.Client, logger *zap.Logger) {
	for {
		select {
		// Write
		case batch := <-b:
			if err := c.Write(batch); err != nil {
				logger.Error("Error writing to Influx:", zap.Error(err))
			}
		// Stop
		case s := <-stop:
			if s {
				logger.Info("Stopping...")
				// Close client resources
				if err := c.Close(); err != nil {
					logger.Error("Error closing Influx:", zap.Error(err))
					return
				}
			} else {
				logger.Warn("Received false stop signal...")
			}
		}
	}
}
