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
var disableindex bool
var searchQuery string

func init() {
	rootCmd.AddCommand(embeddingPromptCmd)
	embeddingPromptCmd.Flags().StringVarP(&inspirationFilesFlag, "inspiration", "i", "", "Comma-separated list of inspiration files")
	embeddingPromptCmd.Flags().IntVarP(&embedsCountFlag, "embeds", "k", 10, "Number of embeddings to use for the prompt generation")
	embeddingPromptCmd.Flags().IntVarP(&nCountFlag, "searchsize", "n", 15, "Number of embeddings to include in the search space before filtering out the matches with inspiration files.")
	embeddingPromptCmd.Flags().BoolVarP(&disableindex, "disableindex", "d", false, "For large projects disabling indexing speeds up the process.")
	embeddingPromptCmd.Flags().StringVarP(&searchQuery, "search", "s", "", "The search query to start the embedding prompt with. Default (<prompt>)")
}

var embeddingPromptCmd = &cobra.Command{
	Aliases: []string{"p", "prompt"},
	Use:     "embeddingprompt <prompt>",
	Short:   "Generate code or document content using code-search",
	Long: `The 'embeddingprompt' command generates content based on code-searching existing files. This enables GPT to be able to generate code with depth. To add breadth, the user can recommend needed Inspiration files like interfaces and types to enhance GPTs general understanding.
The inspiration files should be a comma-separated list of file paths.`,
	Run: func(cmd *cobra.Command, args []string) {
		embedsCount := embedsCountFlag
		nCount := nCountFlag
		if embedsCountFlag > nCountFlag {
			nCount = embedsCountFlag + 5
		}
		var content string
		if promptFile != "" {
			content = util.ReadFileP(promptFile)
		} else {
			if len(args) > 0 {
				content = strings.Join(args[0:], " ")
			}
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
				MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
					if disableindex {
						return t
					}
					cmd.Println("Checking for file changes. (use -d to disable this check)...\n")
					return t.ApplyWorkflow(embed.PrepareEmbedFilesWorkflow(files)).
						WorkTzap(func(t *tzap.Tzap) {
							uncachedEmbeddings, ok := t.Data["uncachedEmbeddings"].(types.Embeddings)
							if !ok {
								panic("Loading embeddings went wrong")
							}
							if len(uncachedEmbeddings.Vectors) > 19 {
								price := float64(len(uncachedEmbeddings.Vectors)*400) * 0.0004 / 1000
								ok := stdin.ConfirmPrompt(fmt.Sprintf("Embeddings - You are about to fetch %d embeddings. Proceed? Estimation tokens: %d. Price is: $0.0004 per 1000 tokens. Estimating %.4f USD", len(uncachedEmbeddings.Vectors), len(uncachedEmbeddings.Vectors)*400, price))
								if !ok {
									panic("commit aborted by user")
								}
							}
						}).
						ApplyWorkflow(embed.FetchOrCachedEmbeddingForFilesWorkflow()).
						ApplyWorkflow(embed.SaveAndLoadEmbeddingsToDB())
				}).
				WorkTzap(func(t *tzap.Tzap) {
					if len(inspirationFiles) == 0 {
						cmd.Println(bold("Inspiration files: None (use --inspiration to add more)"))
						return
					}
					cmd.Println(bold("\nInspiration files:"))
					for _, inspirationFile := range inspirationFiles {
						inspirationFile = strings.TrimSpace(inspirationFile)
						tokens, err := t.CountTokens(util.ReadFileP(inspirationFile))
						if err != nil {
							panic(err)
						}
						cmd.Println(fmt.Sprintf("\tt:%d\t%s", tokens, cyan(inspirationFile)))

					}
				}).
				MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
					if searchQuery == "" {
						if content == "" {
							searchQuery = stdin.GetStdinInput("Enter your task/embedding search? (also available as -s <query>): ")
						} else {
							searchQuery = content
						}
					}
					cmd.Println(bold("\nSearch query: "), yellow(searchQuery))
					return t.ApplyWorkflow(embed.EmbeddingInspirationWorkflow(searchQuery, inspirationFiles, embedsCount, nCount))
				}).
				MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
					searchResults := t.Data["searchResults"].(types.SearchResults)

					t = t.AddSystemMessage(
						"The following file contents are embeddings for the user input:",
					)
					cmd.Println(bold("\nSearch result embeddings:"))

					for _, result := range searchResults.Results {
						t = t.AddSystemMessage(result.Metadata["splitPart"])
						tokens, err := t.CountTokens(result.Metadata["splitPart"])
						if err != nil {
							panic(err)
						}
						cmd.Println(fmt.Sprintf("\tt:%d\t%s", tokens, cyan(cmdutil.FormatVectorToClickable(result))))
					}
					time.Sleep(1 * time.Second)
					cmd.Println()

					return t
				}).
				MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
					if content != "" {
						t = t.AddUserMessage(content).RequestChatCompletion().AsAssistantMessage()
					}
					for {
						input := stdin.GetStdinInput("\n\n(We are slowly deprecating out.file, it's still required as an argument for compatability, but nothing will be written to file. Instead you can now follow up chat with code).\n\n Ask follow up question (or use ctrl+c to exit): ")
						t = t.AddUserMessage(input)
						t = t.RequestChatCompletion().AsAssistantMessage()
					}
				})
		})
		if err != nil {
			panic(err)
		}

	},
}
