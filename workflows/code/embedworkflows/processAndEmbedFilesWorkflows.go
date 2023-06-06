package embedworkflows

import (
	"fmt"

	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/project"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func PrepareEmbedFilesWorkflow(name project.ProjectName, files []types.FileReader, embedder *embed.Embedder) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "prepareEmbedFilesWorkflow",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			rawFileEmbeddings := embedder.PrepareEmbeddingsFromFiles()
			uncachedEmbeddings := embedder.GetUncachedEmbeddings(rawFileEmbeddings)
			data := types.MappedInterface{"rawFileEmbeddings": rawFileEmbeddings, "uncachedEmbeddings": uncachedEmbeddings, "embedder": embedder}
			return t.AddTzap(&tzap.Tzap{Name: "prepareEmbedFilesTzap", Data: data})
		},
	}
}

func FetchOrCachedEmbeddingForFilesWorkflow(files []types.FileReader) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "fetchOrCachedEmbeddingForFilesWorkflow",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			embedder, ok := t.Data["embedder"].(*embed.Embedder)
			if !ok {
				panic(fmt.Errorf("loading embedder went wrong"))
			}
			uncachedEmbeddings, ok := t.Data["uncachedEmbeddings"].(types.Embeddings)
			if !ok {
				panic("Loading embeddings went wrong")
			}
			if len(uncachedEmbeddings.Vectors) > 0 {
				if err := embedder.FetchThenCacheNewEmbeddings(t, files, uncachedEmbeddings); err != nil {
					panic(err)
				}
			}
			rawFileEmbeddings, ok := t.Data["rawFileEmbeddings"].(types.Embeddings)
			if !ok {
				panic("Loading embeddings went wrong")
			}
			cachedEmbeddings, err := embedder.GetCachedEmbeddings(files, rawFileEmbeddings)
			if err != nil {
				panic(err)
			}
			data := types.MappedInterface{"embeddings": cachedEmbeddings}
			return t.AddTzap(&tzap.Tzap{Name: "fetchOrCachedEmbeddingForFilesTzap", Data: data})
		},
	}
}
func SaveAndLoadEmbeddingsToDB(name project.ProjectName) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "saveAndLoadEmbeddingsToDB",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			embeddings, ok := t.Data["embeddings"].(types.Embeddings)
			if !ok {
				panic("Loading embeddings went wrong")
			}
			for _, vector := range embeddings.Vectors {
				err := t.TG.AddEmbeddingDocument(t.C, string(name), vector.ID, vector.Values, vector.Metadata)
				if err != nil {
					panic(err)
				}
			}
			// Store in local vector db (if default tzapconnector)
			return t
		},
	}
}
