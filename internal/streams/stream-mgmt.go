package streams

import (
	"fmt"
	"sort"
	"sync"
)

type StreamProps struct {
	StreamPath string
	StreamId int
	MediaId string
}

var streams []StreamProps
var streamsLock sync.Mutex


func CreateNewStream(streamPath string) StreamProps {
	streamsLock.Lock()
	lastStreamId := 0

	sort.Slice(streams, func(i, j int) bool {
		return streams[i].StreamId > streams[j].StreamId
	})	

	if len(streams) > 0 {
		lastStreamId = streams[0].StreamId
	}

	newStream := StreamProps{StreamId: lastStreamId + 1, StreamPath: streamPath}
	streams = append(streams, newStream)
	
	streamsLock.Unlock()
	return newStream
}

func RemoveStream(stream StreamProps){
	streamsLock.Lock()
	for i, s := range(streams) {
		if s.StreamId == stream.StreamId {
			streams = append(streams[:i], streams[i+1:]... )
		}
	}
	streamsLock.Unlock()
	onStramEnd(stream)
}


/* 
	Use this to mark streams as completed or run custom logic related to stream ends
*/
func onStramEnd(stream StreamProps){
	fmt.Printf("Stream with id %d and media id %s has ended \n", stream.StreamId, stream.MediaId)
}