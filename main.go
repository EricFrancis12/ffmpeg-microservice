package main

import (
	"flag"
	"fmt"
	"log"
)

const defaultHttpPort = ":3003"

func main() {
	httpPort := stringVar(FlagNameHttpPort, defaultHttpPort, "Port the HTTP Server will run on")

	httpServer := NewHTTPServer(httpPort)
	fmt.Println("HTTP Server starting on port " + httpPort)
	log.Fatal(httpServer.Run())
}

func stringVar(name string, value string, usage string) string {
	var result string
	flag.StringVar(&result, name, value, usage)
	flag.Parse()
	return result
}
