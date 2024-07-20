package main

import (
	"fmt"
	"os"

	"github.com/niuguy/langmate/app"
	"github.com/niuguy/langmate/llm"
	"github.com/spf13/cobra"
)

var (
	model string
	lang  string
)

var rootCmd = &cobra.Command{
	Use:   "langmate",
	Short: "Double copy text to translate or rephrase",
	Long:  `Double copy text to translate or rephrase`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("model:", model)
		// fmt.Println("lang:", lang)
		textProcessor := llm.CreateTextProcessor(model)
		app.StartHook(textProcessor, lang)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&model, "model", "m", "gpt", "Specify the model to use (gpt,llama)")
	rootCmd.PersistentFlags().StringVarP(&lang, "lang", "l", "en", "Specify the target language (e.g., en, fr)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
