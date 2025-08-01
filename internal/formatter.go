package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var tokenThreshold = 100 // soft buffer before hitting max

func BuildPrompt(root string, ignoreDirs []string, maxTokens int, chunk bool, ignoreFiles []string) ([]string, error) {
	files, err := WalkDir(root, ignoreDirs, ignoreFiles)
	if err != nil {
		return nil, err
	}

	filtered := FilterFiles(files, ignoreFiles)

	if chunk {
		return buildChunkedPrompt(filtered, maxTokens)
	}

	// Single prompt mode
	var prompt strings.Builder
	for _, file := range filtered {
		content, err := readSanitizedFile(file)
		if err != nil {
			continue
		}

		prompt.WriteString(fmt.Sprintf("\n########## %s ##########\n", file))
		prompt.WriteString(content)
		prompt.WriteString("\n")

		if EstimateTokens(prompt.String()) >= maxTokens-tokenThreshold {
			break
		}
	}

	return []string{prompt.String()}, nil
}

func buildChunkedPrompt(files []string, maxTokens int) ([]string, error) {
	var chunks []string
	var current strings.Builder
	tokenCount := 0

	for _, file := range files {
		content, err := readSanitizedFile(file)
		if err != nil {
			continue
		}

		entry := fmt.Sprintf("\n########## %s ##########\n%s\n", file, content)
		newTokens := EstimateTokens(entry)

		if tokenCount+newTokens >= maxTokens-tokenThreshold {
			chunks = append(chunks, current.String())
			current.Reset()
			tokenCount = 0
		}

		current.WriteString(entry)
		tokenCount += newTokens
	}

	if current.Len() > 0 {
		chunks = append(chunks, current.String())
	}

	return chunks, nil
}

func readSanitizedFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	content := string(data)
	// Replace tabs and newlines
	content = strings.ReplaceAll(content, "\t", "\\t")
	content = strings.ReplaceAll(content, "\n", "\\n")
	return content, nil
}

func FilterFiles(files []string, ignorePatterns []string) []string {
	var filtered []string

	for _, file := range files {
		base := filepath.Base(file)
		skip := false

		for _, pattern := range ignorePatterns {
			match, _ := filepath.Match(pattern, base)
			if match {
				skip = true
				break
			}
		}

		if !skip {
			filtered = append(filtered, file)
		}
	}

	return filtered
}
