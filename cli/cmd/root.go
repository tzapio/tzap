package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/pkg/types/openai"
)

var settings struct {
	Model string
}
var rootCmd = &cobra.Command{
	Use:   "tzap",
	Short: "Tzap Cli!",
	Long:  `tbd`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var modelMap map[string]string = map[string]string{"gpt35": openai.GPT3Dot5Turbo, "gpt4": openai.GPT4}

func init() {
	rootCmd.PersistentFlags().StringVarP(&settings.Model, "model", "m", "gpt35", "Define what openai model to use. (Available gpt35 gpt4).")
}
