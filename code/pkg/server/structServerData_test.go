package slogserver

import (
	"testing"
	"time"
)

func TestToBase64(t *testing.T) {
	wd := NewWireData("some-plaintext-data", "127.0.0.1", "TCP", false)
	sd := NewServerData(wd)

	// Main functionality tests
	b64, chksum, err := sd.ToBase64()

	if err != nil {
		t.Error("Failed to encode to Base64")
	}
	if chksum == 0 {
		t.Error("Checksum is 0?")
	}
	if len(b64) == 0 {
		t.Error("Base64 Encoded data is 0")
	}

	// Update the struct with some more data (encryption/decryption unit test)
	// and then check how fast this can encode 10k packets, should be under 1 second if possible
	sd.WireData.Data = "f38d3b88846d6a392ae897eb736dcc0cb10c7c377bf1e13980a2fd13445d3b7c4399e417a65969596fde"
	sd.DecryptedData = "This is a test"

	start := time.Now()
	for i := 1; i < 10000; i++ {
		sd.ToBase64()
	}
	end := time.Now()
	total := end.Sub(start)

	if total.Seconds() > 1 {
		t.Error("Encoding 10k packets took longer than a second: ", total.Seconds(), "s")
	}

}

func TestFromBase64(t *testing.T) {
	// String representation of the same basic struct as "TestToBase64"
	s := "eyJ3aXJlLWRhdGEiOnsiZGF0YSI6InNvbWUtcGxhaW50ZXh0LWRhdGEiLCJzb3VyY2UtaXAiOiIxMjcuMC4wLjEiLCJyZWNlaXZlZC1hdCI6IjIwMjEtMDItMDZUMTU6MjA6NTAuODEzMzEzKzAxOjAwIiwic3lzbG9nLWRhdGEiOmZhbHNlLCJwcm90b2NvbCI6IlRDUCJ9LCJkZWNyeXB0ZWQtZGF0YSI6IiIsInN1c3BpY2lvdXMtaW4iOmZhbHNlLCJzdXNwaWNpb3VzLXJlIjpmYWxzZSwic2VydmVyLWhvc3RuYW1lIjoiRnJhbmNvcy1NYWNCb29rLUFpci5sb2NhbCIsInByZXZpb3VzLWhhc2giOiIifQo="
	ip := "127.0.0.1"
	wData := "some-plaintext-data"

	sd, err := ServerDataFromBase64(s)
	if err != nil {
		t.Error("Decoding failed")
	}

	// Checks
	if sd.SuspiciousIn {
		t.Error("Default constructor for ServerData doesn't flag data as suspicious")
	}
	if sd.WireData.SrcIP != ip {
		t.Error("Sample data IP doesn't match. Got: ", sd.WireData.SrcIP, " want: ", ip)
	}
	if sd.WireData.Data != wData {
		t.Error("WireData.Data within the ServerData doesn't match. Got: ", sd.WireData.Data, " got: ", wData)
	}

	// String representation of the "long" struct in the test
	s = "eyJ3aXJlLWRhdGEiOnsiZGF0YSI6ImYzOGQzYjg4ODQ2ZDZhMzkyYWU4OTdlYjczNmRjYzBjYjEwYzdjMzc3YmYxZTEzOTgwYTJmZDEzNDQ1ZDNiN2M0Mzk5ZTQxN2E2NTk2OTU5NmZkZSIsInNvdXJjZS1pcCI6IjEyNy4wLjAuMSIsInJlY2VpdmVkLWF0IjoiMjAyMC0xMi0yNlQxMzo1ODozMi44MzY0NjYzMzgrMDE6MDAifSwiZGVjcnlwdGVkLWRhdGEiOiJUaGlzIGlzIGEgdGVzdCIsInN1c3BpY2lvdXMiOmZhbHNlLCJzZXJ2ZXItaG9zdG5hbWUiOiJwb3Atb3MifQo="

	start := time.Now()
	for i := 1; i < 10000; i++ {
		tmp, err := ServerDataFromBase64(s)
		if err != nil {
			t.Error("Decoding error in the 10k test", err, tmp)
		}
	}
	end := time.Now()
	total := end.Sub(start)

	if total.Seconds() > 1 {
		t.Error("Decoding 10k packets took longer than a second: ", total.Seconds(), "s")
	}
}
