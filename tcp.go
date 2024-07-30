package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

const delim byte = '\n'

type TCPServer struct {
	ListenAddr string
	Listener   net.Listener
}

func NewTCPServer(listenAddr string) *TCPServer {
	return &TCPServer{
		ListenAddr: listenAddr,
	}
}

func (ts *TCPServer) Listen() error {
	listener, err := net.Listen("tcp", ts.ListenAddr)
	if err != nil {
		return err
	}

	ts.Listener = listener
	return nil
}

func (ts *TCPServer) Run() error {
	err := ts.Listen()
	if err != nil {
		return err
	}

	for {
		conn, err := ts.Listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleTCP(conn)
	}
}

func handleTCP(conn net.Conn) {
	defer conn.Close()

	// Read the FFmpeg command from the connection (assume it ends with a newline)
	command, err := bufio.NewReader(conn).ReadString(delim)
	if err != nil {
		log.Printf("Failed to read command: %v", err)
		return
	}
	command = command[:len(command)-1] // Remove the newline character

	ffmpegCmd := PrepareCmd(command, conn, conn, os.Stderr)

	// Run FFmpeg command
	if err := ffmpegCmd.Run(); err != nil {
		log.Printf("FFmpeg command failed: %v", err)
	}
}
