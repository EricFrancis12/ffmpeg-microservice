# Listen on port and write output to the file system
ffmpeg -listen 1 -i http://localhost:8080/live/stream -vf "scale=100:50" -c:a copy -c:v libx264 -f flv [path]/output.flv

# Listen on port and stream output to http endpoint
ffmpeg -listen 1 -i http://localhost:8080/live/stream -vf "scale=100:50" -c:a copy -c:v libx264 -f flv http://localhost:3003

# Listen on port and stream output back as an http response (does not work)
ffmpeg -listen 1 -i http://localhost:8080/live/stream -vf "scale=100:50" -c:a copy -c:v libx264 -f flv pipe:1
