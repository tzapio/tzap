package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/action"
	"github.com/tzapio/tzap/cli/cmd/cliworkflows"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/embed/embedstore"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/workflows/code/embedworkflows"
)

func init() {
	RootCmd.AddCommand(findCmd)
	findCmd.Flags().IntVarP(&embedsCountFlag, "embeds", "k", 10, "Number of embeddings to use for the search")
	findCmd.Flags().IntVarP(&nCountFlag, "ncount", "n", 20, "Number of embeddings to use for the search")
	findCmd.Flags().StringSliceVarP(&ignoreFiles, "ignore", "i", []string{}, "Files to exclude from search")
	findCmd.Flags().BoolVarP(&disableIndex, "disableindex", "d", false, "For large projects disabling indexing speeds up the process.")
	findCmd.Flags().StringVarP(&lib, "lib", "l", "", "BETA: select library to search.")
}

var findCmd = &cobra.Command{
	Aliases: []string{"f"},
	Use:     "find <query>",
	Short:   "Find for relevant file using the query",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tl.Logger.Println("Cobra CLI Search start")
		findQuery := strings.Join(args, " ")
		err := tzap.HandlePanic(func() {
			t := cmdutil.GetTzapFromContext(cmd.Context())
			defer t.HandleShutdown()
			actionArgs := action.LoadAndSearchEmbeddingsArgs{
				ExcludeFiles: []string{},
				SearchQuery:  findQuery,
				EmbedsCount:  -1,
				NCount:       -1,
				DisableIndex: disableIndex,
				Yes:          tzapCliSettings.Yes,
			}
			t.
				AddSystemMessage(action.FindChainOfThoughtPrompt()).
				ApplyWorkflow(action.LoadAndSearchEmbeddingsWorkflow(actionArgs)).
				MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
					original := t.Data["searchResults"].(types.SearchResults)
					filenames := map[string]struct{}{}
					for _, result := range original.Results {
						filenames[result.Vector.Metadata.Filename] = struct{}{}
					}
					tmp := ""
					for filename := range filenames {
						tmp += filename + "\n"
					}
					//t = t.AddSystemMessage("These are files that exist:\n" + tmp)
					shortenedSearchResults := types.SearchResults{
						Results: embedstore.TightenSearchResults(original.Results[:embedsCountFlag]),
					}

					return t.
						ApplyWorkflow(cliworkflows.PrintEmbeddings(shortenedSearchResults)).
						ApplyWorkflow(embedworkflows.EmbedWorkflow(shortenedSearchResults))
				}).
				WorkTzap(func(t *tzap.Tzap) {
					println("---")
					t = t.
						AddUserMessage("####Find files for:\n" + findQuery).
						RequestChatCompletion()
					println("\n---\n")
					println(t.Data["content"].(string))

				})

		})

		if err != nil {
			panic(err)
		}
	},
}
