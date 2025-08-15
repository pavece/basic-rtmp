package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/pavece/simple-rtmp/internal/rtmp"
)

func main() {
	nl, err := net.Listen("tcp", ":1935")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Basic RTMP server started")

	for {
        connection, err := nl.Accept()
        if err != nil {
            fmt.Println(err)
            continue
        }

        go handleConnection(connection)
    }
}

func handleConnection(connection net.Conn){
	defer connection.Close()

	rtmp.Handshake(connection)
	rtmp.FfmpegCommand.Stderr = os.Stderr
	rtmp.FfmpegCommand.Stdout = os.Stdout

	err := rtmp.FfmpegCommand.Start()
	
	if err != nil {
		log.Fatal(err)
	}
	
	for {
		rtmp.ReadChunkData(connection)
	}
}