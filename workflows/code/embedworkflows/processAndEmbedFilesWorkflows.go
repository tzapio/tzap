package embedworkflows

import (
	"fmt"

	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func PrepareEmbedFilesWorkflow(files []string, embedder *embed.Embedder) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "prepareEmbedFilesWorkflow",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {

			rawFileEmbeddings := embedder.PrepareEmbeddingsFromFiles(files)
			uncachedEmbeddings := embedder.GetUncachedEmbeddings(rawFileEmbeddings)
			data := types.MappedInterface{"rawFileEmbeddings": rawFileEmbeddings, "uncachedEmbeddings": uncachedEmbeddings, "embedder": embedder}
			return t.AddTzap(&tzap.Tzap{Name: "prepareEmbedFilesTzap", Data: data})
		},
	}
}

func FetchOrCachedEmbeddingForFilesWorkflow() types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "fetchOrCachedEmbeddingForFilesWorkflow",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			embedder, ok := t.Data["embedder"].(*embed.Embedder)
			if !ok {
				if err := embedder.FetchAndCacheNewEmbeddings(t, types.Embeddings{}); err != nil {
					panic(err)
				}
				panic(fmt.Errorf("loading embeddings went wrong. embedder"))
			}
			uncachedEmbeddings, ok := t.Data["uncachedEmbeddings"].(types.Embeddings)
			if !ok {
				panic("Loading embeddings went wrong")
			}
			if len(uncachedEmbeddings.Vectors) > 0 {
				if err := embedder.FetchAndCacheNewEmbeddings(t, uncachedEmbeddings); err != nil {
					panic(err)
				}
			}
			rawFileEmbeddings, ok := t.Data["rawFileEmbeddings"].(types.Embeddings)
			if !ok {
				panic("Loading embeddings went wrong")
			}
			cachedEmbeddings := embedder.GetCachedEmbeddings(rawFileEmbeddings)
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

func saveEmbeddingsToTzap(t *tzap.Tzap, embeddings types.Embeddings) *tzap.Tzap {
	for _, vector := range embeddings.Vectors {
		err := t.TG.AddEmbeddingDocument(t.C, vector.ID, vector.Values, vector.Metadata)
		if err != nil {
			panic(err)
		}
	}
	return t
}