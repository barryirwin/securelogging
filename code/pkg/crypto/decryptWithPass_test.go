package slogcrypto

import (
	"testing"

	logging "github.com/FrancoLoyola/noroff-fdp/code/pkg/logging"
)

func TestDecryptWithPass(t *testing.T) {
	logger, err := logging.NewInfoLogger()
	if err != nil {
		t.Error("Failed to create logger")
	}

	test := "f38d3b88846d6a392ae897eb736dcc0cb10c7c377bf1e13980a2fd13445d3b7c4399e417a65969596fde"
	pass := "my-fancy-password"
	want := "This is a test"

	// Decrypt the data
	plaintext, err := DecryptWithPass(test, pass, logger)
	if err != nil {
		t.Error("Failed to decrypt the data")
	}
	if plaintext != want {
		t.Error(plaintext, want)
		t.Error("Decrypted string is different from the test string")
	}
}
