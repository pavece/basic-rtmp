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

		fmt.Println(decoded0)
	}
	
	sendWindowAckSize(connection, 10000000, protocolStatus)
	sendPeerBandwidth(connection, 10000000, 0, protocolStatus)
	sendStreamBeginCommand(connection, 1, protocolStatus)
	sendConnectionResultCommand(connection, 1, protocolStatus)
}

func createStream(chunk Chunk, protocolStatus *ProtocolStatus, connection net.Conn){
	fmt.Println("Create stream command")

	reader := bytes.NewReader(chunk.Data)
	decoder := amf0.NewDecoder(reader)

	//TODO: This should go inside onMetadata and get metadata properties directly
	_, FfmpegPipe, _ := transcoding.SetupTranscoder() 
	protocolStatus.flvWriter = flv.NewFLVWriter(FfmpegPipe, 500) 

	for {
		var decoded0 interface{}
		err := decoder.Decode(&decoded0)

		if err == io.EOF {
			break;
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(decoded0)
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

		fmt.Println(decoded0)
	}

	sendStreamBeginCommand(connection, 1, protocolStatus)
	sendPublishStart(connection, 1, protocolStatus)
}

func commandNotImplemented(chunk Chunk, protocolStatus *ProtocolStatus, connection net.Conn){
	fmt.Println("Command not implemented")
}