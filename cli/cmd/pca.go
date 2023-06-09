package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/action"
	"github.com/tzapio/tzap/cli/cmd/cliworkflows"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/tzap"
)

func init() {
	RootCmd.AddCommand(pcaCMD)
	pcaCMD.Flags().BoolVarP(&disableIndex, "disableindex", "d", false, "For large projects disabling indexing speeds up the process.")
	pcaCMD.Flags().StringVarP(&lib, "lib", "l", "", "BETA: select library to search.")
}

var pcaCMD = &cobra.Command{
	Use:   "pca",
	Short: "Search for relevant embeddings using the query",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tl.Logger.Println("Cobra CLI Search start")
		searchQuery := args[0]
		err := tzap.HandlePanic(func() {
			t := cmdutil.GetTzapFromContext(cmd.Context())
			defer t.HandleShutdown()

			output := action.LoadAndSearchEmbeddings(t, action.LoadAndSearchEmbeddingsArgs{
				ExcludeFiles: []string{},
				SearchQuery:  searchQuery,
				K:            -1,
				N:            -1,
				DisableIndex: disableIndex,
				Yes:          tzapCliSettings.Yes,
			})
			if tzapCliSettings.ApiMode {
				byte, err := json.MarshalIndent(output, "", "  ")
				if err != nil {
					panic(err)
				}
				embeddingJson := string(byte)
				fmt.Println(embeddingJson)
			} else {
				t.ApplyWorkflow(cliworkflows.PrintSearchResults(output.SearchResults))
			}
		})

		if err != nil {
			panic(err)
		}
	},
}
