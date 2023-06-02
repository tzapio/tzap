package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cliworkflows"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/embed/localdb/singlewait"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/workflows/code/embedworkflows"
)

var printEmbedding bool
var ignoreFiles []string
var zipURL string

func init() {
	RootCmd.AddCommand(searchCmd)
	searchCmd.Flags().IntVarP(&embedsCountFlag, "embeds", "k", 10, "Number of embeddings to use for the search")
	searchCmd.Flags().IntVarP(&nCountFlag, "ncount", "n", 20, "Number of embeddings to use for the search")
	searchCmd.Flags().BoolVarP(&printEmbedding, "printembedding", "p", false, "Output the embeddings themeselves")
	searchCmd.Flags().StringSliceVarP(&ignoreFiles, "ignore", "i", []string{}, "Files to exclude from search")
	searchCmd.Flags().BoolVarP(&disableIndex, "disableindex", "d", false,

		"For large projects disabling indexing speeds up the process.")
	searchCmd.Flags().StringVarP(&zipURL, "zipurl", "z", "", "URL of the zip file to search files from")

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

			queryWait := singlewait.New(func() types.QueryRequest {
				tl.Logger.Println("Query start")

				query, err := embed.GetQuery(t, searchQuery)
				if err != nil {
					panic(err)
				}
				tl.Logger.Println("Query end")
				return query
			})

			t = t.
				ApplyWorkflow(cliworkflows.IndexFilesAndEmbeddings("./", disableIndex, settings.Yes))
			searchResults := t.ApplyWorkflow(embedworkflows.SearchFilesWorkflow(queryWait.GetData(), nil, embedsCountFlag, nCountFlag)).
				Data["searchResults"].(types.SearchResults)
			tl.Logger.Println("Showing results")
			cmd.Println("\nSearch result embeddings:")
			var metadatas []map[string]string
			for _, result := range searchResults.Results {
				tokens, err := t.CountTokens(result.Metadata["splitPart"])
				if err != nil {
					panic("could not count tokens in search result")
				}
				cmd.Printf("\tt:%d\t%s", tokens, cmdutil.Cyan(cmdutil.FormatVectorToClickable(result)))
				if printEmbedding {
					cmd.Printf("\n\n\t%s", result.Metadata["splitPart"])
				}
				cmd.Println()
				if settings.ApiMode {
					metadatas = append(metadatas, result.Metadata)
				}

			}
			if settings.ApiMode {
				byte, err := json.Marshal(metadatas)
				if err != nil {
					panic(err)
				}
				embeddingJson := string(byte)
				fmt.Println(embeddingJson)
			}
			cmd.Println()
		})

		if err != nil {
			panic(err)
		}
	},
}
