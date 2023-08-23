package action

import (
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
)

func Completion(t *tzap.Tzap, completionRequest *actionpb.CompletionRequest) (*actionpb.CompletionResponse, error) {
	t = t.ApplyWorkflow(CompletionWorkflow(completionRequest.CompletionArgs)).AsAssistantMessage()
	messages := t.GetThread()
	return &actionpb.CompletionResponse{
		Thread: ToPBMessage(messages),
	}, nil
}
func CompletionWorkflow(completionWorkflowArgs *actionpb.CompletionArgs) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {

			return t.
				LoadThread(ToTzapMessage(completionWorkflowArgs.Thread)).
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
