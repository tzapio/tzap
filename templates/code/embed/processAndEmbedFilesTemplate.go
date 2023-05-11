package embed

import (
	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func ProcessAndEmbedFilesTzapTemplate(files []string) types.NamedTemplate[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedTemplate[*tzap.Tzap, *tzap.Tzap]{
		Name: "processAndEmbedFilesTzapTemplate",
		Template: func(t *tzap.Tzap) *tzap.Tzap {
			// Generate and load embedding files.
			embeddings := generateEmbeddings(t, files)
			// Store in local vector db (if default tzapconnector)
			return saveEmbeddingsToTzap(t, embeddings)
		},
	}
}

func generateEmbeddings(t *tzap.Tzap, files []string) types.Embeddings {
	// Generate embedding files.
	embed.OutputEmbeddingsToFile(t, files)
	// Load embeddings from file.
	embeddings, err := embed.GetEmbeddingsFromFile()
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
