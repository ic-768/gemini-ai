# Gemini Streaming Server

A Go server that provides an API for streaming responses from Google's Gemini API. The server accepts prompts via HTTP POST requests and returns streamed responses using Server-Sent Events (SSE).

## Features

- Streaming responses from Gemini API using SSE
- Simple web interface for testing
- Support for client disconnection handling
- Environment variable configuration

## Project Structure

```
.
├── api/
│   └── client.go        # Gemini API client wrapper
├── handlers/
│   └── stream_handler.go # HTTP handlers for streaming responses
├── models/
│   └── models.go        # Request/response data structures
├── static/
│   └── index.html       # Web interface for testing
├── go.mod              # Go module definition
├── main.go             # Main application entry point
└── README.md           # This file
```

## Setup

1. Clone the repository
2. Install dependencies: `go mod download`
3. Create a `.env` file with your Gemini API key:
   ```
   GEMINI_API_KEY=your_api_key_here
   PORT=8080              # Optional, defaults to 8080
   GEMINI_MODEL=gemini-1.5-flash  # Optional, defaults to gemini-1.5-flash
   ```

## Running the Server

```bash
go run main.go
```

The server will start on port 8080 (or the port specified in the `.env` file).

## API Endpoints

### `/api/generate` (POST)

Accepts a JSON payload with a prompt and streams back the response:

```json
{
  "prompt": "Your prompt here"
}
```

### `/health` (GET)

Returns a simple health check response.

## Web Interface

Open your browser and navigate to `http://localhost:8080` to use the web interface for testing the API.

## License

MIT 