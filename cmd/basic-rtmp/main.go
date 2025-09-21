package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/pavece/simple-rtmp/internal/rtmp"
)

func main() {
	godotenv.Load()

	err := validateEnv()
	if err != nil {
		log.Fatal(err)
	}

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
	rtmp.Handshake(connection)
	protocolStatus := rtmp.NewProtocolStatus()
	protocolStatus.Socket = connection
	
	for {
		err := rtmp.ReadChunkData(connection, protocolStatus)
		if err != nil {
			break;
		}
	}
}

func validateEnv() error {
	if os.Getenv("LOCAL_MEDIA_DIR") == "" {
		return fmt.Errorf("local media dir not specified")
	}

	return nil
}