package rtmp

import (
	"bytes"
	"encoding/binary"
	"log"
)

//Message header building blocks
func buildBasicHeader(format uint8, chunkStreamId uint32) []byte {
	var buf bytes.Buffer

	switch {
	case chunkStreamId >= 2 && chunkStreamId <= 63:
		b := (format << 6) | uint8(chunkStreamId)
		buf.WriteByte(b)

	case chunkStreamId >= 64 && chunkStreamId <= 319:
		buf.WriteByte(format << 6)
		buf.WriteByte(uint8(chunkStreamId - 64))

	case chunkStreamId >= 320 && chunkStreamId <= 65599:
		buf.WriteByte((format << 6) | 1)
		csId := chunkStreamId - 64
		buf.WriteByte(uint8(csId & 0xFF))
		buf.WriteByte(uint8((csId >> 8) & 0xFF))

	default:
		log.Fatal("invalid chunkStreamId")
	}

	return buf.Bytes()
}

func buildType0MessageHeader(timestamp uint32, messageLength uint32, messageTypeId uint8, messageStreamId uint32) []byte{
	var buf bytes.Buffer

	if timestamp >= 16777215 {
		buf.Write([]byte{0xFF, 0xFF, 0xFF}) //Requires extended timestamp
	}else{
		//Timestamp in 24 bit format
		buf.WriteByte(byte((timestamp >> 16) & 0xFF))
		buf.WriteByte(byte((timestamp >> 8) & 0xFF))
		buf.WriteByte(byte(timestamp & 0xFF))
	}

	//Message length in 24 bit format
	buf.WriteByte(byte((messageLength >> 16) & 0xFF))
	buf.WriteByte(byte((messageLength >> 8) & 0xFF))
	buf.WriteByte(byte(messageLength & 0xFF))

	buf.WriteByte(byte(messageTypeId))

	messageStreamIdBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(messageStreamIdBytes, uint32(messageStreamId))
	buf.Write(messageStreamIdBytes)

	return buf.Bytes()
}

// func buildType1MessageHeader(timestampDelta uint32, messageLength uint32, messageTypeId uint8) []byte{
// 	var buf bytes.Buffer

// 	//Timestamp in 24 bit format
// 	buf.WriteByte(byte((timestampDelta >> 16) & 0xFF))
// 	buf.WriteByte(byte((timestampDelta >> 8) & 0xFF))
// 	buf.WriteByte(byte(timestampDelta & 0xFF))
	
// 	//Message length in 24 bit format
// 	buf.WriteByte(byte((messageLength >> 16) & 0xFF))
// 	buf.WriteByte(byte((messageLength >> 8) & 0xFF))
// 	buf.WriteByte(byte(messageLength & 0xFF))

// 	buf.WriteByte(byte(messageTypeId))

// 	return buf.Bytes()
// }

// func buildType2MessageHeader(timestampDelta uint32) []byte{
// 	var buf bytes.Buffer

// 	//Timestamp in 24 bit format
// 	buf.WriteByte(byte((timestampDelta >> 16) & 0xFF))
// 	buf.WriteByte(byte((timestampDelta >> 8) & 0xFF))
// 	buf.WriteByte(byte(timestampDelta & 0xFF))

// 	return buf.Bytes()
// }

// func buildExtendedTimestamp(timestamp uint32) []byte{
// 	extTimestampBytes := make([]byte, 4)
// 	binary.BigEndian.PutUint32(extTimestampBytes, uint32(timestamp))

// 	return extTimestampBytes
// }

func prepend(data, prefix []byte) []byte {
    result := make([]byte, len(prefix)+len(data))
    copy(result, prefix)
    copy(result[len(prefix):], data)
    return result
}


// Extended timestamp not supported in current version
func buildMessageChunks(body []byte, chunkStreamId uint32, messageTypeId uint8, messageStreamId uint32, protocolStatus *ProtocolStatus)[][]byte {
	chunkSize := int(protocolStatus.chunkSize)
	bodyLen := len(body)

	bodyParts := (bodyLen + chunkSize - 1) / chunkSize
	chunks := make([][]byte, bodyParts)

	firstChunkHeader := make([]byte, 0)
	firstChunkHeader = append(firstChunkHeader, buildBasicHeader(0, chunkStreamId)...)
	firstChunkHeader = append(firstChunkHeader, buildType0MessageHeader(protocolStatus.baseTimestamp, uint32(bodyLen), messageTypeId, messageStreamId)...)

	subsequentChunkHeader := make([]byte, 0)
	subsequentChunkHeader = append(subsequentChunkHeader, buildBasicHeader(4, chunkStreamId)...)

	for i := 0; i < bodyParts; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > bodyLen {
			end = bodyLen
		}
		chunks[i] = body[start:end]

		if i == 0 {
			chunks[i] = prepend(chunks[i], firstChunkHeader)
		}else{
			chunks[i] = prepend(chunks[i], subsequentChunkHeader)
		}
	}
	return chunks
}