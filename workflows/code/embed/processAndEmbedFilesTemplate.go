package embed

import (
	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func PrepareEmbedFilesTzapWorkflow(files []string) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "prepareEmbedFilesTzapWorkflow",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			rawFileEmbeddings := embed.PrepareEmbeddingsFromFiles(t, files)
			uncachedEmbeddings := embed.GetUncachedEmbeddings(rawFileEmbeddings)
			data := types.MappedInterface{"rawFileEmbeddings": rawFileEmbeddings, "uncachedEmbeddings": uncachedEmbeddings}
			return t.AddTzap(&tzap.Tzap{Name: "prepareEmbedFilesTzap", Data: data})
		},
	}
}

func FetchOrCachedEmbeddingForFilesTzapWorkflow() types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "fetchOrCachedEmbeddingForFilesTzapWorkflow",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			uncachedEmbeddings, ok := t.Data["uncachedEmbeddings"].(types.Embeddings)
			if !ok {
				panic("Loading embeddings went wrong")
			}
			if len(uncachedEmbeddings.Vectors) > 0 {
				if err := embed.FetchAndCacheNewEmbeddings(t, uncachedEmbeddings); err != nil {
					panic(err)
				}
			}
			rawFileEmbeddings, ok := t.Data["rawFileEmbeddings"].(types.Embeddings)
			if !ok {
				panic("Loading embeddings went wrong")
			}
			cachedEmbeddings := embed.GetCachedEmbeddings(rawFileEmbeddings)
			embed.SaveEmbeddingToFile(cachedEmbeddings)

			data := types.MappedInterface{"embeddings": cachedEmbeddings}
			return t.AddTzap(&tzap.Tzap{Name: "fetchOrCachedEmbeddingForFilesTzap", Data: data})
		},
	}
}
func SaveAndLoadEmbeddingsToDB() types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "saveAndLoadEmbeddingsToDB",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			embeddings, ok := t.Data["embeddings"].(types.Embeddings)
			if !ok {
				panic("Loading embeddings went wrong")
			}

			// Store in local vector db (if default tzapconnector)
			return saveEmbeddingsToTzap(t, embeddings)
		},
	}
}

func generateEmbeddings(t *tzap.Tzap, files []string) types.Embeddings {
	// Generate embedding files.
	embed.PrepareEmbeddingsFromFiles(t, files)
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
	return t
}