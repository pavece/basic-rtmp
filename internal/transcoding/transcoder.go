package transcoding

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pavece/simple-rtmp/internal/config"
	"github.com/pavece/simple-rtmp/internal/uploader"
)

type Transcoder struct {
    mediaId string
    mediaMetadata map[string]int

    mediaUploader uploader.FileWatcher
}


func (t *Transcoder) createMediaFolder(){
    _, err := os.Stat("./media")
    if os.IsNotExist(err) {
        os.Mkdir("./media", 0777)
    }

    err = os.Mkdir("./media/" + t.mediaId, 0777);
    if err != nil {
        log.Fatal("Failed to create media folder")
    }
}

func (t *Transcoder) validateMediaMetadata() error {
    if t.mediaMetadata["height"] < 480 || t.mediaMetadata["width"] < 852 {
        return errors.New("video must be of at least 480p vres, ending stream")
    }
    return nil
}

func (t *Transcoder) setupRenditionFilters() ([]string, string) {
    height := t.mediaMetadata["height"]
    lastRenditionIndex := 0;
    for i, rendition := range config.Renditions {
        if rendition.Height <= height {
            lastRenditionIndex = i
        }
    }

    options := make([]string, 0)

    //General complex filter definition line
    options = append(options, "-filter_complex")
    filtersDefinition := ""
    namingStreamMap := ""

    for i := 0; i<=lastRenditionIndex; i++{
        filtersDefinition += fmt.Sprintf("[0:v]scale=%d:%d[v%d];", config.Renditions[i].Width, config.Renditions[i].Height, i)
        namingStreamMap += fmt.Sprintf("v:%d,a:%d,name:%dp ", i, i, config.Renditions[i].Height)
    }
    options = append(options, filtersDefinition)
    
    //Definition for each filter
    for i := 0; i<=lastRenditionIndex; i++{
        splitParams := strings.Split(fmt.Sprintf("-map [v%d] -map 0:a:0 -c:v:%d libx264 -b:v:%d %dk -c:a:%d aac", i, i, i, config.Renditions[i].Bitrate, i), " ")
        options = append(options, splitParams...)
    }
    
    
    return options, namingStreamMap
}

func (t *Transcoder) generateMasterList(){
    height := t.mediaMetadata["height"]
    baseUrl := os.Getenv("S3_ENDPOINT")+"/"+os.Getenv("CDN_BUCKET_NAME")+"/"+t.mediaId+"/"

    lastRenditionIndex := 0;
    for i, rendition := range config.Renditions {
        if rendition.Height <= height {
            lastRenditionIndex = i
        }
    }

    masterlistContent := []string{}

    masterlistContent = append(masterlistContent, "#EXTM3U")
    masterlistContent = append(masterlistContent, "#EXT-X-VERSION:3")

    for i := 0; i<=lastRenditionIndex; i++ {
        masterlistContent = append(masterlistContent, fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,AVERAGE-BANDWIDTH=%d,RESOLUTION=%dx%d,CODECS=\"avc1.64001f,mp4a.40.2\"", config.Renditions[i].Bitrate + 500, config.Renditions[i].Bitrate + 500, config.Renditions[i].Width, config.Renditions[i].Height))
        masterlistContent = append(masterlistContent, fmt.Sprintf("%s%dp.m3u8", baseUrl, config.Renditions[i].Height))
    }

    os.WriteFile("./media/"+t.mediaId+"/master.m3u8", []byte(strings.Join(masterlistContent, "\n")), 0777)
}

func printMasterURL(mediaId string){
    s3Endpoint := os.Getenv("S3_ENDPOINT")
    bucketName := os.Getenv("CDN_BUCKET_NAME")
    masterUrl := "Masterlist URL: " + s3Endpoint + "/" + bucketName + "/" + mediaId + "/master.m3u8"

    fmt.Println("\n" + strings.Repeat("*", len(masterUrl) + 4))
    fmt.Printf("* %s *\n", masterUrl)
    fmt.Println(strings.Repeat("*", len(masterUrl) + 4) + "\n")
}


func (t *Transcoder) SetupTranscoder(mediaMetadata map[string]int, mediaId string) (*exec.Cmd, io.WriteCloser, error) {
    t.mediaMetadata = mediaMetadata
    t.mediaId = mediaId

    err := t.validateMediaMetadata()
    if err != nil {
        log.Println(err)
        return nil, nil, err
    }

    t.createMediaFolder()
    printMasterURL(mediaId)

    ffmpegRenditionOptions, namingStreamMap := t.setupRenditionFilters()
    baseUrl := os.Getenv("S3_ENDPOINT")+"/"+os.Getenv("CDN_BUCKET_NAME")+"/"+t.mediaId+"/"

    args := []string{
        "-i", "pipe:0",
    }
    args = append(args, ffmpegRenditionOptions...)
    args = append(args, "-f", "hls",
        "-hls_time", "2",
        "-hls_list_size", "4",
        "-hls_flags", "append_list+delete_segments+independent_segments",
        "-hls_base_url", baseUrl,
        "-hls_segment_filename", "./media/"+mediaId+"/%v-segment-%d.ts",
        "-var_stream_map", namingStreamMap,
        "./media/"+mediaId+"/%v.m3u8",
    )

    ffmpegCommand := exec.Command("ffmpeg", args...)
    ffmpegCommand.Stderr = os.Stderr
    ffmpegPipe, err := ffmpegCommand.StdinPipe()
    if err != nil {
        return nil, nil, err
    }

    err = ffmpegCommand.Start()
    if err != nil {
        return nil, nil, err
    }

	t.mediaUploader.InitFileWatcher(mediaId, t.mediaMetadata["height"])
    t.generateMasterList()

    return ffmpegCommand, ffmpegPipe, nil
}