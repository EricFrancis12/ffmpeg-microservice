#!/bin/bash

# Downloads a sample video for use in the testing suite.

VIDEO_URL="https://sample-videos.com/video321/mp4/720/big_buck_bunny_720p_1mb.mp4"
OUTPUT_FILE="video.mp4"

curl -o "$OUTPUT_FILE" "$VIDEO_URL"

if [ $? -eq 0 ]; then
    echo "Download complete: $OUTPUT_FILE"
else
    echo "Download failed with status $?."
fi
