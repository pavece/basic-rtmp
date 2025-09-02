package transcoding

import (
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

//TODO: Add media metadata
//TODO: Use correct commands and infra for live HLS transcoding
func SetupTranscoder(mediaMetadata map[string]int, mediaId string) (*exec.Cmd, io.WriteCloser, error) {
    createMediaFolder(mediaId)

	//TODO: setup ABR ladder and transcode to specific metadata
    ffmpegCommand := exec.Command("ffmpeg",
        "-i", "pipe:0",
        "-c:v", "copy",
        "-c:a", "aac",
        "-f", "hls",
        "-hls_time", "2",
        "-hls_list_size", "4",
        "-hls_flags", "append_list+delete_segments",
        "-hls_base_url", os.Getenv("S3_ENDPOINT") + "/" + os.Getenv("CDN_BUCKET_NAME") + "/" + mediaId + "/",
        "-hls_segment_filename", "./media/" + mediaId + "/segment-%d.ts",
        "./media/" + mediaId + "/master.m3u8",
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
