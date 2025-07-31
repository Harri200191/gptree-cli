package internal

import (
	"fmt"
	"os"
	"strings"
)

func BuildPrompt(root string, ignore []string, maxTokens int, chunk bool) (string, error) {
	files, err := WalkDir(root, ignore)
	if err != nil {
		return "", err
	}

	if chunk {
		return buildChunkedPrompt(root, files, maxTokens)
	}

	// Single prompt version
	var prompt strings.Builder
	for _, file := range files {
		contentBytes, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		content := string(contentBytes)

		prompt.WriteString(fmt.Sprintf("\n########## %s ##########\n", file))
		prompt.WriteString(content)
		prompt.WriteString("\n")


		if EstimateTokens(prompt.String()) > maxTokens {
			break
		}
	}

	return prompt.String(), nil
}


func GenerateReadme(root string, ignore []string) (string, error) {
	files, err := WalkDir(root, ignore)
	if err != nil {
		return "", err
	}

	var readme strings.Builder
	readme.WriteString("# Project Overview\n\n")
	readme.WriteString("This project contains the following key files:\n\n")

	for _, file := range files {
		readme.WriteString(fmt.Sprintf("- `%s`: _Brief description here_\n", file))
	}

	readme.WriteString("\n> NOTE: Replace file descriptions with meaningful explanations.")

	return readme.String(), nil
}

func buildChunkedPrompt(root string, files []string, maxTokens int) (string, error) {
	var current strings.Builder
	chunkIndex := 1
	tokenCount := 0

	for _, file := range files {
		contentBytes, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		content := string(contentBytes)

		var entry string
		entry = fmt.Sprintf("\n########## %s ##########\n%s\n", file, content)


		newTokenCount := EstimateTokens(entry)
		if tokenCount+newTokenCount > maxTokens {
			// Write the current chunk to file
			filename := fmt.Sprintf("prompt_part_%d.txt", chunkIndex)
			err := WriteToFile(filename, current.String())
			if err != nil {
				return "", err
			}

			fmt.Println("Written:", filename)

			// Reset for next chunk
			chunkIndex++
			current.Reset()
			tokenCount = 0
		}

		current.WriteString(entry)
		tokenCount += newTokenCount
	}

	// Write the final remaining chunk
	if current.Len() > 0 {
		filename := fmt.Sprintf("prompt_part_%d.txt", chunkIndex)
		err := WriteToFile(filename, current.String())
		if err != nil {
			return "", err
		}
		fmt.Println("Written:", filename)
	}

	return fmt.Sprintf("%d prompt chunks written.", chunkIndex), nil
}
