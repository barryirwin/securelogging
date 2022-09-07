package slogconfig

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/FrancoLoyola/noroff-fdp/code/pkg/checks"
	slogstring "github.com/FrancoLoyola/noroff-fdp/code/pkg/string"

	"go.uber.org/zap"
)

// ParseServerConfigFile :
//
// Attempts to parse the given config file
func ParseServerConfigFile(filePath string, logger *zap.Logger) (ServerConfig, error) {
	conf := NewServerConfig()
	file, err := os.Open(filePath)
	if err != nil {
		logger.Error("Failed to open the file: ", zap.Error(err))
		return conf, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = slogstring.TrimWhitespace(line)
		// Skip comments, empty lines and invalid lines
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.Contains(line, "=") || strings.Count(line, "=") > 1 {
			logger.Warn("Invalid line in config file: ", zap.String("->", line))
			continue
		}
		// Parsing different arguments
		l := strings.Split(line, "=")
		switch l[0] {
		case "TCPport":
			if !checks.IsPortValid(l[1], logger) {
				logger.Error("TCP Port is not valid")
				break
			}
			conf.TCPport = l[1]
		case "UDPport":
			if !checks.IsPortValid(l[1], logger) {
				logger.Error("UDP Port is not valid")
				break
			}
			conf.UDPport = l[1]
		case "UDPsyslog":
			if !checks.IsPortValid(l[1], logger) {
				logger.Error("Syslog UDP Port is not valid")
				break
			}
			conf.UDPsyslog = l[1]
		case "TCPsyslog":
			if !checks.IsPortValid(l[1], logger) {
				logger.Error("Syslog TCP Port is not valid")
				break
			}
			conf.TCPsyslog = l[1]
		case "CommsPrivKey":
			if len(l[1]) < 1 {
				logger.Error("Comms private key path is empty")
				break
			}
			conf.CommsPrivKey = l[1]
		case "StoragePubKey":
			if len(l[1]) < 1 {
				logger.Error("Storage public key path is empty")
				break
			}
			conf.StoragePubKey = l[1]
		case "StorageFolder":
			if len(l[1]) < 1 {
				logger.Error("Storage folder path is empty, defaulting to '.'")
				conf.StorageFolder = "."
				break
			}
			lastByte := l[1][len(l[1])-1]
			if lastByte != byte('/') {
				logger.Warn("Appending / to the storage folder path")
				l[1] = l[1] + "/"
			}
			conf.StorageFolder = l[1]
		case "InfluxDBname":
			if len(l[1]) < 1 {
				logger.Error("Influx database name is emtpy, defaulting to 'test'")
				conf.InfluxDBname = "test"
				break
			}
			conf.InfluxDBname = l[1]
		case "InfluxIP":
			if len(l[1]) < 1 {
				logger.Error("Influx IP:Port is emtpy, defaulting to 'localhost:8086'")
				conf.InfluxIP = "localhost:8086"
				break
			}
			conf.InfluxIP = l[1]
		case "PacketBufferSize":
			if !isBuffSizeOk(l[1], logger) {
				logger.Error("PacketBufferSize invalid, defaulting to 1024")
				conf.PacketBufferSize = 1024
				break
			}
			// Don't check for errors here, as this conversion is done and checked at isBuffSizeOk
			i, _ := strconv.Atoi(l[1])
			conf.PacketBufferSize = i
		case "FileBatchSize":
			if !isBuffSizeOk(l[1], logger) {
				logger.Error("FileBatchSize invalid, defaulting to 1000")
				conf.FileBatchSize = 1000
				break
			}
			// Don't check for errors here, as this conversion is done and checked at isBuffSizeOk
			i, _ := strconv.Atoi(l[1])
			conf.FileBatchSize = i
		case "WireChannelBuffer":
			if !isBuffSizeOk(l[1], logger) {
				logger.Error("WireChannelBuffer invalid, defaulting to 100")
				conf.WireChannelBuffer = 100
				break
			}
			// Don't check for errors here, as this conversion is done and checked at isBuffSizeOk
			i, _ := strconv.Atoi(l[1])
			conf.WireChannelBuffer = i
		case "ProcessingChannelBuffer":
			if !isBuffSizeOk(l[1], logger) {
				logger.Error("ProcessingChannelBuffer invalid, defaulting to 100")
				conf.ProcessingChannelBuffer = 100
				break
			}
			// Don't check for errors here, as this conversion is done and checked at isBuffSizeOk
			i, _ := strconv.Atoi(l[1])
			conf.ProcessingChannelBuffer = i
		case "CompressedBlob":
			enabled, err := strconv.ParseBool(l[1])
			if err != nil {
				logger.Error("CompressedBlob invalid, defaulting to false")
				conf.CompressedBlob = false
				break
			}
			conf.CompressedBlob = enabled
		}
	}
	return conf, nil
}
