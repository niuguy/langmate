package main

import (
	"fmt"
	"os"

	"github.com/niuguy/langmate/app"
	"github.com/niuguy/langmate/llm"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "langmate",
	Short: "Double copy text to translate or rephrase",
	Long:  `Double copy text to translate or rephrase`,
	Run: func(cmd *cobra.Command, args []string) {
		llmType := "p" // Default to OpenAI
		if len(args) > 0 {
			llmType = args[0]
		}
		textProcessor := llm.CreateTextProcessor(llmType)
		app.StartHook(textProcessor)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
