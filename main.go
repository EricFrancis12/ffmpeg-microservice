package main

import (
	"flag"
	"fmt"
	"log"
)

const (
	defaultHttpPort = ":3003"
	defaultTcpPort  = ":8080"
)

func main() {
	var (
		httpPort string
		tcpPort  string
	)

	flag.StringVar(&httpPort, "hport", defaultHttpPort, "Port the HTTP Server will run on")
	flag.StringVar(&tcpPort, "tport", defaultTcpPort, "Port the TCP Server will run on")
	flag.Parse()

	httpServer := NewHTTPServer(httpPort)
	go func() {
		fmt.Println("HTTP Server starting on port " + httpPort)
		log.Fatal(httpServer.Run())
	}()

	tcpServer := NewTCPServer(tcpPort)
	fmt.Println("TCP Server starting listening for connections on port " + tcpPort)
	log.Fatal(tcpServer.Run())
}
