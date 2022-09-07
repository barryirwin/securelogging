package slogcrypto

import (
	"testing"
	"time"
)

func TestGenerateRandomPass(t *testing.T) {
	// Verify that we can get at least 10k different passwords of 32 bytes
	// Which is the required length for AES-256
	// (not really, slogcrypt uses a hash of the pass, but to have a "standard length")
	passSlice := []string{}
	nPass := 10000
	secondsLimit := 1.5

	// "Finder"
	find := func(slice []string, val string) (int, bool) {
		for i, item := range slice {
			if item == val {
				return i, true
			}
		}
		return -1, false
	}

	// Generate
	start := time.Now()
	for i := 1; i < nPass; i++ {
		tmp := GenerateRandomPass(32)
		passSlice = append(passSlice, tmp)
	}
	end := time.Now()
	total := end.Sub(start)
	if total.Seconds() > secondsLimit {
		t.Error("Password generation was too slow", total.Seconds(), "s")
	}

	// Verify that there are no duplicates
	for index, pass := range passSlice {
		passSlice[index] = ""
		_, found := find(passSlice, pass)
		if found {
			t.Error("Found a repeated password, not that random")
		}
	}
}
