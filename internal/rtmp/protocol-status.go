package rtmp

type ProtocolStatus struct {
	chunkSize       uint32
	baseTimestamp   uint32
	clientWindowAck uint32
	serverWindowAck uint32
}

var protocolStatus = ProtocolStatus{chunkSize: 128, baseTimestamp: 0, clientWindowAck: 0, serverWindowAck: 0}

type Chunk struct {
	BasicHeader BasicHeaderData
	Header      Type0HeaderData
	Data        []byte
}

var chunkStreams = make(map[int]Chunk)
