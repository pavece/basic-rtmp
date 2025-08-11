package rtmp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/yutopp/go-amf0"
)

// Protocol control message handling
var ControlHandlers = map[int]func(Chunk){
	1: setChunkSize,
	2: abortMessage,
	3: notImplemented, //Ack
	4: userControl,
	5: notImplemented, //Window Acknowledgement Size
	6: notImplemented, //Set Peer Bandwidth
	7: notImplemented, //Audio data
	8: notImplemented, //Video data
	18: notImplemented, //AMF0 encoded metadata
	20: parseAMF0Command,
}

func setChunkSize(chunk Chunk) {
	newSize := binary.BigEndian.Uint32(chunk.Data)
	protocolStatus.chunkSize = newSize
	fmt.Println("Updated chunk size to ", newSize)
}

func abortMessage(chunk Chunk) {
	streamId := binary.BigEndian.Uint32(chunk.Data)
	delete(chunkStreams, int(streamId))
	fmt.Println("Aborted message stream ", streamId)
}

func userControl(chunk Chunk){
	fmt.Println("User control message")
}

func parseAMF0Command(chunk Chunk) {
	fmt.Println("AMF0 data")

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
}

func notImplemented(chunk Chunk){
	fmt.Println("Not implemented")
}