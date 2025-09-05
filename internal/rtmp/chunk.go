package rtmp

import (
	"fmt"
	"log"
	"net"
)


func ReadChunkData(connection net.Conn, protocolStatus *ProtocolStatus) error {
	basicHeaderData, err := ParseBasicHeader(connection)
	if err != nil {
		return err
	}
	

	messageHeaderData := ParseMessageHeader(connection, basicHeaderData.Fmt)
	type0header, ok := messageHeaderData.(Type0HeaderData)

	if ok {
		//First chunk of chunk stream
		protocolStatus.chunkStreams[basicHeaderData.ChunkStreamId] = Chunk{Header: type0header, BasicHeader: basicHeaderData}
	}

	currentChunkStream, ok := protocolStatus.chunkStreams[basicHeaderData.ChunkStreamId]

	if !ok {
		log.Fatal("Chunk stream not found")
	}

	if basicHeaderData.Fmt == 1 {
		type1header := messageHeaderData.(Type1HeaderData)
		currentChunkStream.Header.Timestamp += type1header.TimestampDelta
		currentChunkStream.Header.MessageTypeId = type1header.MessageTypeId
		currentChunkStream.Header.MessageLength = type1header.MessageLength
	}

	currentChunkStream.BasicHeader = basicHeaderData

	bufferSize := min(max(int(currentChunkStream.Header.MessageLength) - len(currentChunkStream.Data), 0), int(protocolStatus.chunkSize))

	chunkData := make([]byte, bufferSize)
	connection.Read(chunkData)
	
	currentChunkStream.Data = append(currentChunkStream.Data, chunkData...)
	
	protocolStatus.chunkStreams[basicHeaderData.ChunkStreamId] = currentChunkStream

	if len(currentChunkStream.Data) >= int(currentChunkStream.Header.MessageLength) {
		//Full message on board
		handler, ok := ControlHandlers[int(currentChunkStream.Header.MessageTypeId)]
		if !ok {
			fmt.Println("Handler not implemented")
			return nil
		}

		handler(currentChunkStream, protocolStatus, connection)

		currentChunkStream.Data = []byte{}
		protocolStatus.chunkStreams[basicHeaderData.ChunkStreamId] = currentChunkStream
	}

	return nil
}


