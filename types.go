package main

type Resolution struct {
	Height int
	Width  int
}

type FFprobeResult struct {
	Programs []any        `json:"programs"`
	Streams  []Resolution `json:"streams"`
}

const (
	ContentTypeApplicationOctetStream string = "application/octet-stream"
	ContentTypeVideoMKV               string = "video/mkv"
)

const (
	FlagNameAllowedOrigins string = "ao"
	FlagNameHttpPort       string = "hport"
)

const (
	FormDataKeyCommand string = "command"
	FormDataKeyFile    string = "file"
)

const (
	HTTPHeaderAccept      string = "Accept"
	HTTPHeaderCommand     string = "X-Command"
	HTTPHeaderContentType string = "Content-Type"
)

const URLQueryParamFormData string = "form-data"
