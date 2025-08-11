package rtmp

import (
	"fmt"
	"log"
	"net"
)


func ReadChunkData(connection net.Conn){
	basicHeaderData := ParseBasicHeader(connection)
	messageHeaderData := ParseMessageHeader(connection, basicHeaderData.Fmt)
	type0header, ok := messageHeaderData.(Type0HeaderData)

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

		handler, ok := ControlHandlers[int(currentChunkStream.Header.MessageTypeId)]
		if !ok {
			log.Fatal("Handler not implemented")
		}

		handler(currentChunkStream)

		delete(chunkStreams, basicHeaderData.ChunkStreamId)
	}
}


