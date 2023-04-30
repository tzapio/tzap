package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	"github.com/tzapio/tzap/templates/code/files"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Tzap code!",
	Long:  ``,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			absoluteFilePath, _ := filepath.Abs(filename)
			fmt.Println("Could not find file:", absoluteFilePath, " err:", err)
			os.Exit(1)
		}

		content := strings.Join(args[1:], " ")
		tzap.
			NewWithConnector(
				tzapconnect.WithConfig(
					config.Configuration{
						SupressLogs: true,
						MD5Rewrites: true,
						OpenAIModel: openai.GPT4,
					})).
			ApplyTemplate(files.InspirationTemplate(
				[]string{
					"pkg/types/interfaces.go",
					"pkg/types/structs.go",
					"pkg/tzap/templates.go",
					//"pkg/tzap/splitter/splitter.go",
					"pkg/tzap/file.go",
					//"pkg/tzap/fetch-chat.go",
					"pkg/tzap/tzap.go",
					//"examples/githubdoc/main.go",
					"examples/refactoring/main.go",
					"templates/code/translate/translatecodefromto.go",
				},
			)).
			AddUserMessage(content).
			LoadTaskOrRequestNewTask(filename)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
