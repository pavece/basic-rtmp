package uploader

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pavece/simple-rtmp/internal/config"
)

//TODO: Function for uploading these results

type DVRGenerator struct {
	DVRLastLine map[string]string
	VideoHeight int
	MediaId string
	Uploader *FileUploader
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


	err = os.WriteFile(mediaPath + "/dvr/" + "dvr-" + filename, masterPlaylistContent, 0777)
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

func (d *DVRGenerator) generateAndUploadMasterlist(){
    baseUrl := os.Getenv("S3_ENDPOINT")+"/"+os.Getenv("CDN_BUCKET_NAME")+"/"+d.MediaId+"/"

    lastRenditionIndex := 0;
    for i, rendition := range config.Renditions {
        if rendition.Height <= d.VideoHeight {
            lastRenditionIndex = i
        }
    }

    masterlistContent := []string{}

    masterlistContent = append(masterlistContent, "#EXTM3U")
    masterlistContent = append(masterlistContent, "#EXT-X-VERSION:3")

    for i := 0; i<=lastRenditionIndex; i++ {
        masterlistContent = append(masterlistContent, fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,AVERAGE-BANDWIDTH=%d,RESOLUTION=%dx%d,CODECS=\"avc1.64001f,mp4a.40.2\"", config.Renditions[i].Bitrate + 500, config.Renditions[i].Bitrate + 500, config.Renditions[i].Width, config.Renditions[i].Height))
        masterlistContent = append(masterlistContent, fmt.Sprintf("%sdvr-%dp.m3u8", baseUrl, config.Renditions[i].Height))
    }

	d.Uploader.UploadFile(bytes.NewReader([]byte(strings.Join(masterlistContent, "\n"))), fmt.Sprintf("%s/dvr-master.m3u8", d.MediaId))
}

func (d *DVRGenerator) WriteDVRPlaylist(mediaPath string, filename string, masterReader io.Reader) {
	if strings.HasPrefix(filename, "master") {
		d.generateAndUploadMasterlist()
		return
	}

	_, err := os.Stat(mediaPath + "/dvr/" + "dvr-" + filename)
	if os.IsNotExist(err) {
		d.createDVRPlaylist(mediaPath,"dvr-" + filename, masterReader)
		return
	}

	d.appendLastSegment(mediaPath, filename, masterReader)
}

func NewDVRGenerator(videoHeight int, mediaId string, uploader *FileUploader) *DVRGenerator{
	return &DVRGenerator{
		DVRLastLine: make(map[string]string),
		VideoHeight: videoHeight,
		MediaId: mediaId,
		Uploader: uploader,
	}
}