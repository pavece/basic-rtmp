package rtmp

import (
	"fmt"
	"io"
	"log"
	"net"
)

type BasicHeaderData struct {
	fmt           int
	chunkStreamId int
}

func ParseBasicHeader(connection net.Conn) BasicHeaderData {
	byte0 := make([]byte, 1)
	if _, err := io.ReadFull(connection, byte0); err != nil{
		log.Fatal("Failed to read first byte from chunk stream basic header")
	}
	
	headerData := BasicHeaderData{}
	headerData.fmt = int((byte0[0] & 0xC0) >> 6)

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

	headerData.chunkStreamId = csIdValue

	fmt.Println(headerData)
	return headerData
}

func ParseMessageHeader(messageHeader []byte) {}