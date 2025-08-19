package rtmp

import (
	"fmt"
	"net"

	"github.com/pavece/simple-rtmp/internal/flv"
	"github.com/pavece/simple-rtmp/internal/transcoding"
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
	sendWindowAckSize(connection, WINDOW_SIZE_BYTE, protocolStatus)
	sendPeerBandwidth(connection, WINDOW_SIZE_BYTE, 0, protocolStatus)
	sendStreamBeginCommand(connection, 1, protocolStatus) //TODO: Implement correct streamId handling
	sendConnectionResultCommand(connection, 1, protocolStatus)
}

func createStream(chunk Chunk, protocolStatus *ProtocolStatus, connection net.Conn){
	_, FfmpegPipe, _ := transcoding.SetupTranscoder() 
	protocolStatus.flvWriter = flv.NewFLVWriter(FfmpegPipe, MEDIA_BUFFER_SIZE_MS) 

	sendCreateStreamResultCommand(connection, 4, 1, protocolStatus)
}

func publish(chunk Chunk, protocolStatus *ProtocolStatus, connection net.Conn){
	sendStreamBeginCommand(connection, 1, protocolStatus)
	sendPublishStart(connection, 1, protocolStatus)
}

func deleteStream(chunk Chunk, protocolStatus *ProtocolStatus, connection net.Conn){
	//TODO: Perform any cleanup related to global stream data
	protocolStatus.StreamClosed = true
}

func commandNotImplemented(chunk Chunk, protocolStatus *ProtocolStatus, connection net.Conn){
	fmt.Println("Command not implemented")
}