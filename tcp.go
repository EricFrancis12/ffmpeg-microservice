package main

import (
	"bufio"
	"log"
	"net"
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

	ffmpegCmd := PrepareCommand(command, conn, conn, conn)

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
