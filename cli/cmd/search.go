package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/workflows/code/embed"
)

var printEmbedding bool

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().IntVarP(&embedsCountFlag, "embeds", "k", 10, "Number of embeddings to use for the search")
	searchCmd.Flags().BoolVarP(&printEmbedding, "printembedding", "p", false, "Output the embeddings themeselves")
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for relevant embeddings using the query",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := strings.Join(args, " ")

		err := tzap.HandlePanic(func() {
			t := cmdutil.GetTzapFromContext(cmd.Context())
			defer t.HandleShutdown()

			searchResults := t.ApplyWorkflow(embed.SearchFilesWorkflow(query, nil, embedsCountFlag, nCountFlag)).
				Data["searchResults"].(types.SearchResults)

			println("\nSearch result embeddings:")
			for _, result := range searchResults.Results {
				tokens, err := t.CountTokens(result.Metadata["splitPart"])
				if err != nil {
					panic("could not count tokens in search result")
				}
				fmt.Printf("\tt:%d\t%s", tokens, cyan(cmdutil.FormatVectorToClickable(result)))
				if printEmbedding {
					fmt.Printf("\n\n\t%s", result.Metadata["splitPart"])
				}
				println()

			}
			println()
		})

		if err != nil {
			panic(err)
		}
	},
}
