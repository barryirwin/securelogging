package decompress

import (
	"testing"

	logging "github.com/FrancoLoyola/noroff-fdp/code/pkg/logging"
)

func TestDecompressGzip(t *testing.T) {
	// https://play.golang.org/p/dsakzksaP1k
	// See the compress.Gzip as well
	logger, err := logging.NewInfoLogger()
	if err != nil {
		t.Error("Failed to create logger")
	}

	test := "H4sIAAAAAAAA/xSRSREAQQjELHHDyOH0L2Fr3/3oVOJde7XgIsGltBuOatjB3r2Uh3i9XZu9S0WxjR2KDddNl84PEq/p2oysMD3t7di7WKm3VDgtzTpbkLWyXeSIjtnL6VplONOVxWm4R7kd2WmT/5+se8qRxA96l0gOGrW7Dw5qGwufa/DwWRynpJdPaeyim12YLK93QxErx2xUasQgNt3xNrvYcUHlyTt4XCvgPVV6N/vOFXhmQstz9/nwnLXNup8MbWer088gZkJDDA94rS4f552gCLtSM2urdUilx42ngSU/q4AWbYLVlq7z3tAVsqop1Eg/uJt+KEtg91K7FNJpN4kSz67wVKLMOU2BxawUIcUZRg2X9ohz8SC6KA4rFJnn9d5LtMqYfl5XoHpV9o6fVkqg4btm6T9qcusr7gx15eDpecO/qynPVyQif8/1SHvvH3kD5LioAYqf0CNZuXrEbX0KReMnPjDZ7M7zjIiABnLdHi5oYgigc9+A3JFL9bjpDioyzlwVngyTxjrv5QSbX8T0ASi1X41CPHxB0cgziYvhyjsq4dYaXwAAAP//E85wAhADAAA="
	want := "7cbefbe074483b52ee871561c837cce2af11fcecbeacee2b28ec1c851c0fcc2fa7390a1fc2fc6626b1a7a69ed69f8e4b9e2b1dc4c35deb0abe4ecb271171ace3a75bb61ddcbab3a61ef2aec8aca6da2b284e77a4f248efbeffa127058beee90f0bec1b197583d3f68f3a4a7b7db58ee1766f864e3e7cc0b23538de8ba588d11c2ff3ec3746f3b0baf49f093be407cdbb5ffde9f7503ddd85b7aee97d3df6c6de77f4d2ecac572583d46642d230903e6bfa93aff41443752c335c56c84ba78fd7a606a396b80c45c20e5c4cbf7ce85e426bbdb158a7f0ffdc914e206f9a5cb50a72eea22a1f6fb1f548b673a6503466b510a4730d561e2ef23ae1f08cb28f1b85133d9cef9a16ba8dc97bfb055fbb69f395ba481619fc34c9f8ea3c59b3ca8575383dcd9d3dbb5db7a9b2444ddcbe78a699d9d33e804f3b2c00b3942924e4fb923c6cf50b2d7f47d0dac3773d9622202d0ae7691e05a1840173cfd04ff274bcd765ed15131ddfbb1f4d3258e73efad8367f88dcf0052c7fbd508919828c13dda1e18753ed54876c58"

	// Decompress the data
	got, err := Gzip(test, logger)
	if err != nil {
		t.Error("Failed to de-base64 and decompress the data")
	}
	if got != want {
		t.Error(got, want)
		t.Error("Decompressed string is different from the test string")
	}
	if len(test) >= len(got) {
		t.Error("Decompressed string is smaller or equal than the test string, not decompressing much here...")
	}
}
