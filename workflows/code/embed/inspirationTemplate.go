package embed

import (
	"fmt"

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
			fmt.Printf("Including static inspiration files: %+v\n", inspirationFiles)
			return t.
				LoadFiles(inspirationFiles).
				Accumulate(func(t *tzap.Tzap) *tzap.Tzap {
					return t.AddSystemMessage(
						"###file: "+t.Data["filepath"].(string),
						t.Data["content"].(string))
				})
		},
	}
}