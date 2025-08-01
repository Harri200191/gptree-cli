package internal

import (
	"fmt"
	"strings"
)

// GenerateReadme generates a README.md by sending prompt chunks incrementally to the LLM,
// preserving conversation context across messages.
func GenerateReadme(promptChunks []string, model string, apiKey string) (string, error) {
	// Build up message history
	var messages []ChatMessage

	// Add system prompt (for GPT-compatible models)
	if !isClaudeModel(model) {
		messages = append(messages, ChatMessage{
			Role:    "system",
			Content: "You are a code assistant helping generate README.md files from codebase summaries.",
		})
	}

	// Feed each chunk to build context
	for i, chunk := range promptChunks {
		userMsg := fmt.Sprintf("Chunk %d of project codebase:\n%s", i+1, chunk)
		messages = append(messages, ChatMessage{
			Role:    "user",
			Content: userMsg,
		})

		// For Claude, send each message directly and discard response (we just need to build context)
		if isClaudeModel(model) {
			_, err := sendToClaudeWithMessages(apiKey, model, messages)
			if err != nil {
				return "", fmt.Errorf("Claude context chunk %d failed: %w", i+1, err)
			}
		}
	}

	finalPrompt := "Now that you have received the complete codebase, generate a complete, well-structured README.md file summarizing the project. Output just the .md code without any explaination or any headings etc. Just code"

	messages = append(messages, ChatMessage{
		Role:    "user",
		Content: finalPrompt,
	})

	var readme string
	var err error

	if isClaudeModel(model) {
		readme, err = sendToClaudeWithMessages(apiKey, model, messages)
	} else {
		readme, err = sendToGPTWithMessages(apiKey, model, messages)
	}

	if err != nil {
		return "", err
	}
	return readme, nil
}

func isClaudeModel(model string) bool {
	return strings.HasPrefix(model, "claude")
}
