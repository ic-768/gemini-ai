package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	
	"gemini/api"
	"gemini/handlers"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	// Check for required API key
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is not set")
	}

	// Get optional configuration from environment
	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = "gemini-1.5-flash" // Default model
	}

	// Create a Gemini client
	geminiClient, err := api.NewGeminiClient(apiKey, model, 4000)
	if err != nil {
		log.Fatal("Client creation error:", err)
	}
	defer geminiClient.Close()

	// Create handlers
	streamHandler := handlers.NewStreamHandler(geminiClient)
	
	// Configure routes
	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	
	// API endpoints
	http.Handle("/api/generate", streamHandler)
	http.HandleFunc("/health", handlers.HealthCheckHandler)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
