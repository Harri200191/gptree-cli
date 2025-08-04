package main

import (
	"net/http"
	"os"

	"github.com/Harri200191/gptree-cli/mcpserver"
)

func main() {
	http.HandleFunc("/tree", mcpserver.HandleTree)
	http.HandleFunc("/prompt", mcpserver.HandlePrompt)
	http.HandleFunc("/summarize", mcpserver.HandleSummarize)
	http.HandleFunc("/readme", mcpserver.HandleReadme)

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	println("🟢 GPTree MCP Server listening on port " + port)
	http.ListenAndServe("0.0.0.0:"+port, nil)
}
