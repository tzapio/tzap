package localdbconnector

import (
	"context"

	"github.com/tzapio/tzap/pkg/embed/cosine"
	"github.com/tzapio/tzap/pkg/types"
)

func (idx *LocalembedTGenerator) ListAllEmbeddingsIds(ctx context.Context) (types.SearchResults, error) {
	allResults := idx.db.GetAll()
	listEmbeddings := types.SearchResults{}
	for _, vector := range allResults {
		listEmbeddings.Results = append(listEmbeddings.Results, vector.Value)
	}
	return listEmbeddings, nil
}
func (idx *LocalembedTGenerator) SearchWithEmbedding(ctx context.Context, embedding types.QueryFilter, k int) (types.SearchResults, error) {
	res := idx.db.GetAll()
	floatVectors := [][1536]float32{}
	vectors := []types.Vector{}
	for _, r := range res {
		vector := r.Value
		floatVectors = append(floatVectors, vector.Values)
		vectors = append(vectors, vector)
	}

	results := cosine.SearchByCosineSimilarity(floatVectors, embedding.Values)
	searchResults := types.SearchResults{}
	if len(results) < k {
		k = len(results)
	}
	for _, r := range results[:k] {
		searchResults.Results = append(searchResults.Results, vectors[r.Index])
	}
	return searchResults, nil
}
func (idx *LocalembedTGenerator) GetEmbeddingDocument(ctx context.Context, docID string) (types.Vector, bool, error) {
	vector, exists := idx.db.Get(docID)

	if !exists {
		return vector, exists, nil
	}

	return vector, true, nil
}
