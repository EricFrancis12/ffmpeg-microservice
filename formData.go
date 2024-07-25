package main

import (
	"fmt"
	"net/http"
)

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
