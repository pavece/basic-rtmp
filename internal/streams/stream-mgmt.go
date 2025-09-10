package streams

import (
	"os"
	"sort"
	"sync"
	"time"

	"github.com/pavece/simple-rtmp/internal/callbacks"
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

	callbacks.OnStramEnd(stream.StreamId, stream.MediaId)
	postStreamCleanup(stream)
}

func postStreamCleanup(stream StreamProps){
	time.Sleep(4 * time.Second)

	_, err := os.Stat("./media/" + stream.MediaId)
	if err == nil {
		os.RemoveAll("./media/" + stream.MediaId)
	}
}