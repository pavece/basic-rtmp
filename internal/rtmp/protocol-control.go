package rtmp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/yutopp/go-amf0"
)

// Protocol control message handling
var ControlHandlers = map[int]func(Chunk, net.Conn){
	1: setChunkSize,
	2: abortMessage,
	3: ack, 
	4: userControl,
	5: windowAckSize,
	6: notImplemented, //Set Peer Bandwidth
	8: getAudio, 
	9: getVideo, 
	18: notImplemented, //AMF0 encoded metadata
	20: parseAMF0Command,
}

func setChunkSize(chunk Chunk, connection net.Conn) {
	newSize := binary.BigEndian.Uint32(chunk.Data)
	protocolStatus.chunkSize = newSize
	fmt.Println("Updated chunk size to ", newSize)
}

func abortMessage(chunk Chunk, connection net.Conn) {
	streamId := binary.BigEndian.Uint32(chunk.Data)
	delete(chunkStreams, int(streamId)) //TODO: incorrect
	fmt.Println("Aborted message stream ", streamId)
}

func ack(chunk Chunk, connection net.Conn){
	totalBytes := binary.BigEndian.Uint32(chunk.Data)
	fmt.Println("Recieved ACK from client, total bytes: ", totalBytes)
}

func userControl(chunk Chunk, connection net.Conn){
	fmt.Println("User control message")
}

func windowAckSize(chunk Chunk, connection net.Conn){
	window := binary.BigEndian.Uint32(chunk.Data)
	protocolStatus.clientWindowAck = window
	fmt.Println("Updated client's ack window to ", window)
}

func getAudio(chunk Chunk, connection net.Conn){
	fmt.Println("Audio chunk")
}

func getVideo(chunk Chunk, connection net.Conn){
	fmt.Println("Video chunk")
}

func parseAMF0Command(chunk Chunk, connection net.Conn) {
	reader := bytes.NewReader(chunk.Data)
	decoder := amf0.NewDecoder(reader)

	command := ""
	decoder.Decode(&command)

	handler, ok := commandHandlers[command]

	if !ok {
		fmt.Println("Command handler not implemented for command ", command)
		return
	}

	fmt.Println("Incoming command: ", command)
	handler(chunk, connection)
}

func notImplemented(chunk Chunk, connection net.Conn){
	fmt.Println("Not implemented")
}