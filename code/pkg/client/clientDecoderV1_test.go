package slogclient

import (
	"os"
	"strings"
	"testing"

	slogcrypto "github.com/FrancoLoyola/noroff-fdp/code/pkg/crypto"
	logging "github.com/FrancoLoyola/noroff-fdp/code/pkg/logging"
)

func TestClientDecoderV1(t *testing.T) {

	logger, err := logging.NewInfoLogger()
	if err != nil {
		t.Error("Failed to create logger")
	}

	// Scaffolding
	pubPath := "unit-test.pub"
	privPath := "unit-test"
	logline := "This is the cleartext"
	logfile := "testfile"

	err = slogcrypto.GenerateRSAkeys(privPath, pubPath, 2048, logger)
	if err != nil {
		t.Error("Failed to generate temporary keys", err.Error())
	}
	// Cleanup
	defer os.Remove(pubPath)
	defer os.Remove(privPath)

	// Load keys
	privKey, err := slogcrypto.ReadRsaPriv(privPath, logger)
	if err != nil {
		t.Error("Failed to load the private key", err.Error())
	}
	pubKey, err := slogcrypto.ReadRsaPub(pubPath, logger)
	if err != nil {
		t.Error("Failed to load the public key", err.Error())
	}

	// Encode
	encodedLine, err := ClientEncoderV1(logline, logfile, &pubKey, logger)
	if err != nil {
		t.Error("Failed to encode the line", err.Error())
	}

	// Test 1 : Something valid
	// Decode
	decodedLine, fishy, err := ClientDecoderV1(encodedLine, &privKey, logger)
	if err != nil {
		t.Error("Failed to decode the line", err.Error())
	}
	if fishy {
		t.Error("Correct payload flagged as fishy")
	}
	// Roundtrip ok?
	elements := strings.Split(decodedLine, ",")
	ts := elements[0]
	host := elements[1]
	file := elements[2]
	checksum := elements[3]
	decodedLog := elements[4]

	if len(elements) != 5 {
		t.Error("Unexpected elements length", len(elements))
	}
	if len(ts) == 0 {
		t.Error("Timestamp is empty")
	}
	if len(host) == 0 {
		t.Error("Host is empty")
	}
	if file != logfile {
		t.Error("Logfile and file do not match", logfile, file)
	}
	_, chksum, _ := slogcrypto.GenerateSHA256Hash(logline, logger)
	if checksum != chksum {
		t.Error("Checksums dont match", checksum, chksum)
	}
	if decodedLog != logline {
		t.Error("Log lines do not match", logline, decodedLog)
	}

	// Test 2 : Missing one part
	e := strings.Split(encodedLine, ",")
	//t.Error(e)
	payload := e[1]
	decodedLine, fishy, err = ClientDecoderV1(payload, &privKey, logger)
	if !fishy {
		t.Error("Invalid payload not flagged as fishy")
	}

	// Test 3 : Invalid aesBlob
	payload = "asdadasdasda," + e[1]
	decodedLine, fishy, err = ClientDecoderV1(payload, &privKey, logger)
	if err == nil {
		t.Error("Invalid aesBlob should return an error")
	}
	if !fishy {
		t.Error("Invalid payload not flagged as fishy")
	}

	// Test 4 : Invalid RSA
	payload = e[0] + ",asdadasdasda"
	decodedLine, fishy, err = ClientDecoderV1(payload, &privKey, logger)
	if err == nil {
		t.Error("Invalid RSA should return an error")
	}
	if !fishy {
		t.Error("Invalid payload not flagged as fishy")
	}

}
