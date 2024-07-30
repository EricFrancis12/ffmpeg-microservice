#!/bin/bash

# Script to keep FFmpeg listener running continuously

$INPUT="http://localhost:8080/live/stream"
OUTPUT="http://localhost:3003"

while true; do
  ffmpeg -listen 1 -i $INPUT -vf "scale=100:50" -c:a copy -c:v libx264 -f flv $OUTPUT
  sleep 1
done
