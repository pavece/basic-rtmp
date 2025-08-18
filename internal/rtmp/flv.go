package rtmp

import (
	"encoding/binary"
	"fmt"
	"io"
	"sort"
)

const (
	FLV_TAG_AUDIO  = 8
	FLV_TAG_VIDEO  = 9
)

type MediaChunk struct {
	Timestamp uint32
	Type      byte
	Payload   []byte
}

type FLVWriter struct {
	w              io.Writer
	buffer         []MediaChunk
	bufferTime     uint32
	lastAudioChunk MediaChunk
	lastVideoChunk MediaChunk
}

func NewFLVWriter(w io.Writer, bufferTime uint32) *FLVWriter {
	writer := &FLVWriter{
		w: w,
		bufferTime: bufferTime,
	}
	writer.WriteHeader()
	return writer
}

func (f *FLVWriter) WriteHeader() error {
	header := []byte{'F', 'L', 'V', 0x01, 0x05, 0, 0, 0, 9}
	if _, err := f.w.Write(header); err != nil {
		return err
	}
	
	return binary.Write(f.w, binary.BigEndian, uint32(0))
}

func (f *FLVWriter) AddChunk(chunk MediaChunk) error {
	if len(chunk.Payload) == 0 {
		return nil
	}

	f.buffer = append(f.buffer, chunk)

	if chunk.Type == FLV_TAG_VIDEO {
		f.lastVideoChunk = chunk
	} else if chunk.Type == FLV_TAG_AUDIO {
		f.lastAudioChunk = chunk
	}

	safeTs := uint32(0)
	if f.lastAudioChunk.Timestamp != 0 && f.lastVideoChunk.Timestamp != 0 {
		if f.lastAudioChunk.Timestamp < f.lastVideoChunk.Timestamp {
			safeTs = f.lastAudioChunk.Timestamp
		} else {
			safeTs = f.lastVideoChunk.Timestamp
		}
	} else if f.lastAudioChunk.Timestamp != 0 {
		safeTs = f.lastAudioChunk.Timestamp
	} else if f.lastVideoChunk.Timestamp != 0 {
		safeTs = f.lastVideoChunk.Timestamp
	}

	return f.FlushUpTo(safeTs)
}

func (f *FLVWriter) FlushUpTo(maxTs uint32) error {
	if len(f.buffer) == 0 {
		return nil
	}

	sort.Slice(f.buffer, func(i, j int) bool {
		if f.buffer[i].Timestamp == f.buffer[j].Timestamp {
			return f.buffer[i].Type == FLV_TAG_VIDEO && f.buffer[j].Type == FLV_TAG_AUDIO
		}
		return f.buffer[i].Timestamp < f.buffer[j].Timestamp
	})

	newBuffer := f.buffer[:0]

	for _, chunk := range f.buffer {
		if chunk.Timestamp <= maxTs {
			if err := f.writeFLVTag(chunk.Type, chunk.Timestamp, chunk.Payload); err != nil {
				return fmt.Errorf("failed to write FLV tag: %w", err)
			}
		} else {
			newBuffer = append(newBuffer, chunk)
		}
	}
	
	f.buffer = newBuffer
	return nil
}

func (f *FLVWriter) FlushAll() error {
	if len(f.buffer) == 0 {
		return nil
	}

	sort.Slice(f.buffer, func(i, j int) bool {
		if f.buffer[i].Timestamp == f.buffer[j].Timestamp {
			return f.buffer[i].Type == FLV_TAG_VIDEO && f.buffer[j].Type == FLV_TAG_AUDIO
		}
		return f.buffer[i].Timestamp < f.buffer[j].Timestamp
	})

	for _, chunk := range f.buffer {
		if err := f.writeFLVTag(chunk.Type, chunk.Timestamp, chunk.Payload); err != nil {
			return fmt.Errorf("failed to write FLV tag: %w", err)
		}
	}

	f.buffer = f.buffer[:0]
	return nil
}

func putUint24BigEndian(b []byte, v uint32) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}

func (f *FLVWriter) writeFLVTag(tagType byte, dts uint32, payload []byte) error {
	if len(payload) == 0 {
		return nil
	}

	if tagType != FLV_TAG_AUDIO && tagType != FLV_TAG_VIDEO {
		return fmt.Errorf("invalid FLV tag type: %d", tagType)
	}

	dataSize := uint32(len(payload))

	header := make([]byte, 11)

	header[0] = tagType

	putUint24BigEndian(header[1:4], dataSize)
	putUint24BigEndian(header[4:7], dts&0xFFFFFF)
	header[7] = byte(dts >> 24)

	header[8] = 0
	header[9] = 0
	header[10] = 0

	if _, err := f.w.Write(header); err != nil {
		return err
	}

	if _, err := f.w.Write(payload); err != nil {
		return err
	}

	prevSize := uint32(11 + len(payload))
	return binary.Write(f.w, binary.BigEndian, prevSize)
}

func (f *FLVWriter) Close() error {
	return f.FlushAll()
}