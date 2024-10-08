# FFmpeg Microservice

FFmpeg Microservice is a scalable and efficient solution designed for processing video & audio files using FFmpeg. This microservice can input/output entirely via streaming - it does not need to store data in memory or on disk in order to function, ensuring robust performance with large files.

It supports various input methods, including direct HTTP streaming, file system input, and multipart form data, and can stream output directly to the client.

## Quickstart

1. Clone the repository:
   
   ```bash
   git clone https://github.com/EricFrancis12/ffmpeg-microservice.git
   ```

2. Navigate to the project directory:
   
   ```bash
   cd ffmpeg-microservice
   ```

3. Build the application:
 
   ```bash
   make build
   ```

   This will create a binary file located at  `/bin/ffmpeg-microservice`.

4. Run the application:
   
   ```bash
   make run
   ```

   The service should now be running at http://localhost:3003 by default.


## Usage

The service accepts input via HTTP Post request, or as Multipart Form Data.

### Stream an input file via HTTP request body

Use `-i -` if you are sending the input file in the request body:

```bash
curl -X POST \
   -H "Content-Type: video/mkv" \
   # Specify the command that will run:
   -H "X-Command: ffmpeg -i - -vf scale=100:50 -c:a copy -c:v libx264 -f flv ./output-A.flv" \
   # The path to the input file:
   --data-binary @./video.mkv \
   http://localhost:3003
```

### Use an input file from the file system

Use `-i [path/to/file]` if you are referencing an input file in the file system:

```bash
curl -X POST \
   -H "X-Command: ffmpeg -i ./video.mkv -vf scale=100:50 -c:a copy -c:v libx264 -f flv ./output-B.flv" \
   http://localhost:3003
```

### Stream the output back as HTTP response

Use the header `"Accept": "application/octet-stream"` and `pipe:` to stream stdout back to the client:

```javascript
fetch("http://localhost:3003", {
   method: "POST",
   headers: {
      "Accept": "application/octet-stream",
      "X-Command": "ffmpeg -i ./video.mkv -vf scale=100:50 -c:a copy -c:v libx264 -f flv pipe:",
   },
})
   .then(res => res.blob())
   .then(blob => {
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.style.display = "none";
      a.href = url;
      a.download = "output-C.flv";
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
   });
```


### Send Multipart Form Data as input

Make sure to use `?form-data=1` to indicate you are sending multipart form data.
The input file needs to have the name `file`.
The command needs to have the name `command`.

```html
<form
   enctype="multipart/form-data"
   method="POST"
   action="http://localhost:3003?form-data=1"
>
   <!-- File Input -->
   <input
      type="file"
      name="file"
      required
   >
   <!-- Command Input -->
   <input
      type="text"
      name="command"
      value="ffmpeg -i - -vf scale=100:50 -c:a copy -c:v libx264 -f flv ./output-D.flv"
      required
   >
   <button type="submit">
      Submit
   </button>
</form>
```


## Testing

1. Download the sample video:
   
   ```bash
   make dl video
   ```

2. Run the test suite:

   ```bash
   make test
   ```


## Find a bug?
If you found an issue or would like to submit an improvement to this project, please submit an issue using the issues tab above. If you would like to submit a PR with a fix, reference the issue you created.
