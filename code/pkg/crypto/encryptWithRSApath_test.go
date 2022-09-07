package slogcrypto

import (
	"os"
	"testing"

	logging "github.com/FrancoLoyola/noroff-fdp/code/pkg/logging"
)

func TestEncryptWithRSApath(t *testing.T) {

	logger, err := logging.NewInfoLogger()
	if err != nil {
		t.Error("Failed to create logger")
	}

	// TEST 1 : Valid key path, created for this unit test
	pubKey := "unit-test.pub"
	privKey := "unit-test"
	//pubKey := "../../../build/keys/comms/slog-comms_rsa.pub"
	//privKey := "../../../build/keys/comms/slog-comms_rsa"
	// SHA-256 hash for "my-fancy-password"
	cleartext := "56b4fbe06471a2a378bb654e94925c39ed04bf70846b7632441500f8c3c4fa5e"

	err = GenerateRSAkeys(privKey, pubKey, 2048, logger)
	if err != nil {
		t.Error("Failed to generate temporary keys", err.Error())
	}

	ciphertext, err := EncryptWithRSApath(cleartext, pubKey, logger)
	if err != nil {
		t.Error("Failed to encrypt cleartext with key" + err.Error())
	}
	if ciphertext == cleartext {
		t.Error("Ciphertext and cleartext are the same: ", ciphertext)
	}

	result, err := DecryptWithRSApath(ciphertext, privKey, logger)
	if err != nil {
		t.Error(result)
		t.Error(err)
	}

	if result == ciphertext {
		t.Error("The result of the decryption is the same as the ciphertext: ", ciphertext)
	}
	if result != cleartext {
		t.Error("The result of the decryption is not the same as the cleartext: ", result, cleartext)
	}
	// Cleanup
	os.Remove(pubKey)
	os.Remove(privKey)
}
