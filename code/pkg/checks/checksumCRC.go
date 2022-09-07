package checks

import (
	"hash/crc32"
)

// CalculateCRC :
//
// Calculates the crc32 of the given input, meant to be used to compare against a received crc.
// Uses crc32.ChecksumIEEE(in)
func CalculateCRC(in []byte) uint32 {
	out := crc32.ChecksumIEEE(in)
	return out
}
