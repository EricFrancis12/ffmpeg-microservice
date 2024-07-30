package main

type Resolution struct {
	Height int
	Width  int
}

type FFprobeResult struct {
	Programs []any        `json:"programs"`
	Streams  []Resolution `json:"streams"`
}

const ContentTypeApplicationOctetStream = "application/octet-stream"

const URLQueryParamFormData = "form-data"

const (
	HTTPHeaderAccept        = "Accept"
	HTTPHeaderContentType   = "Content-Type"
	HTTPHeaderFFmpegCommand = "X-FFmpeg-Command"
)

const (
	FormDataKeyCommand = "command"
	FormDataKeyFile    = "file"
)
