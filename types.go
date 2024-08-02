package main

type Resolution struct {
	Height int
	Width  int
}

type FFprobeResult struct {
	Programs []any        `json:"programs"`
	Streams  []Resolution `json:"streams"`
}

const FlagNameHttpPort string = "hport"

const ContentTypeApplicationOctetStream string = "application/octet-stream"

const URLQueryParamFormData string = "form-data"

const (
	HTTPHeaderAccept      string = "Accept"
	HTTPHeaderContentType string = "Content-Type"
	HTTPHeaderCommand     string = "X-Command"
)

const (
	FormDataKeyCommand string = "command"
	FormDataKeyFile    string = "file"
)
