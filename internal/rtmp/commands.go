package rtmp

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/pavece/simple-rtmp/internal/flv"
	"github.com/pavece/simple-rtmp/internal/transcoding"
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
}


func connect(chunk Chunk, protocolStatus *ProtocolStatus, connection net.Conn){
	reader := bytes.NewReader(chunk.Data)
	decoder := amf0.NewDecoder(reader)

	for {
		var decoded0 interface{}
		err := decoder.Decode(&decoded0)

		if err == io.EOF {
			break;
		}

		if err != nil {
			log.Fatal(err)
		}
	}
	
	sendWindowAckSize(connection, WINDOW_SIZE_BYTE, protocolStatus)
	sendPeerBandwidth(connection, WINDOW_SIZE_BYTE, 0, protocolStatus)
	sendStreamBeginCommand(connection, 1, protocolStatus) //TODO: Implement correct streamId handling
	sendConnectionResultCommand(connection, 1, protocolStatus)
}

func createStream(chunk Chunk, protocolStatus *ProtocolStatus, connection net.Conn){
	fmt.Println("Create stream command")

	reader := bytes.NewReader(chunk.Data)
	decoder := amf0.NewDecoder(reader)

	//TODO: This should go inside onMetadata and get metadata properties directly
	_, FfmpegPipe, _ := transcoding.SetupTranscoder() 
	protocolStatus.flvWriter = flv.NewFLVWriter(FfmpegPipe, MEDIA_BUFFER_SIZE_MS) 

	for {
		var decoded0 interface{}
		err := decoder.Decode(&decoded0)

		if err == io.EOF {
			break;
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	sendCreateStreamResultCommand(connection, 4, 1, protocolStatus)
}

func publish(chunk Chunk, protocolStatus *ProtocolStatus, connection net.Conn){
	reader := bytes.NewReader(chunk.Data)
	decoder := amf0.NewDecoder(reader)

	for {
		var decoded0 interface{}
		err := decoder.Decode(&decoded0)

		if err == io.EOF {
			break;
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	sendStreamBeginCommand(connection, 1, protocolStatus)
	sendPublishStart(connection, 1, protocolStatus)
}

func commandNotImplemented(chunk Chunk, protocolStatus *ProtocolStatus, connection net.Conn){
	fmt.Println("Command not implemented")
}