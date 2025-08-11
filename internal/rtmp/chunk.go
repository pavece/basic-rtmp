package rtmp

import (
	"fmt"
	"log"
	"net"
)


type Chunk struct {
	BasicHeader BasicHeaderData
	Header Type0HeaderData
	Data []byte
}

var chunkStreams = make(map[int]Chunk)


func ReadChunkData(connection net.Conn){
	basicHeaderData := ParseBasicHeader(connection)
	messageHeaderData := ParseMessageHeader(connection, basicHeaderData.Fmt)
	type0header, ok := messageHeaderData.(Type0HeaderData)

	fmt.Println(basicHeaderData)
	fmt.Println(messageHeaderData)

	if ok {
		//First chunk of chunk stream
		chunkStreams[basicHeaderData.ChunkStreamId] = Chunk{Header: type0header, BasicHeader: basicHeaderData}
	}

	currentChunkStream, ok := chunkStreams[basicHeaderData.ChunkStreamId]
	
	if !ok {
		log.Fatal("Chunk stream not found")
	}

	bufferSize := min(max(int(currentChunkStream.Header.MessageLength) - len(currentChunkStream.Data), 0), int(protocolStatus.chunkSize))

	chunkData := make([]byte, bufferSize)
	connection.Read(chunkData)
	
	currentChunkStream.Data = append(currentChunkStream.Data, chunkData...)
	
	chunkStreams[basicHeaderData.ChunkStreamId] = currentChunkStream

	if len(currentChunkStream.Data) >= int(currentChunkStream.Header.MessageLength) {
		//Full message on board
		fmt.Println(currentChunkStream)
	}
}


