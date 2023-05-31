package redisembeddbconnector

import (
	"context"

	"github.com/RediSearch/redisearch-go/v2/redisearch"
)

func (idx *RedisembedTgenerator) AddEmbeddingDocument(ctx context.Context, docID string, embedding [1536]float32, metadata map[string]string) error {
	doc := redisearch.NewDocument(docID, 1.0).
		Set("oaiemb", toBytes(embedding))
	for k, v := range metadata {
		doc.Set(k, v)
	}

	return idx.client.IndexOptions(redisearch.IndexingOptions{Replace: true}, []redisearch.Document{doc}...)
}
func (idx *RedisembedTgenerator) DeleteEmbeddingDocument(ctx context.Context, docID string) error {
	return idx.client.DeleteDocument(docID)
}
