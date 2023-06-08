package embed

import (
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func NewQuery(t *tzap.Tzap, input string) (types.QueryRequest, error) {
	embeddings, err := getEmbeddings(t, input)
	if err != nil {
		return types.QueryRequest{}, err
	}
	queryFilters := CreateQueryFilters(embeddings)
	query := BuildQuery(queryFilters)
	return query, nil
}

func getEmbeddings(t *tzap.Tzap, input string) ([][1536]float32, error) {
	embeddings, err := t.TG.FetchEmbedding(t.C, input)
	if err != nil {
		return nil, err
	}
	return embeddings, nil
}

func CreateQueryFilters(embeddings [][1536]float32) []types.QueryFilter {
	var queryFilters []types.QueryFilter
	for _, embedding := range embeddings {
		queryFilters = append(queryFilters, types.QueryFilter{
			Filter: nil, // No filter is applied; adjust this according to your needs
			Values: embedding,
		})
	}
	return queryFilters
}

func BuildQuery(queryFilters []types.QueryFilter) types.QueryRequest {
	return types.QueryRequest{
		TopK:            10,   // Adjust this value according to your needs
		IncludeMetadata: true, // Include metadata in the response
		Namespace:       "",   // Set the appropriate namespace
		Queries:         queryFilters,
	}
}
