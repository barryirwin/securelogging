package slognet

import (
	"net"

	"github.com/FrancoLoyola/noroff-fdp/code/pkg/checks"
	slogserver "github.com/FrancoLoyola/noroff-fdp/code/pkg/server"
	"go.uber.org/zap"
)

// TCPListener :
//
// Listens for data on a TCP given port.
// Creates a go routine for each incoming connection, never quits if port is binded
func TCPListener(port string, buffSize int, syslog bool, chanWireData chan slogserver.WireData, logger *zap.Logger) {
	// https://coderwall.com/p/wohavg/creating-a-simple-tcp-server-in-go and https://opensource.com/article/18/5/building-concurrent-tcp-server-go
	if len(port) < 0 || !checks.IsPortValid(port, logger) {
		logger.Error("TCP Port is wrong")
		return
	}

	portstr := ":" + port
	l, err := net.Listen("tcp4", portstr)
	if err != nil {
		logger.Fatal("Failed attempting to create the network listener", zap.String("port", port), zap.Error(err))
		return
	}
	defer l.Close()
	logger.Info("Listening on: ", zap.String("port", port))

	// Listen for incoming connections
	n := 0
	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Error("Error accepting connection", zap.String("addr", conn.RemoteAddr().String()), zap.Error(err))
		}
		n++
		// Handle connections in parallel
		go slogserver.ParseTCPRequest(conn, buffSize, syslog, chanWireData, logger)
		if n%1000 == 0 {
			logger.Info("Received 1000 TCP packets", zap.Int("Total", n))
		}
	}
}
