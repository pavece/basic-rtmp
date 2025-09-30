package rtmp

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/pavece/simple-rtmp/internal/flv"
	"github.com/pavece/simple-rtmp/internal/instrumentation"
	"github.com/yutopp/go-amf0"
)

/*
	Generic protocol message handlers
*/
var ControlHandlers = map[int]func(Chunk, *Rtmp){
	1: setChunkSize,
	2: abortMessage,
	3: ack, 
	4: userControl,
	5: windowAckSize,
	6: notImplemented, //Set Peer Bandwidth
	8: getAudio, 
	9: getVideo, 
	18: getMetadata,
	20: parseAMF0Command,
}

func setChunkSize(chunk Chunk, protocolStatus *Rtmp) {
	newSize := binary.BigEndian.Uint32(chunk.Data)
	protocolStatus.chunkSize = newSize
	fmt.Println("Updated chunk size to ", newSize)
}

func abortMessage(chunk Chunk, protocolStatus *Rtmp) {
	streamId := binary.BigEndian.Uint32(chunk.Data)
	delete(protocolStatus.chunkStreams, int(streamId)) //Incorrect
	fmt.Println("Aborted message stream ", streamId)
}

func ack(chunk Chunk, protocolStatus *Rtmp){
	totalBytes := binary.BigEndian.Uint32(chunk.Data)
	fmt.Println("Recieved ACK from client, total bytes: ", totalBytes)
}

func userControl(chunk Chunk, protocolStatus *Rtmp){
	fmt.Println("User control message")
}

func windowAckSize(chunk Chunk, protocolStatus *Rtmp){
	window := binary.BigEndian.Uint32(chunk.Data)
	protocolStatus.clientWindowAck = window
	fmt.Println("Updated client's ack window to ", window)
}

func getAudio(chunk Chunk, protocolStatus *Rtmp){
	protocolStatus.flvWriter.AddChunk(flv.MediaChunk{Type: 8, Timestamp: chunk.Header.Timestamp - protocolStatus.baseTimestamp, Payload: chunk.Data})    
	instrumentation.AudioIngress.Add(float64(len(chunk.Data)))
}

func getVideo(chunk Chunk, protocolStatus *Rtmp){
	protocolStatus.flvWriter.AddChunk(flv.MediaChunk{Type: 9, Timestamp: chunk.Header.Timestamp - protocolStatus.baseTimestamp, Payload: chunk.Data})
	instrumentation.VideoIngress.Add(float64(len(chunk.Data)))
	
}

func getMetadata(chunk Chunk, protocolStatus *Rtmp) {
	reader := bytes.NewReader(chunk.Data)
	decoder := amf0.NewDecoder(reader)

	command := ""
	decoder.Decode(&command)
	decoder.Decode(&command)

	var metadata map[string]int
	decoder.Decode(&metadata)
	protocolStatus.mediaMetadata = metadata

	_, ffmpegBuffer, ffmpegPipe, err := protocolStatus.transcoder.SetupTranscoder(protocolStatus.mediaMetadata, protocolStatus.streamProps.MediaId) 
	if err != nil {
		protocolStatus.Socket.Close()
		return
	}
	
	protocolStatus.flvWriter = flv.NewFLVWriter(ffmpegBuffer, MEDIA_BUFFER_SIZE_MS) 
	protocolStatus.ffmpegPipe = ffmpegPipe
}

func parseAMF0Command(chunk Chunk, protocolStatus *Rtmp) {
	reader := bytes.NewReader(chunk.Data)
	decoder := amf0.NewDecoder(reader)

	command := ""
	decoder.Decode(&command)

	handler, ok := commandHandlers[command]

	if !ok {
		fmt.Println("Command handler not implemented for command ", command)
		return
	}

	fmt.Println("Incoming command: ", command)
	handler(chunk, protocolStatus, protocolStatus.Socket)
}

func notImplemented(chunk Chunk, protocolStatus *Rtmp){
	fmt.Println("Not implemented")
}