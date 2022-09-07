package slogclient

import (
	"fmt"
	"net"

	"go.uber.org/zap"
)

// TCPWriter :
//
// Writes the given string to the address.
//
// addr must be formatted as: "ip:port"
func TCPWriter(message, addr string, logger *zap.Logger) error {
	if len(message) < 1 {
		logger.Warn("Got an empty string to send...")
		return fmt.Errorf("Empty string, nothing to do")
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		logger.Error("ResolveTCPAddr failed:", zap.Error(err))
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		logger.Error("Dial failed:", zap.Error(err))
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(message))
	if err != nil {
		logger.Error("Write to server failed:", zap.Error(err))
		return err
	}

	reply := make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		logger.Error("Write to server failed:", zap.Error(err))
		return err
	}
	logger.Info("Got reply: " + string(reply) + " from: " + addr)

	return nil
}
