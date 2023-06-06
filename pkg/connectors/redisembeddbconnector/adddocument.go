package redisembeddbconnector

import (
	"context"

	"github.com/RediSearch/redisearch-go/v2/redisearch"
	"github.com/tzapio/tzap/pkg/types"
)

func (idx *RedisembedTgenerator) AddEmbeddingDocument(ctx context.Context, projectName string, docID string, embedding [1536]float32, metadata types.Metadata) error {
	doc := redisearch.NewDocument(docID, 1.0).
		Set("oaiemb", toBytes(embedding))

	panic("")
	//doc.Set(k, v)

	return idx.client.IndexOptions(redisearch.IndexingOptions{Replace: true}, []redisearch.Document{doc}...)
}
func (idx *RedisembedTgenerator) DeleteEmbeddingDocument(ctx context.Context, projectName string, docID string) error {
	return idx.client.DeleteDocument(docID)
}
