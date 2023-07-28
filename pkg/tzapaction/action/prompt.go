package action

import (
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
	"github.com/tzapio/tzap/pkg/tzapaction/cliworkflows"
	"github.com/tzapio/tzap/pkg/util"
	"github.com/tzapio/tzap/workflows/code/fileworkflows"
)

func Prompt(t *tzap.Tzap, promptRequest *actionpb.PromptRequest) (*actionpb.PromptResponse, error) {
	t = t.ApplyWorkflow(PromptWorkflow(promptRequest.PromptArgs)).AsAssistantMessage()
	messages := t.GetThread()
	return &actionpb.PromptResponse{
		Thread: ToPBMessage(messages),
	}, nil
}
func PromptWorkflow(promptWorkflowArgs *actionpb.PromptArgs) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
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
				SearchQuery:  promptWorkflowArgs.SearchQuery,
				EmbedsCount:  promptWorkflowArgs.EmbedsCount,
			}

			return t.
				// Include all full length inspiration files.
				AddSystemMessage("You are code assisting AI. Fulfill the user requested by looking at what the user has added as inspiration files and what the search engine has returned as embeddings.").
				ApplyWorkflow(cliworkflows.PrintInspirationFiles(allFiles)).
				ApplyWorkflow(fileworkflows.InspirationWorkflow(allFiles)).
				// Find all embeddings search them, returning the top results.
				ApplyWorkflow(SearchWorkflow(loadAndSearchEmbeddingsArgs)).
				// Append the current conversation thread
				LoadThread(ToTzapMessage(promptWorkflowArgs.Thread)).
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
