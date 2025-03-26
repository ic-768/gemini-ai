package api

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiClient wraps the Gemini API client with configuration
type GeminiClient struct {
	Client *genai.Client
	Model  string
	Tokens int
}

// NewGeminiClient creates a new Gemini client with the provided API key
func NewGeminiClient(apiKey string, model string, maxTokens int) (*GeminiClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	// Default values if not provided
	if model == "" {
		model = "gemini-1.5-flash"
	}
	if maxTokens <= 0 {
		maxTokens = 4000
	}

	return &GeminiClient{
		Client: client,
		Model:  model,
		Tokens: maxTokens,
	}, nil
}

// Close closes the underlying Gemini client
func (g *GeminiClient) Close() {
	if g.Client != nil {
		g.Client.Close()
	}
}

// GenerateContentStream creates a stream for generating content with the specified prompt
func (g *GeminiClient) GenerateContentStream(ctx context.Context, prompt string) *genai.GenerateContentResponseIterator {
	generativeModel := g.Client.GenerativeModel(g.Model)
	generativeModel.SetMaxOutputTokens(int32(g.Tokens))

	textPrompt := genai.Text(prompt)
	return generativeModel.GenerateContentStream(ctx, textPrompt)
}

