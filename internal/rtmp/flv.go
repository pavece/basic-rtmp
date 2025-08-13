package rtmp

import (
	"encoding/binary"
	"io"
)

func WriteFLVHeader(w io.Writer) {
    w.Write([]byte{'F','L','V',0x01,0x05,0x00,0x00,0x00,0x09})
    binary.Write(w, binary.BigEndian, uint32(0))
}

func writeFLVTag(w io.Writer, tagType byte, timestamp uint32, payload []byte) error {
	header := make([]byte, 11)
	header[0] = tagType
	dataSize := uint32(len(payload))
	header[1] = byte(dataSize >> 16)
	header[2] = byte(dataSize >> 8)
	header[3] = byte(dataSize)
	header[4] = byte(timestamp >> 16)
	header[5] = byte(timestamp >> 8)
	header[6] = byte(timestamp)
	header[7] = byte(timestamp >> 24)
	w.Write(header)
	w.Write(payload)

	prevSize := uint32(len(payload) + 11)
	binary.Write(w, binary.BigEndian, prevSize)
	return nil
}
