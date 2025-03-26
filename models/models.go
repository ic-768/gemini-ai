package models

// PromptRequest represents the incoming request structure for text generation
type PromptRequest struct {
	Prompt string `json:"prompt"`
} 