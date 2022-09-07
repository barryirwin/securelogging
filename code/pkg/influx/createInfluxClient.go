package slogdb

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"go.uber.org/zap"
)

// CreateInfluxClient :
//
// Creates an InfluxDB client based on the string. Expected string format: "http://XXX.XXX.XXX.XXX:XXXX"
// Has a defer client.Close()
func CreateInfluxClient(ip string, logger *zap.Logger) (client.Client, error) {

	client, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: ip,
	})
	if err != nil {
		logger.Error("Error creating Influx DB clients:", zap.Error(err))
		return nil, err
	}
	defer client.Close()

	return client, nil
}
