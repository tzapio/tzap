package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	tutil "github.com/tzapio/tzap/pkg/util"
	"github.com/tzapio/tzap/pkg/util/stdin"
	"github.com/tzapio/tzap/workflows/code/embed"
)

var inspirationFilesFlag string
var embedsCountFlag int
var nCountFlag int

func init() {
	rootCmd.AddCommand(embeddingPromptCmd)
	embeddingPromptCmd.Flags().StringVarP(&inspirationFilesFlag, "inspiration", "i", "", "Comma-separated list of inspiration files")
	embeddingPromptCmd.Flags().IntVarP(&embedsCountFlag, "embeds", "k", 10, "Number of embeddings to use for the prompt generation")
	embeddingPromptCmd.Flags().IntVarP(&nCountFlag, "searchsize", "n", 15, "Number of embeddings to include in the search space before filtering out the matches with inspiration files.")
}

var embeddingPromptCmd = &cobra.Command{
	Aliases: []string{"p"},
	Use:     "embeddingprompt <file> <prompt>",
	Short:   "Generate code or document content using code-search",
	Long: `The 'embeddingprompt' command generates content based on code-searching existing files. This enables GPT to be able to generate code with depth. To add breadth, the user can recommend needed Inspiration files like interfaces and types to enhance GPTs general understanding.
The inspiration files should be a comma-separated list of file paths.`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		embedsCount := embedsCountFlag
		nCount := nCountFlag
		if embedsCountFlag > nCountFlag {
			nCount = embedsCountFlag + 5
		}
		content := strings.Join(args[1:], " ")
		var inspirationFiles []string
		if inspirationFilesFlag != "" {
			inspirationFiles = strings.Split(inspirationFilesFlag, ",")
		}
		files, err := tutil.ListFilesInDir("./")
		if err != nil {
			panic(err)
		}
		files = cmdutil.GetNonExcludedFiles(files)
		cmdutil.GetTzapFromContext(cmd.Context()).
			// Process and create embeddings for all files in the current directory
			ApplyWorkflow(embed.PrepareEmbedFilesTzapWorkflow(files)).
			WorkTzap(func(t *tzap.Tzap) {
				uncachedEmbeddings, ok := t.Data["uncachedEmbeddings"].(types.Embeddings)
				if !ok {
					panic("Loading embeddings went wrong")
				}
				if len(uncachedEmbeddings.Vectors) > 20 {
					ok := stdin.ConfirmPrompt(fmt.Sprintf("Embeddings - You are about to fetch %d embeddings. Proceed?", len(uncachedEmbeddings.Vectors)))
					if !ok {
						panic("commit aborted by user")
					}
				}
			}).
			ApplyWorkflow(embed.FetchOrCachedEmbeddingForFilesTzapWorkflow()).
			ApplyWorkflow(embed.SaveAndLoadEmbeddingsToDB()).

			// Search for embeddings in the current directory
			ApplyWorkflow(embed.EmbeddingInspirationWorkflow(content, inspirationFiles, embedsCount, nCount)).
			AddUserMessage(content).
			RequestChatCompletion().
			StoreCompletion(filePath)
		//.ApplyWorkflow(codegeneration.GenerateCodeAndApplyWorkflow())
	},
}
