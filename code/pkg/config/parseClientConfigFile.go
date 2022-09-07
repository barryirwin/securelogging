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

// ParseClientConfigFile :
//
// Attempts to parse the given config file, it crashes the program if some fields are wrong
func ParseClientConfigFile(filePath string, logger *zap.Logger) (ClientConfig, error) {
	conf := NewClientConfig()
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
		case "slogServerIP":
			// Can't check, could be a DNS valid
			conf.ServerIP = l[1]
		case "slogServerPort":
			if !checks.IsPortValid(l[1], logger) {
				logger.Panic("Port is not valid")
				break
			}
			conf.Port = l[1]
		case "slogServerProtocol":
			if !checks.IsProtocolValid(l[1], logger) {
				logger.Panic("slogServerProtocol is not valid")
				break
			}
			conf.Protocol = strings.ToUpper(l[1])
		case "publicCommsKey":
			if len(l[1]) < 1 {
				logger.Panic("Comms public key path is empty, can't encrypt data...")
				break
			}
			if !checks.FileExists(l[1], logger) {
				logger.Panic("Comms public key path doesn't, can't encrypt data...")
				break
			}
			conf.CommsPubKey = l[1]
		case "tailFiles":
			if len(l[1]) < 1 {
				logger.Panic("tailFiles is empty, nothing to do...")
				break
			}
			list := strings.Split(l[1], ",")
			conf.FileList = list
		case "bufferLines":
			if !isBuffSizeOk(l[1], logger) {
				logger.Error("bufferLines invalid, defaulting to 10")
				conf.BufferLines = 10
				break
			}
			// Don't check for errors here, as this conversion is done and checked at isBuffSizeOk
			i, _ := strconv.Atoi(l[1])
			conf.BufferLines = i
		}
	}
	return conf, nil
}
