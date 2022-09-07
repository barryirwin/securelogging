package slogserver

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/FrancoLoyola/noroff-fdp/code/pkg/checks"
	client "github.com/influxdata/influxdb1-client/v2"
)

// ServerData :
//
// Struct for data that arrived from a client, with the enrichment of the server processing
type ServerData struct {
	WireData       WireData `json:"wire-data"`       // What arrived
	DecryptedData  string   `json:"decrypted-data"`  // Decrypted data from WireData (AES + RSA)
	SuspiciousIn   bool     `json:"suspicious-in"`   // Flag suspicious incoming packets
	SuspiciousRe   bool     `json:"suspicious-re"`   // Flag suspicious re-processed packets
	ServerHostname string   `json:"server-hostname"` // More context
	PreviousHash   string   `json:"previous-hash"`   // Hash of the previous log stored
}

// NewServerData :
//
// Contructor for ServerData, data needs to be decrypted
func NewServerData(wireData WireData) ServerData {
	hostname, _ := os.Hostname()
	return ServerData{
		WireData:       wireData,
		DecryptedData:  "",
		SuspiciousIn:   false,
		SuspiciousRe:   false,
		ServerHostname: hostname,
		PreviousHash:   "",
	}
}

// ToBase64 :
//
// Returns the Base64 representation of the struct and the checksum of it.
// This data is not encrypted, needs to be encrypted if it is going to be stored.
func (sd ServerData) ToBase64() (string, uint32, error) {
	// Encoding
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	err := json.NewEncoder(encoder).Encode(sd)
	if err != nil {
		return "", 0, err
	}
	encoder.Close()
	out := buf.String()

	// Checksum
	chksum := checks.CalculateCRC([]byte(out))

	return out, chksum, nil
}

// ServerDataFromBase64 :
//
// Attempts to decode a base64 string representation of a ServerData struct
func ServerDataFromBase64(s string) (ServerData, error) {
	// Preparations, could be a one liner, but this reads better
	strReader := strings.NewReader(s)
	b64Reader := base64.NewDecoder(base64.StdEncoding, strReader)
	decoder := json.NewDecoder(b64Reader)

	var out ServerData
	err := decoder.Decode(&out)
	if err != nil {
		return out, err
	}
	return out, nil
}

// ToInfluxPoint :
//
// Returns the Influx point for this struct, to be used in a batch.
//
// If the data is replayed, set replay to true, so plots can be compared 1-on-1
//
// Only some values are of interest for plots, checksums are not part
// of the point for example.
func (sd ServerData) ToInfluxPoint(replay bool) (*client.Point, error) {

	// These are values that can be used for GROUP BY queries
	tags := map[string]string{
		"serverHostname": sd.ServerHostname,
		"from":           sd.WireData.SrcIP,
		"replayData":     strconv.FormatBool(replay),
		"syslogData":     strconv.FormatBool(sd.WireData.SyslogData),
		"protocol":       sd.WireData.Protocol,
	}

	// Values to plot / show
	// logline is the actual logline in the client
	// if it is a vlid packet, will be the last field
	l := strings.Split(sd.DecryptedData, ",")
	logline := l[len(l)-1]
	fields := map[string]interface{}{
		"suspiciousIncoming":   sd.SuspiciousIn,
		"suspiciousReplay":     sd.SuspiciousRe,
		"encryptedPayloadSize": len(sd.WireData.Data),
		"decryptedPayloadSize": len(sd.DecryptedData),
		"slogline":             sd.DecryptedData,
		"logline":              logline,
	}

	point, err := client.NewPoint("slog-log", tags, fields, sd.WireData.ReceivedAt)
	if err != nil {
		return nil, err
	}
	return point, nil
}
