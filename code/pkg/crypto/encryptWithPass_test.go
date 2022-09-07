package slogcrypto

import (
	"testing"

	logging "github.com/FrancoLoyola/noroff-fdp/code/pkg/logging"
)

func TestEncryptWithPass(t *testing.T) {
	logger, err := logging.NewInfoLogger()
	if err != nil {
		t.Error("Failed to create logger")
	}

	test := "This is a test"
	// "my-fancy-password" SHA-256 hash
	pass := "56b4fbe06471a2a378bb654e94925c39ed04bf70846b7632441500f8c3c4fa5e"

	// Each run will generate a different output, but the length should be the same
	ciphertext, err := EncryptWithPass(test, pass, logger)
	if err != nil {
		t.Error("Encryption should not fail")
	}

	// To be sure, try to decrypt the data
	plaintext, err := DecryptWithPass(ciphertext, pass, logger)
	if err != nil {
		t.Error("Failed to decrypt the data")
	}
	if plaintext != test {
		t.Error("Decrypted string is different from the test string")
	}
}
