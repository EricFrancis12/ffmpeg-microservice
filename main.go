package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
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

func _main() {
	if err := EnsureDir("./output"); err != nil {
		log.Fatalf("Failed ensuring output directory: %v", err)
	}

	// cmd := exec.Command("ffmpeg", "-listen", "1", "-i", "http://localhost:8080/live/stream", "-vf", "scale=100:50", "-c:a", "copy", "-c:v", "libx264", "-f", "flv", "http://localhost:3002")
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err := cmd.Start()
	// if err != nil {
	// 	log.Fatalf("Failed to start ffmpeg: %v", err)
	// }

	// Open the input file
	inputFile, err := os.Open("video.mp4")
	if err != nil {
		log.Fatalf("Failed to open input file: %v", err)
	}
	defer inputFile.Close()

	var (
		output = ""
	)

	// Prepare the ffmpeg command
	cmd := exec.Command("ffmpeg", "-f", "mp4", "-i", "-", "-vf", "scale=100:50", "-c:a", "copy", "-c:v", "libx264", "-f", "flv", output)

	// Set the input file as the standard input for the ffmpeg command
	cmd.Stdin = inputFile

	// Run the command
	cmd.Run()
	if err != nil {
		log.Fatalf("ffmpeg command failed: %s", err)
	}

	// Give ffmpeg a moment to start up
	// time.Sleep(2 * time.Second)
	// go StreamInput()

	// server := NewServer(httpPort)
	// if err := server.Run(); err != nil {
	// 	log.Fatalf("Server failed: %v", err)
	// }
}

func StreamInput() {
	// Open the source video file
	sourceFile, err := os.Open("video.mp4")
	if err != nil {
		log.Fatalf("Failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	// Prepare the HTTP request to stream the file to FFmpeg
	req, err := http.NewRequest("POST", "http://localhost:8080/live/stream", sourceFile)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "video/mp4")

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to stream video to FFmpeg: %v", err)
	}
	defer resp.Body.Close()

	log.Println("Video streaming and processing completed successfully")
}

func EnsureDir(dirName string) error {
	exists, _ := DirExists(dirName)
	if exists {
		return nil
	}

	if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func DirExists(dirName string) (bool, error) {
	info, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}
