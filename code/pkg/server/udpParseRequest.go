package slogserver

import (
	"net"
	"strings"

	slogstring "github.com/FrancoLoyola/noroff-fdp/code/pkg/string"
	"go.uber.org/zap"
)

// ParseUDPrequest :
//
// Parses half of the work required for UDP processing.
func ParseUDPrequest(packetLength int, bytesReceived []byte, addr net.Addr, syslog bool, chanWireData chan WireData, logger *zap.Logger) {

	dataReceived := string(bytesReceived[:packetLength])

	// No spaces allowed all should be a single chunk when received
	dataReceived = strings.TrimSpace(dataReceived)
	if !syslog {
		dataReceived = slogstring.TrimWhitespace(dataReceived)
	}

	tmp := strings.Split(addr.String(), ":")
	addrNoPort := tmp[0]

	// Handle connections in parallel, send the data to the processing go routines and keep listening
	chanWireData <- NewWireData(dataReceived, addrNoPort, "UDP", syslog)
}
