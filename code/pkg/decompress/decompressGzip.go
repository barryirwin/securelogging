package decompress

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"

	"go.uber.org/zap"
)

// Gzip :
//
// Decompress the given base64 encoded string using gzip.
// If the decompression fails, it returns the string as provided.
//
// https://play.golang.org/p/dsakzksaP1k
func Gzip(in string, logger *zap.Logger) (string, error) {
	// Base 64 part
	b64, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		logger.Error("Failed to decode the base64 string", zap.Error(err))
		return in, err
	}
	rb64 := bytes.NewReader(b64)

	// Decompression part
	gz, err := gzip.NewReader(rb64)
	if err != nil {
		logger.Error("Failed to initialize gzip reader", zap.Error(err))
		return in, err
	}
	data, err := ioutil.ReadAll(gz)
	if err != nil {
		logger.Warn("Failed to read gzip data", zap.Error(err))
		return in, err
	}
	out := string(data)
	return out, nil
}
