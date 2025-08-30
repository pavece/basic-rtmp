package transcoding

import (
	"io"
	"os"
	"os/exec"

	"github.com/pavece/simple-rtmp/internal/uploader"
)

//TODO: Add media metadata
//TODO: Use correct commands and infra for live HLS transcoding
func SetupTranscoder(mediaMetadata map[string]int) (*exec.Cmd, io.WriteCloser, error) {
	//TODO: setup ABR ladder and transcode to specific metadata
    ffmpegCommand := exec.Command("ffmpeg",
        "-i", "pipe:0",
        "-c:v", "copy",
        "-c:a", "aac",
        "-f", "hls",
        "-hls_time", "2",
        "-hls_list_size", "4",
        "-hls_flags", "append_list+delete_segments",
        "-hls_base_url", os.Getenv("S3_ENDPOINT") + "/" + os.Getenv("CDN_BUCKET_NAME") + "/", //TODO: Add media id
        "-hls_segment_filename", "./media/segment-%d.ts",
        "./media/master.m3u8",
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

	uploader.SetupFileWatcher()
    return ffmpegCommand, ffmpegPipe, nil
}
