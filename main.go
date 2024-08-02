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
		httpPort = stringVar(FlagNameHttpPort, defaultHttpPort, "Port the HTTP Server will run on")
		tcpPort  = stringVar(FlagNameTcpPort, defaultTcpPort, "Port the TCP Server will run on")
	)

	httpServer := NewHTTPServer(httpPort)
	go func() {
		fmt.Println("HTTP Server starting on port " + httpPort)
		log.Fatal(httpServer.Run())
	}()

	tcpServer := NewTCPServer(tcpPort)
	fmt.Println("TCP Server starting listening for connections on port " + tcpPort)
	log.Fatal(tcpServer.Run())
}

func stringVar(name string, value string, usage string) string {
	var result string
	flag.StringVar(&result, name, value, usage)
	flag.Parse()
	return result
}
