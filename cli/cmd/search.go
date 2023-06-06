package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/action"
	"github.com/tzapio/tzap/cli/cmd/cliworkflows"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

var ignoreFiles []string
var lib string

func init() {
	RootCmd.AddCommand(searchCmd)
	searchCmd.Flags().IntVarP(&embedsCountFlag, "embeds", "k", 10, "Number of embeddings to use for the search")
	searchCmd.Flags().IntVarP(&nCountFlag, "ncount", "n", 20, "Number of embeddings to use for the search")
	searchCmd.Flags().StringSliceVarP(&ignoreFiles, "ignore", "i", []string{}, "Files to exclude from search")
	searchCmd.Flags().BoolVarP(&disableIndex, "disableindex", "d", false, "For large projects disabling indexing speeds up the process.")
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
			var libname types.ProjectName = types.LOCALPROJECTNAME
			if lib != "" {
				libname = types.ProjectName(lib)
			}
			output := action.LoadAndSearchEmbeddings(t, action.LoadAndSearchEmbeddingsArgs{
				ExcludeFiles: []string{},
				SearchQuery:  searchQuery,
				Lib:          libname,
				K:            embedsCountFlag,
				N:            nCountFlag,
				DisableIndex: disableIndex,
				Yes:          settings.Yes,
			})
			if settings.ApiMode {
				var metadatas []types.Metadata
				for _, result := range output.SearchResults.Results {
					metadatas = append(metadatas, result.Vector.Metadata)
				}
				byte, err := json.MarshalIndent(metadatas, "", "  ")
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
