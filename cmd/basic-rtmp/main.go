package main

import (
	"fmt"
	"net"

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

	hsData := rtmp.Handshake(connection)
	rtmp.ParseBasicHeader(connection)

	fmt.Println(hsData)

}