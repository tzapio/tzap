package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapaction/action"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
	"github.com/tzapio/tzap/pkg/tzapaction/cliworkflows"
)

var ignoreFiles []string
var lib string

func init() {
	RootCmd.AddCommand(searchCmd)
	searchCmd.Flags().Int32VarP(&embedsCountFlag, "embeds", "k", 30, "Number of embeddings to use for the search")
	searchCmd.Flags().StringVarP(&lib, "lib", "l", "", "BETA: select library to search.")
}

var searchCmd = &cobra.Command{
	Aliases: []string{"s"},
	Use:     "search <query>",
	Short:   "Search for relevant embeddings using the query",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tl.Logger.Println("Cobra CLI Search start")
		searchQuery := strings.Join(args, " ")

		err := tzap.HandlePanic(func() {
			t := cmdutil.GetTzapFromContext(cmd.Context())
			defer t.HandleShutdown()

			output, err := action.Search(t, &actionpb.SearchRequest{SearchArgs: &actionpb.SearchArgs{
				ExcludeFiles: []string{},
				Lib:          lib,
				SearchQuery:  searchQuery,
				EmbedsCount:  embedsCountFlag,
			}})
			if err != nil {
				panic(err)
			}
			if tzapCliSettings.ApiMode {
				byte, err := json.MarshalIndent(output.Embeddings, "", "  ")
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
