package internal

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"

    "github.com/joho/godotenv"
)

type ChatRequest struct {
    Model    string       `json:"model"`
    Messages []ChatMessage `json:"messages"`
}

type ChatMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type ChatResponse struct {
    Choices []struct {
        Message ChatMessage `json:"message"`
    } `json:"choices"`
}

func GenerateReadmeFromSummary(summary, model, apiKey string) (string, error) {
	prompt := fmt.Sprintf(`Using the following file summaries, write a professional-level README.md with sections like Description, Features, Folder Structure, and How to Use (in markdown format):

		%s

		Only return the markdown content of the README.`, summary)

    return sendToLLM(apiKey, prompt, model)
}


func SummarizeFiles(root string, ignore []string, model string) (string, error) {
    godotenv.Load()

    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        return "", fmt.Errorf("OPENAI_API_KEY not set in environment")
    }

    files, err := WalkDir(root, ignore)
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

func sendToLLM(apiKey, prompt, model string) (string, error) {
    if strings.HasPrefix(model, "claude") {
        return sendToClaude(apiKey, prompt, model)
    }
    return sendToGPT(apiKey, prompt, model)
}

func sendToClaude(apiKey, prompt, model string) (string, error) {
    payload := map[string]interface{}{
        "model": model,
        "messages": []map[string]string{
            {"role": "user", "content": prompt},
        },
        "max_tokens": 1024,
    }

    body, _ := json.Marshal(payload)

    req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(body))
    if err != nil {
        return "", err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-api-key", apiKey)
    req.Header.Set("anthropic-version", "2023-06-01")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    respBody, _ := io.ReadAll(resp.Body)

    var parsed map[string]interface{}
    if err := json.Unmarshal(respBody, &parsed); err != nil {
        return "", err
    }

    if content, ok := parsed["content"].([]interface{}); ok && len(content) > 0 {
        part := content[0].(map[string]interface{})
        return part["text"].(string), nil
    }

    return "", fmt.Errorf("Claude: could not parse response: %s", string(respBody))
}


func sendToGPT(apiKey string, prompt string, model string) (string, error) {
    request := ChatRequest{
        Model: model,
        Messages: []ChatMessage{
            {Role: "system", Content: "You are a code assistant."},
            {Role: "user", Content: prompt},
        },
    }

    requestData, err := json.Marshal(request)
    if err != nil {
        return "", err
    }

    req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestData))
    if err != nil {
        return "", err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+apiKey)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    var result ChatResponse
    if err := json.Unmarshal(body, &result); err != nil {
        return "", err
    }

    if len(result.Choices) == 0 {
        return "", fmt.Errorf("no choices returned from OpenAI")
    }

    return result.Choices[0].Message.Content, nil
}
