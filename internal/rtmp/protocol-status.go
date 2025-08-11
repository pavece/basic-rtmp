package rtmp

type ProtocolStatus struct {
	chunkSize     uint32
	baseTimestamp uint32
}

var protocolStatus = ProtocolStatus{chunkSize: 128, baseTimestamp: 0}