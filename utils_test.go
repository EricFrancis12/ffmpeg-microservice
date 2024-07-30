package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatCommand(t *testing.T) {
	t.Run("Non-empty command", func(t *testing.T) {
		commandStr := "ffmpeg -f mp4 -i - -vf scale=100:50 -c:a copy -c:v libx264 -f flv ./output.flv"
		name, args := FormatCommand(commandStr)

		assert.Equal(t, name, "ffmpeg")
		assert.Equal(t, args, []string{"-f", "mp4", "-i", "-", "-vf", "scale=100:50", "-c:a", "copy", "-c:v", "libx264", "-f", "flv", "./output.flv"})
	})

	t.Run("Empty command", func(t *testing.T) {
		name, args := FormatCommand("")
		assert.Equal(t, name, "")
		assert.Equal(t, args, []string{})
	})

	t.Run("Name, but no args", func(t *testing.T) {
		name, args := FormatCommand("ffmpeg")
		assert.Equal(t, name, "ffmpeg")
		assert.Equal(t, args, []string{})
	})
}

func TestParseResolution(t *testing.T) {
	t.Run("Correct JSON format", func(t *testing.T) {
		resolution, err := ParseResolution("{\"programs\":[],\"streams\":[{\"height\":10,\"width\":20}]}")
		assert.Nil(t, err)
		assert.Equal(t, resolution, Resolution{Height: 10, Width: 20})
	})

	t.Run("Empty string", func(t *testing.T) {
		resolution, err := ParseResolution("")
		assert.NotNil(t, err)
		assert.Equal(t, resolution, Resolution{})
	})
}

func TestParseJSON(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("Test Parse JSON to Map", func(t *testing.T) {
		jsonStr := `{"name": "John", "age": 30}`
		expectedMap := map[string]interface{}{"name": "John", "age": float64(30)}

		resultMap, err := ParseJSON[map[string]interface{}](jsonStr)
		assert.NoError(t, err)
		assert.Equal(t, expectedMap, resultMap)
	})

	t.Run("Test Parse JSON to Struct", func(t *testing.T) {
		jsonStr := `{"name": "John", "age": 30}`
		expectedStruct := Person{Name: "John", Age: 30}

		resultStruct, err := ParseJSON[Person](jsonStr)
		assert.NoError(t, err)
		assert.Equal(t, expectedStruct, resultStruct)
	})

	t.Run("Test Parse Invalid JSON", func(t *testing.T) {
		invalidJsonStr := `{"name": "John", "age": 30`
		var expectedInvalidMap map[string]interface{}

		resultInvalidMap, err := ParseJSON[map[string]interface{}](invalidJsonStr)
		assert.Error(t, err)
		assert.Equal(t, expectedInvalidMap, resultInvalidMap)
	})

	t.Run("Test Parse JSON to Slice", func(t *testing.T) {
		jsonStrSlice := `["apple", "banana", "cherry"]`
		expectedSlice := []string{"apple", "banana", "cherry"}

		resultSlice, err := ParseJSON[[]string](jsonStrSlice)
		assert.NoError(t, err)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("Test Parse JSON to Int", func(t *testing.T) {
		jsonStrInt := `123`
		expectedInt := 123

		resultInt, err := ParseJSON[int](jsonStrInt)
		assert.NoError(t, err)
		assert.Equal(t, expectedInt, resultInt)
	})

	t.Run("Test Parse Empty JSON String", func(t *testing.T) {
		emptyJsonStr := ``
		var expectedEmptyMap map[string]interface{}

		resultEmptyMap, err := ParseJSON[map[string]interface{}](emptyJsonStr)
		assert.Error(t, err)
		assert.Equal(t, expectedEmptyMap, resultEmptyMap)
	})
}
