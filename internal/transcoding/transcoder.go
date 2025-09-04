package transcoding

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pavece/simple-rtmp/internal/uploader"
)

type Rendition struct {
    bitrate int
    height int
    width int
}

var renditions []Rendition = []Rendition{{bitrate: 1000, height: 480, width: 852}, {bitrate: 5000, height: 720, width: 1280}/*, {bitrate: 8000, height: 1080, width: 1920} */}

func createMediaFolder(mediaId string){
    _, err := os.Stat("./media")
    if os.IsNotExist(err) {
        os.Mkdir("./media", 0777)
    }

    err = os.Mkdir("./media/" + mediaId, 0777);
    if err != nil {
        log.Fatal("Failed to create media folder")
    }
}

func validateMediaMetadata(mediaMetadata map[string]int) error {
    if mediaMetadata["height"] < 480 || mediaMetadata["width"] < 852 {
        return errors.New("video must be of at least 480p vres, ending stream")
    }
    return nil
}

func setupRenditionFilters(height int) ([]string, string) {
    lastRenditionIndex := 0;
    for i, rendition := range renditions {
        if rendition.height >= height {
            lastRenditionIndex = i
        }
    }

    options := make([]string, 0)

    //General complex filter definition line
    options = append(options, "-filter_complex")
    filtersDefinition := ""
    namingStreamMap := ""

    for i := 0; i<lastRenditionIndex; i++{
        filtersDefinition += fmt.Sprintf("[0:v]scale=%d:%d[v%d];", renditions[i].width, renditions[i].height, i)
        namingStreamMap += fmt.Sprintf("v:%d,a:%d,name:%dp", i, i, renditions[i].height)
    }
    
    fmt.Println(namingStreamMap)
    options = append(options, filtersDefinition)

    //Definition for each filter
    for i := 0; i<lastRenditionIndex; i++{
        splitParams := strings.Split(fmt.Sprintf("-map [v%d] -map 0:a:0 -c:v:%d libx264 -b:v:%d %dk -c:a:%d aac", i, i, i, renditions[i].bitrate, i), " ")
        options = append(options, splitParams...)
    }

    
    return options, namingStreamMap
}


func SetupTranscoder(mediaMetadata map[string]int, mediaId string) (*exec.Cmd, io.WriteCloser, error) {
    err := validateMediaMetadata(mediaMetadata)
    if err != nil {
        log.Println(err)
        return nil, nil, err
    }

    createMediaFolder(mediaId)

    ffmpegRenditionOptions, namingStreamMap := setupRenditionFilters(mediaMetadata["height"])

    args := []string{
        "-i", "pipe:0",
    }
    args = append(args, ffmpegRenditionOptions...)
    args = append(args, "-f", "hls",
        "-hls_time", "2",
        "-hls_list_size", "4",
        "-hls_flags", "append_list+delete_segments",
        "-hls_base_url", os.Getenv("S3_ENDPOINT")+"/"+os.Getenv("CDN_BUCKET_NAME")+"/"+mediaId+"/",
        "-hls_segment_filename", "./media/"+mediaId+"/%v-segment-%d.ts",
        "-var_stream_map", namingStreamMap,
        "-master_pl_name", "master.m3u8",
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

	uploader.SetupFileWatcher(mediaId)
    return ffmpegCommand, ffmpegPipe, nil
}
