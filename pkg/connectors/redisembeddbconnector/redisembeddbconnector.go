package redisembeddbconnector

import (
	"github.com/RediSearch/redisearch-go/v2/redisearch"
	"github.com/tzapio/tzap/pkg/types"
)

type RedisembedTgenerator struct {
	*types.UnimplementedTGenerator

	client *redisearch.Client
}

func InitiateRedisClient(addr string) (types.TGenerator, error) {
	client, err := NewEmbeddingIndex(addr, "files", "oaiemb", 1536)
	if err != nil {
		println("Cannot connect to " + addr + " - redis disabled")
		return nil, nil
	}

	return &RedisembedTgenerator{client: client}, nil
}

func NewEmbeddingIndex(addr string, indexName, fieldName string, dimensions int) (*redisearch.Client, error) {
	client := redisearch.NewClient(addr, indexName)
	// Create the index
	println("creating index", indexName, "field", fieldName, "dimensions", dimensions)
	client.Drop()
	schema := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewVectorFieldOptions(fieldName, redisearch.VectorFieldOptions{
			Algorithm: redisearch.HNSW,
			Attributes: map[string]interface{}{
				"TYPE":            "FLOAT32",
				"DIM":             1536,
				"DISTANCE_METRIC": "L2",
			},
		}),
		)

	err := client.CreateIndex(schema)
	if err != nil {
		//print failed to create index
		println("failed to create index")
		//return nil, err
	}

	return client, nil
}
