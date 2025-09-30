package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pavece/simple-rtmp/internal/instrumentation"
	"github.com/pavece/simple-rtmp/internal/rtmp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	go servePrometheus()

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
	protocol := rtmp.New(connection)
	protocol.Handshake()
	
	for {
		err := protocol.ReadChunkData()
		if err != nil {
			break;
		}
	}
}

func servePrometheus(){
	port := os.Getenv("PROMETHEUS_PORT")
	if port == "" {
		return
	}

	fmt.Printf("Serving prometheus metrics on port %s\n", port)

	http.Handle("/metrics", promhttp.HandlerFor(instrumentation.Registry, promhttp.HandlerOpts{}))
	http.ListenAndServe(":"+port, nil)
}

func validateEnv() error {
	if os.Getenv("LOCAL_MEDIA_DIR") == "" {
		return fmt.Errorf("local media dir not specified")
	}

	return nil
}