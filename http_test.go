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

const (
	targetHeight = 50
	targetWidth  = 100
	inputPath    = "./video.mkv"
	tmpDir       = "./tmp"
)

func TestHandleHTTP(t *testing.T) {
	assert.Nil(t, MakeDirIfNotExists(tmpDir, os.ModePerm))
	assert.Nil(t, clearDir(tmpDir))

	inputFile, err := os.ReadFile(inputPath)
	assert.Nil(t, err)

	server := httptest.NewServer(http.HandlerFunc(handleHTTP))
	req, err := http.NewRequest("POST", server.URL, bytes.NewReader(inputFile))
	assert.Nil(t, err)

	req.Header.Set("Content-Type", "video/mkv")

	client := &http.Client{}

	t.Run("Write to file system", func(t *testing.T) {
		outputPath := fmt.Sprintf("%s/output-A.flv", tmpDir)
		command := fmt.Sprintf("ffmpeg -i - -vf scale=%d:%d -c:a copy -c:v libx264 -f flv %s", targetWidth, targetHeight, outputPath)

		req.Header.Set(HTTPHeaderCommand, command)

		resp, err := client.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()

		resolution, err := GetVideoResolution(outputPath)
		assert.Nil(t, err)
		assert.Equal(t, targetHeight, resolution.Height)
		assert.Equal(t, targetWidth, resolution.Width)

		assert.Nil(t, os.Remove(outputPath))
	})

	t.Run("Pipe response back to client", func(t *testing.T) {
		outputPath := fmt.Sprintf("%s/output-B.flv", tmpDir)
		command := fmt.Sprintf("ffmpeg -i - -vf scale=%d:%d -c:a copy -c:v libx264 -f flv pipe:1", targetWidth, targetHeight)

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
		assert.Equal(t, targetHeight, resolution.Height)
		assert.Equal(t, targetWidth, resolution.Width)

		assert.Nil(t, os.Remove(outputPath))
	})

	t.Run("Modify a local file", func(t *testing.T) {
		outputPath := fmt.Sprintf("%s/output-C.flv", tmpDir)
		command := fmt.Sprintf("ffmpeg -i %s -vf scale=%d:%d -c:a copy -c:v libx264 -f flv %s", inputPath, targetWidth, targetHeight, outputPath)

		req, err := http.NewRequest("POST", server.URL, nil)
		assert.Nil(t, err)

		req.Header.Set(HTTPHeaderCommand, command)

		resp, err := client.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()

		resolution, err := GetVideoResolution(outputPath)
		assert.Nil(t, err)
		assert.Equal(t, targetHeight, resolution.Height)
		assert.Equal(t, targetWidth, resolution.Width)

		assert.Nil(t, os.Remove(outputPath))
	})
}

func TestHandleFormData(t *testing.T) {
	assert.Nil(t, MakeDirIfNotExists(tmpDir, os.ModePerm))
	assert.Nil(t, clearDir(tmpDir))

	server := httptest.NewServer(http.HandlerFunc(handleFormData))

	client := &http.Client{}

	outputPath := fmt.Sprintf("%s/output-D.flv", tmpDir)
	command := fmt.Sprintf("ffmpeg -i - -vf scale=%d:%d -c:a copy -c:v libx264 -f flv %s", targetWidth, targetHeight, outputPath)

	file, err := os.Open("./video.mkv")
	assert.Nil(t, err)

	// Prepare the reader instances to encode
	values := map[string]io.Reader{}
	values[FormDataKeyFile] = file
	values[FormDataKeyCommand] = strings.NewReader(command)

	assert.Nil(t, uploadFormData(client, server.URL, values))

	assert.Nil(t, os.Remove(outputPath))
}

func uploadFormData(client *http.Client, url string, values map[string]io.Reader) (err error) {
	// Prepare a form that will be submitted to the url
	var b bytes.Buffer

	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add the video file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}

	}
	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	req.Header.Set(HTTPHeaderContentType, w.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}

	return
}
