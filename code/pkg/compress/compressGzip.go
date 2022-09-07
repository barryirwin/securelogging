package compress

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"

	"go.uber.org/zap"
)

// Gzip :
//
// Compress the given string using gzip, then encoding the binary result as a base64 string.
// If the compression fails, it returns the string as provided.
//
// https://play.golang.org/p/WO42wxyy77K
func Gzip(in string, logger *zap.Logger) (string, error) {

	// Compression part
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(in)); err != nil {
		logger.Error("Failed to write the string gzip", zap.Error(err))
		return in, err
	}
	if err := gz.Close(); err != nil {
		logger.Error("Failed to close the gzip writer", zap.Error(err))
		return in, err
	}

	// Base 64 part
	buff := make([]byte, base64.StdEncoding.EncodedLen(len(b.Bytes())))
	base64.StdEncoding.Encode(buff, b.Bytes())

	out := string(buff)
	return out, nil
}
