package gptasfunction

import (
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func GPTAsFunction(task string, content string) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "GPTAsFunction",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			return t.
				AddSystemMessage(task).
				AddUserMessage(content).
				RequestChatCompletion()
		},
	}
}
