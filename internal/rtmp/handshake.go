package rtmp

import (
	"io"
	"log"
	"math/rand/v2"
	"net"
)

type HandshakeData struct {
	Version int
	InitialTimestamp [4]byte
}

func Handshake(connection net.Conn) HandshakeData {
	var handshakeData = HandshakeData{}

	var c0 [1]byte
    var c1 [1536]byte
	var c2 [1536]byte

    if _, err := io.ReadFull(connection, c0[:]); err != nil {
		log.Fatal("Error reading c0")
    }


	s0 := make([]byte, 1)
	s1 := make([]byte, 1536)
	s2 := make([]byte, 1536)

	s0[0] = byte(3)
	handshakeData.Version = 3
	
	copy(s1[0:4], []byte{0, 0, 0, 0})
	copy(s1[4:8], []byte{0, 0, 0, 0})
	for i := 8; i<1536; i++ {
		s1[i] = byte(rand.IntN(254))
	}

	connection.Write(s0)
	connection.Write(s1)
	

	if _, err := io.ReadFull(connection, c1[:]); err != nil {
		log.Fatal("Error reading c1")
    }
	
	copy(s2[0:4], c1[0:4])
	copy(s2[4:8], c1[0:4])
	copy(s2[8:], c1[8:])

	copy(handshakeData.InitialTimestamp[:], c1[0:4])

	connection.Write(s2)

	if _, err := io.ReadFull(connection, c2[:]); err != nil {
		log.Fatal("Error reading c2")
	}

	return handshakeData
}