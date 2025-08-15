package rtmp

import (
	"encoding/binary"
	"io"
	"sort"
)
type MediaChunk struct {
	Timestamp uint32
	Type      byte
	Payload   []byte
}

type FLVWriter struct {
	w            io.Writer
	buffer       []MediaChunk
	bufferTime   uint32
	lastAudioChunk MediaChunk
	lastVideoChunk MediaChunk
}

func NewFLVWriter(w io.Writer, bufferTime uint32) *FLVWriter {
	writer := &FLVWriter{
		w:          w,
		bufferTime: bufferTime,
	}
	writer.WriteHeader()
	return writer
}

func (f *FLVWriter) WriteHeader() {
	f.w.Write([]byte{'F', 'L', 'V', 0x01, 0x05, 0, 0, 0, 9})
	binary.Write(f.w, binary.BigEndian, uint32(0))
}

func (f *FLVWriter) AddChunk(chunk MediaChunk) error {
    f.buffer = append(f.buffer, chunk)

    if chunk.Type == 9 {
        f.lastVideoChunk = chunk
    } else {
        f.lastAudioChunk = chunk
    }

    safeTs := uint32(0)
    if f.lastAudioChunk.Timestamp != 0 {
        safeTs = f.lastAudioChunk.Timestamp
    }

    return f.FlushUpTo(safeTs)
}

func (f *FLVWriter) FlushUpTo(maxTs uint32) error {
    if len(f.buffer) == 0 {
        return nil
    }

    sort.Slice(f.buffer, func(i, j int) bool {
        return f.buffer[i].Timestamp < f.buffer[j].Timestamp
    })

    newBuffer := f.buffer[:0]
	
	for _, c := range f.buffer {

        if c.Timestamp <= maxTs {
            writeFLVTag(f.w, c.Type, c.Timestamp, c.Payload)
        } else {
            newBuffer = append(newBuffer, c)
        }
    }
    f.buffer = newBuffer
    return nil
}

func writeFLVTag(w io.Writer, tagType byte, dts uint32, payload []byte) error {
	header := make([]byte, 11)
	header[0] = tagType

	dataSize := uint32(len(payload))
	header[1] = byte(dataSize >> 16)
	header[2] = byte(dataSize >> 8)
	header[3] = byte(dataSize)

	var pts uint32 = dts

	if tagType == 9 && len(payload) >= 5 {
		compositionTime := int32(payload[2])<<16 |
			int32(payload[3])<<8 |
			int32(payload[4])
		
			if compositionTime&0x800000 != 0 {
			compositionTime |= ^0xffffff
		}
		pts += uint32(compositionTime)
	}

	header[4] = byte(pts >> 16)
	header[5] = byte(pts >> 8)
	header[6] = byte(pts)
	header[7] = byte(pts >> 24)

	header[8], header[9], header[10] = 0, 0, 0

	if _, err := w.Write(header); err != nil {
		return err
	}
	if _, err := w.Write(payload); err != nil {
		return err
	}

	prevSize := uint32(len(payload) + 11)
	return binary.Write(w, binary.BigEndian, prevSize)
}