package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	nl, err := net.Listen("tcp", ":1935")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("Basic RTMP server started")

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
	// defer connection.Close()
	var version [1]byte
    var c1 [1536]byte

    if _, err := io.ReadFull(connection, version[:]); err != nil {
       return
    }

    fmt.Println("Version (C0): ", version[0])

    if _, err := io.ReadFull(connection, c1[:]); err != nil {
       return
    }

	fmt.Println("Version (C0): ", version[0])
	fmt.Println("Timestamp (C1): ", c1[0:4])
	fmt.Println("Zero (C1): ", c1[4:8])
	fmt.Println("Random (C1):", c1[8:])



}