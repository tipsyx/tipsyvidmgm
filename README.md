
#TIPSY Video Management API

This Video Management API is a web service that allows you to upload, transcribe, and manage video files. It provides various endpoints for uploading videos, listing uploaded videos, displaying video details, and more.

## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
 
- [Usage](#usage)
  - [Uploading Videos](#uploading-videos)
  - [Listing Videos](#listing-videos)
  - [Displaying Video Details](#displaying-video-details)
  - [Deleting Videos](#deleting-videos)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

## Getting Started

### Prerequisites

Before using the API, make sure you have the following prerequisites installed:

- Go programming language (at least Go 1.16)
- MySQL database
- RabbitMQ server
- Whisper API (for transcription)

## Usage

### Uploading Videos

To upload a video, send a POST request to `/upload` with the video file as a form field. The API will enqueue a transcription task for the video.

Example:

```bash
curl -F "video=@path/to/video.mp4" http://localhost:8080/upload
```

### Listing Videos

To list all uploaded videos, send a GET request to `/listvideos`. You can specify the page and items per page as query parameters.

Example:

```bash
curl http://localhost:8080/listvideos?page=1&itemsPerPage=10
```

### Displaying Video Details

To display details of a specific video, send a GET request to `/videodetails/{video_id}` where `{video_id}` is the ID of the video.

Example:

```bash
curl http://localhost:8080/videodetails/1
```

### Deleting Videos

To delete a video, send a POST request to `/deletevideo` with the `id` of the video to delete as a form field.

Example:

```bash
curl -X POST -F "id=1" http://localhost:8080/deletevideo
```

## Configuration

- Edit the `config/config.go` file to set up your database and RabbitMQ configurations.

## Contributing

Contributions are welcome! If you would like to contribute to this project, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them.
4. Push your changes to your fork.
5. Create a pull request to the original repository.

## License

This project is licensed under the [MIT License](LICENSE).


