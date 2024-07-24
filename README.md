# FFmpeg Microservice

The goal of this project is to create a microservice for FFmpeg whose input and output rely entirely on streaming, and does not hold any data in memory or write to the disc. This is to maintain functionality and scalability regardless of file size.

The service will accept input via one of three input types:

1. TCP: The client establishes a TCP connection with the service. The input data and the ffmpeg command are streamed from the client to the server, then to FFmpeg, and back to the client on the same TCP connection.
2. HTTP: The server receives HTTP requests from the client that consists of the input data, the ffmpeg command, and an output destination. The output destination is where the output file will be streamed.
3. Form Data: The same functionality as #2, except the input data will be multipart/form-data from an html form element.


## Usage

### Input Type 1: TCP
WIP: Coming soon

### Input Type 2: HTTP
WIP: Coming soon

### Input Type 3: Form Data

To run the Form Data submission of this project (in this primitive state), follow these steps:

1. Run the Go program:

   ```bash

    Go run .

   ```

2. Open the index.html file in a browser

3. Upload an mp4 video file to the input, and click the submit button

The video file will be scaled down and written to /output.flv
