package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal("Client creation error:", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	model.SetMaxOutputTokens(4000)
	prompt := genai.Text("What is the color of the sky?")

	iter := model.GenerateContentStream(ctx, prompt)

	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Print each chunk as it arrives
		if len(resp.Candidates) > 0 {
			for _, part := range resp.Candidates[0].Content.Parts {
				if text, ok := part.(genai.Text); ok {
					fmt.Print(string(text))
				}
			}
		}
	}
	fmt.Println() // Add newline at end
}
