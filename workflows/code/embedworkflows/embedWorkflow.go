package embedworkflows

import (
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func EmbedWorkflow(searchResults types.SearchResults) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "SearchWorkflow",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			if len(searchResults.Results) > 0 {
				t = t.AddSystemMessage(
					"The following file contents are embeddings for the user input:",
				)
				for _, result := range searchResults.Results {
					t = t.AddSystemMessage(result.Vector.Metadata.SplitPart)
				}
			}
			return t

		},
	}
}
