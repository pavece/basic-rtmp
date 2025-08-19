package rtmp

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/yutopp/go-amf0"
)

func sendWindowAckSize(connection net.Conn, window uint32, protocolStatus *ProtocolStatus) {
	var body []byte
	body = binary.BigEndian.AppendUint32(body, window)
	chunks := buildMessageChunks(body, 2, 5, 0, protocolStatus)

	for _, chunk := range(chunks) {
		connection.Write(chunk)
	}
}

func sendPeerBandwidth(connection net.Conn, window uint32, limitType uint8, protocolStatus *ProtocolStatus) {
	var body []byte
	body = binary.BigEndian.AppendUint32(body, window)
	body = append(body, byte(limitType))

	chunks := buildMessageChunks(body, 2, 6, 0, protocolStatus)
	for _, chunk := range(chunks) {
		connection.Write(chunk)
	}
}

func sendStreamBeginCommand(connection net.Conn, streamId uint32, protocolStatus *ProtocolStatus){
	body := make([]byte, 0, 6)
	body = binary.BigEndian.AppendUint16(body, 0) 
	body = binary.BigEndian.AppendUint32(body, streamId)

	chunks := buildMessageChunks(body, 2, 4, 0, protocolStatus)
	for _, chunk := range(chunks) {
		connection.Write(chunk)
	}
}

func sendConnectionResultCommand(connection net.Conn, transactionId int, protocolStatus *ProtocolStatus){
	
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

	chunks := buildMessageChunks(buf.Bytes(), 3, 20, 0, protocolStatus)
	for _, chunk := range(chunks) {
		connection.Write(chunk)
	}
}

func sendCreateStreamResultCommand(connection net.Conn, transactionId int, streamNumber uint32, protocolStatus *ProtocolStatus){
	var buf bytes.Buffer
    encoder := amf0.NewEncoder(&buf)

	encoder.Encode("_result")
	encoder.Encode(transactionId)
	encoder.Encode(nil)
	encoder.Encode(streamNumber)

	chunks := buildMessageChunks(buf.Bytes(), 3, 20, 0, protocolStatus)
	for _, chunk := range(chunks) {
		connection.Write(chunk)
	}
}

func sendPublishStart(connection net.Conn, streamId uint32, protocolStatus *ProtocolStatus) {
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

    chunks := buildMessageChunks(buf.Bytes(), 5, 20, streamId, protocolStatus) 
    for _, chunk := range chunks {
        connection.Write(chunk)
    }
}
