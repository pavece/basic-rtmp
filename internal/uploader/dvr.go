package uploader

import (
	"io"
	"log"
	"os"
)

func createDVRPlaylist(mediaPath string, masterReader io.Reader){
	masterPlaylistContent, err := io.ReadAll(masterReader)
	if err != nil {
		log.Fatal("Failed to read masterlist")
	}

	err = os.Mkdir(mediaPath + "/dvr", 0777)
	err = os.WriteFile(mediaPath + "/dvr/dvr.m3u8", masterPlaylistContent, 0777)
	if err != nil {
		log.Fatal("Couldn't generate DVR: Failed to write DVR file")
	}
}

func WriteDVRPlaylist(mediaPath string, masterReader io.Reader) {
	_, err := os.Stat(mediaPath + "/dvr/dvr.m3u8")
	if os.IsNotExist(err) {
		createDVRPlaylist(mediaPath, masterReader)
		return
	}
}