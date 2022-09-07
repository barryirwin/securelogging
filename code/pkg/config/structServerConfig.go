package slogconfig

// ServerConfig :
//
// Struct for the server config
type ServerConfig struct {
	TCPport                 string // Port to listen for TCP data
	UDPport                 string // Port to listen for UDP data
	TCPsyslog               string // Port to listen for TCP syslog data
	UDPsyslog               string // Port to listen for UDP syslog data
	CommsPrivKey            string // Path to the private key that the clients use
	StoragePrivKey          string // Path to the private key to decrypt data already stored by the server
	StoragePubKey           string // Path to the public key to store data
	StorageFolder           string // Where to store the files
	InfluxDBname            string // DB name for InfluxDB
	InfluxIP                string // InfluxDB server
	StoreToInflux           bool   // Wheter store or not to Influx
	PacketBufferSize        int    // Size of the buffer for incoming network packets
	FileBatchSize           int    // Batch size of logs to be collected before stored
	WireChannelBuffer       int    // Size of the buffer for incoming wire data
	ProcessingChannelBuffer int    // Size of the buffer for all processing channels
	CompressedBlob          bool   // Wheter use or not compression at the resulting blob for storage to file
	Password                string // Password used a runtime, to pass less parameters to functions
	ReplayMode              bool   // Wheter the server is running on replay mode or not
}

// NewServerConfig :
//
// Contructor for ServerConfig, returns an empty one with defaults
func NewServerConfig() ServerConfig {
	return ServerConfig{
		TCPport:                 "",
		UDPport:                 "",
		TCPsyslog:               "514",
		UDPsyslog:               "514",
		CommsPrivKey:            "",
		StoragePrivKey:          "",
		StoragePubKey:           "",
		StorageFolder:           "",
		InfluxDBname:            "",
		InfluxIP:                "",
		StoreToInflux:           false,
		PacketBufferSize:        2048,
		FileBatchSize:           1000,
		WireChannelBuffer:       100,
		ProcessingChannelBuffer: 100,
		CompressedBlob:          false,
		Password:                "",
		ReplayMode:              false,
	}
}
