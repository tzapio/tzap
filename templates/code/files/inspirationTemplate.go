package files

import (
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func InspirationTemplate(inspirationFiles []string) types.NamedTemplate[*tzap.Tzap] {
	return types.NamedTemplate[*tzap.Tzap]{
		Name: "inspirationFiles",
		Template: func(t *tzap.Tzap) *tzap.Tzap {
			return t.
				LoadFiles(inspirationFiles).
				Accumulate(func(t *tzap.Tzap) *tzap.Tzap {
					return t.AddUserMessage(
						"//file: "+t.Data["filepath"].(string),
						t.Data["content"].(string))
				})
		},
	}
}
