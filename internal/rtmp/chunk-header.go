package rtmp

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

type BasicHeaderData struct {
	Fmt           int
	ChunkStreamId int
}

type Type0HeaderData struct {
	Timestamp 		uint32
	MessageLength	uint32
	MessageTypeId 	uint
	MessageStreamId uint32
}

type Type1HeaderData struct {
	TimestampDelta 	uint32
	MessageLength	uint32
	MessageTypeId 	uint
}

type Type2HeaderData struct {
	TimestampDelta 	uint32
}


func ParseBasicHeader(connection net.Conn) BasicHeaderData {
	byte0 := make([]byte, 1)
	if _, err := io.ReadFull(connection, byte0); err != nil{
		log.Fatal("Failed to read first byte from chunk stream basic header")
	}
	
	headerData := BasicHeaderData{}
	headerData.Fmt = int((byte0[0] & 0xC0) >> 6)

	csIdValue := int(byte0[0] & 0x3F)

	if csIdValue == 0 {
		byte1 := make([]byte, 1)
		io.ReadFull(connection, byte1)

		csIdValue = int(byte1[0]) + 64
	}else if csIdValue == 1 {
		byte1and2 := make([]byte, 2)
		io.ReadFull(connection, byte1and2)

		csIdValue = int(byte1and2[1]) * 256 + int(byte1and2[0]) + 64
	}

	headerData.ChunkStreamId = csIdValue

	return headerData
}

func ParseMessageHeader(connection net.Conn, headerType int) interface{} {
	messageHeaderTypeLength := map[int]int{
		0: 11,
		1: 7,
		2: 3,
		3: 0,
	}

	messageHeader := make([]byte, messageHeaderTypeLength[headerType])
	if _, err := io.ReadFull(connection, messageHeader); err != nil {
		log.Fatal("Failed to read message header", err)
	}

	switch headerType {
	case 0:
		headerData := parseType0Header(messageHeader)
		headerData.Timestamp = parseExtendedTimestamp(connection, headerData.Timestamp)
		return headerData
	case 1:
		headerData := parseType1Header(messageHeader)
		headerData.TimestampDelta = parseExtendedTimestamp(connection, headerData.TimestampDelta)
		return headerData
	case 2:
		headerData := parseType2Header(messageHeader)
		headerData.TimestampDelta = parseExtendedTimestamp(connection, headerData.TimestampDelta)
		return headerData
	case 3:
		break
	default:
		break
	}

	return nil
}

func readUint24(b []byte) uint32 {
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}


func parseType0Header(header []byte) Type0HeaderData {
	return Type0HeaderData{
		Timestamp:      readUint24(header[0:3]),
		MessageLength:  readUint24(header[3:6]),
		MessageTypeId:  uint(header[6]),
		MessageStreamId: binary.LittleEndian.Uint32(header[7:11]),
	}
}

func parseType1Header(header []byte) Type1HeaderData {
	return Type1HeaderData{
		TimestampDelta: readUint24(header[0:3]),
		MessageLength:  readUint24(header[3:6]),
		MessageTypeId:  uint(header[6]),
	}
}

func parseType2Header(header []byte) Type2HeaderData {
	return Type2HeaderData{
		TimestampDelta: readUint24(header[0:3]),
	}
}

func parseExtendedTimestamp(connection net.Conn, currentTimestamp uint32) uint32{
	if currentTimestamp != 16777215 {
		return currentTimestamp
	}

	extendedTimestamp := make([]byte, 4)
	io.ReadFull(connection, extendedTimestamp)
	return binary.BigEndian.Uint32(extendedTimestamp)
}