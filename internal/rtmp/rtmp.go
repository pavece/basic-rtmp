package rtmp

import (
	"io"
	"net"

	"github.com/pavece/simple-rtmp/internal/flv"
	"github.com/pavece/simple-rtmp/internal/streams"
	"github.com/pavece/simple-rtmp/internal/transcoding"
)


type Rtmp struct {
	chunkSize       uint32
	baseTimestamp   uint32
	clientWindowAck uint32
	serverWindowAck uint32
	chunkStreams map[int]Chunk
	flvWriter *flv.FLVWriter
	ffmpegPipe io.WriteCloser
	mediaMetadata map[string]int
	streamProps streams.StreamProps
	Socket net.Conn
	transcoder transcoding.Transcoder
}


type Chunk struct {
	BasicHeader BasicHeaderData
	Header      Type0HeaderData
	Data        []byte
}

func New(connection net.Conn) *Rtmp{
	return &Rtmp{chunkSize: 128, baseTimestamp: 0, clientWindowAck: 0, serverWindowAck: 0, chunkStreams: make(map[int]Chunk), flvWriter: nil, Socket: connection}
}