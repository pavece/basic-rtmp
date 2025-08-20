package rtmp

import (
	"crypto/rand"
	"io"
	"log"
	"net"
)

type HandshakeData struct {
	Version int
	InitialTimestamp [4]byte
}

func Handshake(conn net.Conn) error {
    var c0 [1]byte
    var c1 [1536]byte
    var c2 [1536]byte
    var s0 [1]byte
    var s1 [1536]byte
    var s2 [1536]byte

    if _, err := io.ReadFull(conn, c0[:]); err != nil {
       log.Fatal(err)
    }
    
    s0[0] = 3 //Version
    copy(s1[:4], []byte{0,0,0,0})
    copy(s1[4:8], []byte{0,0,0,0}) 
    if _, err := io.ReadFull(rand.Reader, s1[8:]); err != nil {
       log.Fatal(err)
    }
    if _, err := conn.Write(s0[:]); err != nil {
       log.Fatal(err)
    }
    if _, err := conn.Write(s1[:]); err != nil {
       log.Fatal(err)
    }

    if _, err := io.ReadFull(conn, c1[:]); err != nil {
       log.Fatal(err)
    }

   copy(s2[:], c1[:])
    if _, err := conn.Write(s2[:]); err != nil {
       log.Fatal(err)
    }

    if _, err := io.ReadFull(conn, c2[:]); err != nil {
       log.Fatal(err)
    }

    return nil
}
