package slogclient

import (
	"fmt"
	"net"

	"go.uber.org/zap"
)

// UDPWriter :
//
// Writes the given string to the address.
//
// addr must be formatted as: "ip:port"
func UDPWriter(message, addr string, logger *zap.Logger) error {
	if len(message) < 1 {
		logger.Warn("Got an empty string to send...")
		return fmt.Errorf("Empty string, nothing to do")
	}
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		logger.Error("ResolveUDPAddr failed:", zap.Error(err))
		return err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
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
	logger.Info("Sent data to UDP: " + addr)

	return nil
}
