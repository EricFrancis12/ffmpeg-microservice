package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareHttpTest(t *testing.T, handler http.Handler, tmpDir string) (*httptest.Server, *http.Client) {
	assert.Nil(t, makeDirIfNotExists(tmpDir, os.ModePerm))
	assert.Nil(t, clearDir(tmpDir))

	server := httptest.NewServer(handler)
	client := &http.Client{}
	return server, client
}

func TestHandleHTTP(t *testing.T) {
	var (
		targetHeight = 50
		targetWidth  = 100
		inputPath    = "./video.mkv"
		tmpDir       = "./tmp"
	)

	server, client := prepareHttpTest(t, http.HandlerFunc(handleHTTP), tmpDir)

	inputFile, err := os.ReadFile(inputPath)
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", server.URL, bytes.NewReader(inputFile))
	assert.Nil(t, err)

	req.Header.Set("Content-Type", "video/mkv")

	t.Run("Write to file system", func(t *testing.T) {
		outputPath := fmt.Sprintf("%s/output-A.flv", tmpDir)
		command := fmt.Sprintf("ffmpeg -i - -vf scale=%d:%d -c:a copy -c:v libx264 -f flv %s", targetWidth, targetHeight, outputPath)

		req.Header.Set(HTTPHeaderCommand, command)

		resp, err := client.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()

		checkResolution(t, outputPath, targetHeight, targetWidth)

		assert.Nil(t, os.Remove(outputPath))
	})

	t.Run("Modify a local file", func(t *testing.T) {
		outputPath := fmt.Sprintf("%s/output-B.flv", tmpDir)
		command := fmt.Sprintf("ffmpeg -i %s -vf scale=%d:%d -c:a copy -c:v libx264 -f flv %s", inputPath, targetWidth, targetHeight, outputPath)

		req, err := http.NewRequest("POST", server.URL, nil)
		assert.Nil(t, err)

		req.Header.Set(HTTPHeaderCommand, command)

		resp, err := client.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()

		checkResolution(t, outputPath, targetHeight, targetWidth)

		assert.Nil(t, os.Remove(outputPath))
	})

	t.Run("Pipe response back to client", func(t *testing.T) {
		outputPath := fmt.Sprintf("%s/output-C.flv", tmpDir)
		command := fmt.Sprintf("ffmpeg -i - -vf scale=%d:%d -c:a copy -c:v libx264 -f flv pipe:", targetWidth, targetHeight)

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

		checkResolution(t, outputPath, targetHeight, targetWidth)

		assert.Nil(t, os.Remove(outputPath))
	})
}

type FormDataMap = map[string]io.Reader

func TestHandleFormData(t *testing.T) {
	var (
		targetHeight = 50
		targetWidth  = 100
		inputPath    = "./video.mkv"
		tmpDir       = "./tmp"
	)

	server, client := prepareHttpTest(t, http.HandlerFunc(handleFormData), tmpDir)

	t.Run("Write to file system", func(t *testing.T) {
		outputPath := fmt.Sprintf("%s/output-D.flv", tmpDir)
		command := fmt.Sprintf("ffmpeg -i - -vf scale=%d:%d -c:a copy -c:v libx264 -f flv %s", targetWidth, targetHeight, outputPath)

		file, err := os.Open(inputPath)
		assert.Nil(t, err)

		// Prepare the reader instances to encode
		fdm := make(FormDataMap)
		fdm[FormDataKeyFile] = file
		fdm[FormDataKeyCommand] = strings.NewReader(command)

		uploadFormData(t, client, server.URL, fdm)

		checkResolution(t, outputPath, targetHeight, targetWidth)

		assert.Nil(t, os.Remove(outputPath))
	})
}

func checkResolution(t *testing.T, filePath string, height int, width int) {
	resolution, err := GetVideoResolution(filePath)
	assert.Nil(t, err)
	assert.Equal(t, height, resolution.Height)
	assert.Equal(t, width, resolution.Width)
}

func uploadFormData(t *testing.T, client *http.Client, url string, fdm FormDataMap) {
	// Prepare a form that will be submitted to the url
	var b bytes.Buffer

	mpw := multipart.NewWriter(&b)
	for key, rdr := range fdm {
		var wrtr io.Writer
		if clsr, ok := rdr.(io.Closer); ok {
			defer clsr.Close()
		}
		// Add the video file
		if file, ok := rdr.(*os.File); ok {
			w, err := mpw.CreateFormFile(key, file.Name())
			assert.Nil(t, err)
			wrtr = w
		} else {
			// Add other fields
			w, err := mpw.CreateFormField(key)
			assert.Nil(t, err)
			wrtr = w
		}
		_, err := io.Copy(wrtr, rdr)
		assert.Nil(t, err)
	}
	mpw.Close()

	req, err := http.NewRequest("POST", url, &b)
	assert.Nil(t, err)

	req.Header.Set(HTTPHeaderContentType, mpw.FormDataContentType())

	res, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
