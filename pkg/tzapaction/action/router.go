package action

import (
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
	"github.com/tzapio/tzap/pkg/tzapaction/cliworkflows"
	"github.com/tzapio/tzap/pkg/util"
	"github.com/tzapio/tzap/workflows/code/fileworkflows"
)

func RouterWorkflow(promptWorkflowArgs *actionpb.PromptArgs) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
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
				ExcludeFiles: append(promptWorkflowArgs.ExcludeFiles, allFiles...),
				SearchQuery:  promptWorkflowArgs.SearchArgss[0].SearchQuery,
				Lib:          promptWorkflowArgs.SearchArgss[0].Lib,
				EmbedsCount:  promptWorkflowArgs.SearchArgss[0].EmbedsCount,
			}

			return t.
				// Include all full length inspiration files.
				AddSystemMessage("Your task is to call other LLMs using function calls. You will see embeddings from search results and explicitly linked inspiration files. With this in mind you construct one of the function calls in order to deliver code as detailed as possible.").
				ApplyWorkflow(cliworkflows.PrintInspirationFiles(allFiles)).
				ApplyWorkflow(fileworkflows.InspirationWorkflow(allFiles)).
				// Find all embeddings search them, returning the top results.
				ApplyWorkflow(SearchWorkflow(loadAndSearchEmbeddingsArgs)).

				// Append the current conversation thread
				LoadThread(ToTzapMessage(promptWorkflowArgs.Thread)).
				// Get Completion
				MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
					println(cmdutil.Bold("--- Completion"))
					t = t.RequestFunctionCompletion(Tfs)
					println(cmdutil.Bold("\n---"))
					if t.Data["content"].(types.CompletionMessage).FunctionCall != nil {
						tl.Logger.Println(cmdutil.Bold("Function Call:"))
						tl.Logger.Println(t.Data["content"].(types.CompletionMessage).FunctionCall.Name)
						tl.Logger.Println(t.Data["content"].(types.CompletionMessage).FunctionCall.Arguments)
					}
					return t
				})
		},
	}
}
