package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/util"

	"github.com/tzapio/tzap/pkg/util/stdin"
	"github.com/tzapio/tzap/workflows/code/embed"
)

var inspirationFilesFlag string
var embedsCountFlag int
var nCountFlag int
var promptFile string

func init() {
	rootCmd.AddCommand(embeddingPromptCmd)
	embeddingPromptCmd.Flags().StringVarP(&inspirationFilesFlag, "inspiration", "i", "", "Comma-separated list of inspiration files")
	embeddingPromptCmd.Flags().IntVarP(&embedsCountFlag, "embeds", "k", 10, "Number of embeddings to use for the prompt generation")
	embeddingPromptCmd.Flags().IntVarP(&nCountFlag, "searchsize", "n", 15, "Number of embeddings to include in the search space before filtering out the matches with inspiration files.")
}

var embeddingPromptCmd = &cobra.Command{
	Aliases: []string{"p", "prompt"},
	Use:     "embeddingprompt <file> <prompt>",
	Short:   "Generate code or document content using code-search",
	Long: `The 'embeddingprompt' command generates content based on code-searching existing files. This enables GPT to be able to generate code with depth. To add breadth, the user can recommend needed Inspiration files like interfaces and types to enhance GPTs general understanding.
The inspiration files should be a comma-separated list of file paths.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		embedsCount := embedsCountFlag
		nCount := nCountFlag
		if embedsCountFlag > nCountFlag {
			nCount = embedsCountFlag + 5
		}
		var content string
		if promptFile != "" {
			content = util.ReadFileP(promptFile)
		} else {
			if len(args) < 2 {
				panic("Missing file, please provide a file.\n\nUsage: tzap embeddingprompt <file> <prompt>")
			}

			content = strings.Join(args[1:], " ")
		}

		var inspirationFiles []string
		if inspirationFilesFlag != "" {
			inspirationFiles = strings.Split(inspirationFilesFlag, ",")
		}
		files, err := util.ListFilesInDir("./")
		if err != nil {
			panic(err)
		}
		files = cmdutil.GetNonExcludedFiles(files)

		err = tzap.HandlePanic(func() {
			t := cmdutil.GetTzapFromContext(cmd.Context())
			defer t.HandleShutdown()
			t.
				// Process and create embeddings for all files in the current directory
				ApplyWorkflow(embed.PrepareEmbedFilesWorkflow(files)).
				WorkTzap(func(t *tzap.Tzap) {
					uncachedEmbeddings, ok := t.Data["uncachedEmbeddings"].(types.Embeddings)
					if !ok {

						panic("Loading embeddings went wrong")
					}
					if len(uncachedEmbeddings.Vectors) > 19 {
						ok := stdin.ConfirmPrompt(fmt.Sprintf("Embeddings - You are about to fetch %d embeddings. Proceed?", len(uncachedEmbeddings.Vectors)))
						if !ok {
							panic("commit aborted by user")
						}
					}
				}).
				ApplyWorkflow(embed.FetchOrCachedEmbeddingForFilesWorkflow()).
				ApplyWorkflow(embed.SaveAndLoadEmbeddingsToDB()).
				WorkTzap(func(t *tzap.Tzap) {
					if len(inspirationFiles) == 0 {
						println(bold("Inspiration files: None (use --inspiration to add more)"))
						return
					}
					println(bold("\nInspiration files:"))
					for _, inspirationFile := range inspirationFiles {
						inspirationFile = strings.TrimSpace(inspirationFile)
						tokens, err := t.CountTokens(util.ReadFileP(inspirationFile))
						if err != nil {
							panic(err)
						}
						println(fmt.Sprintf("\tt:%d\t%s", tokens, cyan(inspirationFile)))

					}
				}).
				WorkTzap(func(t *tzap.Tzap) {
					println(bold("\nSearch query: "), yellow(content))
				}).
				ApplyWorkflow(embed.EmbeddingInspirationWorkflow(content, inspirationFiles, embedsCount, nCount)).
				MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
					searchResults := t.Data["searchResults"].(types.SearchResults)

					t = t.AddSystemMessage(
						"The following file contents are embeddings for the user input:",
					)
					println(bold("\nSearch result embeddings:"))

					for _, result := range searchResults.Results {
						t = t.AddSystemMessage(result.Metadata["splitPart"])
						tokens, err := t.CountTokens(result.Metadata["splitPart"])
						if err != nil {
							panic(err)
						}
						println(fmt.Sprintf("\tt:%d\t%s", tokens, cyan(result.ID)))
					}
					time.Sleep(1 * time.Second)
					println()
					return t
				}).
				AddUserMessage(content).
				RequestChatCompletion().
				StoreCompletion(filePath)
		})
		if err != nil {
			println(err.Error())
			panic(err)
		}

	},
}
