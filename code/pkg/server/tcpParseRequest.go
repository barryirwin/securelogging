package slogserver

import (
	"net"
	"strings"
	"time"

	slogstring "github.com/FrancoLoyola/noroff-fdp/code/pkg/string"

	"go.uber.org/zap"
)

// ParseTCPRequest :
//
// Read the data received on the connection, then closes it. Or after 5 seconds, regardless of the result
// Data is transformed into a WireData struct and sent to chanWireData
func ParseTCPRequest(conn net.Conn, buffSize int, syslog bool, chanWireData chan WireData, logger *zap.Logger) {
	// https://forum.golangbridge.org/t/tcp-server-to-read-an-unknown-number-of-bytes-but-act-on-it-as-they-come/8661

	// Always close the connection, no matter what, after max 1 seconds, don't want open connections lying around
	err := conn.SetDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		logger.Error("Error setting deadline for the connection", zap.Error(err))
	}
	defer conn.Close()

	// Reference vars
	okResponse := []byte("OK") // Always ok, not to give more hints
	bytesReceived := make([]byte, buffSize)

	addr := conn.RemoteAddr().String()
	tmp := strings.Split(addr, ":")
	addr = tmp[0]

	len, err := conn.Read(bytesReceived)
	if err != nil {
		logger.Warn("Error reading the data from the connection", zap.String("addr", addr), zap.Error(err))
		return
	}
	dataReceived := string(bytesReceived[:len])
	// No spaces allowed all should be a single chunk when received
	dataReceived = strings.TrimSpace(dataReceived)
	if !syslog {
		dataReceived = slogstring.TrimWhitespace(dataReceived)
	}

	conn.Write(okResponse)

	// Send the data back
	chanWireData <- NewWireData(dataReceived, addr, "TCP", syslog)
}
