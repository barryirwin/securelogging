package slogcrypto

import (
	"testing"

	logging "github.com/FrancoLoyola/noroff-fdp/code/pkg/logging"
)

func TestGenerateSHA256Hash(t *testing.T) {
	// This test only verifies against a valid output of SHA256, if the hashing changes, update.
	pass := "my-fancy-password"
	logger, err := logging.NewInfoLogger()
	if err != nil {
		t.Error("Failed to create logger")
	}
	wantb := []byte{86, 180, 251, 224, 100, 113, 162, 163, 120, 187, 101, 78, 148, 146, 92, 57, 237, 4, 191, 112, 132, 107, 118, 50, 68, 21, 0, 248, 195, 196, 250, 94}
	wants := "56b4fbe06471a2a378bb654e94925c39ed04bf70846b7632441500f8c3c4fa5e"

	gotb, gots, err := GenerateSHA256Hash(pass, logger)
	if err != nil {
		t.Error("This should not fail...")
	}
	if wants != gots {
		t.Error("The returned string doesn't match what was expected")
	}
	// Compare that they match byte by byte
	for pos, bite := range wantb {
		if bite != gotb[pos] {
			t.Errorf("Missmatch in the got %b should be %b", gotb[pos], bite)
		}
	}

	// Empty stuff
	pass = ""
	wantb = []byte{227, 176, 196, 66, 152, 252, 28, 20, 154, 251, 244, 200, 153, 111, 185, 36, 39, 174, 65, 228, 100, 155, 147, 76, 164, 149, 153, 27, 120, 82, 184, 85}
	wants = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	gotb, gots, err = GenerateSHA256Hash(pass, logger)
	if err != nil {
		t.Error("This should not fail...")
	}
	if wants != gots {
		t.Error("The returned string doesn't match what was expected")
		t.Error(wants, gots)
	}
	// Compare that they match byte by byte
	for pos, bite := range wantb {
		if bite != gotb[pos] {
			t.Errorf("Missmatch in the got %b should be %b", gotb[pos], bite)
		}
	}
}
