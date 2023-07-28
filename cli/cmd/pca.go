package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapaction/action"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
	"github.com/tzapio/tzap/pkg/tzapaction/cliworkflows"
)

func init() {
	RootCmd.AddCommand(pcaCMD)
	pcaCMD.Flags().BoolVarP(&tzapCliSettings.DisableIndex, "disableindex", "d", false, "For large projects disabling indexing speeds up the process.")
	pcaCMD.Flags().StringVarP(&lib, "lib", "l", "", "BETA: select library to search.")
}

var pcaCMD = &cobra.Command{
	Use:    "pca",
	Short:  "Search for relevant embeddings using the query",
	Hidden: true,
	Args:   cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tl.Logger.Println("Cobra CLI Search start")
		searchQuery := args[0]
		err := tzap.HandlePanic(func() {
			t := cmdutil.GetTzapFromContext(cmd.Context())
			defer t.HandleShutdown()

			output, err := action.Search(t, &actionpb.SearchRequest{SearchArgs: &actionpb.SearchArgs{
				ExcludeFiles: []string{},
				SearchQuery:  searchQuery,
				EmbedsCount:  -1,
			}})
			if err != nil {
				panic(err)
			}
			if tzapCliSettings.ApiMode {
				byte, err := json.MarshalIndent(output, "", "  ")
				if err != nil {
					panic(err)
				}
				embeddingJson := string(byte)
				fmt.Println(embeddingJson)
			} else {
				t.ApplyWorkflow(cliworkflows.PrintEmbeddings(output.Embeddings))
			}
		})

		if err != nil {
			panic(err)
		}
	},
}
