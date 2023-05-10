package localdbconnector

import (
	"context"
	"encoding/json"

	"github.com/tzapio/tzap/pkg/embed/cosine"
	"github.com/tzapio/tzap/pkg/types"
)

func (idx *LocalembedTGenerator) ListAllEmbeddingsIds(ctx context.Context) (types.SearchResults, error) {
	allResults := idx.db.GetAll()
	listEmbeddings := types.SearchResults{}
	for i, doc := range allResults {
		vector := types.Vector{}

		if err := json.Unmarshal([]byte(doc.Value), &vector); err != nil {
			println("error unmarshalling vector", i, err.Error())
			continue
		}
		listEmbeddings.Results = append(listEmbeddings.Results, vector)
	}
	return listEmbeddings, nil
}
func (idx *LocalembedTGenerator) SearchWithEmbedding(ctx context.Context, embedding types.QueryFilter, k int) (types.SearchResults, error) {
	println("searching with embedding")
	res := idx.db.GetAll()
	floatVectors := [][]float32{}
	vectors := []types.Vector{}
	for i, r := range res {
		vector := types.Vector{}
		err := json.Unmarshal([]byte(r.Value), &vector)
		if err != nil {
			println("error unmarshalling vector", i, r.Value, err.Error())
			return types.SearchResults{}, err
		}
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
	vectorString, exists := idx.db.Get(docID)
	vector := types.Vector{}
	if !exists {
		return vector, exists, nil
	}
	if err := json.Unmarshal([]byte(vectorString), &vector); err != nil {
		return vector, false, err
	}
	return vector, true, nil
}
