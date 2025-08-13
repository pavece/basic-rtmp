package rtmp

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/yutopp/go-amf0"
)

var commandHandlers = map[string]func(Chunk, net.Conn){
	"connect": connect,
	"createStream": createStream,
	"play": commandNotImplemented,
	"publish": publish,
}


func connect(chunk Chunk, connection net.Conn){
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
	
	sendWindowAckSize(connection, 10000000)
	sendPeerBandwidth(connection, 10000000, 0)
	sendStreamBeginCommand(connection, 1)
	sendConnectionResultCommand(connection, 1)
}

func createStream(chunk Chunk, connection net.Conn){
	fmt.Println("Create stream command")

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

	sendCreateStreamResultCommand(connection, 4, 1)
}

func publish(chunk Chunk, connection net.Conn){
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

	sendStreamBeginCommand(connection, 1)
	sendPublishStart(connection, 1)
}

func commandNotImplemented(chunk Chunk, connection net.Conn){
	fmt.Println("Command not implemented")
}