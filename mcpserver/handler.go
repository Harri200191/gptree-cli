package mcpserver

import (
	"encoding/json"
	"net/http"

	"github.com/Harri200191/gptree-cli/internal"
)

type Request struct {
	Path        string   `json:"path"`
	IgnoreDirs  []string `json:"ignore_dirs,omitempty"`
	IgnoreFiles []string `json:"ignore_files,omitempty"`
	OutFile     string   `json:"out,omitempty"`
	Model       string   `json:"model,omitempty"`
	LLMKey      string   `json:"llm_key,omitempty"`
}

func HandleTree(w http.ResponseWriter, r *http.Request) {
	var req Request
	json.NewDecoder(r.Body).Decode(&req)

	internal.ShowTree(req.Path, req.IgnoreDirs) 
}

func HandlePrompt(w http.ResponseWriter, r *http.Request) {
	var req Request
	json.NewDecoder(r.Body).Decode(&req)

	chunks, err := internal.BuildPrompt(req.Path, req.IgnoreDirs, 4096, true, req.IgnoreFiles)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, map[string]interface{}{"prompt_chunks": chunks})
}

func HandleSummarize(w http.ResponseWriter, r *http.Request) {
	var req Request
	json.NewDecoder(r.Body).Decode(&req)

	chunks, err := internal.BuildPrompt(req.Path, req.IgnoreDirs, 4096, true, req.IgnoreFiles)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	summary, err := internal.SummarizeFiles(chunks, req.Model, req.LLMKey)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	writeJSON(w, map[string]string{"summary": summary})
}

func HandleReadme(w http.ResponseWriter, r *http.Request) {
	var req Request
	json.NewDecoder(r.Body).Decode(&req)

	chunks, err := internal.BuildPrompt(req.Path, req.IgnoreDirs, 4096, true, req.IgnoreFiles)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	readme, err := internal.GenerateReadme(chunks, req.Model, req.LLMKey)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	writeJSON(w, map[string]string{"readme": readme})
} 