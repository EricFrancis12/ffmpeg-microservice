package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

// HandleTCP handles TCP connections.
func handleTCP(conn net.Conn) {
	defer conn.Close()

	// Read the FFmpeg command from the connection (assume it ends with a newline)
	command, err := bufio.NewReader(conn).ReadString('\n')
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
