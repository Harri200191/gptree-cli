# gptree-cli

GPT-Friendly Codebase Exporter & Summarizer for LLM Prompt Building

## Overview

`gptree-cli` is a command-line tool designed to export and summarize codebases for use with large language models (LLMs) like ChatGPT or Claude. It generates structured outputs, including directory trees and file summaries, formatted for easy integration into LLM prompts. The tool supports token limit management, Markdown output, and professional README generation.

## Features

- **Export Directory Tree**: Visualize project structure with customizable ignore paths (e.g., `.git`, `node_modules`).
- **File Content Summarization**: Summarize files using GPT or Claude models.
- **Markdown Output**: Format output in Markdown for compatibility with ChatGPT.
- **Auto-Generated README**: Create professional `README.md` files based on codebase summaries.
- **Token Management**: Split output into chunks based on specified token limits.
- **Token Estimation**: Estimate token counts for generated prompts.
- **Build with `make`**: Streamlined build process using `make` commands.

## Installation

### 1. Direct Usage (Linux)
```bash
wget https://github.com/Harri200191/gptree-cli/releases/download/<version>/gptree_<architecture>.deb
sudo dpkg -i gptree_<architecture>.deb
gtree --help
```
Note: Choose architecture from: arm, arm64, amd64. Or just download from releases!

Note: Choose release version from those available e.g. v1.0, v1.1, ...

### 2. Clone and Build
```bash
git clone https://github.com/yourname/gptree-cli
cd gptree-cli
make all
```

### 3. Windows Setup
Install `make` via Chocolatey:
```powershell
choco install make
```
Alternatively, use PowerShell to build:
```powershell
powershell -ExecutionPolicy Bypass -File build.ps1
```

## Usage

Run the tool with:
```bash
gptree [path] [flags]
```

### Example
```bash
gptree . --tree -i .git,.build,debuild

gptree . -i .build,.git,debuild --summarize -o temp.txt --llm-key <your-key> --model claude-3-haiku

gptree . -i .build,.git,debuild --readme -o temp.txt --llm-key <your-key> --model claude-3-haiku
```

### Command Flags
| Flag | Description |
|------|-------------|
| `--ignore-dirs` | Comma-separated list of directories to skip |
| `--ignore-files` | Comma-separated list of files to skip |
| `--out` | Output to a specific file |
| `--summarize` | Summarize all files using GPT or Claude |
| `--readme` | Generate a `README.md` from summaries |
| `--model` | Specify LLM model (e.g., `gpt-4`, `claude-3-sonnet`) |
| `--tree` | Display directory tree with ■ icons |
| `--llm-key` | API keys for LLM's (openAI or Anthropic)|

## Make Commands
| Command | Purpose |
|---------|---------|
| `make all` | Build the CLI to `./build/` |
| `make clean` | Delete the build directory |  
