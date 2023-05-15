package embed

import (
	"context"
	"encoding/json"
	"os"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func CreateQueryJSON(t *tzap.Tzap, filename string, input string) error {
	query, err := GetQuery(t, input)
	if err != nil {
		return err
	}
	if err := SaveQueryAsJSON(query, filename); err != nil {
		return err
	}
	return nil
}

func GetQuery(t *tzap.Tzap, input string) (types.QueryRequest, error) {
	println("\n\nGetting query for input: " + input + "\n\n")
	embeddings, err := getEmbeddings(t, input)
	if err != nil {
		return types.QueryRequest{}, err
	}
	queryFilters := CreateQueryFilters(embeddings)
	query := BuildQuery(queryFilters)
	return query, nil
}

func getEmbeddings(t *tzap.Tzap, input string) ([][]float32, error) {
	embeddings, err := t.TG.FetchEmbedding(context.Background(), input)
	if err != nil {
		return nil, err
	}
	return embeddings, nil
}

func CreateQueryFilters(embeddings [][]float32) []types.QueryFilter {
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

func SaveQueryAsJSON(query types.QueryRequest, filename string) error {
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, queryJSON, 0644); err != nil {
		return err
	}
	return nil
}
