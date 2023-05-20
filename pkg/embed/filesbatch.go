package embed

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tzapio/tzap/pkg/embed/localdb"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func GetCachedEmbeddings(embeddings types.Embeddings) types.Embeddings {
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
					fmt.Println("cannot unmarshal", splitPart, err.Error())
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
					continue
				} else {
					fmt.Println("invalid vector length", splitPart, err.Error())
					continue
				}
			}
		}
		println("Warning: %s is uncached.", vector.ID)
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
			fmt.Println("Added", added, "embeddings to cache")
		}
	}
	return nil
}
func SaveEmbeddingToFile(e types.Embeddings) {
	if err := BatchEmbeddings(e); err != nil {
		panic(err)
	}
	if err := SaveVectorsToFile(e, "./.tzap-data/files.json"); err != nil {
		panic(err)
	}
}
func SaveVectorsToFile(e types.Embeddings, filePath string) error {
	embeddingJSON, err := json.Marshal(e)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, embeddingJSON, 0644); err != nil {
		return err
	}
	//fmt.Printf("Upserted files (count: %d): %s\n", len(e.Vectors), filePath)
	return nil
}

func BatchEmbeddings(e types.Embeddings) error {
	batchSize := 100
	batchNumber := 1
	var batch []types.Vector

	for i, vector := range e.Vectors {
		batch = append(batch, vector)
		deletePreviousBatch()
		if (i+1)%batchSize == 0 || i == len(e.Vectors)-1 {
			filePath := fmt.Sprintf("./.tzap-data/files-%d.json", batchNumber)
			batchEmbeddingJson := types.Embeddings{
				Vectors: batch,
			}
			err := SaveVectorsToFile(batchEmbeddingJson, filePath)
			if err != nil {
				return err
			}
			batch = nil
			batchNumber++
		}
	}
	return nil
}
func deletePreviousBatch() error {
	// Remove previous embedding files
	for i := 1; ; i++ {
		filePath := fmt.Sprintf("./.tzap-data/files-%d.json", i)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			break
		}
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}
	return nil
}
