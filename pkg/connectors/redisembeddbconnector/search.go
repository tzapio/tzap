package redisembeddbconnector

import (
	"context"

	"github.com/RediSearch/redisearch-go/v2/redisearch"
	"github.com/tzapio/tzap/pkg/types"
)

func (idx *RedisembedTgenerator) ListAllEmbeddingsIds(ctx context.Context) (types.SearchResults, error) {
	res, _, err := idx.client.Search(redisearch.NewQuery("*").Limit(0, 1000000))
	if err != nil {
		return types.SearchResults{}, err
	}
	results := types.SearchResults{}
	for _, doc := range res {
		result := types.Vector{ID: doc.Id, Metadata: map[string]string{}}
		for key, value := range doc.Properties {
			result.Metadata[key] = value.(string)
		}
		results.Results = append(results.Results, result)
	}
	return results, nil
}
func (idx *RedisembedTgenerator) SearchWithEmbedding(ctx context.Context, embedding types.QueryFilter, k int) (types.SearchResults, error) {
	queryStr := "*=>[KNN 10 @oaiemb $embeddingbytes AS __vec_score]"
	query := redisearch.
		NewQuery(queryStr).
		Limit(0, k)
	query.SetSortBy("__vec_score", true)
	query.SetParams(map[string]interface{}{"embeddingbytes": toBytes(embedding.Values)})
	query.Dialect = 2

	res, _, err := idx.client.Search(query)
	if err != nil {
		return types.SearchResults{}, err
	}

	results := types.SearchResults{}
	for _, doc := range res {
		result := types.Vector{ID: doc.Id, Metadata: map[string]string{}}
		for key, value := range doc.Properties {
			result.Metadata[key] = value.(string)
		}
		results.Results = append(results.Results, result)
	}
	return results, nil
}
