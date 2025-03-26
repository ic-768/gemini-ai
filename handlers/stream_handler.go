package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"

	// Use proper Go module imports
	"gemini/api"
	"gemini/models"
)

// StreamHandler handles SSE streaming of Gemini API responses
type StreamHandler struct {
	GeminiClient *api.GeminiClient
}

// NewStreamHandler creates a new handler with the provided Gemini client
func NewStreamHandler(client *api.GeminiClient) *StreamHandler {
	return &StreamHandler{
		GeminiClient: client,
	}
}

// ServeHTTP implements the http.Handler interface for streaming responses
func (h *StreamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for all responses
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON request
	var req models.PromptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.Prompt == "" {
		http.Error(w, "Prompt cannot be empty", http.StatusBadRequest)
		return
	}

	// Set up SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create a channel to signal when we're done
	// This is useful for handling client disconnections
	done := r.Context().Done()

	// Create a request context that will be canceled if client disconnects
	reqCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Watch for client disconnection
	go func() {
		<-done
		cancel()
	}()

	// Generate content stream
	iter := h.GeminiClient.GenerateContentStream(reqCtx, req.Prompt)

	// Stream back the response
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Send an initial event to confirm connection
	fmt.Fprintf(w, "event: connected\ndata: %s\n\n", time.Now().Format(time.RFC3339))
	flusher.Flush()

	for {
		resp, err := iter.Next()

		if err == iterator.Done {
			// Send a completion event
			fmt.Fprintf(w, "event: done\ndata: \n\n")
			flusher.Flush()
			break
		}

		if err != nil {
			// Send error event
			errMsg := fmt.Sprintf("Error: %v", err)
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", errMsg)
			log.Printf("Generation error: %v", err)
			flusher.Flush()
			break
		}

		// Process each chunk as it arrives
		if len(resp.Candidates) > 0 {
			for _, part := range resp.Candidates[0].Content.Parts {
				if text, ok := part.(genai.Text); ok {
					// Send text chunk as SSE event - properly encode to preserve whitespace
					textContent := string(text)

					// Log the text content for debugging (optional)
					log.Printf("Raw text chunk from API: %q", textContent)

					// Need to escape newlines in the data field of SSE
					// Replace any \n with \n + data: continuation
					textContent = strings.ReplaceAll(textContent, "\n", "\ndata: ")

					// Send the text event with raw content (no formatting that might alter whitespace)
					fmt.Fprintf(w, "event: text\ndata: %s\n\n", textContent)
					flusher.Flush()
				}
			}
		}
	}
}

// HealthCheckHandler provides a simple health check endpoint
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server is running"))
}

