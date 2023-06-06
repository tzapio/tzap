package localdbconnector

import (
	"context"

	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/embed/cosine"
	"github.com/tzapio/tzap/pkg/project"
	"github.com/tzapio/tzap/pkg/types"
)

func (idx *LocalembedTGenerator) ListAllEmbeddingsIds(ctx context.Context, projectName string) (types.SearchResults, error) {
	tl.Logger.Println("ListAllEmbeddings - Project:", projectName)
	allResults := idx.dbs[project.ProjectName(projectName)].GetAll()
	listEmbeddings := types.SearchResults{}
	for _, vector := range allResults {
		listEmbeddings.Results = append(listEmbeddings.Results, types.SearchResult{
			Vector:     vector.Value,
			Similarity: 0.0})
	}
	return listEmbeddings, nil
}
func (idx *LocalembedTGenerator) SearchWithEmbedding(ctx context.Context, projectName string, embedding types.QueryFilter, k int) (types.SearchResults, error) {
	res := idx.dbs[project.ProjectName(projectName)].GetAll()
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
		searchResults.Results = append(searchResults.Results, types.SearchResult{
			Vector:     vectors[r.Index],
			Similarity: r.Similarity})
	}
	return searchResults, nil
}
func (idx *LocalembedTGenerator) GetEmbeddingDocument(ctx context.Context, projectName string, docID string) (types.Vector, bool, error) {
	vector, exists := idx.dbs[project.ProjectName(projectName)].Get(docID)

	if !exists {
		return vector, exists, nil
	}

	return vector, true, nil
}
