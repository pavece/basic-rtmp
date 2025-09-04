package rtmp

import (
	"bytes"
	"fmt"
	"net"

	"github.com/pavece/simple-rtmp/internal/streams"
	"github.com/yutopp/go-amf0"
)

/*
	AMF0 command handling for media streams
*/
const WINDOW_SIZE_BYTE = 8192
const MEDIA_BUFFER_SIZE_MS = 500

var commandHandlers = map[string]func(Chunk, *ProtocolStatus,  net.Conn){
	"connect": connect,
	"createStream": createStream,
	"play": commandNotImplemented,
	"publish": publish,
	"deleteStream": deleteStream,
}


func connect(chunk Chunk, protocolStatus *ProtocolStatus, connection net.Conn){
	reader := bytes.NewReader(chunk.Data)
	decoder := amf0.NewDecoder(reader)
	
	var connectProps map[string]string
	
	decoder.Decode(&connectProps)
	decoder.Decode(&connectProps)
	decoder.Decode(&connectProps)

	protocolStatus.streamProps = streams.CreateNewStream(connectProps["tcUrl"])

	sendWindowAckSize(connection, WINDOW_SIZE_BYTE, protocolStatus)
	sendPeerBandwidth(connection, WINDOW_SIZE_BYTE, 0, protocolStatus)
	sendStreamBeginCommand(connection, uint32(protocolStatus.streamProps.StreamId), protocolStatus)
	sendConnectionResultCommand(connection, protocolStatus.streamProps.StreamId, protocolStatus)
}

func createStream(chunk Chunk, protocolStatus *ProtocolStatus, connection net.Conn){
	sendCreateStreamResultCommand(connection, 4, 1, protocolStatus)
}

func publish(chunk Chunk, protocolStatus *ProtocolStatus, connection net.Conn){
	reader := bytes.NewReader(chunk.Data)
	decoder := amf0.NewDecoder(reader)
	
	streamKey := ""
	decoder.Decode(&streamKey)
	decoder.Decode(&streamKey)
	decoder.Decode(&streamKey)
	decoder.Decode(&streamKey)

	err := streams.ValidateStreamKey(streamKey)
	if err != nil {
		deleteStream(Chunk{}, protocolStatus, connection)
		connection.Close()
	}

	mediaId, err := streams.GenerateMediaId(streamKey)
	if err != nil {
		deleteStream(Chunk{}, protocolStatus, connection)
		connection.Close()
	}
	
	protocolStatus.streamProps.MediaId = mediaId	
	sendStreamBeginCommand(connection, uint32(protocolStatus.streamProps.StreamId), protocolStatus)
	sendPublishStart(connection, uint32(protocolStatus.streamProps.StreamId), protocolStatus)
}

func deleteStream(chunk Chunk, protocolStatus *ProtocolStatus, connection net.Conn){
	protocolStatus.flvWriter.Close()
	protocolStatus.ffmpegPipe.Close()
	streams.RemoveStream(protocolStatus.streamProps)
	protocolStatus.Socket.Close()
}

func commandNotImplemented(chunk Chunk, protocolStatus *ProtocolStatus, connection net.Conn){
	fmt.Println("Command not implemented")
}