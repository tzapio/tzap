package embed

import (
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

type EmbedCleaner struct {
}

func (ec *EmbedCleaner) CleanOldEmbeddings(t *tzap.Tzap, rawFileEmbeddings *types.Embeddings, unchangedFileTimestamps map[string]int64) {
	storedEmbeddings, err := t.TG.ListAllEmbeddingsIds(t.C)
	if err != nil {
		panic(err)
	}
	idsToDelete, err := ec.getNoLongerPresentEmbeddings(storedEmbeddings, rawFileEmbeddings, unchangedFileTimestamps)
	if err != nil {
		panic(err)
	}
	tl.Logger.Println("Removing old embeddings", len(idsToDelete))
	if err := ec.removeNoLongerPresentEmbeddings(t, idsToDelete); err != nil {
		panic(err)
	}
}
func (ec *EmbedCleaner) getNoLongerPresentEmbeddings(storedEmbeddings types.SearchResults, nowEmbeddings *types.Embeddings, unchangedFiles map[string]int64) ([]string, error) {
	nowEmbeddingsIds := make(map[string]struct{})
	for _, vectorID := range nowEmbeddings.Vectors {
		nowEmbeddingsIds[vectorID.ID] = struct{}{}
	}
	missingIds := []string{}
	for _, storedVector := range storedEmbeddings.Results {
		filename := storedVector.Vector.Metadata.Filename
		if _, exists := unchangedFiles[filename]; exists {
			continue
		}
		if _, exists := nowEmbeddingsIds[storedVector.Vector.ID]; !exists {
			missingIds = append(missingIds, storedVector.Vector.ID)
			tl.Logger.Println("Drift: ", storedVector)
		}
	}
	tl.Logger.Println("Drift Check: ", len(storedEmbeddings.Results), len(nowEmbeddings.Vectors), len(missingIds))
	return missingIds, nil
}
func (ec *EmbedCleaner) removeNoLongerPresentEmbeddings(t *tzap.Tzap, deleteIds []string) error {
	if err := t.TG.DeleteEmbeddingDocuments(t.C, deleteIds); err != nil {
		return err
	}
	return nil
}
