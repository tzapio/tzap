package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/util"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	tutil "github.com/tzapio/tzap/pkg/util"
	"github.com/tzapio/tzap/templates/code/embed"
)

var inspirationFilesFlag string

func init() {
	rootCmd.AddCommand(embeddingPromptCmd)
	embeddingPromptCmd.Flags().StringVarP(&inspirationFilesFlag, "inspiration", "i", "", "Comma-separated list of inspiration files")
}

var embeddingPromptCmd = &cobra.Command{
	Use:   "embeddingprompt <filename> <prompt>",
	Short: "Generate code or document content using code-search",
	Long: `The 'embeddingprompt' command generates content based on code-searching existing files. This enables GPT to be able to generate code with depth. To add breadth, the user can recommend needed Inspiration files like interfaces and types to enhance GPTs general understanding.
The inspiration files should be a comma-separated list of file paths.`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		content := strings.Join(args[1:], " ")
		var inspirationFiles []string
		if inspirationFilesFlag != "" {
			inspirationFiles = strings.Split(inspirationFilesFlag, ",")
		}
		files, err := tutil.ListFilesInDir("./")
		if err != nil {
			panic(err)
		}
		files = util.GetNonExcludedFiles(files)

		tzap.NewWithConnector(
			tzapconnect.WithConfig(
				config.Configuration{
					MD5Rewrites: true,
					EnableLogs:  false,
					OpenAIModel: modelMap[settings.Model],
				})).
			WorkTzap(func(t *tzap.Tzap) {

			}).
			// Process and create embeddings for all files in the current directory
			ApplyTemplate(embed.ProcessAndEmbedFilesTzapTemplate(files)).

			// Search for embeddings in the current directory
			ApplyTemplate(embed.EmbeddingInspirationTemplate(content, inspirationFiles)).
			AddUserMessage(content).
			LoadTaskOrRequestNewTaskMD5(filename)
	},
}
