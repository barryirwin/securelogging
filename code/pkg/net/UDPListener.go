package slognet

import (
	"net"

	"github.com/FrancoLoyola/noroff-fdp/code/pkg/checks"

	slogserver "github.com/FrancoLoyola/noroff-fdp/code/pkg/server"
	"go.uber.org/zap"
)

// UDPListener :
//
// Listens for data on a UDP given port.
// Never quits if the port is binded succesfully, only prints the errors.
// Data is transformed into a WireData struct
func UDPListener(port string, buffSize int, syslog bool, chanWireData chan slogserver.WireData, logger *zap.Logger) {
	// https://ops.tips/blog/udp-client-and-server-in-go/#receiving-from-a-udp-connection-in-a-server
	if len(port) < 0 || !checks.IsPortValid(port, logger) {
		logger.Error("UDP Port is wrong")
		return
	}

	portstr := ":" + port
	l, err := net.ListenPacket("udp", portstr)
	if err != nil {
		logger.Fatal("Failed attempting to create the network listener", zap.String("port", port), zap.Error(err))
		return
	}
	defer l.Close()
	logger.Info("Listening on: ", zap.String("port", port))

	n := 0
	// Listen for incoming connections
	for {
		n++
		if n%1000 == 0 {
			logger.Info("Received 1000 UDP packets", zap.Int("Total", n))
		}
		bytesReceived := make([]byte, buffSize)
		len, addr, err := l.ReadFrom(bytesReceived)
		if err != nil {
			logger.Error("Error accepting connection", zap.String("addr", addr.String()), zap.Error(err))
		}
		// Handle connections in parallel, send the data to the processing go routines and keep listening
		go slogserver.ParseUDPrequest(len, bytesReceived, addr, syslog, chanWireData, logger)
	}
}
