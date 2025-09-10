package uploader

import (
	"io"
	"log"
	"os"
	"strings"
)

//TODO: Function for uploading these results

type DVRGenerator struct {
	DVRLastLine map[string]string
}

func reverseSlice[T any](slice []T){
	for i, j := 0, len(slice)-1; i<j; i, j = i+1, j-1{
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func createDvrDirectoryIfNotExists(mediaPath string){
	path := mediaPath + "/dvr"
	_, err := os.Stat(path)
	
	if os.IsNotExist(err) {
		os.Mkdir(path, 0777);
	}
}

//Clones the whole list (initial cloning)
func (d *DVRGenerator) createDVRPlaylist(mediaPath string, filename string, masterReader io.Reader){
	masterPlaylistContent, err := io.ReadAll(masterReader)
	if err != nil {
		log.Fatal("Error reading .m3u8 list for DVR generation")
	}

	createDvrDirectoryIfNotExists(mediaPath)
	listLines := strings.Split(string(masterPlaylistContent), "\n")
	d.DVRLastLine[filename] = listLines[len(listLines)-2] //Skip the end space


	err = os.WriteFile(mediaPath + "/dvr/" + filename, masterPlaylistContent, 0777)
	if err != nil {
		log.Fatal("Couldn't generate DVR: Failed to write DVR file")
	}
}

//Copy the latest segment(s)
func (d *DVRGenerator) appendLastSegment(mediaPath string, filename string, masterReader io.Reader){
	listContent, err := io.ReadAll(masterReader)
	if err != nil {
		log.Fatal("Error reading .m3u8 list for DVR generation")
	}

	listLines := strings.Split(string(listContent), "\n")
	lastSegment := []string{}
	i := len(listLines)-1;

	for d.DVRLastLine[filename] != listLines[i] && i>= 0 {
		lastSegment = append(lastSegment, listLines[i])
		i--
	}

	if len(lastSegment) >= 2 {
		d.DVRLastLine[filename] = lastSegment[1] //Ignore trailing whitespace
	}

	file, err := os.OpenFile(mediaPath + "/dvr/" + filename, os.O_APPEND, 0777)
	if err != nil {
		log.Fatal("Failed to open ", filename, " for last segment appending")
	}

	reverseSlice(lastSegment)
	file.Write([]byte(strings.Join(lastSegment, "\n")))
}


func (d *DVRGenerator) WriteDVRPlaylist(mediaPath string, filename string, masterReader io.Reader) {
	_, err := os.Stat(mediaPath + "/dvr/" + filename)
	if os.IsNotExist(err) {
		d.createDVRPlaylist(mediaPath, filename, masterReader)
		return
	}

	d.appendLastSegment(mediaPath, filename, masterReader)
}

func NewDVRGenerator() *DVRGenerator{
	return &DVRGenerator{
		DVRLastLine: make(map[string]string),
	}
}