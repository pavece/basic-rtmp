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

const DVR_UPLOAD_SEGMENT_WINDOW = 5
type DVRGenerator struct {
	DVRLastLine map[string]string
	VideoHeight int
	MediaId string
	Uploader *FileUploader
	SegmentsSinceLastUpload map[string]int
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


	err = os.WriteFile(mediaPath + "/dvr/" +  filename, masterPlaylistContent, 0777)
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

	file, err := os.OpenFile(mediaPath + "/dvr/" + filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal("Failed to open ", filename, " for last segment appending")
	}

	defer file.Close()

	reverseSlice(lastSegment)
	file.Write([]byte(strings.Join(lastSegment, "\n")))

	segmentsSinceUpload, ok := d.SegmentsSinceLastUpload[filename]
	if !ok {
		d.SegmentsSinceLastUpload[filename] = 1
	}else{
		d.SegmentsSinceLastUpload[filename] = segmentsSinceUpload + 1
	}

	if segmentsSinceUpload + 1 >= DVR_UPLOAD_SEGMENT_WINDOW {
		file.Seek(0, io.SeekStart)
		listData, err := io.ReadAll(file) 
		if err != nil {
			fmt.Printf("Failed to open file %s for reading (dvr list uploading): %s", filename, err)
			return
		}

	 	d.uploadList(filename, listData)
		d.SegmentsSinceLastUpload[filename] = 0
	}
}

func (d *DVRGenerator) uploadList(filename string, data []byte){
	// Append the endlist statement so you can actually play the DVR
	endlist := []byte("\n#EXT-X-ENDLIST")
	data = append(data, endlist...)	

	err := d.Uploader.UploadFile(bytes.NewReader(data), fmt.Sprintf("%s/%s", d.MediaId, filename))
	if err != nil {
		fmt.Printf("error while uploading dvr list: %s \n", err)
	}
}


func (d *DVRGenerator) generateAndUploadMasterlist(){
    baseUrl := os.Getenv("OBJECT_STORE_PUBLIC_PATH")+"/"+os.Getenv("CDN_BUCKET_NAME")+"/"+d.MediaId+"/"

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

	err := d.Uploader.UploadFile(bytes.NewReader([]byte(strings.Join(masterlistContent, "\n"))), fmt.Sprintf("%s/dvr-master.m3u8", d.MediaId))
	if err != nil {
		fmt.Printf("error while uploading dvr list: %s \n", err)
	}
}

func (d *DVRGenerator) WriteDVRPlaylist(mediaPath string, filename string, masterReader io.Reader) {
	if strings.HasPrefix(filename, "master") {
		d.generateAndUploadMasterlist()
		return
	}

	_, err := os.Stat(mediaPath + "/dvr/" + "dvr-" + filename)
	if os.IsNotExist(err) {
		d.createDVRPlaylist(mediaPath, "dvr-" + filename, masterReader)
		return
	}

	d.appendLastSegment(mediaPath, "dvr-" + filename, masterReader)
}

func NewDVRGenerator(videoHeight int, mediaId string, uploader *FileUploader) *DVRGenerator{
	return &DVRGenerator{
		DVRLastLine: make(map[string]string),
		VideoHeight: videoHeight,
		MediaId: mediaId,
		Uploader: uploader,
		SegmentsSinceLastUpload: make(map[string]int),
	}
}