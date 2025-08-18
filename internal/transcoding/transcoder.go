package transcoding

import (
	"io"
	"os"
	"os/exec"
)

//TODO: Add media metadata
//TODO: Use correct commands and infra for live HLS transcoding
func SetupTranscoder() (*exec.Cmd, io.WriteCloser, error){
	ffmpegCommand := exec.Command("ffmpeg",
		"-i", "pipe:0",   
		"-c:v", "copy",   
		"-c:a", "copy",   
		"output.mp4",     
	)

	ffmpegCommand.Stderr = os.Stderr
	ffmpegCommand.Stdout = os.Stdout

	ffmpegPipe, err := ffmpegCommand.StdinPipe()
	if err != nil {
		return nil, nil, err
	}

	err = ffmpegCommand.Start()
	return ffmpegCommand, ffmpegPipe, err
}