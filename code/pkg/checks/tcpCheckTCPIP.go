package checks

import "net"

// CheckTCPIP :
//
// Checks if the provided address is valid
func CheckTCPIP(ip string) error {
	_, err := net.ResolveTCPAddr("tcp", ip)
	if err != nil {
		return err
	}
	return nil
}
