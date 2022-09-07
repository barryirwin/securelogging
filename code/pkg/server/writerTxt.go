package slogserver

import (
	"os"
	"strconv"

	"go.uber.org/zap"
)

// writerTxt :
//
// Writes the ServerData structs into a file without any encryption
func writerTxt(file string, data []ServerData, logger *zap.Logger) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		logger.Error("Failed to open the file: ", zap.Error(err))
		return err
	}
	// Sort by time of arrival (Done in serverDataProcessor)
	// https://stackoverflow.com/a/57095730
	//sort.Slice(data, func(i, j int) bool { return data[i].WireData.ReceivedAt.Before(data[j].WireData.ReceivedAt) })

	// Store data
	for _, item := range data {
		// timestamp, sourceIP, suspicious, data
		str := item.WireData.ReceivedAt.Format("2006-01-02T15:04:05.0000-07")
		str = str + "," + item.WireData.SrcIP
		str = str + "," + strconv.FormatBool(item.SuspiciousIn)
		str = str + "," + strconv.FormatBool(item.SuspiciousRe)
		str = str + "," + item.DecryptedData + "\n"
		f.WriteString(str)
	}
	logger.Info("Text data written", zap.Int("messages", len(data)))
	return nil
}
