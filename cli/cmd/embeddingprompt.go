package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/tzapio/tzap/cli/action"
	"github.com/tzapio/tzap/cli/cmd/cliworkflows"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"

	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/util"
	"github.com/tzapio/tzap/pkg/util/stdin"
)

var inspirationFiles []string
var embedsCountFlag int
var nCountFlag int
var promptFile string
var disableIndex bool
var searchQuery string

func init() {
	RootCmd.AddCommand(embeddingPromptCmd)
	embeddingPromptCmd.Flags().StringSliceVarP(&inspirationFiles,
		"inspiration", "i", []string{}, "Comma-separated list of inspiration files or multiple -i flags.")
	embeddingPromptCmd.Flags().IntVarP(&embedsCountFlag, "embeds", "k", 10,
		"Number of embeddings to use for the prompt generation")
	embeddingPromptCmd.Flags().IntVarP(&nCountFlag, "searchsize", "n", 15,
		"Number of embeddings to include in the search space before filtering out the matches with inspiration files.")
	embeddingPromptCmd.Flags().BoolVarP(&disableIndex, "disableindex", "d", false,
		"For large projects disabling indexing speeds up the process.")
	embeddingPromptCmd.Flags().StringVarP(&searchQuery, "search", "s", "",
		"The search query to start the embedding prompt with. Default (<prompt>)")
	embeddingPromptCmd.Flags().StringVarP(&promptFile, "promptfile", "f", "", "Read from file instead of prompt")
	embeddingPromptCmd.Flags().StringVarP(&lib, "lib", "l", "", "BETA: select library to search.")
}

var embeddingPromptCmd = &cobra.Command{
	Aliases: []string{"p", "prompt"},
	Use:     "embeddingprompt <prompt>",
	Short:   "Generate code or document content using code-search",
	Long: `The 'embeddingprompt' command generates content based on code-searching existing files.
	This enables GPT to be able to generate code with depth. To add breadth, the user can recommend 
	needed Inspiration files like interfaces and types to enhance GPTs general understanding.
	The inspiration files should be a comma-separated list of file paths.`,
	Run: func(cmd *cobra.Command, args []string) {
		tl.EnableUICompletionLogger()
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
		if searchQuery == "" {
			if content == "" {
				if settings.ApiMode {
					panic("search query required in ApiMode")
				}
				searchQuery = stdin.GetStdinInput("Enter your task/embedding search? (also available as -s <query>): ")
			} else {
				searchQuery = content
			}
		}

		err := tzap.HandlePanic(func() {
			t := cmdutil.GetTzapFromContext(cmd.Context())
			defer t.HandleShutdown()

			cmd.Println(cmdutil.Bold("\nSearch query: "), cmdutil.Yellow(searchQuery))

			output := action.LoadAndSearchEmbeddings(t, action.LoadAndSearchEmbeddingsArgs{
				ExcludeFiles: inspirationFiles,
				SearchQuery:  searchQuery,
				K:            embedsCount,
				N:            nCount,
				DisableIndex: disableIndex,
				Yes:          settings.Yes,
			})

			t.
				ApplyWorkflow(cliworkflows.PrintInspirationFiles(inspirationFiles)).
				ApplyWorkflow(cliworkflows.PrintSearchResults(output.SearchResults)).
				MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
					searchResults := output.SearchResults
					if len(searchResults.Results) > 0 {
						t = t.AddSystemMessage(
							"The following file contents are embeddings for the user input:",
						)
						for _, result := range searchResults.Results {
							t = t.AddSystemMessage(result.Vector.Metadata.SplitPart)
						}
					}
					return t
				}).
				MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
					if content != "" {
						cmd.Println(cmdutil.Bold("--- Completion"))
						t = t.AddUserMessage(content).RequestChatCompletion()
						cmd.Println(cmdutil.Bold("\n---"))
						if settings.ApiMode {
							t = t.AsAssistantMessage()
							jsonString, err := t.GetThreadAsJSON()
							if err != nil {
								panic(err)
							}
							fmt.Print(jsonString)
							return t
						}
						t = t.AsAssistantMessage()

					}
					for {
						input := stdin.GetStdinInput("\n\nAsk follow up question (or use ctrl+c to exit): ")
						t = t.AddUserMessage(input)
						cmd.Println(cmdutil.Bold("--- Completion:"))
						t = t.RequestChatCompletion()
						cmd.Println(cmdutil.Bold("\n---"))
						t = t.AsAssistantMessage()
					}
				})
		})
		if err != nil {
			panic(err)
		}

	},
}
