package rtmp

import (
	"bytes"
	"encoding/binary"

	"github.com/yutopp/go-amf0"
)

func (ps *Rtmp) sendWindowAckSize(window uint32) {
	var body []byte
	body = binary.BigEndian.AppendUint32(body, window)
	chunks := ps.buildMessageChunks(body, 2, 5, 0)

	for _, chunk := range(chunks) {
		ps.Socket.Write(chunk)
	}
}

func (ps *Rtmp) sendPeerBandwidth(window uint32, limitType uint8) {
	var body []byte
	body = binary.BigEndian.AppendUint32(body, window)
	body = append(body, byte(limitType))

	chunks := ps.buildMessageChunks(body, 2, 6, 0)
	for _, chunk := range(chunks) {
		ps.Socket.Write(chunk)
	}
}

func (ps *Rtmp) sendStreamBeginCommand(streamId uint32){
	body := make([]byte, 0, 6)
	body = binary.BigEndian.AppendUint16(body, 0) 
	body = binary.BigEndian.AppendUint32(body, streamId)

	chunks := ps.buildMessageChunks(body, 2, 4, 0)
	for _, chunk := range(chunks) {
		ps.Socket.Write(chunk)
	}
}

func (ps *Rtmp) sendConnectionResultCommand(transactionId int){
	
	var buf bytes.Buffer
    encoder := amf0.NewEncoder(&buf)

	encoder.Encode("_result")
	encoder.Encode(transactionId)
	encoder.Encode(nil)

	info := map[string]interface{}{
        "level":       "status",
        "code":        "NetConnection.Connect.Success",
        "description": "Connection succeeded.",
    }
	encoder.Encode(info)

	chunks := ps.buildMessageChunks(buf.Bytes(), 3, 20, 0)
	for _, chunk := range(chunks) {
		ps.Socket.Write(chunk)
	}
}

func (ps *Rtmp) sendCreateStreamResultCommand(transactionId int, streamNumber uint32){
	var buf bytes.Buffer
    encoder := amf0.NewEncoder(&buf)

	encoder.Encode("_result")
	encoder.Encode(transactionId)
	encoder.Encode(nil)
	encoder.Encode(streamNumber)

	chunks := ps.buildMessageChunks(buf.Bytes(), 3, 20, 0)
	for _, chunk := range(chunks) {
		ps.Socket.Write(chunk)
	}
}

func (ps *Rtmp) sendPublishStart(streamId uint32) {
	var buf bytes.Buffer
	encoder := amf0.NewEncoder(&buf)

	encoder.Encode("onStatus")
	encoder.Encode(0)
	encoder.Encode(nil)

    info := map[string]any{
        "level":       "status",
        "code":        "NetStream.Publish.Start",
        "description": "Publishing stream.",
    }
	encoder.Encode(info)

    chunks := ps.buildMessageChunks(buf.Bytes(), 5, 20, streamId) 
    for _, chunk := range chunks {
        ps.Socket.Write(chunk)
    }
}
