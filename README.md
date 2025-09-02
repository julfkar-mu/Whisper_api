# Whisper Streaming API - Golang

A clean, maintainable Golang REST API for streaming audio transcription using the OpenAI Whisper model. The service streams real-time transcription results from Whisper and exposes them via a REST endpoint.

---

## Features

- Golang-based API adhering to SOLID principles for clean architecture and testability  
- Streaming transcription handling using OpenAI Whisper API  
- Docker containerized for easy deployment  
- REST API contract with JSON request/response  
- Uses Gorilla Mux for routing  
- Graceful handling of streaming data and context cancellations  

---

## Project Structure
/whisperapi
/cmd
/server
main.go # Application entrypoint
/internal
/api
handler.go # HTTP handler logic
routes.go # Route registration
dto.go # API request/response types
/service
whisper.go # Business logic for streaming
/client
whisper_client.go # OpenAI Whisper API client implementation
/pkg
/model
transcription.go # Common models
Dockerfile # Container setup
README.md # Project documentation



---

## Getting Started

### Prerequisites

- Go 1.21+ installed  
- Docker (optional, for containerized deployment)  
- OpenAI API key with access to Whisper endpoints  

### Environment

Set the `OPENAI_API_KEY` environment variable with your OpenAI API key.

export OPENAI_API_KEY="your_openai_key_here"


### Build & Run Locally

git clone <repo-url>
cd whisperapi/cmd/server
go build -o server .
./server


The server listens on `http://localhost:8080`.

### API Usage

- **Endpoint:** `POST /api/transcribe`  
- **Content-Type:** `application/json`  
- **Request Body:**  
{
"audio_url": "https://example.com/audio.mp3"
}

- **Response:**  
Streaming JSON responses with partial transcription text and completion status.
{
"text": "transcribed snippet",
"done": false
}

Final response will have `"done": true`.

---

## Running with Docker

Build the Docker image:

docker build -t whisperapi .


Run the container:
docker run -e OPENAI_API_KEY="your_openai_key_here" -p 8080:8080 whisperapi


---

## Testing

Unit tests are included. Run with:

go test ./...


---

## Contributing

Feel free to open issues or pull requests for feature suggestions, bug fixes, or improvements.

---

## License

MIT License - see LICENSE file for details.

---




