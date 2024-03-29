package embed

import (
	"encoding/json"

	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/util/reflectutil"
)

type EmbeddingCache struct {
	embeddingCacheDB types.DBCollectionInterface[string]
}

func NewEmbeddingCache(embeddingCacheDB types.DBCollectionInterface[string]) *EmbeddingCache {
	return &EmbeddingCache{embeddingCacheDB}
}

func (ec *EmbeddingCache) GetCachedEmbeddings(files []types.FileReader, embeddings *types.Embeddings) (*types.Embeddings, error) {
	tl.Logger.Println("Getting cached embeddings", len(embeddings.Vectors))
	var cachedEmbeddings []*types.Vector

	for _, vector := range embeddings.Vectors {
		splitPart := vector.Metadata.SplitPart
		kv, exists := ec.embeddingCacheDB.ScanGet(splitPart)
		if exists {
			if !reflectutil.IsZero(kv.Value) {
				var float32Vector [1536]float32
				err := json.Unmarshal([]byte(kv.Value), &float32Vector)
				if err != nil {
					return nil, err
				}

				if len(float32Vector) == 1536 {
					cachedVector := &types.Vector{
						ID:        vector.ID,
						TimeStamp: 0,
						Metadata:  vector.Metadata,
						Values:    float32Vector,
					}

					cachedEmbeddings = append(cachedEmbeddings, cachedVector)
					continue
				} else {
					println("invalid vector length", splitPart)
					return &types.Embeddings{}, nil
				}
			}
		}

		println("Warning: %s is uncached.", vector.ID)
	}

	return &types.Embeddings{Vectors: cachedEmbeddings}, nil
}

func (ec *EmbeddingCache) GetUncachedEmbeddings(embeddings *types.Embeddings) *types.Embeddings {
	var uncachedEmbeddings []*types.Vector

	for _, vector := range embeddings.Vectors {
		splitPart := vector.Metadata.SplitPart
		kv, exists := ec.embeddingCacheDB.ScanGet(splitPart)
		if !exists || reflectutil.IsZero(kv.Value) {
			uncachedEmbeddings = append(uncachedEmbeddings, vector)
		}
	}

	return &types.Embeddings{Vectors: uncachedEmbeddings}
}

func (ec *EmbeddingCache) FetchThenCacheNewEmbeddings(t *tzap.Tzap, files []types.FileReader, uncachedEmbeddings *types.Embeddings) error {
	storedFiles := map[string]struct{}{}

	if len(uncachedEmbeddings.Vectors) > 0 {
		batchSize := 200

		for i := 0; i < len(uncachedEmbeddings.Vectors); i += batchSize {
			end := i + batchSize
			if end > len(uncachedEmbeddings.Vectors) {
				end = len(uncachedEmbeddings.Vectors)
			}

			batch := uncachedEmbeddings.Vectors[i:end]
			var inputStrings []string
			for _, vector := range batch {
				storedFiles[vector.Metadata.Filename] = struct{}{}
				inputStrings = append(inputStrings, vector.Metadata.SplitPart)
			}

			embeddingsResult, err := t.TG.FetchEmbedding(t.C, inputStrings...)
			if err != nil {
				return err
			}

			cacheKeyVal := make([]types.KeyValue[string], len(embeddingsResult))
			for i, embedding := range embeddingsResult {
				embBytes, err := json.Marshal(embedding)
				if err != nil {
					return err
				}
				cacheKeyVal[i] = types.KeyValue[string]{Key: inputStrings[i], Value: string(embBytes)}
			}

			added, err := ec.embeddingCacheDB.BatchSet(cacheKeyVal)
			if err != nil {
				return err
			}

			tl.UILogger.Println("Added", added, "embeddings to cache")
		}
	}

	return nil
}
