package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Unmarshals a JSON string into a variable of the provided type.
func ParseJSON[T any](jsonStr string) (T, error) {
	var v T
	err := json.Unmarshal([]byte(jsonStr), &v)
	if err != nil {
		return v, err
	}
	return v, nil
}

func DirExists(dirPath string) (bool, error) {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func MakeDirIfNotExists(dirPath string, perm os.FileMode) error {
	dirExists, err := DirExists(dirPath)
	if err != nil {
		return err
	}

	if dirExists {
		return nil
	}

	if err := os.Mkdir(dirPath, perm); err != nil {
		return err
	}

	return nil
}

func PrepareCmd(command string, stdin io.ReadCloser, stdout io.Writer, stderr io.Writer) *exec.Cmd {
	command = strings.TrimSpace(command)
	name, args := FormatCommand(command)
	cmd := exec.Command(name, args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd
}

func FormatCommand(command string) (name string, args []string) {
	parts := strings.Split(command, " ")
	return parts[0], parts[1:]
}

func GetVideoResolution(filePath string) (Resolution, error) {
	command := fmt.Sprintf("ffprobe -v error -select_streams v -show_entries stream=width,height -of json %s", filePath)
	name, args := FormatCommand(command)
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return Resolution{}, err
	}

	resolution, err := ParseResolution(out.String())
	if err != nil {
		return Resolution{}, err
	}

	return resolution, nil
}

func ParseResolution(jsonStr string) (Resolution, error) {
	ffprobeResult, err := ParseJSON[FFprobeResult](jsonStr)
	if err != nil {
		return Resolution{}, err
	}

	if len(ffprobeResult.Streams) < 1 {
		return Resolution{}, fmt.Errorf("no streams in FFprobe output")
	}

	return ffprobeResult.Streams[0], nil
}

func clearDir(dirPath string) error {
	fi, err := os.Stat(dirPath)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("dirPath is not a directory")
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryPath := filepath.Join(dirPath, entry.Name())
		if entry.IsDir() {
			return clearDir(entryPath)
		}
		return os.Remove(entryPath)
	}

	return nil
}
