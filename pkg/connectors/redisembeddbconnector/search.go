package redisembeddbconnector

import (
	"context"

	"github.com/RediSearch/redisearch-go/v2/redisearch"
	"github.com/tzapio/tzap/pkg/types"
)

func (idx *RedisembedTgenerator) ListAllEmbeddingsIds(ctx context.Context, projectName string) (types.SearchResults, error) {
	res, _, err := idx.client.Search(redisearch.NewQuery("*").Limit(0, 1000000))
	if err != nil {
		return types.SearchResults{}, err
	}
	results := types.SearchResults{}
	for _, doc := range res {
		result := types.Vector{ID: doc.Id, Metadata: types.Metadata{}}
		for key, value := range doc.Properties {
			_ = value
			panic(key)
			//result.Metadata[key] = value.(string)
		}
		results.Results = append(results.Results, types.SearchResult{Vector: result})
	}
	return results, nil
}
func (idx *RedisembedTgenerator) SearchWithEmbedding(ctx context.Context, projectName string, embedding types.QueryFilter, k int) (types.SearchResults, error) {
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
		result := types.Vector{ID: doc.Id, Metadata: types.Metadata{}}
		for key, value := range doc.Properties {
			_ = value
			panic(key)
			//result.Metadata[key] = value.(string)
		}
		results.Results = append(results.Results, types.SearchResult{Vector: result, Similarity: doc.Properties["__vec_score"].(float32)})
	}
	return results, nil
}
