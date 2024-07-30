package main

import (
	"io"
	"os/exec"
	"strings"
)

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
