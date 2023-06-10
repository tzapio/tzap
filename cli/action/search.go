package action

import (
	"github.com/tzapio/tzap/cli/cmd/cliworkflows"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/util/singlewait"
	"github.com/tzapio/tzap/workflows/code/embedworkflows"
)

type LoadAndSearchEmbeddingsArgs struct {
	ExcludeFiles []string `json:"exclude_files"`
	SearchQuery  string   `json:"search_query"`
	K            int      `json:"k"`
	N            int      `json:"n"`
	DisableIndex bool     `json:"disable_index"`
	Yes          bool     `json:"yes"`
}

type LoadAndSearchEmbeddingsOutput struct {
	SearchResults types.SearchResults `json:"searchResults"`
	QueryResult   types.QueryRequest  `json:"queryResult"`
}

func LoadAndSearchEmbeddings(t *tzap.Tzap, args LoadAndSearchEmbeddingsArgs) *LoadAndSearchEmbeddingsOutput {
	resultT := t.
		ApplyWorkflow(LoadAndSearchEmbeddingsWorkflow(args))
	searchResult := resultT.Data["searchResults"].(types.SearchResults)
	queryResult := resultT.Data["queryResult"].(types.QueryRequest)
	return &LoadAndSearchEmbeddingsOutput{
		SearchResults: searchResult,
		QueryResult:   queryResult,
	}
}

func LoadAndSearchEmbeddingsWorkflow(args LoadAndSearchEmbeddingsArgs) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "loadAndSearchEmbeddings",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			queryWait := singlewait.New(func() types.QueryRequest {
				tl.Logger.Println("loadAndSearchEmbeddings: Getting query")
				query, err := embed.NewQuery(t, args.SearchQuery)
				if err != nil {
					panic(err)
				}
				tl.Logger.Println("loadAndSearchEmbeddings: Query received")
				return query
			})

			return t.
				ApplyWorkflow(cliworkflows.IndexFilesAndEmbeddings(args.DisableIndex, args.Yes)).
				ApplyWorkflow(embedworkflows.SearchFilesWorkflow(queryWait.GetData(), args.ExcludeFiles, args.K, args.N))
		},
	}
}
