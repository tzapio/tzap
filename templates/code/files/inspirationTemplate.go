package files

import (
	"fmt"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func InspirationTemplate(inspirationFiles []string) types.NamedTemplate[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedTemplate[*tzap.Tzap, *tzap.Tzap]{
		Name: "inspirationFiles",
		Template: func(t *tzap.Tzap) *tzap.Tzap {
			// print inspirationFiles
			fmt.Printf("%+v\n", inspirationFiles)
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
