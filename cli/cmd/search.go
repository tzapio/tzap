package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cliworkflows"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/workflows/code/embed"
)

var printEmbedding bool
var ignoreFiles []string

func init() {
	RootCmd.AddCommand(searchCmd)
	searchCmd.Flags().IntVarP(&embedsCountFlag, "embeds", "k", 20, "Number of embeddings to use for the search")
	searchCmd.Flags().IntVarP(&nCountFlag, "ncount", "n", 30, "Number of embeddings to use for the search")
	searchCmd.Flags().BoolVarP(&printEmbedding, "printembedding", "p", false, "Output the embeddings themeselves")
	searchCmd.Flags().StringSliceVarP(&ignoreFiles, "ignore", "i", []string{}, "Files to exclude from search")
	searchCmd.Flags().BoolVarP(&disableIndex, "disableindex", "d", false,
		"For large projects disabling indexing speeds up the process.")
}

var searchCmd = &cobra.Command{
	Aliases: []string{"s"},
	Use:     "search <query>",
	Short:   "Search for relevant embeddings using the query",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := strings.Join(args, " ")

		fileInDirEvaluator, err := cmdutil.NewFileInDirEvaluator()
		if err != nil {
			panic(err)
		}
		files, err := fileInDirEvaluator.WalkDir("./")
		if err != nil {
			panic(err)
		}
		err = tzap.HandlePanic(func() {
			t := cmdutil.GetTzapFromContext(cmd.Context())
			defer t.HandleShutdown()
			t = t.ApplyWorkflow(cliworkflows.LoadAndFetchEmbeddings(files, disableIndex, settings.Yes))
			searchResults := t.ApplyWorkflow(embed.SearchFilesWorkflow(query, nil, embedsCountFlag, nCountFlag)).
				Data["searchResults"].(types.SearchResults)

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
