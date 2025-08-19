package rtmp

import (
	"github.com/pavece/simple-rtmp/internal/flv"
)


type ProtocolStatus struct {
	StreamClosed 	bool
	chunkSize       uint32
	baseTimestamp   uint32
	clientWindowAck uint32
	serverWindowAck uint32
	chunkStreams map[int]Chunk
	flvWriter *flv.FLVWriter
	mediaMetadata map[string]int
}

type Chunk struct {
	BasicHeader BasicHeaderData
	Header      Type0HeaderData
	Data        []byte
}

func NewProtocolStatus() *ProtocolStatus{
	return &ProtocolStatus{chunkSize: 128, baseTimestamp: 0, clientWindowAck: 0, serverWindowAck: 0, chunkStreams: make(map[int]Chunk), flvWriter: nil, StreamClosed: false}
}
