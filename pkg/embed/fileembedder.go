package embed

import (
	"encoding/json"
	"os"

	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/embed/localdb"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/util"
)

func GetCachedEmbeddings(embeddings types.Embeddings) types.Embeddings {
	var storedFiles map[string]struct{} = map[string]struct{}{}

	tl.Logger.Println("Getting cached embeddings", len(embeddings.Vectors))
	var cachedEmbeddings []types.Vector
	db, err := localdb.NewFileDB("./.tzap-data/embeddingsCache.db")
	if err != nil {
		panic(err)
	}

	for _, vector := range embeddings.Vectors {
		splitPart := vector.Metadata["splitPart"]
		kv, exists := db.ScanGet(splitPart)
		if exists {
			if kv.Value != "" {
				float32Vector := []float32{}
				if err := json.Unmarshal([]byte(kv.Value), &float32Vector); err != nil {
					println("cannot unmarshal", splitPart, err.Error())
					continue
				}
				if len(float32Vector) == 1536 {
					vector := types.Vector{
						ID:        vector.ID,
						TimeStamp: 0,
						Metadata:  vector.Metadata,
						Values:    float32Vector,
					}

					cachedEmbeddings = append(cachedEmbeddings, vector)
					storedFiles[vector.Metadata["filename"]] = struct{}{}
					continue
				} else {
					println("invalid vector length", splitPart, err.Error())
					continue
				}
			}
		}
		println("Warning: %s is uncached.", vector.ID)
	}
	db2, err := localdb.NewFileDB("./.tzap-data/filesmd5.db")
	if err != nil {
		panic(err)
	}
	var keyvals []types.KeyValue
	for file := range storedFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}
		keyvals = append(keyvals, types.KeyValue{Key: file, Value: util.MD5HashByte(content)})
	}
	if len(keyvals) > 0 {
		added, err := db2.BatchSet(keyvals)
		if err != nil {
			panic("failing to store changed files should not happend and has probably caused some kind of corruption")
		}
		tl.Logger.Printf("Added %d files to md5 cache. Total: %d", added, len(storedFiles))
	}
	return types.Embeddings{Vectors: cachedEmbeddings}
}

func GetUncachedEmbeddings(embeddings types.Embeddings) types.Embeddings {
	var uncachedEmbeddings []types.Vector
	db, err := localdb.NewFileDB("./.tzap-data/embeddingsCache.db")
	if err != nil {
		panic(err)
	}

	for _, vector := range embeddings.Vectors {
		splitPart := vector.Metadata["splitPart"]
		kv, exists := db.ScanGet(splitPart)
		if !exists || kv.Value == "" {
			uncachedEmbeddings = append(uncachedEmbeddings, vector)
		}

	}
	return types.Embeddings{Vectors: uncachedEmbeddings}
}

func FetchAndCacheNewEmbeddings(t *tzap.Tzap, uncachedEmbeddings types.Embeddings) error {
	var storedFiles map[string]struct{} = map[string]struct{}{}

	if len(uncachedEmbeddings.Vectors) > 0 {
		batchSize := 100

		db, err := localdb.NewFileDB("./.tzap-data/embeddingsCache.db")
		if err != nil {
			panic(err)
		}

		for i := 0; i < len(uncachedEmbeddings.Vectors); i += batchSize {
			end := i + batchSize
			if end > len(uncachedEmbeddings.Vectors) {
				end = len(uncachedEmbeddings.Vectors)
			}

			batch := uncachedEmbeddings.Vectors[i:end]
			var inputStrings []string
			for _, vector := range batch {
				storedFiles[vector.Metadata["filename"]] = struct{}{}
				inputStrings = append(inputStrings, vector.Metadata["splitPart"])
			}

			embeddingsResult, err := t.TG.FetchEmbedding(t.C, inputStrings...)
			if err != nil {
				return err
			}

			cacheKeyVal := []types.KeyValue{}
			for i, embedding := range embeddingsResult {
				embeddingByte, err := json.Marshal(embedding)
				if err != nil {
					return err
				}
				cacheKeyVal = append(cacheKeyVal, types.KeyValue{Key: inputStrings[i], Value: string(embeddingByte)})
			}

			added, err := db.BatchSet(cacheKeyVal)
			if err != nil {
				return err
			}
			tl.UILogger.Println("Added", added, "embeddings to cache")
		}
		db2, err := localdb.NewFileDB("./.tzap-data/filesmd5.db")
		if err != nil {
			return err
		}
		var keyvals []types.KeyValue
		for file := range storedFiles {
			content, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			keyvals = append(keyvals, types.KeyValue{Key: file, Value: util.MD5HashByte(content)})
		}
		added, err := db2.BatchSet(keyvals)
		if err != nil {
			panic("failing to store changed files should not happend and has probably caused some kind of corruption")
		}
		tl.Logger.Printf("Added %d files to md5 cache", added)
	}
	return nil
}
