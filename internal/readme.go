package internal

import (
	"fmt"
)

/// GenerateReadme generates a README.md by sending prompt chunks incrementally to the LLM,
// preserving conversation context across messages.
func GenerateReadme(promptChunks []string, model string, apiKey string) (string, error) {
	// Build up message history
	var messages []ChatMessage

	// System prompt for context
	messages = append(messages, ChatMessage{
		Role:    "system",
		Content: "You are a code assistant helping generate README.md files from full codebases.",
	})

	// Add each chunk as a user message
	for i, chunk := range promptChunks {
		userMsg := fmt.Sprintf("Code chunk %d of the project:\n\n%s", i+1, chunk)
		messages = append(messages, ChatMessage{
			Role:    "user",
			Content: userMsg,
		})
	}

	// Final instruction prompt
	messages = append(messages, ChatMessage{
		Role:    "user",
		Content: "Now that you have received the complete codebase, generate a professional README.md file. Output only the markdown code, without extra explanation or commentary.",
	})

	// Unified LLM call with full context
	readme, err := sendToLLMWithMessages(apiKey, model, messages)
	if err != nil {
		return "", err
	}

	return readme, nil
}
