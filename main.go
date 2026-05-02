package main

import (
	"fmt"
	"os"

	"github.com/niuguy/langmate/app"
	"github.com/spf13/cobra"
)

var (
	model   string
	lang    string
	preview bool
)

var rootCmd = &cobra.Command{
	Use:   "langmate",
	Short: "Rephrase selected text with Cmd+Ctrl+R",
	Long:  `LangMate runs as a background daemon that listens for Cmd+Ctrl+R hotkey to rephrase selected text in place.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := app.LoadDaemonConfig()
		if cmd.Flags().Changed("model") {
			config.Model = model
		}
		if cmd.Flags().Changed("lang") {
			config.Lang = lang
		}
		if cmd.Flags().Changed("preview") {
			config.Preview = preview
		}
		app.StartDaemon(config)
	},
}

func init() {
	defaultConfig := app.DefaultDaemonConfig()
	rootCmd.PersistentFlags().StringVarP(&model, "model", "m", defaultConfig.Model, "Specify the model to use (gpt,llama)")
	rootCmd.PersistentFlags().StringVarP(&lang, "lang", "l", defaultConfig.Lang, "Specify the target language (e.g., en, fr)")
	rootCmd.PersistentFlags().BoolVar(&preview, "preview", false, "Preview and edit rephrased text before replacing the selection")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
