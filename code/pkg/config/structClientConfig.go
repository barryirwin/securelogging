package slogconfig

// ClientConfig :
//
// Struct for the server config
type ClientConfig struct {
	ServerIP    string   // IP of the slog server
	Port        string   // Port to send data to
	Protocol    string   // Protocol to send data
	CommsPubKey string   // Path to the public key that the clients use
	BufferLines int      // Buffer for the tail file (for bursts)
	FileList    []string // List of files to tail -f
}

// NewClientConfig :
//
// Contructor for ClientConfig, returns an empty one with defaults
func NewClientConfig() ClientConfig {
	return ClientConfig{
		ServerIP:    "127.0.0.1",
		Port:        "3333",
		Protocol:    "tcp",
		BufferLines: 10,
		CommsPubKey: "",
		FileList:    []string{},
	}
}
