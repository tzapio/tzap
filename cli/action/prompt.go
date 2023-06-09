package action

import (
	"github.com/tzapio/tzap/cli/actionpb"
	"github.com/tzapio/tzap/cli/cmd/cliworkflows"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/pkg/embed/embedstore"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/util"
	"github.com/tzapio/tzap/workflows/code/embedworkflows"
	"github.com/tzapio/tzap/workflows/code/fileworkflows"
)

type PromptWorkflowArgs struct {
	InspirationFiles []string
	SearchQuery      string
	EmbedsCount      int32
	NCount           int32
	DisableIndex     bool
	Yes              bool
	MessageThread    []types.Message
}

// PromptWorkflow defines a workflow for generating code based on code-searching existing files and user input.
// This workflow creates a chat interface for the user to prompt for their desired output.
// The inspiration files parameter allows the user to recommend needed inspiration files to enhance GPT's general understanding.
// The search results are generated using a configurable number of embeddings and ranked by relevance.
// The user is then prompted to select a single search result. Once selected, the selected search result is used to generate the final desired output. Finally, the workflow is terminated when the completion of the user's message thread is reached.
func PromptWorkflow(promptWorkflowArgs PromptWorkflowArgs) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {

			var allFiles []string
			for _, ins := range promptWorkflowArgs.InspirationFiles {
				files, err := util.ListGlob(ins)
				if err != nil {
					panic(err)
				}
				allFiles = append(allFiles, files...)

			}

			loadAndSearchEmbeddingsArgs := &actionpb.SearchArgs{
				ExcludeFiles: allFiles,
				SearchQuery:  promptWorkflowArgs.SearchQuery,
				EmbedsCount:  promptWorkflowArgs.EmbedsCount,
				NCount:       promptWorkflowArgs.NCount,
				DisableIndex: promptWorkflowArgs.DisableIndex,
				Yes:          promptWorkflowArgs.Yes,
			}

			return t.
				// Include all full length inspiration files.
				ApplyWorkflow(cliworkflows.PrintInspirationFiles(allFiles)).
				ApplyWorkflow(fileworkflows.InspirationWorkflow(allFiles)).
				// Find all embeddings search them, returning the top results.
				ApplyWorkflow(LoadAndSearchEmbeddingsWorkflow(loadAndSearchEmbeddingsArgs)).
				MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
					searchResult := embedstore.TightenSearchResults(t.Data["searchResults"].(types.SearchResults).Results)
					return t.
						ApplyWorkflow(cliworkflows.PrintEmbeddings(searchResult)).
						ApplyWorkflow(embedworkflows.EmbedWorkflow(searchResult))
				}).
				// Append the current conversation thread
				LoadThread(promptWorkflowArgs.MessageThread).
				// Get Completion
				MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
					println(cmdutil.Bold("--- Completion"))
					t = t.RequestChatCompletion()
					println(cmdutil.Bold("\n---"))
					return t
				})
		},
	}
}
