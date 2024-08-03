test:
	go test ./... -v

build:
	go build -o bin/ffmpeg-microservice .

run:
	./bin/ffmpeg-microservice $(ARGS)

dl video:
	bash scripts/download_sample_video.sh
