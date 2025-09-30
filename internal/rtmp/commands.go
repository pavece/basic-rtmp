package rtmp

import (
	"bytes"
	"fmt"
	"net"

	"github.com/pavece/simple-rtmp/internal/callbacks"
	"github.com/pavece/simple-rtmp/internal/streams"
	"github.com/yutopp/go-amf0"
)

/*
	AMF0 command handling for media streams
*/
const WINDOW_SIZE_BYTE = 8192
const MEDIA_BUFFER_SIZE_MS = 500

var commandHandlers = map[string]func(Chunk, *Rtmp,  net.Conn){
	"connect": connect,
	"createStream": createStream,
	"play": commandNotImplemented,
	"publish": publish,
	"deleteStream": deleteStream,
}


func connect(chunk Chunk, protocolStatus *Rtmp, connection net.Conn){
	reader := bytes.NewReader(chunk.Data)
	decoder := amf0.NewDecoder(reader)
	
	var connectProps map[string]string
	
	decoder.Decode(&connectProps)
	decoder.Decode(&connectProps)
	decoder.Decode(&connectProps)

	protocolStatus.streamProps = streams.CreateNewStream(connectProps["tcUrl"])

	protocolStatus.sendWindowAckSize(WINDOW_SIZE_BYTE)
	protocolStatus.sendPeerBandwidth(WINDOW_SIZE_BYTE, 0)
	protocolStatus.sendStreamBeginCommand(uint32(protocolStatus.streamProps.StreamId))
	protocolStatus.sendConnectionResultCommand(protocolStatus.streamProps.StreamId)
}

func createStream(chunk Chunk, protocolStatus *Rtmp, connection net.Conn){
	protocolStatus.sendCreateStreamResultCommand(4, 1)
}

func publish(chunk Chunk, protocolStatus *Rtmp, connection net.Conn){
	reader := bytes.NewReader(chunk.Data)
	decoder := amf0.NewDecoder(reader)
	
	streamKey := ""
	decoder.Decode(&streamKey)
	decoder.Decode(&streamKey)
	decoder.Decode(&streamKey)
	decoder.Decode(&streamKey)

	err := callbacks.ValidateStreamKey(streamKey)
	if err != nil {
		deleteStream(Chunk{}, protocolStatus, connection)
		connection.Close()
	}

	mediaId, err := callbacks.GenerateMediaId(streamKey)
	if err != nil {
		deleteStream(Chunk{}, protocolStatus, connection)
		connection.Close()
	}
	
	protocolStatus.streamProps.MediaId = mediaId	
	protocolStatus.sendStreamBeginCommand(uint32(protocolStatus.streamProps.StreamId))
	protocolStatus.sendPublishStart(uint32(protocolStatus.streamProps.StreamId))
}

func deleteStream(chunk Chunk, protocolStatus *Rtmp, connection net.Conn){
	protocolStatus.flvWriter.Close()
	protocolStatus.ffmpegPipe.Close()
	streams.RemoveStream(protocolStatus.streamProps)
	protocolStatus.Socket.Close()
}

func commandNotImplemented(chunk Chunk, protocolStatus *Rtmp, connection net.Conn){
	fmt.Println("Command not implemented")
}