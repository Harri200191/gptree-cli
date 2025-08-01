package internal

import (
	"fmt"
)

func SummarizeFiles(prompts []string, model string, apiKey string) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("LLM API key is required")
	}

	messages := []ChatMessage{
		{
			Role:    "system",
			Content: "You are a helpful assistant that summarizes codebases. You'll receive codebase chunks in parts. At the end, you'll be asked to summarize what each file does.",
		},
	}

	for _, chunk := range prompts {
		messages = append(messages, ChatMessage{
			Role:    "user",
			Content: fmt.Sprintf("Here is the next code chunk:\n\n%s", chunk),
		})
	}

	// Final prompt to generate the summary
	messages = append(messages, ChatMessage{
		Role:    "user",
		Content: "Now summarize what each file in the codebase does in 1-3 sentences per file.",
	})

	// Send to LLM using unified method
	var response string
	var sendErr error

	if isClaudeModel(model) {
		response, sendErr = sendToClaudeWithMessages(apiKey, model, messages)
	} else {
		response, sendErr = sendToGPTWithMessages(apiKey, model, messages)
	}

	if sendErr != nil {
		return "", sendErr
	}

	return response, nil
}
