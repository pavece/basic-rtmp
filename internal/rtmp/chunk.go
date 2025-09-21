package rtmp

import (
	"fmt"
	"net"
)


func ReadChunkData(connection net.Conn, protocolStatus *ProtocolStatus) error {
	basicHeaderData, err := ParseBasicHeader(connection)
	if err != nil {
		return err
	}
	
	messageHeaderData := ParseMessageHeader(connection, basicHeaderData.Fmt)

	if basicHeaderData.Fmt == 0 {
        type0header := messageHeaderData.(Type0HeaderData)
        protocolStatus.chunkStreams[basicHeaderData.ChunkStreamId] = Chunk{
            Header: type0header, 
            BasicHeader: basicHeaderData,
        }
    }
    
	currentChunkStream, ok := protocolStatus.chunkStreams[basicHeaderData.ChunkStreamId]
    if !ok {
        return fmt.Errorf("chunk stream %d not found - missing Type 0 header", basicHeaderData.ChunkStreamId)
    }

   switch basicHeaderData.Fmt {
    case 1:
        type1header := messageHeaderData.(Type1HeaderData)
        currentChunkStream.Header.Timestamp += type1header.TimestampDelta
        currentChunkStream.Header.MessageTypeId = type1header.MessageTypeId
        currentChunkStream.Header.MessageLength = type1header.MessageLength
    case 2:
        type2header := messageHeaderData.(Type2HeaderData)
        currentChunkStream.Header.Timestamp += type2header.TimestampDelta
    case 3:
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


