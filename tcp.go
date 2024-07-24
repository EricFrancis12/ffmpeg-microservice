package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
)

// HandleTCP handles TCP connections.
func handleTCP(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Handling TCP connection")

	// Create a buffer to store the incoming data
	// var buf bytes.Buffer

	// Read the FFmpeg command from the connection (assume it ends with a newline)
	command, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Printf("Failed to read command: %v", err)
		return
	}
	command = strings.TrimSpace(command)

	// Setup FFmpeg command
	ffmpegCmd := exec.Command("ffmpeg", strings.Split(command, " ")...)
	ffmpegCmd.Stdin = conn
	ffmpegCmd.Stdout = conn
	ffmpegCmd.Stderr = conn

	// Run FFmpeg command
	if err := ffmpegCmd.Run(); err != nil {
		log.Printf("FFmpeg command failed: %v", err)
	}
}

// Custom listener to differentiate between TCP and HTTP connections.
func customListener(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleTCP(conn)
	}
}
