package files

import (
	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func ProcessAndEmbedFilesTzapTemplate(dir string) types.NamedTemplate[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedTemplate[*tzap.Tzap, *tzap.Tzap]{
		Name: "processAndEmbedFilesTzapTemplate",
		Template: func(t *tzap.Tzap) *tzap.Tzap {
			embeddings := generateEmbeddings(t, dir)
			return saveEmbeddingsToTzap(t, embeddings)
		},
	}
}

func generateEmbeddings(t *tzap.Tzap, dir string) types.Embeddings {
	embed.ProcessDirectory(t, dir)
	embeddings, err := embed.GetEmbeddings(t, dir)
	if err != nil {
		panic(err)
	}
	return embeddings
}

func saveEmbeddingsToTzap(t *tzap.Tzap, embeddings types.Embeddings) *tzap.Tzap {
	for _, vector := range embeddings.Vectors {
		err := t.TG.AddEmbeddingDocument(t.C, vector.ID, vector.Values, vector.Metadata)
		if err != nil {
			panic(err)
		}
	}
	return t.AddSystemMessage("Embeddings saved to Tzap successfully.")
}
