package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	"github.com/tzapio/tzap/templates/code/files"
)

var inspirationFilesFlag string

func init() {
	rootCmd.AddCommand(embeddingPromptCmd)
	embeddingPromptCmd.Flags().StringVarP(&inspirationFilesFlag, "inspiration", "i", "", "Comma-separated list of inspiration files")
}

var embeddingPromptCmd = &cobra.Command{
	Use:   "embeddingprompt <filename> <prompt>",
	Short: "Generate code or document content using embedding inspiration template",
	Long: `The 'embeddingprompt' command generates content based on the prompt and inspiration files provided.
The inspiration files should be a comma-separated list of file paths.`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		content := strings.Join(args[1:], " ")
		var inspirationFiles []string
		if inspirationFilesFlag != "" {
			inspirationFiles = strings.Split(inspirationFilesFlag, ",")
		}
		tzap.NewWithConnector(
			tzapconnect.WithConfig(
				config.Configuration{
					MD5Rewrites: true,
					EnableLogs:  false,
					OpenAIModel: modelMap[settings.Model],
				})).
			WorkTzap(func(t *tzap.Tzap) {

			}).
			ApplyTemplate(files.ProcessAndEmbedFilesTzapTemplate("./")).
			ApplyTemplate(files.EmbeddingInspirationTemplate(content, inspirationFiles)).
			AddUserMessage(content).
			LoadTaskOrRequestNewTaskMD5(filename)
	},
}
