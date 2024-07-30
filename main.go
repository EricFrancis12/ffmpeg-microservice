package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	tcpPort  = ":8080"
	httpPort = ":3003"
)

func main() {
	listener, err := net.Listen("tcp", tcpPort)
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}
	defer listener.Close()

	go func() {
		server := NewServer(httpPort)

		err := server.Run()
		if err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	fmt.Println("Listening for TCP connections on port " + tcpPort)
	customListener(listener)
}

func PrepareCmd(cmdStr string, stdin io.ReadCloser, stdout io.Writer, stderr io.Writer) *exec.Cmd {
	cmdStr = strings.TrimSpace(cmdStr)
	name, args := FormatCmd(cmdStr)
	cmd := exec.Command(name, args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd
}

func FormatCmd(str string) (name string, args []string) {
	parts := strings.Split(str, " ")
	return parts[0], parts[1:]
}

func WithCors(router *mux.Router) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	}).Handler(router)
}
