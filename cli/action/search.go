package action

import (
	"github.com/tzapio/tzap/cli/actionpb"
	"github.com/tzapio/tzap/cli/cmd/cliworkflows"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/embed/embedstore"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/util/singlewait"
	"github.com/tzapio/tzap/workflows/code/embedworkflows"
	"github.com/tzapio/tzap/workflows/code/fileworkflows"
)

type LoadAndSearchEmbeddingsArgs struct {
	ExcludeFiles []string `json:"exclude_files"`
	SearchQuery  string   `json:"search_query"`
	EmbedsCount  int      `json:"k"`
	NCount       int      `json:"n"`
	DisableIndex bool     `json:"disable_index"`
	Yes          bool     `json:"yes"`
}

type LoadAndSearchEmbeddingsOutput struct {
	SearchResults types.SearchResults `json:"searchResults"`
	QueryResult   types.QueryRequest  `json:"queryResult"`
}

func LoadAndSearchEmbeddings(t *tzap.Tzap, args *actionpb.SearchArgs) *LoadAndSearchEmbeddingsOutput {
	resultT := t.
		ApplyWorkflow(LoadAndSearchEmbeddingsWorkflow(args))
	searchResult := resultT.Data["searchResults"].(types.SearchResults)
	queryResult := resultT.Data["queryResult"].(types.QueryRequest)
	return &LoadAndSearchEmbeddingsOutput{
		SearchResults: searchResult,
		QueryResult:   queryResult,
	}
}

func LoadAndSearchEmbeddingsWorkflow(args *actionpb.SearchArgs) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
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
				ApplyWorkflow(embedworkflows.SearchFilesWorkflow(queryWait.GetData(), args.ExcludeFiles, int(args.EmbedsCount), int(args.NCount)))
		},
	}
}

// PromptWorkflow defines a workflow for generating code based on code-searching existing files and user input.
// This workflow creates a chat interface for the user to prompt for their desired output.
// The inspiration files parameter allows the user to recommend needed inspiration files to enhance GPT's general understanding.
// The search results are generated using a configurable number of embeddings and ranked by relevance.
// The user is then prompted to select a single search result. Once selected, the selected search result is used to generate the final desired output. Finally, the workflow is terminated when the completion of the user's message thread is reached.
func SearchWorkflow(promptWorkflowArgs PromptWorkflowArgs) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			loadAndSearchEmbeddingsArgs := &actionpb.SearchArgs{
				ExcludeFiles: promptWorkflowArgs.InspirationFiles,
				SearchQuery:  promptWorkflowArgs.SearchQuery,
				EmbedsCount:  promptWorkflowArgs.EmbedsCount,
				NCount:       promptWorkflowArgs.NCount,
				DisableIndex: promptWorkflowArgs.DisableIndex,
				Yes:          promptWorkflowArgs.Yes,
			}

			return t.
				// Include all full length inspiration files.
				ApplyWorkflow(cliworkflows.PrintInspirationFiles(promptWorkflowArgs.InspirationFiles)).
				ApplyWorkflow(fileworkflows.InspirationWorkflow(promptWorkflowArgs.InspirationFiles)).
				// Find all embeddings search them, returning the top results.
				ApplyWorkflow(LoadAndSearchEmbeddingsWorkflow(loadAndSearchEmbeddingsArgs)).
				MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
					searchResult := embedstore.TightenSearchResults(t.Data["searchResults"].(types.SearchResults).Results)
					return t.
						ApplyWorkflow(cliworkflows.PrintEmbeddings(searchResult)).
						ApplyWorkflow(embedworkflows.EmbedWorkflow(searchResult))
				})
		},
	}
}
