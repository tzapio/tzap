package fileworkflows

import (
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func InspirationWorkflow(inspirationFiles []string) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "inspirationWorkflow",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			if len(inspirationFiles) == 0 {
				return t
			}
			return t.
				AddSystemMessage("####The following files are explictily included by the user for relevance: ").
				LoadFiles(inspirationFiles).
				Reduce(func(t *tzap.Tzap, child *tzap.Tzap) *tzap.Tzap {
					return t.AddSystemMessage(
						"####file: "+child.Data["filepath"].(string),
						child.Data["content"].(string))
				})
		},
	}
}
