package embedstore

import (
	"context"

	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/embed/cosine"
	"github.com/tzapio/tzap/pkg/project"
	"github.com/tzapio/tzap/pkg/types"
)

func (idx *EmbedStore) ListAllEmbeddingsIds(ctx context.Context) (types.SearchResults, error) {
	tl.Logger.Println("ListAllEmbeddings")
	embeddingCollection := project.GetProjectFromContext(ctx).GetEmbeddingCollection()
	allResults := embeddingCollection.GetAll()
	listEmbeddings := types.SearchResults{}
	for _, vector := range allResults {
		listEmbeddings.Results = append(listEmbeddings.Results, types.SearchResult{
			Vector:     vector.Value,
			Similarity: 0.0})
	}
	return listEmbeddings, nil
}
func (idx *EmbedStore) SearchWithEmbedding(ctx context.Context, embedding types.QueryFilter, k int) (types.SearchResults, error) {
	tl.Logger.Println("SearchWithEmbedding")
	embeddingCollection := project.GetProjectFromContext(ctx).GetEmbeddingCollection()
	res := embeddingCollection.GetAll()
	floatVectors := [][1536]float32{}
	vectors := []types.Vector{}
	for _, r := range res {
		vector := r.Value
		floatVectors = append(floatVectors, vector.Values)
		vectors = append(vectors, vector)
	}

	results := cosine.SearchByCosineSimilarity(floatVectors, embedding.Values)
	searchResults := types.SearchResults{}
	if len(results) < k || k <= -1 {
		k = len(results)
	}
	//pcaResult := pca.EmbeddingsTo3D(floatVectors)
	for _, r := range results[:k] {
		searchResults.Results = append(searchResults.Results, types.SearchResult{
			Vector: vectors[r.Index],
			//PCA:        pcaResult[r.Index],
			Similarity: r.Similarity,
		})
	}
	return searchResults, nil
}
func (idx *EmbedStore) GetEmbeddingDocument(ctx context.Context, docID string) (types.Vector, bool, error) {
	embeddingCollection := project.GetProjectFromContext(ctx).GetEmbeddingCollection()
	vector, exists := embeddingCollection.Get(docID)

	if !exists {
		return vector, exists, nil
	}

	return vector, true, nil
}
