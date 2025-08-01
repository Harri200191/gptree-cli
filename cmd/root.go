package cmd

import (
	"fmt"
	"os"
	"strings"
	"github.com/Harri200191/gptree-cli/internal"
	"github.com/spf13/cobra"
)

var (
	ignoreDirs     []string
	ignoreFiles    []string
	outputFile     string
	generateReadme bool
	summarize      bool
	model          string
	showTree       bool
	llmKey         string
)

var modelAliases = map[string]string{
	"claude-3-haiku":  "claude-3-haiku-20240307",
	"claude-3-sonnet": "claude-3-sonnet-20240229",
	"claude-3-opus":   "claude-3-opus-20240229",
	"haiku":           "claude-3-haiku-20240307",
	"sonnet":          "claude-3-sonnet-20240229",
	"opus":            "claude-3-opus-20240229",
	"gpt-3.5": "gpt-3.5-turbo",
	"gpt-4":   "gpt-4",
}

func NormalizeModel(model string) string {
	model = strings.ToLower(model)
	if resolved, ok := modelAliases[model]; ok {
		return resolved
	}
	return model  
}

var rootCmd = &cobra.Command{
	Use:   "gptree [directory]",
	Short: "Export a directory tree and file contents as a GPT prompt",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("❌ Please specify a directory to analyze.")
			os.Exit(1)
		}

		// ✅ Validate incompatible combinations
		if showTree {
			if summarize || generateReadme || outputFile != "" || len(ignoreFiles) > 0 || model != "gpt-3.5-turbo" || llmKey != "" {
				fmt.Println("❌ --tree cannot be combined with --summarize, --readme, --ignore-files, --model, --out, or --llm-key")
				os.Exit(1)
			}
		} 

		// Can not use summarize and readme together
		if summarize && generateReadme {
			fmt.Println("❌ --summarize and --readme cannot be combined")
			os.Exit(1)
		}

		if (generateReadme || summarize) && llmKey == "" {
			fmt.Println("❌ You must provide an --llm-key when using --readme or --summarize")
			os.Exit(1)
		}

		// 🌲 Just show the directory tree
		if showTree {
			err := internal.ShowTree(args[0], ignoreDirs)
			if err != nil {
				fmt.Println("Error generating tree:", err)
			}
			return
		}

		// 🧠 Build prompt
		prompt, err := internal.BuildPrompt(args[0], ignoreDirs, 4096, true, ignoreFiles)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		} 

		// 📘 Generate README
		if generateReadme {
			// ensure that model is specified
			model = NormalizeModel(model)
			if model == "" {
				fmt.Println("❌ You must specify a model when using --readme")
				os.Exit(1)
			}

			readme, err := internal.GenerateReadme(prompt, model, llmKey)
			if err != nil {
				fmt.Println("Error generating README:", err)
				return
			}

			outputPath := outputFile
			if outputPath == "" {
				outputPath = "README.md"
			} else {
				if !strings.HasSuffix(outputPath, ".md") {
					fmt.Println("❌ Output file for README must have a .md extension.")
					return
				}
			}

			// Check if file exists
			if _, err := os.Stat(outputPath); err == nil {
				// File exists — confirm overwrite
				fmt.Printf("⚠️  %s already exists. Overwrite? (Y/n): ", outputPath)
				var response string
				fmt.Scanln(&response)
				response = strings.ToLower(strings.TrimSpace(response))
				if response != "y" && response != "" {
					fmt.Println("❌ Aborted. README not written.")
					return
				}
			}

			err = internal.WriteToFile(outputPath, readme)
			if err != nil {
				fmt.Printf("Failed to write %s: %v\n", outputPath, err)
				return
			}

			fmt.Printf("✅ %s written successfully.\n", outputPath)
		}


		// 📄 Summarize files
		if summarize {
			// ensure that model is specified
			model = NormalizeModel(model)
			if model == "" {
				fmt.Println("❌ You must specify a model when using --summarize")
				os.Exit(1)
			}

			result, err := internal.SummarizeFiles(prompt, model, llmKey)
			if err != nil {
				fmt.Println("Error summarizing files:", err)
				return
			}

			outputPath := outputFile
			if outputPath == "" {
				outputPath = "summaries.txt"
			} else {
				if !strings.HasSuffix(outputPath, ".txt") {
					fmt.Println("❌ Output file for summaries must have a .txt extension.")
					return
				}
			}

			// Check if file exists
			if _, err := os.Stat(outputPath); err == nil {
				// File exists — confirm overwrite
				fmt.Printf("⚠️  %s already exists. Overwrite? (Y/n): ", outputPath)
				var response string
				fmt.Scanln(&response)
				response = strings.ToLower(strings.TrimSpace(response))
				if response != "y" && response != "" {
					fmt.Println("❌ Aborted. Summaries not written.")
					return
				}
			}

			err = internal.WriteToFile(outputPath, result)
			if err != nil {
				fmt.Printf("Failed to write %s: %v\n", outputPath, err)
				return
			}

			fmt.Printf("✅ %s written with GPT-generated descriptions.\n", outputPath)
		}

	},
}

func Execute() {
	rootCmd.PersistentFlags().StringSliceVarP(&ignoreDirs, "ignore-dirs", "i", []string{}, "Directories to ignore")
	rootCmd.PersistentFlags().StringSliceVar(&ignoreFiles, "ignore-files", []string{}, "File patterns to ignore (e.g. *.env, *.csv, gptree)")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "out", "o", "", "Output file to write the prompt/summary")
	rootCmd.PersistentFlags().BoolVar(&generateReadme, "readme", false, "Generate a README.md with GPT summary")
	rootCmd.PersistentFlags().BoolVar(&summarize, "summarize", false, "Use GPT API to summarize file contents into summaries.txt")
	rootCmd.PersistentFlags().StringVar(&model, "model", "gpt-3.5-turbo", "LLM model to use: gpt-3.5-turbo | gpt-4 | claude-3-sonnet | claude-3-haiku")
	rootCmd.PersistentFlags().BoolVar(&showTree, "tree", false, "Only show ASCII directory tree (no prompt or summaries)")
	rootCmd.PersistentFlags().StringVar(&llmKey, "llm-key", "", "API key for LLMs (OpenAI or Anthropic)")

	cobra.CheckErr(rootCmd.Execute())
}
