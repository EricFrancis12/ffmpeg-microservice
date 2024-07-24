package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	listenAddr string
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}

func (s *Server) Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/http", handleHTTP)
	router.HandleFunc("/form-data", handleFormData)
	router.HandleFunc("/demo", demoHandleReq)

	fmt.Println("HTTP Server starting on port " + s.listenAddr)
	return http.ListenAndServe(s.listenAddr, WithCors(router))
}

func WithCors(router *mux.Router) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	}).Handler(router)
}

// HandleHTTP handles HTTP requests.
func handleHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling HTTP request")

	// Read the FFmpeg command from the request header
	command := r.Header.Get("X-FFmpeg-Command")
	if command == "" {
		http.Error(w, "Missing FFmpeg command", http.StatusBadRequest)
		return
	}

	// Setup FFmpeg command
	ffmpegCmd := exec.Command("ffmpeg", strings.Split(command, " ")...)
	ffmpegCmd.Stdin = r.Body
	ffmpegCmd.Stdout = w
	ffmpegCmd.Stderr = w

	// Run FFmpeg command
	if err := ffmpegCmd.Run(); err != nil {
		http.Error(w, fmt.Sprintf("FFmpeg command failed: %v", err), http.StatusInternalServerError)
	}
}

// HandleFormData handles multipart/form-data requests.
func handleFormData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling multipart/form-data request")

	// Parse the multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}
	fmt.Println("~ 1")

	// Get the FFmpeg command from the form data
	// command := r.FormValue("command")
	// if command == "" {
	// 	http.Error(w, "Missing FFmpeg command", http.StatusBadRequest)
	// 	return
	// }

	fmt.Println("~ 2")

	// Get the input file from the form data
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get input file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Println("~ 3")

	// Setup FFmpeg command
	// ffmpegCmd := exec.Command("ffmpeg", strings.Split(command, " ")...)

	var (
		output = "./output.flv"
		// output = "http:// ..."
	)

	// Prepare the ffmpeg command
	ffmpegCmd := exec.Command("ffmpeg", "-f", "mp4", "-i", "-", "-vf", "scale=100:50", "-c:a", "copy", "-c:v", "libx264", "-f", "flv", output)

	ffmpegCmd.Stdin = file
	ffmpegCmd.Stdout = w
	ffmpegCmd.Stderr = w

	fmt.Println("~ 4")

	// Run FFmpeg command
	if err := ffmpegCmd.Run(); err != nil {
		fmt.Println(err.Error())
		http.Error(w, fmt.Sprintf("FFmpeg command failed: %v", err), http.StatusInternalServerError)
	}

	fmt.Println("~ 5")
}

func demoHandleReq(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New Request!!!")

	// Open the target file to save the processed video
	targetFile, err := os.Create("output/" + strconv.Itoa(int(time.Now().Unix())) + ".flv")
	if err != nil {
		log.Fatalf("Failed to create target file: %v", err)
	}
	defer targetFile.Close()

	// Write the response from FFmpeg to the target file
	_, err = io.Copy(targetFile, r.Body)
	if err != nil {
		log.Fatalf("Failed to write to target file: %v", err)
	}

	log.Println("Video streaming and processing completed successfully")
}
