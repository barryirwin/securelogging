package slogserver

import "time"

// WireData :
//
// Struct for data arriving from the wire, TCP or UDP
type WireData struct {
	Data       string    `json:"data"`        // What arrived
	SrcIP      string    `json:"source-ip"`   // Source IP from where it came from
	ReceivedAt time.Time `json:"received-at"` // When it was received
	SyslogData bool      `json:"syslog-data"` // To mark the packets that arrive from Syslog ports
	Protocol   string    `json:"protocol"`    // If the data arrived via TCP or UDP
}

// NewWireData :
//
// Contructor for WireData, assumes that the packet was received when the constructor is called
func NewWireData(data, srcIP, proto string, syslog bool) WireData {
	return WireData{
		Data:       data,
		SrcIP:      srcIP,
		ReceivedAt: time.Now(),
		SyslogData: syslog,
		Protocol:   proto,
	}
}
