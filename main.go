package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

const (
	defaultHttpPort       = ":3003"
	defaultAllowedOrigins = "*"
)

func main() {
	hport := stringFlag(FlagNameHttpPort, defaultHttpPort, "Port the HTTP Server will run on")
	ao := stringsFlag(FlagNameAllowedOrigins, "CORS allowed origins")
	flag.Parse()

	var (
		httpPort       = *hport
		allowedOrigins = []string(*ao)
	)

	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{defaultAllowedOrigins}
	}

	httpServer := NewHTTPServer(httpPort, allowedOrigins)
	fmt.Printf("HTTP Server starting on port %s", httpPort)
	log.Fatal(httpServer.Run())
}

type flagsSlice []string

func (i flagsSlice) String() string {
	return strings.Join([]string(i), ", ")
}

func (i *flagsSlice) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func stringsFlag(name string, usage string) *flagsSlice {
	var result flagsSlice
	flag.Var(&result, name, usage)
	return &result
}

func stringFlag(name string, defaultValue string, usage string) *string {
	var result string
	flag.StringVar(&result, name, defaultValue, usage)
	return &result
}
