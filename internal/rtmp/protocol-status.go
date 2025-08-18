package rtmp

import (
	"github.com/pavece/simple-rtmp/internal/flv"
	"github.com/pavece/simple-rtmp/internal/transcoding"
)


//TODO: All this should go into a connection independent DS to allow multiple streams (This is just for tesing)
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

var _, FfmpegPipe, _ = transcoding.SetupTranscoder() //TODO: Transcoder should be set up on createStream
var FlvWriter = flv.NewFLVWriter(FfmpegPipe, 500) 