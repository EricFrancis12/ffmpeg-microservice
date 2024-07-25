package main

import (
	"fmt"
	"net/http"

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

	router.PathPrefix("/public").Handler(http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))

	fmt.Println("HTTP Server starting on port " + s.listenAddr)
	return http.ListenAndServe(s.listenAddr, WithCors(router))
}

// HandleHTTP handles HTTP requests.
func handleHTTP(w http.ResponseWriter, r *http.Request) {
	// Read the FFmpeg command from the request header
	command := r.Header.Get("X-FFmpeg-Command")
	if command == "" {
		http.Error(w, "Missing FFmpeg command", http.StatusBadRequest)
		return
	}

	ffmpegCmd := PrepareCommand(command, r.Body, w, w)

	// Run FFmpeg command
	if err := ffmpegCmd.Run(); err != nil {
		http.Error(w, fmt.Sprintf("FFmpeg command failed: %v", err), http.StatusInternalServerError)
	}
}

// HandleFormData handles multipart/form-data requests.
func handleFormData(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Get the FFmpeg command from the form data
	command := r.FormValue("command")
	if command == "" {
		http.Error(w, "Missing FFmpeg command", http.StatusBadRequest)
		return
	}

	// Get the input file from the form data
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get input file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ffmpegCmd := PrepareCommand(command, file, w, w)

	// Run FFmpeg command
	if err := ffmpegCmd.Run(); err != nil {
		http.Error(w, fmt.Sprintf("FFmpeg command failed: %v", err), http.StatusInternalServerError)
	}
}

func WithCors(router *mux.Router) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	}).Handler(router)
}
