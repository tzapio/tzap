package embed

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tzapio/tzap/pkg/embed/localdb"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func SplitCachedUncachedEmbeddings(embeddings types.Embeddings) ([]types.Vector, []types.Vector) {
	var cachedEmbeddings []types.Vector
	var uncachedEmbeddings []types.Vector
	db, err := localdb.NewFileDB("./.tzap-data/embeddingsCache.db")
	if err != nil {
		panic(err)
	}

	for _, vector := range embeddings.Vectors {
		splitPart := vector.Metadata["splitPart"]
		kv, exists := db.ScanGet(splitPart)
		if exists {
			float32Vector := []float32{}

			if kv.Value != "" {
				if err := json.Unmarshal([]byte(kv.Value), &float32Vector); err != nil {
					fmt.Println("cannot unmarshal", splitPart, err.Error())
					continue
				}
				if len(float32Vector) == 1536 {
					vector := types.Vector{
						ID:        vector.Metadata["filename"] + "-" + vector.Metadata["start"] + "-" + vector.Metadata["end"],
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
		uncachedEmbeddings = append(uncachedEmbeddings, vector)
	}
	return cachedEmbeddings, uncachedEmbeddings
}
func ProcessEmbeddings(t *tzap.Tzap, embeddings types.Embeddings) error {
	cachedEmbeddings, uncachedEmbeddings := SplitCachedUncachedEmbeddings(embeddings)
	var combinedVectors []types.Vector = cachedEmbeddings

	if len(uncachedEmbeddings) > 0 {

		var inputStrings []string
		for _, vector := range uncachedEmbeddings {
			inputStrings = append(inputStrings, vector.Metadata["splitPart"])
		}
		println("getting embeddings ", len(uncachedEmbeddings))
		embeddingsResult, err := t.TG.FetchEmbedding(t.C, inputStrings...)
		if err != nil {
			panic(err)
		}
		db, err := localdb.NewFileDB("./.tzap-data/embeddingsCache.db")
		if err != nil {
			panic(err)
		}
		var embeddingVectors []types.Vector
		cacheKeyVal := []types.KeyValue{}
		for i, embedding := range embeddingsResult {
			metadata := embeddings.Vectors[i].Metadata
			embeddingByte, err := json.Marshal(embedding)
			if err != nil {
				panic(err)
			}
			cacheKeyVal = append(cacheKeyVal, types.KeyValue{Key: inputStrings[i], Value: string(embeddingByte)})
			embeddingVectors = append(embeddingVectors, types.Vector{
				ID:        metadata["filename"] + "-" + metadata["start"] + "-" + metadata["end"],
				TimeStamp: 0,
				Metadata:  metadata,
				Values:    embedding,
			})
		}
		added, err := db.BatchSet(cacheKeyVal)
		if err != nil {
			panic(err)
		}
		fmt.Println("Added", added, "embeddings to cache out of ", len(inputStrings))
		combinedVectors = append(combinedVectors, embeddingVectors...)
	}

	embeddingJson := types.Embeddings{
		Vectors: combinedVectors,
	}

	if err := BatchEmbeddings(embeddingJson); err != nil {
		panic(err)
	}
	embeddingJSON, err := json.Marshal(embeddingJson)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("./.tzap-data/files.json", embeddingJSON, 0644); err != nil {
		panic(err)
	}
	fmt.Println("Embeddings saved to ./.tzap-data/files.json")
	return nil
}

func SaveBatchToFile(batch []types.Vector, batchNumber int) error {
	embeddingJson := types.Embeddings{
		Vectors: batch,
	}

	embeddingJSON, err := json.Marshal(embeddingJson)
	if err != nil {
		return err
	}

	if err := os.WriteFile(fmt.Sprintf("./.tzap-data/files-%d.json", batchNumber), embeddingJSON, 0644); err != nil {
		return err
	}
	fmt.Printf("Embeddings saved to ./.tzap-data/files-%d.json\n", batchNumber)
	return nil
}

func BatchEmbeddings(e types.Embeddings) error {
	batchSize := 100
	batchNumber := 1
	var batch []types.Vector

	for i, vector := range e.Vectors {
		batch = append(batch, vector)

		if (i+1)%batchSize == 0 || i == len(e.Vectors)-1 {
			err := SaveBatchToFile(batch, batchNumber)
			if err != nil {
				return err
			}
			batch = nil
			batchNumber++
		}
	}
	return nil
}
