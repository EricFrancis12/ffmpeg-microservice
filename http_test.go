package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatCommand(t *testing.T) {
	t.Run("Non-empty command", func(t *testing.T) {
		commandStr := "ffmpeg -f mp4 -i - -vf scale=100:50 -c:a copy -c:v libx264 -f flv ./output.flv"
		name, args := FormatCmd(commandStr)

		assert.Equal(t, name, "ffmpeg")
		assert.Equal(t, args, []string{"-f", "mp4", "-i", "-", "-vf", "scale=100:50", "-c:a", "copy", "-c:v", "libx264", "-f", "flv", "./output.flv"})
	})

	t.Run("Empty command", func(t *testing.T) {
		name, args := FormatCmd("")
		assert.Equal(t, name, "")
		assert.Equal(t, args, []string{})
	})

	t.Run("Name, but no args", func(t *testing.T) {
		name, args := FormatCmd("ffmpeg")
		assert.Equal(t, name, "ffmpeg")
		assert.Equal(t, args, []string{})
	})
}
