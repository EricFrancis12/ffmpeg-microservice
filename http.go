package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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
