package internal

import (
    "fmt"
    "os"
    "strings"
)


func SummarizeFiles(root string, ignoreDirs []string, ignoreFiles []string, model string, apiKey string) (string, error) {
    if apiKey == "" {
        return "", fmt.Errorf("OPENAI_API_KEY not set in environment")
    }

    files, err := WalkDir(root, ignoreDirs, ignoreFiles)
    if err != nil {
        return "", err
    }

    var summaries strings.Builder
    summaries.WriteString("## File Summaries\n\n")

    for _, file := range files {
        contentBytes, err := os.ReadFile(file)
        if err != nil {
            continue
        }
        content := string(contentBytes)
        prompt := fmt.Sprintf("Summarize what this file does in 1-3 sentences:\n\n%s", content)

        summary, err := sendToGPT(apiKey, prompt, model)
        if err != nil {
            summary = fmt.Sprintf("❌ Error summarizing %s: %v\n", file, err)
        }

        summaries.WriteString(fmt.Sprintf("### %s\n%s\n\n", file, summary))
    }

    return summaries.String(), nil
}
