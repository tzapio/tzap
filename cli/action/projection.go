package action

import (
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	// import additional required packages
)

func FutureProjection(priorModel string, priorResult string, newModel string) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "projection",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			return t.
				AddSystemMessage("You are analyzing a prior result based on a prior definition and now you will attempt to project it based on a new definition to get a new result. The user will write the new model and you will output the new result.",
					"####"+"Prior Model:\n"+priorModel,
					"####"+"Prior Result:\n"+priorResult).
				AddUserMessage(newModel).
				RequestChatCompletion()
		},
	}
}
