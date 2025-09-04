package main

import (
	"fmt"
	"net"

	"github.com/joho/godotenv"
	"github.com/pavece/simple-rtmp/internal/rtmp"
	"github.com/pavece/simple-rtmp/internal/uploader"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error while loading .env file: ", err)
		return
	}

	nl, err := net.Listen("tcp", ":1935")
	if err != nil {
		fmt.Println(err)
		return
	}

	uploader.FileUploaderInstance.SetupFileUploader()
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
	protocolStatus := rtmp.NewProtocolStatus()
	protocolStatus.Socket = connection
	
	for {
		rtmp.ReadChunkData(connection, protocolStatus)
	}
}