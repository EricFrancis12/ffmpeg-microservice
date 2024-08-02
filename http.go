package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	ListenAddr string
}

func NewHTTPServer(listenAddr string) *HTTPServer {
	return &HTTPServer{
		ListenAddr: listenAddr,
	}
}

func (hs *HTTPServer) Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/", handlePost).Methods("POST")

	return http.ListenAndServe(hs.ListenAddr, router)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get(URLQueryParamFormData) == "1" {
		handleFormData(w, r)
	} else {
		handleHTTP(w, r)
	}
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	// Read the command from the request header
	command := r.Header.Get(HTTPHeaderCommand)
	if command == "" {
		http.Error(w, "Missing command", http.StatusBadRequest)
		return
	}

	// Check if client is requesting the output to be streamed back as the response.
	// If so, the stdout of the cmd is set to w
	var stdout io.Writer = os.Stderr
	if r.Header.Get(HTTPHeaderAccept) == ContentTypeApplicationOctetStream {
		stdout = w
	}

	cmd := PrepareCmd(command, r.Body, stdout, os.Stderr)

	if err := cmd.Run(); err != nil {
		http.Error(w, fmt.Sprintf("command failed: %v", err), http.StatusInternalServerError)
	}
}

func handleFormData(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Get the command from the form data
	command := r.FormValue(FormDataKeyCommand)
	if command == "" {
		http.Error(w, "Missing command", http.StatusBadRequest)
		return
	}

	// Get the input file from the form data
	file, _, err := r.FormFile(FormDataKeyFile)
	if err != nil {
		http.Error(w, "Failed to get input file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	cmd := PrepareCmd(command, file, w, os.Stderr)

	if err := cmd.Run(); err != nil {
		http.Error(w, fmt.Sprintf("command failed: %v", err), http.StatusInternalServerError)
	}
}
