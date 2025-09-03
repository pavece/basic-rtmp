package transcoding

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/pavece/simple-rtmp/internal/uploader"
)

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

func validateMediaMetadata(mediaMetadata map[string]int) {
    //TODO: Validate metadata
    fmt.Println(mediaMetadata)
}

func setupRenditionFilters(mediaMetadata map[string]int) []string {
    return []string{}
}


func SetupTranscoder(mediaMetadata map[string]int, mediaId string) (*exec.Cmd, io.WriteCloser, error) {
    createMediaFolder(mediaId)
    validateMediaMetadata(mediaMetadata)
    // setupRenditionFilters(mediaMetadata)

    ffmpegCommand := exec.Command("ffmpeg",
        "-i", "pipe:0",

        "-filter_complex", "[0:v]scale=854:480[v0];[0:v]scale=1280:720[v1];[0:v]scale=1920:1080[v2]",

        "-map", "[v0]", "-map", "0:a:0",
        "-c:v:0", "libx264", "-b:v:0", "800k", "-c:a:0", "aac",

        "-map", "[v1]", "-map", "0:a:0",
        "-c:v:1", "libx264", "-b:v:1", "2500k", "-c:a:1", "aac",

        "-map", "[v2]", "-map", "0:a:0",
        "-c:v:2", "libx264", "-b:v:2", "5000k", "-c:a:2", "aac",

        "-f", "hls",
        "-hls_time", "2",
        "-hls_list_size", "4",
        "-hls_flags", "append_list+delete_segments",
        "-hls_base_url", os.Getenv("S3_ENDPOINT")+"/"+os.Getenv("CDN_BUCKET_NAME")+"/"+mediaId+"/",
        "-hls_segment_filename", "./media/"+mediaId+"/%v-segment-%d.ts",
        "-var_stream_map", "v:0,a:0,name:480p v:1,a:1,name:720p v:2,a:2,name:1080p",
        "-master_pl_name", "master.m3u8",
        "./media/"+mediaId+"/%v.m3u8",
    )

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
