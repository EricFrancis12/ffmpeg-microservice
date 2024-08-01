package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleHTTP(t *testing.T) {
	const (
		height = 50
		width  = 100
	)

	err := MakeDirIfNotExists("./tmp", os.ModePerm)
	assert.Nil(t, err)

	inputFile, err := os.ReadFile("./video.mp4")
	assert.Nil(t, err)

	server := httptest.NewServer(http.HandlerFunc(handleHTTP))
	req, err := http.NewRequest("POST", server.URL, bytes.NewReader(inputFile))
	assert.Nil(t, err)

	req.Header.Set("Content-Type", "video/mp4")

	client := &http.Client{}

	t.Run("Write to file system", func(t *testing.T) {
		outputPath := "./tmp/output-A.flv"
		command := fmt.Sprintf("ffmpeg -f mp4 -i - -vf scale=%d:%d -c:a copy -c:v libx264 -f flv %s", width, height, outputPath)

		req.Header.Set(HTTPHeaderCommand, command)

		resp, err := client.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()

		resolution, err := GetVideoResolution(outputPath)
		assert.Nil(t, err)
		assert.Equal(t, height, resolution.Height)
		assert.Equal(t, width, resolution.Width)

		assert.Nil(t, os.Remove(outputPath))
	})

	t.Run("Pipe response back to client", func(t *testing.T) {
		outputPath := "./tmp/output-B.flv"
		command := fmt.Sprintf("ffmpeg -f mp4 -i - -vf scale=%d:%d -c:a copy -c:v libx264 -f flv pipe:1", width, height)

		req.Header.Set(HTTPHeaderCommand, command)
		req.Header.Set(HTTPHeaderAccept, ContentTypeApplicationOctetStream)

		resp, err := client.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()

		file, err := os.Create(outputPath)
		assert.Nil(t, err)

		_, err = io.Copy(file, resp.Body)
		assert.Nil(t, err)
		file.Close()

		resolution, err := GetVideoResolution(outputPath)
		assert.Nil(t, err)
		assert.Equal(t, height, resolution.Height)
		assert.Equal(t, width, resolution.Width)

		assert.Nil(t, os.Remove(outputPath))
	})
}
