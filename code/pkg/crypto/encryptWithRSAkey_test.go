package slogcrypto

import (
	"os"
	"testing"

	logging "github.com/FrancoLoyola/noroff-fdp/code/pkg/logging"
)

func TestEncryptWithRSAkey(t *testing.T) {

	logger, err := logging.NewInfoLogger()
	if err != nil {
		t.Error("Failed to create logger")
	}

	// TEST 1 : Valid key path, created for this unit test
	pubPath := "unit-test.pub"
	privPath := "unit-test"
	cleartext := "56b4fbe06471a2a378bb654e94925c39ed04bf70846b7632441500f8c3c4fa5e"

	err = GenerateRSAkeys(privPath, pubPath, 2048, logger)
	if err != nil {
		t.Error("Failed to generate temporary keys", err.Error())
	}
	// Cleanup
	defer os.Remove(pubPath)
	defer os.Remove(privPath)

	privKey, err := ReadRsaPriv(privPath, logger)
	if err != nil {
		t.Error("Failed to read private key", err.Error())
	}
	pubKey, err := ReadRsaPub(pubPath, logger)
	if err != nil {
		t.Error("Failed to read public key", err.Error())
	}

	// Encrypt
	ciphertext, err := EncryptWithRSAkey(cleartext, pubKey, logger)
	if err != nil {
		t.Error("Failed to encrypt cleartext with key" + err.Error())
	}

	if ciphertext == cleartext {
		t.Error("Ciphertext and cleartext are the same: ", ciphertext)
	}

	// Decrypt
	result, err := DecryptWithRSAkey(ciphertext, privKey, logger)
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
}
