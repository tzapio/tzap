package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapaction/action"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
)

func init() {
	RootCmd.AddCommand(findCmd)
	findCmd.Flags().Int32VarP(&embedsCountFlag, "embeds", "k", 10, "Number of embeddings to use for the search")
	findCmd.Flags().StringSliceVarP(&ignoreFiles, "ignore", "i", []string{}, "Files to exclude from search")
	findCmd.Flags().BoolVarP(&tzapCliSettings.DisableIndex, "disableindex", "d", false, "For large projects disabling indexing speeds up the process.")
	findCmd.Flags().StringVarP(&lib, "lib", "l", "", "BETA: select library to search.")
}

var findCmd = &cobra.Command{
	Aliases: []string{"f"},
	Use:     "find <query>",
	Short:   "Find for relevant file using the query",
	Hidden:  true,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tl.Logger.Println("Cobra CLI Search start")
		findQuery := strings.Join(args, " ")
		err := tzap.HandlePanic(func() {
			t := cmdutil.GetTzapFromContext(cmd.Context())
			defer t.HandleShutdown()
			actionArgs := &actionpb.SearchArgs{
				ExcludeFiles: []string{},
				SearchQuery:  findQuery,
				EmbedsCount:  embedsCountFlag,
				Lib:          lib,
			}
			t.
				AddSystemMessage(action.FindChainOfThoughtPrompt()).
				ApplyWorkflow(action.SearchWorkflow(actionArgs)).
				WorkTzap(func(t *tzap.Tzap) {
					println("---")
					t = t.
						AddUserMessage("####Find files for:\n" + findQuery).
						RequestChatCompletion()
					println("\n---\n")
					println(t.Data["content"].(types.CompletionMessage).Content)

				})

		})

		if err != nil {
			panic(err)
		}
	},
}
