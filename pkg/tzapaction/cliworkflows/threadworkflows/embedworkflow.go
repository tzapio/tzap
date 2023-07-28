package threadworkflows

import (
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
)

func EmbedWorkflow(embeddings []*actionpb.Embedding) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "SearchWorkflow",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			if len(embeddings) > 0 {
				t = t.AddSystemMessage(
					"The following file contents are embeddings for the user input:",
				)
				for _, embedding := range embeddings {
					t = t.AddSystemMessage(embedding.Content)
				}
			}
			return t

		},
	}
}
