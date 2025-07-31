package cmd

import (
	"fmt"
	"os"

	"github.com/Harri200191/gptree-cli/internal"
	"github.com/spf13/cobra"
)

var (
	ignoreDirs     []string
	outputFile     string
	generateReadme bool
	summarize      bool
	model          string
	showTree       bool
	ignoreFiles []string
)

var rootCmd = &cobra.Command{
	Use:   "gptree",
	Short: "Export a directory tree and file contents as a GPT prompt",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify a directory to analyze.")
			return
		}

		if showTree {
			err := internal.ShowTree(args[0], ignoreDirs)
			if err != nil {
				fmt.Println("Error generating tree:", err)
			}
			return
		}

		prompt, err := internal.BuildPrompt(args[0], ignoreDirs, 4096, true, ignoreFiles)
		fmt.Println(prompt)
		return

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if outputFile != "" {
			internal.WriteToFile(outputFile, prompt)
		} else {
			fmt.Println(prompt)
		}

		if generateReadme {
			summary, err := internal.SummarizeFiles(args[0], ignoreDirs, ignoreFiles, model)
			if err != nil {
				fmt.Println("Error during summarization for README:", err)
				return
			}

			readme, err := internal.GenerateReadmeFromSummary(summary, model, os.Getenv("OPENAI_API_KEY"))
			if err != nil {
				fmt.Println("Error generating README:", err)
				return
			}

			err = internal.WriteToFile("README.md", readme)
			if err != nil {
				fmt.Println("Failed to write README.md")
				return
			}
		}

		if summarize {
			result, err := internal.SummarizeFiles(args[0], ignoreDirs, ignoreFiles, model)
			if err != nil {
				fmt.Println("Error summarizing files:", err)
				return
			}

			err = internal.WriteToFile("summaries.txt", result)
			if err != nil {
				fmt.Println("Failed to write summaries.txt")
				return
			}

			fmt.Println("✅ summaries.txt written with GPT-generated descriptions.")
			return
		}

	},
}

func Execute() {
	rootCmd.PersistentFlags().StringSliceVarP(&ignoreDirs, "ignore-dirs", "i", []string{}, "Directories to ignore")
	rootCmd.PersistentFlags().StringSliceVar(&ignoreFiles, "ignore-files", []string{}, "File patterns to ignore (e.g. *.env, *.csv, gptree)")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "out", "o", "", "Output file")
	rootCmd.PersistentFlags().BoolVar(&generateReadme, "readme", false, "Generate a README.md with summary of project files")
	rootCmd.PersistentFlags().BoolVar(&summarize, "summarize", false, "Use GPT API to summarize file contents")
	rootCmd.PersistentFlags().StringVar(&model, "model", "gpt-3.5-turbo", "Model to use: gpt-3.5-turbo | gpt-4 | claude-3-sonnet | claude-3-haiku")
	rootCmd.PersistentFlags().BoolVar(&showTree, "tree", false, "Only show tree structure of the directory (no content or summaries)")

	cobra.CheckErr(rootCmd.Execute())
}
